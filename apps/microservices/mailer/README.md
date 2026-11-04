# mailer

The service that sends the Book a Demo emails: for every `demo-request.received` event on the Redis Stream `demo-requests:events` it sends the confirmation to the visitor, in the visitor's language, and the lead notification to the team, through Resend, and records each outcome on the stream as `demo-request.email-sent` or `demo-request.email-failed`, which the archive of `demo-requests` then stores. It is a Nest.js hybrid application on Express: a health route, plus the `mailer` consumer group on the Redis Streams transport of `@langdrift/events`. The [Nest rules](../../../.claude/rules/11-nestjs.rule.md) say how it is laid out.

## Behavior

- Each send carries an idempotency key made of the received event's id and the email kind, so a redelivered event sends nothing twice.
- A failed send is recorded as an `email-failed` event with the provider's reason and does not stop the other email; only a failure to record an outcome leaves the entry pending, to be claimed again. An entry claimed too many times is parked on `demo-requests:events:dead-letter`.
- The copy of both emails lives in `src/i18n/messages/demo-emails.json`, in `en`, `pt-BR`, `zh-Hant`, `ja` and `vi`, formatted as ICU; the catalog test keeps the locales in parity.
- `GET /health` answers `{ "status": "ok" }`: the readiness signal of the Dockerfile's `HEALTHCHECK` and of the microruntime probe.

## Environment

Keys, never values. `.env.development` holds the non-secret defaults and is tracked; `drift env setup` creates the ignored `.env` beside it for the rest. Every value is read where it is used and a missing one throws.

| Key | Meaning |
| --- | --- |
| `PORT` | The port the health route listens on (3403 in development). |
| `REDIS_URL` | The Redis connection string of the event log. |
| `RESEND_API_KEY` | The Resend key the sender uses (`.env` only). |
| `DEMO_LEAD_EMAIL` | The team address that receives the lead notification (`.env` only). |

## Development

The service becomes ready once Redis answers: its consumer group is created on the first command, so with Redis away the process waits and the health route does not listen yet. Once up, a Redis outage makes the recording of an outcome fail at once instead of queueing it, so the entry stays pending and is handled again when Redis returns; the idempotency key keeps the email from going out twice.

From the repository root, with Redis listening on the `REDIS_URL` of `.env.development` and `demo-requests` producing events:

```sh
pnpm turbo run dev --filter=mailer
```

Tests construct the handler with a fake sender and a fake producer and never boot Nest or reach Resend:

```sh
pnpm turbo run test typecheck --filter=mailer
```

## Image

The Dockerfile builds from the repository root (`docker build -f apps/microservices/mailer/Dockerfile .`) the same way as `demo-requests`, runs as `node`, exposes 3403 and declares a `HEALTHCHECK` on `/health`. CI builds no image.
