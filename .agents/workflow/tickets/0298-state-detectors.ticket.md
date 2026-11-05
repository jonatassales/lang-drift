# 0298: drift measures the state rules

**Spec:** .agents/workflow/specs/0042-state-rules-compliance.spec.md

**What to build:** `drift quality check jotai` gains a row for an atom module imported outside the `@state` barrel and one for a `useState` seeded with an object literal.

**Blocked by:** 0292, 0293, 0294, 0295, 0296, 0297.

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check jotai` reports nothing.
- [ ] `platform/cli/drift quality coverage jotai` stays at 100%.

## Dev part

**Override:** No part qualifies: detector rows are patterns in the existing table, configuration rather than a problem with its own goal (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check jotai`
