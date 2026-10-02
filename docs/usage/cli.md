# Command line

## Usage

```shell
swarm-cronjob [options]
```

## Options

```
$ swarm-cronjob --help
Usage: swarm-cronjob [flags]

Create jobs on a time-based schedule on Swarm. More info:
https://github.com/crazy-max/swarm-cronjob

Flags:
  -h, --help                Show context-sensitive help.
      --version
      --log-level="info"    Set log level ($LOG_LEVEL).
      --log-json            Enable JSON logging output ($LOG_JSON).
      --service-label=SERVICE-LABEL,...
                            Additional service label filter. Repeat for multiple
                            labels ($SERVICE_LABELS).
```

## Environment variables

The following environment variables can be used instead:

| Name        | Default | Description                |
|-------------|---------|----------------------------|
| `LOG_LEVEL` | `info`  | Log level                  |
| `LOG_JSON`  | `false` | Enable JSON logging output |
| `SERVICE_LABELS` | | Comma separated service label filters |

## Service label filters

Use `--service-label` to select the services an instance schedules.
Filters accept a label name or `name=value`. Every filter must match.
The existing cronjob labels are still required.

For example, limit scheduling to one Docker stack:

```shell
swarm-cronjob --service-label com.docker.stack.namespace=production
```

Repeat the flag to require another label:

```shell
swarm-cronjob --service-label com.docker.stack.namespace=production --service-label scheduler=one
```

Set `SERVICE_LABELS=com.docker.stack.namespace=production,scheduler=one` for the same selection through environment configuration.
Command line filters replace environment filters.
Services that lose a required label are removed from the schedule when their update event arrives or during periodic reconciliation.
Without filters, the instance continues scheduling cronjobs across the whole cluster.
