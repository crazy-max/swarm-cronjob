# Command Line

## Usage

```shell
$ swarm-cronjob [options]
```

## Options

```
$ swarm-cronjob --help
Usage: swarm-cronjob

Create jobs on a time-based schedule on Swarm. More info:
https://github.com/crazy-max/swarm-cronjob

Flags:
  --help                Show context-sensitive help.
  --version
  --timezone="UTC"      Timezone assigned to swarm-cronjob ($TZ).
  --log-level="info"    Set log level ($LOG_LEVEL).
  --log-json            Enable JSON logging output ($LOG_JSON).
```

## Environment variables

Following environment variables can be used in place:

| Name               | Default       | Description   |
|--------------------|---------------|---------------|
| `TZ`               | `UTC`         | Timezone assigned |
| `LOG_LEVEL`        | `info`        | Log level output |
| `LOG_JSON`         | `false`       | Enable JSON logging output |
