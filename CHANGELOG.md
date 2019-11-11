# Changelog

## 1.7.1 (2019/11/11)

* Update libs
* Cache go modules

## 1.7.0 (2019/10/30)

* Seconds field is now optional
* Docker client v19.03.4

## 1.6.0 (2019/10/13)

* Allow to set more replicas (Issue #16)
* Docker client v19.03.3
* Update libs

## 1.5.0 (2019/09/27)

* Update libs
* Go 1.12.10

## 1.4.0 (2019/09/22)

* Log removed/disabled services
* Docker client v19.03.2
* Use GOPROXY
* :warning: Stop publishing Docker image on Quay
* Multi-platform Docker image
* Switch to GitHub Actions
* Add instructions to create a Linux service

## 1.3.0 (2019/07/19)

* Docker client v18.09.8

## 1.3.0-beta.1 (2019/07/18)

* Add support for global mode (Issue #7)
* Use v3 robfig/cron
* Docker client v18.09.7

## 1.2.1 (2019/05/30)

* Fix nil pointer (Issue #7)

## 1.2.0 (2019/05/01)

* Skip completed tasks while checking status (Issue #4)
* Update Docker client and some libs
* Go 1.12.4

## 1.1.0 (2019/03/21)

* Go 1.12.1

## 1.0.0 (2019/02/17)

* Add JSON log output
* Deliver artifacts through goreleaser
* Review project structure

## 0.2.1 (2019/01/24)

* Go 1.11.5
* Update `go.sum` after go@1.11.4 symlink fix (golang/go#29278)

## 0.2.0 (2019/01/22)

* Add support for Docker API [1.38](https://docs.docker.com/engine/api/v1.38/) (Issue #3)
* ldflags -X not properly applied

## 0.1.2 (2019/01/14)

* Fix non-cronjob services added to cronjob list (Issue #2)
* Handle removed services
* NPE while checking service

## 0.1.1 (2018/12/13)

* Fix build args
* Checksum mismatch on Go 1.11.4

## 0.1.0 (2018/12/13)

* Initial version based on Docker API [1.26](https://docs.docker.com/engine/api/v1.26/)
