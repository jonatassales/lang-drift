# 0237: console-stubs' Go applies the general rules

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** Each free function of `apps/microservices/console-stubs` moves to a kebab-case file named after it in its package; methods of one type may span files grouped by concern. Literals take one field per line, one-letter names become words, a discarded value is named or its error handled, and every `return` that follows a statement gets its blank line. The stub API answers exactly as before.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `platform/cli/drift quality check general` prints no finding under `apps/microservices/console-stubs`.
- [ ] Every file declares one free function or one type with its methods, as spec 0035 reads the rule, so the wider Go detector of 0246 finds nothing here later.
- [ ] Every new file is kebab-case and named after its declaration.
- [ ] `go vet` and `go test -race` pass in the module, and the console reads it as before.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `cd apps/microservices/console-stubs && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
