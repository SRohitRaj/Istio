#!/bin/bash

# Copyright 2017 Istio Authors

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.


#######################################################
# e2e-suite triggered after istio/presubmit succeeded #
#######################################################

# Exit immediately for non zero status
set -e
# Check unset variables
set -u
# Print commands
set -x

PROJECT_NAME=istio-testing
ZONE=us-east4-c
CLUSTER_VERSION=1.7.4
MACHINE_TYPE=n1-standard-4
NUM_NODES=1
CLUSTER_NAME=e2e-yutongz-$(uuidgen | cut -c1-6)

CLUSTER_CREATED=false

delete_cluster () {
    if [ "${CLUSTER_CREATED}" = true ]; then
        gcloud container clusters delete ${CLUSTER_NAME} --zone ${ZONE} --project ${PROJECT_NAME} --quiet \
            || echo "Failed to delete cluster ${CLUSTER_CREATED}"
    fi
}
trap delete_cluster EXIT

if [ "${CI:-}" == 'bootstrap' ]; then
  # Test harness will checkout code to directory $GOPATH/src/github.com/istio/istio
  # but we depend on being at path $GOPATH/src/istio.io/istio for imports
  ln -sf ${GOPATH}/src/github.com/istio ${GOPATH}/src/istio.io
  cd ${GOPATH}/src/istio.io/istio

  # bootsrap upload all artifacts in _artifacts to the log bucket.
  ARTIFACTS_DIR=${ARTIFACTS_DIR:-"${GOPATH}/src/istio.io/istio/_artifacts"}
  LOG_HOST="stackdriver"
  PROJ_ID=${PROJECT_NAME}
  E2E_ARGS+=(--test_logs_path="${ARTIFACTS_DIR}" --log_provider=${LOG_HOST} --project_id=${PROJ_ID})
fi

if [ -f /home/bootstrap/.kube/config ]; then
  sudo chmod 666 /home/bootstrap/.kube/config
fi

gcloud container clusters create ${CLUSTER_NAME} --zone ${ZONE} --project ${PROJECT_NAME} --cluster-version ${CLUSTER_VERSION} \
  --machine-type ${MACHINE_TYPE} --num-nodes ${NUM_NODES} --enable-kubernetes-alpha --quiet \
  || { echo "Failed to create a new cluster"; exit 1; }
CLUSTER_CREATED=true

echo 'Running Integration Tests'
./tests/e2e.sh ${E2E_ARGS[@]:-} ${@}
