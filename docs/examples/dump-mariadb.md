# Dump MariaDB database

```yaml
version: "3.2"

services:
  db:
    image: mariadb:10.4
    environment:
      - "MYSQL_ROOT_PASSWORD=sup"
      - "MYSQL_DATABASE=foo"
      - "MYSQL_USER=foo"
      - "MYSQL_PASSWORD=bar"
    volumes:
      - "db:/var/lib/mysql"

  dump:
    image: mariadb:10.4
    command: bash -c "mkdir -p /dumps && /usr/bin/mysqldump -v -h db -u root --password=sup foo | gzip -9 > /dumps/backup-$$(date +%Y%m%d-%H%M%S).sql.gz && ls -al /dumps/"
    depends_on:
      - db
    volumes:
      - "dumps:/dumps"
    deploy:
      labels:
        - "swarm.cronjob.enable=true"
        - "swarm.cronjob.schedule=* * * * *"
        - "swarm.cronjob.skip-running=true"
      replicas: 0
      restart_policy:
        condition: none

volumes:
  db:
  dumps:
```
