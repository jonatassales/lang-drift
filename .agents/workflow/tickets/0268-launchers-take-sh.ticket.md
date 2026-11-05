# 0268: the drift and ex launchers take .sh

**Spec:** .agents/workflow/specs/0037-tooling-rules-compliance.spec.md

**What to build:** `platform/cli/drift` becomes `platform/cli/drift.sh`, and `drift setup` links it under the name `drift`. `lib/cli/bin/ex` becomes `lib/cli/bin/ex.sh`; the package's `bin` map keeps the command `ex` pointing at it, `files` lists it, and the launcher test reads it by the new name. Every caller follows: the root `postinstall`, the husky hooks, the workflows, the widgets, `drift doctor`'s required paths and the handbook.

**Blocked by:** 0265.

**Status:** ready-for-agent

- [ ] No extensionless automation file remains in `platform/cli/` or a package `bin/`.
- [ ] `drift` and `ex` run under their old command names; `pnpm pack --dry-run` in `lib/cli` lists `bin/ex.sh`.
- [ ] Whether the launcher change ships as a patch release is recorded for the owner.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift.sh doctor`
- `cd lib/cli && pnpm pack --dry-run`
