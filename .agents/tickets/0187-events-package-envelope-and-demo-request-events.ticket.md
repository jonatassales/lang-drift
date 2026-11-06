# 0187: `@langdrift/events`: envelope and demo-request events

**Spec:** .agents/specs/0031-book-a-demo-microservices.spec.md

**What to build:** A developer can import the event envelope and the three demo-request event types from one workspace package and validate any entry with Zod. The package `@langdrift/events` under `packages/` exports the envelope schema and types, a function that creates an envelope (id, occurred-at, correlation and causation ids), the stream and consumer-group names, the submission schema (company email, company, position, locale) and the schemas of `demo-request.received`, `demo-request.email-sent` and `demo-request.email-failed` in version 1, plus the discriminated union over the type. It builds to ES modules, typechecks and tests through Turborepo like the other packages.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `@langdrift/events` is a workspace package with `build`, `typecheck` and `test` scripts that call `tsc` and `vitest` directly, and its exports resolve from the built output.
- [x] The envelope schema accepts `{ id, type, version, aggregate: { type, id }, occurredAt, correlationId, causationId?, payload }` and refuses a missing or malformed field with a reason.
- [x] The three event schemas accept their payloads as the spec lists them and refuse an unknown type.
- [x] Creating an envelope for a received event sets the correlation id to the aggregate id; creating a mailer event sets the causation id to the received event's id.
- [x] Each schema and the envelope factory has a colocated unit test whose description names the package and the file path.
- [x] `drift quality dead-code` reports nothing new for the package once the services consume it (verified at ticket 0192, when every export has a consumer).

### Verification

- `pnpm turbo run build typecheck test --filter=@langdrift/events`
- `pnpm turbo run check`
- `platform/cli/drift quality check 00-general 01-filesystem 04-typescript 10-testing`
