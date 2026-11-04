# demo-requests

The service that receives a Book a Demo request from the website, records it as the first event of an event-sourced request on the Redis Stream `demo-requests:events`, and (spec 0031) exposes the stream over Server-Sent Events, folds a request's events into its state and archives every event into MinIO. It is a Nest.js application on Express, consuming `@langdrift/events` for the event schemas and the Redis Streams transport. The [Nest rules](../../../.claude/rules/11-nestjs.rule.md) say how it is laid out.

## Endpoints

| Route | Auth | Answer |
| --- | --- | --- |
| `POST /v1/demo-requests` | bearer | 202 `{ "id": … }` with `Location: /v1/demo-requests/<id>` once `demo-request.received` is on the stream; 400 `{ "code": "VALIDATION_ERROR", "fieldErrors": { … } }`; 401 without the token; 503 when the event log cannot take the event |
| `GET /health` | none | 200 `{ "status": "ok" }` |

The bearer token is `Authorization: Bearer <DEMO_REQUESTS_TOKEN>`, compared in constant time. Only the website holds it.

## Environment

Keys, never values. `.env.development` holds the non-secret defaults and is tracked; `drift env setup` creates the ignored `.env` beside it for the rest. Every value is read where it is used and a missing one throws.

| Key | Meaning |
| --- | --- |
| `PORT` | The port the HTTP server listens on (3402 in development). |
| `REDIS_URL` | The Redis connection string of the event log. |
| `MINIO_ENDPOINT`, `MINIO_PORT`, `MINIO_USE_SSL` | Where MinIO listens. |
| `MINIO_BUCKET` | The bucket the archive writes into. |
| `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY` | The archive's credentials (`.env` only). |
| `DEMO_REQUESTS_TOKEN` | The bearer token the website presents (`.env` only, at least 32 characters). |

## Development

From the repository root, with Redis listening on the `REDIS_URL` of `.env.development` (the `redis` component's Dockerfile builds one: `docker build -t langdrift-redis:local apps/microservices/redis && docker run -d -p 6379:6379 langdrift-redis:local`):

```sh
pnpm turbo run dev --filter=demo-requests
```

Tests construct the controller, the guard and the pipe with fakes and never boot Nest or reach Redis:

```sh
pnpm turbo run test typecheck --filter=demo-requests
```

## Image

The Dockerfile builds from the repository root (`docker build -f apps/microservices/demo-requests/Dockerfile .`): pnpm arrives as a checksum-pinned static binary, Turborepo prunes the workspace to this app, the app is built and deployed with its production dependencies, and the runtime image runs as `node`, exposes 3402 and declares a `HEALTHCHECK` on `/health`, which `ex` turns into the microruntime's readiness probe (spec 0032). CI builds no image.
