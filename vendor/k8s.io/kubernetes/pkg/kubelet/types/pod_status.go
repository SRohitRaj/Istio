/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"k8s.io/api/core/v1"
)

// PodConditionsByKubelet is the list of pod conditions owned by kubelet
var PodConditionsByKubelet = []v1.PodConditionType{
	v1.PodScheduled,
	v1.PodReady,
	v1.PodInitialized,
	v1.PodReasonUnschedulable,
	v1.ContainersReady,
}

// PodConditionByKubelet returns if the pod condition type is owned by kubelet
func PodConditionByKubelet(conditionType v1.PodConditionType) bool {
	for _, c := range PodConditionsByKubelet {
		if c == conditionType {
			return true
		}
	}
	return false
}
