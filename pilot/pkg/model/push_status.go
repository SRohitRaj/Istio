// Copyright 2018 Istio Authors
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

package model

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	networking "istio.io/api/networking/v1alpha3"
)

// PushStatus tracks the status of a mush - metrics and errors.
// Metrics are reset after a push - at the beginning all
// values are zero, and when push completes the status is reset.
// The struct is exposed in a debug endpoint - fields public to allow
// easy serialization as json.
type PushStatus struct {
	// TODO: will be renamed as PushContext

	mutex sync.Mutex

	// ProxyStatus is keyed by the error code, and holds a map keyed
	// by the ID.
	ProxyStatus map[string]map[string]PushStatusEvent

	// Start represents the time of last config change that reset the
	// push status.
	Start time.Time

	End time.Time

	// ContextMutex is used to sync the data cache.
	// Currently it is used directly - to avoid making copies of the
	// structs. All data is set when the PushStatus object is populated,
	// from a single thread - only read locks are needed, data should not
	// be changed by plugins
	Mutex sync.RWMutex

	// Services list all services in the system at the time push started.
	Services []*Service

	//ServicesByName map[string]*Service
	//
	//ServiceAttributes map[string]*ServiceAttributes
	//
	//ConfigsByType map[string][]*Config

	// TODO: add the remaining O(n**2) model, deprecate/remove all remaining
	// uses of model:

	//Endpoints map[string][]*ServiceInstance
	//ServicesForProxy map[string][]*ServiceInstance
	//ManagementPorts map[string]*PortList
	//WorkloadHealthCheck map[string]*ProbeList

	// ServiceAccounts represents the list of service accounts
	// for a service.
	//	ServiceAccounts map[string][]string
	// Temp: the code in alpha3 should use VirtualService directly
	VirtualServiceConfigs []Config
	//TODO: gateways              []*networking.Gateway
}

// PushStatusEvent represents an event captured by push status.
// It may contain additional message and the affected proxy.
type PushStatusEvent struct {
	Proxy   string `json:"proxy,omitempty"`
	Message string `json:"message,omitempty"`
}

// PushMetric wraps a prometheus metric.
type PushMetric struct {
	Name  string
	gauge prometheus.Gauge
}

func newPushMetric(name, help string) *PushMetric {
	pm := &PushMetric{
		gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}),
		Name: name,
	}
	prometheus.MustRegister(pm.gauge)
	metrics = append(metrics, pm)
	return pm
}

// Add will add an case to the metric.
func (ps *PushStatus) Add(metric *PushMetric, key string, proxy *Proxy, msg string) {
	if ps == nil {
		log.Infof("Metric without context %s %v %s", key, proxy, msg)
		return
	}
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	metricMap, f := ps.ProxyStatus[metric.Name]
	if !f {
		metricMap = map[string]PushStatusEvent{}
		ps.ProxyStatus[metric.Name] = metricMap
	}
	ev := PushStatusEvent{Message: msg}
	if proxy != nil {
		ev.Proxy = proxy.ID
	}
	metricMap[key] = ev
}

var (
	// ProxyStatusNoService represents proxies not selected by any service
	// This can be normal - for workloads that act only as client, or are not covered by a Service.
	// It can also be an error, for example in cases the Endpoint list of a service was not updated by the time
	// the sidecar calls.
	// Updated by GetProxyServiceInstances
	ProxyStatusNoService = newPushMetric(
		"pilot_no_ip",
		"Pods not found in the endpoint table, possibly invalid.",
	)

	// ProxyStatusEndpointNotReady represents proxies found not be ready.
	// Updated by GetProxyServiceInstances. Normal condition when starting
	// an app with readiness, error if it doesn't change to 0.
	ProxyStatusEndpointNotReady = newPushMetric(
		"pilot_endpoint_not_ready",
		"Endpoint found in unready state.",
	)

	// ProxyStatusConflictOutboundListenerTCPOverHTTP metric tracks number of
	// wildcard TCP listeners that conflicted with existing wildcard HTTP listener on same port
	ProxyStatusConflictOutboundListenerTCPOverHTTP = newPushMetric(
		"pilot_conflict_outbound_listener_tcp_over_current_http",
		"Number of conflicting wildcard tcp listeners with current wildcard http listener.",
	)

	// ProxyStatusConflictOutboundListenerTCPOverTCP metric tracks number of
	// TCP listeners that conflicted with existing TCP listeners on same port
	ProxyStatusConflictOutboundListenerTCPOverTCP = newPushMetric(
		"pilot_conflict_outbound_listener_tcp_over_current_tcp",
		"Number of conflicting tcp listeners with current tcp listener.",
	)

	// ProxyStatusConflictOutboundListenerHTTPOverTCP metric tracks number of
	// wildcard HTTP listeners that conflicted with existing wildcard TCP listener on same port
	ProxyStatusConflictOutboundListenerHTTPOverTCP = newPushMetric(
		"pilot_conflict_outbound_listener_http_over_current_tcp",
		"Number of conflicting wildcard http listeners with current wildcard tcp listener.",
	)

	// ProxyStatusConflictInboundListener tracks cases of multiple inbound
	// listeners - 2 services selecting the same port of the pod.
	ProxyStatusConflictInboundListener = newPushMetric(
		"pilot_conflict_inbound_listener",
		"Number of conflicting inbound listeners.",
	)

	// ProxyStatusClusterNoInstances tracks clusters (services) without workloads.
	ProxyStatusClusterNoInstances = newPushMetric(
		"pilot_eds_no_instances",
		"Number of clusters without instances.",
	)

	// DuplicatedDomains tracks rejected VirtualServices due to duplicated hostname.
	DuplicatedDomains = newPushMetric(
		"pilot_vservice_dup_domain",
		"Virtual services with dup domains.",
	)

	// LastPushStatus preserves the metrics and data collected during lasts global push.
	// It can be used by debugging tools to inspect the push event. It will be reset after each push with the
	// new version.
	LastPushStatus *PushStatus

	// All metrics we registered.
	metrics []*PushMetric
)

// NewStatus creates a new PushStatus structure to track push status.
func NewStatus() *PushStatus {
	// TODO: detect push in progress, don't update status if set
	return &PushStatus{
		ProxyStatus: map[string]map[string]PushStatusEvent{},
		Start:       time.Now(),
	}
}

// JSON implements json.Marshaller, with a lock.
func (ps *PushStatus) JSON() ([]byte, error) {
	if ps == nil {
		return []byte{'{', '}'}, nil
	}
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	return json.MarshalIndent(ps, "", "    ")
}

// OnConfigChange is called when a config change is detected.
func (ps *PushStatus) OnConfigChange() {
	LastPushStatus = ps
	ps.UpdateMetrics()
}

// UpdateMetrics will update the prometheus metrics based on the
// current status of the push.
func (ps *PushStatus) UpdateMetrics() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	for _, pm := range metrics {
		mmap, f := ps.ProxyStatus[pm.Name]
		if f {
			pm.gauge.Set(float64(len(mmap)))
		} else {
			pm.gauge.Set(0)
		}
	}
}

// VirtualServices lists all virtual services bound to the specified gateways
// This replaces store.VirtualServices
func (ps *PushStatus) VirtualServices(gateways map[string]bool) []Config {
	configs := ps.VirtualServiceConfigs
	sortConfigByCreationTime(configs)
	out := make([]Config, 0)
	for _, config := range configs {
		rule := config.Spec.(*networking.VirtualService)
		if len(rule.Gateways) == 0 {
			// This rule applies only to IstioMeshGateway
			if gateways[IstioMeshGateway] {
				out = append(out, config)
			}
		} else {
			for _, g := range rule.Gateways {
				// note: Gateway names do _not_ use wildcard matching, so we do not use Hostname.Matches here
				if gateways[ResolveShortnameToFQDN(g, config.ConfigMeta).String()] {
					out = append(out, config)
					break
				} else if g == IstioMeshGateway && gateways[g] {
					// "mesh" gateway cannot be expanded into FQDN
					out = append(out, config)
					break
				}
			}
		}
	}

	// Need to parse each rule and convert the shortname to FQDN
	for _, r := range out {
		rule := r.Spec.(*networking.VirtualService)
		// resolve top level hosts
		for i, h := range rule.Hosts {
			rule.Hosts[i] = ResolveShortnameToFQDN(h, r.ConfigMeta).String()
		}
		// resolve gateways to bind to
		for i, g := range rule.Gateways {
			if g != IstioMeshGateway {
				rule.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
			}
		}
		// resolve host in http route.destination, route.mirror
		for _, d := range rule.Http {
			for _, m := range d.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range d.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
			if d.Mirror != nil {
				d.Mirror.Host = ResolveShortnameToFQDN(d.Mirror.Host, r.ConfigMeta).String()
			}
		}
		//resolve host in tcp route.destination
		for _, d := range rule.Tcp {
			for _, m := range d.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range d.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
		}
		//resolve host in tls route.destination
		for _, tls := range rule.Tls {
			for _, m := range tls.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range tls.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
		}
	}

	return out
}
