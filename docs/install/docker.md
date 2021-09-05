# Installation with Docker

## About

swarm-cronjob provides automatically updated Docker :whale: images within several registries:

| Registry                                                                                         | Image                           |
|--------------------------------------------------------------------------------------------------|---------------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/swarm-cronjob/)                             | `crazymax/swarm-cronjob`                 |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/swarm-cronjob)  | `ghcr.io/crazy-max/swarm-cronjob`        |

It is possible to always use the latest stable tag or to use another service that handles updating Docker images.

!!! note
    Want to be notified of new releases? Check out :bell: [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun) project!

Following platforms for this image are available:

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
```

## Usage

```yaml
version: "3.2"

services:
  swarm-cronjob:
    image: crazymax/swarm-cronjob
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
    environment:
      - "TZ=Europe/Paris"
      - "LOG_LEVEL=info"
      - "LOG_JSON=false"
    deploy:
      placement:
        constraints:
          - node.role == manager
```

Edit this example with your preferences and deploy the stack:

```shell
docker stack deploy -c swarm_cronjob.yml swarm_cronjob
```

Or use the following command:

```shell
docker service create --name swarm_cronjob \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  --env "LOG_LEVEL=info" \
  --env "LOG_JSON=false" \
  --constraint "node.role == manager" \
  crazymax/swarm-cronjob
```

You are now ready to [deploy cronjob based services with swarm](../usage/get-started.md).
