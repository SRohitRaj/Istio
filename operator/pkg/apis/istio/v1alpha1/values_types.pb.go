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


syntax = 'proto3';

import "google/protobuf/duration.proto";
import "github.com/gogo/protobuf/protobuf/google/protobuf/wrappers.proto";
import "gogoproto/gogo.proto";

package v1alpha1;

// Package-wide variables from generator "generated".
option go_package = "v1alpha1";

option (gogoproto.marshaler_all) = false;
option (gogoproto.unmarshaler_all) = false;
option (gogoproto.sizer_all) = false;

// ArchConfig specifies the pod scheduling target architecture(amd64, ppc64le, s390x) for all the Istio control plane components.
message ArchConfig {
  // Sets pod scheduling weight for amd64 arch
  uint32 amd64 = 1;

  // Sets pod scheduling weight for ppc64le arch.
  uint32 ppc64le = 2;

  // Sets pod scheduling weight for s390x arch.
  uint32 s390x = 3;
}

// Configuration for CNI.
message CNIConfig {
  // Controls whether CNI is enabled.
  google.protobuf.BoolValue enabled = 1;

  string hub = 2;

  TypeInterface tag = 3;

  string image = 4;

  string pullPolicy = 5;

  string cniBinDir = 6;

  string cniConfDir = 7;

  string cniConfFileName = 8;

  repeated string excludeNamespaces = 9;

  TypeMapStringInterface podAnnotations = 10 [deprecated=true];

  string psp_cluster_role = 11;

  string logLevel = 12;

  CNIRepairConfig repair = 13;

  google.protobuf.BoolValue chained = 14;
}

message CNIRepairConfig {
  // Controls whether repair behavior is enabled.
  google.protobuf.BoolValue enabled = 1;

  string hub = 2;

  TypeInterface tag = 3;

  string image = 4;

  // Controls whether various repair behaviors are enabled.
  bool labelPods = 5;

  string createEvents = 6 [deprecated=true];

  bool deletePods = 7;

  string brokenPodLabelKey = 8;

  string brokenPodLabelValue = 9;

  string initContainerName = 10;
}

// Configuration for CPU target utilization for HorizontalPodAutoscaler target.
message CPUTargetUtilizationConfig {
  // K8s utilization setting for HorizontalPodAutoscaler target.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  int32 targetAverageUtilization = 1;
}

// Mirrors Resources for unmarshaling.
message Resources {
  map<string, string> limits = 1;
  map<string, string> requests = 2;
}

// Configuration for Core DNS.
message CoreDNSConfig {
  // Controls whether CoreDNS is enabled.
  google.protobuf.BoolValue enabled = 1;

  // Image for Core DNS.
  string coreDNSImage = 2;

  string coreDNSTag = 3;

  string coreDNSPluginImage = 4;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 5 [deprecated=true];

  // Number of replicas for Core DNS.
  uint32 replicaCount = 6 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 7 [deprecated=true];

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 8 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge = 9 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 10 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 11 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 12 [deprecated=true];

  TypeSliceOfMapStringInterface tolerations = 13 [deprecated=true];

  // Controls whether auto scaling with a HorizontalPodAutoscaler is enabled.
  google.protobuf.BoolValue autoscaleEnabled = 14;

  // maxReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMax = 15;

  // minReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMin = 16;

  // K8s utilization setting for HorizontalPodAutoscaler target.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 17 [deprecated=true];
}

// DefaultPodDisruptionBudgetConfig specifies the default pod disruption budget configuration.
//
// See https://kubernetes.io/docs/concepts/workloads/pods/disruptions/
message DefaultPodDisruptionBudgetConfig {
  // Controls whether a PodDisruptionBudget with a default minAvailable value of 1 is created for each deployment.
  google.protobuf.BoolValue enabled = 1;
}

// DefaultResourcesConfig specifies the default k8s resources settings for all Istio control plane components.
message DefaultResourcesConfig {
  // k8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  ResourcesRequestsConfig requests = 1;
}

// Configuration for an egress gateway.
message EgressGatewayConfig {
  // Controls whether auto scaling with a HorizontalPodAutoscaler is enabled.
  google.protobuf.BoolValue autoscaleEnabled = 1;

  // maxReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMax = 2;

  // minReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMin = 3;

  string connectTimeout = 4;

  // K8s utilization setting for HorizontalPodAutoscaler target.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 5 [deprecated=true];

  google.protobuf.Duration drainDuration = 6;

  // Controls whether an egress gateway is enabled.
  google.protobuf.BoolValue enabled = 7;

  // Environment variables passed to the proxy container.
  TypeMapStringInterface env = 8;

  GatewayLabelsConfig labels = 9;

  string name = 25;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 10 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 11 [deprecated=true];

  // Pod anti-affinity label selector.
  //
  // Specify the pod anti-affinity that allows you to constrain which nodes
  // your pod is eligible to be scheduled based on labels on pods that are
  // already running on the node rather than based on labels on nodes.
  // There are currently two types of anti-affinity:
  //    "requiredDuringSchedulingIgnoredDuringExecution"
  //    "preferredDuringSchedulingIgnoredDuringExecution"
  // which denote “hard” vs. “soft” requirements, you can define your values
  // in "podAntiAffinityLabelSelector" and "podAntiAffinityTermLabelSelector"
  // correspondingly.
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
  //
  // Examples:
  // podAntiAffinityLabelSelector:
  //  - key: security
  //    operator: In
  //    values: S1,S2
  //    topologyKey: "kubernetes.io/hostname"
  //  This pod anti-affinity rule says that the pod requires not to be scheduled
  //  onto a node if that node is already running a pod with label having key
  //  “security” and value “S1”.
  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 12 [deprecated=true];

  // See PodAntiAffinityLabelSelector.
  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 13 [deprecated=true];

  // Ports Configuration for the egress gateway service.
  repeated PortsConfig ports = 14;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 15 [deprecated=true];

  // Config for secret volume mounts.
  repeated SecretVolume secretVolumes = 16;

  // Annotations to add to the egress gateway service.
  TypeMapStringInterface serviceAnnotations = 17;

  // Service type.
  //
  // See https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
  string type = 18;

  // Enables cross-cluster access using SNI matching.
  ZeroVPNConfig zvpn = 19;

  TypeSliceOfMapStringInterface tolerations = 20 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge = 21 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 22 [deprecated=true];

  TypeSliceOfMapStringInterface configVolumes = 23;

  TypeSliceOfMapStringInterface additionalContainers = 24;

  google.protobuf.BoolValue runAsRoot = 26;

  // Next available 27.
}

// GatewayLabelsConfig is a set of Configuration for gateway labels.
message GatewayLabelsConfig {
  string app = 1;
  string istio = 2;
}

// Configuration for gateways.
message GatewaysConfig {
  // Configuration for an egress gateway.
  EgressGatewayConfig istio_egressgateway = 1;

  // Controls whether any gateways are enabled.
  google.protobuf.BoolValue enabled = 2;

  // Configuration for an ingress gateway.
  IngressGatewayConfig istio_ingressgateway = 4;
}

// Global Configuration for Istio components.
message GlobalConfig {
  // Specifies pod scheduling arch(amd64, ppc64le, s390x) and weight as follows:
  //   0 - Never scheduled
  //   1 - Least preferred
  //   2 - No preference
  //   3 - Most preferred
  ArchConfig arch = 1;

  // Specifies the namespace for the configuration and validation component.
  string configNamespace = 2;

  string configRootNamespace = 50;

  // Controls whether the server-side validation is enabled.
  google.protobuf.BoolValue configValidation = 3;

  // Controls whether the MTLS for communication between the control plane components is enabled.
  google.protobuf.BoolValue controlPlaneSecurityEnabled = 4;

  repeated string defaultConfigVisibilitySettings = 52;
  // Default k8s node selector for all the Istio control plane components
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface defaultNodeSelector = 6 [deprecated=true];

  // Specifies the default pod disruption budget configuration.
  DefaultPodDisruptionBudgetConfig defaultPodDisruptionBudget = 7 [deprecated=true];

  // Default k8s resources settings for all Istio control plane components.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  DefaultResourcesConfig defaultResources = 9 [deprecated=true];

  TypeSliceOfMapStringInterface defaultTolerations = 55 [deprecated=true];

  // Controls whether the helm test templates are enabled.
  google.protobuf.BoolValue enableHelmTest = 10;

  // Controls whether the distributed tracing for the applications is enabled.
  //
  // See https://opentracing.io/docs/overview/what-is-tracing/
  google.protobuf.BoolValue enableTracing = 11;

  // Specifies the docker hub for Istio images.
  string hub = 12;

  // Specifies the image pull policy for the Istio images. one of Always, Never, IfNotPresent.
  // Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated.
  //
  // More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
  string imagePullPolicy = 13;
  // ImagePullPolicy             v1.PullPolicy                 `json:"imagePullPolicy,omitempty"`

  repeated string imagePullSecrets = 37;

  // Specifies the default namespace for the Istio control plane components.
  string istioNamespace = 14;

  // Specifies the global locality load balancing settings.
  // Locality-weighted load balancing allows administrators to control the distribution of traffic to
  // endpoints based on the localities of where the traffic originates and where it will terminate.
  // Please set either failover or distribute configuration but not both.
  //
  // localityLbSetting:
  //   distribute:
  //   - from: "us-central1/*"
  //     to:
  //       "us-central1/*": 80
  //       "us-central2/*": 20
  //
  // localityLbSetting:
  //   failover:
  //   - from: us-east
  //     to: eu-west
  //   - from: us-west
  //     to: us-east
  TypeMapStringInterface localityLbSetting = 15;
  // 	LocalityLbSetting map[string]interface{}    `json:"localityLbSetting"`

  google.protobuf.BoolValue logAsJson = 36;

  // Specifies the global logging level settings for the Istio control plane components.
  GlobalLoggingConfig logging = 17;

  // Specifies the Configuration for Istio mesh expansion to bare metal.
  MeshExpansionConfig meshExpansion = 18;

  string meshID = 53;

  // Configure the mesh networks to be used by the Split Horizon EDS.
  //
  // The following example defines two networks with different endpoints association methods.
  // For `network1` all endpoints that their IP belongs to the provided CIDR range will be
  // mapped to network1. The gateway for this network example is specified by its public IP
  // address and port.
  // The second network, `network2`, in this example is defined differently with all endpoints
  // retrieved through the specified Multi-Cluster registry being mapped to network2. The
  // gateway is also defined differently with the name of the gateway service on the remote
  // cluster. The public IP for the gateway will be determined from that remote service (only
  // LoadBalancer gateway service type is currently supported, for a NodePort type gateway service,
  // it still need to be configured manually).
  //
  // meshNetworks:
  //   network1:
  //     endpoints:
  //     - fromCidr: "192.168.0.1/24"
  //     gateways:
  //     - address: 1.1.1.1
  //       port: 80
  //   network2:
  //     endpoints:
  //     - fromRegistry: reg1
  //     gateways:
  //     - registryServiceName: istio-ingressgateway.istio-system.svc.cluster.local
  //       port: 443
  //
  TypeMapStringInterface meshNetworks = 19;

  // Specifies the monitor port number for all Istio control plane components.
  uint32 monitoringPort = 20;

  // Specifies the Configuration for Istio mesh across multiple clusters through Istio gateways.
  MultiClusterConfig multiCluster = 22;

  string network = 39;

  // Custom DNS config for the pod to resolve names of services in other
  // clusters. Use this to add additional search domains, and other settings.
  // see https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#dns-config
  // This does not apply to gateway pods as they typically need a different
  // set of DNS settings than the normal application pods (e.g. in multicluster scenarios).
  repeated string podDNSSearchNamespaces = 43;

  google.protobuf.BoolValue omitSidecarInjectorConfigMap = 38;

  // Controls whether to restrict the applications namespace the controller manages;
  // If set it to false, the controller watches all namespaces.
  google.protobuf.BoolValue oneNamespace = 23;

  google.protobuf.BoolValue operatorManageWebhooks = 41;

  // Controls whether to allow traffic in cases when the mixer policy service cannot be reached.
  google.protobuf.BoolValue policyCheckFailOpen = 25;

  // Specifies the namespace for the policy component.
  string policyNamespace = 26;

  // Specifies the k8s priorityClassName for the istio control plane components.
  //
  // See https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/#priorityclass
  string priorityClassName = 27 [deprecated=true];

  string prometheusNamespace = 51;

  // Specifies how proxies are configured within Istio.
  ProxyConfig proxy = 28;

  // Specifies the Configuration for proxy_init container which sets the pods' networking to intercept the inbound/outbound traffic.
  ProxyInitConfig proxyInit = 29;

  // Specifies the Configuration for the SecretDiscoveryService instead of using K8S secrets to mount the certificates.
  SDSConfig sds = 30;

  // Specifies the tag for the Istio docker images.
  TypeInterface tag = 31;

  // Specifies the namespace for the telemetry component.
  string telemetryNamespace = 32;

  // Specifies the Configuration for each of the supported tracers.
  TracerConfig tracer = 33;

  // Specifies the trust domain that corresponds to the root cert of CA.
  string trustDomain = 34;

  // The trust domain aliases represent the aliases of trustDomain.
  repeated string trustDomainAliases = 42;

  // Controls whether to use of Mesh Configuration Protocol to distribute configuration.
  google.protobuf.BoolValue useMCP = 35;

  // Settings for remote cluster.
  // Controls whether to use the Istio remote control plane
  google.protobuf.BoolValue istioRemote = 44;

  google.protobuf.BoolValue createRemoteSvcEndpoints = 45;
  // If set, a selector-less service and endpoint for istio-pilot are created with the remotePilotAddress IP,
  // which ensures the istio-pilot. is DNS resolvable in the remote cluster.
  google.protobuf.BoolValue remotePilotCreateSvcEndpoint = 46;
  // Specifies the Istio control plane’s policy Pod IP address or remote cluster DNS resolvable hostname.
  string remotePolicyAddress = 47;
  // Specifies the Istio control plane’s pilot Pod IP address or remote cluster DNS resolvable hostname.
  string remotePilotAddress = 48;
  // Specifies the Istio control plane’s telemetry Pod IP address or remote cluster DNS resolvable hostname
  string remoteTelemetryAddress = 49;

  // Specifies the configution of istiod
  IstiodConfig istiod = 54;

  // Configure the Pilot certificate provider.
  // Currently, two providers are supported: "kubernetes" and "citadel".
  string pilotCertProvider = 56;

  // Configure the policy for validating JWT.
  // Currently, two options are supported: "third-party-jwt" and "first-party-jwt".
  string jwtPolicy = 57;

  // Specifies the configuration for Security Token Service.
  STSConfig sts = 58;

  // Configures the revision this control plane is a part of
  string revision = 59;

  // Controls whether the in-cluster MTLS key and certs are loaded from the secret volume mounts.
  google.protobuf.BoolValue mountMtlsCerts = 60;

  // The address of the CA for CSR.
  string caAddress = 61;

  // Controls whether one central istiod is enabled.
  google.protobuf.BoolValue centralIstiod = 62;
  // The next available key is 63
}

// Configuration for Security Token Service (STS) server.
//
// See https://tools.ietf.org/html/draft-ietf-oauth-token-exchange-16
message STSConfig {
  uint32 servicePort = 1;
}

message IstiodConfig {
  // If enabled, istiod will perform config analysis
  google.protobuf.BoolValue enableAnalysis = 2;
}

// GlobalLoggingConfig specifies the global logging level settings for the Istio control plane components.
message GlobalLoggingConfig {
  // Comma-separated minimum per-scope logging level of messages to output, in the form of <scope>:<level>,<scope>:<level>
  // The control plane has different scopes depending on component, but can configure default log level across all components
  // If empty, default scope and level will be used as configured in code
  string level = 1;
}

// Configuration for an ingress gateway.
message IngressGatewayConfig {
  // Controls whether auto scaling with a HorizontalPodAutoscaler is enabled.
  google.protobuf.BoolValue autoscaleEnabled = 1;

  // maxReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMax = 2;

  // minReplicas setting for HorizontalPodAutoscaler.
  uint32 autoscaleMin = 3;

  string connectTimeout = 4;

  // K8s utilization setting for HorizontalPodAutoscaler target.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 5 [deprecated=true];

  google.protobuf.BoolValue customService = 6;

  string debug = 7;

  string domain = 8;

  google.protobuf.Duration drainDuration = 9;

  // Controls whether an ingress gateway is enabled.
  google.protobuf.BoolValue enabled = 10;

  // Environment variables passed to the proxy container.
  TypeMapStringInterface env = 11;

  repeated string externalIPs = 12;

  google.protobuf.BoolValue k8sIngress = 13;

  google.protobuf.BoolValue k8sIngressHttps = 14;

  GatewayLabelsConfig labels = 15;

  string loadBalancerIP = 16;

  repeated string loadBalancerSourceRanges = 17;

  repeated PortsConfig meshExpansionPorts = 18;

  string name = 44;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 19 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 20 [deprecated=true];

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 21 [deprecated=true];

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 22 [deprecated=true];

  // Port Configuration for the ingress gateway.
  repeated PortsConfig ports = 23;

  // Number of replicas for the ingress gateway Deployment.
  uint32 replicaCount = 24 [deprecated=true];

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  TypeMapStringInterface resources = 25 [deprecated=true];

  // Secret Discovery Service (SDS) Configuration for ingress gateway.
  IngressGatewaySdsConfig sds = 26;

  // Config for secret volume mounts.
  repeated SecretVolume secretVolumes = 27;

  // Annotations to add to the egress gateway service.
  TypeMapStringInterface serviceAnnotations = 28;

  // Service type.
  //
  // See https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
  string type = 29;

  // Enables cross-cluster access using SNI matching.
  IngressGatewayZvpnConfig zvpn = 30;

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge = 31 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 32 [deprecated=true];

  // Ports to explicitly check for readiness
  string applicationPorts = 33;

  string externalTrafficPolicy = 34;

  TypeSliceOfMapStringInterface tolerations = 35 [deprecated=true];

  TypeSliceOfMapStringInterface ingressPorts = 36;

  TypeSliceOfMapStringInterface additionalContainers = 37;

  TypeSliceOfMapStringInterface configVolumes = 38;

  google.protobuf.BoolValue certificates = 39;

  google.protobuf.BoolValue tls = 40;

  TypeMapStringInterface telemetry_addon_gateways = 41;

  TypeSliceOfMapStringInterface hosts = 42;

  string telemetry_domain_name = 43;

  google.protobuf.BoolValue runAsRoot = 45;

  // Next available 46.
}

// Secret Discovery Service (SDS) Configuration for ingress gateway.
message IngressGatewaySdsConfig {
  // If true, ingress gateway fetches credentials from SDS server to handle TLS connections.
  google.protobuf.BoolValue enabled = 1;

  // SDS server that watches kubernetes secrets and provisions credentials to ingress gateway.
  // This server runs in the same pod as ingress gateway.
  string image = 2;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 3 [deprecated=true];
}

// IngressGatewayZvpnConfig enables cross-cluster access using SNI matching.
message IngressGatewayZvpnConfig {
  // Controls whether ZeroVPN is enabled.
  google.protobuf.BoolValue enabled = 1;

  string suffix = 2;
}

// Configuration for Kubernetes environment adapter in mixer.
message KubernetesEnvMixerAdapterConfig {
  // Enables the Kubernetes env adapter in Mixer.
  //
  // See: https://istio.io/docs/reference/config/policy-and-telemetry/adapters/kubernetesenv/
  google.protobuf.BoolValue enabled = 1;
}

// Configuration for when mixer starts rejecting grpc requests.
message LoadSheddingConfig {
  string latencyThreshold = 1;
  mode mode = 2;
}

// Throttling behavior for mixer.
enum mode {
  // Removes throttling behavior for mixer.
  disabled = 0;
  // Enables an advisory mode for throttling behavior for mixer.
  log_only = 1;
  // Turn on throttling behavior for mixer.
  enforce = 2;
}

// Configuration for Istio mesh expansion to bare metal.
message MeshExpansionConfig {
  // Exposes Pilot and Citadel mTLS on the ingress gateway.
  google.protobuf.BoolValue enabled = 1;

  // Exposes Pilot and Citadel mTLS and the plain text Pilot ports on an internal gateway.
  google.protobuf.BoolValue useILB = 2;
}

// Configuration for Mixer Telemetry adapters.
message MixerTelemetryAdaptersConfig {
  // Configuration for Kubernetes environment adapter in mixer.
  KubernetesEnvMixerAdapterConfig kubernetesenv = 1;

  // Configuration for Prometheus adapter in mixer.
  PrometheusMixerAdapterConfig prometheus = 2;

  // Configuration for stdio adapter in mixer, recommended for debug usage only.
  StdioMixerAdapterConfig stdio = 3;

  //
  StackdriverMixerAdapterConfig stackdriver = 4;

  // Sets the --useAdapterCRDs mixer startup argument.
  google.protobuf.BoolValue useAdapterCRDs = 5;
}

// Configuration for Mixer Policy adapters.
message MixerPolicyAdaptersConfig {
  // Configuration for Kubernetes environment adapter in mixer.
  KubernetesEnvMixerAdapterConfig kubernetesenv = 1;

  // Configuration for Prometheus adapter in mixer.
  PrometheusMixerAdapterConfig prometheus = 2;

  // Configuration for stdio adapter in mixer, recommended for debug usage only.
  StdioMixerAdapterConfig stdio = 3;

  //
  StackdriverMixerAdapterConfig stackdriver = 4;

  // Sets the --useAdapterCRDs mixer startup argument.
  google.protobuf.BoolValue useAdapterCRDs = 5;
}

// Configuration for Mixer.
message MixerConfig {
  // MixerPolicyConfig is set of configurations for Mixer Policy
  MixerPolicyConfig policy = 1;

  // MixerTelemetryConfig is set of configurations for Mixer Telemetry
  MixerTelemetryConfig telemetry = 2;

  // Configuration for different mixer adapters.
  MixerTelemetryAdaptersConfig adapters = 3;
}

// Configuration for Mixer Policy.
message MixerPolicyConfig {
  // Controls whether a HorizontalPodAutoscaler is installed for Mixer Policy.
  google.protobuf.BoolValue autoscaleEnabled = 1;

  // Maximum number of replicas in the HorizontalPodAutoscaler for Mixer Policy.
  uint32 autoscaleMax = 2;

  // Minimum number of replicas in the HorizontalPodAutoscaler for Mixer Policy.
  uint32 autoscaleMin = 3;

  // Target CPU utilization used in HorizontalPodAutoscaler.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 4 [deprecated=true];

  // Controls whether Mixer Policy is enabled
  google.protobuf.BoolValue enabled = 5;

  // Image name used for Mixer Policy.
  //
  // This can be set either to image name if hub is also set, or can be set to the full hub:name string.
  //
  // Examples: custom-mixer, docker.io/someuser:custom-mixer
  string image = 6;

  // K8s annotations to attach to mixer policy deployment
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 7 [deprecated=true];

  // Number of replicas in the Mixer Policy Deployment
  uint32 replicaCount = 8 [deprecated=true];

  // Configuration for different mixer adapters.
  MixerPolicyAdaptersConfig adapters = 9;

  // Controls whether to enable the sticky session setting when choosing backend pods.
  google.protobuf.BoolValue sessionAffinityEnabled = 10;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 11 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge = 12 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 13 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 14 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 15 [deprecated=true];

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 16 [deprecated=true];

  // Environment variables passed to the Mixer container.
  //
  // Examples:
  // env:
  //   ENV_VAR_1: value1
  //   ENV_VAR_2: value2
  TypeMapStringInterface env = 17;

  TypeSliceOfMapStringInterface tolerations = 18 [deprecated=true];

  string hub = 19;

  TypeInterface tag = 20;
}

// Configuration for Mixer Telemetry.
message MixerTelemetryConfig {
  // Controls whether a HorizontalPodAutoscaler is installed for Mixer Telemetry.
  google.protobuf.BoolValue autoscaleEnabled = 2;

  // Maximum number of replicas in the HorizontalPodAutoscaler for Mixer Telemetry.
  uint32 autoscaleMax = 3;

  // Minimum number of replicas in the HorizontalPodAutoscaler for Mixer Telemetry.
  uint32 autoscaleMin = 4;

  // Target CPU utilization used in HorizontalPodAutoscaler.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 5 [deprecated=true];

  // Controls whether Mixer Telemetry is enabled.
  google.protobuf.BoolValue enabled = 6;

  // Environment variables passed to the Mixer container.
  //
  // Examples:
  // env:
  //   ENV_VAR_1: value1
  //   ENV_VAR_2: value2
  TypeMapStringInterface env = 7;

  // Image name used for Mixer Telemetry.
  //
  // This can be set either to image name if hub is also set, or can be set to the full hub:name string.
  //
  // Examples: custom-mixer, docker.io/someuser:custom-mixer
  string image = 8;

  // LoadSheddingConfig configs when mixer starts rejecting grpc requests.
  LoadSheddingConfig loadshedding = 9;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 10 [deprecated=true];

  // K8s annotations to attach to mixer telemetry deployment
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 11 [deprecated=true];

  // Number of replicas in the Mixer Telemetry Deployment.
  uint32 replicaCount = 12 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge  = 15 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 16 [deprecated=true];

  // Controls whether to use of Mesh Configuration Protocol to distribute configuration.
  google.protobuf.BoolValue useMCP = 17;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 13 [deprecated=true];

  // Controls whether to enable the sticky session setting when choosing backend pods.
  google.protobuf.BoolValue sessionAffinityEnabled = 14;

  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 20 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 21 [deprecated=true];

  TypeSliceOfMapStringInterface tolerations = 22 [deprecated=true];

  string hub = 23;

  TypeInterface tag = 24;
}

// MultiClusterConfig specifies the Configuration for Istio mesh across multiple clusters through the istio gateways.
message MultiClusterConfig {
  // Enables the connection between two kubernetes clusters via their respective ingressgateway services.
  // Use if the pods in each cluster cannot directly talk to one another.
  google.protobuf.BoolValue enabled = 1;

  string clusterName = 2;

  string globalDomainSuffix = 3;

  google.protobuf.BoolValue includeEnvoyFilter = 4;
}

// OutboundTrafficPolicyConfig controls the default behavior of the sidecar for handling outbound traffic from the application.
message OutboundTrafficPolicyConfig {
  // Specifies the sidecar's default behavior when handling outbound traffic from the application.
  enum Mode {
    // Outbound traffic to unknown destinations will be allowed, in case there are no services or ServiceEntries for the destination port
    ALLOW_ANY = 0;
    // Restrict outbound traffic to services defined in the service registry as well as those defined through ServiceEntries
    REGISTRY_ONLY = 1;
  }
  Mode mode = 2;
}

// Configuration for Pilot.
message PilotConfig {
  // Controls whether Pilot is enabled.
  google.protobuf.BoolValue enabled = 1;

  // Controls whether a HorizontalPodAutoscaler is installed for Pilot.
  google.protobuf.BoolValue autoscaleEnabled = 2;

  // Minimum number of replicas in the HorizontalPodAutoscaler for Pilot.
  uint32 autoscaleMin = 3;

  // Maximum number of replicas in the HorizontalPodAutoscaler for Pilot.
  uint32 autoscaleMax = 4;

  // Number of replicas in the Pilot Deployment.
  uint32 replicaCount = 5 [deprecated=true];

  // Image name used for Pilot.
  //
  // This can be set either to image name if hub is also set, or can be set to the full hub:name string.
  //
  // Examples: custom-pilot, docker.io/someuser:custom-pilot
  string image = 6;

  // Controls whether a sidecar proxy is installed in the Pilot pod.
  //
  // Setting to true installs a proxy in the Pilot pod, used primarily for collecting Pilot telemetry.
  google.protobuf.BoolValue sidecar = 7;

  // Trace sampling fraction.
  //
  // Used to set the fraction of time that traces are sampled. Higher values are more accurate but add CPU overhead.
  //
  // Allowed values: 0.0 to 1.0
  double traceSampling = 8;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 9 [deprecated=true];

  // Namespace that the configuration management feature is installed into, if different from Pilot namespace.
  string configNamespace = 10;

  // Target CPU utilization used in HorizontalPodAutoscaler.
  //
  // See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
  CPUTargetUtilizationConfig cpu = 11 [deprecated=true];

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 12 [deprecated=true];

  // Maximum duration that a sidecar can be connected to a pilot.
  //
  // This setting balances out load across pilot instances, but adds some resource overhead.
  //
  // Examples: 300s, 30m, 1h
  google.protobuf.Duration keepaliveMaxServerConnectionAge = 13;

  // Labels that are added to Pilot pods.
  //
  // See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
  TypeMapStringInterface deploymentLabels = 14;

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 16 [deprecated=true];

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 17 [deprecated=true];

  // Configuration settings passed to Pilot as a ConfigMap.
  //
  // This controls whether the mesh config map, generated from values.yaml is generated.
  // If false, pilot wil use default values or user-supplied values, in that order of preference.
  google.protobuf.BoolValue configMap = 18;

  // Controls whether Pilot is configured through the Mesh Control Protocol (MCP).
  //
  // If set to true, Pilot requires an MCP server (like Galley) to be installed.
  google.protobuf.BoolValue useMCP = 20;

  // Environment variables passed to the Pilot container.
  //
  // Examples:
  // env:
  //   ENV_VAR_1: value1
  //   ENV_VAR_2: value2
  TypeMapStringInterface env = 21;

  // Controls whether Istio policy is applied to Pilot.
  PilotPolicyConfig policy = 22;

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxSurge  = 24 [deprecated=true];

  // K8s rolling update strategy
  TypeIntOrStringForPB rollingMaxUnavailable = 25 [deprecated=true];

  //
  TypeSliceOfMapStringInterface tolerations = 26 [deprecated=true];

  TypeSliceOfMapStringInterface appNamespaces = 27;

  // if protocol sniffing is enabled for outbound
  google.protobuf.BoolValue enableProtocolSniffingForOutbound = 28;
  // if protocol sniffing is enabled for inbound
  google.protobuf.BoolValue enableProtocolSniffingForInbound = 29;

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 30 [deprecated=true];

  // ConfigSource describes a source of configuration data for networking
  // rules, and other Istio configuration artifacts. Multiple data sources
  // can be configured for a single control plane.
  PilotConfigSource configSource = 31;

  string jwksResolverExtraRootCA = 32;

  TypeSliceString plugins = 33;

  string hub = 34;

  TypeInterface tag = 35;
}

// Controls legacy k8s ingress. Only one pilot profile should enable ingress support.
message PilotIngressConfig {
  // Sets the type ingress service for Pilot.
  //
  // If empty, node-port is assumed.
  //
  // Allowed values: node-port, istio-ingressgateway, ingress
  string ingressService = 1;

  ingressControllerMode ingressControllerMode = 2;
  // If mode is STRICT, this value must be set on "kubernetes.io/ingress.class" annotation to activate.
  string ingressClass = 3;
}

// Mode for the ingress controller.
enum ingressControllerMode {
  // Unspecified Istio ingress controller.
  UNSPECIFIED = 0;
  // Selects all Ingress resources, with or without Istio annotation.
  DEFAULT = 1;
  // Selects only resources with istio annotation.
  STRICT = 2;
  // No ingress or sync.
  OFF = 3;
}

// Controls whether Istio policy is applied to Pilot.
message PilotPolicyConfig {
  // Controls whether Istio policy is applied to Pilot.
  google.protobuf.BoolValue enabled = 1;
}

// Controls telemetry configuration
message TelemetryConfig {
  // Controls whether telemetry is exported for Pilot.
  google.protobuf.BoolValue enabled = 1;

  // Use telemetry v1.
  TelemetryV1Config v1 = 2;

  // Use telemetry v2.
  TelemetryV2Config v2 = 3;
}

// Controls whether pilot will configure telemetry v1.
message TelemetryV1Config {
  // Controls whether pilot will configure telemetry v1.
  google.protobuf.BoolValue enabled = 1;
}

// Controls whether pilot will configure telemetry v2.
message TelemetryV2Config {
  // Controls whether pilot will configure telemetry v2.
  google.protobuf.BoolValue enabled = 1;

  TelemetryV2MetadataExchangeConfig metadata_exchange = 4;

  TelemetryV2PrometheusConfig prometheus = 2;

  TelemetryV2StackDriverConfig stackdriver = 3;

  TelemetryV2AccessLogPolicyFilterConfig access_log_policy = 5;
}

message TelemetryV2MetadataExchangeConfig {
  // Controls whether enabled WebAssembly runtime for metadata exchange filter.
  google.protobuf.BoolValue wasmEnabled = 2;
}

// Conrols telemetry v2 prometheus settings.
message TelemetryV2PrometheusConfig {
  // Controls whether stats envoyfilter would be enabled or not.
  google.protobuf.BoolValue enabled = 1;

  // Controls whether enabled WebAssembly runtime for stats filter.
  google.protobuf.BoolValue wasmEnabled = 2;

  message ConfigOverride {
    // Overrides default gateway telemetry v2 configuration.
    TypeMapStringInterface gateway = 1;

    // Overrides default inbound sidecar telemetry v2 configuration.
    TypeMapStringInterface inboundSidecar = 2;

    // Overrides default outbound sidecar telemetry v2 configuration.
    TypeMapStringInterface outboundSidecar = 3;
  }

  // Overrides default telemetry v2 filter configuration.
  ConfigOverride config_override = 3;
}

// Conrols telemetry v2 stackdriver settings.
message TelemetryV2StackDriverConfig {
  google.protobuf.BoolValue enabled = 1;

  google.protobuf.BoolValue logging = 2;

  google.protobuf.BoolValue monitoring = 3;

  google.protobuf.BoolValue topology = 4;

  google.protobuf.BoolValue disableOutbound = 6;

  TypeMapStringInterface configOverride = 5;
}

// Conrols telemetry v2 access log policy filter settings.
message TelemetryV2AccessLogPolicyFilterConfig {
  google.protobuf.BoolValue enabled = 1;

  google.protobuf.Duration logWindowDuration = 2;
}

// PilotConfigSource describes information about a configuration store inside a
// mesh. A single control plane instance can interact with one or more data
// sources.
message PilotConfigSource {
  // Describes the source of configuration, if nothing is specified default is MCP.
  repeated string subscribedResources = 1;
}

// Configuration for a port.
message PortsConfig {
  // Port name.
  string name = 1;

  // Port number.
  int32 port = 2;

  // NodePort number.
  int32 nodePort = 3;

  // Target port number.
  int32 targetPort = 4;
}

// Configuration for Prometheus.
message PrometheusConfig {
  google.protobuf.BoolValue createPrometheusResource = 1;

  google.protobuf.BoolValue enabled = 2;

  uint32 replicaCount = 3 [deprecated=true];

  string hub = 4;

  TypeInterface tag = 5;

  string retention = 6;

  TypeMapStringInterface nodeSelector = 7 [deprecated=true];
  // GOSTRUCT: NodeSelector             map[string]interface{}    `json:"nodeSelector,omitempty"`

  google.protobuf.Duration scrapeInterval = 8;

  string contextPath = 9;

  PrometheusServiceConfig service = 11;

  PrometheusSecurityConfig security = 12;

  TypeSliceOfMapStringInterface tolerations = 13 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 14 [deprecated=true];

  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 15 [deprecated=true];

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 16 [deprecated=true];

  TypeSliceOfMapStringInterface datasources = 17 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 18 [deprecated=true];

  string image = 19 [deprecated=true];

  // Configure whether provisions a certificate to Prometheus through Istio Agent.
  // When this option is set as true, a sidecar is deployed along Prometheus to
  // provision a certificate through Istio Agent to Prometheus. The provisioned certificate
  // is shared with Prometheus through mounted files.
  // When this option is set as false, this certificate provisioning mechanism is disabled.
  google.protobuf.BoolValue  provisionPrometheusCert = 20;
}

// Configuration for Prometheus adapter in mixer.
message PrometheusMixerAdapterConfig {
  // Enables the Prometheus adapter in Mixer.
  google.protobuf.BoolValue enabled = 1;

  // Sets the duration after which Prometheus registry purges a metric.
  //
  // See: https://istio.io/docs/reference/config/policy-and-telemetry/adapters/prometheus/#Params
  google.protobuf.Duration metricsExpiryDuration = 2;
}

// Configuration for Prometheus adapter security.
message PrometheusSecurityConfig {
  // Controls whether Prometheus security is enabled.
  google.protobuf.BoolValue enabled = 1;
}

// Configuration for Prometheus adapter service.
message PrometheusServiceConfig {
  TypeMapStringInterface annotations = 1;
  PrometheusServiceNodePortConfig nodePort = 2;
}

// Configuration for Prometheus Service NodePort.
message PrometheusServiceNodePortConfig {
  // Controls whether Prometheus NodePort config is enabled.
  google.protobuf.BoolValue enabled = 1;

  uint32 port = 2;
}

// Configures the access log for sidecar to JSON or TEXT
enum accessLogEncoding {
  JSON = 0;
  TEXT = 1;
}

// Configuration for Proxy.
message ProxyConfig {
  string autoInject = 4;

  // Domain for the cluster, default: "cluster.local".
  //
  // K8s allows this to be customized, see https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/
  string clusterDomain = 5;

  // Per Component log level for proxy, applies to gateways and sidecars.
  //
  // If a component level is not set, then the global "logLevel" will be used. If left empty, "misc:error" is used.
  string componentLogLevel = 6;

  // Enables core dumps for newly injected sidecars.
  //
  // If set, newly injected sidecars will have core dumps enabled.
  google.protobuf.BoolValue enableCoreDump = 9;

  // Specifies the Istio ingress ports not to capture.
  string excludeInboundPorts = 12;

  // Lists the excluded IP ranges of Istio egress traffic that the sidecar captures.
  string excludeIPRanges = 13;

  // Image name or path for the proxy, default: "proxyv2".
  //
  // If registry or tag are not specified, global.hub and global.tag are used.
  //
  // Examples: my-proxy (uses global.hub/tag), docker.io/myrepo/my-proxy:v1.0.0
  string image = 14;

  // Lists the IP ranges of Istio egress traffic that the sidecar captures.
  //
  // Example: "172.30.0.0/16,172.20.0.0/16"
  // This would only capture egress traffic on those two IP Ranges, all other outbound traffic would # be allowed by the sidecar."
  string includeIPRanges = 16;

  // Log level for proxy, applies to gateways and sidecars. If left empty, "warning" is used. Expected values are: trace\|debug\|info\|warning\|error\|critical\|off
  string logLevel = 18;

  // Enables privileged securityContext for the istio-proxy container.
  //
  // See https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
  google.protobuf.BoolValue privileged = 19;

  // Sets the initial delay for readiness probes in seconds.
  uint32 readinessInitialDelaySeconds = 20;

  // Sets the interval between readiness probes in seconds.
  uint32 readinessPeriodSeconds = 21;

  // Sets the number of successive failed probes before indicating readiness failure.
  uint32 readinessFailureThreshold = 22;

  // Default port used for the Pilot agent's health checks.
  uint32 statusPort = 23;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 24 [deprecated=true];

  tracer tracer = 25;

  string excludeOutboundPorts = 28;

  TypeMapStringInterface lifecycle = 36;

  // Controls if sidecar is injected at the front of the container list and blocks the start of the other containers until the proxy is ready
  google.protobuf.BoolValue holdApplicationUntilProxyStarts = 37;
}

// Specifies which tracer to use.
enum tracer {
  zipkin = 0;
  lightstep = 1;
  datadog = 2;
  stackdriver = 3;
}

// Configuration for proxy_init container which sets the pods' networking to intercept the inbound/outbound traffic.
message ProxyInitConfig {
  // Specifies the image for the proxy_init container.
  string image = 1;
  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 5 [deprecated=true];
}

// Configuration for K8s resource requests.
message ResourcesRequestsConfig {
  string cpu = 1;

  string memory = 2;
}

// Configuration for the SecretDiscoveryService instead of using K8S secrets to mount the certificates.
message SDSConfig {
  TypeMapStringInterface token = 5 [deprecated=true];
}

// Configuration for secret volume mounts.
//
// See https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets.
message SecretVolume {
  string mountPath = 1;

  string name = 2;

  string secretName = 3;
}

// ServiceConfig is described in istio.io documentation.
message ServiceConfig {
  TypeMapStringInterface annotations = 1;

  uint32 externalPort = 2;

  string name = 3;

  string type = 18;
}

// SidecarInjectorConfig is described in istio.io documentation.
message SidecarInjectorConfig {
  // Enables sidecar auto-injection in namespaces by default.
  google.protobuf.BoolValue enableNamespacesByDefault = 2;

  // Instructs Istio to not inject the sidecar on those pods, based on labels that are present in those pods.
  //
  // Annotations in the pods have higher precedence than the label selectors.
  // Order of evaluation: Pod Annotations → NeverInjectSelector → AlwaysInjectSelector → Default Policy.
  // See https://istio.io/docs/setup/kubernetes/additional-setup/sidecar-injection/#more-control-adding-exceptions
  TypeSliceOfMapStringInterface neverInjectSelector = 11;

  // See NeverInjectSelector.
  TypeSliceOfMapStringInterface alwaysInjectSelector = 12;

  //  If true, webhook or istioctl injector will rewrite PodSpec for liveness health check to redirect request to sidecar. This makes liveness check work even when mTLS is enabled.
  google.protobuf.BoolValue rewriteAppHTTPProbe = 16;

  string injectLabel = 18;

  // injectedAnnotations are additional annotations that will be added to the pod spec after injection
  // This is primarily to support PSP annotations.
  TypeMapStringInterface injectedAnnotations = 19;

  // Enable objectSelector to filter out pods with no need for sidecar before calling istio-sidecar-injector.
  TypeMapStringInterface objectSelector = 21;

  // Configure the injection url for sidecar injector webhook
  string injectionURL = 22;
}

// Configuration for stdio adapter in mixer, recommended for debug usage only.
message StdioMixerAdapterConfig {
  // Enable stdio adapter to output logs and metrics to local machine.
  google.protobuf.BoolValue enabled = 1;

  // Whether to output a console-friendly or json-friendly format.
  google.protobuf.BoolValue outputAsJson = 2;
}

// Configuration for stackdriver adapter in mixer.
message StackdriverMixerAdapterConfig {
  google.protobuf.BoolValue enabled = 1;

  StackdriverAuthConfig auth = 2;

  StackdriverTracerConfig tracer = 3;

  StackdriverContextGraph contextGraph = 4;

  message EnabledConfig {
    google.protobuf.BoolValue enabled = 1;
  }

  EnabledConfig logging = 5;

  EnabledConfig metrics = 6;
}

message StackdriverAuthConfig {
  google.protobuf.BoolValue appCredentials = 1;

  string apiKey = 2;

  string serviceAccountPath = 3;
}

message StackdriverTracerConfig {
  google.protobuf.BoolValue enabled = 1;

  uint32 sampleProbability = 2;
}

message StackdriverContextGraph {
  google.protobuf.BoolValue enabled = 1;
}

// Configuration for each of the supported tracers.
message TracerConfig {
  // Configuration for the datadog tracing service.
  TracerDatadogConfig datadog = 1;

  // Configuration for the lightstep tracing service.
  TracerLightStepConfig lightstep = 2;

  // Configuration for the zipkin tracing service.
  TracerZipkinConfig zipkin = 3;

  // Configuration for the stackdriver tracing service.
  TracerStackdriverConfig stackdriver = 4;
}

// Configuration for the datadog tracing service.
message TracerDatadogConfig {
  // Address in host:port format for reporting trace data to the Datadog agent.
  string address = 1;
}

// Configuration for the lightstep tracing service.
message TracerLightStepConfig {
  // Sets the lightstep satellite pool address in host:port format for reporting trace data.
  string address = 1;

  // Sets the lightstep access token.
  string accessToken = 2;
}

// Configuration for the zipkin tracing service.
message TracerZipkinConfig {
  // Address of zipkin instance in host:port format for reporting trace data.
  //
  // Example: <zipkin-collector-service>.<zipkin-collector-namespace>:941
  string address = 1;
}

// Configuration for the stackdriver tracing service.
message TracerStackdriverConfig {
  // enables trace output to stdout.
  google.protobuf.BoolValue debug = 1;

  // The global default max number of attributes per span.
  uint32 maxNumberOfAttributes = 2;

  // The global default max number of annotation events per span.
  uint32 maxNumberOfAnnotations = 3;

  // The global default max number of message events per span.
  uint32 maxNumberOfMessageEvents = 4;
}

// Configurations for different tracing system to be installed.
message TracingConfig {
  // Enables tracing systems installation.
  google.protobuf.BoolValue enabled = 1;

  // Defines Configuration for addon Jaeger tracing.
  TracingJaegerConfig jaeger = 3;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 4 [deprecated=true];

  // Configures which tracing system to be installed.
  string provider = 5;

  // Controls K8s service for addon tracing components.
  ServiceConfig service = 6;

  // Defines Configuration for addon Zipkin tracing.
  TracingZipkinConfig zipkin = 7;

  TracingOpencensusConfig opencensus = 8;

  string contextPath = 9;

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 13 [deprecated=true];

  // See EgressGatewayConfig.
  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 14 [deprecated=true];

  TypeSliceOfMapStringInterface tolerations = 15 [deprecated=true];
}

message TracingOpencensusConfig {
  // Image hub for Opencensus tracing deployment.
  string hub = 1;

  // Image tag for Opencensus tracing deployment.
  TypeInterface tag = 2;

  TracingOpencensusExportersConfig exporters  = 3;
  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  TypeMapStringInterface resources = 5 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 6 [deprecated=true];
}

message TracingOpencensusExportersConfig {
  TypeMapStringInterface stackdriver = 1;
}

// Configuration for addon Jaeger tracing.
message TracingJaegerConfig {
  // Image hub for Jaeger tracing deployment.
  string hub = 1;

  // Image tag for Jaeger tracing deployment.
  TypeInterface tag = 2;

  string image = 10;

  // Configures Jaeger in-memory storage setting.
  TracingJaegerMemoryConfig memory = 3;

  string spanStorageType = 4;

  google.protobuf.BoolValue persist = 5;

  string storageClassName = 6;

  string accessMode = 7;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  TypeMapStringInterface resources = 8 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 9 [deprecated=true];
}

// Configuration for Jaeger in-memory storage setting.
message TracingJaegerMemoryConfig {
  // Set limit of the amount of traces stored in memory for Jaeger
  uint32 max_traces = 1;
}

// Configuration for Zipkin.
message TracingZipkinConfig {
  // Image hub for Zipkin tracing deployment.
  string hub = 1;

  // Image tag for Zipkin tracing deployment.
  TypeInterface tag = 2;

  string image = 9;

  // InitialDelaySeconds of readiness probe for Zipkin deployment
  uint32 probeStartupDelay = 3;

  // InitialDelaySeconds of liveness probe for Zipkin deployment
  uint32 livenessProbeStartupDelay = 11;

  // Container port for Zipkin deployment
  uint32 queryPort = 4;

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 5 [deprecated=true];

  // Configure java heap opts for Zipkin deployment
  uint32 javaOptsHeap = 6;

  // Configures number of max spans to keep in Zipkin memory storage.
  //
  // Example: A safe estimate is 1K of memory per span (each span with 2 annotations + 1 binary annotation), plus 100 MB for a safety buffer
  uint32 maxSpans = 7;

  // Configures GC values of JAVA_OPTS for Zipkin deployment
  TracingZipkinNodeConfig node = 8;

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 10 [deprecated=true];
}

// Configuration for GC values of JAVA_OPTS for Zipkin deployment
message TracingZipkinNodeConfig {
  // Configures -XX:ConcGCThreads value of JAVA_OPTS for Zipkin deployment
  uint32 cpus = 1;
}

message KialiSecurityConfig {
  google.protobuf.BoolValue enabled = 1;

  string cert_file = 2;

  string private_key_file = 3;
}

message KialiServiceConfig {
  TypeMapStringInterface annotations = 1;

  // Service type.
  //
  // See https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
  string type = 18;
}

message KialiDashboardConfig {
  string secretName = 1;

  string usernameKey = 2;

  string passphraseKey = 3;

  google.protobuf.BoolValue viewOnlyMode = 4;

  string grafanaURL = 5;

  string jaegerURL = 6;

  TypeMapStringInterface auth = 7;

  string grafanaInClusterURL = 8;

  string jaegerInClusterURL = 9;
}

// Configuration for Kiali addon.
message KialiConfig {
  google.protobuf.BoolValue enabled = 1;

  google.protobuf.BoolValue createDemoSecret = 2;

  // Image hub for kiali deployment.
  string hub = 3;

  // Image tag for kiali deployment.
  TypeInterface tag = 4;

  // Number of replicas for Kiali.
  uint32 replicaCount = 5 [deprecated=true];

  string prometheusNamespace = 6;

  KialiSecurityConfig security = 7;

  KialiDashboardConfig dashboard = 8;

  string contextPath = 15;

  // K8s node selector.
  //
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
  TypeMapStringInterface nodeSelector = 10 [deprecated=true];

  // K8s annotations for pods.
  //
  // See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
  TypeMapStringInterface podAnnotations = 11 [deprecated=true];

  // Pod anti-affinity label selector.
  //
  // Specify the pod anti-affinity that allows you to constrain which nodes
  // your pod is eligible to be scheduled based on labels on pods that are
  // already running on the node rather than based on labels on nodes.
  // There are currently two types of anti-affinity:
  //    "requiredDuringSchedulingIgnoredDuringExecution"
  //    "preferredDuringSchedulingIgnoredDuringExecution"
  // which denote “hard” vs. “soft” requirements, you can define your values
  // in "podAntiAffinityLabelSelector" and "podAntiAffinityTermLabelSelector"
  // correspondingly.
  // See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
  //
  // Examples:
  // podAntiAffinityLabelSelector:
  //  - key: security
  //    operator: In
  //    values: S1,S2
  //    topologyKey: "kubernetes.io/hostname"
  //  This pod anti-affinity rule says that the pod requires not to be scheduled
  //  onto a node if that node is already running a pod with label having key
  //  “security” and value “S1”.
  TypeSliceOfMapStringInterface podAntiAffinityLabelSelector = 12 [deprecated=true];

  // See PodAntiAffinityLabelSelector.
  TypeSliceOfMapStringInterface podAntiAffinityTermLabelSelector = 13 [deprecated=true];

  TypeSliceOfMapStringInterface tolerations = 14 [deprecated=true];

  string image = 16 [deprecated=true];

  // K8s resources settings.
  //
  // See https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container
  Resources resources = 17 [deprecated=true];

  string prometheusAddr = 18 [deprecated=true];

  KialiServiceConfig service = 19;
}

message BaseConfig {
  // For Helm2 use, adds the CRDs to templates.
  google.protobuf.BoolValue enableCRDTemplates = 1;

  // URL to use for validating webhook.
  string validationURL = 2;
}

message IstiodRemoteConfig {
  // URL to use for sidecar injector webhook.
  string injectionURL = 1;
}

message Values {
  CNIConfig cni = 2;

  CoreDNSConfig istiocoredns = 3;

  GatewaysConfig gateways = 5;

  GlobalConfig global = 6;

  TypeMapStringInterface grafana = 7;

  MixerConfig mixer = 8;

  PilotConfig pilot = 10;

  // Controls whether telemetry is exported for Pilot.
  TelemetryConfig telemetry = 23;

  PrometheusConfig prometheus = 11;

  SidecarInjectorConfig sidecarInjectorWebhook = 13;

  TracingConfig tracing = 14;

  KialiConfig kiali = 15;

  // Deprecated.
  string version = 16;

  google.protobuf.BoolValue clusterResources = 17;

  // TODO: populate these.
  TypeMapStringInterface prometheusOperator = 18;
  CNIConfig istio_cni = 19;

  google.protobuf.BoolValue kustomize = 20;

  string revision = 21;

  // TODO can this import the real mesh config API?
  TypeInterface meshConfig = 36;

  BaseConfig base = 37;

  IstiodRemoteConfig istiodRemote = 38;
}

// GOTYPE: map[string]interface{}
message TypeMapStringInterface {}

// GOTYPE: []map[string]interface{}
message TypeSliceOfMapStringInterface {}

// GOTYPE: *IntOrStringForPB
message TypeIntOrStringForPB {}

// GOTYPE: []string
message TypeSliceString {}

// ZeroVPNConfig enables cross-cluster access using SNI matching.
message ZeroVPNConfig {
  // Controls whether ZeroVPN is enabled.
  google.protobuf.BoolValue enabled = 1;

  string suffix = 2;
}

// GOTYPE: interface{}
message TypeInterface {}

