# Safe Seat in the run: the person writes the dev ticket

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E14.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 16 of 19, after spec 0060 and before spec 0063. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the hero simulation, implement stops at the reserved part and the run hands the dev ticket to a person: "waiting dev-ticket · the person writes the reserved part". Once the person has written it, the run closes the stage, "written without an agent · no paid call". Only then does the agent's code review run. After the pull request's checks pass, the run states "Safe Seat met". The Workflow card shows the Dev ticket as driven by the person, with no checkbox to switch it off, and a SAFE SEAT field that reads "met". That field shows whether the floor was met, never a share.

The LangDrift CLI has none of it. There is no dev-ticket stage, no waiting or closed line, no Safe Seat line or check, and nothing on the card. The reserved part exists only as a judge answer kind in the LangDrift Core contracts. The closed spec 0033 went the other way and planned more than the simulation shows: SEAT-01 to SEAT-03 rules, overrides, a `harness safe-seat set` command, and the reserved part on the Plan card.

## Solution

The dev ticket is a stage of every run, between Implement and Review, and it always belongs to the person. After the implement line, the run prints that it waits for the person to write the reserved part. While the person writes, the run makes no call to the coding agent. When the person signals they are done, inside the same `dft agent run`, the run closes the stage, stating that it was written without an agent and with no paid call. After the client's checks pass, the run prints "Safe Seat met", and the card's SAFE SEAT field reads "met". Nothing else from the closed spec 0033 is built.

## Responsibility

This element runs the dev-ticket stage, guarantees that no paid call happens during it, decides and prints Safe Seat, and records both for the card.

It does not:

- run implement or report that it stopped at the reserved part; specs 0058 and 0059 do;
- place the stage in the workflow or show the Dev ticket row; spec 0057 places the stage and owns HUMAN IN THE LOOP, and this element supplies the row's states;
- run the checks Safe Seat is printed after; spec 0064 does;
- run code-review after the stage; spec 0059 does.

## Produces

- **The dev-ticket stage state, in the run's record (spec 0065), within spec 0057's stage states.** It holds:
  - `state`: `pending`, `active` from the `waiting dev-ticket` line, `done` from the `closed dev-ticket` line;
  - `closedBy`: `person`;
  - `paidCalls`: `0`, a known zero, since the run made no stage-driver call between the two lines.

  It feeds the Dev ticket row: `the person writes the reserved part` while active, and the done mark after.
- **The Safe Seat result, in the run's record.** It holds `safeSeat: met` and `decidedAt`, written when the line prints. It feeds the SAFE SEAT card field.

No share, percentage, line count or reserved-part path is recorded.

## Shows

Phase 2, `dft agent run`. The two dev-ticket lines follow `console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part` (spec 0059) and come before `console › code-review 2 findings · fixed` (spec 0059):

```
console › waiting    dev-ticket · the person writes the reserved part
✓ closed     dev-ticket · written without an agent · no paid call
```

They print in every run that reaches them; no choice adds or removes them.

The Safe Seat line follows `✓ checks     4 of 4 passed` (spec 0064) and comes before `runtime    stopped · 5 Deployments removed · nothing left` (spec 0062):

```
✓ Safe Seat  met
```

Phase-2 card fields:

- **SAFE SEAT:** unset (a dash) until `✓ Safe Seat  met` has printed, then `met`, in the success tone. It never shows a share.
- **The Dev ticket row of HUMAN IN THE LOOP**, supplied to spec 0057: `driven by the person` throughout, with no box; `the person writes the reserved part` while the stage is active; a done mark after `closed`.

## User Stories

1. As a developer, I want every run to stop after implement and hand me the reserved part, so that I keep writing a part of every ticket by hand.
2. As a developer, I want the session to say "the person writes the reserved part", so that I know the run is waiting for my code, not for an approval.
3. As a developer, I want no checkbox that switches the dev ticket off, so that Safe Seat cannot be skipped by a choice on the card.
4. As a developer, I want the run to make no call to the coding agent while I write, so that the reserved part is mine.
5. As a developer, I want to tell the run I am done inside the same `dft agent run`, so that I need no second command.
6. As a developer, I want the run to close the dev ticket with "written without an agent · no paid call", so that the record states how the part was written.
7. As a FinOps reader, I want the dev ticket recorded with zero paid calls, a known zero, so that the usage record shows the person's part cost no model spend.
8. As the owner, I want code-review to run only after the dev ticket closed, so that the agent reviews the part the person wrote.
9. As the owner, I want "Safe Seat met" printed once the checks have passed, so that I see the run kept the person's seat before the pull request is ready.
10. As a team member watching the card, I want the Dev ticket row to say the person is writing while it is active, so that I see who the run is waiting for.
11. As a team member, I want the SAFE SEAT field to read "met" and never a share, so that Safe Seat reads as a floor kept, not a score.
12. As an auditor, I want the run's record to keep the dev ticket's closing and the Safe Seat result, so that I can tell afterwards that a person wrote part of the ticket.
13. As the console, I want the Safe Seat result and the dev ticket's states to come from `dft` as data, so that I only show them.
14. As the owner, I want nothing about Safe Seat on the Plan card and no Safe Seat rules or overrides to configure, so that the product shows exactly what the simulation shows.
15. As a developer, I want a run whose input ends while I write to stop there, so that nothing after the dev ticket runs without my part.

## Implementation Decisions

**Remove**, each beyond what the simulation shows, from what the closed spec 0033 planned and from the code:

- The SEAT-01, SEAT-02 and SEAT-03 rules, and the floor checked as a share of the diff. Serves `✓ Safe Seat  met` and the SAFE SEAT field, which never show a share.
- Overrides of the reserved part and their records.
- `harness safe-seat set`, and any Safe Seat input command (spec 0048 keeps only `init`, `plan` and `run`).
- The reserved part on the Plan card, and any reserved-part figure in the plan's hand-off. Serves the Plan card's fields, which have none.
- The judge's reserved-part answer kind and its fields in the LangDrift Core contracts (spec 0049 removes the contract). The reserved part is not something the judge decides in the simulation.

**Build**, each shown and missing:

- The dev-ticket stage between Implement and Review, always run, never a checkbox. Serves `console › waiting    dev-ticket · the person writes the reserved part`.
- The person's "done", read inside `dft agent run` after the `waiting` line. Until the owner decides otherwise, it is the same input as a checkpoint confirmation: one line on standard input, read without echo in a terminal (spec 0057). Serves the moment `closed` prints.
- No stage-driver call between the two lines, enforced and recorded as `paidCalls: 0`. Serves `✓ closed     dev-ticket · written without an agent · no paid call`.
- The Safe Seat decision. Until the owner decides otherwise, it is met when the run's dev ticket closed with no paid call during it, printed after the checks pass. Serves `✓ Safe Seat  met`.
- The run record's dev-ticket state and Safe Seat result. Serves the Dev ticket row and the SAFE SEAT field.

**Change:** none. Nothing of this element exists in `lib/cli` today.

**Interfaces with neighbours:**

- From spec 0059: implement's report with `stoppedAtReservedPart`. The stage starts only after it.
- From spec 0057: the stage's turn, and the shared reading of a line on standard input.
- To spec 0058: no call while the stage is active.
- To spec 0059: the closing, after which `code-review` runs.
- From spec 0064: the moment `✓ checks     4 of 4 passed` has printed.
- To spec 0065: the stage's state and the Safe Seat result, kept in the run's record.

**Failure paths the simulation does not show:**

- When standard input ends while the dev ticket waits, the run stops after the `waiting dev-ticket` line. It prints nothing more, makes no call, and records the stage as `active`.
- A run that stops before the checks pass prints no Safe Seat line, and SAFE SEAT stays unset.
- Spec 0048's single stop rule applies.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** The test feeds one line on standard input after `waiting    dev-ticket`, as it does for each confirmation (spec 0057). This spec adds:
  - the two dev-ticket lines between `implement` and `code-review`, and `✓ Safe Seat  met` between `✓ checks     4 of 4 passed` and `runtime    stopped …`, byte for byte;
  - that the fixture's fake coding agent is not called between the two dev-ticket lines;
  - the run's record: the dev ticket `done`, `closedBy: person`, `paidCalls: 0`, and `safeSeat: met`;
  - a variant whose input ends at the dev ticket: the run prints through `console › waiting    dev-ticket · the person writes the reserved part` and nothing after, with no code-review call and no Safe Seat line.
- A unit test of the Safe Seat rule: met when the dev ticket closed with no paid call during it, and no result at all when the run did not reach the checks. Prior art: the run tests with a fake stage driver.

## Out of Scope

- What the reserved part is and how implement stops at it: the client's skills and their report (spec 0059).
- The checks (spec 0064), the runtime lines around Safe Seat (spec 0062), and the pull request (spec 0063).
- This repository's own Safe Seat policy for its tickets, which is not product behavior.
- The console and its stub API.

## Open Questions

- **How `dft` knows the person finished.** The simulation shows the wait and the close, not the signal between them. Until the owner decides, the person enters a line in the run's terminal, as for a confirmation. A commit on the branch, a change in the reserved part's paths, or another signal are all possible.
- **How "met" is decided.** The simulation prints "met" after the checks pass and shows "whether the floor was met, never a share", but never shows the measure. It also never shows a run where Safe Seat is not met, or what that run would print. Until the owner decides, met means the dev ticket closed with no paid call during it. A floor measured on the diff, the checks' result, or both may be wanted.
- **Who writes the dev ticket.** The line says "the person", while confirmations say "the owner". Whether the dev ticket must be written by a particular person, and whether `dft` should record who, is open. `dft` records `person`.

## Further Notes

- "Written without an agent · no paid call" states what the run did: it called no coding agent and paid for nothing during the stage. It does not observe the person's own tools.
- The glossary's **Safe Seat** entry, "keeps a person writing a part of every ticket by hand", is what this element makes true in the run.
