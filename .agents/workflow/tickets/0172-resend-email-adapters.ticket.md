# 0172: Resend email adapters for confirmation and notification

**Spec:** .agents/workflow/specs/0029-resend-book-a-demo-integration.spec.md

**What to build:** Install `resend` in `apps/website`. Create `apps/website/src/infra/resend/send-demo-confirmation.adapter.ts` to send a localized confirmation email from `no-reply@langdrift.md` to the applicant. Create `apps/website/src/infra/resend/send-demo-lead-notification.adapter.ts` to send an internal notification email to `jonatas.langdrift@gmail.com` with the submission details (email, company, position, locale, timestamp). Read `process.env.RESEND_API_KEY` at the boundary and document it in `apps/website/README.md`. Colocate unit tests mocking Resend calls.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `resend` package installed in `apps/website`.
- [x] `apps/website/src/infra/resend/send-demo-confirmation.adapter.ts` sends confirmation from `no-reply@langdrift.md` to applicant.
- [x] `apps/website/src/infra/resend/send-demo-lead-notification.adapter.ts` sends lead details to `jonatas.langdrift@gmail.com`.
- [x] Unit tests in `send-demo-confirmation.adapter.unit.ts` and `send-demo-lead-notification.adapter.unit.ts`.
- [x] Environment variable `RESEND_API_KEY` documented in `apps/website/README.md`.

### Verification

- `pnpm turbo run check typecheck test --filter=website`
- `platform/cli/drift check biome`
- `platform/cli/drift quality check 00-general 01-filesystem 04-typescript`
