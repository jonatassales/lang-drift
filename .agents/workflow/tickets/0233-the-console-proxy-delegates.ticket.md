# 0233: the console proxy delegates to its implementation

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** The console's `src/proxy.ts` keeps the matcher and delegates to one function in `src/proxy/` that composes the existing computes: the document policy with its nonce, the choice between the missing-screen rewrite and the locale middleware, the hardening of the locale preference cookie as one map, and the response headers. The locale middleware instance is a constant of that function's file. Behaviour, headers and cookie attributes are unchanged.

**Blocked by:** 0230.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] `apps/console/src/proxy.ts` holds the matcher and a delegation, and sets no header or cookie itself.
- [ ] The composed function has a test for the cookie hardening and the headers.
- [ ] The content-security-policy and matcher tests keep passing.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
