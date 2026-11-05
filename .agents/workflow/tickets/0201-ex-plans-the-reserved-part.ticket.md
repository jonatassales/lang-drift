# 0201: ex plans the reserved part and refuses a run without one

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** With Safe Seat on, `ex agent plan` proposes the ticket's reserved part:
- paths, a goal, acceptance criteria and a category (`core-domain`, `hard-part` or `new-stack`);
- `plannedLines` and the ticket's `totalPlannedLines`;
- a reason, a confidence and `chosenBy`.

The judge suggests a part only if it qualifies: its planned share reaches the team's percentage, and it is not boilerplate, configuration, generated code, a rename or a tool-run migration. When nothing qualifies, the plan records why. Without a judge, the part stays null with the reason. The lead can replace the suggestion with paths, a goal, acceptance criteria and a category of their own; the plan then shows that pick's planned-share check as met, below or unknown, without blocking it. The reserved part lives in the plan record and the `langdrift-plan` fence.

`ex agent run` authorizes the plan with its reserved part. Gate `SEAT-01` refuses the run, with nothing paid, when Safe Seat is on and the plan has neither a reserved part nor an override. `ex agent safe-seat override <runtime-id> --reason … --by …` records a `safe-seat-override` input record that answers `SEAT-01`, and the run can then start.

**Blocked by:** 0200.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A scripted judge's qualifying suggestion lands in the plan fence with every field. A suggestion below the team's percentage is not offered, and the plan says why.
- [ ] The lead's own pick replaces the suggestion, is recorded with `chosenBy` set to the lead, and shows its planned-share check.
- [ ] `ex agent plan --no-judge` leaves the reserved part null with the reason.
- [ ] With Safe Seat on, `ex agent run` without a reserved part or an override exits as a paused gate naming `SEAT-01`, and nothing is paid or reserved.
- [ ] After `ex agent safe-seat override` with a reason, the same run starts. An override without a reason is refused.
- [ ] With Safe Seat off or not recorded, planning and running behave as today.
- [ ] The `ex` README documents the reserved part, `SEAT-01` and the override. The contracts README documents the fence fields and the override record.

## Dev part

**Category:** core-domain. What makes a part worth reserving.

**Paths:**
- `lib/cli/src/domain/run/qualify-reserved-part.go`
- `lib/cli/src/domain/run/qualify-reserved-part_test.go`

**Boundary:** The planner calls `run.QualifyReservedPart(part run.ReservedPart, percent int) run.Qualification`, where `Qualification` carries the verdict (`met`, `below` or `unknown`) and the part's `excludedKind` when there is one. `ReservedPart.excludedKind` is the judge's assessment and is null for the lead's pick. The agent declares `ReservedPart` and `Qualification` in their own files and writes the callers.

**Goal:** Decide whether a proposed part qualifies:
- `unknown` when either planned size is missing or the total is zero;
- `below` when the planned share falls short of the percentage;
- a part with an `excludedKind` is reported with it and never qualifies.

- [ ] Table tests cover met, below, unknown, a zero total and each excluded kind.

Until the owner writes these files, the build fails on `QualifyReservedPart` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/cli`
- `pnpm turbo run check`
