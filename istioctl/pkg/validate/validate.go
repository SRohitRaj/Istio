// Copyright 2018 Istio Authors.
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

package validate

import (
	"errors"
	"fmt"
	"io"
	"os"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	mixercrd "istio.io/istio/mixer/pkg/config/crd"
	mixerstore "istio.io/istio/mixer/pkg/config/store"
	mixervalidate "istio.io/istio/mixer/pkg/validate"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pilot/pkg/model"
)

var (
	mixerAPIVersion = "config.istio.io/v1alpha2"

	errMissingFilename = errors.New(`error: you must specify resources by --filename.
Example resource specifications include:
   '-f rsrc.yaml'
   '--filename=rsrc.json'`)
	errKindNotSupported = errors.New("kind is not supported")
)

type validator struct {
	mixerValidator mixerstore.BackendValidator
}

func (v *validator) validateResource(un *unstructured.Unstructured) error {
	schema, exists := model.IstioConfigTypes.GetByType(crd.CamelCaseToKebabCase(un.GetKind()))
	if exists {
		obj, err := crd.ConvertObjectFromUnstructured(schema, un, "")
		if err != nil {
			return fmt.Errorf("cannot parse proto message: %v", err)
		}
		return schema.Validate(obj.Name, obj.Namespace, obj.Spec)
	}

	if v.mixerValidator != nil && un.GetAPIVersion() == mixerAPIVersion {
		if !v.mixerValidator.SupportsKind(un.GetKind()) {
			return errKindNotSupported
		}
		return v.mixerValidator.Validate(&mixerstore.BackendEvent{
			Type: mixerstore.Update,
			Key: mixerstore.Key{
				Name:      un.GetName(),
				Namespace: un.GetNamespace(),
				Kind:      un.GetKind(),
			},
			Value: mixercrd.ToBackEndResource(un),
		})
	}

	return nil
}

func (v *validator) validateFile(reader io.Reader) error {
	decoder := yaml.NewDecoder(reader)
	var errs error
	for {
		out := make(map[string]interface{})
		err := decoder.Decode(&out)
		if err == io.EOF {
			return errs
		}
		if err != nil {
			errs = multierror.Append(errs, err)
			return errs
		}
		if len(out) == 0 {
			continue
		}
		un := unstructured.Unstructured{Object: out}
		err = v.validateResource(&un)
		if err != nil {
			errs = multierror.Append(errs, multierror.Prefix(err, fmt.Sprintf("%s/%s/%s:",
				un.GetKind(), un.GetNamespace(), un.GetName())))
		}
	}
}

func validateFiles(filenames []string, referential bool) error {
	if len(filenames) == 0 {
		return errMissingFilename
	}

	v := &validator{
		mixerValidator: mixervalidate.NewDefaultValidator(referential),
	}

	var errs error
	for _, filename := range filenames {
		reader, err := os.Open(filename)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("cannot read file %q: %v", filename, err))
			continue
		}
		err = v.validateFile(reader)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// NewValidateCommand creates a new command for validating Istio k8s resources.
func NewValidateCommand() *cobra.Command {
	var filenames []string
	var referential bool

	c := &cobra.Command{
		Use:     "validate -f FILENAME [options]",
		Short:   "Validate Istio policy and rules",
		Example: `istioctl validate -f bookinfo-gateway.yaml`,
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return validateFiles(filenames, referential)
		},
	}

	flags := c.PersistentFlags()
	flags.StringSliceVarP(&filenames, "filename", "f", nil, "Names of files to validate")
	flags.BoolVarP(&referential, "referential", "x", true, "Enable structural validation for policy and telemetry")

	return c
}
