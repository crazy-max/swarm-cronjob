<img src="assets/logo.png" alt="swarm-cronjob" width="128px" style="display: block; margin-left: auto; margin-right: auto"/>

<p align="center">
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/swarm-cronjob.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/swarm-cronjob/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/actions?workflow=build"><img src="https://img.shields.io/github/workflow/status/crazy-max/swarm-cronjob/build?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/stars/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/pulls/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/swarm-cronjob"><img src="https://goreportcard.com/badge/github.com/crazy-max/swarm-cronjob?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/swarm-cronjob"><img src="https://img.shields.io/codacy/grade/1edb80b0f97b4195b7bb50cfb35a37d2.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

---

## What is swarm-cronjob?

**swarm-cronjob** creates jobs on a time-based schedule on [Swarm](https://docs.docker.com/engine/swarm/) with a
dedicated service in a distributed manner that configures itself automatically and dynamically through
[labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l---label)
and Docker API.

## Features

* Continuously updates its configuration (no restart)
* Cron implementation through go routines
* Allow to skip a job if the service is currently running
* Timezone can be changed for the scheduler

## License

This project is licensed under the terms of the MIT license.<br />
Icon credit to [Laurel](https://twitter.com/laurelcomics).
