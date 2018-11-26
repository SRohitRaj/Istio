#!/bin/bash

# Copyright 2018 Istio Authors
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

set -e
set -u
set -o pipefail

REPORT_PATH=${GOPATH}/out/codecov
BASELINE_PATH=${GOPATH}/out/codecov_baseline

# First run codecov from current workspace (PR)
OUT_DIR="${REPORT_PATH}" MAXPROCS=${MAXPROCS:-} ./bin/codecov.sh || echo "Some tests have failed"

if [[ -n "${CIRCLE_PR_NUMBER:-}" ]]; then
  TMP_GITHUB_TOKEN=$(mktemp /tmp/XXXXX.github)
  openssl version
  openssl aes-256-cbc -d -in .circleci/accounts/istio-github.enc -out "${TMP_GITHUB_TOKEN}" -k "${GCS_BUCKET_TOKEN}" -md sha256

  # Backup codecov.sh since the base SHA may not have this copy.
  TMP_CODECOV_SH=$(mktemp /tmp/XXXXX.codecov)
  cp ./bin/codecov.sh "${TMP_CODECOV_SH}"

  go get -u istio.io/test-infra/toolbox/githubctl
  BASE_SHA=$("${GOPATH}"/bin/githubctl --token_file="${TMP_GITHUB_TOKEN}" --op=getBaseSHA --repo=istio --pr_num="${CIRCLE_PR_NUMBER}")
  git checkout "${BASE_SHA}"

  cp "${TMP_CODECOV_SH}" ./bin/codecov.sh


  OUT_DIR="${BASELINE_PATH}" MAXPROCS="${MAXPROCS:-}" ./bin/codecov.sh || echo "Some tests have failed"

  go get -u istio.io/test-infra/toolbox/pkg_check
  "${GOPATH}"/bin/pkg_check  --report_file="${REPORT_PATH}/codecov.report" --baseline_file="${BASELINE_PATH}/codecov.report"
fi

