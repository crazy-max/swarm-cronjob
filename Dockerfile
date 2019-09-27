# syntax=docker/dockerfile:experimental
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.12.10-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN printf "I am running on ${BUILDPLATFORM:-linux/amd64}, building for ${TARGETPLATFORM:-linux/amd64}\n$(uname -a)\n"

RUN [ "$TARGETPLATFORM" = "linux/amd64"   ] && echo GOOS=linux GOARCH=amd64 > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/arm/v6"  ] && echo GOOS=linux GOARCH=arm GOARM=6 > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/arm/v7"  ] && echo GOOS=linux GOARCH=arm GOARM=7 > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/arm64"   ] && echo GOOS=linux GOARCH=arm64 > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/386"     ] && echo GOOS=linux GOARCH=386 > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/ppc64le" ] && echo GOOS=linux GOARCH=ppc64le > /tmp/.env || true
RUN [ "$TARGETPLATFORM" = "linux/s390x"   ] && echo GOOS=linux GOARCH=s390x > /tmp/.env || true
RUN env $(cat /tmp/.env | xargs) go env

RUN apk --update --no-cache add \
    build-base \
    gcc \
    git \
  && rm -rf /tmp/* /var/cache/apk/*

WORKDIR /app

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
COPY go.mod .
COPY go.sum .
RUN env $(cat /tmp/.env | xargs) go mod download
COPY . ./

ARG VERSION=dev
RUN env $(cat /tmp/.env | xargs) go build -ldflags "-w -s -X 'main.version=${VERSION}'" -v -o swarm-cronjob cmd/main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest

LABEL maintainer="CrazyMax" \
  org.label-schema.name="swarm-cronjob" \
  org.label-schema.description="Create jobs on a time-based schedule on Swarm" \
  org.label-schema.url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vcs-url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vendor="CrazyMax" \
  org.label-schema.schema-version="1.0"

RUN apk --update --no-cache add \
    ca-certificates \
    libressl \
    shadow \
  && addgroup -g 1000 swarm-cronjob \
  && adduser -u 1000 -G swarm-cronjob -s /sbin/nologin -D swarm-cronjob \
  && rm -rf /tmp/* /var/cache/apk/*

COPY --from=builder /app/swarm-cronjob /usr/local/bin/swarm-cronjob
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
RUN swarm-cronjob --version

USER swarm-cronjob

ENTRYPOINT [ "swarm-cronjob" ]
