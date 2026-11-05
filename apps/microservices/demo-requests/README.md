# demo-requests

The service that receives a Book a Demo request from the website, records it as the first event of an event-sourced request on the Redis Stream `demo-requests:events`, exposes the stream over Server-Sent Events, folds a request's events into its state and archives every event into MinIO. It is a Nest.js hybrid application on Express, consuming `infra/streaming` for the event schemas and the Redis Streams transport: the HTTP routes and the `archive` consumer group run in one process. The [Nest rules](../../../.claude/rules/11-nestjs.rule.md) say how it is laid out.

## Endpoints

| Route | Auth | Answer |
| --- | --- | --- |
| `POST /v1/demo-requests` | bearer | 202 `{ "id": … }` with `Location: /v1/demo-requests/<id>` once `demo-request.received` is on the stream; 400 `{ "code": "VALIDATION_ERROR", "fieldErrors": { … } }`; 401 without the token; 503 when the event log cannot take the event |
| `GET /v1/demo-requests/events` | bearer | Server-Sent Events: each entry with its stream id as `id`, its type as `event` and the envelope as `data`; starts at `Last-Event-ID` when the client sends one, at the live end otherwise; `?from=0-0` replays from the start; a `ping` event on every empty 15-second block |
| `GET /v1/demo-requests/:id` | bearer | 200 `{ id, state, submission, emails }`, the fold of the request's events (`received`, `confirmed` once the confirmation was sent, `failed` when it failed); 404 when the request has no events |
| `GET /health` | none | 200 `{ "status": "ok" }` |

The bearer token is `Authorization: Bearer <DEMO_REQUESTS_TOKEN>`, compared in constant time. Only the website holds it.

## Archive

The `archive` consumer group copies every event of the stream into the bucket `MINIO_BUCKET` as `demo-request/<request id>/<stream id>.json`, with `Content-Type: application/json` and the event type as `X-Amz-Meta-Event-Type`, so a trimmed stream never loses a request's history. The bucket is created when the service starts if it is absent. An entry the archive cannot write stays pending and is claimed again; an entry claimed too many times is parked on `demo-requests:events:dead-letter` with the reason.

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

The service becomes ready once Redis answers: its consumer group is created on the first command, so with Redis away the process waits and the health route does not listen yet. Once up, a Redis outage makes `POST` answer 503 at once and queues nothing, so "nothing was stored" stays true.

From the repository root, with Redis listening on the `REDIS_URL` of `.env.development` (the `messaging` component's Dockerfile builds one: `docker build -t langdrift-messaging:local packages/infra/messaging && docker run -d -p 6379:6379 langdrift-messaging:local`):

```sh
pnpm turbo run dev --filter=demo-requests
```

Tests construct the controller, the guard and the pipe with fakes and never boot Nest or reach Redis:

```sh
pnpm turbo run test typecheck --filter=demo-requests
```

## Image

The Dockerfile builds from the repository root (`docker build -f apps/microservices/demo-requests/Dockerfile .`): pnpm arrives as a checksum-pinned static binary, Turborepo prunes the workspace to this app, the app is built and deployed with its production dependencies, and the runtime image runs as `node`, exposes 3402 and declares a `HEALTHCHECK` on `/health`, which `ex` turns into the microruntime's readiness probe (spec 0032). CI builds no image.
