# 0280: every translator is named translate

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Every `useTranslations` and `getTranslations` binding is `translate`, including the destructured `[{ id }, translate]` form; a scope with two translators names them `translate` plus the namespace in PascalCase. Message keys do not change.

**Blocked by:** 0229, 0230, 0277.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The translator-name detector reports nothing.
- [ ] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift quality check 04-typescript`
