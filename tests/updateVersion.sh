#!/bin/bash

# Copyright 2017 Istio Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="${ROOT}/istio.VERSION"
GIT_COMMIT=false

set -o errexit
set -o pipefail
set -x

function usage() {
  [[ -n "${1}" ]] && echo "${1}"

  cat <<EOF
usage: ${BASH_SOURCE[0]} [options ...]"
  options:
    -i ... URL to download istioctl binaries
    -m ... <hub>,<tag> for the manager docker images
    -x ... <hub>,<tag> for the mixer docker images
    -c ... create a git commit for the changes
EOF
  exit 2
}

source "$VERSION_FILE" || error_exit "Could not source versions"

while getopts :i:m:x:c arg; do
  case ${arg} in
    i) ISTIOCTL_URL="${OPTARG}";;
    m) MANAGER_HUB_TAG="${OPTARG}";; # Format: "<hub>,<tag>"
    x) MIXER_HUB_TAG="${OPTARG}";; # Format: "<hub>,<tag>"
    c) GIT_COMMIT=true;;
    *) usage;;
  esac
done

if [[ -n ${MANAGER_HUB_TAG} ]]; then
    MANAGER_HUB="$(echo ${MANAGER_HUB_TAG}|cut -f1 -d,)"
    MANAGER_TAG="$(echo ${MANAGER_HUB_TAG}|cut -f2 -d,)"
fi

if [[ -n ${MIXER_HUB_TAG} ]]; then
    MIXER_HUB="$(echo ${MIXER_HUB_TAG}|cut -f1 -d,)"
    MIXER_TAG="$(echo ${MIXER_HUB_TAG}|cut -f2 -d,)"
fi

function error_exit() {
  # ${BASH_SOURCE[1]} is the file name of the caller.
  echo "${BASH_SOURCE[1]}: line ${BASH_LINENO[0]}: ${1:-Unknown Error.} (exit ${2:-1})" 1>&2
  exit ${2:-1}
}

function set_git() {
  if [[ ! -e "${HOME}/.gitconfig" ]]; then
    cat > "${HOME}/.gitconfig" << EOF
[user]
  name = istio-testing
  email = istio.testing@gmail.com
EOF
  fi
}


function create_commit() {
  set_git
  # If nothing to commit skip
  check_git_status && return

  echo 'Creating a commit'
  git commit -a -m "Updating istio version" \
    || error_exit 'Could not create a commit'

}

function check_git_status() {
  local git_files="$(git status -s)"
  [[ -z "${git_files}" ]] && return 0
  return 1
}

# Generated merge yaml files for easy installation, and adjust for 1.6 RBAC
function merge_files() {
  SRC=$ROOT/kubernetes/istio-install

  OUT=$ROOT/kubernetes/istio-1.6.yaml
  echo "# GENERATED FILE. Use for 1.6+ clusters." > $OUT
  echo "# TO UPDATE, modify files in istio-install and run updateVersion.sh" >> $OUT
  cat $SRC/istio-mixer.yaml >> $OUT
  cat $SRC/../istio-rbac/istio-rbac.yaml >> $OUT
  sed 's/# RBAC: //' $SRC/istio-manager.yaml >> $OUT
  sed 's/# RBAC: //' $SRC/istio-ingress-controller.yaml >> $OUT

  OUT=$ROOT/kubernetes/istio-1.6-auth.yaml
  echo "# GENERATED FILE. Used for 1.6+ clusters, will enable authentication." > $OUT
  echo "# TO UPDATE, modify files in istio-install and run updateVersion.sh" >> $OUT
  cat $SRC/istio-ca.yaml >> $OUT
  cat $SRC/istio-mixer.yaml >> $OUT
  cat $SRC/../istio-rbac/istio-rbac.yaml >> $OUT
  # Service account and authPolicy:MUTUAL_TLS
  sed 's/# RBAC: //' $SRC/istio-manager.yaml | sed 's/# authPolicy/authPolicy/' >> $OUT
  sed 's/# RBAC: //' $SRC/istio-ingress-controller.yaml >> $OUT

  OUT=$ROOT/kubernetes/istio-1.5.yaml
  echo "# GENERATED FILE" > $OUT
  echo "# TO UPDATE, modify files in istio-install and run updateVersion.sh" >> $OUT
  cat $SRC/istio-mixer.yaml >> $OUT
  cat $SRC/istio-manager.yaml >> $OUT
  cat $SRC/istio-ingress-controller.yaml >> $OUT
}

function update_version_file() {
  cat <<EOF > "${VERSION_FILE}"
# DO NOT EDIT THIS FILE MANUALLY instead use
# tests/updateVersion.sh (see tests/README.md)
export MIXER_HUB="${MIXER_HUB}"
export MIXER_TAG="${MIXER_TAG}"
export ISTIOCTL_URL="${ISTIOCTL_URL}"
export MANAGER_HUB="${MANAGER_HUB}"
export MANAGER_TAG="${MANAGER_TAG}"
EOF
}

function update_istio_install() {
  pushd $ROOT/kubernetes/istio-install
  sed -i "s|image: .*/\(.*\):.*|image: $MANAGER_HUB/\1:$MANAGER_TAG|" istio-manager.yaml
  sed -i "s|image: .*/\(.*\):.*|image: $MANAGER_HUB/\1:$MANAGER_TAG|" istio-ingress-controller.yaml
  sed -i "s|image: .*/\(.*\):.*|image: $MIXER_HUB/\1:$MIXER_TAG|" istio-mixer.yaml
  popd
}

if [[ ${GIT_COMMIT} == true ]]; then
    check_git_status \
      || error_exit "You have modified files. Please commit or reset your workspace."
fi

update_version_file
update_istio_install
merge_files

if [[ ${GIT_COMMIT} == true ]]; then
    create_commit
fi
