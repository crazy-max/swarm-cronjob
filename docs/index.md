<img src="assets/logo.png" alt="swarm-cronjob" width="128px" style="display: block; margin-left: auto; margin-right: auto"/>

<p align="center">
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/swarm-cronjob.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/swarm-cronjob/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/actions?workflow=build"><img src="https://img.shields.io/github/actions/workflow/status/crazy-max/swarm-cronjob/build.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/stars/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/pulls/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/swarm-cronjob"><img src="https://goreportcard.com/badge/github.com/crazy-max/swarm-cronjob?style=flat-square" alt="Go Report"></a>
  <a href="https://codecov.io/gh/crazy-max/swarm-cronjob"><img src="https://img.shields.io/codecov/c/github/crazy-max/swarm-cronjob?logo=codecov&style=flat-square" alt="Codecov"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate PayPal"></a>
</p>

---

## What is swarm-cronjob?

**swarm-cronjob** is a scheduler for [Docker Swarm](https://docs.docker.com/engine/swarm/).
You define a cron-style schedule with [service labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l---label),
and swarm-cronjob turns that service into a recurring job. It watches the
Docker API, keeps its configuration in sync as services are added or updated,
and runs scheduled workloads directly inside your Swarm cluster.

## Features

* Label-driven configuration, define schedules directly on Swarm services
* Automatic service discovery and configuration reloads, no restart required
* Cron-style scheduling for recurring jobs inside your Swarm cluster
* Overlap control, can skip a run if the target service is already running
* Configurable scheduler time zone

## License

This project is licensed under the terms of the MIT license.<br />
Icon credit to [Laurel](https://twitter.com/laurelcomics).
