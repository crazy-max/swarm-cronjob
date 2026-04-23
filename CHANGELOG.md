# Changelog

## 1.16.0 (2026/04/23)

* Skip cron runs while services are reconciling by @crazy-max in #472
* Fix missed schedule updates with periodic cronjob reconciliation by @crazy-max in #473
* Migrate to Moby client and api modules, updating to 29.4.1 by @crazy-max in #466
* Simplify shutdown lifecycle by @crazy-max in #471
* Go 1.26 by @crazy-max in #459
* MkDocs Materials 9.7.5 by @crazy-max in #469
* Bump github.com/alecthomas/kong to 1.15.0 in #452
* Bump github.com/go-viper/mapstructure/v2 to 2.5.0 in #421
* Bump github.com/rs/zerolog to 1.35.1 in #443 #475
* Bump golang.org/x/sys to 0.43.0 in #457

**Full Changelog**: [`v1.15.0...v1.16.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.15.0...v1.16.0)

## 1.15.0 (2025/12/31)

* Go 1.25 by @crazy-max in #406
* Alpine Linux 3.23 by @crazy-max in #410
* Switch to `github.com/go-viper/mapstructure` by @crazy-max in #413
* Bump github.com/alecthomas/kong to 1.10.0 in #400
* Bump github.com/docker/cli to 28.5.2 by @crazy-max in #411
* Bump github.com/docker/docker to 28.5.2 by @crazy-max in #411
* Bump github.com/prometheus/client_golang to 1.11.1 in #379
* Bump github.com/rs/zerolog to 1.34.0 in #399
* Bump golang.org/x/crypto to 0.45.0 in #408
* Bump golang.org/x/sys to 0.39.0 in #396

**Full Changelog**: [`v1.14.0...v1.15.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.14.0...v1.15.0)

## 1.14.0 (2024/12/24)

* Add `tzdata` package to Docker image by @fl4shback in #337
* Go 1.23 by @crazy-max in #372
* Alpine Linux 3.21 by @crazy-max in #372
* Bump github.com/alecthomas/kong to 1.6.0 in #325 #374
* Bump github.com/distribution/reference to 0.6.0 in #373
* Bump github.com/docker/cli to 27.4.1+incompatible by @crazy-max in #378
* Bump github.com/docker/docker to 27.4.1+incompatible by @crazy-max in #378
* Bump github.com/rs/zerolog to 1.33.0 in #320 #345
* Bump golang.org/x/sys to 0.28.0 in #333 #346 #375

**Full Changelog**: [`v1.13.0...v1.14.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.13.0...v1.14.0)

## 1.13.0 (2024/02/01)

* Enable automatic API version negotiation for the Docker client by @crazy-max in #313
* Go 1.21 by @crazy-max in #285 #284
* Alpine Linux 3.19 by @crazy-max in #317
* Bump github.com/alecthomas/kong to 0.8.1 in #289
* Bump github.com/distribution/distribution to 2.8.3+incompatible in #259 #260 #302
* Bump github.com/docker/cli to 24.0.7+incompatible in #250 #290
* Bump github.com/docker/docker to 24.0.7+incompatible in #250 #290
* Bump github.com/opencontainers/image-spec to 1.0.2 in #239
* Bump github.com/prometheus/client_golang to 1.11.1 in #234
* Bump github.com/rs/zerolog to 1.31.0 in #245 #288
* Bump golang.org/x/crypto to 0.17.0 in #237 #303
* Bump golang.org/x/net to 0.17.0 in #238 #287
* Bump golang.org/x/sys to 0.16.0 in #252 #279 #301 #306

**Full Changelog**: [`v1.12.0...v1.13.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.12.0...v1.13.0)

## 1.12.0 (2023/02/14)

* Go 1.19 by @crazy-max in #214
* Alpine Linux 3.17 by @crazy-max in #223
* Improve workflow by @crazy-max in #215
* Bump github.com/alecthomas/kong to 0.7.1 in #220
* Bump github.com/docker/cli to 20.10.22+incompatible in #221
* Bump github.com/docker/distribution to 2.8.0+incompatible in #233
* Bump github.com/docker/docker to 20.10.22+incompatible in #222
* Bump github.com/rs/zerolog to 1.29.0 in #207 #227
* Bump golang.org/x/sys to 0.5.0 by @crazy-max in #224 #230

**Full Changelog**: [`v1.11.0...v1.12.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.11.0...v1.12.0)

## 1.11.0 (2022/07/17)

* Add an option to query the registry on service update by @crazy-max in #201
* Fix a possible nil pointer in `ElectAuthServer` from Docker CLI by @crazy-max in #202
* Go 1.18 by @crazy-max in #204
* Alpine Linux 3.16 by @crazy-max in #166 #206
* goreleaser-xx 1.2.5
* Move `syscall` to `golang.org/x/sys`
* MkDocs Material 8.3.9 by @crazy-max in #205
* Improve Dockerfiles by @crazy-max in #163
* Bump github.com/alecthomas/kong to 0.6.1 in #152 #162 #165 #176 #200
* Bump github.com/docker/cli to 20.10.17 in #149 #160 #197
* Bump github.com/docker/docker to 20.10.17 in #148 #159 #196
* Bump github.com/mitchellh/mapstructure to 1.5.0 in #145 #155 #181
* Bump github.com/rs/zerolog to 1.27.0 in #144 #150 #162 #198

**Full Changelog**: [`v1.10.0...v1.11.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.10.0...v1.11.0)

## 1.10.0 (2021/09/05)

* Docker client v20.10.8 in #109 #110 #134 #135
* Go 1.17 by @crazy-max in #114 #140
* Add `darwin/amd64`, `darwin/arm64`, `linux/riscv64`, `windows/arm64` artifacts by @crazy-max in #141
* Alpine Linux 3.14
* MkDocs Materials 7.2.6 by @crazy-max in #143
* Remove `linux/s390x` Docker platform support (for now)
* Switch to goreleaser-xx by @crazy-max in #111
* Bump github.com/alecthomas/kong to 0.2.16 in #106 #112
* Bump github.com/mitchellh/mapstructure to 1.4.1 in #103
* Bump github.com/rs/zerolog to 1.24.0 in #131 #137 #142

**Full Changelog**: [`v1.9.0...v1.10.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.9.0...v1.10.0)

## 1.9.0 (2021/01/03)

* Refactor CI and development workflow with Buildx Bake by @crazy-max in #99
    * Upload artifacts
    * Add `image-local` target
    * Single job for artifacts and image
    * Add `armv5` artifact
* Use embedded tzdata package and remove `--timezone` flag by @crazy-max in #98
* Go 1.15 by @crazy-max in #97
* Send registry authentication details to Swarm agents by @crazy-max in #96
* Docker client v20.10.1
* Remove support for `freebsd/*` (moby/moby#38818)
* Handle registry auth from the spec by @crazy-max in #92
* Docker image also available on [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/swarm-cronjob)
* Add docs website with MkDocs
* Add notes about time zones by @clburlison in #43
* Add renovate example by @decentral1se in #42
* Add MariaDB dump example in #35
* Bump github.com/alecthomas/kong to 0.2.12 in #89
* Bump github.com/mitchellh/mapstructure to 1.4.0 in #88
* Bump github.com/rs/zerolog to 1.20.0 in #68

**Full Changelog**: [`v1.8.0...v1.9.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.8.0...v1.9.0)

## 1.8.0 (2020/04/06)

* Switch to the Kong command-line parser
* Go 1.13
* Docker client v19.03.8
* Use Open Container Specification labels, label-schema.org labels are deprecated
* Update dependencies

**Full Changelog**: [`v1.7.1...v1.8.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.7.1...v1.8.0)

## 1.7.1 (2019/11/11)

* Update dependencies
* Cache Go modules

**Full Changelog**: [`v1.7.0...v1.7.1`](https://github.com/crazy-max/swarm-cronjob/compare/v1.7.0...v1.7.1)

## 1.7.0 (2019/10/30)

* Seconds field is now optional
* Docker client v19.03.4

**Full Changelog**: [`v1.6.0...v1.7.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.6.0...v1.7.0)

## 1.6.0 (2019/10/13)

* Allow setting more replicas in #16
* Docker client v19.03.3
* Update dependencies

**Full Changelog**: [`v1.5.0...v1.6.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.5.0...v1.6.0)

## 1.5.0 (2019/09/27)

* Update dependencies
* Go 1.12.10

**Full Changelog**: [`v1.4.0...v1.5.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.4.0...v1.5.0)

## 1.4.0 (2019/09/22)

* Log removed/disabled services
* Docker client v19.03.2
* Use GOPROXY
* :warning: Stop publishing Docker image on Quay
* Multi-platform Docker image
* Switch to GitHub Actions
* Add instructions to create a Linux service

**Full Changelog**: [`v1.3.0...v1.4.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.3.0...v1.4.0)

## 1.3.0 (2019/07/19)

* Docker client v18.09.8

**Full Changelog**: [`v1.3.0-beta.1...v1.3.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.3.0-beta.1...v1.3.0)

## 1.3.0-beta.1 (2019/07/18)

* Add support for global mode in #7
* Use v3 robfig/cron
* Docker client v18.09.7

**Full Changelog**: [`v1.2.1...v1.3.0-beta.1`](https://github.com/crazy-max/swarm-cronjob/compare/v1.2.1...v1.3.0-beta.1)

## 1.2.1 (2019/05/30)

* Fix nil pointer in #7

**Full Changelog**: [`v1.2.0...v1.2.1`](https://github.com/crazy-max/swarm-cronjob/compare/v1.2.0...v1.2.1)

## 1.2.0 (2019/05/01)

* Skip completed tasks while checking status in #4
* Update Docker client and some libraries
* Go 1.12.4

**Full Changelog**: [`v1.1.0...v1.2.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.1.0...v1.2.0)

## 1.1.0 (2019/03/21)

* Go 1.12.1

**Full Changelog**: [`v1.0.0...v1.1.0`](https://github.com/crazy-max/swarm-cronjob/compare/v1.0.0...v1.1.0)

## 1.0.0 (2019/02/17)

* Add JSON log output
* Deliver artifacts through goreleaser
* Review project structure

**Full Changelog**: [`v0.2.1...v1.0.0`](https://github.com/crazy-max/swarm-cronjob/compare/v0.2.1...v1.0.0)

## 0.2.1 (2019/01/24)

* Go 1.11.5
* Update `go.sum` after symlink fix in Go 1.11.4

**Full Changelog**: [`v0.2.0...v0.2.1`](https://github.com/crazy-max/swarm-cronjob/compare/v0.2.0...v0.2.1)

## 0.2.0 (2019/01/22)

* Add support for Docker API [1.38](https://docs.docker.com/engine/api/v1.38/) in #3
* Fix `ldflags -X` not being applied properly

**Full Changelog**: [`v0.1.2...v0.2.0`](https://github.com/crazy-max/swarm-cronjob/compare/v0.1.2...v0.2.0)

## 0.1.2 (2019/01/14)

* Fix non-cronjob services added to the cronjob list in #2
* Handle removed services
* Fix an NPE while checking a service

**Full Changelog**: [`v0.1.1...v0.1.2`](https://github.com/crazy-max/swarm-cronjob/compare/v0.1.1...v0.1.2)

## 0.1.1 (2018/12/13)

* Fix build arguments
* Checksum mismatch on Go 1.11.4

**Full Changelog**: [`v0.1.0...v0.1.1`](https://github.com/crazy-max/swarm-cronjob/compare/v0.1.0...v0.1.1)

## 0.1.0 (2018/12/13)

* Initial version based on Docker API [1.26](https://docs.docker.com/engine/api/v1.26/)
