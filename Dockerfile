# syntax=docker/dockerfile:1.3
ARG GO_VERSION

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine3.14 AS base
COPY --from=goreleaser-xx / /
RUN apk add --no-cache ca-certificates gcc file git linux-headers musl-dev tar
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
ARG GIT_REF
RUN --mount=type=bind,target=/src,rw \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=target=/go/pkg/mod,type=cache \
  goreleaser-xx --debug \
    --name "swarm-cronjob" \
    --dist "/out" \
    --hooks="go mod tidy" \
    --hooks="go mod download" \
    --main="./cmd/main.go" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="CHANGELOG.md" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifacts
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM alpine:3.14

RUN apk --update --no-cache add \
  ca-certificates \
  openssl

COPY --from=build /usr/local/bin/swarm-cronjob /usr/local/bin/swarm-cronjob
RUN swarm-cronjob --version

ENTRYPOINT [ "swarm-cronjob" ]
