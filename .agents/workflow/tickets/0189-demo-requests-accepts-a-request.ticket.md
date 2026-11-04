# 0189: demo-requests accepts a request

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** A caller with the bearer token posts a demo request to the `demo-requests` service and gets 202 with the request id, and the `demo-request.received` event lands on the stream. The first Nest app under `apps/microservices/` exists with its health route, the bearer-token guard (constant-time comparison, token read from the environment at the check, never falling open), the Zod validation pipe answering 400 with field errors, the log producer behind a token-injected provider, a README with the environment contract by key, a tracked development environment file, a Dockerfile with a `HEALTHCHECK`, and the `redis` and `minio` component folders with their Dockerfiles so the service can be run locally against them.

**Blocked by:** 0186, 0188.

**Status:** ready-for-agent

- [x] `pnpm turbo run dev --filter=demo-requests` starts the service on port 3402 with its environment loaded from the app's `.env`, and `GET /health` answers `{ "status": "ok" }`.
- [x] `POST /v1/demo-requests` without the bearer token answers 401; with it and a valid body answers 202 with `{ id }` and a `Location` header; with an invalid body answers 400 with field errors; with the log down answers 503.
- [x] The appended entry carries the `type`, `aggregateId` and `envelope` fields and validates against the received-event schema.
- [x] The service's `tsconfig` has decorator metadata off and every constructor parameter is injected by token; `nest build` bundles the app and `node dist/main.js` starts it.
- [x] The README documents `PORT`, `REDIS_URL`, `MINIO_*`, `DEMO_REQUESTS_TOKEN` by key; a missing value throws where it is read.
- [x] The `redis` Dockerfile (append-only, `HEALTHCHECK` on `redis-cli ping`) and the `minio` Dockerfile (pinned release, `HEALTHCHECK` on its liveness) build and run with `docker build` and `docker run`.
- [x] The service Dockerfile builds from the repository root with pnpm as a checksum-pinned static binary and Turborepo's prune, runs as `node` and declares a `HEALTHCHECK` on `/health`.
- [x] Unit tests construct the controller, the guard and the pipe with fakes and never boot Nest.

### Verification

- `pnpm turbo run build typecheck test --filter=demo-requests`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
- `platform/cli/drift quality dead-code`
- `platform/cli/drift quality check`
- Manual: start Redis from its Dockerfile, run the service, `curl` the three replies.
