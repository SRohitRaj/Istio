// Copyright 2016 Istio Authors
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

package template

// InterfaceTemplate defines the template used to generate the adapter
// interfaces for Mixer for a given aspect.
var InterfaceTemplate = `// Copyright 2017 Istio Authors
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

package template

import (
	"github.com/golang/protobuf/proto"
	"fmt"
	"istio.io/mixer/pkg/attribute"
	rpc "github.com/googleapis/googleapis/google/rpc"
	"github.com/hashicorp/go-multierror"
	"istio.io/mixer/pkg/expr"
	"github.com/golang/glog"
	"istio.io/mixer/pkg/status"
	"istio.io/mixer/pkg/adapter"
	"istio.io/api/mixer/v1/config/descriptor"
	adptTmpl "istio.io/mixer/pkg/adapter/template"
	{{range .}}
		"{{.PackageImportPath}}"
	{{end}}
)

var (
	SupportedTmplInfo = map[string]Info {
	{{range .}}
		{{.GoPackageName}}.TemplateName: {
			CtrCfg:  &{{.GoPackageName}}.InstanceParam{},
			Variety:   adptTmpl.{{.VarietyName}},
			BldrName:  "{{.PackageImportPath}}.{{.Name}}ProcessorBuilder",
			HndlrName: "{{.PackageImportPath}}.{{.Name}}Processor",
			SupportsTemplate: func(hndlrBuilder adapter.HandlerBuilder) bool {
				_, ok := hndlrBuilder.({{.GoPackageName}}.{{.Name}}ProcessorBuilder)
				return ok
			},
			HandlerSupportsTemplate: func(hndlr adapter.Handler) bool {
				_, ok := hndlr.({{.GoPackageName}}.{{.Name}}Processor)
				return ok
			},
			InferType: func(cp proto.Message, tEvalFn TypeEvalFn) (proto.Message, error) {
				var err error = nil
				cpb := cp.(*{{.GoPackageName}}.InstanceParam)
				infrdType := &{{.GoPackageName}}.Type{}

				{{range .TemplateMessage.Fields}}
					{{if isPrimitiveValueType .GoType}}
						infrdType.{{.GoName}} = {{primitiveToValueType .GoType}}
					{{end}}
					{{if isValueType .GoType}}
						if infrdType.{{.GoName}}, err = tEvalFn(cpb.{{.GoName}}); err != nil {
							return nil, err
						}
					{{end}}
					{{if isStringValueTypeMap .GoType}}
						infrdType.{{.GoName}} = make(map[string]istio_mixer_v1_config_descriptor.ValueType)
						for k, v := range cpb.{{.GoName}} {
							if infrdType.{{.GoName}}[k], err = tEvalFn(v); err != nil {
								return nil, err
							}
						}
					{{end}}
				{{end}}
				_ = cpb
				return infrdType, err
			},
			ConfigureType: func(types map[string]proto.Message, builder *adapter.HandlerBuilder) error {
				// Mixer framework should have ensured the type safety.
				castedBuilder := (*builder).({{.GoPackageName}}.{{.Name}}ProcessorBuilder)
				castedTypes := make(map[string]*{{.GoPackageName}}.Type)
				for k, v := range types {
					// Mixer framework should have ensured the type safety.
					v1 := v.(*{{.GoPackageName}}.Type)
					castedTypes[k] = v1
				}
				return castedBuilder.Configure{{.Name}}(castedTypes)
			},
			{{if eq .VarietyName "TEMPLATE_VARIETY_REPORT"}}
				ProcessReport: func(insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler) rpc.Status {
					result := &multierror.Error{}
					var instances []*{{.GoPackageName}}.Instance

					castedInsts := make(map[string]*{{.GoPackageName}}.InstanceParam)
					for k, v := range insts {
						v1 := v.(*{{.GoPackageName}}.InstanceParam)
						castedInsts[k] = v1
					}
					for name, md := range castedInsts {
						{{range .TemplateMessage.Fields}}
							{{if isStringValueTypeMap .GoType}}
								{{.GoName}}, err := evalAll(md.{{.GoName}}, attrs, mapper)
							{{else}}
								{{.GoName}}, err := mapper.Eval(md.{{.GoName}}, attrs)
							{{end}}
								if err != nil {
									result = multierror.Append(result, fmt.Errorf("failed to eval {{.GoName}} for instance '%s': %v", name, err))
									continue
								}
						{{end}}

						instances = append(instances, &{{.GoPackageName}}.Instance{
							Name:       name,
							{{range .TemplateMessage.Fields}}
								{{if isPrimitiveValueType .GoType}}
									{{.GoName}}: {{.GoName}}.({{.GoType}}),
								{{else}}
									{{.GoName}}: {{.GoName}},
								{{end}}
							{{end}}
						})
					}

					if err := handler.({{.GoPackageName}}.{{.Name}}Processor).Report{{.Name}}(instances); err != nil {
						result = multierror.Append(result, fmt.Errorf("failed to report all values: %v", err))
					}

					err := result.ErrorOrNil()
					if err != nil {
						return status.WithError(err)
					}

					return status.OK
				},
				ProcessCheck: nil,
				ProcessQuota: nil,
			{{else if eq .VarietyName "TEMPLATE_VARIETY_CHECK"}}
				ProcessCheck: func(insts map[string]proto.Message, attrs attribute.Bag, mapper expr.Evaluator,
				handler adapter.Handler) (rpc.Status, adapter.CacheabilityInfo) {
					var found bool
					var err error

					var instances []*{{.GoPackageName}}.Instance
					castedInsts := make(map[string]*{{.GoPackageName}}.InstanceParam)
					for k, v := range insts {
						v1 := v.(*{{.GoPackageName}}.InstanceParam)
						castedInsts[k] = v1
					}
					for name, md := range castedInsts {
						{{range .TemplateMessage.Fields}}
							{{if isStringValueTypeMap .GoType}}
								{{.GoName}}, err := evalAll(md.{{.GoName}}, attrs, mapper)
							{{else}}
								{{.GoName}}, err := mapper.Eval(md.{{.GoName}}, attrs)
							{{end}}
								if err != nil {
									return status.WithError(err), adapter.CacheabilityInfo{}
								}
						{{end}}

						instances = append(instances, &{{.GoPackageName}}.Instance{
							Name:       name,
							{{range .TemplateMessage.Fields}}
								{{if isPrimitiveValueType .GoType}}
									{{.GoName}}: {{.GoName}}.({{.GoType}}),
								{{else}}
									{{.GoName}}: {{.GoName}},
								{{end}}
							{{end}}
						})
					}
					var cacheInfo adapter.CacheabilityInfo
					if found, cacheInfo, err = handler.({{.GoPackageName}}.{{.Name}}Processor).Check{{.Name}}(instances); err != nil {
						return status.WithError(err), adapter.CacheabilityInfo{}
					}

					if found {
						return status.OK, cacheInfo
					}

					return status.WithPermissionDenied(fmt.Sprintf("%s rejected", instances)), adapter.CacheabilityInfo{}
				},
				ProcessReport: nil,
				ProcessQuota: nil,
			{{else}}
				ProcessQuota: func(quotaName string, inst proto.Message, attrs attribute.Bag, mapper expr.Evaluator, handler adapter.Handler,
				qma adapter.QuotaRequestArgs) (rpc.Status, adapter.CacheabilityInfo, adapter.QuotaResult) {
					castedInst := inst.(*{{.GoPackageName}}.InstanceParam)
					{{range .TemplateMessage.Fields}}
						{{if isStringValueTypeMap .GoType}}
							{{.GoName}}, err := evalAll(castedInst.{{.GoName}}, attrs, mapper)
						{{else}}
							{{.GoName}}, err := mapper.Eval(castedInst.{{.GoName}}, attrs)
						{{end}}
							if err != nil {
								msg := fmt.Sprintf("failed to eval {{.GoName}} for instance '%s': %v", quotaName, err)
								glog.Error(msg)
								return status.WithInvalidArgument(msg), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
							}
					{{end}}

					instance := &{{.GoPackageName}}.Instance{
						Name:       quotaName,
						{{range .TemplateMessage.Fields}}
							{{if isPrimitiveValueType .GoType}}
								{{.GoName}}: {{.GoName}}.({{.GoType}}),
							{{else}}
								{{.GoName}}: {{.GoName}},
							{{end}}
						{{end}}
					}

					var qr adapter.QuotaResult
					var cacheInfo adapter.CacheabilityInfo
					if qr, cacheInfo, err = handler.({{.GoPackageName}}.{{.Name}}Processor).Alloc{{.Name}}(instance, qma); err != nil {
						glog.Errorf("Quota allocation failed: %v", err)
						return status.WithError(err), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
					}
					if qr.Amount == 0 {
						msg := fmt.Sprintf("Unable to allocate %v units from quota %s", qma.QuotaAmount, quotaName)
						glog.Warning(msg)
						return status.WithResourceExhausted(msg), adapter.CacheabilityInfo{}, adapter.QuotaResult{}
					}
					if glog.V(2) {
						glog.Infof("Allocated %v units from quota %s", qma.QuotaAmount, quotaName)
					}
					return status.OK, cacheInfo, qr
				},
				ProcessReport: nil,
				ProcessCheck: nil,
			{{end}}
		},
	{{end}}
	}
)
{{range .}}
{{end}}

`
