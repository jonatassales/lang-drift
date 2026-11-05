# 0005: Adopt @langdrift/ex into the monorepo

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** the pristine Exphere copy lands as the first commit, then the package joins the toolchain: `pnpm install` succeeds without Go (no postinstall), `pnpm turbo run build/test/typecheck` covers ex (test with the race detector), and PR CI's Harness job adds `go vet` and a pinned `govulncheck`. AGENTS.md records the package.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] First commit is the unmodified copy (npm name still @exphere/ex), so later diffs are reviewable.
- [ ] npm postinstall removed; build.sh kept for standalone use.
- [ ] Package scripts: build, test (-race), typecheck; all green through turbo from the repo root.
- [ ] Harness CI job runs gofmt/vet and a version-pinned govulncheck for packages/ex.
- [ ] AGENTS.md status and layout mention the package.
