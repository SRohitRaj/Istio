#!/bin/sh

# Copyright 2019 Istio Authors
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

# Find Dockerfiles running an apt-get update but no uprade

BASE_DIR="$(cd "$(dirname "${0}")" && pwd -P)"
ISTIO_ROOT="$(cd "$(dirname "${BASE_DIR}")" && pwd -P)"
CD_TMPFILE=$(mktemp /tmp/check_dockerfile.XXXXXX)
HL_TMPFILE=$(mktemp /tmp/hadolint.XXXXXX)

find "${ISTIO_ROOT}" -name 'Dockerfile*' | \
while read -r f; do
  docker run --rm -i hadolint/hadolint < "$f" > "${HL_TMPFILE}"
  if [ "" != "$(cat "${HL_TMPFILE}")" ]
  then
    {
      echo "$f:"
      cut -d":" -f2- < "${HL_TMPFILE}"
      echo
    } >> "${CD_TMPFILE}"
  fi
done

rm -f "${HL_TMPFILE}"
if [ "" != "$(cat "${CD_TMPFILE}")" ]; then
  cat "${CD_TMPFILE}"
  rm -f "${CD_TMPFILE}"
  exit 1
fi
rm -f "${CD_TMPFILE}"
