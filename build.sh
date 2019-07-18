#!/usr/bin/env bash
set -e

PROJECT=swarm-cronjob
VERSION=${TRAVIS_TAG:-dev}
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_TAG=docker_build
BUILD_WORKINGDIR=${BUILD_WORKINGDIR:-.}
DOCKERFILE=${DOCKERFILE:-Dockerfile}
VCS_REF=${TRAVIS_COMMIT::8}
RUNNING_TIMEOUT=120
RUNNING_LOG_CHECK="Number of cronjob tasks: 5"

PUSH_LATEST=${PUSH_LATEST:-true}
DOCKER_USERNAME=${DOCKER_USERNAME:-crazymax}
DOCKER_LOGIN=${DOCKER_LOGIN:-crazymax}
DOCKER_REPONAME=${DOCKER_REPONAME:-swarm-cronjob}
QUAY_USERNAME=${QUAY_USERNAME:-crazymax}
QUAY_LOGIN=${QUAY_LOGIN:-crazymax}
QUAY_REPONAME=${QUAY_REPONAME:-swarm-cronjob}

# Check dev or travis
BRANCH=${TRAVIS_BRANCH:-local}
if [[ ${TRAVIS_PULL_REQUEST} == "true" ]]; then
  BRANCH=${TRAVIS_PULL_REQUEST_BRANCH}
fi
DOCKER_TAG=${BRANCH:-local}
if [[ "$BRANCH" == "local" ]]; then
  BUILD_DATE=
else
  DOCKER_TAG=latest
  VERSION=${VERSION#v}
fi

echo "PROJECT=${PROJECT}"
echo "VERSION=${VERSION}"
echo "BUILD_DATE=${BUILD_DATE}"
echo "BUILD_TAG=${BUILD_TAG}"
echo "BUILD_WORKINGDIR=${BUILD_WORKINGDIR}"
echo "DOCKERFILE=${DOCKERFILE}"
echo "VCS_REF=${VCS_REF}"
echo "PUSH_LATEST=${PUSH_LATEST}"
echo "DOCKER_LOGIN=${DOCKER_LOGIN}"
echo "DOCKER_USERNAME=${DOCKER_USERNAME}"
echo "DOCKER_REPONAME=${DOCKER_REPONAME}"
echo "QUAY_LOGIN=${QUAY_LOGIN}"
echo "QUAY_USERNAME=${QUAY_USERNAME}"
echo "QUAY_REPONAME=${QUAY_REPONAME}"
echo "TRAVIS_BRANCH=${TRAVIS_BRANCH}"
echo "TRAVIS_PULL_REQUEST=${TRAVIS_PULL_REQUEST}"
echo "BRANCH=${BRANCH}"
echo "DOCKER_TAG=${DOCKER_TAG}"
echo

echo "### Goreleaser"
if [[ -n "$TRAVIS_TAG" ]]; then
  goreleaser release --skip-publish --rm-dist
  else
  goreleaser release --snapshot --rm-dist
fi

echo "### Docker build"
docker build \
  --build-arg BUILD_DATE=${BUILD_DATE} \
  --build-arg VCS_REF=${VCS_REF} \
  --build-arg VERSION=${VERSION} \
  -t ${BUILD_TAG} -f ${DOCKERFILE} ${BUILD_WORKINGDIR}
echo

echo "### Test"
docker swarm leave --force > /dev/null 2>&1 || true
docker swarm init --advertise-addr ${ADVERTISE_ADDR:-192.168.99.100}
docker stack deploy date -c .res/example/date.yml
docker stack deploy sleep -c .res/example/sleep.yml
docker stack deploy error -c .res/example/error.yml
docker stack deploy configs -c .res/example/configs.yml
docker stack deploy global -c .res/example/global.yml
docker service create --name ${PROJECT} \
  --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock \
  --env "LOG_LEVEL=debug" \
  --env "LOG_JSON=false" \
  --constraint 'node.role == manager' \
  ${BUILD_TAG}
echo

echo "### Waiting for ${PROJECT} to be up..."
TIMEOUT=$((SECONDS + RUNNING_TIMEOUT))
while read LOGLINE; do
  echo ${LOGLINE}
  if [[ ${LOGLINE} == *"${RUNNING_LOG_CHECK}"* ]]; then
    echo "Service up!"
    break
  fi
  if [[ $SECONDS -gt ${TIMEOUT} ]]; then
    >&2 echo "ERROR: Failed to run ${PROJECT} container"
    docker swarm leave --force > /dev/null 2>&1 || true
    exit 1
  fi
done < <(docker service logs -f ${PROJECT} 2>&1)
echo
docker swarm leave --force > /dev/null 2>&1 || true

if [ "${VERSION}" == "dev" -o "${TRAVIS_PULL_REQUEST}" == "true" ]; then
  echo "INFO: This is a PR or an untagged build, skipping push..."
  exit 0
fi
if [[ ! -z ${DOCKER_PASSWORD} ]]; then
  echo "### Push to Docker Hub..."
  echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_LOGIN" --password-stdin > /dev/null 2>&1
  if [ "${DOCKER_TAG}" == "latest" -a "${PUSH_LATEST}" == "true" ]; then
    docker tag ${BUILD_TAG} ${DOCKER_USERNAME}/${DOCKER_REPONAME}:${DOCKER_TAG}
  fi
  if [[ "${VERSION}" != "latest" ]]; then
    docker tag ${BUILD_TAG} ${DOCKER_USERNAME}/${DOCKER_REPONAME}:${VERSION}
  fi
  docker push ${DOCKER_USERNAME}/${DOCKER_REPONAME}
  echo
fi
if [[ ! -z ${QUAY_PASSWORD} ]]; then
  echo "### Push to Quay..."
  echo "$QUAY_PASSWORD" | docker login quay.io --username "$QUAY_LOGIN" --password-stdin > /dev/null 2>&1
  if [ "${DOCKER_TAG}" == "latest" -a "${PUSH_LATEST}" == "true" ]; then
    docker tag ${BUILD_TAG} quay.io/${QUAY_USERNAME}/${QUAY_REPONAME}:${DOCKER_TAG}
  fi
  if [[ "${VERSION}" != "latest" ]]; then
    docker tag ${BUILD_TAG} quay.io/${QUAY_USERNAME}/${QUAY_REPONAME}:${VERSION}
  fi
  docker push quay.io/${QUAY_USERNAME}/${QUAY_REPONAME}
  echo
fi
