# 0269: next typegen runs as its own Turborepo task

**Spec:** .agents/workflow/specs/0037-tooling-rules-compliance.spec.md

**What to build:** The console and the website gain a `typegen` script (`next typegen`) and `typecheck` becomes `tsc --noEmit`. The root Turborepo configuration gains a `typegen` task that depends on `transit` and caches `.next/types/**`, and `typecheck` depends on `typegen`. Tickets and documents cite `pnpm turbo run typecheck --filter=<app>`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] No package script chains two tools.
- [x] `pnpm turbo run typecheck` passes from a clean clone and replays `typegen` from cache.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run typecheck`
