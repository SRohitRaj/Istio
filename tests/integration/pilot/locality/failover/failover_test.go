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

package failover

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"text/template"

	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/apps"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/galley"
	"istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/tests/integration/pilot/locality"
)

// This test allows for Locality Load Balancing testing without needing Kube nodes in multiple regions.
// We do this by overriding the calling service locality with the istio-locality label and the serving
// services with a service entry.
//
// The created Service Entry for fake-external-service.com points the domain at services
// that exist internally in the mesh. In the Service Entry we set service B to be in the same locality
// as the caller (A) and service C to be in a different locality. We then verify that all request go to
// service B. For details check the Service Entry configuration at the bottom of the page.
//
//  CDS Test
//
//                                                  100% +-> B (notregion.zone.subzone)
//                                                       |
//                                                       |
//                                                       |
// A (region.zone.subzone) -> fake-external-service.com -|-> C (othernotregion.zone.subzone)
//                                                       |
//                                                       |
//                                                       |
//                                                       +-> NonExistantService (region.zone.subzone)
//
//
//  EDS Test
//
//                                                       +-> 10.28.1.138 (B -> region.zone.subzone)
//                                                       |
//                                                       |
//                                                       |
// A (region.zone.subzone) -> fake-external-service.com -|
//                                                       |
//                                                       |
//                                                       |
//                                                       +-> 10.28.1.139 (C -> notregion.notzone.notsubzone)

const (
	testDuration = 5 * time.Second
)

var serviceBHostname = regexp.MustCompile("^b-.*$")

type FailoverServiceConfig struct {
	Name               string
	Host               string
	Namespace          string
	Resolution         string
	ServiceBAddress    string
	ServiceCAddress    string
	NonExistantService string
}

func TestLocalityFailover(t *testing.T) {
	ctx := framework.NewContext(t)
	defer ctx.Done(t)

	ctx.RequireOrSkip(t, environment.Kube)
	ctx.Environment().(*kube.Environment).Settings()

	// Share Istio Control Plane
	g := galley.NewOrFail(t, ctx, galley.Config{})
	p := pilot.NewOrFail(t, ctx, pilot.Config{Galley: g})

	rand.Seed(time.Now().UnixNano())
	t.Run("TestCDS", func(t *testing.T) {
		testCDS(t, ctx, g, p)
	})

	// t.Run("TestFailoverEDS", func(t *testing.T) {
	// 	testFailoverEDS(t, ctx, g, p)
	// })

}

func testCDS(t *testing.T, ctx resource.Context, g galley.Instance, p pilot.Instance) {
	instance := apps.NewOrFail(ctx, t, apps.Config{Pilot: p})
	a := instance.GetAppOrFail("a", t).(apps.KubeApp)

	fakeHostname := fmt.Sprintf("fake-cds-external-service-%v.com", rand.Int())
	se := FailoverServiceConfig{
		"lplb-failover-cds-service-entry",
		fakeHostname,
		instance.Namespace().Name(),
		"DNS",
		"b",
		"c",
		"nonexistantservice",
	}

	tmpl, _ := template.New("CDSServiceConfig").Parse(fakeExternalServiceConfig)
	var buf bytes.Buffer
	tmpl.Execute(&buf, se)
	g.ApplyConfigOrFail(t, instance.Namespace(), buf.String())

	// TODO: find a better way to do this!
	// This sleep allows config to propagate
	time.Sleep(10 * time.Second)

	// Send traffic to service B via a service entry for the test duration.
	log.Infof("Sending traffic to local service (CDS) via %v", fakeHostname)
	locality.SendTraffic(t, testDuration, a, fakeHostname, serviceBHostname)
}

func testFailoverEDS(t *testing.T, ctx resource.Context, g galley.Instance, p pilot.Instance) {
	// 	t.Parallel()
	// 	instance := apps.NewOrFail(ctx, t, apps.Config{Pilot: p})
	// 	a := instance.GetAppOrFail("a", t).(apps.KubeApp)
	// 	b := instance.GetAppOrFail("b", t).(apps.KubeApp)
	// 	c := instance.GetAppOrFail("c", t).(apps.KubeApp)

	// 	se := ServiceConfig{
	// 		"STATIC",
	// 		b.EndpointForPort(80).NetworkEndpoint().Address,
	// 		c.EndpointForPort(80).NetworkEndpoint().Address,
	// 	}
	// 	tmpl, _ := template.New("EDSServiceConfig").Parse(fakeExternalFailoverServiceConfig)
	// 	var buf bytes.Buffer
	// 	tmpl.Execute(&buf, se)
	// 	g.ApplyConfigOrFail(t, instance.Namespace(), buf.String())

	// 	// TODO: find a better way to do this!
	// 	// This sleep allows config to propagate
	// 	time.Sleep(20 * time.Second)

	// 	// Send traffic to service B via a service entry for the test duration.
	// 	wg := &sync.WaitGroup{}
	// 	wg.Add(numSendTasks)
	// 	log.Info("Sending traffic to local service (EDS)")
	// 	for i := 0; i < numSendTasks; i++ {
	// 		go sendTraffic(t, testDuration, a, "fake-external-service.com", wg)
	// 	}
	// 	wg.Wait()
}

var fakeExternalServiceConfig = `
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  hosts:
  - {{.Host}}
  exportTo:
  - "."
  ports:
  - number: 80
    name: http
    protocol: HTTP
  resolution: {{.Resolution}}
  location: MESH_EXTERNAL
  endpoints:
  - address: {{.NonExistantService}}
    locality: region/zone/subzone
  - address: {{.ServiceBAddress}}
    locality: notregion/zone/subzone
  - address: {{.ServiceCAddress}}
    locality: othernotregion/zone/subzone
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: fake-service-route
  namespace: {{.Namespace}}
spec:
  hosts:
  - {{.Host}}
  http:
  - route:
    - destination:
        host: {{.Host}}
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: fake-service-destination
  namespace: {{.Namespace}}
spec:
  host: {{.Host}}
  trafficPolicy:
    outlierDetection:
      consecutiveErrors: 100
      interval: 1s
      baseEjectionTime: 3m
      maxEjectionPercent: 100
`
