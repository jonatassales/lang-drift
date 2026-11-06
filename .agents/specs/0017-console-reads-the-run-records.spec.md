# The console reads a project's run records through a connector

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The console's Tickets, Analysis, Plans and Workflows listings (spec 0014) show what `ex` does with a ticket, but they read a deterministic stub, extub, and the stub's figures. Since spec 0016, `ex` writes the real thing inside each project: `runtime.md` with the `langdrift-plan` fence, `run.json`, the laps, the ledger and the memory. A team that plans and runs tickets with `ex` sees none of it in the console, and the console's four cards and the CLI's records describe one ticket in two shapes that nothing keeps aligned.

## Solution

A connector reads a project's `.langdrift/` records and serves them to the console in the shape its ticket phases already consume. Each runtime folder becomes one ticket: the plan fence fills the Ticket, Analysis and Plan cards, the run record and the ledger fill the Workflow card, and a phase is reached when the record that fills it exists. The console shows a card's unknown values as dashes, exactly as it does for the stub. The connector is the first LangDrift Vision connector the pivot of 2026-10-28 allowed: it brings a tool's data into the console without producing artifacts for it.

## User Stories

1. As a team member, I want the console's Tickets listing to show the tickets `ex` planned in my project, so that the console and the CLI tell one story.
2. As a team member, I want the Analysis card to show the documents the plan read, the episodes it kept and dropped, and the model and effort it suggested, from the plan fence.
3. As a team member, I want the Plan card to show the components, their proposed shares and the envelope, and whether the plan is proposed or authorized.
4. As a team member, I want the Workflow card to show each stage's progress, the estimate and the reservation, what the lap observed, the pull request and its checks, and how the run ended.
5. As a team member, I want a value nothing measured to show as a dash, never as zero, so that an open lap does not read as a free one.
6. As a team member, I want the listing to refresh when `ex` writes a new line of the run, so that I can follow the session from the console.
7. As a team member, I want the console to say where the records come from, so that a stub ticket and a real one are never confused.
8. As a frontend developer, I want one mapping from the records to the console's `Ticket`, tested against the fixtures `ex` pins, so that a change to either side fails a test.
9. As a platform developer, I want the connector to read only, so that the console never writes into a project's folder.
10. As a platform developer, I want the records' field names to stay the ones spec 0016 fixed, so that the connector never needs the CLI's code.

## Implementation Decisions

- **Source.** The connector reads `.langdrift/microruntimes/runtimes/<id>/runtime.md`, `run.json`, `laps/*.lap.md`, `run.md` and `.langdrift/finops/ledger/<id>.ledger.json`, through the project's root the console already knows from its project record. It reads the fences strictly and treats an unreadable record as a ticket at the last phase it can fill.
- **Phases.** `ticket` once `runtime.md` has a plan fence; `analysis` once the fence's `retrieval` or `facts` carry anything; `plan` once `components` and `totals` exist; `workflow` once `run.json` lists a lap.
- **The Ticket.** `number`, `slug`, `title` and `status` from the plan fence's `ticket`; `artifacts` from the sources whose `inferredRole` is `prd`, `design-doc`, `spec` or `guardrail`, keyed by that role; `createdAt` from the fence's first input's `setAt`.
- **The Analysis.** `documentsRead` from the sources' roles; `useCase` from the fact of topic `use-case`, when one exists; `bounded` from `inputs.cpuMillis` and `inputs.memoryMiB`; `budgetCents` from `inputs.modelBudgetCents`; `specComponents` from the facts of topic `component`; `acceptanceCriteria` from the facts of topic `acceptance`; `criteria` from `retrieval.criteria`; `threshold` from `retrieval.threshold`; `episodes` from `retrieval.considered`, each with its `run` as the id and slug, its `similarity`, its `verdict` and, for a dropped one, its `reason` mapped to the console's drop reasons or `null`; `measuredResources` from the peaks of the episodes `sizing.measuredRunsCited` names, as the memory fences record them; `model` from `suggestion.model`, its name and the list prices of the plan's price record; `effort` from `suggestion.effort`. A field whose source is `null` stays `null`.
- **The Plan.** `bounded` from `totals.bounded`; `components` from `components[]`, each with its `path`, the first of its `technologies`, its `kind`, its `port` and its `proposed` shares; `suggested` from `totals.proposed`; `validationErrors` from the `schema` check's `errors`; `revision` from `runtime.revision`; `authorization` `pending` while the plan's `state` is `proposed`, `authorized` after.
- **The Workflow.** The lap is the open one, else the last. `model` and `effort` from its `choice`; `stages` from `run.json`'s `stages`, their `skills` from the workflow mirror and their `progress` from each `state`; `budgetCents` from the ledger's `bounded.perRunCents`; `estimateCents` and `reserveCents` from the lap's ledger line; `inputTokensK` and `outputTokensK` from that line's estimate tokens; `observedCents` from the line's `observed` once the lap ended, else `null`; `releasedCents` from the line's `recovered` once the lap ended, else `null`; `pullRequest` from `delivery.pullRequest`, `delivery.commit` and the counts of `delivery.checks`, or `null`; `outcome` from `run.json`'s `outcome`; `priceDate` from the price record's `readOn`.
- **Transport.** The connector is a route of the console's server that reads the folder on request, with the project's root as the only input, and a file watcher that invalidates the listing when a record changes. extub keeps serving the stub for projects without a folder, and every response says `source: records` or `source: stub`.
- **Writes.** None. Authorizing a plan or confirming a stage stays `ex agent run` and `ex agent confirm`; the console shows the command to run.

## Testing Decisions

A good test observes what the console shows or what a route returns, never the connector's structure.

1. **Mapping** (Node test runner, beside the console's other tests): the fixtures `ex` pins, copied into the console's test data, map to a `Ticket` whose four cards carry the sample's figures; a plan without a run maps to phase `plan` with `workflow` null; an open lap maps to `observedCents` null; an unreadable record maps to the last phase it fills.
2. **Route**: a project folder with the fixtures returns the tickets with `source: records`; a folder without one falls back to the stub with `source: stub`.
3. **Repository checks**: `pnpm turbo run check`, `pnpm turbo run typecheck test`, `drift check dead-code`.

## Out of Scope

- Writing into `.langdrift/` from the console.
- Showing a lap's steps table or the memory's prose; the cards keep the four panels of spec 0014.
- Other connectors (GitHub, Linear, Obsidian), which follow the same reading-only shape later.

## Further Notes

Spec 0016 fixed the records' fields and their names; this spec names which of them each card reads. The fixtures under `packages/ex/src/testdata/langdrift/` are the contract between the two sides: a change to a field changes a fixture, and the mapping test tells.

Spec 0019 moved the four panels of spec 0014 from the listings to the ticket page: a listing entry now shows the ticket and its phase's summary line, and the ticket page shows the panels' bodies. The mapping to `Ticket` this spec describes is unchanged.
