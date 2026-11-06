# The effort hook records the run's effort and unblocks implement

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E13.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 15 of 19, after spec 0059 and before spec 0061. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the hero simulation, right before implement the run prints "hook effort recorded · implement unblocked". The harness this mirrors never lets `implement` start until an effort analysis has been recorded: the workflow policy lists "the effort hook before `implement`" among its gates. In the run, the analysis is the model and effort the plan suggested and the team applied, so the hook records it and implement may start.

The LangDrift CLI has no such step. Its hooks are generic committed `sh` scripts that a stage lists. Each runs before the stage's skills and prints "hook <script> · exit N". A non-zero exit stops the stage. Effort is not in a hook's environment, and nothing records the run's effort or makes implement depend on it.

## Solution

The effort hook is a built-in step of `dft agent run`, run once at the start of Implement: after Implement's confirmation when Implement is ticked, and before the implement call. It records the run's effort in the run's record: the applied model and effort beside the pair the judge suggested. `dft` starts the implement call only for a run whose effort record matches the applied choice. Once the record is written and the check holds, the hook prints its one line. Generic script hooks go.

## Responsibility

This element records the run's effort before implement, enforces that implement never starts without it, and prints the `hook` line.

It does not:

- choose the model or the effort; spec 0057 applies the choice and spec 0051 suggests it;
- run implement or pass the effort to the coding agent; spec 0058 does;
- print the implement line; spec 0059 does;
- decide where the stage sits or whether Implement waits for a person; spec 0057 does.

## Produces

- **The effort record, in the run's record (spec 0065).** It holds:
  - `model`, the applied price-list id (`claude-sonnet-5-5`);
  - `effort`, the applied effort (`medium`);
  - `suggestedModel` and `suggestedEffort`, the judge's suggestion from the plan (`claude-sonnet-5-5`, `medium`);
  - `recordedAt`, a timestamp.

  It feeds the `hook` line's "effort recorded" and the implement precondition. It is token-free evidence of the effort implement ran at.

## Shows

Phase 2, `dft agent run`, once per run. The `hook` line comes after `console › to-tickets 3 tickets · blocking edges recorded` (spec 0059) and before `console › implement  tdd · …` (spec 0059):

```
console › hook       effort recorded · implement unblocked
```

When Implement is ticked, the line follows the Implement pair (spec 0057):

```
console › waiting    implement · the owner confirms
✓ confirmed  implement · by the owner
console › hook       effort recorded · implement unblocked
```

No card field. On the card the line belongs to Implement's span of lines (spec 0057).

## User Stories

1. As the owner, I want the run's effort recorded before implement starts, so that the most expensive stage always runs at a recorded effort.
2. As the owner, I want implement to be blocked until that record exists, so that no run implements at an effort nobody chose.
3. As the owner, I want the record to hold both the applied pair and the judge's suggestion, so that I can see whether the team followed the recommendation.
4. As the owner, I want the session to say "effort recorded · implement unblocked", so that I see the gate pass before the paid work.
5. As the owner, I want the hook to run after my Implement confirmation when I ticked Implement, so that the effort is recorded only once I let the stage go.
6. As the owner, I want the hook to run whether or not Implement is ticked, so that the gate does not depend on a person being in the loop.
7. As a FinOps reader, I want the recorded effort to match the effort priced in the estimate, so that the gate, the estimate and the work agree.
8. As an auditor, I want the effort record kept in the run's record, so that the run shows at what effort it implemented.
9. As a team member, I want no generic script hooks run by `dft agent run`, so that nothing runs in a stage that the session does not show.
10. As a team member, I want no "hook <script> · exit N" lines, so that the session reads exactly as the simulation.
11. As the console, I want the hook's line to come from `dft` as data, so that I only show it.
12. As the owner, I want a hook that cannot record the effort to stop the run without printing its line, so that the session never claims a gate passed when it did not.

## Implementation Decisions

**Remove**, each beyond what the simulation shows:

- Stage `hooks`, the lists of committed scripts in the workflow. Serves the `hook` line, which is the effort hook alone.
- Running a hook as an `sh` script through the stage driver, with its `LANGDRIFT_HOOK` variable and its private hook log.
- The line "hook <script> · exit N", and the failure "hook <script> of stage <stage> exited N".
- `dft runtime hook` and the hook executor behind it, which run consumer extensions. Spec 0048 removes the command; its adapter goes with it.

**Build**, each shown and missing:

- The effort hook step, run at the start of Implement, after Implement's checkpoint when ticked. Serves `console › hook       effort recorded · implement unblocked`.
- The effort record in the run's record. Serves "effort recorded".
- The implement precondition: the implement call starts only when the run's effort record exists and its model and effort equal the applied choice. Serves "implement unblocked".

**Change**, each different:

- The `hook` line reports the effort gate, not a script's exit status, and it prints once per run, always before `implement`.

**Interfaces with neighbours:**

- From spec 0057: the applied model and effort, and Implement's turn after its checkpoint.
- From spec 0050: the judge's suggestion carried in the plan's hand-off.
- To spec 0058: permission to start the implement call.
- To spec 0065: the effort record's fields, kept in the run's record.

**Failure paths the simulation does not show:**

- If the record cannot be written, or the precondition does not hold, the run stops before the `hook` line under spec 0048's single rule.
- Implement is never called in that case, and nothing is printed instead.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** This spec adds:
  - the `hook` line between `to-tickets` and `implement` in the default run;
  - in the variant with Implement ticked, the line right after `✓ confirmed  implement · by the owner`;
  - the run's effort record: `claude-sonnet-5-5`, `medium`, suggested `claude-sonnet-5-5`, `medium`;
  - that the fixture's fake coding agent finds the effort record already written when it is called for Implement.
- **A unit test of the precondition:** implement is refused without an effort record, and with one whose model or effort differs from the applied choice. Prior art: the run tests that drive a fake stage driver.

## Out of Scope

- The choice and the stage order (spec 0057), and the suggestion (spec 0051).
- The implement call and its line (specs 0058 and 0059).
- The run's record as a file (spec 0065).
- The console and its stub API.

## Open Questions

- **A guard inside the client's coding agent.** The harness the simulation mirrors also guards `implement` inside the coding agent: this repository's guard reads an effort log keyed by the agent's session. Should `dft` also satisfy such a client-side guard for the implement call, and how would the call's session be bound to the record? The simulation shows only `dft`'s own gate.
- **The repository's effort analyzer.** In this repository the effort before `implement` is recorded by the `llm-effort-analyzer` skill and unblocked by a guard hook. Whether the effort hook ships that skill and guard into `.langdrift/` as part of the same workflow, instead of the built-in step this spec describes, is the owner's decision.

## Further Notes

- The hook is named for what it guards, as in the workflow policy's gates. It is not a hook system: no other step of the run can be hooked.
