# 0279: the messages barrel only re-exports

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** The website's messages barrel loses its side-effect import. The catalog-shape assertion becomes a pure computation module that the messages loader calls for each catalog it assembles, so the guard still runs at startup, explicitly.

**Blocked by:** 0229.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The barrel holds named re-exports only.
- [ ] The catalogs test keeps proving parity and the guard runs for every catalog.
- [ ] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift check localization`
