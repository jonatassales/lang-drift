# 0245: the general rules name the Go test idioms they allow

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `00-general.rule.md` gains the amendment the owner approved on 2026-11-04: under Readable names, a Go test keeps the `t *testing.T` parameter the toolchain's convention names; under Multi-field literals, each case of a Go test table may stay on one line. Nothing else changes, and spec 0027's exception for tests is recorded as withdrawn.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Readable names item 1 and Multi-field literals item 1 carry the two sentences of the spec's amendment, word for word.
- [x] No other rule text changes.

## Dev part

**Override:** No part qualifies: edits to rules, documents and references, with no problem of their own to solve (policy, Dev part). 2026-11-04.

## Verification

- `git diff .claude/rules/00-general.rule.md`
