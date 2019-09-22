# Logs

Here is a sample output:

```
$ docker service logs swarm_cronjob_app
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Starting swarm-cronjob v1.2.0
swarm_cronjob_app.1.nvsjbhdhiagl@default    | Thu, 13 Dec 2018 20:04:37 UTC INF Add cronjob with schedule 0/10 * * * * * service=date_test
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
