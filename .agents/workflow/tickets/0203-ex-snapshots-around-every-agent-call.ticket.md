# 0203: ex snapshots around every agent call and guards the reserved paths

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** `ex` records the run's base commit at start. It snapshots the working tree (tracked and untracked files minus ignored ones) as a git tree written through a private index:
- before and after every agent call;
- when the dev ticket opens and closes.

The person's index, worktree, branches and refs are never touched. Each agent call's stage evidence gains the paths it changed, with lines added and removed, and the two tree ids.

Gate `SEAT-02` pauses the run right after an agent call that changed a reserved path, naming the paths. The person then chooses:
- `ex agent safe-seat restore <runtime-id>` writes the reserved paths back from the snapshot taken before that call, and removes reserved files the call created;
- `ex agent safe-seat accept <runtime-id>` keeps the change.

Both record who and when in the lap's evidence and resume the run. The seeded test project becomes a git repository, so snapshots run against real git.

**Blocked by:** 0202.

**Status:** ready-for-agent

- [ ] A run records its base commit. Each scripted agent call has before and after tree ids, and its changed paths with line counts, in the stage evidence. The dev ticket has opening and closing tree ids.
- [ ] After a run, `git status`, the index, `HEAD` and every ref of the person's repository are unchanged, apart from the agent's own edits to the working tree.
- [ ] A scripted call that writes a reserved path pauses the run at `SEAT-02`, naming the path.
- [ ] `restore` brings back the before-snapshot content, removes a reserved file the call created, and resumes. `accept` keeps the change and resumes. Both are recorded with `--by` and the time.
- [ ] In a project that is not a git repository, or whose base commit cannot be read, the run takes no snapshot, says that `SEAT-02` cannot guard the reserved paths, and records the reason, so that ticket 0204 reads the ticket as `unknown`.
- [ ] The `ex` README documents snapshots, `SEAT-02`, `restore` and `accept`. The contracts README documents the new stage evidence fields.

## Dev part

**Category:** new-stack. First contact with git's plumbing.

**Paths:**
- `packages/ex/src/infrastructure/snapshots/take.go`
- `packages/ex/src/infrastructure/snapshots/take_test.go`

**Boundary:** The execution use cases reach snapshots through a port the agent declares, whose adapter calls `snapshots.Take(ctx context.Context, root string) (TreeID, error)`. The agent declares `TreeID` in its own file, wires the port and writes the callers.

**Goal:** Write the working tree as a git tree object without touching the person's repository state:
- include tracked and untracked files, minus ignored ones;
- go through a private index, never the person's;
- leave `HEAD`, refs and the person's index unchanged;
- return the same tree id for the same content.

- [ ] Tests against a temporary repository cover: a clean tree; a modified tracked file; a new untracked file; an ignored file left out; staged changes in the person's index left as they were; the same content giving the same id.

Until the owner writes these files, the build fails on `Take` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/ex`
- `pnpm turbo run check`
