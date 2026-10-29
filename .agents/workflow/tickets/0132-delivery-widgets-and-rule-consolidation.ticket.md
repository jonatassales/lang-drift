# 0132: Delivery widgets and rule consolidation

**Spec:** .agents/workflow/specs/0023-enforce-repository-rules-and-code-review.spec.md

**What to build:** Consolidate repository rules into V3 structure, relocate dead-code allowlist to platform CLI exceptions, configure lint-staged in package.json, implement the 5 delivery widgets (Workflow, FinOps, Microruntimes, Code Quality, Agentic) in .github/widgets, and update CI with Harness - Guardrails, Harness - Score, and Delivery - Dashboard.

**Status:** ready-for-agent

- [x] Consolidate repository rules into `.claude/rules/*.rule.md` V3.
- [x] Relocate `dead-code.allow` to `platform/cli/src/exceptions/dead-code.allow`.
- [x] Remove `platform/quality` and configure `lint-staged` in `package.json`.
- [x] Implement `.github/widgets/workflow.sh`.
- [x] Implement `.github/widgets/finops.sh`.
- [x] Implement `.github/widgets/microruntimes.sh`.
- [x] Implement `.github/widgets/code-quality.sh`.
- [x] Implement `.github/widgets/agentic.sh`.
- [x] Implement `.github/widgets/dashboard.sh`.
- [x] Wire `drift widget ...` subcommands in `platform/cli/src/drift.sh`.
- [x] Update `.github/workflows/pull-request.yml` with `Harness - Guardrails`, `Harness - Score`, and `Delivery - Dashboard`.

### Verification

- `platform/cli/drift check rules`
- `platform/cli/drift check workflows`
- `platform/cli/drift widget dashboard`
- `platform/cli/drift git pre-commit`
- `pnpm turbo run check`
