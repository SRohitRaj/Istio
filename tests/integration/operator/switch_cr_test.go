// Copyright 2020 Istio Authors
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

package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"

	api "istio.io/api/operator/v1alpha1"
	"istio.io/istio/operator/pkg/object"
	"istio.io/istio/operator/pkg/util"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/istioctl"
	"istio.io/istio/pkg/test/framework/image"
	"istio.io/istio/pkg/test/framework/resource/environment"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/pkg/log"
)

const (
	IstioNamespace = "istio-system"
	pollInterval   = time.Second
	pollTimeOut    = 100 * time.Second
)

var (
	ManifestTestDataPath = filepath.Join(env.IstioSrc, "operator/cmd/mesh/testdata/manifest-generate/input")
	ProfilesPath         = filepath.Join(env.IstioSrc, "operator/data/profiles")
)

func TestController(t *testing.T) {
	framework.
		NewTest(t).
		RequiresEnvironment(environment.Kube).
		Run(func(ctx framework.TestContext) {
			istioCtl := istioctl.NewOrFail(ctx, ctx, istioctl.Config{})
			workDir, err := ctx.CreateTmpDirectory("operator-controller-test")
			if err != nil {
				t.Fatal("failed to create test directory")
			}

			checkControllerInstallation(t, ctx, istioCtl, workDir, path.Join(ManifestTestDataPath, "all_on.yaml"))
			checkControllerInstallation(t, ctx, istioCtl, workDir, path.Join(ProfilesPath, "default.yaml"))
			checkControllerInstallation(t, ctx, istioCtl, workDir, path.Join(ProfilesPath, "demo.yaml"))
		})
}

// checkInstallStatus check the status of IstioOperator CR from the cluster
func checkInstallStatus(cs kube.Cluster) error {
	log.Info("checking IstioOperator CR status")
	gvr := schema.GroupVersionResource{
		Group:    "install.istio.io",
		Version:  "v1alpha1",
		Resource: "istiooperators",
	}
	conditionF := func() (done bool, err error) {
		us, err := cs.GetUnstructured(gvr, "istio-system", "test-istiocontrolplane")
		if err != nil {
			return false, fmt.Errorf("failed to get istioOperator resource: %v", err)
		}
		usIOPStatus := us.UnstructuredContent()["status"].(map[string]interface{})
		iopStatusString, err := json.Marshal(usIOPStatus)
		if err != nil {
			return false, fmt.Errorf("failed to marshal istioOperator status: %v", err)
		}
		status := &api.InstallStatus{}
		jspb := jsonpb.Unmarshaler{AllowUnknownFields: true}
		if err := jspb.Unmarshal(bytes.NewReader(iopStatusString), status); err != nil {
			return false, fmt.Errorf("failed to unmarshal istioOperator status: %v", err)
		}
		if status.Status != api.InstallStatus_HEALTHY {
			return false, fmt.Errorf("expect IstioOperator status to be healthy, but got: %v", status.Status)
		}
		var errs util.Errors
		for cn, cnstatus := range status.ComponentStatus {
			if cnstatus.Status != api.InstallStatus_HEALTHY {
				errs = util.AppendErr(errs, fmt.Errorf("expect component: %s status to be healthy,"+
					" but got: %v", cn, cnstatus.Status))
			}
		}
		return true, nil
	}
	if errPoll := wait.Poll(pollInterval, pollTimeOut, conditionF); errPoll != nil {
		return fmt.Errorf("failed to poll IstioOperator status: %v", errPoll)
	}
	return nil
}

func checkControllerInstallation(t *testing.T, ctx framework.TestContext, istioCtl istioctl.Instance, workDir string, iopFile string) {
	scopes.CI.Infof(fmt.Sprintf("=== Checking istio installation by operator controller with iop file: %s===\n", iopFile))
	s, err := image.SettingsFromCommandLine()
	if err != nil {
		t.Fatal(err)
	}
	iop, err := ioutil.ReadFile(iopFile)
	if err != nil {
		t.Fatalf("failed to read iop file: %v", err)
	}
	metadataYAML := `
metadata:
  name: test-istiocontrolplane
  namespace: istio-system
`
	iopcr, err := util.OverlayYAML(string(iop), metadataYAML)
	if err != nil {
		t.Fatalf("failed to overlay iop with metadata: %v", err)
	}
	iopCRFile := filepath.Join(workDir, "iop_cr.yaml")
	if err := ioutil.WriteFile(iopCRFile, []byte(iopcr), os.ModePerm); err != nil {
		t.Fatalf("failed to write iop cr file: %v", err)
	}
	initCmd := []string{
		"operator", "init",
		"--wait",
		"-f", iopCRFile,
		"--hub=" + s.Hub,
		"--tag=" + s.Tag,
	}
	// install istio using operator controller
	istioCtl.InvokeOrFail(t, initCmd)
	cs := ctx.Environment().(*kube.Environment).KubeClusters[0]

	// takes time for reconciliation to be done
	scopes.CI.Infof("waiting for reconciliation to be done")
	if err := checkInstallStatus(cs); err != nil {
		t.Fatalf("IstioOperator status not healthy: %v", err)
	}
	if _, err := cs.CheckPodsAreReady(cs.NewPodFetch(IstioNamespace)); err != nil {
		t.Fatalf("pods are not ready: %v", err)
	}

	if err := compareInClusterAndGeneratedResources(t, istioCtl, iopFile, cs); err != nil {
		t.Fatalf("in cluster resources does not match with the generated ones: %v", err)
	}
	scopes.CI.Infof("=== Succeeded ===")
}

func compareInClusterAndGeneratedResources(t *testing.T, istioCtl istioctl.Instance,
	iopFile string, cs kube.Cluster) error {
	// get manifests by running `manifest generate`
	generateCmd := []string{
		"manifest", "generate",
		"-f", iopFile,
	}
	genManifests := istioCtl.InvokeOrFail(t, generateCmd)
	genK8SObjects, err := object.ParseK8sObjectsFromYAMLManifest(genManifests)
	if err != nil {
		return fmt.Errorf("failed to parse generated manifest: %v", err)
	}
	efgvr := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "envoyfilters",
	}
	var errors util.Errors
	for _, genK8SObject := range genK8SObjects {
		kind := genK8SObject.Kind
		ns := genK8SObject.Namespace
		name := genK8SObject.Name
		log.Infof("checking kind: %s, namespace: %s, name: %s", kind, ns, name)
		switch kind {
		case "Service":
			if _, err := cs.GetService(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected service: %s from cluster", name))
			}
		case "ServiceAccount":
			if _, err := cs.GetServiceAccount(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected serviceAccount: %s from cluster", name))
			}
		case "Deployment":
			if _, err := cs.GetDeployment(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected deployment: %s from cluster", name))
			}
		case "ConfigMap":
			if _, err := cs.GetConfigMap(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected configMap: %s from cluster", name))
			}
		case "ValidatingWebhookConfiguration":
			if exist := cs.ValidatingWebhookConfigurationExists(name); !exist {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected ValidatingWebhookConfiguration: %s from cluster", name))
			}
		case "MutatingWebhookConfiguration":
			if exist := cs.MutatingWebhookConfigurationExists(name); !exist {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected MutatingWebhookConfiguration: %s from cluster", name))
			}
		case "CustomResourceDefinition":
			if _, err := cs.GetCustomResourceDefinition(name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected CustomResourceDefinition: %s from cluster", name))
			}
		case "ClusterRole":
			if _, err := cs.GetClusterRole(name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected ClusterRole: %s from cluster", name))
			}
		case "ClusterRoleBinding":
			if _, err := cs.GetClusterRoleBinding(name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected ClusterRoleBinding: %s from cluster", name))
			}
		case "EnvoyFilter":
			if _, err := cs.GetUnstructured(efgvr, ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected Envoyfilter: %s from cluster", name))
			}
		case "PodDisruptionBudget":
			if _, err := cs.GetPodDisruptionBudget(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected CustomResourceDefinition: %s from cluster", name))
			}
		case "HorizontalPodAutoscaler":
			if _, err := cs.GetHorizontalPodAutoscaler(ns, name); err != nil {
				errors = util.AppendErr(errors,
					fmt.Errorf("failed to get expected CustomResourceDefinition: %s from cluster", name))
			}
		}
	}
	return errors.ToError()
}
