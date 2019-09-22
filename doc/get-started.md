# Get started

> :warning: Before starting, you must have a swarm-cronjob instance up and running using [docker](install/docker.md) or an available [binary](install/binary.md) for your platform.

When swarm-cronjob is ready, create a new stack based on [this one (date)](../.res/example/date.yml). You can include any configuration as long as you abide with the following conditions:

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

To do so, create a new stack based on [this one (global)](../.res/example/global.yml). Same conditions have to be applied as `replicated` mode excepted :

* Set `mode` to `global`
* Remove `replicas` field as this is only used with `replicated` mode

Once ready, deploy your global cron stack on the swarm cluster :

`docker stack deploy -c global.yml global`

> :bulb: More examples can be found [here](../.res/example)
