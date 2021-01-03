module github.com/crazy-max/swarm-cronjob

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/Microsoft/hcsshim/test v0.0.0-20201218223536-d3e5debf77da // indirect
	github.com/Shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d // indirect
	github.com/agl/ed25519 v0.0.0-00010101000000-000000000000 // indirect
	github.com/alecthomas/kong v0.2.12
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bugsnag/bugsnag-go v1.8.0 // indirect
	github.com/bugsnag/panicwrap v1.2.2 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cloudflare/cfssl v1.5.0 // indirect
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/docker/cli v0.0.0-20200303215952-eb310fca4956
	github.com/docker/docker v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go v1.5.1-1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/fvbommel/sortorder v1.0.2 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/mapstructure v1.4.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/opencontainers/selinux v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/zerolog v1.20.0
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/theupdateframework/notary v0.6.1 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/grpc v1.34.0 // indirect
	gopkg.in/dancannon/gorethink.v3 v3.0.5 // indirect
	gopkg.in/fatih/pool.v2 v2.0.0 // indirect
	gopkg.in/gorethink/gorethink.v3 v3.0.5 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

// github.com/docker/engine v19.03.8
replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200309214505-aa6a9891b09c

// latest docker/distribution
replace github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible

// latest commit before the author shutdown the repo
replace github.com/agl/ed25519 => github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
