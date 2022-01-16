# syntax=docker/dockerfile:1-labs

ARG GO_VERSION="1.17"
ARG GORELEASER_XX_VERSION="1.2.5"

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:${GORELEASER_XX_VERSION} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache file git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name "swarm-cronjob" \
    --dist "/out" \
    --main="./cmd" \
    --flags="-trimpath" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="CHANGELOG.md" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM scratch AS binary
COPY --from=build /usr/local/bin/swarm-cronjob* /

FROM alpine:3.15
RUN apk --update --no-cache add ca-certificates openssl
COPY --from=build /usr/local/bin/swarm-cronjob /usr/local/bin/swarm-cronjob
ENTRYPOINT [ "swarm-cronjob" ]
