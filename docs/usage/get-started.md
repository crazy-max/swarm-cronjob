# Get started

!!! warning
    Before starting, you must have a swarm-cronjob instance up and running using [docker](../install/docker.md)
    or an available [binary](../install/binary.md) for your platform.

When swarm-cronjob is ready, create a new stack to be scheduled like this one:

```yaml
version: "3.2"

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

You can include any configuration as long as you abide with the following conditions:

* Set `command` to run the task command
* Set `mode` to `replicated` (default)
* Set `replicas` to `0` to avoid running task as soon as the service is deployed
* Set `restart_policy.condition` to `none`. This is needed for a cronjob, otherwise the task will restart automatically
* Add [Docker labels](docker-labels.md) to tell *swarm-cronjob* that your service is a cronjob

Once ready, deploy your scheduled stack on the swarm cluster:

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

You can also use global mode services with swarm-cronjob. A typical use-case would be to remove unused data on your
nodes using `docker system prune` command periodically.

To do so, create a new global stack:

```yaml
version: "3.2"

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

Same conditions have to be applied as `replicated` mode excepted:

* Set `mode` to `global`
* Remove `replicas` field as this is only used with `replicated` mode

Once ready, deploy your global cron stack on the swarm cluster:

```shell
docker stack deploy -c global.yml global
```

You can also use replicated-job mode services with swarm-cronjob. This mode is ideal for job-based tasks that need to run with specific concurrency limits and completion tracking.

To create a replicated-job service, create a new stack:

```yaml
version: "3.8"

services:
  backup-job:
    image: busybox
    command: ["sh", "-c", "echo 'Backup job executed at:' && date && sleep 5"]
    deploy:
      mode: replicated-job
      replicas: 0
      labels:
        - "swarm.cronjob.enable=true"
        - "swarm.cronjob.schedule=0 */30 * * * *"
        - "swarm.cronjob.skip-running=false"
        - "swarm.cronjob.replicas=2"
      restart_policy:
        condition: none
```

Similar conditions apply as with other modes, but with key differences:

* Set `mode` to `replicated-job`
* **Must** set `replicas: 0` to ensure the job starts with `MaxConcurrent: 0` (required for proper initialization)
* Use `swarm.cronjob.replicas` to set the desired `MaxConcurrent` value for job execution
* Requires Docker Compose version 3.8 or higher for replicated-job support
* The `swarm.cronjob.replicas` label controls how many tasks can run concurrently when the job is triggered

!!! warning "Initial MaxConcurrent Requirement"
    Replicated-job services **must** be created with `replicas: 0` initially. This sets `MaxConcurrent: 0` in the Docker service, which swarm-cronjob requires to properly manage job scheduling. The `swarm.cronjob.replicas` label will be used to set the actual concurrency limit when jobs are triggered.

Once ready, deploy your replicated-job stack on the swarm cluster:

```shell
docker stack deploy -c backup-job.yml backup
```
