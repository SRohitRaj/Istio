# Test-specific targets, included from top Makefile
ifeq (${TEST_ENV},minikube)

# In minikube env we don't need to push the images to dockerhub or gcr, it is all local,
# but we need to use the minikube's docker env.
export KUBECONFIG=${ISTIO_OUT}/minikube.conf
export TEST_ENV=minikube
MINIKUBE_FLAGS=-use_local_cluster -cluster_wide
.PHONY: minikube

# Prepare minikube
minikube:
	minikube update-context
	@echo "Minikube started ${KUBECONFIG}"
	minikube docker-env > ${ISTIO_OUT}/minikube.dockerenv

e2e_docker: minikube docker

else ifeq (${TEST_ENV},minikube-none)

# In minikube env we don't need to push the images to dockerhub or gcr, it is all local,
# but we need to use the minikube's docker env.
export KUBECONFIG=${ISTIO_OUT}/minikube.conf
export TEST_ENV=minikube-none
MINIKUBE_FLAGS=-use_local_cluster -cluster_wide
.PHONY: minikube

# Prepare minikube
minikube:
	minikube update-context
	@echo "Minikube started ${KUBECONFIG}"

e2e_docker: minikube docker


else

# All other test environments require the docker images to be pushed to a repo.
# The HUB is defined in user-specific .istiorc, TAG can be set or defaults to git version
e2e_docker: push

endif

E2E_ARGS ?=
E2E_ARGS += ${MINIKUBE_FLAGS}
E2E_ARGS += --istioctl ${ISTIO_OUT}/istioctl

EXTRA_E2E_ARGS = --mixer_tag ${TAG}
EXTRA_E2E_ARGS += --pilot_tag ${TAG}
EXTRA_E2E_ARGS += --ca_tag ${TAG}
EXTRA_E2E_ARGS += --mixer_hub ${HUB}
EXTRA_E2E_ARGS += --pilot_hub ${HUB}
EXTRA_E2E_ARGS += --ca_hub ${HUB}

# Simple e2e test using fortio, approx 2 min
e2e_simple: istioctl installgen
	go test -v -timeout 20m ./tests/e2e/tests/simple -args ${E2E_ARGS} ${EXTRA_E2E_ARGS}

e2e_mixer: istioctl installgen
	go test -v -timeout 20m ./tests/e2e/tests/mixer -args ${E2E_ARGS} ${EXTRA_E2E_ARGS}

e2e_bookinfo: istioctl installgen
	go test -v -timeout 20m ./tests/e2e/tests/bookinfo -args ${E2E_ARGS} ${EXTRA_E2E_ARGS}

e2e_all: e2e_simple e2e_mixer e2e_bookinfo
