# 0266: the Go files take kebab-case

**Spec:** .agents/workflow/specs/0036-filesystem-rules-compliance.spec.md

**What to build:** Every Go file that remains in snake_case after spec 0035's splits is renamed to kebab-case, keeping the toolchain's `_test.go` and build-constraint suffixes (`cluster_create.go` becomes `cluster-create.go`, `local_source_test.go` becomes `local-source_test.go`), in `lib/cli`, `lib/core/contracts/go` and `apps/microservices/console-stubs`. Record files under `testdata/` keep their names.

**Blocked by:** 0237, 0247, 0248, 0249, 0250, 0251, 0252, 0253, 0254, 0255, 0256, 0257, 0258, 0259, 0260.

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check filesystem` reports no snake_case Go file.
- [ ] `go vet` and `go test -race` pass in the three modules.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `cd lib/core/contracts/go && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check filesystem`
