# 0247: the Microruntimes contracts declare one unit per file

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `lib/core/contracts/go` applies the general rules: one free function or one type with its methods per file, named after it in kebab-case; the nine unnamed conditions of `Validate` named; one field per line; no discarded marshal error; blank lines before returns and between phases. Exported identifiers and import paths do not change.

**Blocked by:** 0246, 0222.

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check general` prints no finding under `lib/core/contracts/go`.
- [ ] No exported identifier changes and the CLI builds against it unchanged.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd lib/core/contracts/go && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
