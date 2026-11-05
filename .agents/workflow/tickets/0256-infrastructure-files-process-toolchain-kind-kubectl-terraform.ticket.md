# 0256: the infrastructure packages files, process, toolchain, kind, kubectl and terraform declare one unit per file

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** The `lib/cli/src/infrastructure` packages files, process, toolchain, kind, kubectl and terraform apply the general rules, tests included: one unit per file in kebab-case, no `else` after an exit, compound conditions named, spacing, one field per line, readable names and no discarded result; the lock's release returns its error and callers join it into theirs.

**Blocked by:** 0246, 0222, 0223, 0224, 0225, 0226, 0227.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `platform/cli/drift quality check general` prints no finding in these packages.
- [ ] No exported identifier or import path changes.
- [ ] `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
