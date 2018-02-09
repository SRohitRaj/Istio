// Copyright 2017 Istio Authors
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

package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"istio.io/istio/pkg/collateral"
	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/version"
	"istio.io/istio/security/cmd/node_agent/na"
	"istio.io/istio/security/pkg/cmd"
)

const (
	// The default path/file of root cert.
	defaultRoot = "/etc/certs/root-cert.pem"
)

var (
	naConfig = na.NewConfig()

	rootCmd = &cobra.Command{
		Use:   "node_agent",
		Short: "Istio security per-node agent",

		Run: func(cmd *cobra.Command, args []string) {
			runNodeAgent()
		},
	}
)

func init() {
	rootCmd.AddCommand(version.CobraCommand())
	rootCmd.AddCommand(collateral.CobraCommand(rootCmd, &doc.GenManHeader{
		Title:   "Istio Node Agent",
		Section: "node_agent CLI",
		Manual:  "Istio Node Agent",
	}))

	flags := rootCmd.Flags()

	flags.StringVar(&naConfig.ServiceIdentityOrg, "org", "", "Organization for the cert")
	flags.DurationVar(&naConfig.WorkloadCertTTL, "workload-cert-ttl", 19*time.Hour,
		"The requested TTL for the workload")
	flags.IntVar(&naConfig.RSAKeySize, "key-size", 2048, "Size of generated private key")
	flags.StringVar(&naConfig.IstioCAAddress,
		"ca-address", "istio-ca:8060", "Istio CA address")
	flags.StringVar(&naConfig.Env, "env", "onprem", "Node Environment : onprem | gcp | aws")

	flags.StringVar(&naConfig.PlatformConfig.OnPremConfig.CertChainFile, "onprem-cert-chain",
		"/etc/certs/cert-chain.pem", "Node Agent identity cert file in on premise environment")
	flags.StringVar(&naConfig.PlatformConfig.OnPremConfig.KeyFile,
		"onprem-key", "/etc/certs/key.pem", "Node identity private key file in on premise environment")
	flags.StringVar(&naConfig.PlatformConfig.OnPremConfig.RootCACertFile, "onprem-root-cert",
		defaultRoot, "Root Certificate file in on premise environment")

	flags.StringVar(&naConfig.PlatformConfig.GcpConfig.RootCACertFile, "gcp-root-cert",
		defaultRoot, "Root Certificate file in GCP environment")
	flags.StringVar(&naConfig.PlatformConfig.GcpConfig.CAAddr, "gcp-ca-address",
		"istio-ca:8060", "Istio CA address in GCP environment")

	flags.StringVar(&naConfig.PlatformConfig.AwsConfig.RootCACertFile, "aws-root-cert",
		defaultRoot, "Root Certificate file in AWS environment")

	naConfig.LoggingOptions.AttachCobraFlags(rootCmd)
	cmd.InitializeFlags(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}

func runNodeAgent() {
	if err := log.Configure(naConfig.LoggingOptions); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
	nodeAgent, err := na.NewNodeAgent(naConfig)
	if err != nil {
		log.Errora(err)
		os.Exit(-1)
	}

	log.Infof("Starting Node Agent")
	if err := nodeAgent.Start(); err != nil {
		log.Errorf("Node agent terminated with error: %v.", err)
		os.Exit(-1)
	}
}
