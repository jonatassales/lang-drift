# 0238: demo-requests and mailer apply the general rules

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** The two Nest microservices close their general-rule findings: a test's helpers move to `testdata/` beside it, one function per file; files are named after their export; a declaration and the guard that tests it are separated by a blank line; every `return` in a test that follows a statement gets its blank line; no `_` parameter and no awaited result consumed before it is named. The Nest-specific items stay with the Nest compliance spec, not yet written. A demo request is accepted, forwarded, emailed and recorded as today.

**Blocked by:** 0236.

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check general` prints no finding under `apps/microservices/demo-requests` or `apps/microservices/mailer`, except the Nest-specific items the spec leaves out.
- [ ] The vitest suites of both services pass unchanged in their assertions.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test --filter=demo-requests --filter=mailer`
- `platform/cli/drift quality check general`
