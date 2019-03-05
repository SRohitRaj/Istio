package pilot

import (
	"path/filepath"
	"testing"
	"time"

	"istio.io/istio/pkg/test/framework2"
	"istio.io/istio/pkg/test/framework2/components/galley"
	pilot2 "istio.io/istio/pkg/test/framework2/components/pilot"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	xdscore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/structpath"
)

func TestSidecarListeners(t *testing.T) {
	// Call Requires to explicitly initialize dependencies that the test needs.
	ctx := framework2.NewContext(t)
	// TODO - remove prior to checkin
	scopes.Framework.SetOutputLevel(log.DebugLevel)

	// TODO: Limit to Native environment until the Kubernetes environment is supported in the Galley
	// component

	galley := galley.NewOrFail(t, ctx)
	pilotInst := pilot2.NewOrFail(t, ctx, &pilot2.Config{Galley: galley})

	// Simulate proxy identity of a sidecar ...
	nodeID := &model.Proxy{
		ClusterID:   "integration-test",
		Type:        model.SidecarProxy,
		IPAddresses: []string{"10.2.0.1"},
		ID:          "app3.testns",
		DNSDomains:  []string{"testns.cluster.local"},
	}

	// ... and get listeners from Pilot for that proxy
	req := &xdsapi.DiscoveryRequest{
		Node: &xdscore.Node{
			Id: nodeID.ServiceNode(),
		},
		TypeUrl: "type.googleapis.com/envoy.api.v2.Listener",
	}
	// Start the xDS stream
	err := pilotInst.StartDiscovery(req)
	if err != nil {
		t.Fatalf("Failed to test as no resource accepted: %v", err)
	}

	// Test the empty case where no config is loaded
	err = pilotInst.WatchDiscovery(time.Second*10,
		func(response *xdsapi.DiscoveryResponse) (b bool, e error) {
			validator := structpath.AssertThatProto(t, response)
			if !validator.Accept("{.resources[?(@.address.socketAddress.portValue==%v)]}", 15001) {
				return false, nil
			}
			validateListenersNoConfig(t, validator)
			return true, nil
		})
	if err != nil {
		t.Fatalf("Failed to test as no resource accepted: %v", err)
	}

	// Apply some config
	path, err := filepath.Abs("../../testdata/config")
	if err != nil {
		t.Fatalf("No such directory: %v", err)
	}
	err = galley.ApplyConfigDir(path)
	if err != nil {
		t.Fatalf("Error applying directory: %v", err)
	}

	// Now continue to watch on the same stream
	err = pilotInst.WatchDiscovery(time.Second*10,
		func(response *xdsapi.DiscoveryResponse) (b bool, e error) {
			validator := structpath.AssertThatProto(t, response)
			if !validator.Accept("{.resources[?(@.address.socketAddress.portValue==27018)]}") {
				return false, nil
			}
			validateMongoListener(t, validator)
			return true, nil
		})
	if err != nil {
		t.Fatalf("Failed to test as no resource accepted: %v", err)
	}
}

func validateListenersNoConfig(t *testing.T, response *structpath.Structpath) {
	t.Run("validate-legacy-port-3333", func(t *testing.T) {
		// Deprecated: Should be removed as no longer needed
		response.ForTest(t).
			Select("{.resources[?(@.address.socketAddress.portValue==3333)]}").
			Equals("10.2.0.1", "{.address.socketAddress.address}").
			Equals("envoy.tcp_proxy", "{.filterChains[0].filters[*].name}").
			Equals("inbound|3333|http|mgmtCluster", "{.filterChains[0].filters[*].config.cluster}").
			Equals(false, "{.deprecatedV1.bindToPort}").
			NotExists("{.useOriginalDst}")
	})
	t.Run("validate-legacy-port-9999", func(t *testing.T) {
		// Deprecated: Should be removed as no longer needed
		response.ForTest(t).
			Select("{.resources[?(@.address.socketAddress.portValue==9999)]}").
			Equals("10.2.0.1", "{.address.socketAddress.address}").
			Equals("envoy.tcp_proxy", "{.filterChains[0].filters[*].name}").
			Equals("inbound|9999|custom|mgmtCluster", "{.filterChains[0].filters[*].config.cluster}").
			Equals(false, "{.deprecatedV1.bindToPort}").
			NotExists("{.useOriginalDst}")
	})
	t.Run("iptables-forwarding-listener", func(t *testing.T) {
		response.ForTest(t).
			Select("{.resources[?(@.address.socketAddress.portValue==15001)]}").
			Equals("virtual", "{.name}").
			Equals("0.0.0.0", "{.address.socketAddress.address}").
			Equals("envoy.tcp_proxy", "{.filterChains[0].filters[*].name}").
			Equals("BlackHoleCluster", "{.filterChains[0].filters[0].config.cluster}").
			Equals("BlackHoleCluster", "{.filterChains[0].filters[0].config.stat_prefix}").
			Equals(true, "{.useOriginalDst}")
	})
}

func validateMongoListener(t *testing.T, response *structpath.Structpath) {
	t.Run("validate-mongo-listener", func(t *testing.T) {
		mixerListener := response.ForTest(t).
			Select("{.resources[?(@.address.socketAddress.portValue==%v)]}", 27018)

		mixerListener.
			Equals("0.0.0.0", "{.address.socketAddress.address}").
			// Example doing a struct comparison, note the pain with oneofs....
			Equals(&xdscore.SocketAddress{
				Address: "0.0.0.0",
				PortSpecifier: &xdscore.SocketAddress_PortValue{
					PortValue: uint32(27018),
				},
			}, "{.address.socketAddress}").
			Select("{.filterChains[0].filters[0]}").
			Equals("envoy.mongo_proxy", "{.name}").
			Select("{.config}").
			Exists("{.stat_prefix}")
	})
}

// Capturing TestMain allows us to:
// - Do cleanup before exit
// - process testing specific flags
func TestMain(m *testing.M) {
	framework2.RunSuite("sidecar_api_test", m, nil)
}
