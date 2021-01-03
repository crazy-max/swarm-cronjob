# FAQ

## Timezones

By default, all interpretation and scheduling is done with your local timezone (`TZ` environment variable).

Individual cron schedules may also override the time zone they are to be interpreted in by providing an additional
space-separated field at the beginning of the cron spec, of the form `CRON_TZ=Asia/Tokyo`.

For example:

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
        - "swarm.cronjob.schedule=CRON_TZ=Asia/Tokyo * * * * *"
        - "swarm.cronjob.skip-running=false"
      restart_policy:
        condition: none
```
