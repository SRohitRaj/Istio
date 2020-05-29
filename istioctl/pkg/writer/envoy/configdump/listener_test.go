package configdump

import (
	"testing"

	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
)

func TestListenerFilter_Verify(t *testing.T) {
	tests := []struct {
		desc       string
		inFilter   *ListenerFilter
		inListener *listener.Listener
		expect     bool
	}{
		{
			desc: "filter-fields-empty",
			inFilter: &ListenerFilter{
				Address: "",
				Port:    0,
				Type:    "",
			},
			inListener: &listener.Listener{},
			expect:     true,
		},
		{
			desc: "addrs-dont-match",
			inFilter: &ListenerFilter{
				Address: "0.0.0.0",
			},
			inListener: &listener.Listener{
				Address: &v3.Address{
					Address: &v3.Address_SocketAddress{
						SocketAddress: &v3.SocketAddress{Address: "1.1.1.1"},
					},
				},
			},
			expect: false,
		},
		{
			desc: "ports-dont-match",
			inFilter: &ListenerFilter{
				Port: 10,
			},
			inListener: &listener.Listener{
				Address: &v3.Address{
					Address: &v3.Address_SocketAddress{
						SocketAddress: &v3.SocketAddress{
							PortSpecifier: &v3.SocketAddress_PortValue{
								PortValue: 11,
							},
						},
					},
				},
			},
			expect: false,
		},
		{
			desc: "http-type-match",
			inFilter: &ListenerFilter{
				Type: "HTTP",
			},
			inListener: &listener.Listener{
				FilterChains: []*listener.FilterChain{{
					Filters: []*listener.Filter{{
						Name: "envoy.http_connection_manager",
					},
					},
				},
				},
			},
			expect: true,
		},
		{
			desc: "http-tcp-type-match",
			inFilter: &ListenerFilter{
				Type: "HTTP+TCP",
			},
			inListener: &listener.Listener{
				FilterChains: []*listener.FilterChain{{
					Filters: []*listener.Filter{{
						Name: "envoy.tcp_proxy",
					},
						{
							Name: "envoy.tcp_proxy",
						},
						{
							Name: "envoy.http_connection_manager",
						}},
				}},
			},
			expect: true,
		},
		{
			desc: "tcp-type-match",
			inFilter: &ListenerFilter{
				Type: "TCP",
			},
			inListener: &listener.Listener{
				FilterChains: []*listener.FilterChain{{
					Filters: []*listener.Filter{{
						Name: "envoy.tcp_proxy",
					}},
				}},
			},
			expect: true,
		},
		{
			desc: "unknown-type",
			inFilter: &ListenerFilter{
				Type: "UNKNOWN",
			},
			inListener: &listener.Listener{
				FilterChains: []*listener.FilterChain{{
					Filters: []*listener.Filter{},
				}},
			},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := tt.inFilter.Verify(tt.inListener); got != tt.expect {
				t.Errorf("%s: expect %v got %v", tt.desc, tt.expect, got)
			}
		})
	}
}
