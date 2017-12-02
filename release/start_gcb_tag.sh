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

# This is an example file for how to start an istio release publish
# using Cloud Builder. To run it you need a Google Cloud project and
# the key file for a service account that had been granted access to start a build.

set -o errexit
set -o nounset
set -o pipefail
set -x

SCRIPTPATH=$( cd "$(dirname "$0")" ; pwd -P )
source ${SCRIPTPATH}/gcb_build_lib.sh

PROJECT_ID=""
KEY_FILE_PATH=""
SVC_ACCT=""
SUBS_FILE="$(mktemp /tmp/build.subs.XXXX)"
VER_STRING="0.0.0"
REPO=""
REPO_FILE=default.xml
REPO_FILE_VER=master
WAIT_FOR_RESULT="false"

GCS_SRC=""
GCS_GITHUB_SECRET="istio-secrets/github.txt"
REL_ORG="istio"
USER_EMAIL=""
USER_NAME=""

function usage() {
  echo "$0
    -p <name> project ID                                        (required)
    -a        service account for login                         (optional, defaults to project's cloudbuild@ )
    -k <file> path to key file for service account              (required)
    -v <ver>  version string                                    (optional, defaults to $VER_STRING )
    -u <url>  URL to git repo with manifest file                (required)
    -m <file> name of manifest file in repo specified by -u     (optional, defaults to $REPO_FILE )
    -t <tag>  commit tag or branch for manifest repo in -u      (optional, defaults to $REPO_FILE_VER )
    -w        specify that script should wait until build done  (optional)  

    -s <name> GCS bucket to read build artifacts                (required)
    -g <path> GCS bucket&path to file with github secret        (optional, detaults to $GCS_GITHUB_SECRET )
    -h <name> github org to tag                                 (optional, defaults to $REL_ORG )
    -e <email> email of submitter for tags                      (required)
    -n <name>  name of submitter for tags                       (required)"
  exit 1
}

while getopts a:e:g:h:k:m:n:p:s:t:u:v:w arg ; do
  case "${arg}" in
    a) SVC_ACCT="${OPTARG}";;
    e) USER_EMAIL="${OPTARG}";;
    g) GCS_GITHUB_SECRET="${OPTARG}";;
    h) REL_ORG="${OPTARG}";;
    k) KEY_FILE_PATH="${OPTARG}";;
    m) REPO_FILE="${OPTARG}";;
    n) USER_NAME="${OPTARG}";;
    p) PROJECT_ID="${OPTARG}";;
    s) GCS_SRC="${OPTARG}";;
    t) REPO_FILE_VER="${OPTARG}";;
    u) REPO="${OPTARG}";;
    v) VER_STRING="${OPTARG}";;
    w) WAIT_FOR_RESULT="true";;
    *) usage;;
  esac
done

[[ -z "${PROJECT_ID}"    ]] && usage
[[ -z "${KEY_FILE_PATH}" ]] && usage
[[ -z "${REPO}"          ]] && usage
[[ -z "${REPO_FILE}"     ]] && usage
[[ -z "${REPO_FILE_VER}" ]] && usage
[[ -z "${VER_STRING}"    ]] && usage

[[ -z "${USER_EMAIL}" ]] && usage
[[ -z "${USER_NAME}"  ]] && usage
[[ -z "${GCS_SRC}"    ]] && usage

DEFAULT_SVC_ACCT="cloudbuild@${PROJECT_ID}.iam.gserviceaccount.com"

if [[ -z "${SVC_ACCT}"  ]]; then
  SVC_ACCT="${DEFAULT_SVC_ACCT}"
fi

# generate the substitutions file
echo "  \"substitutions\": {
    \"_VER_STRING\": \"${VER_STRING}\",
    \"_GCS_SOURCE\": \"${GCS_SRC}\",
    \"_GCS_SECRET\": \"${GCS_GITHUB_SECRET}\",
    \"_ORG\": \"${REL_ORG}\",
    \"_USER_EMAIL\": \"${USER_EMAIL}\",
    \"_USER_NAME\": \"${USER_NAME}\",
  }" >> "${SUBS_FILE}"

run_build "${REPO}" "${REPO_FILE}" "${REPO_FILE_VER}" "cloud_tag.template.json" \
  "${SUBS_FILE}" "${PROJECT_ID}" "${SVC_ACCT}" "${KEY_FILE_PATH}" "${WAIT_FOR_RESULT}"

# cleanup
rm -f "${SUBS_FILE}"
exit $BUILD_FAILED
