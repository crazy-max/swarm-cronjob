<p align="center"><a href="https://crazymax.dev/swarm-cronjob/" target="_blank"><img height="128" src="https://raw.githubusercontent.com/crazy-max/swarm-cronjob/master/.github/swarm-cronjob.png"></a></p>

<p align="center">
  <a href="https://crazymax.dev/swarm-cronjob/"><img src="https://img.shields.io/badge/doc-mkdocs-02a6f2?style=flat-square&logo=read-the-docs" alt="Documentation"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/swarm-cronjob.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/swarm-cronjob/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/swarm-cronjob/actions?workflow=build"><img src="https://img.shields.io/github/actions/workflow/status/crazy-max/swarm-cronjob/build.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/stars/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/pulls/crazymax/swarm-cronjob.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/swarm-cronjob"><img src="https://goreportcard.com/badge/github.com/crazy-max/swarm-cronjob?style=flat-square" alt="Go Report"></a>
  <a href="https://codecov.io/gh/crazy-max/swarm-cronjob"><img src="https://img.shields.io/codecov/c/github/crazy-max/swarm-cronjob?logo=codecov&style=flat-square" alt="Codecov"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**swarm-cronjob** lets you run recurring jobs on [Docker Swarm](https://docs.docker.com/engine/swarm/)
by defining cron-style schedules with [service labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l---label).
It watches your services through the Docker API, keeps its configuration in
sync automatically, and creates scheduled job runs across the cluster without
requiring a separate scheduler.

> [!TIP] 
> Want to be notified of new releases? Check out 🔔 [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun)
> project!

## Documentation

Documentation can be found on https://crazymax.dev/swarm-cronjob/

## Contributing

Want to contribute? Awesome! The most basic way to show your support is to star
the project, or to raise issues. You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max)
or by making a [PayPal donation](https://www.paypal.me/crazyws) to ensure this
journey continues indefinitely!

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.<br/>
Icon credit to [Laurel](https://twitter.com/laurelcomics).
