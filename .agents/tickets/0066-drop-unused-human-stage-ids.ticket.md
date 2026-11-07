# 0066: Drop the unused `humanStageIds` export

**Spec:** None

**What to build:** ticket 0058 left `export const humanStageIds` in `microruntimes-workflow.ts` with no consumer, so `drift check dead-code` fails on `main` and on every branch cut from it. The Workflow card derives the stages from `workflowStageIds` and the `HumanStageId` type, so the list is removed instead of being added to `platform/quality/dead-code.allow`.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `humanStageIds` no longer exists; `HumanStageId`, `workflowStageIds` and their consumers do not change.
- [ ] `drift check dead-code` passes.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check dead-code`.
