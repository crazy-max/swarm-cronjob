# Installation from binary

## Download

swarm-cronjob binaries are available on [releases]({{ config.repo_url }}releases/latest) page.

Choose the archive matching the destination platform:

* [`swarm-cronjob_{{ git.tag | trim('v') }}_darwin_i386.tar.gz`]({{ config.repo_url }}releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_darwin_i386.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_darwin_x86_64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_darwin_x86_64.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_linux_arm64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_arm64.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_linux_armv6.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_armv6.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_linux_armv7.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_armv7.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_linux_i386.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_i386.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_linux_x86_64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_x86_64.tar.gz)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_windows_i386.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_windows_i386.zip)
* [`swarm-cronjob_{{ git.tag | trim('v') }}_windows_x86_64.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_windows_x86_64.zip)

And extract swarm-cronjob:

```shell
$ wget -qO- {{ config.repo_url }}releases/download/v{{ git.tag | trim('v') }}/swarm-cronjob_{{ git.tag | trim('v') }}_linux_x86_64.tar.gz | tar -zxvf - swarm-cronjob
```

After getting the binary, it can be tested with [`./swarm-cronjob --help`](../usage/cli.md) command and moved to a permanent location.

## Server configuration

Steps below are the recommended server configuration.

### Prepare environment

Create user to run swarm-cronjob (ex. `swarm-cronjob`)

```shell
$ groupadd swarm-cronjob
$ useradd -s /bin/false -d /bin/null -g swarm-cronjob swarm-cronjob
```

### Copy binary to global location

```shell
$ cp swarm-cronjob /usr/local/bin/swarm-cronjob
```

## Running swarm-cronjob

After the above steps, two options to run swarm-cronjob:

### 1. Creating a service file (recommended)

See how to create [Linux service](linux-service.md) to start swarm-cronjob automatically.

### 2. Running from terminal

```shell
$ /usr/local/bin/swarm-cronjob
```

!!! note
    When launched manually, swarm-cronjob can be killed using `Ctrl+C`

## Updating to a new version

You can update to a new version of swarm-cronjob by stopping it, replacing the binary at `/usr/local/bin/swarm-cronjob` and restarting the instance.

If you have carried out the installation steps as described above, the binary should have the generic name `swarm-cronjob`. Do not change this, i.e. to include the version number.
