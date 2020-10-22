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

package cmd

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_labels "k8s.io/apimachinery/pkg/labels"
	apimachinery_schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"istio.io/istio/operator/cmd/mesh"
	operator_istio "istio.io/istio/operator/pkg/apis/istio"
	iopv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	"istio.io/istio/operator/pkg/manifest"
	"istio.io/istio/operator/pkg/util"
	"istio.io/istio/operator/pkg/util/clog"
	"istio.io/istio/pkg/config"
)

type istiodListArgs struct {
	// manifestsPath is a path to a charts and profiles directory in the local filesystem, or URL with a release tgz.
	manifestsPath string
}

var (
	istioOperatorGVR = apimachinery_schema.GroupVersionResource{
		Group:    iopv1alpha1.SchemeGroupVersion.Group,
		Version:  iopv1alpha1.SchemeGroupVersion.Version,
		Resource: "istiooperators",
	}
)

func istiodListCmd() *cobra.Command {
	kubeConfigFlags := &genericclioptions.ConfigFlags{
		Context:    strPtr(""),
		Namespace:  strPtr(""),
		KubeConfig: strPtr(""),
	}
	listArgs := &istiodListArgs{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists Istio control planes",
		Long:  "The list subcommand displays installed Istio control planes to the console",
		// nolint: lll
		Example: `  # List Istio installations
  istioctl experimental istiod list

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := clog.NewConsoleLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), scope)
			return istiodList(cmd.OutOrStdout(), listArgs, kubeConfigFlags, l)
		}}

	listCmd.PersistentFlags().StringVarP(&listArgs.manifestsPath, "manifests", "d", "", mesh.ManifestsFlagHelpStr)
	return listCmd
}

func istiodList(writer io.Writer, listArgs *istiodListArgs, restClientGetter genericclioptions.RESTClientGetter, l clog.Logger) error {
	restConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	iops, err := getIOPs(restConfig)
	if err != nil {
		return err
	}
	return printIOPs(writer, iops, listArgs.manifestsPath, restConfig, l)
}

func printIOPs(writer io.Writer, iops []*iopv1alpha1.IstioOperator, manifestsPath string, restConfig *rest.Config, l clog.Logger) error {
	if len(iops) == 0 {
		_, err := fmt.Fprintf(writer, "No IstioOperators present.\n")
		return err
	}

	sort.Slice(iops, func(i, j int) bool {
		return iops[i].Spec.Revision < iops[j].Spec.Revision
	})

	w := new(tabwriter.Writer).Init(writer, 0, 8, 1, ' ', 0)
	fmt.Fprintln(w, "VERSION\tREVISION\tPROFILE\tCUSTOMIZATIONS\tPODS\tAGE")
	for _, iop := range iops {
		podCount, deploymentAge, err := getControlPlaneDeployment(iop, manifestsPath, restConfig)
		if err != nil {
			podCount = "<error>"
			deploymentAge = "<error>"
		}

		diffs, err := getDiffs(iop, manifestsPath, effectiveProfile(iop.Spec.Profile), l)
		if err != nil {
			return err
		}
		for i, diff := range diffs {
			if i == 0 {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					iop.Spec.Tag,
					renderWithDefault(iop.Spec.Revision, "master"),
					renderWithDefault(iop.Spec.Profile, "default"),
					diff, podCount, deploymentAge)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					"",
					"",
					"",
					diff, "", "")
			}
		}
	}
	return w.Flush()
}

func getIOPs(restConfig *rest.Config) ([]*iopv1alpha1.IstioOperator, error) {
	client, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	ul, err := client.
		Resource(istioOperatorGVR).
		Namespace(istioNamespace).
		List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	iops := []*iopv1alpha1.IstioOperator{}
	for _, un := range ul.Items {
		un.SetCreationTimestamp(meta_v1.Time{}) // UnmarshalIstioOperator chokes on these
		by := util.ToYAML(un.Object)
		iop, err := operator_istio.UnmarshalIstioOperator(by, true)
		if err != nil {
			return nil, err
		}
		iops = append(iops, iop)
	}
	return iops, nil
}

func getDiffs(installed *iopv1alpha1.IstioOperator, manifestsPath, profile string, l clog.Logger) ([]string, error) {
	setFlags := []string{"profile=" + profile}
	if manifestsPath != "" {
		setFlags = append(setFlags, fmt.Sprintf("installPackagePath=%s", manifestsPath))
	}

	_, base, err := manifest.GenerateConfig([]string{}, setFlags, true, nil, l)
	if err != nil {
		return []string{}, err
	}

	if err != nil {
		return []string{}, err
	}
	mapInstalled, err := config.ToMap(installed.Spec)
	if err != nil {
		return []string{}, err
	}
	mapBase, err := config.ToMap(base.Spec)
	if err != nil {
		return []string{}, err
	}

	return diffIOPs(mapInstalled, mapBase)
}

func diffIOPs(installed, base map[string]interface{}) ([]string, error) {
	setflags, err := diffWalk("", "", installed, base)
	if err != nil {
		return []string{}, err
	}
	sort.Strings(setflags)
	return setflags, nil
}

func diffWalk(path, separator string, obj interface{}, orig interface{}) ([]string, error) {
	switch v := obj.(type) {
	case map[string]interface{}:
		accum := make([]string, 0)
		typedOrig, ok := orig.(map[string]interface{})
		if ok {
			for key, vv := range v {
				childwalk, err := diffWalk(fmt.Sprintf("%s%s%s", path, separator, pathComponent(key)), ".", vv, typedOrig[key])
				if err != nil {
					return accum, err
				}
				accum = append(accum, childwalk...)
			}
		}
		return accum, nil
	case []interface{}:
		accum := make([]string, 0)
		typedOrig, ok := orig.([]interface{})
		if ok {
			for idx, vv := range v {
				indexwalk, err := diffWalk(fmt.Sprintf("%s[%d]", path, idx), ".", vv, typedOrig[idx])
				if err != nil {
					return accum, err
				}
				accum = append(accum, indexwalk...)
			}
		}
		return accum, nil
	case string:
		if v != orig && orig != nil {
			return []string{fmt.Sprintf("%s=%q", path, v)}, nil
		}
	default:
		if v != orig && orig != nil {
			return []string{fmt.Sprintf("%s=%v", path, v)}, nil
		}
	}
	return []string{}, nil
}

func renderWithDefault(s, def string) string {
	if s != "" {
		return s
	}
	return fmt.Sprintf("<%s>", def)
}

func effectiveProfile(profile string) string {
	if profile != "" {
		return profile
	}
	return "default"
}

func getControlPlaneDeployment(iop *iopv1alpha1.IstioOperator, manifestsPath string, restConfig *rest.Config) (string, string, error) {
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", "", err
	}
	deploymentName := getDeploymentName(iop, manifestsPath)
	if err != nil {
		return "", "", err
	}
	deployment, err := client.AppsV1().
		Deployments(istioNamespace).
		Get(context.TODO(), deploymentName, meta_v1.GetOptions{})
	if err != nil {
		return "", "", err
	}
	pods, err := client.CoreV1().
		Pods(istioNamespace).
		List(context.TODO(), meta_v1.ListOptions{
			LabelSelector: k8s_labels.Set(deployment.Spec.Selector.MatchLabels).AsSelector().String(),
		})
	if err != nil {
		return "", "", err
	}

	podCount := strconv.Itoa(len(pods.Items))
	deploymentAge := translateTimestampSince(deployment.CreationTimestamp)
	return podCount, deploymentAge, nil
}

// Human-readable age.  (This is from kubectl pkg/describe/describe.go)
func translateTimestampSince(timestamp meta_v1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

func getDeploymentName(iop *iopv1alpha1.IstioOperator, manifestsPath string) string {
	// It would be difficult but perhaps better to render the manifest, extract the PilotComponentName manifest,
	// parse the yaml, and scrape it for the deployment name.  Instead we take a shortcut, as we know
	// the naming convention.
	_ = manifestsPath

	if iop.Spec.Revision == "" {
		return "istiod"
	}
	return fmt.Sprintf("istiod-%s", iop.Spec.Revision)
}

func strPtr(val string) *string {
	return &val
}

func pathComponent(component string) string {
	if !strings.Contains(component, util.PathSeparator) {
		return component
	}
	return strings.ReplaceAll(component, util.PathSeparator, util.EscapedPathSeparator)
}
