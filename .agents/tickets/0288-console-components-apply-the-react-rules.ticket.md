# 0288: the console's components apply the React rules

**Spec:** .agents/specs/0040-react-rules-compliance.spec.md

**What to build:** The nineteen pass-through console templates are deleted after a route-by-route navigation check; the mobile navigation's open state is derived in render (`openPath === pathname`) with no effect; a page whose `loading.tsx` renders the same fallback drops its page-level `<Suspense>`; `'use client'` leaves every component that needs no client feature; functions in component scope are declarations; props interfaces sit above their components; scopes follow the rule's order; every effect that stays says why.

**Blocked by:** 0230, 0242.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] No `template.tsx` returns its children.
- [ ] Navigation parity per route is recorded under `.audits/`.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
