# Docker labels

You can configure your service using swarm-cronjob through Docker labels:

| Name                           | Default | Description                                                                                                        |
| ------------------------------ | ------- | ------------------------------------------------------------------------------------------------------------------ |
| `swarm.cronjob.enable`         |         | Set to true to enable the cronjob. **required**                                                                    |
| `swarm.cronjob.schedule`       |         | [CRON expression format](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) to use. **required** |
| `swarm.cronjob.skip-running`   | `false` | Do not start a job if the service is currently running.                                                            |
| `swarm.cronjob.replicas`       | `1`     | Number of replicas to set on schedule in `replicated` mode, or `MaxConcurrent` limit in `replicated-job` mode.     |
| `swarm.cronjob.registry-auth`  | `false` | Send registry authentication details to Swarm agents.                                                              |
| `swarm.cronjob.query-registry` |         | Indicates whether the service update requires contacting a registry                                                |
