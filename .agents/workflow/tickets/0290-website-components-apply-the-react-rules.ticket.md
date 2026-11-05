# 0290: the website's components apply the React rules

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** The website header's drawer modal, body scroll lock, focus return and desktop close move into a navigation drawer hook. Functions in component scope are declarations, props interfaces sit above their components, scopes follow the rule's order, and every effect that stays says why.

**Blocked by:** 0239.

**Status:** ready-for-agent

- [ ] No effect derives state in the website or the console.
- [ ] The drawer hook has a test with a fake dialog element.
- [ ] `pnpm turbo run check typecheck test build --filter=website` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
