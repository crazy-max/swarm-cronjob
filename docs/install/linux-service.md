# Run as service on Debian based distro

## Using systemd

!!! warning
    Make sure to follow the instructions to [install from binary](binary.md) before.

To create a new service, paste this content in `/etc/systemd/system/swarm-cronjob.service`:

```
[Unit]
Description=swarm-cronjob
Documentation={{ config.site_url }}
After=syslog.target
After=network.target

[Service]
RestartSec=2s
Type=simple
User=swarmcronjob
Group=swarmcronjob
ExecStart=/usr/local/bin/swarm-cronjob
Restart=always
#Environment=TZ=Europe/Paris

[Install]
WantedBy=multi-user.target
```

Change the user, group, and other required startup values following your needs.

Enable and start swarm-cronjob at boot:

```shell
sudo systemctl enable swarm-cronjob
sudo systemctl start swarm-cronjob
```

To view logs:

```shell
journalctl -fu swarm-cronjob.service
```
