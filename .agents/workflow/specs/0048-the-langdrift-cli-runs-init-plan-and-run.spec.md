# The LangDrift CLI runs init, plan and run

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E1.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 2 of 19, after spec 0049 and before spec 0050. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The session on the website's home is the definition of LangDrift's main feature: a person types `dft agent plan 0031-TICKET-book-a-demo-runtime`, then `dft agent run`, and nothing else. The LangDrift CLI that ships today answers to some thirty commands and flag modes (`doctor`, `status`, `reindex`, `confirm`, `cluster …`, `runtime …`, `harness …`, `finops …`, `gates report`, `connect`, `supervise`, `--json` and more), takes a file path where the simulation types a bare ticket id, requires a runtime argument where the simulation types `dft agent run` alone, prints the typed command a second time, names every actor line `agent ›` where the simulation prints `console ›`, and opens with a scope comment that says the planner's calls are paid. A person who watched the session and then types the same commands gets a different tool. `dft init` asks for a monthly budget, a per-run cap, a reserve margin, gates, an artifacts folder and a CI workflow, none of which the session needs, because the envelope and the budget come from the design doc.

There is also no shared definition of how a line is laid out, so every element that prints a line re-invents its columns, and no shared fixture and golden test against which every element can prove that it prints exactly the simulation.

## Solution

`dft` has three commands: `dft init`, `dft agent plan <ticket id>` and `dft agent run`. Each prints exactly the simulation's lines for its phase, in the simulation's order, with the simulation's layout, and nothing else. The typed command is never echoed: the person's shell shows it once. Every actor line starts with `console › `. The terminal's title reads `dft — project`. `dft init` creates the client's `.langdrift/` folder with only what plan and run need.

This spec also fixes three things every other element spec relies on:

- **The line format**: the kinds, the prefixes and the column layout every printed line follows.
- **The stop rule**: how any failure the simulation does not show ends a command without printing a line the simulation does not have.
- **The fixture and the golden test**: a project that reproduces the simulation's ticket, documents, guardrail, six episodes and five components, against which one end-to-end test compares `dft agent plan 0031-TICKET-book-a-demo-runtime` and `dft agent run` with the simulation byte for byte. Every other element spec adds its assertions to this test.

## Responsibility

This element is the command surface of the LangDrift CLI and the frame of its session:

- which commands and flags exist, and how their arguments resolve to a ticket and to a recorded intent;
- finding the project root and creating `.langdrift/` with `dft init`;
- the comment lines that open and close the plan's three acts;
- the terminal title;
- the line format, the line kinds and their terminal tones;
- the stop rule and the exit status;
- the shared fixture and the golden end-to-end test.

It does not:

- read the ticket or its documents (spec 0053), judge (spec 0051), search episodes (spec 0052), state guardrails and the envelope (spec 0054), size `runtime.md` (spec 0055), inspect components (spec 0066), or price, gate and release (spec 0056);
- hand off the plan or read the authorization (spec 0050), even though it decides which recorded intent `dft agent run` picks;
- wait for or confirm a human stage (spec 0057), drive the coding agent (spec 0058) or report skills (spec 0059);
- bring up the runtime (spec 0062), deliver (spec 0063), follow checks (spec 0064) or keep the usage record and the episode (spec 0065).

## Produces

`.langdrift/` is LangDrift's own folder in the client's project. It holds the opinionated spec-driven development workflow that `dft init` sets up, the skills that workflow runs, and the FinOps harness, prices and records that plan and run use. The client's own artifacts stay outside it: the tickets that take part in the workflow, with the PRD, design doc and spec they link, live in the client's `.agents/` or `.claude/` folder (spec 0053).

`dft init` writes, in the working directory:

| Path | Content | Feeds |
| --- | --- | --- |
| `.langdrift/` | The folder itself. Its presence marks the project root for plan and run. | Every command's project discovery |
| `.langdrift/skills/` | The skills the workflow runs, one folder per skill as spec 0059 defines them: `grilling`, `to-spec`, `to-tickets`, `implement`, `tdd` and `code-review`, taken from the skills this repository runs from `.claude/skills/`. | Every stage line of spec 0059 |
| `.langdrift/finops/prices.json` | The three models' list prices, read on 2026-10-06, with the content spec 0056 defines. | The Workflow MODEL price line, every figure of spec 0056 |
| `.langdrift/finops/forecasts.json` | The expected tokens per effort, with the content spec 0056 defines. | The estimate line and ESTIMATED SPEND |
| `.langdrift/harness/workflow.json` | The seven stages with their skills and persons, and the default choice `defaults` (model `claude-sonnet-5-5`, effort `medium`, human Planning and Review), with the content spec 0057 defines. | HUMAN IN THE LOOP, the default ticks, the default MODEL and EFFORT, every `waiting` and `confirmed` pair |
| `.langdrift/harness/guardrails/0001-finops.gr.md` | The FinOps guardrail in the shape spec 0054 reads, with the simulation's three criteria FIN-02, FIN-03 and FIN-04, as a review rule with no automatic enforcement. From then on it is the team's document. | `reads harness/guardrails/0001-finops.gr.md …`, the guardrail row and its FIN line |
| `.langdrift/templates/` | Four empty document templates, `ticket.md`, `prd.md`, `dd.md` and `spec.md`, carrying only the field lines and sections spec 0053 reads (`**Status:**`, `**Use case:**`, `**cpuMillis:**`, `**memoryMiB:**`, `**Budget per run:**`, `## Components`, `## Acceptance criteria`). Plan and run never read them; see Open Questions. | The team's documents, which feed spec 0053's lines |

Rules:

- `dft init` writes exactly the files an element spec declares as written by `dft init`, and nothing else. The table above lists the files declared so far; a spec that needs another init file adds it here, with its content and the shown line or field it feeds.
- `harness/driver.json` (spec 0058) is not written: until the owner decides which coding agent `dft` ships a configuration for, the team writes it. A run without it, or with one that cannot pass the model and effort, stops before it prints anything (spec 0058).
- An existing file is kept byte for byte; `dft init` never replaces one, and running it twice changes nothing.
- It writes no project identity (`config.json`), no input records, no mirrors, no budget, no gates, no artifacts folder record, no execution rules and no CI workflow.

`dft agent plan` and `dft agent run` write nothing of this element's own: every record they leave is declared by the element that produces it.

## Shows

### The terminal title

`dft — project`: the word `dft`, an em dash between single spaces, and the name of the project folder (the folder that holds `.langdrift/`). `dft agent plan` and `dft agent run` set it when they start, only when standard output is a terminal.

### Phase 1, `dft agent plan`

The prompt and the command are the person's shell; `dft` prints neither:

```
~/project %
dft agent plan 0031-TICKET-book-a-demo-runtime
```

The comment lines this element prints, all of kind comment, each at its moment:

1. `# plan only · nothing is built, started or paid for in this step`: the first line `dft` prints, before `console › tool call: read 0031-TICKET-book-a-demo-runtime` (spec 0053).
2. `# analysis · reading the linked documents, then the project's episodes`: after `console › ready-for-agent · links 0002-PRD · 0005-DD · 0011-SPEC · guardrail 0001-finops` (spec 0053), before `console › tool call: read 0002-PRD · use case: book-a-demo request` (spec 0053).
3. `# plan · resolving the components, then sizing them inside the envelope`: after `✓ suggests   claude-sonnet-5-5 · effort medium` (spec 0051), before `console › resolving 5 components named in the spec · paths, images, ports` (spec 0066).
4. `# next: once the team authorizes, dft agent run applies the recorded intent`: the last line, after `✓ plan       runtime.md draft, revision 1 · proposed to the console` (spec 0050).

### Phase 2, `dft agent run`

The prompt and the command are the person's shell; `dft` prints neither:

```
~/project %
dft agent run
```

The first line `dft` prints is `# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team` (spec 0050). This element prints no other line in phase 2. Every variant of phase 2 (over budget, review unchecked, stages checked) starts the same way.

### The Workflow card's controls

`Play workflow` and `Play again` are the website's playback controls, with the states `Starting…`, `Playing…` and `available once the plan is drafted`. In the product, `Play workflow` is the moment a person types `dft agent run`, and "available once the plan is drafted" is the rule that `dft agent run` needs a recorded intent. `dft` produces no field for them. A second run of the same intent, the product's counterpart of `Play again`, is an open question.

## User Stories

1. As a developer who watched the session, I want `dft` to have exactly `init`, `agent plan` and `agent run`, so that the tool I install is the tool I saw.
2. As a developer, I want to type `dft agent plan 0031-TICKET-book-a-demo-runtime` with the ticket's bare id, so that I do not need to know where the ticket file lives.
3. As a developer, I want `dft agent run` to take no argument, so that I run what the team authorized without naming it again.
4. As a developer, I want `dft` never to print the command I typed, so that the terminal shows it once, as in the session.
5. As a developer, I want the plan to open with `# plan only · nothing is built, started or paid for in this step`, so that I know the plan spends nothing.
6. As a developer, I want the plan's three acts opened by their comments, so that I can follow the ticket, the analysis and the plan as the session does.
7. As a developer, I want the plan to end with `# next: once the team authorizes, dft agent run applies the recorded intent`, so that I know what happens next and who decides it.
8. As a developer, I want every actor line to start with `console › `, so that the terminal reads like the session.
9. As a developer, I want results, statuses and stage lines aligned in the session's columns, so that values line up as I read down the terminal.
10. As a developer, I want the five component rows aligned in three columns, so that paths, stacks and shares read as a table.
11. As a developer, I want success lines green, warnings yellow and comments dimmed in my terminal, so that the session's meaning carries its tone.
12. As a developer piping `dft` into a file, I want plain text without escape sequences, so that the file holds exactly the session's lines.
13. As a developer, I want the terminal's title to read `dft — project`, so that the window names the tool and the project.
14. As a developer, I want `dft init` to create `.langdrift/` with the price list, the forecasts, the workflow and the FinOps guardrail, without asking me anything, so that the envelope and the budget come from the design doc, where the team set them.
15. As a developer, I want `dft init` to keep every file I already have, so that running it twice is harmless.
16. As a developer, I want `dft agent plan` and `dft agent run` to find the project from any folder inside it, so that I can run them where I work.
17. As a developer, I want a plan for a ticket id that matches no file, or several, to stop without a made-up line, so that the terminal never shows something the session would not.
18. As a developer, I want `dft agent run` to stop when there is no authorized intent, so that nothing runs before the team authorizes it.
19. As a developer, I want any stop to say its cause on standard error and exit non-zero, so that a script and a person both see that the command did not finish.
20. As a developer, I want the lines already printed to stay as they were when a command stops, so that the session reads truthfully up to the point it reached.
21. As a team lead, I want the over-budget pause to end the run as the session shows it, with `# needs a decision …`, so that a paused run is a normal ending, not a crash.
22. As the owner, I want every command and flag the session does not show removed, so that the product does exactly what its definition shows.
23. As a maintainer, I want one golden test that plays the plan and the run against a fixture of the simulation's project, so that any drift from the session fails a test.
24. As a maintainer of an element, I want to add my element's assertions and fixture content to that one test, so that every element is proven against the same session.
25. As a maintainer, I want the golden test to cover the over-budget, review-unchecked and stages-checked variants, so that every variant the simulation holds is proven.
26. As a maintainer, I want the golden test to compare the hand-off records with the four cards' fields, so that the data, not only the text, matches the simulation.

## Implementation Decisions

### Remove

Every command and flag mode beyond the three, under the CLI's command presentation and its catalog:

- The commands `doctor`, `agent status`, `agent reindex`, `agent confirm`, `gates report`, `version` and `help`, and the generated `completion` command. Serve no shown line.
- The `runtime` group and every runtime command the catalog still holds: `up`, `down`, `inspect`, `logs`, `forward`, `hook`, `init`, `validate`, and `agent migrate`, with the retirement message of `dft runtime`. Serve no shown line; the runtime's lifecycle inside a run is spec 0062's.
- The `cluster` group: `create`, `delete`, `inspect`, `load-image`. Serve no shown line.
- The connected-agent commands `agent connect`, `agent project`, `agent supervise`, and any `stop` command. Serve no shown line; the hand-off is spec 0050's files.
- The `harness` input commands (`envelope set`, `workflow set`, `model set`, `driver set`) and the `finops` input commands (`budget set`, `gates set`, `prices set`, `forecasts set`), with their `--by`, `--status`, `--owner` and `--file` flags. Serve no shown line: the envelope and the budget come from the design doc (specs 0053, 0054), the choices from the recorded intent (spec 0050).
- The global flags `--project`, `--json`, `--timeout`, `--version` and `-V`. The project is found from the working directory; `--json` and its NDJSON events go because the hand-off is files (spec 0050 decides it); a finite deadline is no shown behavior.
- The `agent plan` flags `--runtime`, `--name`, `--component`, `--cpu-millis`, `--memory-mib`, `--max-depth`, `--max-files`, `--max-bytes` and `--no-judge`. Serve no shown line: the runtime id is the ticket's slug (spec 0055), the components are named in the spec (spec 0066), the envelope is the design doc's (spec 0054), the reading stays inside the ticket's links (spec 0053), and the judge always judges (spec 0051).
- The `agent run` flags `--model`, `--effort`, `--acknowledge-overspend`, `--confirm`, `--checks-wait` and `--by`. Serve no shown line: the choices are the recorded intent's (spec 0050), confirmations are answered inside the run (spec 0057), the wait for checks is spec 0064's.
- The `init` flags and questions `--ask`, `--budget`, `--per-run`, `--reserve-margin`, `--default-gates`, `--artifacts`, `--ci-gates` and `--by`, the month's budget, the per-run cap, the reserve margin, the default gates, the artifacts folder and the CI Gates workflow (which runs the removed `gates report`). Serve no shown line: the budget is the design doc's and the margin is fixed (spec 0056).
- What `dft init` writes today beyond the table under Produces: the project identity in `config.json`, the input records and their lock under the private state, the mirrors, the execution rules `microruntimes/AGENTS.md`, and the CI Gates workflow. Serve no shown line. A spec that needs one of them declares it as an init file.
- The private state that only the removed commands use: connections, credentials and migrations. Serve no shown line.
- The echo of the typed command as the plan's and the run's first printed line. Serves the command, which the shell shows once.
- The project-relative file path as the plan's argument. Serves the bare id.

### Build

- **`dft agent plan <ticket id>`.** Exactly one argument, the ticket's bare id, as in `0031-TICKET-book-a-demo-runtime`; no argument, or more than one, is a stop. The id is handed to spec 0053, which resolves it to the one `<id>.ticket.md` of the project and reads it; it is also the id the run's header names (spec 0050). Serves the typed command and `console › tool call: read 0031-TICKET-book-a-demo-runtime`.
- **`dft agent run` with no argument.** It applies the one recorded intent (spec 0050) whose status is authorized and that has no completed run. None is a stop; more than one is a stop. Serves the typed command and `# run · 0031-TICKET-book-a-demo-runtime · …`.
- **The project root.** The nearest folder, from the working directory up, that holds `.langdrift/`. None is a stop for plan and run. `dft init` uses the working directory.
- **The terminal title**, `dft — <project folder name>`, set with the standard terminal title sequence when standard output is a terminal, by plan and run. Serves `dft — project`.
- **The four comment lines of phase 1**, quoted under Shows, each emitted at its moment.
- **The line kinds and their tones.** Every printed line carries the simulation's kind: comment, output, success or warning (command is the person's typing and is never printed). When standard output is a terminal, a comment is dimmed, a success line is green, a warning is yellow, an output line keeps the terminal's default color. Otherwise the bytes are the plain text, with no escape sequence. Serves the session's tones.
- **Lines print when their work is done.** No line is delayed to imitate the simulation's pacing.
- **The stop rule**, below.
- **The fixture and the golden test**, under Testing Decisions.

### Change

- The actor prefix `agent › ` becomes `console › ` on every actor line (the CLI's session format). Serves every `console › …` line of both phases.
- The plan's scope comment changes from "nothing is built or started · the planner's calls are paid and gated" to `# plan only · nothing is built, started or paid for in this step`.
- The run's argument `RUNTIME` is dropped; the run picks the recorded intent as Build says.
- `dft init` changes from an interactive setup of inputs to the silent creation of `.langdrift/` and its declared files. On success it prints nothing and exits 0, because the simulation shows no `dft init` line.
- The plan's last result line, `✓ plan … · proposed · authorize with dft agent run`, is spec 0050's to change.

### The line format every element uses

Every line is one UTF-8 line of English, ended by a newline, with no trailing space. Columns count characters, and `›`, `·`, `✓` and `—` count one each. Items inside a value are separated by ` · `. Numbers print without thousands separators (`1450m`, `2304 MiB`); money and token figures follow spec 0056 (`$25.00`, `3.5M`, `250K`). There are six shapes:

| Shape | Layout | Lines that take it |
| --- | --- | --- |
| Comment | `# ` + text | `# plan only …`, `# analysis …`, `# plan · resolving …`, `# next: …`, `# run · …`, `# needs a decision …` |
| Narrative actor line | `console › ` + text | `tool call: read …`, `ready-for-agent · links …`, `judging · …`, `RAG over project episodes …`, `embedding …`, `vector search …`, `6 candidates …`, `dropped …`, `3 episodes cited …`, `reasoning · …`, `weighing effort · …`, `weighing model · …`, `resolving …`, `tool call: inspect …`, `cross-checking · …`, `tool call: validate …` |
| Labeled actor line | `console › ` + label right-padded with spaces to 10 columns + one space + value | `intent`, `estimate`, `budget`, `gate`, `waiting`, `grilling`, `to-spec`, `to-tickets`, `hook`, `implement`, `code-review` |
| Result line | `✓ ` + label right-padded to 10 columns + one space + value | `read`, `suggests`, `size`, `plan`, `confirmed`, `closed`, `completed`, `pushed`, `pr #<number>`, `checks`, `Safe Seat`, `ready` |
| Status line | label right-padded to 10 columns + one space + value | `reads`, `bounded`, `budget`, `limits`, `estimate`, `runtime`, `checks` |
| Component row | `✓ ` + path right-padded to 19 columns + three spaces + `<stack> · <role>` right-padded to 15 columns + one space + `<cpu>m · <memory> MiB` | the five component rows |

A label of 10 columns or more is followed by exactly one space. With this rule:

- `console › to-tickets 3 tickets · blocking edges recorded` and `console › code-review 2 findings · fixed` keep one space;
- `✓ read       4 documents · nothing outside the ticket's links` and `✓ pr #19     feat(demo): run book a demo as a microruntime` reach the value at the same column;
- `reads      harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04` and `runtime    stopped · 5 Deployments removed · nothing left` align with each other;
- `✓ services/signup-api   Go · api        250m · 256 MiB` and `✓ apps/signup           Next.js · web   500m · 1024 MiB` align in three columns.

Each line's kind is the simulation's, as the website's session gives it. Of note: `limits     …` is a warning; `console › gate       reserved … · opened` is a success line although it has no mark; `console › gate       reservation … · paused` is a warning; `✓ …` lines are success; `# …` lines are comments; every other line is output.

The CLI keeps one catalog of these lines keyed by the website's line ids, as it does today; each element spec defines the words of its own lines in it.

### The stop rule every element uses

This is the single rule every element spec refers to. A failure, refusal or interruption that the simulation does not show is a stop:

- The command prints no further session line, and no warning or error line of its own. The lines already printed stay.
- A plan that stops writes no draft and no recorded intent; its card data stays as far as it got (specs 0055, 0050).
- If the run started the runtime, it is taken down without a line (spec 0062).
- The records reflect what happened, never more: a paid call is never left out of the usage record, a stage keeps the state it reached in the run's record (spec 0065), and money settles as spec 0056 says.
- Nothing is retried and nothing resumes inside the command.
- One plain sentence naming the cause goes to standard error, the CLI's error channel. It names no secret. Standard error is not the session; standard output stays a prefix of the simulation.
- The command exits with status 1.

A command that ends on its last simulation line exits 0. The over-budget pause, which ends on `# needs a decision: raise the budget or lower the effort · nothing was paid for`, is a shown ending and exits 0. An interrupt (Ctrl-C) is a stop.

### Interfaces with the neighbours

- To spec 0053: the bare ticket id and the project root; spec 0053 resolves the file.
- From spec 0050: the list of recorded intents with their status and whether each has a completed run; to spec 0050: the one intent the run applies.
- To every element: the session sink, which takes a line id, its kind and its data and prints it in the format above; and the stop rule.
- To spec 0057: the run's standard input, on which the human checkpoints are answered.

## Testing Decisions

### The seam

One end-to-end golden test in the LangDrift CLI's top-level test package, the highest seam: it runs `dft` in process, as the existing plan and run tests do with their injected dependencies, against a fresh copy of the fixture project. It runs `dft init` on a copy, then `dft agent plan 0031-TICKET-book-a-demo-runtime`, authorizes the recorded intent as spec 0050 says, then `dft agent run`. A good assertion compares what a person sees and what a receiver reads: the printed bytes, each line's kind, the exit status and the hand-off records. It never asserts an internal call.

### The fixture

A project folder named `project`, frozen in the CLI's test data, so that the shell prompt and the title read `~/project %` and `dft — project`. Each element spec lists the exact content it needs under its own Testing Decisions; this list gathers them, and a spec that needs more adds it there and here.

- **The ticket**, `.agents/0031-TICKET-book-a-demo-runtime.ticket.md`, in the client's own artifact folder: the heading `Run the book-a-demo request use case as a microruntime`, the line `**Status:** ready-for-agent`, and Markdown links to the PRD, the design doc and the spec beside it and to LangDrift's guardrail in `.langdrift/harness/guardrails/` (spec 0053).
- **Its three documents**, in `.agents/`, in the field-line convention spec 0053 reads:
  - `0002-PRD-book-a-demo.prd.md`, with `**Use case:** book-a-demo request`;
  - `0005-DD-book-a-demo-runtime.dd.md`, with `**cpuMillis:** 2000`, `**memoryMiB:** 3072` and `**Budget per run:** $25.00`;
  - `0011-SPEC-book-a-demo-runtime.spec.md`, whose `## Components` section lists `services/signup-api` (with `:8080` beside it), `services/notify-api`, `platform/nats`, `platform/redis` and `apps/signup`, in that order, and whose `## Acceptance criteria` section has three items.
- **The guardrail**, `.langdrift/harness/guardrails/0001-finops.gr.md`, as `dft init` scaffolds it: one criteria table with the columns `ID`, the rule and `Checked at`, whose rows are FIN-02 (`plan, run`), FIN-03 (`run`) and FIN-04 (`plan`), so that `reads harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04` follows (spec 0054).
- **Five components**, with no Dockerfile in any of them (spec 0066):
  - `services/signup-api/go.mod`;
  - `services/notify-api/package.json`, whose start script sets `PORT=8081`;
  - `platform/nats/nats.conf` with `port: 4222`;
  - `platform/redis/redis.conf` with `port 6379`;
  - `apps/signup/next.config.ts` beside a `package.json`.
- **Six episode records** under `.langdrift/microruntimes/episodes/`, as `ep-NNNN-<slug>.json` in spec 0052's format, all with `measured: null`, and the episode index `.langdrift/agent/state/rag/index.json` built by the local embedder's id:

  | Id | Slug | Similarity | Verdict |
  | --- | --- | --- | --- |
  | ep-0412 | checkout-preview-sizing | 0.81 | kept |
  | ep-0388 | notify-api-rollout | 0.74 | kept |
  | ep-0351 | signup-redesign-preview | 0.66 | kept |
  | ep-0297 | docs-search-index | 0.52 | dropped · docs only |
  | ep-0240 | batch-export-worker | 0.47 | dropped · other runtime |
  | ep-0188 | storefront-cache-warmup | 0.42 | dropped · no broker |

  The test builds the six vectors from the fixture ticket's query vector, so that their cosines are exactly these values (spec 0052). ep-0412, ep-0388 and ep-0351 share the stacks and roles; ep-0297 has no components; ep-0240 runs another workload; ep-0188 has no broker.
- **What `dft init` writes**, the table under Produces, except that the guardrail is the one above.
- **`.langdrift/harness/driver.json`**, pointing at the fixture's fake coding agent (spec 0058).
- **No `envelope.json`**, no input record and no mirror: the envelope and the budget come from the design doc.
- **A Git repository** with one commit holding all of the above, and a remote.

### The scripted boundaries

Every boundary outside the project is replaced at its port, never reached:

- the judge, by the recorded fixture judge (spec 0051);
- the image resolver, by the fixture resolver that resolves all five components (spec 0066);
- the coding agent, by the fixture's fake coding agent, a POSIX shell script that writes each stage's report and prints usage summing to the observed tokens of the applied effort (3,700,000 input and 240,000 output at medium), and leaves the branch's commits, the first of them with the subject `feat(demo): run book a demo as a microruntime` (specs 0058, 0059, 0063);
- the cluster, by a scripted cluster that reports five Deployments ready and removed (spec 0062);
- Git push and the forge, by a scripted forge that reports the push at `3f9c2e1`, pull request #19 and its four checks `build`, `check`, `test` and `guardrails` (specs 0063, 0064);
- the person, by one line on standard input per expected confirmation and one after `waiting    dev-ticket` (specs 0057, 0061);
- the clock.

The embedder is the local one and runs for real (spec 0052).

### The goldens

Frozen in the CLI's test data, transcribed from the English catalog of the simulation with its values filled:

- `plan`: the 42 lines phase 1 prints, from `# plan only …` to `# next: once the team authorizes …`; exit 0.
- `run`: the 25 lines phase 2 prints with the default intent (Claude Sonnet 5.5, medium, Planning and Review checked), from `# run · …` to `# next: the owner reviews #19 and decides whether to merge`; exit 0.
- `run-over-budget`: Claude Opus 5.5 at high effort, ending on `console › gate       reservation $40.80 exceeds the $25.00 budget · paused` and `# needs a decision: raise the budget or lower the effort · nothing was paid for`; exit 0.
- `run-to-merge`: Review unchecked, ending on `✓ ready      #19 to merge · $1.60 released to the envelope` and `# next: the merge stays the owner's`.
- `run-all-stages`: Spec, Tickets and Implement also checked, each adding its `waiting` and `confirmed` pair before its own work.
- `run-haiku-low`: Claude Haiku 5.5 at low effort, the applied pair reaching every stage call (spec 0058).
- The golden stops each element spec lists under its own Testing Decisions, among them the run whose input ends at `waiting    planning` (spec 0057) and the run whose input ends at `waiting    dev-ticket` (spec 0061): each prints a prefix of a golden above, writes one sentence to standard error and exits 1.
- The hand-off records after the plan and after the run, compared field by field with the four cards (spec 0050 defines them).
- Each line's kind, read from the session sink, compared with the simulation's kinds.

This spec's own assertions: the four comment lines at their positions; no line repeats the typed command; every actor line starts with `console › `; every line follows the six shapes; `dft init` on an empty folder writes exactly the files under Produces and prints nothing; a second `dft init` changes no byte.

### Unit tests

- The line format: padding of labels shorter than, equal to and longer than 10 columns, and the component row's three columns, counted in characters.
- The plan's argument: none, one, two.
- The intent selection: none, one, several authorized, one with a completed run.
- The stop rule: a scripted failure after some lines leaves standard output as a prefix of the golden, writes one sentence to standard error and exits 1.

Prior art: the CLI's plan and run tests, which copy a frozen fixture project and inject scripted dependencies, and the catalog tests of the session package.

## Out of Scope

- Every line and field owned by specs 0049 to 0066; this spec only frames them.
- The console and its stub API, including how the console renders `Play workflow` and `Play again`.
- The receiving surface's playback controls: Pause, Resume, Replay, `Play workflow` with `Starting…` and `Playing…`, `Play again`, the note `available once the plan is drafted`, which card is open, the scroll into view and the reduced-motion transcript. They are how a surface plays the session back, not engine behavior: `dft` neither pauses nor replays a run, and produces no field for them.
- The website's simulation itself; the golden is a frozen copy of it.
- Publishing and versioning of `@langdrift/cli`.

## Open Questions

- **Which intent `dft agent run` applies when several are authorized.** The simulation shows one. Until the owner decides, the run stops when more than one authorized intent has no completed run.
- **A second `dft agent run` of the same intent**, the product's counterpart of `Play again`, after a completed run or after a stop. Until the owner decides, a completed intent is not runnable again, and a stopped run's intent stays runnable and starts again from the gate.
- **`dft --help` and `dft --version`.** They are conventions of every CLI, but the session shows neither, so this spec removes them. The owner may keep them as conventions outside the session.
- **The document templates.** Spec 0053 relies on templates `dft init` scaffolds so that the team's documents carry the fields it reads, and the simulation's data says its numbers "follow the scaffold dft init produces". Plan and run never read the templates. Until the owner decides, `dft init` writes them under `.langdrift/templates/`.
- **The stop rule.** One sentence on stderr, exit 1, no retry and a silent teardown is this spec's choice for every failure the simulation never shows. The owner confirms it, or the simulation shows a failure first.
- **The fixture's file names.** The fixture reproduces the client's documents with the simulation's names (`0031-TICKET-…ticket.md`, `0002-PRD-…prd.md`, `0005-DD-…dd.md`, `0011-SPEC-…spec.md`, `0001-finops.gr.md`). Their uppercase parts and the `.dd.md` and `.gr.md` suffixes break the filesystem rules 1 and 2 for files in this repository. The owner states an exception for fixtures that reproduce a client project, or the fixture renames them.

## Further Notes

- The glossary entries **Connected agent**, **Input record**, **Mirror**, **Component Dockerfile** and **Lap** describe mechanisms these specs remove or replace; they become stale once these specs are implemented.
- The golden is a copy of the website's catalog, not a read of it, because a CLI test cannot reach the website's source without climbing paths. A change to the simulation's catalog updates the golden in the same change; review keeps them together.
- `lib/cli/README.md` describes the removed commands, flags and flows; it is rewritten when this spec's tickets land.
