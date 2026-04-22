# Get started

!!! warning
    Before starting, you must have a swarm-cronjob instance up and running using [docker](../install/docker.md)
    or a [binary](../install/binary.md) for your platform.

When swarm-cronjob is ready, create a stack to be scheduled like this:

```yaml
services:
  test:
    image: busybox
    command: date
    deploy:
      mode: replicated
      replicas: 0
      labels:
        - "swarm.cronjob.enable=true"
        - "swarm.cronjob.schedule=* * * * *"
        - "swarm.cronjob.skip-running=false"
      restart_policy:
        condition: none
```

You can include any configuration as long as it meets the following conditions:

* Set `command` to run the task
* Set `mode` to `replicated` (default)
* Set `replicas` to `0` to avoid running the task as soon as the service is deployed
* Set `restart_policy.condition` to `none`. This is required for a scheduled job, otherwise the task will restart automatically
* Add [Docker labels](docker-labels.md) to tell *swarm-cronjob* that your service is scheduled

Once ready, deploy your scheduled stack on the Swarm cluster:

`docker stack deploy -c date.yml date`

!!! example "Logs"
    ```
    $ docker service logs swarm_cronjob_app
    swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Starting swarm-cronjob v1.2.0
    swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Add cronjob with schedule * * * * * service=date_test
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

You can also use global-mode services with swarm-cronjob. A typical use case is to remove unused data from your
nodes periodically with the `docker system prune` command.

To do so, create a new global stack:

```yaml
services:
  prune-nodes:
    image: docker
    command: ["docker", "system", "prune", "-f"]
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
    deploy:
      mode: global
      labels:
        - "swarm.cronjob.enable=true"
        - "swarm.cronjob.schedule=0 */5 * * * *"
        - "swarm.cronjob.skip-running=false"
      restart_policy:
        condition: none
```

The same conditions apply as in `replicated` mode, except:

* Set `mode` to `global`
* Remove the `replicas` field, it is only used with `replicated` mode

Once ready, deploy your global cron stack on the Swarm cluster:

```shell
docker stack deploy -c global.yml global
```
