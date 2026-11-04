# 0200: ex records the Safe Seat policy

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** A team lead records the project's Safe Seat policy with `ex harness safe-seat set`: on or off, the team's percentage and the team's exclusion patterns, with `--by`. The policy is stored as a `safe-seat` input record with its history and mirrored read-only, like the budget. Turning Safe Seat on without a percentage records the 5% floor. Validation refuses:
- a percentage below 5 or above 100;
- an absolute, empty or climbing pattern.

At a terminal, `ex init` offers Safe Seat, on by default, and records the answer as the command would. A flag gives the same answer without a terminal. Without a terminal and without the flag, `init` records nothing about Safe Seat.

`ex agent plan` shows the policy in force, or "Safe Seat not recorded". `ex agent run` reads the record at authorization, and the run record names the record it used. The CI Gates report gains a Safe Seat line marked unknown, with the reason that CI has no Safe Seat source.

**Blocked by:** 0199; 0196, the last open ticket changing `ex`'s plan and run, so that two writers never share those files.

**Status:** ready-for-agent

- [ ] `ex harness safe-seat set` records and mirrors the policy:
  - `--by`, the time and the previous record are kept;
  - a second record nests the first;
  - turning Safe Seat on without a percentage records 5;
  - turning it off is recorded like any other change.
- [ ] Percentages 4 and 101 are refused with the reason, and so are absolute, empty and climbing patterns, all before anything is recorded.
- [ ] `ex init` at a terminal offers Safe Seat on by default. The flag records it without a terminal. Non-interactive `init` without the flag behaves exactly as today.
- [ ] `ex agent plan` prints the policy in force, or "Safe Seat not recorded" when there is no record.
- [ ] The run record names the `safe-seat` record read at authorization, and a policy change after authorization does not alter a running run.
- [ ] `ex gates report` prints a Safe Seat line marked unknown, with its reason.
- [ ] The `ex` README documents the command, the floor and the init question. The contracts README documents the input record kind.

## Dev part

**Category:** core-domain. LangDrift's principle, the floor, as code.

**Paths:**
- `packages/ex/src/domain/harness/validate-safe-seat.go`
- `packages/ex/src/domain/harness/validate-safe-seat_test.go`

**Boundary:** The input service calls `harness.ValidateSafeSeat(policy harness.SafeSeat) error` before recording, and the floor is the exported constant `harness.SafeSeatFloor`. The agent declares the `SafeSeat` type and the constant in their own file and writes the callers; it writes nothing in these paths.

**Goal:** Decide whether a Safe Seat policy is valid. The checks:
- the percentage runs from the floor to 100;
- each pattern is non-empty, project-relative and never climbs out of the project, including traversal embedded in a longer pattern.

- [ ] A table test covers 4, 5, 100, 101, an empty pattern, an absolute pattern, a leading `../`, an embedded `/../` and a valid glob.
- [ ] Each refusal names its reason.

Until the owner writes these files, the build fails on `ValidateSafeSeat` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/ex`
- `pnpm turbo run check`
