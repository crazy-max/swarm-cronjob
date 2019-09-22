# Installation from binary

## Download

swarm-cronjob binaries are available in [releases](https://github.com/crazy-max/swarm-cronjob/releases) page.

Choose the archive matching the destination platform and extract swarm-cronjob:

```
wget -qO- https://github.com/crazy-max/swarm-cronjob/releases/download/1.3.0/swarm-cronjob_1.3.0_linux_x86_64.tar.gz | tar -zxvf - swarm-cronjob
```

## Test

After getting the binary, it can be tested with `./swarm-cronjob --help` or moved to a permanent location.

```
$ ./swarm-cronjob --help
usage: swarm-cronjob [<flags>]

Create jobs on a time-based schedule on Swarm. More info on
https://github.com/crazy-max/swarm-cronjob

Flags:
  --help              Show context-sensitive help (also try --help-long and
                      --help-man).
  --timezone="UTC"    Timezone assigned to the scheduler.
  --log-level="info"  Set log level.
  --log-json          Enable JSON logging output.
  --version           Show application version.
```

## Server configuration

Steps below are the recommended server configuration.

### Prepare environment

Create user to run swarm-cronjob (ex. `swarmcronjob`)

```
groupadd swarmcronjob
useradd -s /bin/false -d /bin/null -g swarmcronjob swarmcronjob
```

### Copy binary to global location

```
cp swarm-cronjob /usr/local/bin/swarm-cronjob
```

## Running swarm-cronjob

After the above steps, two options to run swarm-cronjob:

### 1. Creating a service file (recommended)

See how to create [Linux service](linux-service.md) to start swarm-cronjob automatically.

### 2. Running from command-line/terminal

```
/usr/local/bin/swarm-cronjob
```

> :bulb: When launched manually, swarm-cronjob can be killed using `Ctrl+C`

## Updating to a new version

You can update to a new version of swarm-cronjob by stopping it, replacing the binary at `/usr/local/bin/swarm-cronjob` and restarting the instance.

If you have carried out the installation steps as described above, the binary should have the generic name `swarm-cronjob`. Do not change this, i.e. to include the version number.

## Next

You are now ready to [deploy cronjob based services with swarm](../get-started.md).
