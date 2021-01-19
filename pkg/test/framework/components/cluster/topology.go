//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package cluster

import (
	"fmt"

	"istio.io/istio/pkg/test/framework/resource"
)

// Map can be given as a shared reference to multiple Topology/Cluster implemetations
// allowing clusters to find each other for lookups of Primary, ConfigCluster, etc.
type Map = map[string]resource.Cluster

func NewTopology(
	name string,
	networkName string,
	controlPlaneCluster string,
	configCluster string,
	clusters Map,
) Topology {
	return Topology{
		name:                    name,
		networkName:             networkName,
		controlPlaneClusterName: controlPlaneCluster,
		configClusterName:       configCluster,
		clusters:                clusters,
	}
}

// Topology gives information about the relationship between clusters.
// Cluster implementations can embed this struct to include common functionality.
type Topology struct {
	name                    string
	networkName             string
	controlPlaneClusterName string
	configClusterName       string
	// clusters should contain all clusters in the context
	clusters map[string]resource.Cluster
}

// NetworkName the cluster is on
func (c Topology) NetworkName() string {
	return c.networkName
}

// Name provides the name this cluster used by Istio.
func (c Topology) Name() string {
	return c.name
}

func (c Topology) IsPrimary() bool {
	return c.Primary().Name() == c.Name()
}

func (c Topology) IsConfig() bool {
	return c.Config().Name() == c.Name()
}

func (c Topology) IsRemote() bool {
	return !c.IsPrimary()
}

func (c Topology) Primary() resource.Cluster {
	cluster, ok := c.clusters[c.controlPlaneClusterName]
	if !ok || cluster == nil {
		panic(fmt.Errorf("cannot find %s, the primary cluster for %s", c.controlPlaneClusterName, c.Name()))
	}
	return cluster
}

func (c Topology) Config() resource.Cluster {
	cluster, ok := c.clusters[c.configClusterName]
	if !ok || cluster == nil {
		panic(fmt.Errorf("cannot find %s, the config cluster for %s", c.configClusterName, c.Name()))
	}
	return cluster
}

func (c Topology) WithPrimary(primaryClusterName string) Topology {
	// TODO remove this, should only be provided by external config
	c.controlPlaneClusterName = primaryClusterName
	return c
}

func (c Topology) WithConfig(configClusterName string) Topology {
	// TODO remove this, should only be provided by external config
	c.configClusterName = configClusterName
	return c
}
