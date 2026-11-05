# 0243: the website's statements are spaced by phase

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** A blank line separates every declaration from the guard that tests it in the website, a `return` that follows a statement gets its blank line in tests too, no callback names a parameter `_`, and every awaited result in a test is named before it is consumed.

**Blocked by:** 0239, 0240.

**Status:** ready-for-agent

- [x] `platform/cli/drift quality check general` prints no finding under `apps/website`.
- [x] `pnpm turbo run check typecheck test build --filter=website` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check general`
