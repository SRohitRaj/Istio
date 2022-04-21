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

	"go.uber.org/atomic"

	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/kube"
)

// Cluster defines cluster struct
type Cluster struct {
	// ID of the cluster.
	ID cluster.ID
	// SyncTimeout is marked after features.RemoteClusterTimeout.
	SyncTimeout *atomic.Bool
	// Client for accessing the cluster.
	Client kube.Client

	kubeConfigSha [sha256.Size]byte

	stop chan struct{}
	// initialSync is marked when RunAndWait completes
	initialSync *atomic.Bool
}

// Stop channel which is closed when the cluster is removed or the Controller that created the client is stopped.
// Client.RunAndWait is called using this channel.
func (r *Cluster) Stop() <-chan struct{} {
	return r.stop
}

// Run starts the cluster's informers and waits for caches to sync. Once caches are synced, we mark the cluster synced.
// This should be called after each of the handlers have registered informers, and should be run in a goroutine.
func (r *Cluster) Run() {
	r.Client.RunAndWait(r.Stop())
	r.initialSync.Store(true)
}

func (r *Cluster) HasSynced() bool {
	return r.initialSync.Load() || r.SyncTimeout.Load()
}

func (r *Cluster) SyncDidTimeout() bool {
	return r.SyncTimeout.Load() && !r.HasSynced()
}
