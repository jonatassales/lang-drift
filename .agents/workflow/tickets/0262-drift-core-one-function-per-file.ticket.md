# 0262: drift's core shell declares one function per file

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** The functions of `platform/cli/src/core/common.sh`, `output.sh` and `commit-rules.sh` move to `platform/cli/src/core/<function>.sh`, one each, named after the function without its `drift_` prefix; `common.sh` and `output.sh` disappear; a `core/load.sh` with no function of its own sources every sibling, and every command sources it. `drift.sh`'s `drift_usage_error` moves to `core/usage-error.sh`. `drift quality dead-code` keeps tracking the functions across the new layout.

**Blocked by:** 0246.

**Status:** ready-for-agent

- [ ] No file under `platform/cli/src/core/` declares more than one function.
- [ ] `drift doctor`, `drift check biome`, `drift check localization`, `drift quality check`, `drift quality coverage` and `drift quality dead-code` behave as before.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift doctor`
- `platform/cli/drift quality dead-code`
- `platform/cli/drift quality check general`
