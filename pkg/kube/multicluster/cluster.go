// Copyright Istio Authors
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

package multicluster

import (
	"crypto/sha256"
	"time"

	"go.uber.org/atomic"
	corev1 "k8s.io/api/core/v1"

	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	filter "istio.io/istio/pkg/kube/namespace"
	"istio.io/istio/pkg/log"
)

// Cluster defines cluster struct
type Cluster struct {
	// ID of the cluster.
	ID cluster.ID
	// Client for accessing the cluster.
	Client kube.Client

	kubeConfigSha [sha256.Size]byte

	stop chan struct{}
	// initialSync is marked when RunAndWait completes
	initialSync *atomic.Bool
	// initialSyncTimeout is set when RunAndWait timed out
	initialSyncTimeout *atomic.Bool
}

// Run starts the cluster's informers and waits for caches to sync. Once caches are synced, we mark the cluster synced.
// This should be called after each of the handlers have registered informers, and should be run in a goroutine.
func (r *Cluster) Run(mesh mesh.Watcher, handlers []handler, buildComponents func(cluster *Cluster)) {
	if features.RemoteClusterTimeout > 0 {
		time.AfterFunc(features.RemoteClusterTimeout, func() {
			if !r.initialSync.Load() {
				log.Errorf("remote cluster %s failed to sync after %v", r.ID, features.RemoteClusterTimeout)
				timeouts.With(clusterLabel.Value(string(r.ID))).Increment()
			}
			r.initialSyncTimeout.Store(true)
		})
	}
	// Build a namespace watcher. This must have no filter, since this is our input to the filter itself.
	// This must be done before we build components, so they can access the filter.
	namespaces := kclient.New[*corev1.Namespace](r.Client)
	filter := filter.NewDiscoveryNamespacesFilter(namespaces, mesh, r.stop)
	kube.SetObjectFilter(r.Client, filter)

	buildComponents(r)

	if !r.Client.RunAndWait(r.stop) {
		log.Warnf("remote cluster %s failed to sync", r.ID)
		return
	}
	for _, h := range handlers {
		if !kube.WaitForCacheSync("cluster"+string(r.ID), r.stop, h.HasSynced) {
			log.Warnf("remote cluster %s failed to sync handler", r.ID)
			return
		}
	}

	r.initialSync.Store(true)
}

// Stop closes the stop channel, if is safe to be called multi times.
func (r *Cluster) Stop() {
	select {
	case <-r.stop:
		return
	default:
		close(r.stop)
	}
}

func (r *Cluster) HasSynced() bool {
	// It could happen when a wrong credential provide, this cluster has no chance to run.
	// In this case, the `initialSyncTimeout` will never be set
	// In order not block istiod start up, check close as well.
	if r.Closed() {
		return true
	}
	return r.initialSync.Load() || r.initialSyncTimeout.Load()
}

func (r *Cluster) Closed() bool {
	select {
	case <-r.stop:
		return true
	default:
		return false
	}
}

func (r *Cluster) SyncDidTimeout() bool {
	return !r.initialSync.Load() && r.initialSyncTimeout.Load()
}
