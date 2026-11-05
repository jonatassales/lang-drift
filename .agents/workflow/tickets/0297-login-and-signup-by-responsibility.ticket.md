# 0297: the login and signup forms keep one responsibility per file

**Spec:** .agents/workflow/specs/0042-state-rules-compliance.spec.md

**What to build:** Login and signup each get a types file for their credentials, a schema mirroring the inputs' constraints (signup's twelve-character minimum included) and a submit handler that posts through the transport and navigates home through the callback it receives.

**Blocked by:** 0230, 0295.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] Each submit handler has a test with a fake transport and navigation.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
