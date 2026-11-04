# Book a Demo on microservices

**Status:** ready-for-agent

## Tickets

- [x] 0186: Nest rule and nested-workspace mechanics
- [x] 0187: `@langdrift/events`: envelope and demo-request events
- [x] 0188: `@langdrift/events/transport`: Redis Streams strategy and client
- [x] 0189: demo-requests accepts a request
- [x] 0190: demo-requests reads the log
- [ ] 0191: demo-requests archives every event to MinIO
- [ ] 0192: mailer sends and records
- [ ] 0193: website edge: interceptors and forwarding

## Problem Statement

A visitor who books a demo on the website gets a confirmation email and the team gets a lead notification, both sent by the website's own route handler through Resend. Spec 0029 chose that shape on purpose: "directly within the Next.js website application, without spinning up an extra microservice". It leaves three problems.

The website holds email credentials and sends mail itself, which its README says it never does. The request is not recorded anywhere: when the confirmation goes out and the notification fails, the visitor sees an error although their request was received, and nothing keeps what was sent, to whom, or when. And the email copy for five locales lives in TypeScript maps, which the localization rules forbid.

The platform also has no service of its own yet. LangDrift FinOps records executions as evidence with provenance; its first backend should be built the way the platform asks its consumers to build, as bounded services whose history is an event log, and it should run as a microruntime planned and executed by `ex`, the platform's own engine. Today nothing in the repository is run by `ex`.

## Solution

Email dispatch leaves the website and becomes two Nest.js microservices under a new `apps/microservices/` folder, with an event log between them.

The website route keeps every check it performs today, now as one interceptor per concern in an `interceptors/` folder that reuses the form's Zod schema, and forwards the accepted submission to the **demo-requests** service with a bearer token. While that service is unreachable or not configured, the route answers that requests are paused and the form shows its return-soon notice, which is what production shows until the services are hosted.

**demo-requests** accepts a submission, appends a `demo-request.received` event to a Redis Stream, exposes the stream as Server-Sent Events, lets a client read one request's state by folding its events, and archives every event of the stream into a MinIO bucket as an object, so trimming the stream never loses history. **mailer** consumes the same stream in its own consumer group, sends the confirmation to the visitor and the lead notification to the team through Resend, and records each outcome as an event (`demo-request.email-sent` or `demo-request.email-failed`), which the archive stores like any other: that is how every sent email is kept. The copy of both emails moves to JSON catalogs in the five locales.

The two services share one workspace package, `@langdrift/events`: the event envelope, the demo-request event schemas and a Nest transport for Redis Streams, so Nest's microservice features (`@EventPattern` handlers, hybrid applications) run on a durable log instead of pub/sub.

Every component of the future microruntime has a Dockerfile: the two services, and `redis` and `minio` as infrastructure components. Running them under `ex` is spec 0032.

## User Stories

1. As a visitor, I want to book a demo and receive a confirmation in my language, so that I know the team will be in touch.
2. As a visitor, I want the form to tell me when demo requests are paused and that nothing was stored, so that I do not wait for an email that will not come.
3. As a visitor, I want a validation mistake answered by the website before anything is sent, so that I get the field errors immediately.
4. As a team member, I want a lead notification for every accepted request at the address the team configured, so that nobody has to poll a mailbox or a dashboard.
5. As a team member, I want to read one request's state (received, confirmed, failed) and the emails sent for it, so that a visitor who writes in can be answered with facts.
6. As a team member, I want the history of every request kept as objects in a bucket, so that a trimmed stream or a restarted Redis never loses a lead.
7. As an operator, I want the website to hold no email credential and send no mail, so that the website's README stays true and the credential lives in one service.
8. As an operator, I want the website's route to refuse a foreign origin, a wrong content type, an oversized or malformed body and an invalid submission exactly as it does today, so that moving the dispatch does not weaken the edge.
9. As an operator, I want the website to forward to the service with a bearer token read from its environment, so that only the website can create requests.
10. As an operator, I want the service to refuse a request without the token or with a wrong one, comparing in constant time, so that the token cannot be guessed by timing.
11. As an operator, I want an unconfigured or unreachable service to surface as the paused state, not as a server error, so that production degrades to a truthful notice.
12. As an operator, I want each service to expose a health endpoint, so that a runtime can probe its readiness.
13. As an operator, I want the mailer's handler to be idempotent per request through Resend's idempotency key, so that a redelivered event never sends a second email.
14. As an operator, I want an event the mailer cannot handle after a bounded number of deliveries parked on a dead-letter stream, so that a poison entry never blocks the group.
15. As an operator, I want a consumer that dies mid-handling to leave its entry pending and another consumer of the group to claim it, so that a crash loses no event.
16. As an operator, I want both services to shut down gracefully on SIGTERM, finishing the batch in flight before closing the connections, so that a rollout acknowledges nothing it did not handle.
17. As an operator, I want every environment variable of each service documented by key in its README and read where it is used, with no default standing in for a missing value, so that a misconfiguration fails loudly.
18. As a developer, I want to run both services locally against Redis and MinIO started from their own Dockerfiles, so that the local environment matches the runtime's components.
19. As a developer, I want to watch the stream over Server-Sent Events with `curl`, from the beginning or from a given entry, so that I can see an event-sourced request move through the system.
20. As a developer, I want the event envelope and every event payload validated with Zod on both the producing and the consuming side, so that a malformed entry is refused with a reason instead of crashing a consumer.
21. As a developer, I want the Redis Streams transport to be a Nest transport (a custom server strategy and a client), so that handlers are ordinary `@EventPattern` methods and the services stay Nest microservices.
22. As a developer, I want unit tests that construct a class with fakes and never boot Nest, so that the suite runs on Vitest without decorator metadata and finishes in seconds.
23. As a developer, I want the email copy in JSON catalogs with the same keys, placeholders and valid ICU in every locale, enforced by a parity test, so that a missing translation is a failing test.
24. As a maintainer, I want a rule file for Nest apps that names their tree, suffixes, injection, validation, catalogs and tests, so that the second service looks like the first.
25. As a maintainer, I want the repository's mechanics (workspace globs, dead-code scan, localization check, launch configurations) to know the nested `apps/microservices/*` layout, so that the new apps are checked like the others.
26. As a maintainer, I want spec 0029 marked as superseded, so that nobody reads it as the current design.
27. As the owner, I want every component to have a Dockerfile with a `HEALTHCHECK`, so that spec 0032 can run them as a microruntime with readiness derived from the component itself.

## Implementation Decisions

**Services and package.** Two Nest.js applications on the Express adapter under `apps/microservices/`, `demo-requests` on port 3402 and `mailer` on port 3403, and one capability package `@langdrift/events` under `packages/`, built to ES modules, which Nest 12 requires, and consumed through its exports; the package exists because both services consume it from the first day. `redis` and `minio` are folders under `apps/microservices/` holding only a Dockerfile and a README line; pnpm ignores them.

**Event log.** One Redis Stream, `demo-requests:events`, trimmed approximately at 100,000 entries, with a dead-letter stream beside it. Each entry carries three fields: `type`, `aggregateId` (both flat, so a range read filters without parsing) and `envelope` (JSON). The envelope is `{ id, type, version, aggregate: { type: 'demo-request', id }, occurredAt, correlationId, causationId?, payload }`; `correlationId` is the aggregate id and a mailer event's `causationId` is the id of the received event. Three event types in version 1:

| Type | Producer | Payload |
| --- | --- | --- |
| `demo-request.received` | demo-requests | email, company, position, locale, submittedAt |
| `demo-request.email-sent` | mailer | kind (`confirmation` or `lead-notification`), provider, messageId, to, sentAt |
| `demo-request.email-failed` | mailer | kind, provider, reason, failedAt |

A new version of a type is a new schema and a new member of the discriminated union. Two consumer groups: `archive`, run inside demo-requests, and `mailer`; the consumer name is the host name, so a restarted pod claims its own pending entries.

**Transport.** `@langdrift/events` ships a custom Nest transport strategy for Redis Streams on `ioredis` and a matching client. The server creates the group with `MKSTREAM`, reads with `XREADGROUP … BLOCK` on a dedicated blocking connection, dispatches each entry to the handler registered for its `type` or to the wildcard handler, acknowledges only after the handler's result completes, leaves the entry pending on a throw, claims idle pending entries with `XAUTOCLAIM` at startup and on an interval, parks an entry past the maximum deliveries on the dead-letter stream, and closes by disconnecting the blocking connection, awaiting the batch in flight and quitting. An entry with no handler in the group is acknowledged and skipped. The client appends with `XADD … MAXLEN ~` and does not support request-response. Nest's built-in Redis transporter is pub/sub, not durable, and is not used.

**demo-requests.** A hybrid application: HTTP routes plus the `archive` consumer. `POST /v1/demo-requests` validates with the shared submission schema through a Zod pipe, appends the received event and answers 202 with the id and a `Location`; 400 with field errors; 503 when the log is unavailable. `GET /v1/demo-requests/events` is a Nest `@Sse()` route that starts at `Last-Event-ID` or at the live end, replays from a given entry on request, and sends a ping on every empty block. `GET /v1/demo-requests/:id` folds the aggregate's events into a view (`received`, `confirmed`, `failed`, the submission and the emails) and answers 404 when none exist; the fold reads the stream, and the MinIO archive is the named upgrade path for a projection. Every route under `/v1` is behind a bearer-token guard that reads the token from the environment at the check and compares with `timingSafeEqual`; a missing token throws, never falls open. The archive consumer handles every type, ensures the bucket on module init and writes each envelope as `demo-request/<aggregateId>/<streamId>.json` with the event type as object metadata.

**mailer.** A hybrid application: a health route plus the `mailer` consumer. One handler for `demo-request.received` sends the confirmation to the visitor and the lead notification to the configured team address, each through Resend with an idempotency key derived from the event id and the email kind, and appends a sent or failed event for each. The two Resend adapters move from the website; the sender stays `no-reply@langdrift.md`; the copy comes in as an argument formatted from JSON catalogs with ICU through `intl-messageformat`.

**Website edge.** The route keeps its order of checks as one interceptor per concern in an `interceptors/` folder of the domain: method, origin, content type, body size, body (bounded read and strict JSON), schema. Each interceptor answers the refusal itself or yields its value; the route chains them and only then forwards. An adapter under the website's infrastructure folder posts the accepted submission to the service with the bearer token, a bounded timeout, no redirects and no cache, and returns the status. Upstream 202 becomes 200 `OK`; a missing URL or token, a network failure, a timeout or any other status becomes 503 `UNAVAILABLE` with `Retry-After`, the state the form already renders. The `DELIVERY_FAILED` reply code goes. The website's Resend adapters, their tests and the `resend` dependency are deleted; its README documents the two new keys and stops saying the form answers that requests are paused.

**Nest conventions and the rules.** Injection by token on every constructor parameter with decorator metadata off, so Vitest needs no plugin and a missing `@Inject` fails at startup. `nest build` bundles with Nest's own webpack configuration so the app's path alias resolves at runtime; one alias per app. Providers that need the environment are factories that read `process.env` at creation and throw when a value is missing; no configuration module. Catalogs live in each service's `i18n/messages` folder as JSON; loaders hold no copy. These decisions become `.claude/rules/11-nestjs.rule.md`, drafted here for the owner to promote, and the `.interceptor` suffix joins the TypeScript suffix table.

**Dockerfiles.** Every component has one with a `HEALTHCHECK`. The services build from the repository root context in four stages with pnpm as a checksum-pinned static binary and Turborepo's `prune --docker`, run as the `node` user and expose their port. `redis` is `redis:7-alpine` with append-only persistence; `minio` is a pinned `minio/minio` release whose root credentials arrive at run time. A root `.dockerignore` bounds the context. CI does not build images.

**Workspace mechanics.** The workspace globs gain `apps/microservices/*`; the dead-code scan and the localization check learn the nested layout and the Nest catalog path; the launch configurations gain both services; the context document lists the new layout and the "Event log" term.

**Environment contract.** demo-requests: `PORT`, `REDIS_URL`, `MINIO_ENDPOINT`, `MINIO_PORT`, `MINIO_USE_SSL`, `MINIO_BUCKET`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `DEMO_REQUESTS_TOKEN`. mailer: `PORT`, `REDIS_URL`, `RESEND_API_KEY`, `DEMO_LEAD_EMAIL`. website: `DEMO_REQUESTS_URL`, `DEMO_REQUESTS_TOKEN`. Non-secret defaults are tracked in each app's development environment file; secrets only in the ignored one.

## Testing Decisions

A good test exercises a seam from the outside and asserts observable behavior: the reply of a route, the entries appended to the log, the objects written, the Resend call made, never the internal call order.

Seams, highest first:

1. **The website route.** `POST` through the route handler with a stubbed transport and stubbed environment, as the existing route tests do: every refusal case stays, plus forwarding with the bearer header, 200 on 202, 503 on upstream failure, on a thrown fetch and on each missing key.
2. **The Redis Streams transport.** The server strategy and the client against a recording fake of the `ioredis` subset they use (`xgroup`, `xreadgroup`, `xack`, `xautoclaim`, `xpending`, `xadd`, `disconnect`, `quit`): creation of the group, dispatch by type and wildcard, acknowledgement after completion, pending on throw, dead-lettering past the maximum deliveries, graceful close.
3. **Controllers and services constructed directly with fakes.** The demo-requests controller with a fake producer and reader (202 with id and `Location`, 400 field errors, 503 on log failure, the SSE observable from `Last-Event-ID`, the 404 fold); the archive service with a fake MinIO client (bucket ensured once, object key and metadata); the mailer handler with a fake Resend client and producer (two sends with their idempotency keys, a sent event per success, a failed event per failure, the team address read from the environment).
4. **Pure functions.** Envelope creation, stream-entry parsing, the fold, the SSE mapping, the dead-letter rule, the message formatters, each with its own unit test.
5. **Catalog parity.** The mailer catalog test mirrors the website's: same keys, placeholders and ICU validity in every locale, the counts 0, 1, 2, 3 and 12 format everywhere.

Prior art: the website's route tests, Resend adapter tests, early-access submission tests and catalog tests, all Vitest with colocated `*.unit.ts` files whose descriptions name the path of the file. No test boots a Nest application and none reaches Redis, MinIO or Resend; the end-to-end path is run by hand and summarized in the pull request.

## Out of Scope

- Hosting the services in production; the website stays paused in production until then.
- Running the services as a microruntime under `ex`, and the contract changes that requires (spec 0032).
- A read-model projection in MinIO; the fold over the stream is enough for a marketing form.
- A browser consumer of the SSE endpoint (the console) and the account that would authenticate it.
- Rate limiting or a challenge on the service; the website edge is the only producer.
- Moving the position and locale enumerations out of the website into the package; both repeat them for now.
- Any change to the early-access flow.

## Further Notes

- This spec supersedes spec 0029's decision to dispatch email inside the website; 0029 keeps the confirmation page and the toast it introduced.
- The sent/failed events are the token-free record the platform asks for: what was sent, to whom, when, with the provider's message id, and never the API key.
- Resend's `delivered@resend.dev` passes the form's free-mail refinement and is the recipient for the manual end-to-end check.
