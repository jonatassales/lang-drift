# 0220: open tickets and documents follow the new layout

**Spec:** .agents/workflow/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The backlog and the living documents point at the tree as it now is. Open tickets that cite moved paths or the old package name are updated: the console microruntime ticket's link to the stubs module, the Book a Demo microruntime ticket (which pins `@langdrift/cli` 0.4.0 at the root and names the components `messaging` and `storage` by their new paths), the Safe Seat CLI tickets' filters and paths, and the four open stub tickets' verification steps. Closed tickets, specs and the CLI's consumer fixtures stay as written. Whatever prose in `AGENTS.md` and the context document still describes the old layout is corrected.

**Blocked by:** 0212, 0213, 0215, 0217, 0219.

**Status:** ready-for-agent

- [ ] Tickets 0070, 0077, 0085, 0091, 0116, 0198 and 0200 to 0205 cite only paths and package names that exist.
- [ ] No closed ticket, spec or CLI fixture changed.
- [ ] `git grep -n -E '@langdrift/(ex|web|tokens|events|contracts|extub)|apps/extub|microservices/(redis|minio)|EXTUB_|packages/(ex|contracts|events)'` outside `.agents/workflow` and the CLI fixtures returns nothing, and inside `.agents/workflow` matches only closed tickets and specs.
- [ ] `AGENTS.md` and the context document describe the four workspaces, the naming rule and the console-stubs contract with no stale path.
- [ ] The agent memory note on component paths names `infrastructure/messaging` and `infrastructure/storage`.

## Dev part

**Override:** No part qualifies: edits to tickets and documents, with no problem of their own to solve (policy, Dev part: a problem with its own goal and acceptance criteria). 2026-11-04.

## Verification

- The `git grep` above.
- `platform/cli/drift quality check 01-filesystem`
