#!/usr/bin/env bash
set -e

TARGETPLATFORM=${1:-linux/amd64}
VERSION=${2:-snapshot}

case "$TARGETPLATFORM" in
  "linux/amd64")
    export GOOS=linux
    export GOARCH=amd64
    export GOARM=
    ;;
  "linux/arm/v6")
    export GOOS=linux
    export GOARCH=arm
    export GOARM=6
    ;;
  "linux/arm/v7")
    export GOOS=linux
    export GOARCH=arm
    export GOARM=7
    ;;
  "linux/arm64")
    export GOOS=linux
    export GOARCH=arm64
    export GOARM=
    ;;
  "linux/386")
    export GOOS=linux
    export GOARCH=386
    export GOARM=
    ;;
  "linux/ppc64le")
    export GOOS=linux
    export GOARCH=ppc64le
    export GOARM=
    ;;
  "linux/s390x")
    export GOOS=linux
    export GOARCH=s390x
    export GOARM=
    ;;
esac

echo "TARGETPLATFORM=${TARGETPLATFORM}"
echo "VERSION=${VERSION}"
echo "GOOS=${GOOS}"
echo "GOARCH=${GOARCH}"
echo "GOARM=${GOARM}"

export CGO_ENABLED=0
export GO111MODULE=on
export GOPROXY=https://goproxy.io

go env
go build -ldflags "-w -s -X 'main.version=${VERSION}'" -v -o swarm-cronjob cmd/main.go
