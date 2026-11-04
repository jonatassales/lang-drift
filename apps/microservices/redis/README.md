# redis

The event log component of the microservices: a Redis image with append-only persistence, holding the stream `demo-requests:events` that `demo-requests` appends to and `mailer` and the archive consume. It is only a Dockerfile; pnpm ignores the folder. Build and run it locally with `docker build -t langdrift-redis:local apps/microservices/redis && docker run -d -p 6379:6379 langdrift-redis:local`, or point `REDIS_URL` at any Redis 7 or later. Under `ex` (spec 0032) it runs as the `redis` component and its `HEALTHCHECK` is the readiness probe.
