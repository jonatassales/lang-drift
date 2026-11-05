# The run leaves evidence and a usage record

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E18.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 4 of 19, after spec 0050 and before spec 0053. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation, once the runtime is ready, the run closes its implementation with one line: `✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record`. Its figures come from what the coding agent's calls actually used, 3,700K input and 240K output tokens at medium effort on Claude Sonnet 5.5. Beside it, the Workflow card after the run shows every stage's state, the pull request and Safe Seat, which some record must hold, and the next plan's RAG cites episodes that some earlier run must have left.

The LangDrift CLI prints the completed line in the simulation's format, but it writes much more than any line or field shows: a planning spend entry with its own reservation, laps with their own records, fences and stage evidence files, provisional billing, a proposed ticket status, a run lock, observed resource peaks, facts, and a per-runtime memory with prose that the next plan re-embeds. It keeps no single record of the run that the Workflow card can be filled from, and the elements that now produce run data (stage states, the effort record, the dev ticket and Safe Seat, the stage reports, the delivery, the runtime) have nowhere agreed to put it.

## Solution

The run leaves two records of its own in the runtime's folder, and nothing else:

- **The usage record**, `usage.json`: what each paid stage call used, as the stage driver reports it, and the run's totals. Spec 0056 prices the totals into the observed spend, and this element prints the completed line with that figure and the reservation.
- **The run's record**, `run.json`: the run's evidence and the Workflow card after the run, in one place. Each element that produces run data writes its section through this element as its line prints: the applied choice and stage states (spec 0057), the stage reports' figures (spec 0059), the effort record (spec 0060), the dev ticket and Safe Seat (spec 0061), the runtime (spec 0062), the delivery (spec 0063), and the outcome.

The money ledger (spec 0056), the execution ledger (spec 0062) and the episode (spec 0052) are their owners' records; together with these two they are all the run leaves. Laps, planning spend, billing, ticket status proposals and the per-runtime memory are removed.

## Responsibility

This element:

- keeps the usage record from the stage driver's call results;
- keeps the run's record, with the sections its neighbours fill, and its outcome;
- prints `✓ completed  implementation · … · usage record`;
- tells spec 0052 when a run finished, so that the episode is written, and spec 0048 whether an intent has a completed run.

It does not:

- run the coding agent or read its output (spec 0058);
- price tokens, reserve or release (spec 0056);
- decide any section's content (specs 0057 and 0059 to 0063);
- write the episode or its index (spec 0052);
- write the plan's card data or the recorded intent (spec 0050).

## Produces

Under `.langdrift/microruntimes/runtimes/<runtime id>/`, the runtime's existing folder (`book-a-demo-runtime`):

### The run's id

Every `dft agent run` gets one id when it starts, before its first line: the runtime id and a four-digit sequence, `<runtime id>-<NNNN>`, numbered from the highest run of that runtime in the money ledger plus one (`book-a-demo-runtime-0001` for the first). It names the money ledger's file `finops/ledger/<run>.ledger.json` (spec 0056), the private folder `agent/state/runs/<run>/` with the stage reports, call logs and execution ledger (specs 0058, 0059, 0062), and the stage driver's `{{run}}` and `LANGDRIFT_RUN` (spec 0058). It prints nowhere and is kept in `run.json` as `run`.

### `usage.json`, the usage record

One entry per `dft agent run` that made at least one paid call, appended, matching the money ledger's entry for the same run (spec 0056). Each entry holds:

- `model` and `effort`, the applied choice (`claude-sonnet-5-5`, `medium`);
- `calls`: one per stage-driver call, in order, with `stage`, `skill` (`grilling`, `to-spec`, `to-tickets`, `implement · tdd`, `code-review`), `inputTokens` and `outputTokens` (integers, or `null` when the call reported none);
- `inputTokens` and `outputTokens`, the run's totals: `3700000` and `240000` in the fixture at medium effort; `null` when any call's tokens are `null`.

The totals feed spec 0056's observed spend, which feeds the completed line, OBSERVED and RELEASED. The dev ticket makes no call and adds no entry to `calls` (spec 0061 records its `paidCalls: 0`). A run paused at the gate makes no call and adds no entry.

### `run.json`, the run's record

Replaced atomically each time one of its sections changes, so that a receiver can follow the Workflow card live (spec 0050). It holds the latest run of the runtime:

| Section | Written by | Feeds |
| --- | --- | --- |
| `intent` | this element, from spec 0050: the applied ticket, revision, model, effort and human choices | the Workflow card's choices after the run |
| `stages` | spec 0057: per stage `id`, `human`, `state`, `confirmedBy`, and `reviewed` | HUMAN IN THE LOOP's states |
| `stages[].report` | spec 0059: the stage report's figures (for example `files: 14`) | the evidence of each stage line |
| `effort` | spec 0060: `model`, `effort`, `suggestedModel`, `suggestedEffort`, `recordedAt` | the evidence of the `hook` line |
| `stages[dev-ticket]` and `safeSeat` | spec 0061: `state`, `closedBy`, `paidCalls: 0`; `safeSeat: met`, `decidedAt` | the Dev ticket row, SAFE SEAT |
| `runtime` | spec 0062: `revision`, `deployments`, totals, `state`, `removed`, `left` | the evidence of the runtime lines |
| `delivery` | spec 0063: `branch`, `commit`, `pullRequest` with `number`, `title`, `state`; spec 0064: `pullRequest.checks` with `state` (`running`, `passed`), `passed` and `total` | the Workflow header, PULL REQUEST and its state |
| `outcome` | this element: `running`, `paused` (at the gate), `stopped` (spec 0048's rule) or `completed` (the ready line printed) | spec 0048's intent selection, spec 0052's episode. The card's paused result reads the money ledger's `gate` (spec 0056), not this field. |

The money is not copied here: RESERVED, OBSERVED and RELEASED, and the paused result, come from the run's entry in the money ledger (spec 0056). Every value is token-free and in units; what is not known is `null`, never zero.

## Shows

Phase 2 only; it does not print in the over-budget variant. After `runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB` (spec 0062) and before `✓ pushed     feat/0031-book-a-demo-runtime · 3f9c2e1` (spec 0063), a success line:

```
✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record
```

The two figures are spec 0056's observed and reserved for the applied choice: $9.80 of $11.40 for Claude Sonnet 5.5 at medium; for Claude Opus 5.5 at medium, `✓ completed  implementation · $19.60 observed of $22.80 reserved · usage record`. The line is the same with Review ticked or not and with any stages checked, since the dev ticket adds no paid call. This element owns no card field; OBSERVED is spec 0056's.

## User Stories

1. As a team lead, I want the completed line to show what the run observed against what it reserved, so that I see the spend before the delivery.
2. As a team lead, I want the observed spend to come from the tokens the coding agent reported, so that the figure is measured, not estimated.
3. As a team lead, I want the line to name the usage record, so that I know there is a record behind the figure.
4. As an auditor, I want each paid call's stage, skill and tokens in the usage record, so that I can see where the spend went.
5. As an auditor, I want an unreported call to make the run's usage unknown, never zero, so that a missing measurement is never shown as free work.
6. As the owner, I want a run whose usage is unknown to stop before the completed line, so that the line never prints a figure that was not measured.
7. As the owner, I want the dev ticket to add no call, so that the person's work is recorded as no paid call.
8. As the owner, I want a run paused at the gate to leave no usage entry, so that "nothing was paid for" holds in the records.
9. As the owner, I want a run that stops after paid calls to keep those calls in the usage record, so that no spend goes unrecorded.
10. As a console developer, I want one run record holding every Workflow field after the run, so that the card fills from one place.
11. As a console developer, I want the run record replaced as each section changes, so that the card can follow the run live.
12. As a console developer, I want the run's outcome in the record, so that I can tell whether the run is running, paused, stopped or done.
13. As a maintainer of an element, I want to write my section through one writer, so that no two elements race on the same file.
14. As an auditor, I want every stage's report figures kept with the stage, so that each stage line has its evidence.
15. As an auditor, I want the effort record, the dev ticket's state and Safe Seat's result kept, so that each of those lines has its evidence.
16. As the RAG, I want to be told when a run completed, so that its episode is written once, after the ready line.
17. As the owner, I want laps, planning spend, billing, a proposed ticket status and the per-runtime memory removed, so that the run leaves only what the simulation needs.
18. As a developer, I want the records token-free, so that no secret or credential lands in `.langdrift/`.

## Implementation Decisions

### Remove

- The planning spend entry and its reservation in the ledger, and the planning record of the run state. Serve nothing shown; the plan pays for nothing (specs 0051, 0056).
- Laps: the lap records under the private state, the `laps/NNNN.lap.md` files with their `langdrift-lap` fences, the per-lap stage evidence files, the lap steps and the "lap … open · … · resumed" line. Serve nothing shown; a run has no laps and nothing resumes (spec 0048).
- The run lock that refuses a plan or a run while another is open. Serves nothing shown.
- The provisional billing entry (FG-06) in the ledger (spec 0056 removes the gate).
- The proposed ticket status (`in-review`, `ready-to-merge`) and its prose line. Serves nothing shown.
- The per-runtime memory `run.md`, its "Worth remembering" prose and its `langdrift-run` fence with facts, observed peaks and measured tokens, read as the episode source. Serve nothing shown; spec 0052's episode record replaces it.
- The sampling of resource peaks for the run record and the memory. Serves nothing shown (spec 0052 records `measured: null`).
- The run record's fields no section above holds: `kind`, `files`, the lap list, and the stage `summary` text. The checks keep only the state PULL REQUEST shows (spec 0064).
- The sample's target layout `microruntimes/runs/NNNN-<slug>/` with its own `plan.json`, `run.json`, run memory and laps, which no code reads.

### Build

- The usage record, appended per run from the call results spec 0058 hands over, with totals `null` when any call's tokens are `null`. Serves the completed line through spec 0056.
- The run record's sections and one writer through which each neighbour writes its section, replacing the file atomically. Serves the Workflow card after the run.
- The outcome, set to `running` when the header prints, `paused` when the gate pauses (spec 0056), `completed` when the ready line prints (spec 0063), and `stopped` on any stop (spec 0048).
- The completed line, after the runtime's ready line, with spec 0056's observed and reserved figures. Serves `✓ completed  implementation · … · usage record`.
- The finished signal to spec 0052 when the outcome becomes `completed`.

### Change

- The run record moves from a lap-oriented execution record to the Workflow card's data plus each line's evidence.
- Observed usage moves from a lap's figure to the usage record's totals, summed over the run's calls.

### Interfaces

- From spec 0058: each call's stage, skill and tokens (`null` when unreported). To spec 0056: the run's token totals; from spec 0056: `observedCents` and `reservedCents`.
- From specs 0057, 0059, 0060, 0061, 0062 and 0063: their sections, at the moment their lines print.
- From spec 0050: the applied intent.
- To spec 0048: whether an intent's latest run is `completed`. To spec 0052: the finished run, with its ticket and components.
- From spec 0062: the moment the ready line has printed. To spec 0063: the moment the completed line has printed.

### Failure paths

Every one follows spec 0048's stop rule, with no line of its own:

- the run's usage is unknown (a call reported no tokens): the run stops before the completed line, with the usage record holding the calls as reported and `null` totals;
- a record cannot be written: the run stops before the line that needed it;
- any stop sets `outcome: stopped` and keeps every paid call already in the usage record.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** The fixture's fake coding agent reports tokens summing to 3,700,000 input and 240,000 output at medium effort (spec 0058). This spec adds:
  - the completed line, byte for byte, at its position, in every variant that passes the gate;
  - `usage.json` after the run: one entry, Sonnet at medium, five calls (grilling, to-spec, to-tickets, implement · tdd, code-review) and the totals;
  - `run.json` after each variant, compared field by field with the Workflow card after the run, with `outcome: completed`;
  - after the over-budget variant: no usage entry and `outcome: paused`.
- Golden stop: one call reporting no tokens. Standard output ends before the completed line, `usage.json` holds the calls with `null` totals, `run.json` holds `outcome: stopped`, one sentence goes to standard error and the exit status is 1.
- Unit tests only for the totals rule: integers summed, `null` when any call is `null`.
- Prior art: the CLI's records round-trip tests.

## Out of Scope

- Pricing and the money ledger (spec 0056), the episode and its index (spec 0052), the execution ledger (spec 0062).
- The content of each section (specs 0057, 0059 to 0063).
- The console and its stub API.

## Open Questions

- **The run's id.** The simulation shows no run id, yet the money ledger, the private stage folder and the stage driver each need one. This spec chooses `<runtime id>-<NNNN>`, numbered per runtime. The owner may prefer another identity, such as a timestamp or one tied to the intent's revision.
- **Measured resources for the next plan.** No line of the run measures CPU or memory, so every episode the product writes has `measured: null`, and MEASURED RESOURCES will always read "none cited". The owner decides whether the run should measure, which would need a line the simulation does not have yet.

## Further Notes

- The glossary's **Evidence** entry ("the token-free record the LangDrift CLI keeps of every lifecycle operation and inspection") is met by the run record, the usage record and the two ledgers; no separate evidence folder remains.
- The stage driver's private call logs (spec 0058) are not evidence and are not kept by this element.
