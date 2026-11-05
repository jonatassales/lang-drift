# 0016: Remove e2e testing

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** the ex CLI's infrastructure e2e suite, its CI workflow and every reference to them are deleted. E2E testing returns at another point in time.

**Blocked by:** none.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] .github/workflows/e2e.yml and the `e2e`-tagged test files are deleted.
- [x] The ex `typecheck` script runs `go vet ./...` without the `e2e` tag.
- [x] README, AGENTS.md, the spec and the tickets no longer mention the suite.
