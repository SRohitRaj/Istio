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

package main

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/api/routing/v1alpha1"
	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/test/util"
)

// sortedConfigStore lets us facade any ConfigStore (such as memory.Make()'s) providing
// a stable List() which helps with testing `istioctl get` output.
type sortedConfigStore struct {
	store model.ConfigStore
}

type testCase struct {
	configs []model.Config
	args    []string

	// Typically use one of the three
	expectedOutput string         // Expected constant output
	expectedRegexp *regexp.Regexp // Expected regexp output
	goldenFilename string         // Expected output stored in golden file

	wantException bool
}

type testCaseGoldenFile struct {
	configs []model.Config
	args    []string

	filename       string         // Expected constant output
	expectedRegexp *regexp.Regexp // Expected regexp output

	wantException bool
}

var (
	testRouteRules = []model.Config{
		{
			ConfigMeta: model.ConfigMeta{
				Name:      "d",
				Namespace: "default",
				Type:      model.RouteRule.Type,
				Group:     model.RouteRule.Group,
				Version:   model.RouteRule.Version,
			},
			Spec: &v1alpha1.RouteRule{
				Precedence: 2,
				Destination: &v1alpha1.IstioService{
					Name: "d",
				},
			},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "b",
				Namespace: "default",
				Type:      model.RouteRule.Type,
				Group:     model.RouteRule.Group,
				Version:   model.RouteRule.Version,
			},
			Spec: &v1alpha1.RouteRule{
				Precedence: 3,
				Destination: &v1alpha1.IstioService{
					Name: "b",
				},
			},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "c",
				Namespace: "istio-system",
				Type:      model.RouteRule.Type,
				Group:     model.RouteRule.Group,
				Version:   model.RouteRule.Version,
			},
			Spec: &v1alpha1.RouteRule{
				Precedence: 2,
				Destination: &v1alpha1.IstioService{
					Name: "c",
				},
			},
		},
		{
			ConfigMeta: model.ConfigMeta{Name: "a",
				Namespace: "default",
				Type:      model.RouteRule.Type,
				Group:     model.RouteRule.Group,
				Version:   model.RouteRule.Version,
			},
			Spec: &v1alpha1.RouteRule{
				Precedence: 1,
				Destination: &v1alpha1.IstioService{
					Name: "a",
				},
			},
		},
	}

	testGateways = []model.Config{
		{
			ConfigMeta: model.ConfigMeta{
				Name:      "bookinfo-gateway",
				Namespace: "default",
				Type:      model.Gateway.Type,
				Group:     model.Gateway.Group,
				Version:   model.Gateway.Version,
			},
			Spec: &networking.Gateway{
				Selector: map[string]string{"istio": "ingressgateway"},
				Servers: []*networking.Server{
					{
						Port: &networking.Port{
							Number:   80,
							Protocol: "HTTP",
						},
						Hosts: []string{"*"},
					},
				},
			},
		},
	}

	testVirtualServices = []model.Config{
		{
			ConfigMeta: model.ConfigMeta{
				Name:      "bookinfo",
				Namespace: "default",
				Type:      model.VirtualService.Type,
				Group:     model.VirtualService.Group,
				Version:   model.VirtualService.Version,
			},
			Spec: &networking.VirtualService{
				Hosts:    []string{"*"},
				Gateways: []string{"bookinfo-gateway"},
				Http: []*networking.HTTPRoute{
					{
						Match: []*networking.HTTPMatchRequest{
							{
								Uri: &networking.StringMatch{
									&networking.StringMatch_Exact{"/productpage"},
								},
							},
							{
								Uri: &networking.StringMatch{
									&networking.StringMatch_Exact{"/login"},
								},
							},
							{
								Uri: &networking.StringMatch{
									&networking.StringMatch_Exact{"/logout"},
								},
							},
							{
								Uri: &networking.StringMatch{
									&networking.StringMatch_Prefix{"/api/v1/products"},
								},
							},
						},
						Route: []*networking.DestinationWeight{
							{
								Destination: &networking.Destination{
									Host: "productpage",
									Port: &networking.PortSelector{
										Port: &networking.PortSelector_Number{80},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	testDestinationRules = []model.Config{
		{
			ConfigMeta: model.ConfigMeta{
				Name:      "googleapis",
				Namespace: "default",
				Type:      model.DestinationRule.Type,
				Group:     model.DestinationRule.Group,
				Version:   model.DestinationRule.Version,
			},
			Spec: &networking.DestinationRule{
				Host: "*.googleapis.com",
				TrafficPolicy: &networking.TrafficPolicy{
					Tls: &networking.TLSSettings{
						Mode: networking.TLSSettings_SIMPLE,
					},
				},
			},
		},
	}

	testServiceEntries = []model.Config{
		{
			ConfigMeta: model.ConfigMeta{
				Name:      "googleapis",
				Namespace: "default",
				Type:      model.ServiceEntry.Type,
				Group:     model.ServiceEntry.Group,
				Version:   model.ServiceEntry.Version,
			},
			Spec: &networking.ServiceEntry{
				Hosts: []string{"*.googleapis.com"},
				Ports: []*networking.Port{
					{
						Name:     "https",
						Number:   443,
						Protocol: "HTTP",
					},
				},
			},
		},
	}
)

func TestGet(t *testing.T) {
	cases := []testCase{
		{ // case 0
			[]model.Config{},
			strings.Split("get routerules", " "),
			`No resources found.
`,
			nil,
			"",
			false,
		},
		{ // case 1
			testRouteRules,
			strings.Split("get routerules --all-namespaces", " "),
			`NAME      KIND                        NAMESPACE
a         RouteRule.config.v1alpha2   default
b         RouteRule.config.v1alpha2   default
c         RouteRule.config.v1alpha2   istio-system
d         RouteRule.config.v1alpha2   default
`,
			nil,
			"",
			false,
		},
		{ // case 2
			testGateways,
			strings.Split("get gateways -n default", " "),
			`GATEWAY NAME       HOSTS     NAMESPACE
bookinfo-gateway   *         default
`,
			nil,
			"",
			false,
		},
		{ // case 3
			testVirtualServices,
			strings.Split("get virtualservices -n default", " "),
			`VIRTUAL-SERVICE NAME   GATEWAYS           HOSTS     #HTTP     #TCP      NAMESPACE
bookinfo               bookinfo-gateway   *             1        0      default
`,
			nil,
			"",
			false,
		},
		{ // case 4 invalid type
			[]model.Config{},
			strings.Split("get invalid", " "),
			"",
			regexp.MustCompile("^Usage:.*"),
			"",
			true, // "istioctl get invalid" should fail
		},
		{ // case 5 all
			append(testRouteRules, testVirtualServices...),
			strings.Split("get all", " "),
			`VIRTUAL-SERVICE NAME   GATEWAYS           HOSTS     #HTTP     #TCP      NAMESPACE
bookinfo               bookinfo-gateway   *             1        0      default

NAME      KIND                        NAMESPACE
a         RouteRule.config.v1alpha2   default
b         RouteRule.config.v1alpha2   default
c         RouteRule.config.v1alpha2   istio-system
d         RouteRule.config.v1alpha2   default
`,
			nil,
			"",
			false,
		},
		{ // case 6 all with no data
			[]model.Config{},
			strings.Split("get all", " "),
			"No resources found.\n",
			nil,
			"",
			false,
		},
		{ // case 7
			testDestinationRules,
			strings.Split("get destinationrules", " "),
			`DESTINATION-RULE NAME   HOST               SUBSETS   NAMESPACE
googleapis              *.googleapis.com             default
`,
			nil,
			"",
			false,
		},
		{ // case 8
			testServiceEntries,
			strings.Split("get serviceentries", " "),
			`SERVICE-ENTRY NAME   HOSTS              PORTS      NAMESPACE
googleapis           *.googleapis.com   HTTP/443   default
`,
			nil,
			"",
			false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}

func TestCreate(t *testing.T) {
	cases := []testCase{
		{ // case 0 -- invalid doesn't provide -f filename
			[]model.Config{},
			strings.Split("create routerules", " "),
			"",
			regexp.MustCompile("^Usage:.*"),
			"",
			true,
		},
		{ // case 1
			[]model.Config{},
			strings.Split("create -f convert/testdata/v1alpha1/route-rule-80-20.yaml", " "),
			"",
			regexp.MustCompile("^Warning: route-rule is deprecated and will not be supported" +
				" in future Istio versions \\(route-rule-80-20\\).\n" +
				"Created config route-rule/default/route-rule-80-20.*"),
			"",
			false,
		},
		{ // case 2
			[]model.Config{},
			strings.Split("create -f convert/testdata/v1alpha3/route-rule-80-20.yaml.golden", " "),
			"",
			regexp.MustCompile("^Created config virtual-service/default/helloworld-service.*"),
			"",
			false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}

func TestReplace(t *testing.T) {
	cases := []testCase{
		{ // case 0 -- invalid doesn't provide -f
			[]model.Config{},
			strings.Split("replace routerules", " "),
			"",
			regexp.MustCompile("^Usage:.*"),
			"",
			true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []testCase{
		{ // case 0
			[]model.Config{},
			strings.Split("delete routerule unknown", " "),
			"",
			regexp.MustCompile("^Error: 1 error occurred:\n\n\\* cannot delete unknown: item not found\n$"),
			"",
			true,
		},
		{ // case 1
			testRouteRules,
			strings.Split("delete routerule a", " "),
			`Deleted config: routerule a
`,
			nil,
			"",
			false,
		},
		{ // case 2
			testRouteRules,
			strings.Split("delete routerule a b", " "),
			`Deleted config: routerule a
Deleted config: routerule b
`,
			nil,
			"",
			false,
		},
		{ // case 3 - delete by filename of istio config which doesn't exist
			testRouteRules,
			strings.Split("delete -f convert/testdata/v1alpha1/route-rule-80-20.yaml", " "),
			"",
			regexp.MustCompile("^Error: 1 error occurred:\n\n\\* cannot delete route-rule/default/route-rule-80-20: item not found\n$"),
			"",
			true,
		},
		{ // case 4 - "all" not valid for delete
			[]model.Config{},
			strings.Split("delete all foo", " "),
			"",
			regexp.MustCompile("^Error: configuration type all not found"),
			"",
			true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}

// mockClientFactoryGenerator creates a factory for model.ConfigStore preloaded with data
func mockClientFactoryGenerator(configs []model.Config) func() (model.ConfigStore, error) {
	outFactory := func() (model.ConfigStore, error) {
		// Initialize the real client to get the supported config types
		realClient, err := newClient()
		if err != nil {
			return nil, err
		}

		// Initialize memory based model.ConfigStore with configs
		outConfig := memory.Make(realClient.ConfigDescriptor())
		for _, config := range configs {
			if _, err := outConfig.Create(config); err != nil {
				return nil, err
			}
		}

		// Wrap the memory ConfigStore so List() is sorted
		return sortedConfigStore{store: outConfig}, nil
	}

	return outFactory
}

func (cs sortedConfigStore) Create(config model.Config) (string, error) {
	return cs.store.Create(config)
}

func (cs sortedConfigStore) Get(typ, name, namespace string) (*model.Config, bool) {
	return cs.store.Get(typ, name, namespace)
}

func (cs sortedConfigStore) Update(config model.Config) (string, error) {
	return cs.store.Update(config)
}
func (cs sortedConfigStore) Delete(typ, name, namespace string) error {
	return cs.store.Delete(typ, name, namespace)
}

func (cs sortedConfigStore) ConfigDescriptor() model.ConfigDescriptor {
	return cs.store.ConfigDescriptor()
}

// List() is a facade that always returns cs.store items sorted by name/namespace
func (cs sortedConfigStore) List(typ, namespace string) ([]model.Config, error) {
	out, err := cs.store.List(typ, namespace)
	if err != nil {
		return out, err
	}

	// Sort by name, namespace
	sort.Slice(out, func(i, j int) bool {
		iName := out[i].ConfigMeta.Name
		jName := out[j].ConfigMeta.Name
		if iName == jName {
			return out[i].ConfigMeta.Namespace < out[j].ConfigMeta.Namespace
		}
		return iName < jName
	})

	return out, nil
}

func verifyOutput(t *testing.T, c testCase) {
	t.Helper()

	// Override the client factory used by main.go
	clientFactory = mockClientFactoryGenerator(c.configs)

	var out bytes.Buffer
	rootCmd.SetOutput(&out)
	rootCmd.SetArgs(c.args)

	file = "" // Clear, because we re-use

	fErr := rootCmd.Execute()
	output := out.String()

	if c.expectedOutput != "" && c.expectedOutput != output {
		t.Fatalf("Unexpected output for 'istioctl %s'\n got: %q\nwant: %q",
			strings.Join(c.args, " "), output, c.expectedOutput)
	}

	if c.expectedRegexp != nil && !c.expectedRegexp.MatchString(output) {
		t.Fatalf("Output didn't match for 'istioctl %s'\n got %v\nwant: %v",
			strings.Join(c.args, " "), output, c.expectedRegexp)
	}

	if c.goldenFilename != "" {
		util.CompareContent([]byte(output), c.goldenFilename, t)
	}

	if c.wantException {
		if fErr == nil {
			t.Fatalf("Wanted an exception for 'istioctl %s', didn't get one, output was %q",
				strings.Join(c.args, " "), output)
		}
	} else {
		if fErr != nil {
			t.Fatalf("Unwanted exception for 'istioctl %s': %v", strings.Join(c.args, " "), fErr)
		}
	}
}

func TestKubeInject(t *testing.T) {
	cases := []testCase{
		{ // case 0
			[]model.Config{},
			strings.Split("kube-inject", " "),
			"Error: filename not specified (see --filename or -f)\n",
			nil,
			"",
			true,
		},
		{ // case 1
			[]model.Config{},
			strings.Split("kube-inject -f missing.yaml", " "),
			"Error: open missing.yaml: no such file or directory\n",
			nil,
			"",
			true,
		},
		{ // case 2
			[]model.Config{},
			strings.Split("kube-inject --meshConfigFile testdata/mesh-config.yaml --injectConfigFile testdata/inject-config.yaml -f testdata/deployment/hello.yaml", " "),
			"",
			nil,
			"testdata/deployment/hello.yaml.injected",
			false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}
