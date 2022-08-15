#!/bin/bash
#
# Copyright Istio Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

set -ox errexit

if [ "$#" -ne 2 ]; then
    echo "Incorrect parameters"
    echo "Usage: build-services.sh <version> <prefix>"
    exit 1
fi

VERSION=$1
PREFIX=$2
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# Docker build variables
ENABLE_MULTIARCH_IMAGES=${ENABLE_MULTIARCH_IMAGES:-"false"}

if [ "${ENABLE_MULTIARCH_IMAGES}" == "true" ]; then
  PLATFORMS="linux/arm64,linux/amd64"	
  DOCKER_BUILD_ARGS="docker buildx build --platform ${PLATFORMS} --push"
  # Install QEMU emulators
  docker run --rm --privileged tonistiigi/binfmt --install all
  docker buildx rm multi-builder || :
  docker buildx create --use --name multi-builder --platform ${PLATFORMS}
  docker buildx use multi-builder
else
  DOCKER_BUILD_ARGS="docker build"	
fi

pushd "$SCRIPTDIR/productpage"
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-productpage-v1:${VERSION}" -t "${PREFIX}/examples-bookinfo-productpage-v1:latest" .
  #flooding
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-productpage-v-flooding:${VERSION}" -t "${PREFIX}/examples-bookinfo-productpage-v-flooding:latest" --build-arg flood_factor=100 .
popd

pushd "$SCRIPTDIR/details"
  #plain build -- no calling external book service to fetch topics
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-details-v1:${VERSION}" -t "${PREFIX}/examples-bookinfo-details-v1:latest" --build-arg service_version=v1 .
  #with calling external book service to fetch topic for the book
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-details-v2:${VERSION}" -t "${PREFIX}/examples-bookinfo-details-v2:latest" --build-arg service_version=v2 \
	 --build-arg enable_external_book_service=true .
popd


pushd "$SCRIPTDIR/reviews"
  #java build the app.
  docker run --rm -u root -v "$(pwd)":/home/gradle/project -w /home/gradle/project gradle:4.8.1 gradle clean build
  
  # Temporarily clone https://github.com/OpenLiberty/ci.docker
  # to allow this project to build multiarch(linux/arm64) images not provided by the OpenLiberty CI
  # Remove this workaround once https://github.com/OpenLiberty/ci.docker/issues/241 is fixed
  git clone https://github.com/OpenLiberty/ci.docker.git openliberty-base
  pushd openliberty-base
    ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/websphere-liberty:22.0.0.8-full-java8-ibmjava" -f releases/22.0.0.8/full/Dockerfile.ubuntu.openjdk8 releases/22.0.0.8/full 
  popd
  rm -rf openliberty-base 
  
  pushd reviews-wlpcfg
    #plain build -- no ratings
    ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-reviews-v1:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v1:latest" --build-arg service_version=v1 \
	    --build-arg REGISTRY="${PREFIX}" .
    #with ratings black stars
    ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-reviews-v2:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v2:latest" --build-arg service_version=v2 \
	   --build-arg enable_ratings=true --build-arg REGISTRY="${PREFIX}" .
    #with ratings red stars
    ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-reviews-v3:${VERSION}" -t "${PREFIX}/examples-bookinfo-reviews-v3:latest" --build-arg service_version=v3 \
	   --build-arg enable_ratings=true --build-arg star_color=red --build-arg REGISTRY="${PREFIX}" .
  popd
popd

pushd "$SCRIPTDIR/ratings"
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v1:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v1:latest" --build-arg service_version=v1 .
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v2:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v2:latest" --build-arg service_version=v2 .
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v-faulty:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v-faulty:latest" --build-arg service_version=v-faulty .
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v-delayed:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v-delayed:latest" --build-arg service_version=v-delayed .
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v-unavailable:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v-unavailable:latest" --build-arg service_version=v-unavailable .
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-ratings-v-unhealthy:${VERSION}" -t "${PREFIX}/examples-bookinfo-ratings-v-unhealthy:latest" --build-arg service_version=v-unhealthy .
popd

pushd "$SCRIPTDIR/mysql"
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-mysqldb:${VERSION}" -t "${PREFIX}/examples-bookinfo-mysqldb:latest" .
popd

pushd "$SCRIPTDIR/mongodb"
  ${DOCKER_BUILD_ARGS} --pull -t "${PREFIX}/examples-bookinfo-mongodb:${VERSION}" -t "${PREFIX}/examples-bookinfo-mongodb:latest" .
popd
