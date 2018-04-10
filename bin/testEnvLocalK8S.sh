#!/usr/bin/env bash

set -euo pipefail

# testEnv will setup a local test environment, for running Istio unit tests.

# Based on circleCI config - used to reproduce the environment and to improve local testing

# expect istio scripts to be under $GOPATH/src/istio.io/istio/bin/...

export TOP=${TOP:-${GOPATH}}
export ISTIO_GO=${TOP}/src/istio.io/istio

if [[ "$OSTYPE" == "darwin"* ]]; then 
   export GOOS_LOCAL=darwin
else
  export GOOS_LOCAL=${GOOS_LOCAL:-linux}
fi

export GOPATH=${TOP}
export PATH=${GOPATH}/bin:${PATH}
export OUT=${TOP}/out
export ISTIO_OUT=${ISTIO_OUT:-${TOP}/out/${GOOS_LOCAL}_amd64/release}

# components used in the test (starting with circleci for consistency, eventually ci will use this)
export K8S_VER=${K8S_VER:-v1.9.2}
export ETCD_VER=${ETCD_VER:-v3.2.15}

export MASTER_IP=127.0.0.1
export MASTER_CLUSTER_IP=10.99.0.1

# TODO: customize the ports and generate a local config
export KUBECONFIG=${TOP}/src/istio.io/istio/.circleci/config

${ISTIO_GO}/bin/init.sh

# Checked in certificates, to avoid regenerating them
CERTDIR=${CERTDIR:-${ISTIO_GO}/.circleci/pki/istio-certs}
LOG_DIR=${LOG_DIR:-${OUT}/log}
ETCD_DATADIR=${ETCD_DATADIR:-${OUT}/etcd-data}

EASYRSA_DIR=$OUT/easy-rsa-master/easyrsa3
EASYRSA=$EASYRSA_DIR/easyrsa

# Ensure k8s certificats - if not found, download easy-rsa and create k8s certs
function ensureK8SCerts() {
    if [ -f ${CERTDIR}/apiserver.key ] ; then
        return
    fi

    mkdir -p ${CERTDIR}
    pushd $OUT
    curl -L -O https://storage.googleapis.com/kubernetes-release/easy-rsa/easy-rsa.tar.gz
    tar xzf easy-rsa.tar.gz
    cd easy-rsa-master/easyrsa3

    ./easyrsa init-pki > /dev/null
    ./easyrsa --batch "--req-cn=${MASTER_IP}@`date +%s`" build-ca nopass > /dev/null
    ./easyrsa --subject-alt-name="IP:${MASTER_IP},""IP:${MASTER_CLUSTER_IP},""DNS:kubernetes,""DNS:kubernetes.default,""DNS:kubernetes.default.svc,""DNS:kubernetes.default.svc.cluster,""DNS:kubernetes.default.svc.cluster.local" \
        --days=10000 build-server-full server nopass > /dev/null

    cp pki/private/ca.key ${CERTDIR}/k8sca.key
    cp pki/ca.crt ${CERTDIR}/k8sca.crt
    cp pki/issued/server.crt ${CERTDIR}/apiserver.crt
    cp pki/private/server.key ${CERTDIR}/apiserver.key
    popd
}

# Get dependencies needed for tests. Only needed once.
# The docker builder image should include them.
function getDeps() {
   mkdir -p $TOP/bin
   if [ ! -f $TOP/bin/kubectl ] ; then
     if [ -f /usr/local/bin/kubectl ] ; then
       ln -s /usr/local/bin/kubectl $TOP/bin/kubectl
     else
       curl -Lo $TOP/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${K8S_VER}/bin/${GOOS_LOCAL}/amd64/kubectl && chmod +x $TOP/bin/kubectl
     fi
   fi
   if [ ! -f $TOP/bin/kube-apiserver ] ; then
     if [ -f /usr/local/bin/kube-apiserver ] ; then
       ln -s /usr/local/bin/kube-apiserver $TOP/bin/
     elif [ -f /tmp/apiserver/kube-apiserver ] ; then
       ln -s /tmp/apiserver/kube-apiserver $TOP/bin/
     else
       # bucket doesn't contain a kube-apiserver for darwin
       curl -Lo ${TOP}/bin/kube-apiserver https://storage.googleapis.com/kubernetes-release/release/${K8S_VER}/bin/${GOOS_LOCAL}/amd64/kube-apiserver && chmod +x ${TOP}/bin/kube-apiserver
     fi
   fi
   if [ ! -f $TOP/bin/etcd ] ; then
     if [ -f /usr/local/bin/etcd ] ; then
        ln -s /usr/local/bin/etcd $TOP/bin/
     else
       if [ "${GOOS_LOCAL}" == "darwin" ]; then
	   # I tried using unzip -p <(curl) but curl is launched async and unzip doesn't wait
           ETC_TEMP=$(mktemp)
           curl -L https://github.com/coreos/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-darwin-amd64.zip > ${ETC_TEMP}
           unzip -p ${ETC_TEMP} etcd-${ETCD_VER}-darwin-amd64/etcd > ${TOP}/bin/etcd
           chmod +x ${TOP}/bin/etcd
           rm ${ETC_TEMP}
       else
	   curl -L https://github.com/coreos/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz | tar xz -O etcd-${ETCD_VER}-linux-amd64/etcd > ${TOP}/bin/etcd && chmod +x ${TOP}/bin/etcd
       fi
     fi
   fi
   if [ ! -f $TOP/bin/envoy ] ; then
     # Init should be run after cloning the workspace
     ${ISTIO_GO}/bin/init.sh
   fi
}

# No root required, run local etcd and kube apiserver for tests.
function startLocalApiserver() {
    ensureK8SCerts
    getDeps

    mkdir -p ${LOG_DIR}
    mkdir -p ${ETCD_DATADIR}
    ${TOP}/bin/etcd --data-dir ${ETCD_DATADIR} > ${LOG_DIR}/etcd.log 2>&1 &
    echo $! > $LOG_DIR/etcd.pid
    # make sure etcd is actually alive
    kill -0 $(cat $LOG_DIR/etcd.pid)

    ${TOP}/bin/kube-apiserver --etcd-servers http://127.0.0.1:2379 \
        --client-ca-file ${CERTDIR}/k8sca.crt \
        --requestheader-client-ca-file ${CERTDIR}/k8sca.crt \
        --tls-cert-file ${CERTDIR}/apiserver.crt \
        --tls-private-key-file ${CERTDIR}/apiserver.key \
        --service-cluster-ip-range 10.99.0.0/16 \
        --port 8080 -v 2 --insecure-bind-address 0.0.0.0 \
        > ${LOG_DIR}/apiserver.log 2>&1 &
    echo $! > $LOG_DIR/apiserver.pid
    # make sure apiserver is actually alive
    kill -0 $(cat $LOG_DIR/apiserver.pid)

    # Really need to make sure that API Server is up before proceed further
    waitForApiServer
    printf "Started local etcd and apiserver!\n"
}

function ensureLocalApiServer() {
    kubectl get nodes 2>/dev/null || startLocalApiserver
}

function createIstioConfigmap() {
  helm template ${ISTIO_GO}/install/kubernetes/helm/istio --namespace=istio-system \
     --execute=templates/configmap.yaml --values install/kubernetes/helm/istio/values.yaml  > ${LOG_DIR}/istio-configmap.yaml
  kubectl create -f ${LOG_DIR}/istio-configmap.yaml 
  helm template ${ISTIO_GO}/install/kubernetes/helm/istio --namespace=istio-system \
     --execute=charts/ingress/templates/service.yaml --values install/kubernetes/helm/istio/values.yaml  > ${LOG_DIR}/istio-ingress.yaml
  kubectl create -f ${LOG_DIR}/istio-ingress.yaml
}

function startIstio() {
    ensureLocalApiServer
    createIstioConfigmap
    startPilot
    startEnvoy
    startMixer
}

function stopIstio() {
  if [[ -f $LOG_DIR/pilot.pid ]] ; then
    kill -9 $(cat $LOG_DIR/pilot.pid) || true
    rm $LOG_DIR/pilot.pid
   fi
  if [[ -f $LOG_DIR/mixer.pid ]] ; then
    kill -9 $(cat $LOG_DIR/mixer.pid) || true
    rm $LOG_DIR/mixer.pid
  fi
  if [[ -f $LOG_DIR/envoy4.pid ]] ; then
    kill -9 $(cat $LOG_DIR/envoy4.pid) || true
    rm $LOG_DIR/envoy4.pid
  fi
}

function startPilot() {
  printf "Pilot starting...\n"
  POD_NAME=pilot POD_NAMESPACE=istio-system \
  ${ISTIO_OUT}/pilot-discovery discovery --port 18080 \
                                         --monitoringPort 19093 \
                                         --log_target ${LOG_DIR} \
                                         --kubeconfig ${ISTIO_GO}/.circleci/config &
  echo $! > $LOG_DIR/pilot.pid
}

function startMixer() {
  printf "Mixer starting...\n"
  ${ISTIO_OUT}/mixs server --configStoreURL=fs:${ISTIO_GO}/mixer/testdata/configroot & 
  echo $! > $LOG_DIR/mixer.pid
}

function startEnvoy() {
    printf "Envoy starting...\n"
    ${ISTIO_OUT}/envoy -c tests/testdata/envoy_local.json \
        --base-id 4 --service-cluster unittest --service-node local.test
    echo $! > $LOG_DIR/envoy4.pid
}

function stopLocalApiserver() {
  if [[ -f $LOG_DIR/etcd.pid ]]; then
    kill -9 $(cat $LOG_DIR/etcd.pid)
    kill -9 $(cat $LOG_DIR/apiserver.pid)
    rm $LOG_DIR/{etcd,apiserver}.pid
  fi
  if [[ -d "${ETCD_DATADIR}" ]]; then
    rm -rf ${ETCD_DATADIR}
  fi
}

function startLocalServers() {
    startLocalApiserver
    startPilot
    startEnvoy
}

function waitForApiServer() {
count=0
set +xe

  while true; do
    status=$(kubectl get pod 2>&1 | grep resources | wc -l)
    if [ $status -ne 1 ]; then
      if [ $count -gt 30 ]; then
        printf "API Server failed to come up\n"
        exit -1
      fi
      count=$((count+1))
      sleep 1
    else
      printf "API Server ready\n"
      break
    fi
  done
}

CMD=${1:-help}
case "$1" in
    start) startLocalApiserver ;;
    stop) stopLocalApiserver ;;
    startIstio) startIstio ;;
    stopIstio) stopIstio ;;
    ensure) ensureLocalApiServer ;;
    *) printf "start stop ensure\n"
esac
