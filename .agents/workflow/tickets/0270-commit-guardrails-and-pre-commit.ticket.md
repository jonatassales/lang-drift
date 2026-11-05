# 0270: pre-commit runs the three checks and the commit header has its rules

**Spec:** .agents/workflow/specs/0037-tooling-rules-compliance.spec.md

**What to build:** `drift git pre-commit` runs lint-staged's Biome, `drift check localization` and `drift quality check`, in that order. The shared commit rules gain a header check: `type(scope): subject` or `type: subject`, optionally with `!`, all lowercase, the scope one kebab-case word, at most 72 bytes and no trailing period, its error naming the violated part. `drift git commit-msg` runs it before commitlint and `drift git lint-history` applies it to every commit in the range; commitlint raises the blank-line rules to errors. `drift check commit-rules` runs fixture messages through the guardrails.

**Blocked by:** 0263, 0265.

**Status:** ready-for-agent

- [ ] `drift check commit-rules` passes with accepted and rejected fixtures, each rejection naming its part.
- [ ] A commit with a capitalized subject, two scopes or a trailing period is refused by the hook.
- [ ] The pre-commit runs the localization check.

## Dev part

**Override:** No part qualifies: the header check is one pattern the spec writes out field by field, and the rest is hook order, commitlint configuration and a fixture runner in the shape of `drift check biome` (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift check commit-rules`
- `platform/cli/drift git lint-history origin/main`
