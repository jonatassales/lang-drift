# The harness states the guardrails and the envelope

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E7.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 6 of 19, after spec 0053 and before spec 0056. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation, once the ticket's documents are read, the plan states the guardrail the ticket links and the criteria it reads, and states the envelope as the plan's bound, "set by the team in the design doc". The Analysis card shows the guardrail as a review rule with no automatic enforcement, with its three criteria in short form. The LangDrift CLI disagrees on where the bound comes from. It takes the envelope from a recorded input or from `--cpu-millis`/`--memory-mib` flags, mirrors it in `harness/envelope.json` with a per-execution ceiling, and prints "set by X in the CLI · record Y". It gives the guardrail its own `tool call: read` line. The `reads` line lists every `FIN-\d{2}` id found in the file, in file order. No record says how a guardrail is enforced, so the card's "review rule, no automatic enforcement" has no data behind it. A team that wrote its envelope in the design doc sees the plan bound by a number typed somewhere else.

## Solution

The harness is the part of the client's `.langdrift/` folder that holds the team's guardrails. When the plan reaches the guardrail the ticket links, it reads that document from `.langdrift/harness/guardrails/`, prints the `reads` line with the document's path and the criteria the plan step answers to, and records the guardrail as a review rule that `dft` never enforces. After the judge has stated the criteria (spec 0051), the plan prints the `bounded` line. The bound is exactly the envelope the design doc sets (read by spec 0053), with "set by the team in the design doc" as its provenance. No CLI input, recorded input, mirror or ceiling takes part.

## Responsibility

This element:
- reads the guardrail document the ticket links (the link is found by spec 0053): its criteria, in document order;
- prints the `reads` line;
- records every guardrail as a review rule with no automatic enforcement;
- states the envelope extracted from the design doc as the plan's bound, with its provenance, and prints the `bounded` line;
- owns the card's guardrail document row and the FIN line under it;
- defines the shape of a guardrail document that the plan reads.

It does not:
- find or name the ticket's links, or extract the envelope from the design doc (spec 0053);
- state what each criterion requires: the short form of each criterion is the judge's (spec 0051), and this element only shows it on the card;
- size anything against the bound or show BOUNDED and HEADROOM (spec 0055);
- state the budget or enforce money gates, even those that cite FIN-02 and FIN-03 (spec 0056);
- scaffold the guardrail document into a new project (spec 0048 does, in this shape).

## Produces

No new file. The guardrail documents stay where they are, under `.langdrift/harness/guardrails/<number>-<name>.gr.md`. This element reads them and hands its values to its neighbours:

- **Guardrail** (into the plan's card data, spec 0050):
  - `path`, relative to `.langdrift/` (`harness/guardrails/0001-finops.gr.md`), feeding the `reads` line;
  - `file` (`0001-finops.gr.md`), feeding the card row;
  - `criteria`: ordered ids (`FIN-02`, `FIN-03`, `FIN-04`), each with its rule text and the steps it is checked at, feeding the card row;
  - `planCriteria`: the criteria checked at the plan step (`FIN-02`, `FIN-04`), feeding the `reads` line;
  - `enforcement` (`review`), feeding "review rule, no automatic enforcement".
- **Criteria text** handed to spec 0051, which returns one short statement per criterion. This element shows those statements in the FIN line.
- **Bound** (into the `runtime.md` draft through spec 0055, and into the card data): `cpuMillis` (integer, millicores, `2000`), `memoryMiB` (integer, MiB, `3072`), `setBy` (`the team`), `source` (the design doc's path). Feeds the `bounded` line, and through spec 0055 the BOUNDED field, the `✓ size … of 2000m … of 3072 MiB` figures and HEADROOM. Through spec 0066 it also feeds the components' share bars.

## Shows

Phase 1, `dft agent plan`. The first line follows `✓ read       4 documents · nothing outside the ticket's links` (spec 0053):

```
reads      harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04
```

The judge's `console › judging · FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load` follows (spec 0051), then:

```
bounded    cpuMillis 2000 · memoryMiB 3072 · set by the team in the design doc
```

The line after it is `budget     $25.00 per run for model calls · set by the team in the design doc` (spec 0056).

Card fields, in the Analysis card's DOCUMENTS list, after the PRD, design doc and spec rows:

- **Guardrail document row**, set when the `reads` line prints: `0001-finops.gr.md ✓ Read — FIN-02 · FIN-03 · FIN-04 · review rule, no automatic enforcement`.
- **Its FIN line**, set when the judge's FIN line prints: `FIN-02: estimate before a paid call · FIN-03: a budget before a paid run · FIN-04: size cpu and memory to the load`.

BOUNDED on the Plan card, though it shows this element's bound once the `bounded` line prints, belongs to spec 0055.

## User Stories

1. As a team lead, I want the plan to read the guardrail the ticket links from the harness, so that our FinOps rules are part of every plan.
2. As a developer, I want the `reads` line to show the guardrail's path under `.langdrift/`, so that I can open the exact file the plan read.
3. As a team lead, I want the `reads` line to list the criteria the plan step answers to, so that I see which rules this step must respect.
4. As a reviewer, I want the card to list every criterion of the guardrail, so that nothing the ticket links is hidden.
5. As a reviewer, I want the card to say the guardrail is a review rule with no automatic enforcement, so that nobody mistakes a written rule for a technical gate.
6. As a reviewer, I want each criterion shown in a few words under the guardrail row, so that I can read what FIN-02, FIN-03 and FIN-04 ask without opening the file.
7. As a team lead, I want the plan's bound to be the envelope our design doc sets, so that the limit we wrote is the limit the plan respects.
8. As a team lead, I want the `bounded` line to say "set by the team in the design doc", so that provenance names our document and not a flag someone typed.
9. As a developer, I want `cpuMillis` and `memoryMiB` printed as the integers the design doc holds, so that the terminal matches the document character for character.
10. As a platform owner, I want no CLI flag, recorded input or mirror able to replace the design doc's envelope, so that there is one source of the bound.
11. As a platform owner, I want no per-execution ceiling beside the envelope, so that there is one bound, not two.
12. As a developer, I want a ticket whose guardrail cannot be read, or holds no criteria, to stop the plan with nothing written, so that a plan never claims rules it did not read.
13. As a developer, I want a design doc without an envelope to stop the plan with the field named, so that the plan never sizes against an invented bound.
14. As a maintainer of `dft init`, I want one documented guardrail shape, so that the scaffold, the fixture and every team's guardrails read the same way.
15. As an implementer of sizing, I want the bound handed over with its provenance, so that BOUNDED, HEADROOM and the size line use one value.
16. As an auditor, I want the bound recorded in the `runtime.md` draft with its source, so that the plan's evidence says where the limit came from.
17. As a team member, I want the guardrail's row to show ✓ Read like the other documents, so that the four linked documents read as one list.

## Implementation Decisions

**Remove**
- The guardrail's own `tool call: read` line. The `reads` line is the guardrail's only line.
- Listing every `FIN-\d{2}` match in the file in the `reads` line. Serves `reads … · FIN-02 · FIN-04`: the line lists the criteria checked at the plan step.
- The envelope taken from a recorded input, from the `harness/envelope.json` mirror or from `--cpu-millis`/`--memory-mib`. Removed with them: the envelope's input record (record id, set by, via, set at), the default envelope, the per-execution ceiling and its within-ceiling check, and the `notRuntimeFields` list in `envelope.json` (the limits line is spec 0055's). Serves `bounded … · set by the team in the design doc`.
- The text "set by X in the CLI · record Y".
- The judge's figure comparisons against the input, and the "discrepancy … the input holds" lines. Spec 0051 records their removal from the judge. Here, nothing holds against the design doc's envelope.
- The plan record's input-provenance fields for the envelope (record, via, set at) and its ceiling.

**Build**
- A guardrail document's shape: a Markdown document in the harness's guardrails folder (see Produces), holding one criteria table whose columns are `ID`, the rule, and `Checked at` (`plan`, `run`, or `plan, run`). The fixture's `0001-finops.gr.md` has three rows:
  - FIN-02, `plan, run`;
  - FIN-03, `run`;
  - FIN-04, `plan`.

  Serves the card row's three criteria and the `reads` line's two.
- The `reads` line: the guardrail's path relative to `.langdrift/`, then the ids whose `Checked at` includes `plan`, in table order.
- The enforcement record: every guardrail document is recorded as `review`, because `dft` enforces no guardrail criterion. The money gates of spec 0056 are their own mechanism and cite FIN ids only as references. Serves "review rule, no automatic enforcement".
- The FIN line on the card: each criterion as `<id>: <short statement>`, joined by ` · `, in table order, with the statements spec 0051 returns.
- The bound: the design doc's envelope as spec 0053 extracts it, with `setBy` "the team" and the design doc as source, recorded in the `runtime.md` draft (spec 0055) and the card data (spec 0050).
- Safe stops: no guardrail link, an unreadable guardrail, a guardrail without a criteria table or with a row lacking `Checked at`, and a design doc without the envelope. The plan stops before the line this element would print, writes no draft and exits non-zero. The reason goes through the CLI's error channel (spec 0048). No line the simulation lacks is printed.

**Change**
- The envelope's source moves from CLI input to the team's design doc. The Core contract's statement that "a fact collected from a source never sets one" no longer holds for the envelope. Spec 0049 carries the contract side.
- `harness/envelope.json` is no longer read or written. The harness keeps the guardrails, and the workflow file of spec 0057.

**Interfaces with neighbours**
- From spec 0053: the guardrail link and the design doc's envelope with its path.
- To spec 0051: the criteria with their rule text. From spec 0051: one short statement per criterion.
- To spec 0055: the bound, recorded in the draft and used for the size line, BOUNDED and HEADROOM.
- To spec 0066: the bound's `cpuMillis`, the denominator of each component's share bar.
- To spec 0050: the guardrail row, the FIN line and the bound in the card data.

## Testing Decisions

- The seam is the golden test of spec 0048 over the shared fixture. This element adds the `reads` and `bounded` lines, the guardrail row and the FIN line to the comparison, byte for byte. The fixture's `.langdrift/harness/guardrails/0001-finops.gr.md` holds the three-row criteria table above. Its design doc holds the envelope that spec 0053 extracts. No `envelope.json` exists in the fixture.
- Golden cases for the safe stops: a ticket without a guardrail link, a guardrail without the table, and a design doc without the envelope. Each prints up to the stop, writes no draft and exits non-zero.
- One unit test for the `reads` rule: which ids a table with mixed `Checked at` values lists, and in what order.
- Tests observe the printed lines, the card data and the draft's bound only.

## Out of Scope

- Reading the ticket and the documents, and extracting the envelope and the budget: spec 0053.
- The judge's short statements of the criteria and the `judging · FIN-…` line: spec 0051.
- BOUNDED, HEADROOM, the size line and the limits line: spec 0055.
- The budget line, the money gates and their FIN references: spec 0056.
- Scaffolding the guardrail in a new project: spec 0048.
- The console and its stub API.

## Open Questions

- **`reads … FIN-02 · FIN-04` omits FIN-03.** The `reads` line lists FIN-02 and FIN-04. The next line, the judge's, and the card both list FIN-02, FIN-03 and FIN-04. This spec follows the simulation literally. Its provisional mechanism is a `Checked at` column that scopes FIN-03 to the run. Two decisions are open:
  - Is the omission meant? If it is a slip, the `reads` line lists every criterion and the column is dropped.
  - If it is meant, is a `Checked at` column the way a guardrail says which step a criterion governs?

## Further Notes

- The sample guardrail in `lib/cli` testdata carries FIN-01 to FIN-08 without a `Checked at` column. The simulation's guardrail carries FIN-02, FIN-03 and FIN-04. The scaffold `dft init` writes follows the simulation (spec 0048).
- "The team" is the design doc's author as a group, not a person or an account. No identity is read or recorded for it.
