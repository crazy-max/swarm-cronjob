# Run as service on Debian based distro

## Using systemd

> :warning: Make sure to follow the instructions to [install from binary](binary.md) before.

Run the below command in a terminal:

```
sudo vim /etc/systemd/system/swarm-cronjob.service
```

Copy the sample [swarm-cronjob.service](../../.res/systemd/swarm-cronjob.service).

Change the user, group, and other required startup values following your needs.

Enable and start swarm-cronjob at boot:

```
sudo systemctl enable swarm-cronjob
sudo systemctl start swarm-cronjob
```

To view logs:

```
journalctl -fu swarm-cronjob.service
```
