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

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TESTS_DIR="${ROOT}/tests"
. ${TESTS_DIR}/commonUtils.sh || { echo "Cannot load common utilities"; exit 1; }

function create_rule() {
    ${ISTIOCLI} -n ${NAMESPACE} --istioNamespace ${NAMESPACE} create -f ${1} \
      || error_exit 'Could not create rule'
}

function replace_rule() {
    ${ISTIOCLI} -n ${NAMESPACE} --istioNamespace ${NAMESPACE} replace -f ${1} \
      || error_exit 'Could not replace rule'
}

function cleanup_all_rules() {
    print_block_echo "Cleaning up rules"
    local rules=($(${ISTIOCLI} -n ${NAMESPACE} --istioNamespace ${NAMESPACE} get route-rule \
      | grep "name:" | awk '{print $2}'))
    for r in ${rules[@]}; do
      ${ISTIOCLI} -n ${NAMESPACE} --istioNamespace ${NAMESPACE} delete route-rule "${r}"
    done
}

function delete_rule() {
    ${ISTIOCLI} -n ${NAMESPACE} --istioNamespace ${NAMESPACE} delete -f ${1}
    return $?
}
