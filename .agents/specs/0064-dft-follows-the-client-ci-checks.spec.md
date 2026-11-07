# dft follows the client's CI checks

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E17.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 19 of 19, after spec 0062, the last. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation, once the pull request is open, the run waits for the client's own CI: `checks     build · check · test · guardrails · running`, then `✓ checks     4 of 4 passed`. Only then does Safe Seat print and the runtime stop.

The LangDrift CLI already prints these two lines in the simulation's format, but it surrounds them with what the simulation never shows: a `--checks-wait` of ten minutes after which the checks are left as `checks     unknown · still running after …`, other `checks unknown` lines when there is no pull request, `gh` is missing or the wait is interrupted, a `✗ checks … failed` line, and a run that goes on to "for review" when a check failed. The checks it read are also written into the run's record, where no field shows them.

## Solution

After the pull request is open, `dft agent run` reads the pull request's checks from the forge with `gh`, prints the running line once it sees them, with their names in the forge's order, waits until every one has completed, and prints `✓ checks     N of N passed` when all passed. It waits as long as the checks take. Anything else (a failed, cancelled or skipped check, an interrupt, a forge that cannot be read) is a stop under spec 0048's rule, with no line of its own.

## Responsibility

This element reads the client's CI checks of the run's pull request, prints the two checks lines, records the checks' state that PULL REQUEST shows, and decides whether the run may go on to Safe Seat.

It does not:

- open the pull request or mark it ready (spec 0063);
- define, run or change the client's CI: the checks are the client's own;
- decide Safe Seat (spec 0061) or take the runtime down (spec 0062);
- own PULL REQUEST: spec 0063 does, and shows the state this element records;
- keep the checks' names or per-check results: no shown field holds them.

## Produces

- **The checks' state, in the run's record (spec 0065)**, beside the pull request: `pullRequest.checks` with `state` `running` from the running line, then `passed` with `passed` (`4`) and `total` (`4`) from the passed line. It feeds PULL REQUEST's state, `checks running` then `4 of 4 checks passed` (spec 0063).

The checks' names and per-check results are read, printed and used to decide, and are not kept, because no card field shows them.

## Shows

Phase 2 only; neither prints in the over-budget variant. After `✓ pr #19     feat(demo): run book a demo as a microruntime` (spec 0063), an output line:

```
checks     build · check · test · guardrails · running
```

Then, a success line, before `✓ Safe Seat  met` (spec 0061):

```
✓ checks     4 of 4 passed
```

Both are the same with Review ticked or not and with any stages checked. The names and the count are the forge's: the fixture's pull request has four checks named `build`, `check`, `test` and `guardrails`. On the Workflow card, PULL REQUEST's state reads `checks running` from the running line and `4 of 4 checks passed`, in the success tone, from the passed line, until the ready line replaces it (spec 0063). The card has no field of its own for the checks.

## User Stories

1. As a developer, I want the run to wait for my CI's checks on the pull request, so that the delivery is only called ready when my own CI agrees.
2. As a developer, I want the running line to name each check, so that I know which checks the run is waiting on.
3. As a developer, I want the checks named as my forge names them, in its order, so that the terminal matches what I see on the pull request.
4. As a developer, I want `✓ checks     4 of 4 passed` only when every check passed, so that the line is a fact.
5. As a developer, I want the run to wait as long as my checks take, so that a slow CI is not turned into an unknown.
6. As a developer, I want to stop the wait with an interrupt, so that I decide when waiting is no longer useful.
7. As a team lead, I want a failing check to stop the run before Safe Seat and before the pull request is marked ready, so that a red pull request is never presented as ready.
8. As a team lead, I want the stopped run to leave the pull request a draft, so that the team sees it is not ready.
9. As the owner, I want no `checks unknown` or `✗ checks` line, so that the terminal shows only the simulation's two lines.
10. As the owner, I want no `--checks-wait` flag, so that the run has no flag the simulation does not show.
11. As a developer, I want the checks read through the same forge tool as the pull request, so that I install nothing more.
12. As an auditor, I want only the checks' state that PULL REQUEST shows recorded, so that the records stay the simulation's.
13. As a developer, I want the runtime taken down when the checks stop the run, so that a failed CI never leaves a cluster running.

## Implementation Decisions

### Remove

- The `--checks-wait` flag and its ten-minute default (spec 0048 removes the flag), and the deadline in the wait. Serves the passed line, reached however long the checks take.
- `checks     unknown · …` in every form: still running after the wait, the wait interrupted, no pull request, `gh` not installed, the forge's answer unreadable. Serve nothing shown; each is a stop.
- `✗ checks     N of M passed · … failed`. Serves nothing shown; a failed check is a stop.
- The checks' names and per-check results kept in the run's record's delivery. Serve no field; only the state PULL REQUEST shows is kept.
- The "for review when delivery is not clean" outcome after failed checks (spec 0063 removes the rule).

### Build

- After the pr line, polling the pull request's checks through the forge port (`gh` reading each check's name and state) at the existing poll interval, until the first answer that lists at least one check. Then the running line, with every listed name in the forge's order, joined by ` · `.
- Polling on until every listed check has completed, then the passed line with the number passed and the number listed, printed only when they are equal.
- A check that completes as failed, cancelled or skipped is a stop under spec 0048's rule: no line, the runtime taken down (spec 0062), the pull request left a draft (spec 0063).
- An interrupt, or a forge answer that cannot be read, is a stop under the same rule.

### Change

- The wait has no deadline: it ends when the checks end or the person interrupts.
- The run goes on to Safe Seat only after every check passed.

### Interfaces

- From spec 0063: the pull request's number, once the pr line has printed; the forge port, extended with the read of checks.
- To spec 0061: the moment `✓ checks     4 of 4 passed` has printed.
- To spec 0065: the checks' state, written into the run's record through its writer. To spec 0063: that state, for PULL REQUEST.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** The scripted forge first answers with no check, then with the four checks pending, then with them all passed. This spec adds:
  - the running line, printed once, with `build · check · test · guardrails`, and the passed line, byte for byte, at their positions;
  - the checks' state in the run's record, `running` at the running line and `passed`, 4 of 4, at the passed line, matching PULL REQUEST's state;
  - no read of checks in the over-budget variant.
- Golden stops: one check failed, and one check skipped. Each leaves standard output ending on the running line, takes the runtime down, leaves the pull request a draft, writes one sentence to standard error and exits 1.
- No unit test: the two rules (names in order, passed only when all passed) are covered by the golden cases.
- Prior art: the existing delivery observer tests of `gh` output, with a scripted process runner.

## Out of Scope

- The client's CI itself: which checks exist and what they run.
- Opening, readying and merging the pull request (spec 0063), Safe Seat (spec 0061), the runtime (spec 0062).
- The console and its stub API.

## Open Questions

- **A check that fails, is cancelled or is skipped.** The simulation shows only four passing checks. Until the owner decides, any of these stops the run before Safe Seat, with the pull request left a draft; the owner may prefer the run to report and go on, or the coding agent to fix and push again.
- **A pull request with no checks.** A project without CI never reports a check, and the run would wait until interrupted. Until the owner decides, it waits; the owner may prefer the passed line with `0 of 0` or a stop.
- **Checks that appear after the running line.** A forge may register checks one by one. The running line names those listed at its first answer, and the passed line counts all of them at the end, so the two could differ. Until the owner decides, the lines follow the forge as it answers.

## Further Notes

- The website's data calls the runtime "the microruntime the checks use", but the checks run on the forge and nothing shows them reaching the local runtime; spec 0062 records that question.
- The `guardrails` check is a check of the client's CI, not something `dft` runs; the former CI Gates workflow that `dft init` offered is removed (spec 0048).
