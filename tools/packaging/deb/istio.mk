# Creates the 2 packages. BUILD_WITH_CONTAINER=1 or in CI/CD
#
# Development/manual testing:
#    BUILD_WITH_CONTAINER=1 make deb
#    make deb/docker # will create the istio_deb container with istio installed
#

deb: ${ISTIO_OUT_LINUX}/release/istio-sidecar.deb ${ISTIO_OUT_LINUX}/release/istio.deb

# Base directory for istio binaries. Likely to change !
ISTIO_DEB_BIN=/usr/local/bin

ISTIO_DEB_DEPS:=pilot-discovery istioctl
ISTIO_FILES:=
$(foreach DEP,$(ISTIO_DEB_DEPS),\
        $(eval ${ISTIO_OUT_LINUX}/release/istio.deb: $(ISTIO_OUT_LINUX)/$(DEP)) \
        $(eval ISTIO_FILES+=$(ISTIO_OUT_LINUX)/$(DEP)=$(ISTIO_DEB_BIN)/$(DEP)) )

SIDECAR_DEB_DEPS:=envoy pilot-agent node_agent istio-iptables istio-clean-iptables
SIDECAR_FILES:=
$(foreach DEP,$(SIDECAR_DEB_DEPS),\
        $(eval ${ISTIO_OUT_LINUX}/release/istio-sidecar.deb: $(ISTIO_OUT_LINUX)/$(DEP)) \
        $(eval SIDECAR_FILES+=$(ISTIO_OUT_LINUX)/$(DEP)=$(ISTIO_DEB_BIN)/$(DEP)) )

ISTIO_DEB_DEST:=${ISTIO_DEB_BIN}/istio-start.sh \
		/lib/systemd/system/istio.service \
		/var/lib/istio/envoy/sidecar.env

$(foreach DEST,$(ISTIO_DEB_DEST),\
        $(eval ${ISTIO_OUT_LINUX}/istio-sidecar.deb:   tools/packaging/common/$(notdir $(DEST))) \
        $(eval SIDECAR_FILES+=${REPO_ROOT}/tools/packaging/common/$(notdir $(DEST))=$(DEST)))

SIDECAR_FILES+=${REPO_ROOT}/tools/packaging/common/envoy_bootstrap_v2.json=/var/lib/istio/envoy/envoy_bootstrap_tmpl.json


# original name used in 0.2 - will be updated to 'istio.deb' since it now includes all istio binaries.
ISTIO_DEB_NAME ?= istio-sidecar

# TODO: rename istio-sidecar.deb to istio.deb

# Note: adding --deb-systemd ${REPO_ROOT}/tools/packaging/common/istio.service will result in
# a /etc/systemd/system/multi-user.target.wants/istio.service and auto-start. Currently not used
# since we need configuration.
# --iteration 1 adds a "-1" suffix to the version that didn't exist before
${ISTIO_OUT_LINUX}/release/istio-sidecar.deb: | ${ISTIO_OUT_LINUX} deb/fpm

# Package the sidecar deb file.
deb/fpm:
	rm -f ${ISTIO_OUT_LINUX}/release/istio-sidecar.deb
	fpm -s dir -t deb -n ${ISTIO_DEB_NAME} -p ${ISTIO_OUT_LINUX}/release/istio-sidecar.deb --version $(PACKAGE_VERSION) -f \
		--url http://istio.io  \
		--license Apache \
		--vendor istio.io \
		--maintainer istio@istio.io \
		--after-install tools/packaging/deb/postinst.sh \
		--config-files /var/lib/istio/envoy/envoy_bootstrap_tmpl.json \
		--config-files /var/lib/istio/envoy/sidecar.env \
		--description "Istio Sidecar" \
		--depends iproute2 \
		--depends iptables \
		$(SIDECAR_FILES)

${ISTIO_OUT_LINUX}/release/istio.deb:
	rm -f ${ISTIO_OUT_LINUX}/release/istio.deb
	fpm -s dir -t deb -n istio -p ${ISTIO_OUT_LINUX}/release/istio.deb --version $(PACKAGE_VERSION) -f \
		--url http://istio.io  \
		--license Apache \
		--vendor istio.io \
		--maintainer istio@istio.io \
		--description "Istio" \
		$(ISTIO_FILES)

GEN_CERT ?= go run istio.io/istio/security/tools/generate_cert
# TODO: use k8s style - /etc/pki/istio/...
PKI_DIR ?= tests/testdata/certs/cacerts
VM_PKI_DIR ?= tests/testdata/certs/vm

testcert-gen:
	mkdir -p ${PKI_DIR}
	mkdir -p ${VM_PKI_DIR}
	${GEN_CERT} -ca  --out-priv ${PKI_DIR}/ca-key.pem --out-cert ${PKI_DIR}/ca-cert.pem  -organization "istio ca"
	${GEN_CERT} -signer-cert ${PKI_DIR}/ca-cert.pem -signer-priv ${PKI_DIR}/ca-key.pem \
	 	-out-cert ${VM_PKI_DIR}/cert-chain.pem -out-priv ${VM_PKI_DIR}/key.pem \
	 	-host spiffe://cluster.local/ns/vmtest/sa/default --mode signer
     cp ${PKI_DIR}/ca-cert.pem ${VM_PKI_DIR}/root-cert.pem

# Install the deb in a docker image, for testing the install process.
# Will use a minimal base image, install all that is needed.
deb/docker: testcert-gen
	mkdir -p ${ISTIO_OUT_LINUX}/deb
	cp tools/packaging/deb/Dockerfile tools/packaging/deb/deb_test.sh ${ISTIO_OUT_LINUX}/deb
	# Istio configs, for testing istiod running in the VM.
	cp tests/testdata/config/*.yaml ${ISTIO_OUT_LINUX}/deb
	# Test certificates - can be used to verify connection with an istiod running on the host or
	# in a separate container.
	cp -a tests/testdata/certs ${ISTIO_OUT_LINUX}/deb
	cp ${ISTIO_OUT_LINUX}/release/istio-sidecar.deb ${ISTIO_OUT_LINUX}/deb/istio-sidecar.deb
	cp ${ISTIO_OUT_LINUX}/release/istio.deb ${ISTIO_OUT_LINUX}/deb/istio.deb
	docker build -t istio_deb -f ${ISTIO_OUT_LINUX}/deb/Dockerfile  ${ISTIO_OUT_LINUX}/deb/

# For the test, by default use a local pilot.
# Set it to 172.18.0.1 to run against a pilot running in IDE.
# You may need to enable 15007 in the local machine firewall for this to work.
DEB_PILOT_IP ?= 127.0.0.1
DEB_CMD ?= /bin/bash
ISTIO_NET ?= 172.18
DEB_IP ?= ${ISTIO_NET}.0.3
DEB_PORT_PREFIX ?= 1600

# TODO: docker compose ?


# Run the docker image including the installed debian, with access to all source
# code. Useful for debugging/experiments with iptables.
#
# Before running:
# docker network create --subnet=172.18.0.0/16 istiotest
# The IP of the docker matches the byon-docker service entry in the static configs, if testing without k8s.
#
# On host, run istiod (can be standalone), using Kind or real K8S cluster:
#
# export TOKEN_ISSUER=https://localhost # Dummy, to ignore missing token. Can be real OIDC server.
# export MASTER_ELECTION=false
# istiod discovery -n istio-system
#
deb/run/docker:
	docker run --cap-add=NET_ADMIN --rm \
	  -v ${GO_TOP}:${GO_TOP} \
      -w ${PWD} \
      --mount type=bind,source="$(HOME)/.kube",destination="/home/.kube" \
      --mount type=bind,source="$(GOPATH)",destination="/ws" \
      --net istiotest --ip ${DEB_IP} \
      --add-host echo:10.1.1.1 \
      --add-host byon.test.istio.io:10.1.1.2 \
      --add-host byon-docker.test.istio.io:10.1.1.2 \
      --add-host istiod.istio-system.svc:${DEB_PILOT_IP} \
      ${DEB_ENV} -e ISTIO_SERVICE_CIDR=10.1.1.0/24 \
      -e ISTIO_INBOUND_PORTS=7070,7072,7073,7074,7075 \
      -e PILOT_CERT_DIR=/var/lib/istio/pilot \
      -p 127.0.0.1:${DEB_PORT_PREFIX}1:15007 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}2:7070 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}3:7072 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}4:7073 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}5:7074 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}6:7075 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}7:15012 \
      -p 127.0.0.1:${DEB_PORT_PREFIX}8:15010 \
      -e GOPATH=${GOPATH} \
      -it istio_deb ${DEB_CMD}

deb/run/debug:
	$(MAKE) deb/run/docker DEB_ENV="-e DEB_PILOT_IP=172.18.0.1"

deb/run/tproxy:
	$(MAKE) deb/run/docker DEB_PORT_PREFIX=1610 DEB_IP=172.18.0.4 DEB_ENV="-e ISTIO_INBOUND_INTERCEPTION_MODE=TPROXY"

deb/run/mtls:
	$(MAKE) deb/run/docker DEB_PORT_PREFIX=1620 -e DEB_PILOT_IP=172.18.0.1 DEB_IP=172.18.0.5 DEB_ENV="-e ISTIO_PILOT_PORT=15005 -e ISTIO_CP_AUTH=MUTUAL_TLS"

# Similar with above, but using a pilot running on the local machine
deb/run/docker-debug:
	$(MAKE) deb/run/docker PILOT_IP=

#
deb/docker-run: deb/docker deb/run/docker

.PHONY: \
	deb \
	deb/build-in-docker \
	deb/docker \
	deb/docker-run \
	deb/run/docker \
	deb/fpm \
	deb/test \
	sidecar.deb
