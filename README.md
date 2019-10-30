<p align="center"><a href="https://github.com/crazy-max/swarm-cronjob" target="_blank"><img height="292" src="https://raw.githubusercontent.com/crazy-max/swarm-cronjob/master/.res/swarm-cronjob.jpg"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/swarm-cronjob.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/swarm-cronjob/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/actions"><img src="https://github.com/crazy-max/swarm-cronjob/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/stars/crazymax/swarm-cronjob.svg?style=flat-square" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/pulls/crazymax/swarm-cronjob.svg?style=flat-square" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/swarm-cronjob"><img src="https://goreportcard.com/badge/github.com/crazy-max/swarm-cronjob?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/swarm-cronjob"><img src="https://img.shields.io/codacy/grade/1edb80b0f97b4195b7bb50cfb35a37d2.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://www.patreon.com/crazymax"><img src="https://img.shields.io/badge/donate-patreon-f96854.svg?logo=patreon&style=flat-square" alt="Support me on Patreon"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**swarm-cronjob** creates jobs on a time-based schedule on [Swarm](https://docs.docker.com/engine/swarm/) with a dedicated service in a distributed manner that configures itself automatically and dynamically through [labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l---label) and Docker API.

💡 Want to be notified of new releases? Check out 🔔 [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun) project!

## Features

* Continuously updates its configuration (no restart)
* Cron implementation through go routines
* Allow to skip a job if the service is currently running
* Timezone can be changed for the scheduler

## Documentation

* [Get started](doc/get-started.md)
* Install
  * [With Docker](doc/install/docker.md)
  * [From binary](doc/install/binary.md)
  * [Linux service](doc/install/linux-service.md)
* [Logs](doc/logs.md)

## How can I help ?

All kinds of contributions are welcome :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Support me on Patreon](.res/patreon.png)](https://www.patreon.com/crazymax) 
[![Paypal Donate](.res/paypal.png)](https://www.paypal.me/crazyws)

## License

MIT. See `LICENSE` for more details.
