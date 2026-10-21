# 0016: E2E workflow in CI

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** a separate GitHub Actions workflow runs the e2e-tagged suite on ubuntu with a pinned kind, on manual dispatch and a weekly schedule. PR validation stays untouched.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] .github/workflows/e2e.yml: workflow_dispatch + weekly cron, pinned actions/images, read-only token.
- [ ] Runs GOTOOLCHAIN=local go test -tags e2e for the package with kind + docker available.
- [ ] actionlint clean.
