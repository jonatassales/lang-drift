# 0260: the CLI's presentation declares one unit per file and documents its environment

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `lib/cli/src/presentation/cli` applies the general rules; the seven discarded `MarkFlagRequired` errors go through one function that panics with the flag's name. `lib/cli/README.md` lists, by key and purpose and never by value, every environment variable `dft` reads when the ticket runs, each read staying at its operation.

**Blocked by:** 0246, 0222, 0223, 0224, 0225, 0226, 0227.

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check general` prints no finding under `lib/cli`.
- [ ] No `_ =` discards an error in `lib/cli`.
- [ ] The README lists every `os.Getenv` key the CLI reads, checked by `git grep -o 'Getenv("[A-Z_]*")' lib/cli/src`.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
