# syntax=docker/dockerfile:experimental
FROM --platform=amd64 golang:1.12.4 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go version
RUN go mod download
COPY . ./
RUN cp /usr/local/go/lib/time/zoneinfo.zip ./

ARG TARGETPLATFORM
ARG VERSION

RUN bash gobuild.sh ${TARGETPLATFORM} ${VERSION}

FROM --platform=$TARGETPLATFORM scratch

LABEL maintainer="CrazyMax" \
  org.label-schema.name="swarm-cronjob" \
  org.label-schema.description="Create jobs on a time-based schedule on Swarm" \
  org.label-schema.url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vcs-url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vendor="CrazyMax" \
  org.label-schema.schema-version="1.0"

COPY --from=builder /app/swarm-cronjob /usr/local/bin/swarm-cronjob
COPY --from=builder /app/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
CMD [ "swarm-cronjob" ]
