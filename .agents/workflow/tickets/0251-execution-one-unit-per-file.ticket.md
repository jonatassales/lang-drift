# 0251: the execution package declares one unit per file

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `lib/cli/src/application/execution` applies the general rules, after spec 0045's run phases and stop have landed; `confirm.go` is named after what it declares and `finish.go`'s maps take one field per line. Each file declares one free function or one type with its methods, named after it in kebab-case; no `else` after an exit and at most two conditional levels; compound conditions named; blank lines before returns and between phases; multi-field literals one field per line; no one-letter name but the test's `t`; no discarded result. The package's tests follow the same rules, one `Test` per file with its fakes and helpers in files of their own. Import paths and exported identifiers do not change.

**Blocked by:** 0246, 0222, 0223, 0224, 0225, 0226, 0227.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `platform/cli/drift quality check general` prints no finding in the package.
- [ ] No exported identifier or import path changes.
- [ ] `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
