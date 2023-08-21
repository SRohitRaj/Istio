package xdstest

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"sort"
	"testing"
)

type LbEpInfo struct {
	Address string
	// nolint: structcheck
	Weight uint32
}

type LocLbEpInfo struct {
	LbEps  []LbEpInfo
	Weight uint32
}

func (i LocLbEpInfo) GetAddrs() []string {
	addrs := make([]string, 0)
	for _, ep := range i.LbEps {
		addrs = append(addrs, ep.Address)
	}
	return addrs
}

func CompareEndpointsOrFail(t *testing.T, cluster string, got []*endpointv3.LocalityLbEndpoints, want []LocLbEpInfo) {
	if err := CompareEndpoints(cluster, got, want); err != nil {
		t.Error(err)
	}
}

func CompareEndpoints(cluster string, got []*endpointv3.LocalityLbEndpoints, want []LocLbEpInfo) error {
	if len(got) != len(want) {
		return fmt.Errorf("unexpected number of filtered endpoints for %s: got %v, want %v", cluster, len(got), len(want))
	}

	sort.Slice(got, func(i, j int) bool {
		addrI := got[i].LbEndpoints[0].GetEndpoint().Address.GetSocketAddress().Address
		addrJ := got[j].LbEndpoints[0].GetEndpoint().Address.GetSocketAddress().Address
		return addrI < addrJ
	})

	for i, ep := range got {
		if len(ep.LbEndpoints) != len(want[i].LbEps) {
			return fmt.Errorf("unexpected number of LB endpoints within endpoint %d: %v, want %v",
				i, getLbEndpointAddrs(ep), want[i].GetAddrs())
		}

		if ep.LoadBalancingWeight.GetValue() != want[i].Weight {
			return fmt.Errorf("unexpected weight for endpoint %d: got %v, want %v", i, ep.LoadBalancingWeight.GetValue(), want[i].Weight)
		}

		for _, lbEp := range ep.LbEndpoints {
			addr := lbEp.GetEndpoint().Address.GetSocketAddress().Address
			found := false
			for _, wantLbEp := range want[i].LbEps {
				if addr == wantLbEp.Address {
					found = true

					// Now compare the weight.
					if lbEp.GetLoadBalancingWeight().Value != wantLbEp.Weight {
						return fmt.Errorf("unexpected weight for endpoint %s: got %v, want %v",
							addr, lbEp.GetLoadBalancingWeight().Value, wantLbEp.Weight)
					}
					break
				}
			}
			if !found {
				return fmt.Errorf("unexpected address for endpoint %d: %v", i, addr)
			}
		}
	}
	return nil
}

func getLbEndpointAddrs(ep *endpointv3.LocalityLbEndpoints) []string {
	addrs := make([]string, 0)
	for _, lbEp := range ep.LbEndpoints {
		addrs = append(addrs, lbEp.GetEndpoint().Address.GetSocketAddress().Address)
	}
	return addrs
}
