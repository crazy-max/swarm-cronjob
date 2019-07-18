<p align="center"><a href="https://github.com/crazy-max/swarm-cronjob" target="_blank"><img height="292"src="https://raw.githubusercontent.com/crazy-max/swarm-cronjob/master/.res/swarm-cronjob.jpg"></a></p>

<p align="center">
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/badge/dynamic/json.svg?label=version&query=$.results[1].name&url=https://hub.docker.com/v2/repositories/crazymax/swarm-cronjob/tags&style=flat-square" alt="Latest Version"></a>
  <a href="https://travis-ci.com/crazy-max/swarm-cronjob"><img src="https://img.shields.io/travis/com/crazy-max/swarm-cronjob/master.svg?style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/stars/crazymax/swarm-cronjob.svg?style=flat-square" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/swarm-cronjob/"><img src="https://img.shields.io/docker/pulls/crazymax/swarm-cronjob.svg?style=flat-square" alt="Docker Pulls"></a>
  <a href="https://quay.io/repository/crazymax/swarm-cronjob"><img src="https://quay.io/repository/crazymax/swarm-cronjob/status?style=flat-square" alt="Docker Repository on Quay"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/swarm-cronjob"><img src="https://goreportcard.com/badge/github.com/crazy-max/swarm-cronjob?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/swarm-cronjob"><img src="https://img.shields.io/codacy/grade/1edb80b0f97b4195b7bb50cfb35a37d2.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=YZ64W5KHJGGZG"><img src="https://img.shields.io/badge/donate-paypal-7057ff.svg?style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**swarm-cronjob** creates jobs on a time-based schedule on [Swarm](https://docs.docker.com/engine/swarm/) with a dedicated service in a distributed manner that configures itself automatically and dynamically through [labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l---label) and Docker API.

## Features

* Continuously updates its configuration (no restart)
* Cron implementation through go routines
* Allow to skip a job if the service is currently running
* Timezone can be changed for the scheduler

## Docker

### Environment variables

* `TZ` : The timezone assigned to the scheduler (default `UTC`)
* `LOG_LEVEL` : Log level (default `info`)
* `LOG_JSON` : Enable JSON logging output (default `false`)

## Quickstart

### Swarm cluster

Create a service that uses the swarm-cronjob image :

```
$ docker service create --name swarm_cronjob \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  --env "LOG_LEVEL=info" \
  --env "LOG_JSON=false" \
  --constraint "node.role == manager" \
  crazymax/swarm-cronjob
```

Alternatively, you can deploy the stack [swarm_cronjob.yml](.res/example/swarm_cronjob.yml) :

`docker stack deploy -c swarm_cronjob.yml swarm_cronjob`

Now that we have a swarm-cronjob instance up and running, we will deploy new services.

Create a new stack based on [this one (date)](.res/example/date.yml). You can include any configuration as long as you abide with the following conditions :

* Set `command` to run the task command
* Set `mode` to `replicated` (default)
* Set `replicas` to `0` to avoid running task as soon as the service is deployed
* Set `restart_policy.condition` to `none`. This is needed for a cronjob, otherwise the task will restart automatically
* Add labels to tell *swarm-cronjob* that your service is a cronjob :
  * `swarm.cronjob.enable` : Set to true to enable the cronjob (**required**)
  * `swarm.cronjob.schedule` : [CRON expression format](https://godoc.org/github.com/crazy-max/cron#hdr-CRON_Expression_Format) to use (**required**)
  * `swarm.cronjob.skip-running` : Do not start a job if the service is currently running (**optional**)

Once ready, deploy your cron stack on the swarm cluster :

`docker stack deploy -c date.yml date`

You can also use global mode services with swarm-cronjob. A typical use-case would be to remove unused data on your nodes using `docker system prune` command periodically.

To do so, create a new stack based on [this one (global)](.res/example/global.yml). Same conditions have to be applied as `replicated` mode excepted :

* Set `mode` to `global`
* Remove `replicas` field as this is only used with `replicated` mode

Once ready, deploy your global cron stack on the swarm cluster :

`docker stack deploy -c global.yml global`

> :bulb: More examples can be found [here](.res/example)

### Without Docker

swarm-cronjob binaries are available in [releases](https://github.com/crazy-max/swarm-cronjob/releases) page.

Choose the archive matching the destination platform and extract swarm-cronjob:

```
$ cd /opt
$ wget -qO- https://github.com/crazy-max/swarm-cronjob/releases/download/v1.1.0/swarm-cronjob_1.1.0_linux_x86_64.tar.gz | tar -zxvf - swarm-cronjob
```

After getting the binary, it can be tested with `./swarm-cronjob --help` or moved to a permanent location.
When launched manually, swarm-cronjob can be killed using `Ctrl+C`:

```
$ ./swarm-cronjob --help
usage: swarm-cronjob [<flags>]

Create jobs on a time-based schedule on Swarm. More info on
https://github.com/crazy-max/swarm-cronjob

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man).
  --timezone="UTC"    Timezone assigned to the scheduler.
  --log-level="info"  Set log level.
  --log-json          Enable JSON logging output.
  --version           Show application version.
```

## Logs

Here is a sample output:

```
$ docker service logs swarm_cronjob_app
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Starting swarm-cronjob v1.2.0
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Add cronjob with schedule 0/10 * * * * * service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:05:00 UTC INF Start job last_status=n/a service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:06:00 UTC INF Start job last_status=n/a service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:07:00 UTC INF Start job last_status=n/a service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:08:00 UTC INF Start job last_status=n/a service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:09:00 UTC INF Start job last_status=n/a service=date_test
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:10:00 UTC INF Start job last_status=n/a service=date_test
$ docker service logs date_test
date_test.1.o1d5mn4gjff3@default    | Thu Dec 13 20:11:01 UTC 2018
date_test.1.5askx244las2@default    | Thu Dec 13 20:09:00 UTC 2018
date_test.1.4lz5ez2waekk@default    | Thu Dec 13 20:12:00 UTC 2018
date_test.1.135qzpxd1ui3@default    | Thu Dec 13 20:13:01 UTC 2018
date_test.1.hngject056n3@default    | Thu Dec 13 20:10:00 UTC 2018
```

## How can I help ?

All kinds of contributions are welcome :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Paypal](.res/paypal-donate.png)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=YZ64W5KHJGGZG)

## License

MIT. See `LICENSE` for more details.
