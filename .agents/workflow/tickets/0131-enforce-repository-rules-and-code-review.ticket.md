# 0131: Enforce repository rules and code-review in pre-commit and CI

**Spec:** .agents/workflow/specs/0023-enforce-repository-rules-and-code-review.spec.md

**What to build:** Add mandatory two-axis `code-review` before every commit to `.claude/rules/00-general.rule.md`, update `code-review` and `implement` skills to explicitly evaluate all 10 rule files, and implement `drift check rules` in `platform/cli/src/commands/check-rules.sh` to mechanically enforce general, filesystem, frontend CSS whitelist/ratchet, Next.js, and Jotai rules. Wire `check rules` and `check biome` into `platform/cli/src/commands/git-pre-commit.sh` and `.github/workflows/validate.yml`.

**Status:** ready-for-agent

- [x] `.claude/rules/00-general.rule.md` adds Rule 4 mandating two-axis `code-review` before every commit.
- [x] `.claude/skills/code-review/SKILL.md` mandates `.claude/rules/*.rule.md` as the authoritative source of standards.
- [x] `.claude/skills/implement/SKILL.md` requires running `code-review` against all rules before delivery.
- [x] `platform/cli/src/commands/check-rules.sh` implements mechanical checks for rules 00, 01, 02, 03, 04, 05, 06, and 07.
- [x] `platform/cli/src/commands/check-biome.sh` verifies Biome configuration and import constraints.
- [x] `platform/cli/src/drift.sh` and `platform/cli/src/commands/help.sh` register `check biome` and `check rules`.
- [x] `platform/cli/src/commands/git-pre-commit.sh` executes `drift check rules` on every local commit.
- [x] `.github/workflows/validate.yml` adds `drift check rules` and `drift check biome` to the `check` job.

### Verification

- `platform/cli/drift check rules`
- `platform/cli/drift check biome`
- `platform/cli/drift check workflows`
- `platform/cli/drift git pre-commit`
- `pnpm turbo run check`
- `pnpm turbo run typecheck`
