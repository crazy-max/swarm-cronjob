# Prune Swarm nodes

```yaml
version: "3.2"

services:
  prune-nodes:
    image: docker
    command: ["docker", "system", "prune", "-f"]
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock"
    deploy:
      mode: global
      labels:
        - "swarm.cronjob.enable=true"
        - "swarm.cronjob.schedule=0 */5 * * * *"
        - "swarm.cronjob.skip-running=false"
      restart_policy:
        condition: none
```
