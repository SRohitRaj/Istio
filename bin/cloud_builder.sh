#!/bin/bash
# Copyright 2017 Istio Authors. All Rights Reserved.
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
#
################################################################################

set -o errexit
set -o nounset
set -o pipefail
set -x

OUTPUT_PATH=""
# The default for PROXY_PATH (which indicates where the proxy path is located
# relative to the istio repo) is based on repo manifest that places istio at:
# go/src/istio.io/istio
# and proxy at:
# src/proxy
PROXY_PATH="../../../../src/proxy"
TAG_NAME="0.0.0"

function usage() {
  echo "$0
    -o        path to store build artifacts
    -p        path to proxy repo (relative to istio repo, defaults to ${PROXY_PATH} ) 
    -t <tag>  tag to use (optional, defaults to ${TAG_NAME} )"
  exit 1
}

while getopts o:p:t: arg ; do
  case "${arg}" in
    o) OUTPUT_PATH="${OPTARG}";;
    p) PROXY_PATH="${PROXY_PATH}";;
    t) TAG_NAME="${OPTARG}";;
    *) usage;;
  esac
done

[[ -z "${OUTPUT_PATH}" ]] && usage
[[ -z "${PROXY_PATH}" ]] && usage

# switch to the root of the istio repo
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd $ROOT

export GOPATH="$(cd "$ROOT/../../.." && pwd)":${ROOT}/vendor
echo gopath is $GOPATH

# clean slate here
../../../../repo status

if [ ! -d "${PROXY_PATH}" ]; then
  echo "proxy dir not detected at ${PROXY_PATH}"
  usage
fi

# Proxy has some specific requirements for Bazel's
# config (plus it's nicely places bazel in batch
# mode) so this component gets built first.

pushd "${PROXY_PATH}"

# Use this file for Cloud Builder specific settings.
# This file sets RAM sizes and also specifies batch
# mode that should shutdown bazel after each call.
echo 'Setting bazel.rc'
cp tools/bazel.rc.cloudbuilder "${HOME}/.bazelrc"

./script/push-debian.sh -c opt -v "${TAG_NAME}" -o "${OUTPUT_PATH}"
popd

# clean slate here
../../../../repo status

# Pilot likes checking if the source tree is 'clean'
# but some other components like littering the tree
# so it's better to build pilot sooner than later.

pushd pilot
# Build istioctl binaries
touch platform/kube/config

# XXX seeing if this makes more consistent results for generated_files and lintconfig.json
# bazel build //pilot/...
bazel build //...
popd

# slate is TBD here
../../../../repo status

cp ./generated_files /output/generated_files.before
cp ./lintconfig.json /output/lintconfig.json.before

# bazel_to_go likes to run from dir with WORKSPACE file
./bin/bazel_to_go.py
# Remove doubly-vendorized k8s dependencies that confuse go
rm -rf vendor/k8s.io/*/vendor

cp ./generated_files /output/generated_files.after
cp ./lintconfig.json /output/lintconfig.json.after

pushd pilot
./bin/upload-istioctl -r -o "${OUTPUT_PATH}"
# An empty hub skips the tag and push steps.  -h "" provokes unset var error msg so using " "
./bin/push-docker -h " " -t "${TAG_NAME}" -b -o "${OUTPUT_PATH}"
./bin/push-debian.sh -c opt -v "${TAG_NAME}" -o "${OUTPUT_PATH}"
popd

# pilot dirties (both from bazel_to_go, it seems):
# generated_files
# lintconfig.json
../../../../repo status
ls -ld *

pushd mixer
./bin/push-docker           -h " " -t "${TAG_NAME}" -b -o "${OUTPUT_PATH}"
popd

# clean slate here
../../../../repo status

pushd security
./bin/push-docker           -h " " -t "${TAG_NAME}" -b -o "${OUTPUT_PATH}"
./bin/push-debian.sh -c opt -v "${TAG_NAME}" -o "${OUTPUT_PATH}"
popd

# seen dirty:
# mixer/template/apikey/go_default_library_tmpl.pb.go
# mixer/template/template.gen.go   
../../../../repo status

# bazel build //... dirties what's listed above plus:
# broker/pkg/model/config/mock_store.go
# broker/pkg/platform/kube/crd/types.go
# mixer/template/apikey/go_default_library_handler.gen.go

# store artifacts that are used by a separate cloud builder step to generate tar files
cp istio.VERSION LICENSE README.md CONTRIBUTING.md "${OUTPUT_PATH}/"
find samples install -type f \( -name "*.yaml" -o -name "cleanup*" -o -name "*.md" \) \
  -exec cp --parents {} "${OUTPUT_PATH}" \;
find install/tools -type f -exec cp --parents {} "${OUTPUT_PATH}" \;

../../../../repo status
