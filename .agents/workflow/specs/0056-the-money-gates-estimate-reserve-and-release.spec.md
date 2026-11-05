# The money gates estimate, reserve and release a run's model spend

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E9.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 7 of 19, after spec 0054 and before spec 0052. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The hero simulation shows one money mechanism and only one. The plan states the run's budget, "set by the team in the design doc", and that the runtime itself has no price source. The Workflow card prices the chosen model and effort at dated list prices, shows the estimate and what the gate would reserve of the budget, and warns when the reservation would exceed it. The run prints the estimate, the budget and the gate, opens the gate or pauses with nothing paid, and ends with what it observed and what the reservation released.

The LangDrift CLI already computes the estimate, the reservation and the release with the simulation's formulas, but it surrounds them with gates, inputs and records the simulation never shows: six FinOps gates (FG-01 to FG-06), ticket allocations, a period envelope, provisional billing, a reservation for planning, a margin the owner sets at init, a budget "set by owner" from a CLI input, and a refusal when a forecast is missing. It also accepts any model in a price file and never records, at plan time, the estimate and reservation the Workflow card shows. A team that reads the card cannot trust that the run's gate is the one the card describes.

## Solution

The money gates do what the simulation shows and nothing else. Three Claude models carry their list prices, read on 2026-10-06 and uncached. Expected and observed tokens per effort are known. The budget is the one the design doc sets, $25.00 per run for model calls. The plan prices every model and effort pair, so the Workflow card only looks up the pair the team picks. The run applies the recorded pair: it prints the estimate and the budget, reserves the estimate plus a fixed 20% margin, opens the gate when the reservation fits the budget and pauses with nothing paid when it does not. At the end it prices what the stages observed and releases the difference. Every other gate, input and record goes.

## Responsibility

This element computes and records a run's money: the estimate, the budget, the reservation, the gate decision, the observed spend and the release. It owns the price list and the token forecasts under `.langdrift/finops/`, the spend table the plan hands off, and the run's money record.

It does not:

- read the budget out of the design doc; spec 0053 (project documents) extracts it, and this element states it in the plan's `budget` line, after spec 0054's `bounded` line;
- measure the tokens the stages use; the stage driver (spec 0058) reads each call's usage and spec 0065 (evidence and ledger) keeps the usage record;
- carry the hand-off to the card; spec 0050 does;
- choose the model and effort; the Workflow card's choice and its default are spec 0057's, and the suggestion is the judge's (spec 0051);
- print the lines that only quote its figures: `✓ completed implementation …` (spec 0065), `✓ ready … released to the envelope` and the phase-2 Workflow header (spec 0063).

## Produces

Under the client's `.langdrift/`:

- **`finops/prices.json`, the price list.** It replaces the numbered `finops/prices/` folder and its mirror header and holds exactly:
  - `readOn`: `2026-10-06`; feeds the MODEL price line's date ("list price, Oct 2026");
  - `cache`: `uncached`; feeds "uncached list prices";
  - `currency`: `USD`, and `unit`: cents per million tokens;
  - `models`: three entries with `id`, `name`, `shortName`, `inputCents` and `outputCents`:
    - `claude-opus-5-5`, Claude Opus 5.5, Opus 5.5, 400, 2000;
    - `claude-sonnet-5-5`, Claude Sonnet 5.5, Sonnet 5.5, 200, 1000;
    - `claude-haiku-5-5`, Claude Haiku 5.5, Haiku 5.5, 10, 50.

  `id` feeds the `intent` line (spec 0050) and the judge's suggestion; `name` feeds the MODEL choice (spec 0057) and the Analysis card; `shortName` feeds the Workflow header; the cents feed the MODEL price line and every figure below.
- **`finops/forecasts.json`, the expected tokens.** It keeps its place, drops its mirror header and its `status` and `illustrative` flags, and holds, in thousands of tokens, per effort: low 2,000 input and 120 output; medium 3,500 and 250; high 6,000 and 500. It feeds the estimate and the "{input} input · {output} output" figures.
- **The spend table, in the plan's hand-off.** `dft agent plan` writes, in the `workflow` section of the plan's card data (`plan.json`, spec 0050), before its first line prints, one row per model and effort pair with the pair's expected tokens, `estimateCents`, `reservedCents` and `overBudget`. It also writes `budgetCents` (2500) with its source, the design doc's path, and `marginPercent` (20). Every phase-1 Workflow figure is a lookup in this table, and all of them exist from the moment the panel appears, before any line, as the simulation shows; the design doc's budget is read that early by spec 0053. For the simulation's ticket it reads:

  | Model | Effort | Estimate | Reserved | Gate |
  | --- | --- | --- | --- | --- |
  | Opus 5.5 | low | $10.40 | $12.48 | opens |
  | Opus 5.5 | medium | $19.00 | $22.80 | opens |
  | Opus 5.5 | high | $34.00 | $40.80 | pauses |
  | Sonnet 5.5 | low | $5.20 | $6.24 | opens |
  | Sonnet 5.5 | medium | $9.50 | $11.40 | opens |
  | Sonnet 5.5 | high | $17.00 | $20.40 | opens |
  | Haiku 5.5 | low | $0.26 | $0.32 | opens |
  | Haiku 5.5 | medium | $0.48 | $0.58 | opens |
  | Haiku 5.5 | high | $0.85 | $1.02 | opens |

- **The run's money record, in `finops/ledger/<run>.ledger.json`**, where `<run>` is the run's id (spec 0065). One entry per `dft agent run`, with `model`, `effort`, the expected tokens, `estimateCents`, `budgetCents`, `marginPercent`, `reservedCents`, `gate` (`opened` or `paused`), and, once the run reaches them, the observed tokens, `observedCents` and `releasedCents`. It feeds RESERVED, OBSERVED and RELEASED, the paused result, the `completed` line's two figures (spec 0065), and the `ready` line's released figure and the phase-2 Workflow header's observed figure (spec 0063). The observed tokens come from the usage record that spec 0065 keeps from the stage driver's calls (spec 0058). With the simulation's fixture, medium on Sonnet observes 3,700K input and 240K output: $9.80 observed, $1.60 released.

The formulas, in integer cents:

- **Spend:** `round((inputTokens × inputCents + outputTokens × outputCents) / 1,000,000)`, rounding a half up. With tokens in thousands it is the simulation's `round((inputThousands × inputCents + outputThousands × outputCents) / 1000)`. It gives the estimate from the forecast and the observed spend from the usage record, both at the chosen model's price.
- **Reservation:** `ceil(estimate × (100 + 20) / 100)`.
- **Gate:** paused when the reservation exceeds the budget, opened otherwise; a reservation equal to the budget opens.
- **Release:** `reserved − observed`.

## Shows

Phase 1, `dft agent plan`. The budget line follows `bounded …` (spec 0054), and the estimate line follows `limits …` (spec 0055):

```
budget     $25.00 per run for model calls · set by the team in the design doc
estimate   runtime has no price source · cpu and memory stay in millicores and MiB
```

Phase-1 card fields:

- **Plan, PRICE:** `runtime: no price source · units only`, revealed by the `estimate` line.
- **Workflow header**, from the moment the panel appears and through the run until spec 0063's ready line replaces it: `Sonnet 5.5 · medium · $9.50 of $25.00`, that is `{short name} · {effort} · {estimate} of {budget}` for the chosen pair, which is `Opus 5.5 · high · $34.00 of $25.00` for Opus at high effort.
- **The MODEL price line**, under the choice: `$2.00 / $10.00 per MTok · list price, Oct 2026` for Sonnet; `$4.00 / $20.00 per MTok · list price, Oct 2026` for Opus; `$0.10 / $0.50 per MTok · list price, Oct 2026` for Haiku.
- **ESTIMATED SPEND:** `$9.50`, then `3.5M input · 250K output · uncached list prices`, then `reserves $11.40 of the $25.00 budget`, with a bar at the reservation's share of the budget. Token counts are compact: `2M` and `120K` at low effort, `6M` and `500K` at high.
- **ESTIMATED SPEND, over budget:** the estimate, the reservation and the bar take the warning tone, and `over the budget · the gate will pause this run` follows. For Opus at high effort it reads `$34.00`, `6M input · 500K output · uncached list prices`, `reserves $40.80 of the $25.00 budget` and `over the budget · the gate will pause this run`.

Phase 2, `dft agent run`. These lines follow the `intent` line (spec 0050), here for Sonnet at medium effort:

```
console › estimate   $9.50 · 3.5M input · 250K output · list prices · FIN-02
console › budget     $25.00 per run · set by the team · FIN-03
console › gate       reserved $11.40 of $25.00 · opened
```

The waiting line for planning (spec 0057) comes next.

Phase 2, over budget, here for Opus at high effort. The gate's two lines replace everything from the gate on, and the run ends after them:

```
console › estimate   $34.00 · 6M input · 500K output · list prices · FIN-02
console › budget     $25.00 per run · set by the team · FIN-03
console › gate       reservation $40.80 exceeds the $25.00 budget · paused
# needs a decision: raise the budget or lower the effort · nothing was paid for
```

Phase-2 card fields:

- **RESERVED:** `$11.40` once the `gate … opened` line has printed; unset (a dash) before.
- **OBSERVED:** `$9.80` once `✓ completed implementation …` (spec 0065) has printed.
- **RELEASED:** `$1.60` once `✓ ready …` (spec 0063) has printed.
- **Paused run:** the Workflow header reads `paused at the gate`, and the result reads `paused at the gate · needs a decision`, then `reserves $40.80 of the $25.00 budget · nothing paid`. RESERVED, OBSERVED and RELEASED are not shown.

## User Stories

1. As a team member reading the plan, I want the run's budget stated with its source, "$25.00 per run for model calls · set by the team in the design doc", so that I know which document bounds the spend.
2. As a team member, I want the plan to say that the runtime has no price source and that cpu and memory stay in millicores and MiB, so that nobody mistakes the resource figures for money.
3. As a team member reading the Plan card, I want PRICE to read "runtime: no price source · units only", so that the card says the same as the terminal.
4. As a team member choosing a model, I want each model's input and output list price per million tokens and the date the prices were read, so that I know how current they are.
5. As a team member, I want only Claude Opus 5.5, Sonnet 5.5 and Haiku 5.5 offered, priced from one list, so that the choice and the price come from the same source.
6. As a team member choosing an effort, I want the expected input and output tokens beside the estimate, so that I see what the dollar figure is made of.
7. As a team member, I want the estimate and the reservation it would hold against the $25.00 budget before I play the run, so that I can pick a pair the gate will open for.
8. As a team member, I want a pair whose reservation exceeds the budget flagged with "over the budget · the gate will pause this run", so that I am not surprised by a pause.
9. As a team member, I want the Workflow header to summarize my pair as "Sonnet 5.5 · medium · $9.50 of $25.00", so that the card's title tells me the cost of my choice.
10. As the console, I want every pair's estimate, reservation and gate outcome already computed by the plan, so that I only show data and never reason about money.
11. As the owner running `dft agent run`, I want the estimate printed with its tokens and the FIN-02 criterion it satisfies, so that I see the estimate came before any paid call.
12. As the owner, I want the budget printed with "set by the team" and FIN-03, so that I see the run is bounded before it spends.
13. As the owner, I want the gate to reserve the estimate plus 20% and say "reserved $11.40 of $25.00 · opened", so that I know how much of the budget the run holds.
14. As the owner, I want a reservation over the budget to stop the run at the gate with "nothing was paid for" and the decision it needs, so that no model call is made outside the budget.
15. As the owner, I want the paused run to tell me to raise the budget or lower the effort, so that I know the two ways forward.
16. As a team member watching the card, I want RESERVED to appear when the gate opens, OBSERVED when the usage record closes and RELEASED when the run is ready, so that each figure appears when it becomes true.
17. As a team member, I want what was observed priced at the same list price as the estimate, so that estimate and observation compare like for like.
18. As a team member, I want the release to be the reservation minus what was observed, so that "$1.60 released to the envelope" is exact.
19. As an auditor, I want each run's money record kept in the project's ledger with the model, the effort, the tokens and every figure in cents, so that I can recompute every line.
20. As a developer of `dft`, I want one spend formula with one rounding rule for the estimate and the observed spend, so that the card and the terminal never disagree by a cent.
21. As the owner, I want planning to reserve nothing and no gate the simulation does not show to stand between me and the run, so that the money rules I read on the card are all the rules there are.

## Implementation Decisions

**Remove**, each beyond what the simulation shows:

- FG-01, the price-source gate, and its "gate blocked" line. The plan's spend table and the price list's three models make a separate gate pointless. Serves the `gate` line, which has only the opened and paused forms.
- FG-02 and FG-03's ticket-allocation and period-remainder checks, the allocations, `unallocatedCents` and the period envelope. The only bound is the per-run budget. Serves `gate … opened` and `gate … paused`.
- FG-04, the during-run pause when observed spend reaches the reservation or a call reports no usage, with its line and its "allow unmeasured" setting. Serves the run's line sequence, which has no pause after the gate.
- FG-05, earlier overspend, and `--acknowledge-overspend`.
- FG-06, provisional billing, with the ledger's billed, recovered, overspend and variance figures and its provisional state. Serves RELEASED, which is `reserved − observed`.
- The planning spend and its reservation in the ledger and the run record. Serves "nothing is built, started or paid for in this step"; who pays for the judge is spec 0051's open question.
- The reserve margin asked at `dft init` and stored in the budget input. Serves `reserves $11.40 of the $25.00 budget`, which uses a fixed 20%.
- The budget as an input record "set by owner", the `finops/budgets.json` and `finops/gates.json` mirrors, and the mirror and input-record machinery for prices and forecasts. Serves "set by the team" and "set by the team in the design doc".
- The refusal "no forecast is recorded for effort … set it with dft finops forecasts set" and every `dft finops` input command (spec 0048 removes the commands). Serves the estimate line.
- Accepting any model a price file lists. The price list holds the three models and the MODEL choice offers those three.
- The clamp of the release at zero. RELEASED is `reserved − observed`.
- The gate's error text printed after a pause ("…; nothing was paid"). The pause prints its two lines and nothing else.

**Build**, each shown and missing:

- The plan's spend table, its budget with source and its margin, written into `plan.json`'s `workflow` section before the plan's first line. Serves the Workflow header, the MODEL price line and ESTIMATED SPEND, including the over-budget form.
- The plan's `budget` line, from the design doc's budget that spec 0053 extracts. Serves `budget     $25.00 per run for model calls · set by the team in the design doc`.
- The plan's `estimate` line and the Plan card's PRICE field, which state that the runtime has no price source. Serves `estimate   runtime has no price source · cpu and memory stay in millicores and MiB` and `runtime: no price source · units only`.
- `shortName` on each price-list model. Serves the Workflow header.
- The money record's `gate` field, and the paused form of the Workflow result. Serves `paused at the gate`, `paused at the gate · needs a decision` and `reserves $40.80 of the $25.00 budget · nothing paid`.

**Change**, each different:

- The run's `budget` line reads the recorded design-doc budget and prints "set by the team · FIN-03", not "set by owner" from the budget input. Serves `console › budget     $25.00 per run · set by the team · FIN-03`.
- The run takes its estimate, reservation and budget from the plan's spend-table row for the applied pair, so the gate is the one the card described. Serves the `estimate` and `gate` lines.
- The margin is fixed at 20% in the money gates, not read from an input.
- `finops/prices.json` replaces the numbered price-list folder, and the price list and forecasts lose their mirror headers.
- The ledger keeps one money entry per run with the fields under Produces, not laps, billing and totals.

**Interfaces with neighbours:**

- From spec 0053: the design doc's per-run model-call budget in cents and the document's path, before the plan's first line. From spec 0054: the ids of the guardrail criteria the estimate and the budget satisfy (FIN-02, FIN-03).
- To spec 0050: the price list's models, the forecasts, the spend table, the budget and the margin, carried in the plan's hand-off. From spec 0050 and spec 0057: the applied model and effort.
- From spec 0058, through spec 0065's usage record: the input and output tokens of each paid stage call, summed for the run.
- To spec 0065: `observedCents` and `reservedCents` for `✓ completed implementation`. To spec 0063: `releasedCents` for `✓ ready` and `observedCents` for the phase-2 header.
- To spec 0057: the gate's outcome. A paused run has no stage, so no waiting line.

**Failure paths the simulation does not show:**

- A missing or unreadable price list or forecast, a design doc without a budget, and a chosen model or effort missing from the table each stop the command before the line that would need them. Nothing further prints, and the stop follows the single rule spec 0048 sets for every stop.
- Nothing is reserved and nothing is paid before the gate opens.
- A run that stops after the gate keeps its money record with `reservedCents` and whatever usage was recorded; its release stays unknown, never zero.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test**: `dft agent plan 0031-TICKET-book-a-demo-runtime`, then `dft agent run`, against the shared fixture project, comparing the printed lines and the produced card data with the simulation byte for byte. This spec adds:
  - the two phase-1 lines and the Plan PRICE field;
  - the spend table in the hand-off, all nine rows as listed under Produces, with `budgetCents` 2500 and `marginPercent` 20;
  - the default run's `estimate`, `budget` and `gate … opened` lines and the money record: reserved 1140, observed 980 and released 160 cents, with the fixture's stage driver reporting 3,700K input and 240K output tokens in total;
  - an over-budget variant, recording Opus at high effort as the applied pair: the four lines quoted under Shows, then nothing; no stage driver call; a money record with `gate: paused` and no observed spend.
- **A unit test of `domain/finops`** checks the formulas against the nine pairs of the table, with the observed figures for medium effort: Opus $19.60, Sonnet $9.80, Haiku $0.49. It also checks:
  - half-up rounding, where Haiku at medium effort is 47.5 cents and rounds to 48;
  - the ceiling of the reservation, where 26 cents reserves 32;
  - the gate's boundary, where a reservation equal to the budget opens.
- A good test asserts printed lines and recorded figures, never internal gate structures. Prior art: the catalog tests in `presentation/cli/session` and `domain/finops/arithmetic_test.go`.

## Out of Scope

- The design doc's extraction of the budget (spec 0053) and the guardrail's FIN criteria (spec 0054).
- The judge's suggestion and who pays for planning (spec 0051).
- How the applied pair is chosen and recorded (spec 0057), and the hand-off that carries it (spec 0050).
- Reading usage from the coding agent (spec 0058) and the usage record (spec 0065).
- The `completed`, `ready` and phase-2 header lines that quote these figures (specs 0065 and 0063).
- The console and its stub API, which only receive this data.

## Open Questions

- **Who keeps the price list current.** The simulation shows one list read on 2026-10-06. Whether `dft` ships updated lists, or the team may edit or extend `finops/prices.json`, is the owner's call. Until then `dft init` writes the three models and `dft` reads only them.
- **Where a ticket's expected tokens come from.** The simulation's data says they follow "from the spec's five components; no episode measured tokens", while the terminal shows only the result. A team-wide forecast per effort, a figure derived from the components, or episodes that measured tokens are each possible. Until the owner decides, `finops/forecasts.json` is the team's forecast per effort.
- **The fixed 20% margin.** The simulation never says who sets it. Whether a team may change it, and where, is open.
- **Raising the budget after a pause.** "Raise the budget" means editing the design doc. Whether a raised budget needs a new plan (revision 2) or is read again by the next `dft agent run`, the simulation does not show.
- **Observed spend above the reservation.** The simulation never shows observed above reserved. What RELEASED and the `ready` line say then, a negative release or something else, is open. No during-run gate is added.

## Further Notes

- `round` is half up, as in the simulation's `Math.round`, and all money is integer cents until it is formatted. The formatter is US dollars with two decimals in every terminal line. The card formats in the viewer's locale, which is the console's concern.
- The website's accessible status sentences ("The simulated run ended…", "The simulated run paused at the gate…") are the site's announcements, built from these same figures. They are not product data.
- The glossary's **Input record** and **Mirror** entries describe the machinery this spec removes for money. `domain-modeling` should update them when the tickets land.
