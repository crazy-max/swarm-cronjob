ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

FROM golang:1.11 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./
RUN cp /usr/local/go/lib/time/zoneinfo.zip ./ \
  && CGO_ENABLED=0 GOOS=linux go build \
  -ldflags "-w -s -X 'github.com/crazy-max/swarm-cronjob/app/config.AppVersion=${VERSION}'" \
  -v -o swarm-cronjob

FROM scratch

LABEL maintainer="CrazyMax" \
  org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.name="swarm-cronjob" \
  org.label-schema.description="Create jobs on a time-based schedule on Swarm" \
  org.label-schema.version=$VERSION \
  org.label-schema.url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.vcs-url="https://github.com/crazy-max/swarm-cronjob" \
  org.label-schema.vendor="CrazyMax" \
  org.label-schema.schema-version="1.0"

COPY --from=builder /app/swarm-cronjob /usr/local/bin/swarm-cronjob
COPY --from=builder /app/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
CMD [ "swarm-cronjob" ]
