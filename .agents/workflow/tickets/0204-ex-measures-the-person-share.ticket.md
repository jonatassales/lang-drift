# 0204: ex measures the person's share at delivery

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** When `ex` observes delivery, it measures the ticket. The count covers the lines that the final diff adds, from the run's base commit to the delivered head; a modified line counts once. It excludes:
- the product's lockfiles;
- files marked `linguist-generated`;
- files with the `Code generated … DO NOT EDIT` header;
- the team's patterns.

Each counted line is attributed to the last snapshot interval that introduced it:
- **agent:** the interval is inside an agent call;
- **person:** any other interval, when the dev ticket's attestation is "without an agent";
- **unknown:** any other interval otherwise, and any line the chain cannot trace.

The person's share is person lines over all counted lines, with unknown lines in the denominator. The ticket's state for the lap:
- `not-applicable` when nothing counts;
- `unknown` when delivery was not observed or the base commit cannot be read;
- `met` at or above the floor;
- `needs-override` below it.

`ex agent status` prints the counts (agent, person, unknown, excluded), the share and the floor. They live in the project's private state. The committed run record carries only the state and the policy record it used.

**Blocked by:** 0203.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] One scenario covers all four sources of a line:
  - the scripted agent writes the boundary and other files;
  - the test writes the reserved paths and attests "without an agent";
  - the agent's review rewrites one person line;
  - an excluded lockfile changes.

  `ex agent status` prints the expected counts, share and state.
- [ ] Attesting "with an agent" turns the person lines into unknown lines, which stay in the denominator.
- [ ] A diff of only excluded files is `not-applicable`. A ticket below the floor is `needs-override`. A lap without observed delivery is `unknown`.
- [ ] The committed run record holds the state and the policy record and no count or share. The counts exist only in the private state.
- [ ] The person's index, worktree and refs are unchanged by the measurement.
- [ ] The `ex` README documents the measurement and `ex agent status`. The contracts README documents the state field.

## Dev part

**Category:** hard-part. Following every line back to the interval that introduced it.

**Paths:**
- `lib/cli/src/infrastructure/snapshots/attribute.go`
- `lib/cli/src/infrastructure/snapshots/attribute_test.go`

**Boundary:** The measurement reaches attribution through the snapshot port of ticket 0203. Its adapter calls `snapshots.Attribute(ctx context.Context, request AttributionRequest) (Attribution, error)`. The types:
- `AttributionRequest` carries the project root, the `base` and `head` commits as `CommitID`, the interval chain and the exclusion predicate;
- `Interval` names its opening and closing `TreeID` and its owner: agent, or outside an agent call;
- `Attribution` carries the agent, person-eligible, unknown and excluded line counts.

The agent declares `AttributionRequest`, `CommitID`, `Interval` and `Attribution` in their own files, extends the port and writes the callers.

**Goal:** Attribute every line the final diff adds:
- chain the interval trees into private, unreferenced commits;
- find, for each added line, the interval that last introduced it;
- count lines introduced after the last snapshot as outside an agent call;
- count a line the chain cannot trace as unknown;
- never create a ref or touch the index.

- [ ] Tests against a temporary repository cover:
  - lines from an agent interval;
  - lines from a person interval;
  - a person line later rewritten by an agent interval;
  - a line added after the last snapshot;
  - an excluded file;
  - a line the chain cannot trace.

Until the owner writes these files, the build fails on `Attribute` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/cli`
- `pnpm turbo run check`
