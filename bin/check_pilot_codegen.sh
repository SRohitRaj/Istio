#!/bin/bash

# Make sure the types.go is generated by scripts.

set -ex

SCRIPTPATH=$( cd "$(dirname "$0")" ; pwd -P )
ROOTDIR=$SCRIPTPATH/..
pushd $ROOTDIR

expected='/tmp/types.go'

go run pilot/tools/generate_config_crd_types.go --template pilot/tools/types.go.tmpl --output $expected

actual='pilot/pkg/config/kube/crd/types.go'
dout=$(diff $actual $expected)

if [[ $dout ]]; then
  exit 1
fi

