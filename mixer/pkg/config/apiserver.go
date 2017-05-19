// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rpc "github.com/googleapis/googleapis/google/rpc"

	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/config/descriptor"
	pb "istio.io/mixer/pkg/config/proto"
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/status"
)

type validateFunc func(cfg map[string]string) (rt *Validated, desc descriptor.Finder, ce *adapter.ConfigErrors)

type readBodyFunc func(r io.Reader) ([]byte, error)

// API defines and implements the configuration API.
// The server constructs and uses a validator for validations
// The server uses KeyValueStore to persist keys.
type API struct {
	version  string
	rootPath string

	// used at the back end for validation and storage
	store    KeyValueStore
	validate validateFunc

	// house keeping
	handler http.Handler
	server  *http.Server

	// fault injection
	readBody readBodyFunc
}

// MsgOk defines the text of the OK message in rpc.Status.Message.
const msgOk = "ok"

// APIResponse defines the shape of the api response.
type APIResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Status rpc.Status  `json:"status,omitempty"`
}

// register routes
func (a *API) register(c *restful.Container) {
	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON, "application/yaml", "application/x-yaml")
	ws.Produces(restful.MIME_JSON)
	ws.Path(a.rootPath)

	// List the scopes
	ws.Route(ws.
		GET("/scopes").
		To(a.getScopes).
		Doc("Gets scopes associated with Mixer").
		Writes(APIResponse{}))

	// Create a policy
	ws.Route(ws.
		POST("/scopes/{scope}/subjects/{subject}").
		To(a.createPolicy).
		Doc("Create a policy").
		// TODO Reads(...).
		Writes(APIResponse{}))

	// List the rules for a scope and subject
	ws.Route(ws.
		GET("/scopes/{scope}/subjects/{subject}/rules").
		To(a.getRules).
		Doc("Gets rules associated with the given scope and subject").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("subject", "subject").DataType("string")).
		Writes(APIResponse{}))

	// Delete the list of rules for a scope and subject
	ws.Route(ws.
		DELETE("/scopes/{scope}/subjects/{subject}/rules").
		To(a.deleteRules).
		Doc("Deletes rules associated with the given scope and subject").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("subject", "subject").DataType("string")).
		Writes(APIResponse{}))

	ws.Route(ws.
		PUT("/scopes/{scope}/subjects/{subject}/rules").
		To(a.putRules).
		Doc("Replaces rules associated with the given scope and subject").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("subject", "subject").DataType("string")).
		Reads(&pb.ServiceConfig{}).
		Writes(APIResponse{}))

	// Delete a rule
	ws.Route(ws.
		DELETE("/scopes/{scope}/subjects/{subject}/rules/{ruleid}").
		To(a.deleteRule).
		Doc("Replaces rules associated with the given scope and subject").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("subject", "subject").DataType("string")).
		Param(ws.PathParameter("ruleid", "rule id").DataType("string")))

	// Creates or replaces a rule's list of aspects
	ws.Route(ws.
		PUT("/scopes/{scope}/subjects/{subject}/rules/{ruleid}/aspects/{aspect}").
		To(a.putAspect).
		Doc("Creates or replaces a rule’s list of aspects.").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("subject", "subject").DataType("string")).
		Param(ws.PathParameter("ruleid", "rule id").DataType("string")).
		Param(ws.PathParameter("aspect", "aspect").DataType("string")).
		Reads(&pb.Aspect{}).
		Writes(APIResponse{}))

	// Creates or replaces a named adapter configuration
	ws.Route(ws.
		PUT("/scopes/{scope}/adapters/{adapter_name}/{config_name}").
		To(a.putAdapter).
		Doc("Creates or replaces a named adapter configuration.").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("adapter_name", "adapter name").DataType("string")).
		Param(ws.PathParameter("config_name", "config name").DataType("string")).
		Reads(&pb.Adapter{}).
		Writes(APIResponse{}))

	// Gets a descriptor
	ws.Route(ws.
		GET("/scopes/{scope}/descriptors/{descriptor_type}/{descriptor_name}").
		To(a.getDescriptor).
		Doc("Gets a descriptor.").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("descriptor_type", "descriptor type").DataType("string")).
		Param(ws.PathParameter("descriptor_name", "descriptor name").DataType("string")).
		Writes(APIResponse{}))

	// Creates or replaces a descriptor
	ws.Route(ws.
		PUT("/scopes/{scope}/descriptors/{descriptor_type}/{descriptor_name}").
		To(a.putDescriptor).
		Doc("Creates or replaces a descriptor.").
		Param(ws.PathParameter("scope", "scope").DataType("string")).
		Param(ws.PathParameter("descriptor_type", "descriptor type").DataType("string")).
		Param(ws.PathParameter("descriptor_name", "descriptor name").DataType("string")).
		// TODO Reads(&pb.Descriptor{}).
		Writes(APIResponse{}))

	c.Add(ws)
}

// NewAPI creates a new API server
func NewAPI(version string, port uint16, tc expr.TypeChecker, aspectFinder AspectValidatorFinder,
	builderFinder BuilderValidatorFinder, findAspects AdapterToAspectMapper, store KeyValueStore) *API {
	c := restful.NewContainer()
	a := &API{
		version:  version,
		rootPath: fmt.Sprintf("/api/%s", version),
		store:    store,
		readBody: ioutil.ReadAll,
		validate: func(cfg map[string]string) (*Validated, descriptor.Finder, *adapter.ConfigErrors) {
			v := newValidator(aspectFinder, builderFinder, findAspects, true, tc)
			rt, ce := v.validate(cfg)
			return rt, v.descriptorFinder, ce
		},
	}
	a.register(c)
	a.server = &http.Server{Addr: ":" + strconv.Itoa(int(port)), Handler: c}
	a.handler = c
	// ensure that we always send back an APIResponse object.
	c.ServiceErrorHandler(func(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
		writeErrorResponse(err.Code, err.Message, resp)
	})
	return a
}

// Run calls listen and serve on the API server
func (a *API) Run() {
	glog.Warning(a.server.ListenAndServe())
}

// getScopes returns the scopes
// "/scopes"
func (a *API) getScopes(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "Listing scopes not implemented", resp)
}

// getRules returns the rules document for the scope and the subject.
// "/scopes/{scope}/subjects/{subject}/rules"
func (a *API) getRules(req *restful.Request, resp *restful.Response) {
	funcPath := req.Request.URL.Path[len(a.rootPath):]
	st, msg, data := getRules(a.store, funcPath)
	writeResponse(st, msg, data, resp)
}

// deleteRules deletes the rules document for the scope and the subject.
// "/scopes/{scope}/subjects/{subject}/rules"
func (a *API) deleteRules(req *restful.Request, resp *restful.Response) {
	funcPath := req.Request.URL.Path[len(a.rootPath):]
	if err := a.store.Delete(funcPath); err != nil {
		// This should only happen if user asks to delete
		// rules that the store *could not* delete
		writeErrorResponse(http.StatusInternalServerError, err.Error(), resp)
		return
	}

	writeResponse(http.StatusOK, fmt.Sprintf("Deleted %s", funcPath), nil, resp)
}

func getRules(store KeyValueStore, path string) (statusCode int, msg string, data *pb.ServiceConfig) {
	var val string
	var found bool

	if val, _, found = store.Get(path); !found {
		return http.StatusNotFound, fmt.Sprintf("no rules for %s", path), nil
	}

	m := &pb.ServiceConfig{}
	if err := yaml.Unmarshal([]byte(val), m); err != nil {
		msg := fmt.Sprintf("unable to parse rules at '%s': %v", path, err)
		glog.Warning(msg)
		return http.StatusInternalServerError, msg, nil
	}
	return http.StatusOK, msgOk, m
}

// putRules replaces the entire rules document for the scope and subject
// "/scopes/{scope}/subjects/{subject}/rules"
func (a *API) putRules(req *restful.Request, resp *restful.Response) {

	key := req.Request.URL.Path[len(a.rootPath):]
	var data map[string]string
	var err error
	// TODO optimize only read descriptors and adapters
	if data, _, _, err = readdb(a.store, "/"); err != nil {
		writeErrorResponse(http.StatusInternalServerError, err.Error(), resp)
		return
	}
	// TODO send index back to the client

	var bval []byte
	if bval, err = a.readBody(req.Request.Body); err != nil {
		writeErrorResponse(http.StatusInternalServerError, err.Error(), resp)
		return
	}
	val := string(bval)
	data[key] = val
	/*
		rt *Validated, desc descriptor.Finder, ce *adapter.ConfigErrors
	*/
	var vd *Validated
	var cerr *adapter.ConfigErrors
	if vd, _, cerr = a.validate(data); cerr != nil {
		glog.Warningf("Validation failed with %s\n %s", cerr.Error(), val)
		writeErrorResponse(http.StatusBadRequest, cerr.Error(), resp)
		return
	}

	if _, err = a.store.Set(key, val); err != nil {
		writeErrorResponse(http.StatusInternalServerError, err.Error(), resp)
		return
	}
	// TODO send index back to the client
	writeResponse(http.StatusOK, fmt.Sprintf("Created %s", key),
		vd.rule[*parseRulesKey(key)], resp)
}

// createPolicy creates a policy
// "/scopes/{scope}/subjects/{subject}"
func (a *API) createPolicy(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "create policy not implemented", resp) // TODO
}

// putAspect creates or replaces a rule’s list of aspects.
// "/scopes/{scope}/subjects/{subject}/rules/{ruleid}/aspects/{aspect}"
func (a *API) putAspect(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "put aspect not implemented", resp) // TODO
}

// deleteRule deletes a rule
// "/scopes/{scope}/subjects/{subject}/rules/{ruleid}"
func (a *API) deleteRule(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "delete rule not implemented", resp) // TODO
}

// putAdapter creates or replaces an adapter configuration
// "/scopes/{scope}/adapters/{adapter_name}/{config_name}"
func (a *API) putAdapter(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "put adapter not implemented", resp) // TODO
}

// getDescriptor returns a descriptor
// "/scopes/{scope}/descriptors/{descriptor_type}/{descriptor_name}"
func (a *API) getDescriptor(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "get descriptor not implemented", resp) // TODO
}

// putDescriptor creates or replaces a descriptor
// "/scopes/{scope}/descriptors/{descriptor_type}/{descriptor_name}"
func (a *API) putDescriptor(req *restful.Request, resp *restful.Response) {
	writeErrorResponse(http.StatusNotImplemented, "put descriptor not implemented", resp) // TODO
}

// a subset of restful.Response
type response interface {
	// WriteHeaderAndJson is a convenience method for writing the status and a value in Json with a given Content-Type.
	WriteHeaderAndJson(status int, value interface{}, contentType string) error
}

func writeResponse(httpStatus int, msg string, data interface{}, resp response) {
	if err := resp.WriteHeaderAndJson(
		httpStatus,
		&APIResponse{
			Data: data,
			Status: status.WithMessage(
				httpStatusToRPC(httpStatus), msg),
		},
		restful.MIME_JSON,
	); err != nil {
		glog.Warning(err)
	}
}

func writeErrorResponse(httpStatus int, msg string, resp response) {
	writeResponse(httpStatus, msg, nil, resp)
}

func httpStatusToRPC(httpStatus int) (code rpc.Code) {
	var ok bool
	if code, ok = httpStatusToRPCMap[httpStatus]; !ok {
		code = rpc.UNKNOWN
	}
	return code
}

// httpStatusToRpc limited mapping from proto documentation.
var httpStatusToRPCMap = map[int]rpc.Code{
	http.StatusOK:                 rpc.OK,
	http.StatusNotFound:           rpc.NOT_FOUND,
	http.StatusConflict:           rpc.ALREADY_EXISTS,
	http.StatusForbidden:          rpc.PERMISSION_DENIED,
	http.StatusUnauthorized:       rpc.UNAUTHENTICATED,
	http.StatusPreconditionFailed: rpc.FAILED_PRECONDITION,
	http.StatusBadRequest:         rpc.INVALID_ARGUMENT,
}
