#!/usr/bin/env bash

# Copyright Istio Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

WD=$(dirname "$0")
WD=$(cd "$WD"; pwd)

set -eux

# This script sets up the plain text rendered deployments for addons
# See samples/addons/README.md for more information

ADDONS="${WD}/../../samples/addons"
DASHBOARDS="${WD}/dashboards"
mkdir -p "${ADDONS}"
TMP=$(mktemp -d)
LOKI_VERSION=${LOKI_VERSION:-"2.9.9"}

# Set up kiali
{
helm3 template kiali-server \
  --namespace istio-system \
  --version 1.63.1 \
  --set deployment.image_version=v1.63 \
  --include-crds \
  --set nameOverride=kiali \
  --set fullnameOverride=kiali \
  kiali-server \
  --repo https://kiali.org/helm-charts \
  -f "${WD}/values-kiali.yaml"
} > "${ADDONS}/kiali.yaml"

# Set up prometheus
helm3 template prometheus prometheus \
  --namespace istio-system \
  --version 19.6.1 \
  --repo https://prometheus-community.github.io/helm-charts \
  -f "${WD}/values-prometheus.yaml" \
  > "${ADDONS}/prometheus.yaml"

function compressDashboard() {
  < "${DASHBOARDS}/$1" jq -c  > "${TMP}/$1"
}

# Set up grafana
{
  helm3 template grafana grafana \
    --namespace istio-system \
    --version 6.31.1 \
    --repo https://grafana.github.io/helm-charts \
    -f "${WD}/values-grafana.yaml"

  # Set up grafana dashboards. Split into 2 and compress to single line json to avoid Kubernetes size limits
  compressDashboard "pilot-dashboard.json"
  compressDashboard "istio-performance-dashboard.json"
  compressDashboard "istio-workload-dashboard.json"
  compressDashboard "istio-service-dashboard.json"
  compressDashboard "istio-mesh-dashboard.json"
  compressDashboard "istio-extension-dashboard.json"
  echo -e "\n---\n"
  kubectl create configmap -n istio-system istio-grafana-dashboards \
    --dry-run=client -oyaml \
    --from-file=pilot-dashboard.json="${TMP}/pilot-dashboard.json" \
    --from-file=istio-performance-dashboard.json="${TMP}/istio-performance-dashboard.json"

  echo -e "\n---\n"
  kubectl create configmap -n istio-system istio-services-grafana-dashboards \
    --dry-run=client -oyaml \
    --from-file=istio-workload-dashboard.json="${TMP}/istio-workload-dashboard.json" \
    --from-file=istio-service-dashboard.json="${TMP}/istio-service-dashboard.json" \
    --from-file=istio-mesh-dashboard.json="${TMP}/istio-mesh-dashboard.json" \
    --from-file=istio-extension-dashboard.json="${TMP}/istio-extension-dashboard.json"
} > "${ADDONS}/grafana.yaml"

# Set up loki
{
  helm3 template loki loki-stack \
    --namespace istio-system \
    --version "${LOKI_VERSION}" \
    --repo https://grafana.github.io/helm-charts \
    -f "${WD}/values-loki.yaml"
} > "${ADDONS}/loki.yaml"
