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

package mesh

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"istio.io/istio/operator/pkg/cache"
	"istio.io/istio/operator/pkg/helmreconciler"
	"istio.io/istio/operator/pkg/manifest"
	"istio.io/istio/operator/pkg/name"
	"istio.io/istio/operator/pkg/util/clog"
	"istio.io/istio/operator/pkg/util/progress"
	proxyinfo "istio.io/istio/pkg/proxy"
	"istio.io/pkg/log"
)

type uninstallArgs struct {
	// kubeConfigPath is the path to kube config file.
	kubeConfigPath string
	// context is the cluster context in the kube config.
	context string
	// skipConfirmation determines whether the user is prompted for confirmation.
	// If set to true, the user is not prompted and a Yes response is assumed in all cases.
	skipConfirmation bool
	// force proceeds even if there are validation errors
	force bool
	// revision is the Istio control plane revision the command targets.
	revision string
	// istioNamespace is the target namespace of istio control plane.
	istioNamespace string
	// filename is the path of input IstioOperator CR.
	filename string
	// set is a string with element format "path=value" where path is an IstioOperator path and the value is a
	// value to set the node at that path to.
	set []string
	// manifestsPath is a path to a charts and profiles directory in the local filesystem, or URL with a release tgz.
	manifestsPath string
}

func addUninstallFlags(cmd *cobra.Command, args *uninstallArgs) {
	cmd.PersistentFlags().StringVarP(&args.kubeConfigPath, "kubeconfig", "c", "", "Path to kube config.")
	cmd.PersistentFlags().StringVar(&args.context, "context", "", "The name of the kubeconfig context to use.")
	cmd.PersistentFlags().BoolVarP(&args.skipConfirmation, "skip-confirmation", "y", false, skipConfirmationFlagHelpStr)
	cmd.PersistentFlags().BoolVar(&args.force, "force", false, "Proceed even with validation errors.")
	cmd.PersistentFlags().StringVarP(&args.revision, "revision", "r", "", revisionFlagHelpStr)
	cmd.PersistentFlags().StringVar(&args.istioNamespace, "istioNamespace", istioDefaultNamespace,
		"The namespace of Istio Control Plane.")
	cmd.PersistentFlags().StringVarP(&args.filename, "fileName", "f", "",
		"The filename of the IstioOperator CR.")
	cmd.PersistentFlags().StringVarP(&args.manifestsPath, "manifests", "d", "", ManifestsFlagHelpStr)
	cmd.PersistentFlags().StringArrayVarP(&args.set, "set", "s", nil, setFlagHelpStr)
}

func UninstallCmd(logOpts *log.Options) *cobra.Command {
	rootArgs := &rootArgs{}
	uiArgs := &uninstallArgs{}
	uicmd := &cobra.Command{
		Use:   "uninstall --revision foo",
		Short: "uninstall the control plane by revision",
		Long:  "The uninstall command uninstall the control plane by revision or IstioOperator CR",
		Args: func(cmd *cobra.Command, args []string) error {
			if uiArgs.revision == "" && uiArgs.filename == "" {
				return fmt.Errorf("at least one of the --revision or --filename flags must be set")
			}
			if len(args) > 0 {
				return fmt.Errorf("istioctl uninstall does not take arguments")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstall(cmd, rootArgs, uiArgs, logOpts)
		}}
	addUninstallFlags(uicmd, uiArgs)
	return uicmd
}

// uninstall uninstalls control plane by either pruning by target revision or deleting specified manifests.
func uninstall(cmd *cobra.Command, rootArgs *rootArgs, uiArgs *uninstallArgs, logOpts *log.Options) error {
	l := clog.NewConsoleLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), installerScope)
	if err := configLogs(logOpts); err != nil {
		return fmt.Errorf("could not configure logs: %s", err)
	}
	restConfig, _, clt, err := K8sConfig(uiArgs.kubeConfigPath, uiArgs.context)
	if err != nil {
		return err
	}

	cache.FlushObjectCaches()
	opts := &helmreconciler.Options{DryRun: rootArgs.dryRun, Log: l, ProgressLog: progress.NewLog()}
	h, err := helmreconciler.NewHelmReconciler(clt, restConfig, nil, opts)
	if err != nil {
		return fmt.Errorf("failed to create reconciler: %v", err)
	}

	// If only revision flag is set, we would prune resources by the revision label.
	// Otherwise we would merge the revision flag and the filename flag and delete resources by generated manifests.
	if uiArgs.filename == "" {
		objectsList, tp, err := h.GetPrunedResourcesByRevision(uiArgs.revision)
		if err != nil {
			return err
		}
		if err := preCheckWarnings(cmd, rootArgs, uiArgs, uiArgs.revision, tp, l); err != nil {
			return fmt.Errorf("failed to do preuninstall check: %v", err)
		}

		if err := h.DeleteControlPlaneByRevision(uiArgs.revision, objectsList); err != nil {
			return fmt.Errorf("failed to prune control plane resources by revision: %v", err)
		}
	}
	manifestMap, iops, err := manifest.GenManifests([]string{uiArgs.filename}, applyFlagAliases(uiArgs.set, uiArgs.manifestsPath, uiArgs.revision), uiArgs.force, restConfig, l)
	if err != nil {
		return err
	}
	if iops == nil {
		return fmt.Errorf("istioOperatorSpec is nil")
	}
	if err := preCheckWarnings(cmd, rootArgs, uiArgs, iops.Revision, nil, l); err != nil {
		return fmt.Errorf("failed to do preuninstall check: %v", err)
	}
	cpManifests := manifestMap[name.PilotComponentName]
	if err := h.DeleteControlPlaneByManifests(strings.Join(cpManifests, "---"), iops.Revision); err != nil {
		return fmt.Errorf("failed to delete control plane by manifests: %v", err)
	}
	opts.ProgressLog.SetState(progress.StateComplete)
	return nil
}

// preCheckWarnings checks possible breaking changes and issue warnings to users, it checks the following:
// 1. checks proxies still pointing to the target control plane revision.
// 2. lists to be pruned resources if user uninstall by --revision flag.
func preCheckWarnings(cmd *cobra.Command, rootArgs *rootArgs, uiArgs *uninstallArgs,
	rev string, resourcesList []string, l *clog.ConsoleLogger) error {
	pids, err := proxyinfo.GetIDsFromProxyInfo(uiArgs.kubeConfigPath, uiArgs.context, rev, uiArgs.istioNamespace)
	if err != nil {
		return err
	}
	needConfirmation := false
	var message string
	if resourcesList != nil {
		needConfirmation = true
		message += fmt.Sprintf("To be pruned resources from the cluster: %s\n",
			strings.Join(resourcesList, ","))
	}
	if len(pids) != 0 {
		needConfirmation = true
		message += fmt.Sprintf("There are still proxies pointing to the control plane revision %s:\n%s."+
			" If you proceed with the uninstall, these proxies will become detached from any control plane"+
			" and will not function correctly.. ",
			uiArgs.revision, strings.Join(pids, " \n"))
	}
	if uiArgs.skipConfirmation || rootArgs.dryRun {
		l.LogAndPrint(message)
		return nil
	}
	message += "Proceed? (y/N)"
	if needConfirmation && !confirm(message, cmd.OutOrStdout()) {
		cmd.Print("Cancelled.\n")
		os.Exit(1)
	}
	return nil
}
