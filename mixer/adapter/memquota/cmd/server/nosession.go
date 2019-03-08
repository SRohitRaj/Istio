// Copyright 2019 Istio Authors
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

// THIS FILE IS AUTOMATICALLY GENERATED.

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	memquota "istio.io/istio/mixer/adapter/memquota"
	config "istio.io/istio/mixer/adapter/memquota/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime/handler"
	"istio.io/istio/mixer/template/quota"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run()
	}

	// NoSession models nosession adapter backend.
	NoSession struct {
		listener net.Listener
		shutdown chan error
		server   *grpc.Server

		builder     adapter.HandlerBuilder
		env         adapter.Env
		builderLock sync.RWMutex
		handlerMap  map[string]adapter.Handler
	}

	// Cert includes cert config for adapter server
	Cert struct {
		credentialFile    string
		privateKeyFile    string
		caCertificateFile string
		enableTLS         bool
		requireClientAuth bool
	}
)

// DefaultCertOption return default cert option for adapter service
func DefaultCertOption() *Cert {
	return &Cert{
		credentialFile:    "/etc/certs/cert-chain.pem",
		privateKeyFile:    "/etc/certs/key.pem",
		caCertificateFile: "/etc/certs/root-cert.pem",
		enableTLS:         false,
		requireClientAuth: false,
	}
}

// AttachCobraFlags attach certs related flags.
func (c *Cert) AttachCobraFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&c.credentialFile, "certFile", "", c.credentialFile,
		"The location of the certificate file for TLS")
	cmd.PersistentFlags().StringVarP(&c.privateKeyFile, "keyFile", "", c.privateKeyFile,
		"The location of the key file for TLS")
	cmd.PersistentFlags().StringVarP(&c.caCertificateFile, "caCertFile", "", c.caCertificateFile,
		"The location of the certificate file for the root certificate authority")
	cmd.PersistentFlags().BoolVarP(&c.enableTLS, "enableTLS", "", c.enableTLS,
		"")
	cmd.PersistentFlags().BoolVarP(&c.requireClientAuth, "requireClientAuth", "", c.requireClientAuth,
		"")
}

var _ quota.HandleQuotaServiceServer = &NoSession{}

func (s *NoSession) updateHandlers(rawcfg []byte) (adapter.Handler, error) {
	cfg := &config.Params{}

	if err := cfg.Unmarshal(rawcfg); err != nil {
		return nil, err
	}

	s.builderLock.Lock()
	defer s.builderLock.Unlock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		return handler, nil
	}

	s.env.Logger().Infof("Loaded handler with: %v", cfg)
	s.builder.SetAdapterConfig(cfg)

	if ce := s.builder.Validate(); ce != nil {
		return nil, ce
	}

	h, err := s.builder.Build(context.Background(), s.env)
	if err != nil {
		s.env.Logger().Errorf("could not build: %v", err)
		return nil, err
	}
	s.handlerMap[string(rawcfg)] = h
	return h, nil
}

func (s *NoSession) getQuotaHandler(rawcfg []byte) (quota.Handler, error) {
	s.builderLock.RLock()
	if handler, ok := s.handlerMap[string(rawcfg)]; ok {
		h := handler.(quota.Handler)
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()
	h, err := s.updateHandlers(rawcfg)
	if err != nil {
		return nil, err
	}

	// establish session
	return h.(quota.Handler), nil
}

func quotaInstance(inst *quota.InstanceMsg) *quota.Instance {
	return &quota.Instance{
		Name: inst.Name,

		Dimensions: transformValueMap(inst.Dimensions),
	}
}

// nolint:deadcode
func transformValueMap(in map[string]*v1beta1.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = transformValue(v.GetValue())
	}
	return out
}

// nolint:deadcode
func transformValueSlice(in []interface{}) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, inst := range in {
		out = append(out, transformValue(inst))
	}
	return out
}

func transformValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *v1beta1.Value_StringValue:
		return t.StringValue
	case *v1beta1.Value_Int64Value:
		return t.Int64Value
	case *v1beta1.Value_DoubleValue:
		return t.DoubleValue
	case *v1beta1.Value_BoolValue:
		return t.BoolValue
	case *v1beta1.Value_IpAddressValue:
		return t.IpAddressValue.Value
	case *v1beta1.Value_EmailAddressValue:
		return t.EmailAddressValue.Value
	case *v1beta1.Value_UriValue:
		return t.UriValue.Value
	default:
		return fmt.Sprintf("%v", in)
	}
}

// HandleQuota handles 'Quota' instances.
func (s *NoSession) HandleQuota(ctx context.Context, r *quota.HandleQuotaRequest) (*adptModel.QuotaResult, error) {
	if r.AdapterConfig == nil {
		return nil, errors.New("adapter config cannot be empty")
	}
	h, err := s.getQuotaHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	qi := quotaInstance(r.Instance)
	resp := adptModel.QuotaResult{
		Quotas: make(map[string]adptModel.QuotaResult_Result),
	}
	for qt, p := range r.QuotaRequest.Quotas {
		qa := adapter.QuotaArgs{
			DeduplicationID: r.DedupId,
			QuotaAmount:     p.Amount,
			BestEffort:      p.BestEffort,
		}
		qr, err := h.HandleQuota(ctx, qi, qa)
		if err != nil {
			return nil, err
		}
		resp.Quotas[qt] = adptModel.QuotaResult_Result{
			ValidDuration: qr.ValidDuration,
			GrantedAmount: qr.Amount,
		}
	}
	if err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	}
	return &resp, nil
}

// Addr returns the listening address of the server
func (s *NoSession) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *NoSession) Run() {
	s.shutdown = make(chan error, 1)
	go func() { //nolint:adapterlinter
		err := s.server.Serve(s.listener)

		// notify closer we're done
		s.shutdown <- err
	}()
}

// Wait waits for server to stop
func (s *NoSession) Wait() error {
	if s.shutdown == nil {
		return fmt.Errorf("server not running")
	}

	err := <-s.shutdown
	s.shutdown = nil
	return err
}

// Close gracefully shuts down the server
func (s *NoSession) Close() error {
	if s.shutdown != nil {
		s.server.GracefulStop()
		_ = s.Wait()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

func getServerTLSOption(c *Cert) (grpc.ServerOption, error) {
	certificate, err := tls.LoadX509KeyPair(
		c.credentialFile,
		c.privateKeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load key cert pair")
	}
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(c.caCertificateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	if c.requireClientAuth {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return grpc.Creds(credentials.NewTLS(tlsConfig)), nil
}

// NewMemquotaNoSessionServer creates a new no session server based on given args.
func NewMemquotaNoSessionServer(addr uint16, poolSize int, c *Cert) (*NoSession, error) {
	saddr := fmt.Sprintf(":%d", addr)

	gp := pool.NewGoroutinePool(poolSize, false)
	inf := memquota.GetInfo()
	s := &NoSession{
		builder:    inf.NewBuilder(),
		env:        handler.NewEnv(0, "memquota-nosession", gp),
		handlerMap: make(map[string]adapter.Handler),
	}
	var err error
	if s.listener, err = net.Listen("tcp", saddr); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v\n", s.listener.Addr())
	if c.enableTLS {
		so, err := getServerTLSOption(c)
		if err != nil {
			return nil, err
		}
		s.server = grpc.NewServer(so)
	} else {
		s.server = grpc.NewServer()
	}

	quota.RegisterHandleQuotaServiceServer(s.server, s)

	return s, nil
}
