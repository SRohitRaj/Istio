// Copyright Istio Authors.
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

package install

import (
	"fmt"

	"istio.io/istio/pkg/config/schema"
	appsv1 "k8s.io/api/apps/v1"
)

func verifyDeploymentStatus(deployment *appsv1.Deployment) error {
	cond := getDeploymentCondition(deployment.Status, appsv1.DeploymentProgressing)
	if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
		return fmt.Errorf("deployment %q exceeded its progress deadline", deployment.Name)
	}
	if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		return fmt.Errorf("waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated",
			deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas)
	}
	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		return fmt.Errorf("waiting for deployment %q rollout to finish: %d old replicas are pending termination",
			deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas)
	}
	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		return fmt.Errorf("waiting for deployment %q rollout to finish: %d of %d updated replicas are available",
			deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
	}
	return nil
}

func getDeploymentCondition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func findResourceInSpec(kind string) string {
	for _, c := range schema.MustGet().KubeCollections().All() {
		if c.Resource().Kind() == kind {
			return c.Resource().Plural()
		}
	}
	return ""
}
