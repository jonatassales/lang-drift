# 0285: every await statement has its blank lines, in tests too

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Every `await` statement in source and tests gets its blank line before and after, two consecutive `await` statements are separated, and the one `await` shared with optional chaining lands in a binding first.

**Blocked by:** 0284.

**Status:** ready-for-agent

- [ ] The await-spacing detector, now in the `all` scope, reports nothing.
- [ ] `platform/cli/drift quality coverage 04-typescript` shows 100%.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test`
- `platform/cli/drift quality check 04-typescript`
