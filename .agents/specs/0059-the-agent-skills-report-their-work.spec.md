# The agent skills report their work in one structured line per stage

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E12.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 14 of 19, after spec 0058 and before spec 0060. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the hero simulation each stage of the run says what it did in one line of figures:

- grilling settled 3 open questions;
- to-spec revised 0011-SPEC to 3 acceptance criteria;
- to-tickets wrote 3 tickets with their blocking edges;
- implement worked test-first on 5 components and 14 files, wrote only the boundary and stopped at the reserved part;
- code-review found 2 findings and fixed them.

`implement · tdd` is one stage and one line. `code-review` runs after the dev ticket, under Review, without waiting for anyone.

The LangDrift CLI prints `<skill> <summary>`, where the summary is whatever the agent printed last, cut to 120 characters. It prints one line, and makes one paid call, per skill, so `implement` and `tdd` are separate. It runs `code-review` wherever the workflow puts it, which in the sample is inside Implement, before any dev ticket. Nothing in a stage line is a figure `dft` can check or keep, and nothing tells it that implement stopped at the boundary or at the reserved part.

## Solution

Each stage's skill ends by writing a small JSON report at the path the stage driver gives it, holding the figures its line shows. `dft` reads the report after the call ends and checks it against the stage's shape. It prints the stage's one line from the figures, in the simulation's exact wording, and keeps the figures as evidence. Implement runs as `implement · tdd` in one call with one line, and `code-review` runs under Review after the dev ticket closes.

## Responsibility

This element defines the skills the workflow runs and the report each one writes, checks it, and prints the stage's one line: `grilling`, `to-spec`, `to-tickets`, `implement` and `code-review`. LangDrift ships the skills: `dft init` installs them under `.langdrift/skills/`, one folder per skill (`grilling`, `to-spec`, `to-tickets`, `implement`, `tdd`, `code-review`), as part of the opinionated spec-driven workflow (spec 0057). They are taken from the skills this repository runs on its own work, from `.claude/skills/`, outside any `.langdrift/` folder: LangDrift ships the workflow it uses. They work on the client's own tickets and documents, which stay in the client's `.agents/` or `.claude/` folder (spec 0053).

It does not:

- start the coding agent or read its usage; spec 0058 does;
- decide the stages' order or the checkpoints before them; spec 0057 does;
- run the effort hook printed before `implement`; spec 0060 does;
- run the dev ticket between `implement` and `code-review`; spec 0061 does;
- keep the evidence file; spec 0065 does, with these figures in it.

## Produces

- **The skills, under `.langdrift/skills/`.** `dft init` installs them (spec 0048): `grilling` for Planning, `to-spec` for Spec, `to-tickets` for Tickets, `implement` and `tdd` for Implement, `code-review` for Review. Each skill ends by writing its stage report below. The stage driver (spec 0058) runs them with the chosen model and effort.
- **The stage report, written by the skill at `LANGDRIFT_REPORT`.** The path is `.langdrift/agent/state/runs/<run>/stages/<stage>.report.json`, where `<run>` is the run's id (spec 0065), private and ignored. It holds one JSON object per stage:

  | Stage | Fields | Values in the simulation |
  | --- | --- | --- |
  | `planning` | `openQuestionsSettled` (integer) | 3 |
  | `spec` | `spec` (the path of the spec the ticket links), `revised` (true), `acceptanceCriteria` (integer) | the 0011 spec, true, 3 |
  | `tickets` | `tickets` (integer), `blockingEdgesRecorded` (true) | 3, true |
  | `implement` | `method` (`tdd`), `components` (integer), `files` (integer), `boundaryOnly` (true), `stoppedAtReservedPart` (true) | tdd, 5, 14, true, true |
  | `review` | `findings` (integer), `fixed` (true) | 2, true |

  The shapes join the LangDrift Core contracts (spec 0049) as the stage-report schemas.
- **The figures, handed to spec 0065** for the stage's evidence, in place of the free-text summary.

No file under `.langdrift/` is versioned by this element.

## Shows

Phase 2, `dft agent run`, one line per stage, each printed when its call has ended and its report checked.

`grilling` follows `✓ confirmed  planning · by the owner` (spec 0057):

```
console › grilling   3 open questions settled
```

`to-spec` follows `grilling`, or the Spec pair when Spec is ticked:

```
console › to-spec    0011-SPEC revised · 3 acceptance criteria
```

`to-tickets` follows `to-spec`, or the Tickets pair when Tickets is ticked:

```
console › to-tickets 3 tickets · blocking edges recorded
```

`implement` follows `console › hook       effort recorded · implement unblocked` (spec 0060):

```
console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part
```

`code-review` follows `✓ closed     dev-ticket · written without an agent · no paid call` (spec 0061) and comes before `runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB` (spec 0062). It never has a waiting line before it, whether or not Review is ticked:

```
console › code-review 2 findings · fixed
```

The templates, with the label's column set by spec 0048:

- `grilling   {openQuestionsSettled} open questions settled`;
- `to-spec    {spec short id} revised · {acceptanceCriteria} acceptance criteria`, where the short id (`0011-SPEC`) follows spec 0053's rule for linked documents;
- `to-tickets {tickets} tickets · blocking edges recorded`;
- `implement  tdd · {components} components · {files} files · boundary only · stops at the reserved part`;
- `code-review {findings} findings · fixed`.

No card field. The card's stage states (spec 0057) turn on these lines.

## User Stories

1. As the owner, I want grilling to tell me how many open questions it settled, so that I know planning closed its unknowns.
2. As the owner, I want to-spec to name the spec it revised and how many acceptance criteria it now holds, so that I know what the implementation is held to.
3. As the owner, I want to-tickets to tell me how many tickets it wrote and that their blocking edges are recorded, so that I know the work is broken down in order.
4. As the owner, I want implement to say it worked test-first, on how many components and files, so that I can judge the size of the change.
5. As the owner, I want implement to say it wrote the boundary only and stopped at the reserved part, so that I know the part I write is still mine.
6. As the owner, I want `implement · tdd` to be one stage with one line, so that the session reads as the card names the stage.
7. As the owner, I want code-review to run after I wrote the dev ticket, so that it reviews my part too.
8. As the owner, I want code-review to run without waiting for me, so that my review starts on a pull request the agent already reviewed and fixed.
9. As the owner, I want code-review to tell me how many findings it had and that they were fixed, so that I know what changed after my part.
10. As a skill author, I want one small JSON report per stage at a path `dft` gives me, so that my skill reports figures, not prose.
11. As a skill author, I want the report's shape published as a contract, so that I can validate my skill's output before a run.
12. As a team member, I want every stage line built from checked figures, so that a line never repeats whatever the agent happened to print last.
13. As a team member, I want a stage whose report is missing or does not match its shape to print no line, so that the session never claims work it cannot show.
14. As an auditor, I want each stage's figures kept in the run's evidence, so that the lines can be traced afterwards.
15. As the console, I want the stage lines to come from `dft` as data, so that I only show them.
16. As a FinOps reader, I want one paid call behind each stage line, so that the usage record maps to the lines I read.

## Implementation Decisions

**Remove**, each beyond what the simulation shows:

- The stage line `<skill> <summary>` built from the agent's last output line, cut to 120 characters. Serves each of the five lines, built from the report instead.
- One line and one paid call per skill, which split `implement` from `tdd`. Serves `console › implement  tdd · …`.
- `code-review` as a skill of Implement, run before the dev ticket. Serves `console › code-review 2 findings · fixed` after `✓ closed     dev-ticket …`.
- The free-text `summary` kept in a stage's evidence.

**Build**, each shown and missing:

- The five stage-report shapes in the LangDrift Core contracts, and their check in `dft`. Serves every stage line.
- The line template for each stage, from the report's figures. Serves the five lines quoted under Shows.
- The `revised` flag and the spec's short id in the to-spec report. Serves `0011-SPEC revised`.
- The `boundaryOnly` and `stoppedAtReservedPart` flags in the implement report. Serves `boundary only · stops at the reserved part`.
- The `blockingEdgesRecorded` and `fixed` flags. Serve `blocking edges recorded` and `fixed`.
- `code-review` as Review's skill, called after the dev ticket closes, with no checkpoint before it (with spec 0057's workflow). Serves the line's position.

**Change**, each different:

- A stage line is printed only after its call ended with a zero exit status and its report matched the shape, not on any exit.
- A stage's evidence carries its report's figures, not a summary.

**Interfaces with neighbours:**

- From spec 0058: the report path given to the call, and the call's exit status.
- From spec 0053: the short-id rule for a linked document (`0011-SPEC`), and the spec the ticket links.
- To spec 0065: the checked figures, per stage.
- To spec 0057: the moment each stage's line printed, which ends the stage on the card.
- With spec 0061: the implement report's `stoppedAtReservedPart` is what hands the reserved part to the person.

**Failure paths the simulation does not show:**

- A report that is missing, is not JSON, lacks a field, or holds `false` where the line prints fixed text fails the stage. So does a `method` other than `tdd` or a `spec` other than the one the ticket links.
- A failed stage prints no line and stops the run under spec 0048's single rule.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** The fixture's fake coding agent (spec 0058) writes, for each stage, the report with the simulation's figures. This spec adds:
  - the five lines, byte for byte, at their places in the default run: `grilling` after the Planning pair, `implement` after `hook`, `code-review` after `closed     dev-ticket`;
  - that no `waiting` line precedes `code-review` when Review is ticked;
  - the five stages' figures in the run's evidence (spec 0065).
- **A unit test per stage** formats a report into its line and refuses each malformed report: a missing field, a wrong type, a `false` flag, an unlinked spec. Prior art: the session catalog tests in `presentation/cli/session`.

## Out of Scope

- The coding agent's invocation and usage (spec 0058).
- The order of the stages and their checkpoints (spec 0057), the effort hook (spec 0060) and the dev ticket (spec 0061).
- The evidence file (spec 0065).
- What the skills do inside the client's project, beyond the report they write.
- The console and its stub API.

## Open Questions

- **Outcomes the simulation does not show.** The simulation never shows findings left unfixed, an implement that went past the boundary, a spec written rather than revised, or tickets without blocking edges. What the line says then is open. Until the owner decides, such a report fails the stage.
- **Checking the reported figures.** Under "evidence over assertion", should `dft` check figures it can measure itself rather than trust the report? It could count the files from the branch's diff, or the acceptance criteria in the spec. The simulation does not say.

## Further Notes

- Counts print as plain integers in the shown wording. A count of 1 takes the singular noun ("1 ticket", "1 finding", "1 acceptance criterion"), which is English grammar, not a product rule.
- The report path sits in the private, ignored state folder because only the figures need keeping, and spec 0065 keeps them.
