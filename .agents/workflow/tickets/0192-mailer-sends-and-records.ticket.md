# 0192: mailer sends and records

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** When a demo request is received, the visitor gets the confirmation in their language and the team gets the lead notification, and both outcomes are on the stream. The second Nest app under `apps/microservices/`, `mailer`, is a hybrid application with a health route and the `mailer` consumer group; its handler for `demo-request.received` sends the confirmation and the lead notification through Resend with an idempotency key per event and kind, and appends `demo-request.email-sent` or `demo-request.email-failed` for each. The two Resend adapters move from the website (sender `no-reply@langdrift.md`, the team address from `DEMO_LEAD_EMAIL`), and the copy of both emails lives in a JSON catalog in the five locales, formatted with ICU and guarded by a parity test.

**Blocked by:** 0186, 0188.

**Status:** ready-for-agent

- [x] `pnpm turbo run dev --filter=mailer` starts the service on port 3403; `GET /health` answers `{ "status": "ok" }`.
- [x] A received event on the stream produces one confirmation to the visitor's address in the request's locale and one notification to the configured team address, each with an idempotency key, and two sent events with the provider's message ids.
- [x] A Resend failure produces a failed event with the reason for that kind and does not prevent the other send.
- [x] A redelivered received event reuses the same idempotency keys, so no second email goes out.
- [x] The catalog holds `en`, `pt-BR`, `zh-Hant`, `ja` and `vi` with the same keys and placeholders; the parity test formats every message in every locale; the loaders hold no copy; `drift check localization` passes.
- [x] The README documents `PORT`, `REDIS_URL`, `RESEND_API_KEY`, `DEMO_LEAD_EMAIL` by key; the Dockerfile builds from the root and declares a `HEALTHCHECK` on `/health`.
- [x] Unit tests construct the handler with fake Resend and producer clients and cover both outcomes per kind; none reaches Resend.

### Verification

- `pnpm turbo run build typecheck test --filter=mailer`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
- `platform/cli/drift quality check`
- Manual: with a Resend test key, post a request and watch the two sent events on the SSE stream.
