# 0205: the tech lead overrides a ticket below the floor

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** Gate `SEAT-03` marks a delivered ticket below the floor as `needs-override` before the owner's stage. `ex agent safe-seat override <runtime-id> --reason … --by …` then answers `SEAT-03`:
- it shows the measured share to whoever decides;
- its private record keeps the share;
- its mirror carries the state and the reason, never the share.

The ticket's latest lap decides its state:
- an override of that lap makes it `overridden`;
- a lap that meets the floor makes it `met`;
- a lap that was not observed makes it `unknown`, whatever an earlier lap showed;
- otherwise it stays `needs-override`.

The run and the README say plainly that `ex` records a merge without an override but cannot prevent one on the hosting service.

**Blocked by:** 0204.

**Status:** ready-for-agent

- [ ] Below the floor at delivery, the run prints `SEAT-03` and `needs-override`, and the committed run record holds that state.
- [ ] `ex agent safe-seat override` prints the measured share and the floor, records the reason, the `--by` and the time, and turns the state to `overridden`. An override without a reason is refused.
- [ ] The override's private record holds the share, and its mirror holds the state and the reason and no number.
- [ ] A later lap at or above the floor makes the ticket `met`. A later lap without observed delivery makes it `unknown`, even after an earlier `met`.
- [ ] The `ex` README documents `SEAT-03`, the override at delivery and the limit that `ex` cannot block a merge on the hosting service.

## Dev part

**Category:** core-domain. The ticket's Safe Seat state across its history.

**Paths:**
- `packages/ex/src/domain/run/resolve-safe-seat-state.go`
- `packages/ex/src/domain/run/resolve-safe-seat-state_test.go`

**Boundary:** The run record writer calls `run.ResolveSafeSeatState(laps []run.SafeSeatLap, overrides []run.SafeSeatOverride) run.SafeSeatState`. Each lap keeps the state that ticket 0204 measured for it; this function decides the ticket's state from them. The agent declares the three types in their own files and switches the ticket's state to this function.

**Goal:** Resolve the ticket's state from its laps and overrides:
- the latest lap decides;
- an override answers only the lap it names;
- a latest lap that is `unknown` makes the ticket `unknown`, even after an earlier `met`.

- [ ] Table tests cover:
  - a single lap in each state;
  - an override of a below-floor lap;
  - a below-floor lap followed by a met lap;
  - a met lap followed by an unknown lap;
  - an override naming an earlier lap;
  - only unknown laps.

Until the owner writes these files, the build fails on `ResolveSafeSeatState` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/ex`
- `pnpm turbo run check`
