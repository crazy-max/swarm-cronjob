# Docker labels

You can configure a service for swarm-cronjob with Docker labels:

| Name                           | Default | Description                                                                                                        |
|--------------------------------|---------|--------------------------------------------------------------------------------------------------------------------|
| `swarm.cronjob.enable`         |         | Set to `true` to enable scheduling for the service. **required**                                                   |
| `swarm.cronjob.schedule`       |         | [Cron expression format](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to use. **required** |
| `swarm.cronjob.skip-running`   | `false` | Do not start a new job if the service is already running.                                                          |
| `swarm.cronjob.replicas`       | `1`     | Number of replicas to set for a scheduled run in `replicated` mode.                                                |
| `swarm.cronjob.registry-auth`  | `false` | Send registry authentication details to Swarm agents.                                                              |
| `swarm.cronjob.query-registry` |         | Whether the service update requires contacting a registry.                                                         |
