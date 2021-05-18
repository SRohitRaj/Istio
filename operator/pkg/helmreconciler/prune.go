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

package helmreconciler

import (
	"context"
	"fmt"
	"strings"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	klabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"istio.io/api/label"
	"istio.io/api/operator/v1alpha1"
	v1alpha12 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	"istio.io/istio/operator/pkg/cache"
	"istio.io/istio/operator/pkg/metrics"
	"istio.io/istio/operator/pkg/name"
	"istio.io/istio/operator/pkg/object"
	"istio.io/istio/operator/pkg/translate"
	"istio.io/istio/operator/pkg/util"
	"istio.io/istio/pkg/proxy"
)

var (
	// NamespacedResources orders non cluster scope resources types which should be deleted, first to last
	NamespacedResources = []schema.GroupVersionKind{
		{Group: "autoscaling", Version: "v2beta1", Kind: name.HPAStr},
		{Group: "policy", Version: "v1beta1", Kind: name.PDBStr},
		{Group: "apps", Version: "v1", Kind: name.DeploymentStr},
		{Group: "apps", Version: "v1", Kind: name.DaemonSetStr},
		{Group: "", Version: "v1", Kind: name.ServiceStr},
		{Group: "", Version: "v1", Kind: name.CMStr},
		{Group: "", Version: "v1", Kind: name.PVCStr},
		{Group: "", Version: "v1", Kind: name.PodStr},
		{Group: "", Version: "v1", Kind: name.SecretStr},
		{Group: "", Version: "v1", Kind: name.SAStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.RoleBindingStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.RoleStr},
		{Group: name.NetworkingAPIGroupName, Version: "v1alpha3", Kind: name.DestinationRuleStr},
		{Group: name.NetworkingAPIGroupName, Version: "v1alpha3", Kind: name.EnvoyFilterStr},
		{Group: name.NetworkingAPIGroupName, Version: "v1alpha3", Kind: name.GatewayStr},
		{Group: name.NetworkingAPIGroupName, Version: "v1alpha3", Kind: name.VirtualServiceStr},
		{Group: name.SecurityAPIGroupName, Version: "v1beta1", Kind: name.PeerAuthenticationStr},
	}

	// ClusterResources are resource types the operator prunes, ordered by which types should be deleted, first to last.
	ClusterResources = []schema.GroupVersionKind{
		{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.MutatingWebhookConfigurationStr},
		{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.ValidatingWebhookConfigurationStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleBindingStr},
		// Cannot currently prune CRDs because this will also wipe out user config.
		// {Group: "apiextensions.k8s.io", Version: "v1beta1", Kind: name.CRDStr},
	}
	// ClusterCPResources lists cluster scope resources types which should be deleted during uninstall command.
	ClusterCPResources = []schema.GroupVersionKind{
		{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.MutatingWebhookConfigurationStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleBindingStr},
	}
	// AllClusterResources lists all cluster scope resources types which should be deleted in purge case, including CRD.
	AllClusterResources = append(ClusterResources,
		schema.GroupVersionKind{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.MutatingWebhookConfigurationStr},
		schema.GroupVersionKind{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.ValidatingWebhookConfigurationStr},
		schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: name.CRDStr},
	)
)

// Prune removes any resources not specified in manifests generated by HelmReconciler h.
func (h *HelmReconciler) Prune(manifests name.ManifestMap, all bool) error {
	return h.runForAllTypes(func(labels map[string]string, objects *unstructured.UnstructuredList) error {
		var errs util.Errors
		if all {
			errs = util.AppendErr(errs, h.deleteResources(nil, labels, "", objects, all))
		} else {
			for cname, manifest := range manifests.Consolidated() {
				errs = util.AppendErr(errs, h.deleteResources(object.AllObjectHashes(manifest), labels, cname, objects, all))
			}
		}
		return errs.ToError()
	})
}

// PruneControlPlaneByRevisionWithController is called to remove specific control plane revision
// during reconciliation process of controller.
// It returns the install status and any error encountered.
func (h *HelmReconciler) PruneControlPlaneByRevisionWithController(iopSpec *v1alpha1.IstioOperatorSpec) (*v1alpha1.InstallStatus, error) {
	ns := v1alpha12.Namespace(iopSpec)
	if ns == "" {
		ns = name.IstioDefaultNamespace
	}
	errStatus := &v1alpha1.InstallStatus{Status: v1alpha1.InstallStatus_ERROR}
	enabledComponents, err := translate.GetEnabledComponents(iopSpec)
	if err != nil {
		return errStatus,
			fmt.Errorf("failed to get enabled components: %v", err)
	}
	pids, err := proxy.GetIDsFromProxyInfo("", "", iopSpec.Revision, ns)
	if err != nil {
		return errStatus,
			fmt.Errorf("failed to check proxy infos: %v", err)
	}
	// TODO(richardwxn): add warning message together with the status
	if len(pids) != 0 {
		msg := fmt.Sprintf("there are proxies still pointing to the pruned control plane: %s.",
			strings.Join(pids, " "))
		st := &v1alpha1.InstallStatus{Status: v1alpha1.InstallStatus_ACTION_REQUIRED, Message: msg}
		return st, nil
	}
	var allUslist []*unstructured.UnstructuredList
	for _, c := range enabledComponents {
		uslist, err := h.GetPrunedResources(iopSpec.Revision, false, c)
		if err != nil {
			return errStatus, err
		}
		allUslist = append(allUslist, uslist...)
	}
	if err := h.DeleteObjectsList(allUslist); err != nil {
		return errStatus, err
	}
	return &v1alpha1.InstallStatus{Status: v1alpha1.InstallStatus_HEALTHY}, nil
}

// DeleteObjectsList removed resources that are in the slice of UnstructuredList.
func (h *HelmReconciler) DeleteObjectsList(objectsList []*unstructured.UnstructuredList) error {
	var errs util.Errors
	deletedObjects := make(map[string]bool)
	for _, objects := range objectsList {
		for _, o := range objects.Items {
			obj := object.NewK8sObject(&o, nil, nil)
			oh := obj.Hash()
			if h.opts.DryRun {
				h.opts.Log.LogAndPrintf("Not deleting object %s because of dry run.", oh)
				continue
			}
			// kube client does not differentiate API version when listing, added this check to deduplicate.
			if deletedObjects[oh] {
				continue
			}
			if o.GetKind() == name.IstioOperatorStr {
				o.SetFinalizers([]string{})
				if err := h.client.Patch(context.TODO(), &o, client.Merge); err != nil {
					scope.Errorf("failed to patch IstioOperator CR: %s, %v", o.GetName(), err)
				}
			}
			err := h.client.Delete(context.TODO(), &o,
				client.PropagationPolicy(metav1.DeletePropagationBackground))
			if err != nil {
				if !kerrors.IsNotFound(err) {
					errs = util.AppendErr(errs, err)
				} else {
					// do not return error if resources are not found
					h.opts.Log.LogAndPrintf("object: %s is not being deleted because it no longer exists",
						obj.Hash())
				}
			} else {
				deletedObjects[oh] = true
				objGvk := o.GroupVersionKind()
				metrics.ResourceDeletionTotal.
					With(metrics.ResourceKindLabel.Value(util.GKString(objGvk.GroupKind()))).
					Increment()
				h.addPrunedKind(objGvk.GroupKind())
				metrics.RemoveResource(obj.FullName(), objGvk.GroupKind())
			}
			h.opts.Log.LogAndPrintf("  Removed %s.", oh)
		}
	}

	return errs.ToError()
}

// GetPrunedResources get the list of resources to be removed
// 1. if includeClusterResources is false, we list the namespaced resources by matching revision and component labels.
// 2. if includeClusterResources is true, we list the namespaced and cluster resources by component labels only.
// If componentName is not empty, only resources associated with specific components would be returned
// UnstructuredList of objects and corresponding list of name kind hash of k8sObjects would be returned
func (h *HelmReconciler) GetPrunedResources(revision string, includeClusterResources bool, componentName string) (
	[]*unstructured.UnstructuredList, error) {
	var usList []*unstructured.UnstructuredList
	labels := make(map[string]string)
	if revision != "" {
		labels[label.IoIstioRev.Name] = revision
	}
	if componentName != "" {
		labels[IstioComponentLabelStr] = componentName
	}
	selector := klabels.Set(labels).AsSelectorPreValidated()
	gvkList := append(NamespacedResources, ClusterCPResources...)
	if includeClusterResources {
		gvkList = append(NamespacedResources, AllClusterResources...)
		if ioplist := h.getIstioOperatorCR(); ioplist.Items != nil {
			usList = append(usList, ioplist)
		}
	}
	for _, gvk := range gvkList {
		objects := &unstructured.UnstructuredList{}
		objects.SetGroupVersionKind(gvk)
		componentRequirement, err := klabels.NewRequirement(IstioComponentLabelStr, selection.Exists, nil)
		if err != nil {
			return usList, err
		}
		if includeClusterResources {
			s := klabels.NewSelector()
			err = h.client.List(context.TODO(), objects,
				client.MatchingLabelsSelector{Selector: s.Add(*componentRequirement)})
		} else {
			// do not prune base components or unknown components
			includeCN := []string{
				string(name.PilotComponentName), string(name.IstiodRemoteComponentName),
				string(name.IngressComponentName), string(name.EgressComponentName),
				string(name.CNIComponentName), string(name.IstioOperatorComponentName),
			}
			includeRequirement, err := klabels.NewRequirement(IstioComponentLabelStr, selection.In, includeCN)
			if err != nil {
				return usList, err
			}
			if err = h.client.List(context.TODO(), objects,
				client.MatchingLabelsSelector{
					Selector: selector.Add(*includeRequirement, *componentRequirement),
				},
			); err != nil {
				continue
			}
		}
		if err != nil {
			continue
		}
		for _, obj := range objects.Items {
			objName := fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
			metrics.AddResource(objName, gvk.GroupKind())
		}
		usList = append(usList, objects)
	}

	return usList, nil
}

// getIstioOperatorCR is a helper function to get IstioOperator CR during purge,
// otherwise the resources would be reconciled back later if there is in-cluster operator deployment.
// And it is needed to remove the IstioOperator CRD.
func (h *HelmReconciler) getIstioOperatorCR() *unstructured.UnstructuredList {
	objects := &unstructured.UnstructuredList{}
	objects.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "install.istio.io",
		Version: "v1alpha1", Kind: name.IstioOperatorStr,
	})
	if err := h.client.List(context.TODO(), objects); err != nil {
		scope.Errorf("failed to list IstioOperator CR: %v", err)
	}
	return objects
}

// DeleteControlPlaneByManifests removed resources by manifests with matching revision label.
// If purge option is set to true, all manifests would be removed regardless of labels match.
func (h *HelmReconciler) DeleteControlPlaneByManifests(manifestMap name.ManifestMap,
	revision string, includeClusterResources bool) error {
	labels := map[string]string{
		operatorLabelStr: operatorReconcileStr,
	}
	cpManifestMap := make(name.ManifestMap)
	if revision != "" {
		labels[label.IoIstioRev.Name] = revision
	}
	if !includeClusterResources {
		// only delete istiod resources if revision is empty and --purge flag is not true.
		cpManifestMap[name.PilotComponentName] = manifestMap[name.PilotComponentName]
		manifestMap = cpManifestMap
	}
	for cn, mf := range manifestMap.Consolidated() {
		if cn == string(name.IstioBaseComponentName) && !includeClusterResources {
			continue
		}
		objects, err := object.ParseK8sObjectsFromYAMLManifest(mf)
		if err != nil {
			return fmt.Errorf("failed to parse k8s objects from yaml: %v", err)
		}
		if objects == nil {
			continue
		}
		unstructuredObjects := unstructured.UnstructuredList{}
		for _, obj := range objects {
			if h.opts.DryRun {
				h.opts.Log.LogAndPrintf("Not deleting object %s because of dry run.", obj.Hash())
				continue
			}
			obju := obj.UnstructuredObject()
			if err := h.applyLabelsAndAnnotations(obju, cn); err != nil {
				return err
			}
			unstructuredObjects.Items = append(unstructuredObjects.Items, *obju)
		}
		if err := h.deleteResources(nil, labels, cn, &unstructuredObjects, includeClusterResources); err != nil {
			return fmt.Errorf("failed to delete resources: %v", err)
		}
	}
	return nil
}

// pruneAllTypes will collect all existing resource types we care about. For each type, the callback function
// will be called with the labels used to select this type, and all objects.
// This is in internal function meant to support prune and delete
func (h *HelmReconciler) runForAllTypes(callback func(labels map[string]string, objects *unstructured.UnstructuredList) error) error {
	var errs util.Errors
	// Ultimately, we want to prune based on component labels. Each of these share a common set of labels
	// Rather than do N List() calls for each component, we will just filter for the common subset here
	// and each component will do its own filtering
	// Because we are filtering by the core labels, List() will only return items that some components will care
	// about, so we are not querying for an overly broad set of resources.
	labels, err := h.getCoreOwnerLabels()
	if err != nil {
		return err
	}
	selector := klabels.Set(labels).AsSelectorPreValidated()
	for _, gvk := range append(NamespacedResources, ClusterResources...) {
		// First, we collect all objects for the provided GVK
		objects := &unstructured.UnstructuredList{}
		objects.SetGroupVersionKind(gvk)
		componentRequirement, err := klabels.NewRequirement(IstioComponentLabelStr, selection.Exists, nil)
		if err != nil {
			return err
		}
		selector = selector.Add(*componentRequirement)
		if err := h.client.List(context.TODO(), objects, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			// we only want to retrieve resources clusters
			scope.Warnf("retrieving resources to prune type %s: %s not found", gvk.String(), err)
			continue
		}
		for _, obj := range objects.Items {
			objName := fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
			metrics.AddResource(objName, gvk.GroupKind())
		}
		errs = util.AppendErr(errs, callback(labels, objects))
	}
	return errs.ToError()
}

// deleteResources delete any resources from the given component that are not in the excluded map. Resource
// labels are used to identify the resources belonging to the component.
func (h *HelmReconciler) deleteResources(excluded map[string]bool, coreLabels map[string]string,
	componentName string, objects *unstructured.UnstructuredList, all bool) error {
	var errs util.Errors
	labels := h.addComponentLabels(coreLabels, componentName)
	selector := klabels.Set(labels).AsSelectorPreValidated()
	for _, o := range objects.Items {
		obj := object.NewK8sObject(&o, nil, nil)
		oh := obj.Hash()
		if !all {
			// Label mismatch. Provided objects don't select against the component, so this likely means the object
			// is for another component.
			if !selector.Matches(klabels.Set(o.GetLabels())) {
				continue
			}
			if excluded[oh] {
				continue
			}
		}
		if h.opts.DryRun {
			h.opts.Log.LogAndPrintf("Not pruning object %s because of dry run.", oh)
			continue
		}
		err := h.client.Delete(context.TODO(), &o, client.PropagationPolicy(metav1.DeletePropagationBackground))
		scope.Infof("Deleting %s (%s/%v)", obj.Hash(), h.iop.Name, h.iop.Spec.Revision)
		objGvk := o.GroupVersionKind()
		if err != nil {
			if !kerrors.IsNotFound(err) {
				errs = util.AppendErr(errs, err)
			} else {
				// do not return error if resources are not found
				h.opts.Log.LogAndPrintf("object: %s is not being deleted because it no longer exists", obj.Hash())
				continue
			}
		}
		if !all {
			h.removeFromObjectCache(componentName, oh)
		}
		metrics.ResourceDeletionTotal.
			With(metrics.ResourceKindLabel.Value(util.GKString(objGvk.GroupKind()))).
			Increment()
		h.addPrunedKind(objGvk.GroupKind())
		metrics.RemoveResource(obj.FullName(), objGvk.GroupKind())
		h.opts.Log.LogAndPrintf("  Removed %s.", oh)
	}
	if all {
		cache.FlushObjectCaches()
	}

	return errs.ToError()
}

// RemoveObject removes object with objHash in componentName from the object cache.
func (h *HelmReconciler) removeFromObjectCache(componentName, objHash string) {
	crHash, err := h.getCRHash(componentName)
	if err != nil {
		scope.Error(err.Error())
	}
	cache.RemoveObject(crHash, objHash)
	scope.Infof("Removed object %s from Cache.", objHash)
}
