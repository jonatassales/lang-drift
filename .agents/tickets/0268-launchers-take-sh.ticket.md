# 0268: the drift and ex launchers take .sh

**Spec:** .agents/specs/0037-tooling-rules-compliance.spec.md

**What to build:** `platform/cli/drift` becomes `platform/cli/drift.sh`, and `drift setup` links it under the name `drift`. `lib/cli/bin/ex` becomes `lib/cli/bin/ex.sh`; the package's `bin` map keeps the command `ex` pointing at it, `files` lists it, and the launcher test reads it by the new name. Every caller follows: the root `postinstall`, the husky hooks, the workflows, the widgets, `drift doctor`'s required paths and the handbook.

**Blocked by:** 0265.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] No extensionless automation file remains in `platform/cli/` or a package `bin/`.
- [x] `drift` and `ex` run under their old command names; `pnpm pack --dry-run` in `lib/cli` lists `bin/ex.sh`.
- [x] Whether the launcher change ships as a patch release is recorded for the owner.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift.sh doctor`
- `cd lib/cli && pnpm pack --dry-run`

> **Note, 2026-11-05.** The drift CLI moved to `packages/drift/` after this ticket was written, so `platform/cli/drift` is now `packages/drift/drift.sh`. The hooks in `packages/drift/.husky/` keep their extensionless names, because husky finds each one by its Git hook name. The pre-commit hook used to reach back to the launcher through `../`; it now runs `quality-check.sh` directly. Consumers get the `bin/ex.sh` rename with the next `@langdrift/cli` release. Until then, the root's pinned 0.4.0 still ships `bin/ex`, and the `ex` command does not change. Whether that release is a patch (0.4.1) is the owner's call.
