# 0210: the home hero session stops at the dev ticket

**Spec:** .agents/workflow/specs/0043-safe-seat-on-the-website.spec.md

**What to build:** The simulated run in the home hero, which the Microruntimes page plays too, shows Safe Seat at work:

- **Implement.** It prints that it stops at the reserved part and writes only the boundary.
- **Dev ticket.** The terminal waits while the person writes the reserved part, then shows the dev ticket closed, written without an agent, with no paid call.
- **Review.** The agent's code review runs after the person, in the review stage.
- **Delivery.** After the checks pass, the run prints "Safe Seat met", with no share.

The Workflow card shows the dev ticket as a fixed row driven by the person, beside the merge row, and a Safe Seat cell once "Safe Seat met" has printed.

The gate-pause branch, the plan session and every money figure stay as they are. The `console ›` lines stay in English; the `#` comments, the card copy and the hero summary are translated, and the summary still says that every step is simulated.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] The run's stages are planning, spec, tickets, implement, dev-ticket, review and merge.
- [ ] `dev-ticket` is not one of the stages a person can switch on or off. Implement's skills read `implement · tdd`, and review's read `code-review`.
- [ ] At default choices the run prints the dev-ticket wait and close after implement, the code review after them, and "Safe Seat met" after the checks and before `ready`. Stage progress moves the dev-ticket row from pending to active to done.
- [ ] The Workflow card shows the person-driven dev-ticket row and the Safe Seat cell, and no percentage appears anywhere.
- [ ] Under reduced motion the whole run, dev ticket included, prints at once. The screen-reader summary mentions the dev ticket.
- [ ] The gate-pause branch and the plan session are unchanged, and the plan duration test still passes.
- [ ] The `microruntimes` catalog gains the new run and card keys in five locales with `en`'s key order, and the parity test passes.
- [ ] The workflow and session-card tests cover the new order, the row and the cell.

## Dev part

**Category:** `core-domain`. How the simulated run tells the Safe Seat flow: the order of the work, where the person takes the seat and when the agent comes back.

**Paths:**
- `apps/website/src/app/lib/data/work-lines.compute.ts`
- `apps/website/src/app/lib/data/work-lines.compute.unit.ts`

**Boundary:** `runSession` calls `workLines(choices: WorkflowChoices): RunLine[]` in place of its `workStageIds.flatMap(...)`, and receives every line between `gateOpened` and `runtimeReady`. The agent writes everything around it:
- the `dev-ticket` stage id;
- the new run messages `devTicketWaiting`, `devTicketClosed` and `safeSeatMet`, and their catalog copy;
- the wider `RunLine.stage`;
- `stageLines`, which `workLines` may reuse;
- the call in `runSession`.

**Goal:**
- Return the work lines in order: each confirmed stage's lines through implement, then the dev ticket's wait and close, then the code review tagged `review`.
- Honor the person flags of the stages a person can switch.
- Give the dev ticket fixed lines that no choice turns off.
- Tag every line with its stage so stage progress follows.

- [ ] Unit tests cover:
  - the default choices;
  - every switchable stage off;
  - the dev-ticket lines always present after implement;
  - the code review after the dev ticket;
  - unique line ids.

Until these files exist, the build fails on `workLines` alone.

**Override:** on 2026-11-04 the owner asked the agent to write this dev part ("Escreva você"), so the agent wrote `workLines` and its tests.

## Verification

- `pnpm --filter website test`
- `pnpm --filter website typecheck`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
