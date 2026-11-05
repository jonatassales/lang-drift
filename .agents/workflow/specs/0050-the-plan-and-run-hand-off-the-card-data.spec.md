# The plan and the run hand off the card data

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E3.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 3 of 19, after spec 0048 and before spec 0065. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation, `dft agent plan` ends with `✓ plan       runtime.md draft, revision 1 · proposed to the console`: the plan is proposed, and the Plan card reads "pending team authorization" with the EVIDENCE "runtime.md draft · revision 1 · proposed to the console". `dft agent run` then opens with `# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team` and `console › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1`: it applies an intent the team recorded and authorized, not one it decided itself. Beside the terminal, four cards fill with data that `dft` and the client's `.langdrift/` produce.

Today nothing in the plan hands the draft anywhere: it writes `runtime.md` locally and the CLI posts MCP tools to a console endpoint no one serves. The run turns a proposed plan into an authorized one by itself, recording the owner through the CLI, and prints "authorized by owner". The plan's last line adds "authorize with dft agent run". The protocol between the CLI and a console carries none of the cards' data, and an NDJSON mode prints events with card and field ids instead. A console that wants to show the session's panel has no record to read, and a team has no way to see what it is asked to authorize.

## Solution

The hand-off is three records under the runtime's folder in the client's `.langdrift/`, which a console only receives:

- **The plan's card data**, written by `dft agent plan` as its lines print: every field of the Ticket, Analysis and Plan cards and the phase-1 Workflow card.
- **The recorded intent**, written by `dft agent plan` as proposed when `✓ plan …` prints: the ticket, the `runtime.md` revision, the model, the effort and the human-in-the-loop choices. The team authorizes it. `dft agent run` applies only an authorized intent, and never authorizes one itself.
- **The run's record**, kept by spec 0065 as the run's lines print: the Workflow card after the run.

Nothing is sent over a network: the console's way of reading these records is outside these specs. The card is the specification of the data, and each field's value belongs to the element that owns it; this element owns the records' structure, the moments they are written, the authorization state and the three lines that state the hand-off.

## Responsibility

This element:

- writes the plan's card data, gathering each field from its owner as the plan's lines print;
- writes the recorded intent as proposed, at the end of the plan;
- prints `✓ plan       runtime.md draft, revision 1 · proposed to the console`;
- sets the Plan card's "pending team authorization" and EVIDENCE;
- reads the recorded intent the run applies (spec 0048 chooses which one), refuses one that is not authorized, and prints the run's header and `intent` line;
- names the run's record (spec 0065) as the third part of the hand-off.

It does not:

- choose which intent the run applies when there are several, or print the comments (spec 0048);
- compute any card value: the ticket and documents (spec 0053), the guardrail and bound (spec 0054), the judge's statements and suggestion (spec 0051), the episodes (spec 0052), the components (spec 0066), the sizing and the draft (spec 0055), the money (spec 0056), the workflow choices and stage states (spec 0057);
- write `runtime.md` (spec 0055) or the run's record (spec 0065);
- authorize the intent: the team does (see Open Questions);
- render anything: the console receives.

## Produces

Under `.langdrift/microruntimes/runtimes/<runtime id>/`, where the runtime id is the ticket's slug (`book-a-demo-runtime`, spec 0055):

### `plan.json`, the plan's card data

Written by `dft agent plan`, replaced atomically each time a line that fills a field prints, so that a receiver can fill the cards live as the simulation does. A field is absent until its line has printed, with two exceptions the simulation shows from the moment the panel appears, before any line: the ticket's header (number, slug, status, title) and the `workflow` section. The plan writes both first, before `# plan only …` prints, from the ticket spec 0053 resolves and the Workflow card data specs 0056 and 0057 supply. The values are units, never formatted text: millicores and MiB as integers, money as integer cents, similarities as numbers with two decimals; a text field holds only text that varies with the reasoning (titles, reasons, statements). The cards' labels and copy ("Linked artifacts", "recommendation · the team authorizes", the thousands separators) are the receiver's.

| Section | Fields and their owners | Card |
| --- | --- | --- |
| `ticket` | id, number, slug, status, title, written before the first line; the four links with kind, short id and file, at the links line (spec 0053) | Ticket |
| `analysis.documents` | the PRD, design doc and spec rows with their extracted facts, and the documents-read count (spec 0053); the guardrail row with its criteria, plan criteria and enforcement (spec 0054) and the criteria statements (spec 0051) | Analysis DOCUMENTS |
| `analysis.episodes` | threshold, candidates with id, slug, similarity, verdict and reason, the counts (spec 0052, with reasons from spec 0051) | Analysis EPISODES, MEASURED RESOURCES |
| `analysis.suggestion` | header state, model, effort, the weighing reasons, `by: recommendation` (spec 0051) | Analysis header, SUGGESTED MODEL, SUGGESTED EFFORT, "recommendation · the team authorizes" |
| `plan.components` | path, manifest, stack, role, port, image, share, and the cross-check (spec 0066) | Plan COMPONENTS, ports line |
| `plan.sizing` | running and final totals, bound with its source, headroom, suggested resource, limits, validation (spec 0055) | Plan figures, bars, total line, BOUNDED, HEADROOM, SUGGESTED RESOURCE, GPU · DISK, VALIDATION |
| `plan.price` | the runtime has no price source (spec 0056) | Plan PRICE |
| `plan.authorization` | `pending`, set when `✓ plan …` prints (this element) | "pending team authorization" in the Plan header and PROPOSED PLAN |
| `plan.evidence` | `file: runtime.md`, `state: draft`, `revision: 1`, `proposedTo: console`, set when `✓ plan …` prints (this element) | Plan EVIDENCE |
| `workflow` | the stages and the default choice (spec 0057); the price list's models and read date, the expected tokens, the spend table, the budget with its source and the margin (spec 0056); written before the first line | Workflow, phase 1, including its header from the moment the panel appears |

`plan.authorization` and `plan.evidence` are written once, by the plan, and never rewritten: the run does not touch `plan.json` (see Open Questions).

### `intent.json`, the recorded intent

Written by `dft agent plan` once, when `✓ plan …` prints, with:

- `ticket`: the ticket's bare id, `0031-TICKET-book-a-demo-runtime` (spec 0053); feeds the run's header;
- `runtimeId` and `revision`: the draft's, `book-a-demo-runtime` and `1` (spec 0055); feeds "runtime.md revision 1";
- `model`, `effort` and `human`: the workflow's default choice (spec 0057), `claude-sonnet-5-5` at `medium` with Planning and Review ticked, the same choice the Workflow card shows; feeds the `intent` line and every run element;
- `status`: `proposed`.

The team authorizes the intent by setting `status` to `authorized`, and may change `model`, `effort` and `human` before it does; the Workflow card is the run's input. `dft` never sets `authorized`. A plan of the same ticket that writes a new revision writes a new proposed intent for it.

### The run's record

`run.json`, kept by spec 0065 beside `plan.json`, is the third part of the hand-off: the stage states (spec 0057), the pull request and its checks' state (specs 0063, 0064) and Safe Seat (spec 0061), as the run's lines print. RESERVED, OBSERVED and RELEASED, and the paused result, come from the run's entry in the money ledger `finops/ledger/<run>.ledger.json` (spec 0056). This element only names these two as the records a receiver reads for the Workflow card after the run.

## Shows

### Phase 1

After `console › tool call: validate runtime.md draft against the schema · 0 errors` (spec 0055) and before `# next: once the team authorizes, dft agent run applies the recorded intent` (spec 0048), a success line:

```
✓ plan       runtime.md draft, revision 1 · proposed to the console
```

When it prints, the Plan card gains:

- the header's suffix: `PLAN — 1,450m · 2,304 MiB · pending team authorization` (the figures are spec 0055's);
- PROPOSED PLAN's suffix: `1,450m · 2,304 MiB · pending team authorization`;
- EVIDENCE: `runtime.md draft · revision 1 · proposed to the console`.

### Phase 2

The first two lines `dft agent run` prints, in every variant (default, over budget, review unchecked, stages checked):

```
# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team
console › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1
```

The header is a comment; the `intent` line is output. The `intent` line carries the authorized intent's model and effort: with Claude Opus 5.5 at high effort, the over-budget variant's line is `console › intent     claude-opus-5-5 · effort high · runtime.md revision 1`. The `estimate` line follows (spec 0056).

After phase 2 the Plan card still reads "pending team authorization" and its EVIDENCE is unchanged.

## User Stories

1. As a team lead, I want the plan to end by proposing a draft, not by running anything, so that the team decides before anything is built or paid for.
2. As a team lead, I want the plan's last result to say `runtime.md draft, revision 1 · proposed to the console`, so that I know which revision I am asked to authorize.
3. As a team lead, I want the Plan card to read "pending team authorization", so that nobody mistakes the proposal for an authorized plan.
4. As a team lead, I want the Plan card's EVIDENCE to name the draft, its revision and that it was proposed, so that the proposal points at the file it stands on.
5. As a team, I want the plan to record a proposed intent with the workflow's default model, effort and human-in-the-loop choices, so that authorizing it is a decision on concrete values.
6. As a team, I want to change the model, the effort and the human-in-the-loop choices before authorizing, so that the Workflow card is the run's input, as in the simulation.
7. As a team, I want `dft` never to authorize an intent itself, so that a recommendation is never an authorization.
8. As a developer, I want `dft agent run` to open with `# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team`, so that I see which ticket runs and on whose authority.
9. As a developer, I want the `intent` line to show the authorized model, effort and `runtime.md` revision, so that I see exactly what the run applies before money moves.
10. As a developer, I want a run of an intent that is still proposed to stop without a line, so that nothing runs before the team authorizes it.
11. As a developer, I want a run whose intent names a revision that is no longer the draft's to stop, so that the run never applies a plan the team did not see.
12. As a console developer, I want every card field in a record under `.langdrift/`, so that the console only receives and never reasons.
13. As a console developer, I want the plan's card data replaced as each line prints, so that the cards can fill live as the simulation's do.
14. As a console developer, I want a field to be absent until its line has printed, so that the console can show "Not set yet" and "Resolving" honestly.
15. As a console developer, I want values in units (millicores, MiB, cents), so that every locale can format them.
16. As a console developer, I want the Workflow card after the run in the run's record, so that one place holds the run's state as it happens.
17. As a platform engineer, I want the hand-off to be files, so that no connection, credential or endpoint stands between `dft` and the data.
18. As the owner, I want the MCP client, the connected agent and the NDJSON mode removed, so that the product carries no channel the simulation does not show.
19. As an auditor, I want the proposed intent and its authorized status kept beside the draft, so that what ran can be traced to what the team authorized.
20. As the owner, I want the Plan card to keep "pending team authorization" after the run, as the simulation shows, until I decide the contradiction, so that the product follows its definition literally.

## Implementation Decisions

### Remove

- The run's own authorization: `dft agent run` turning a proposed plan into an authorized one, recorded as `By` owner and `Via` cli, and the header text "authorized by owner". Serves `# run · … the recorded intent, authorized by the team`.
- The suffix "· proposed · authorize with dft agent run" of the plan's last result. Serves `✓ plan       runtime.md draft, revision 1 · proposed to the console`.
- The MCP client that posts tools to `<console-url>/api/mcp`, the connected agent's infrastructure, and the commands that drive it (spec 0048 removes the commands). Serves nothing shown; the hand-off is files.
- The NDJSON renderer, its records with card and field ids, the line-to-card reveal map, and `--json` (spec 0048 removes the flag). Serves nothing shown: the cards' data is in the records, and the terminal's text is the human renderer's.
- The protocol's tools, events, leases, bindings, claims, measurement windows and sync (spec 0049 removes the contract). Serves nothing shown.
- The `langdrift-plan` fence inside `runtime.md`, which carried collection limits, context hashes, source hashes, input records, discrepancies and contract gaps. Serves nothing shown; the card data moves to `plan.json` with only the cards' fields, and `runtime.md` keeps its `langdrift` fence (spec 0055).

### Build

- `plan.json`, structured per card as Produces says, replaced atomically as each line that fills a field prints. Serves every card field of phase 1.
- `plan.authorization: pending` and `plan.evidence`, set when `✓ plan …` prints. Serves "pending team authorization" in the Plan header and PROPOSED PLAN, and EVIDENCE.
- `intent.json`, written as proposed when `✓ plan …` prints. Serves the run's header and `intent` line.
- The `✓ plan       runtime.md draft, revision 1 · proposed to the console` line, printed after the draft is written (spec 0055) and the intent recorded.
- At the start of `dft agent run`, the intent spec 0048 selects is read and checked: `status` must be `authorized`, `revision` must be the draft's current revision, and `model` and `effort` must be in the price list and the efforts (specs 0056, 0057). Any failure is a stop under spec 0048's rule, before the header prints.
- The run's header and `intent` line, from the intent's `ticket`, `model`, `effort` and `revision`. Serves both phase-2 lines.
- The applied intent is handed to the run's elements: the choice to specs 0056, 0057 and 0058, the revision to spec 0062.

### Change

- "Authorized by" is no longer a recorded actor's name: the header always says "the team", because only an authorized intent runs and authorization is the team's.
- The plan's hand-off changes from a local `runtime.md` plus an unserved MCP post to `runtime.md` plus `plan.json` and `intent.json`.
- The `intent` line's prefix becomes `console › ` (spec 0048).

### Records and interfaces

- The shapes of `plan.json`, `intent.json` and `run.json` are LangDrift Core contracts (spec 0049), replacing the protocol package. Each owner fills its fields through one card-data writer this element provides; no owner writes the files directly.
- From spec 0048: the intent the run applies. To spec 0048: every intent under `.langdrift/microruntimes/runtimes/` with its status and whether its run completed (spec 0065 says when a run completed).
- From specs 0051 to 0057 and 0066: their card fields, at the moment their lines print; before the first line, the ticket's header from spec 0053 and the Workflow card data from specs 0056 and 0057.
- From spec 0055: the draft's path and revision, once written.
- To spec 0065: the applied intent, which the run's record keeps as applied (spec 0057).
- A plan that stops leaves `plan.json` as far as it got and writes no intent (spec 0048's rule).

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** Between `dft agent plan` and `dft agent run`, the test authorizes the recorded intent by setting its `status` to `authorized`, and for the variants also its choices: Claude Opus 5.5 at high effort for over budget, Review unticked for to-merge, Spec, Tickets and Implement ticked for all-stages. This spec adds:
  - the `✓ plan …` line at its position, and the run's header and `intent` line in every variant, byte for byte;
  - `plan.json` after the plan, compared field by field with the Ticket, Analysis, Plan and phase-1 Workflow cards of the simulation, including `authorization: pending` and the EVIDENCE fields;
  - `plan.json` unchanged after the run;
  - `intent.json` after the plan: the ticket id, `book-a-demo-runtime`, revision 1, `claude-sonnet-5-5`, `medium`, Planning and Review ticked, `proposed`;
  - `run.json` after each run variant, compared with the Workflow card after the run (spec 0065 keeps it).
- Golden stops: `dft agent run` with the intent still `proposed`, and with an intent whose revision is not the draft's. Each prints nothing on standard output, writes one sentence on standard error and exits 1.
- One test that captures `plan.json` after every printed line of the plan and checks that each field appears at the line the simulation reveals it (for example "pending team authorization" and EVIDENCE only at `✓ plan …`), and that the ticket's header and the `workflow` section are present before the first line.
- No unit test beyond these: the records' shapes are tested in spec 0049.

## Out of Scope

- Every field value: specs 0051 to 0057, 0065 and 0066.
- The comments, the commands and the choice of intent (spec 0048).
- The contracts' definitions (spec 0049).
- The console, its stub API, how it reads these records and how it shows them.

## Open Questions

- **How the team authorizes without a console that sends.** The simulation shows the intent arriving "authorized by the team" but not how. The console only receives; `dft` keeps three commands, none of which authorizes. Until the owner decides, authorization is the change of `intent.json`'s `status` to `authorized`, made outside `dft`, and the team's changes to the choices are made in the same record.
- **The Plan card after the run** (research note §1.6, contradiction 1). The run's header says "authorized by the team" while the Plan card still reads "pending team authorization". This spec follows the simulation: `plan.json` is never rewritten after the plan. The owner decides whether the card should change once the team authorizes.
- **"Proposed to the console" without a console connection.** The line says the draft is proposed to the console, but nothing is sent; the console reads the records. The owner may keep the words or change them in the simulation first.

## Further Notes

- The records sit in the runtime's existing folder, beside `runtime.md`, so the layout `.langdrift/microruntimes/runtimes/<runtime id>/` stays.
- Rule 12, item 2, asks that an execution started with `dft` appear in the console's runtime view; these records are what makes that possible without a connection.
