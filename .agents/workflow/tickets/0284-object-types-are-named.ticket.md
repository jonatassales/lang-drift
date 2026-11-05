# 0284: every object type in a signature is a named interface

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Each object type in a signature, generic argument or assertion becomes an `interface` directly above its use, or, for the console read envelopes, in the data layer's types by concern, imported with `import type`. A one-field object counts. Unions of object shapes stay `type` aliases over named interfaces.

**Blocked by:** 0281, 0282, 0283.

**Status:** ready-for-agent

- [ ] The inline-object-type detector reports nothing.
- [ ] `pnpm turbo run check typecheck test` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test`
- `platform/cli/drift quality check 04-typescript`
