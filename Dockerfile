FROM --platform=${BUILDPLATFORM:-linux/amd64} tonistiigi/xx:golang AS xgo
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.13-alpine AS builder

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION=dev

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io,direct
COPY --from=xgo / /

RUN apk --update --no-cache add \
    build-base \
    gcc \
    git \
  && rm -rf /tmp/* /var/cache/apk/*

WORKDIR /app

COPY . ./
RUN go mod download

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
RUN go env
RUN go build -ldflags "-w -s -X 'main.version=${VERSION}'" -v -o swarm-cronjob cmd/main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

LABEL maintainer="CrazyMax" \
  org.opencontainers.image.created=$BUILD_DATE \
  org.opencontainers.image.url="https://github.com/crazy-max/swarm-cronjob" \
  org.opencontainers.image.source="https://github.com/crazy-max/swarm-cronjob" \
  org.opencontainers.image.version=$VERSION \
  org.opencontainers.image.revision=$VCS_REF \
  org.opencontainers.image.vendor="CrazyMax" \
  org.opencontainers.image.title="swarm-cronjob" \
  org.opencontainers.image.description="Create jobs on a time-based schedule on Swarm" \
  org.opencontainers.image.licenses="MIT"

RUN apk --update --no-cache add \
    ca-certificates \
    libressl \
  && rm -rf /tmp/* /var/cache/apk/*

COPY --from=builder /app/swarm-cronjob /usr/local/bin/swarm-cronjob
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
RUN swarm-cronjob --version

ENTRYPOINT [ "swarm-cronjob" ]
