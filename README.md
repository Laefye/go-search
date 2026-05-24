## Docker compose

Services start only after Redis and RabbitMQ are healthy via healthchecks.

Build and start:

```sh
docker compose up -d --build
```

Stop:

```sh
docker compose down
```
