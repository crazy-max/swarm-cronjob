# Installation with Docker

swarm-cronjob provides automatically updated Docker :whale: images within [Docker Hub](https://hub.docker.com/r/crazymax/swarm-cronjob). It is possible to always use the latest stable tag or to use another service that handles updating Docker images.

Following multi-platform images are available:

```
$ docker run --rm mplatform/mquery crazymax/swarm-cronjob:latest
Image: crazymax/swarm-cronjob:latest
 * Manifest List: Yes
 * Supported platforms:
   - linux/amd64
   - linux/arm/v6
   - linux/arm/v7
   - linux/arm64
   - linux/386
   - linux/ppc64le
   - linux/s390x
```

Environment variables can be used within your service:

* `TZ` : The timezone assigned to the scheduler (default `UTC`)
* `LOG_LEVEL` : Log level (default `info`)
* `LOG_JSON` : Enable JSON logging output (default `false`)

Create a service that uses the swarm-cronjob image :

```
$ docker service create --name swarm_cronjob \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  --env "LOG_LEVEL=info" \
  --env "LOG_JSON=false" \
  --constraint "node.role == manager" \
  crazymax/swarm-cronjob
```

Alternatively, you can deploy the stack [swarm_cronjob.yml](../../.res/example/swarm_cronjob.yml) :

`docker stack deploy -c swarm_cronjob.yml swarm_cronjob`
