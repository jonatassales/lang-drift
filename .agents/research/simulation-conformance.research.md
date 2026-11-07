# How far the engine is from the hero simulation

**Date:** 2026-11-05. **Question:** what does the home hero simulation show, line by line and field by field, and where do the LangDrift CLI, the LangDrift Core contracts and the client's `.langdrift/` mechanisms do more, less or something different? **Method:** the simulation was played on https://langdrift.md/en/ at 1440×900 and its terminal and panel text recorded every 250 ms; its source was read in `apps/website/src/app/[locale]/(pages)/platform/microruntimes/content/` and `apps/website/src/app/template/i18n/messages/microruntimes.json`. Five read-only audits then compared it with `lib/cli`, `lib/core`, the console's APIs and the workflow documents at commit `e892b4d`. ADR 0005 (Perfect Simulation Alignment) makes the simulation the definition; this note is the evidence the specs 0048 to 0066 cite. The console and its stub API are out of scope (ADR 0005), so their findings are not carried here.

Vocabulary: **MORE** the code does something the simulation does not show; **LESS** the simulation shows something the code does not do; **DIFFERENT** both exist and differ. Paths are under `lib/cli/src/` unless they start with `lib/core/`. "cat" is `presentation/cli/session/catalog.go`, the CLI's line catalog.

## 1. The simulation

### 1.1 Phase 1, `dft agent plan` (plays on load)

The command is typed into a terminal titled `dft — project`; the panel fades in once it is typed, and each card fills as the terminal reaches its step.

```
~/project %
dft agent plan 0031-TICKET-book-a-demo-runtime
# plan only · nothing is built, started or paid for in this step
console › tool call: read 0031-TICKET-book-a-demo-runtime
console › ready-for-agent · links 0002-PRD · 0005-DD · 0011-SPEC · guardrail 0001-finops
# analysis · reading the linked documents, then the project's episodes
console › tool call: read 0002-PRD · use case: book-a-demo request
console › tool call: read 0005-DD · envelope cpuMillis 2000 · memoryMiB 3072 · set by the team
console › tool call: read 0011-SPEC · 5 components named · 3 acceptance criteria
✓ read       4 documents · nothing outside the ticket's links
reads      harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04
console › judging · FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load
bounded    cpuMillis 2000 · memoryMiB 3072 · set by the team in the design doc
budget     $25.00 per run for model calls · set by the team in the design doc
console › RAG over project episodes · a model and effort for this ticket
console › embedding the ticket intent · 1 query vector
console › vector search · project episodes · top 6 by cosine similarity
console › 6 candidates · similarity 0.81 to 0.42
console › judging · 3 share the stack and scale · 3 differ: other runtime, no broker, docs only
console › dropped 3 below 0.55 · kept 3
console › 3 episodes cited · none with measured resources
console › reasoning · no cited episode measured cpu or memory · sizing follows the spec, not history
console › weighing effort · low skips the broker wiring · high is not warranted by 5 components
console › weighing model · Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome
✓ suggests   claude-sonnet-5-5 · effort medium
# plan · resolving the components, then sizing them inside the envelope
console › resolving 5 components named in the spec · paths, images, ports
console › tool call: inspect services/signup-api · go.mod · :8080
✓ services/signup-api   Go · api        250m · 256 MiB
console › tool call: inspect services/notify-api · package.json · :8081
✓ services/notify-api   Node.js · api   300m · 512 MiB
console › tool call: inspect platform/nats · nats.conf · :4222
✓ platform/nats         NATS · broker   200m · 256 MiB
console › tool call: inspect platform/redis · redis.conf · :6379
✓ platform/redis        Redis · cache   200m · 256 MiB
console › tool call: inspect apps/signup · next.config.ts · :3000
✓ apps/signup           Next.js · web   500m · 1024 MiB
console › cross-checking · 5 ports distinct · 5 images resolved · no collisions
✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares
limits     gpu and disk are not runtime.md fields · none requested
estimate   runtime has no price source · cpu and memory stay in millicores and MiB
console › tool call: validate runtime.md draft against the schema · 0 errors
✓ plan       runtime.md draft, revision 1 · proposed to the console
# next: once the team authorizes, dft agent run applies the recorded intent
```

### 1.2 The panel after phase 1

```
TICKET — 0031 · book-a-demo-runtime
  0031 · book-a-demo-runtime · ready-for-agent
  Run the book-a-demo request use case as a microruntime
  LINKED ARTIFACTS: PRD 0002-PRD-book-a-demo.prd.md · DESIGN DOC 0005-DD-book-a-demo-runtime.dd.md
                    SPEC 0011-SPEC-book-a-demo-runtime.spec.md · GUARDRAIL 0001-finops.gr.md
ANALYSIS — Claude Sonnet 5.5 · effort medium
  DOCUMENTS — 4 documents read
    0002-PRD-book-a-demo.prd.md ✓ Read — use case: book-a-demo request
    0005-DD-book-a-demo-runtime.dd.md ✓ Read — envelope 2,000m · 3,072 MiB · budget $25.00 per run · set by the team
    0011-SPEC-book-a-demo-runtime.spec.md ✓ Read — 5 components · 3 acceptance criteria
    0001-finops.gr.md ✓ Read — FIN-02 · FIN-03 · FIN-04 · review rule, no automatic enforcement
      FIN-02: estimate before a paid call · FIN-03: a budget before a paid run · FIN-04: size cpu and memory to the load
  EPISODES — threshold 0.55
    ep-0412 checkout-preview-sizing 0.81 kept
    ep-0388 notify-api-rollout 0.74 kept
    ep-0351 signup-redesign-preview 0.66 kept
    ep-0297 docs-search-index 0.52 dropped · docs only
    ep-0240 batch-export-worker 0.47 dropped · other runtime
    ep-0188 storefront-cache-warmup 0.42 dropped · no broker
  MEASURED RESOURCES — none cited · sizing follows the spec, not history
  SUGGESTED MODEL — Claude Sonnet 5.5 (claude-sonnet-5-5)
  SUGGESTED EFFORT — medium
  recommendation · the team authorizes
PLAN — 1,450m · 2,304 MiB · pending team authorization
  PROPOSED PLAN 1,450m · 2,304 MiB · pending team authorization
  CPU 1,450m of 2,000m · Memory 2,304 MiB of 3,072 MiB
  COMPONENTS (CPU / MEMORY)
    services/signup-api  Go       api     250m  256 MiB
    services/notify-api  Node.js  api     300m  512 MiB
    platform/nats        NATS     broker  200m  256 MiB
    platform/redis       Redis    cache   200m  256 MiB
    apps/signup          Next.js  web     500m  1,024 MiB
  5 components · 1,450m · 2,304 MiB
  ports 8080 · 8081 · 4222 · 6379 · 3000 · distinct · no collisions
  BOUNDED 2,000m · 3,072 MiB | HEADROOM 550m · 768 MiB | SUGGESTED RESOURCE 5 one-replica Deployments
  GPU · DISK not runtime.md fields | PRICE runtime: no price source · units only | VALIDATION schema · 0 errors
  EVIDENCE runtime.md draft · revision 1 · proposed to the console
WORKFLOW — Sonnet 5.5 · medium · $9.50 of $25.00
  HUMAN IN THE LOOP
    Planning (grilling) [x, always]   Spec (to-spec) [ ]   Tickets (to-tickets) [ ]   Implement (implement · tdd) [ ]
    Dev ticket (driven by the person)  Review (code-review) [x]   Merge ✓ always the owner's
  MODEL [Claude Opus 5.5 | Claude Sonnet 5.5 | Claude Haiku 5.5] — $2.00 / $10.00 per MTok · list price, Oct 2026
  EFFORT [low | medium | high]
  ESTIMATED SPEND $9.50 — 3.5M input · 250K output · uncached list prices — reserves $11.40 of the $25.00 budget
  [Play workflow]
```

### 1.3 The Workflow choices and the money

The Workflow card is the run's input. Its data (`microruntimes-workflow.data.ts`):

| Choice | Values | Default |
| --- | --- | --- |
| Model | Claude Opus 5.5 ($4.00 / $20.00 per MTok), Claude Sonnet 5.5 ($2.00 / $10.00), Claude Haiku 5.5 ($0.10 / $0.50); list prices read 2026-10-06, uncached | Sonnet 5.5 |
| Effort | low, medium, high | medium |
| Human in the loop | Planning (always checked), Spec, Tickets, Implement, Review | Planning and Review |
| Fixed | Dev ticket is always the person's; Merge is always the owner's | |

| Effort | Expected tokens (input · output) | Observed tokens |
| --- | --- | --- |
| low | 2,000K · 120K | 2,100K · 125K |
| medium | 3,500K · 250K | 3,700K · 240K |
| high | 6,000K · 500K | 6,300K · 520K |

- Spend in cents is `round((inputThousands × inputCents + outputThousands × outputCents) / 1000)`. Medium on Sonnet gives $9.50 estimated and $9.80 observed.
- The reservation is `ceil(estimate × (100 + 20) / 100)` with a 20% margin, giving $11.40. The budget is $25.00 per run, from the design doc.
- Released is reserved minus observed, $1.60. The run pauses at the gate when the reservation exceeds the budget.
- Each stage the person must confirm prints a `waiting` line and a `confirmed` line before its own work. An unchecked stage runs without waiting.
- Review checked ends with "ready for review"; unchecked ends with "to merge".

### 1.4 Phase 2, `dft agent run` (plays on "Play workflow")

```
~/project %
dft agent run
# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team
console › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1
console › estimate   $9.50 · 3.5M input · 250K output · list prices · FIN-02
console › budget     $25.00 per run · set by the team · FIN-03
console › gate       reserved $11.40 of $25.00 · opened
console › waiting    planning · the owner confirms
✓ confirmed  planning · by the owner
console › grilling   3 open questions settled
console › to-spec    0011-SPEC revised · 3 acceptance criteria
console › to-tickets 3 tickets · blocking edges recorded
console › hook       effort recorded · implement unblocked
console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part
console › waiting    dev-ticket · the person writes the reserved part
✓ closed     dev-ticket · written without an agent · no paid call
console › code-review 2 findings · fixed
runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB
✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record
✓ pushed     feat/0031-book-a-demo-runtime · 3f9c2e1
✓ pr #19     feat(demo): run book a demo as a microruntime
checks     build · check · test · guardrails · running
✓ checks     4 of 4 passed
✓ Safe Seat  met
runtime    stopped · 5 Deployments removed · nothing left
✓ ready      #19 for review · $1.60 released to the envelope
# next: the owner reviews #19 and decides whether to merge
```

Variants the data holds:
- **Over budget**, replacing everything from the gate on:
  - `console › gate       reservation {reserved} exceeds the {budget} budget · paused`
  - `# needs a decision: raise the budget or lower the effort · nothing was paid for`
- **Review unchecked**, replacing the last two lines:
  - `✓ ready      #{number} to merge · {released} released to the envelope`
  - `# next: the merge stays the owner's`
- **Stages checked**: each checked stage among Spec, Tickets and Implement adds its own `waiting {stage}` and `confirmed {stage}` pair.

### 1.5 The panel after phase 2

Only the Workflow card changes:

```
WORKFLOW — #19 · ready for review · $9.80
  HUMAN IN THE LOOP: Planning ✓ · Spec ✓ · Tickets ✓ · Implement ✓ · Dev ticket ✓ · Review "waiting for the owner" · Merge "always the owner's"
  RESERVED $11.40 | OBSERVED $9.80 | RELEASED $1.60
  PULL REQUEST #19 · 3f9c2e1 · ready for review
  SAFE SEAT met
```

The Plan card still reads "pending team authorization".

### 1.6 Contradictions inside the simulation

These stay open for the owner. Until the owner decides, the specs follow the simulation literally.

1. The Plan card reads "pending team authorization" after the run, whose header says "the recorded intent, authorized by the team".
2. The terminal's "read 4 documents" counts the ticket, the PRD, the DD and the SPEC, then reads the guardrail on its own line. The card's "4 documents read" lists the PRD, the DD, the SPEC and the guardrail.
3. `reads … FIN-02 · FIN-04` omits FIN-03. The next line, and the card, judge FIN-02, FIN-03 and FIN-04.

## 2. Who owns each line and field

Each terminal line and card field belongs to one element, and its spec defines it. An element may own nothing visible and still be required. The LangDrift Core contracts (E2) and the stage driver (E11) are two such elements.

| Element | Spec | Terminal lines (phase) | Card fields |
| --- | --- | --- | --- |
| E1 LangDrift CLI | 0048 | `~/project %` and the typed commands, `# plan only …`, `# analysis …`, `# plan · resolving …`, `# next: once the team authorizes …` (1); the command (2) | the terminal title `dft — project`, Play workflow / Play again |
| E2 LangDrift Core | 0049 | none: the shapes behind every other element | none |
| E3 Data hand-off | 0050 | `✓ plan runtime.md draft, revision 1 · proposed to the console` (1); `# run · … the recorded intent, authorized by the team`, `intent …` (2) | Plan "pending team authorization" (header and proposed plan), Plan EVIDENCE |
| E4 Judge | 0051 | `judging · FIN-02 …`, `judging · 3 share …`, `reasoning …`, `weighing effort …`, `weighing model …`, `✓ suggests …` (1) | Analysis header, SUGGESTED MODEL, SUGGESTED EFFORT, "recommendation · the team authorizes" |
| E5 RAG episode index | 0052 | `RAG over project episodes …`, `embedding …`, `vector search …`, `6 candidates …`, `dropped 3 below 0.55 · kept 3`, `3 episodes cited …` (1) | EPISODES (threshold, six rows), MEASURED RESOURCES |
| E6 Project documents | 0053 | `tool call: read 0031-TICKET…`, `ready-for-agent · links …`, the three `tool call: read` document lines, `✓ read 4 documents …` (1) | the whole Ticket card; DOCUMENTS header and the PRD, DD and SPEC rows |
| E7 Harness and guardrails | 0054 | `reads harness/guardrails/…`, `bounded …` (1) | the guardrail document row and its FIN line |
| E8 `runtime.md` | 0055 | `✓ size …`, `limits …`, `tool call: validate runtime.md draft …` (1) | Plan header figures, PROPOSED PLAN figures, CPU and Memory bars, the components total line, BOUNDED, HEADROOM, SUGGESTED RESOURCE, GPU · DISK, VALIDATION |
| E9 Money gates | 0056 | `budget $25.00 per run …`, `estimate runtime has no price source …` (1); `estimate`, `budget`, `gate … opened`, the over-budget `gate … paused` and `# needs a decision …` (2) | Plan PRICE; Workflow header (phase 1), the MODEL price line, ESTIMATED SPEND and its reservation; RESERVED, OBSERVED, RELEASED |
| E10 Workflow and human checkpoints | 0057 | every `waiting {stage}` and `✓ confirmed {stage}`; `# next: the owner reviews …` or `# next: the merge stays the owner's` (2) | HUMAN IN THE LOOP, the MODEL and EFFORT choices, stage states after the run |
| E11 Stage driver | 0058 | none: runs every skill line of E12 | none |
| E12 Agent skills | 0059 | `grilling …`, `to-spec …`, `to-tickets …`, `implement …`, `code-review …` (2) | none |
| E13 Hooks | 0060 | `hook effort recorded · implement unblocked` (2) | none |
| E14 Safe Seat | 0061 | `waiting dev-ticket …`, `✓ closed dev-ticket …`, `✓ Safe Seat met` (2) | SAFE SEAT |
| E15 Local cluster | 0062 | `runtime … 5 Deployments ready …`, `runtime stopped · 5 Deployments removed · nothing left` (2) | none |
| E16 Delivery observation | 0063 | `✓ pushed …`, `✓ pr #19 …`, `✓ ready #19 for review …` or `… to merge …` (2) | Workflow header (phase 2), PULL REQUEST |
| E17 Client CI checks | 0064 | `checks … running`, `✓ checks 4 of 4 passed` (2) | none |
| E18 Evidence and ledger | 0065 | `✓ completed implementation · … · usage record` (2) | none |
| E19 Client components | 0066 | `resolving 5 components …`, the five `inspect` lines, the five `✓ <path>` lines, `cross-checking …` (1) | the five COMPONENTS rows, the ports line |

## 3. Findings by element

### E1 LangDrift CLI

- **DIFFERENT:** every actor line is prefixed `agent › `, not `console › ` (`presentation/cli/session/format.go:15`).
- **MORE:** commands beyond `init`, `plan` and `run`, under `presentation/cli/` and `catalog/`, all removed by ADR 0005:
  - `doctor`, `status`, `reindex`, `confirm`, `report`, `validate`, `set`;
  - `cluster` `create`, `delete`, `inspect`, `load-image`;
  - `runtime` `up`, `down`, `inspect`, `logs`, `forward`, `hook`, `migrate`;
  - `connect`, `supervise`, `project`, `stop`;
  - the `harness` and `finops` input commands and `--json` NDJSON.
- **DIFFERENT:** `agent plan` takes a file path relative to the project; a bare ticket id resolves to nothing (`application/plan/helpers.go:50`, `infrastructure/collection/collector.go:32`).
- **MORE:** the plan echoes the command as its first output line, so it prints twice (`application/plan/planning.go:88`, cat:22).
- **DIFFERENT:** the scope line reads "nothing is built or started · the planner's calls are paid and gated", not "nothing is built, started or paid for in this step" (cat:26).
- **DIFFERENT:** `agent run` requires `RUNTIME`; the simulation types `dft agent run` (`presentation/cli/catalog/agent-run.go:29-31`).
- **MORE:** the plan's handoff line adds "authorize with dft agent run" (cat:129).
- **MATCH:** the analysis, plan and next comments (cat:35, 131).

### E2 LangDrift Core

`lib/core/contracts/go/`:
- **LESS:** a component has no `cpuMillis`, `memoryMiB`, stack or role (`configuration/component.go:17-27`). A role is a free `Kind` string.
- **LESS:** validation checks unique paths and names, not unique ports, and no envelope sum or headroom (`configuration/configuration.go:109-133`).
- **LESS:** no replica count or workload type; no draft, proposed or authorized status beside `Revision` (`configuration.go:19-34`).
- **DIFFERENT:** validation is Go `Validate`, which stops at the first error. There is no `runtime.md` JSON Schema and no error count (`configuration/parse.go:18-57`).
- **DIFFERENT:** `llmBudgetUsd` is "a per-execution display envelope, never a billing gate" (`configuration.go:20-21`). The simulation gates on it.
- **MORE:** `runtimeUid`, `projectId`, `name` (a required DNS label), `envFile`, `passEnv`, the env-file parser, native migration, Dockerfile-derived image tags (`configuration.go:22-25`, `component.go:19,25-26`, `native.go`, `image-tag.go`).
- **Judge contract:**
  - **LESS:** an episode answer carries only `run` and `reason`, with no similarity, keep or drop, threshold or measured resources (`judge/episode-reason.go:4-7`). Drop reasons are free text, not the closed set.
  - **DIFFERENT:** episode ids are runs, not `ep-NNNN` with a slug (`judge/episode-reason.go:5`).
  - **DIFFERENT:** a source verdict is `kept`/`discarded` with a confidence (`judge/source-verdict.go`), where the simulation shows `✓ Read`.
  - **LESS:** figures are resource-only, so the budget is not a number (`judge/figure.go`).
  - **LESS:** a criterion has no enforcement mode (`judge/criterion.go`).
  - **DIFFERENT:** every judgment carries a cost (`judge/cost.go`) against "nothing … paid for in this step".
  - **MORE:** `usage`, `provenance`, `illustrative`, `judgeVersion`, the request's prompt fields, `refused`, and the reserved-part answer's fields.
  - **DIFFERENT:** "A fact collected from a source never sets one" (README:81) against the envelope and budget set in the design doc.
- **Protocol:**
  - **LESS:** no tool proposes a plan (`protocol/tool-names.go:4-19`).
  - **LESS:** events carry no plan steps and no ticket id (`protocol/event.go`).
  - **LESS:** the intent has no human-in-the-loop choices (`protocol/intent.go`).
  - **LESS:** no estimate, reservation, release, gate state, stage or confirmation, pull request, commit, checks or Safe Seat.
  - **DIFFERENT:** `llmCostUsd` is a float, not cents with tokens (`protocol/observation.go`).
  - **DIFFERENT:** agent states have no "ready" and no removed count (`protocol/agent-state.go`).
  - **MORE:** measurement windows, `executionId`, leases and revocation, `bind`, repository identity, laps and source-hash records.
- **MATCH:**
  - model and effort (low, medium, high);
  - the runtime-wide `cpuMillis` and `memoryMiB`;
  - no gpu, disk or price field;
  - the suggestion's `effortReason`, `modelReason` and `reasons`;
  - `criteria[{id, applies}]`;
  - the intent's fence and revision.

### E3 Data hand-off

- **LESS:** nothing in `application/plan` sends the draft anywhere. It writes `runtime.md` locally (`planning.go:662`). The CLI posts MCP tools to `<url>/api/mcp` (`infrastructure/connected/mcp.go:43-44,146-171`), which no endpoint serves.
- **DIFFERENT:** `agent run` itself turns a proposed plan into an authorized one, `By` `owner` and `Via` `cli`, and prints "authorized by owner" (`application/execution/start.go:181-189`). The simulation's intent arrives "authorized by the team".
- **LESS:** the protocol carries none of the cards' data (see E2).
- **MATCH:** "the recorded intent" (fence plus revision) exists as `Intent`.

### E4 Judge

- **DIFFERENT:**
  - the judge is called through the Anthropic SDK (`infrastructure/judge/anthropic.go`) or the console's `/api/judge` (spec 0045, closed);
  - its planning calls are priced, reserved and gated FG-01/02/03, and can refuse the plan (`planning.go:119-143`), against "nothing … paid for in this step";
  - the guardrail judging line is the judge's free "id applies" text, and it is skipped without a judge (`planning.go:207-209, 261-263`).
- **DIFFERENT:** the judge rewrites only the reason text; it never keeps or drops an episode, and "share the stack and scale" is fixed text (`application/plan/retrieval.go:100-140`, cat:81).
- **DIFFERENT:** "sizing follows the inputs, not history", against "the spec" (`application/plan/retrieval.go:27`).
- **LESS (failure-path):** the weighing lines appear only with a judge and a price list (`planning.go:349-372`).
- **MORE:**
  - judge-unavailable and judge-stopped warnings, the `--no-judge` mode (`planning.go:134-138, 295-303`);
  - the judge's figure discrepancies, "discrepancy … the input holds" (`planning.go:214-216, 319-337`).
- **MATCH:** `✓ suggests <model> · effort <effort>`; models limited to the price list; low, medium and high.
- **Open:** where the judge runs and who pays for planning, given "nothing … paid for in this step".

### E5 RAG episode index

- **DIFFERENT:** it embeds the ticket's title and first 4 KiB with a local hashed TF-IDF embedder (`domain/rag/embedder.go:28-57`, `application/plan/retrieval.go:173-181`). "1 query vector" and cosine hold.
- **MORE:** retrieval also scores same component path, same stack, measured resources and measured tokens (`domain/rag/retrieval.go:94-115`).
- **DIFFERENT:** an episode is kept at similarity ≥ 0.55 **or** with ≥ 2 matched criteria, so one below 0.55 can be kept (`domain/rag/retrieval.go:79`).
- **DIFFERENT:** episodes are identified by runtime id, not `ep-NNNN` and a slug (`infrastructure/rag/index.go:97-122`).
- **MORE:** when cited episodes measured resources, every component gets an equal share of the peak × 1.5 (`application/plan/retrieval.go:83-88,186-191`).
- **MORE:** the episode memory and the index rebuild run at the end of `run` (`application/execution/memory.go:15-45`). This is required, because the next plan needs episodes, but it prints nothing.
- **MATCH:** `top 6` (`rag.TopK`), the 0.55 threshold, and the candidate, drop and "none with measured resources" line formats (cat:63-88).

### E6 Project documents

- **MORE:** links are followed breadth-first to depth 3, with `--max-depth/--max-files/--max-bytes` (`infrastructure/collection/collector.go:66-77`, `domain/collection/limits.go:21`), against "nothing outside the ticket's links".
- **DIFFERENT:** links print as full file names, without short ids or the "guardrail" label (`helpers.go:184-193`).
- **DIFFERENT:** each document's summary is the judge's ≤ 8-word free text, or "read". Nothing deterministic extracts "use case", the envelope, the budget, "5 components named" or "3 acceptance criteria" (`planning.go:186-197`).
- **DIFFERENT:** "✓ read N documents" counts the ticket too (`planning.go:199`). A "N refused" variant exists.
- **LESS:** roles are set only for the ticket and the guardrail; PRD, design doc and spec roles are free text from the judge (`records.go:110-125`).
- **LESS:** the source record keeps no summary (`records.go:97-128`).
- **MORE:** the guardrail gets its own `tool call: read` line (`planning.go:186`).
- **MATCH:** ticket status from `**Status:**` and title from the heading (`domain/collection/ticket.go`).

### E7 Harness and guardrails

- **DIFFERENT (unverified):** the `reads` line lists every `FIN-\d{2}` id in file order (`helpers.go:170-182`).
- **DIFFERENT:** the envelope comes from recorded input or `--cpu-millis/--memory-mib` and prints "set by X in the CLI · record Y", not "set by the team in the design doc" (`planning.go:305-317`, `application/plan/guardrails.go:34`, cat:55).
- **MATCH:** FIN rules have no automatic enforcement; their text only reaches the judge.

### E8 `runtime.md`

- **DIFFERENT:** what is recorded to execute is an equal split of the whole envelope (400m and 614 MiB each), not the proposed shares (`records.go:19,258`, `domain/run/plan.go:282-302`, `infrastructure/connected/native.go:137`).
- **LESS (failure-path):** nothing fails when the proposal exceeds the envelope; headroom is clamped to 0 (`helpers.go:151-163`).
- **MORE:** a blocked plan, "✗ plan … blocked", is still written (`planning.go:578-581`, cat:140). A lock refuses to plan while a run is open (`service.go:117-129`).
- **MORE:** model and effort are always taken from a pre-recorded choice input (`planning.go:613-627`).
- **MATCH:** the size, limits and validation line formats; revision 1 on the first plan; the fence's headroom, `Materialization{Deployments, ReplicasEach:1}` and `Suggestion{By:"recommendation", AuthorizedBy:"pending"}` (`records.go:35-73`).

### E9 Money gates

- **MATCH:**
  - the estimate formula and its rounding;
  - the reserve formula with ceil, $11.40 at 20%;
  - released = reserved − observed;
  - the estimate, gate-opened and over-budget pause lines, and "# needs a decision…" (`domain/finops/arithmetic.go:10-26`, `finish.go:367`, `presentation/cli/session/catalog_test.go:48-79`).
- **DIFFERENT:** the reserve margin is whatever the owner set at init, with no default (`presentation/cli/catalog/init-guide.go:62,85-98`). The simulation's is 20%.
- **DIFFERENT:** the budget prints "set by owner" from the budget input (`start.go:212-213`), not "set by the team" from the design doc.
- **MORE:**
  - FG-01, a price source is required;
  - FG-02, a ticket allocation is required;
  - FG-03 checks the ticket allocation and the period remainder;
  - FG-04 pauses when observed spend reaches the reservation or usage goes unreported;
  - FG-05, earlier overspend, with `--acknowledge-overspend`;
  - FG-06, provisional billing in the ledger.

  All of these sit in `domain/finops/arithmetic.go:55-120`, `start.go:215-228` and `play.go:275-303`. Their "gate blocked" lines print too.
- **MORE:**
  - the planning spend and its reservation (`records.go:185-234,354-370`);
  - a missing forecast refuses the run (`start.go:201-204`);
  - any `prices.json` model is accepted, not just the three.
- **LESS:** the plan records no estimate or reservation for the run (`records.go:131-166`).

### E10 Workflow and human checkpoints

- **DIFFERENT:**
  - which stages wait comes from `workflow.json` `defaults.humanAt` and `human: required`. No per-run choice exists (`application/plan/records.go:169-181`, `play.go:94`);
  - model and effort come from the plan's suggestion or `--model`/`--effort` (`start.go:315-330`);
  - the sample workflow has no Tickets or Dev ticket stage, and puts `code-review` inside Implement (`testdata/langdrift/sample/harness/workflow.json`, `lib/cli/README.md`).
- **MORE:**
  - a terminal `Confirm planning? [y/N]` prompt on stderr;
  - without a terminal, an exit in the waiting state with the runtime up, then `dft agent confirm` and a second `run`, which prints "lap 0001 open · … · resumed" (`agent-run.go:35-39`, `play.go:138-148`, `confirm.go:34-93`, `start.go:164-175`).
- **DIFFERENT:** "by owner", from the `--by` default, not "by the owner" (cat:186-188).
- **DIFFERENT:** "for review" versus "to merge" also depends on delivery being clean (`finish.go:411-452`). The simulation decides on the Review checkbox alone.
- **MORE:** a proposed ticket status, `in-review` or `ready-to-merge`, is written and never printed (`finish.go:373-403`).
- **MATCH:** Planning always human; Merge never played (`play.go:111-124`); efforts low, medium and high.

### E11 Stage driver

- **DIFFERENT:** the chosen model and effort never reach the coding agent. `StageCall` has no model or effort, and the driver's placeholders are skill, ticket, run, stage, lap and prompt (`application/execution/ports.go:126-134`, `infrastructure/stages/driver.go:36-44,128-130,164-196`). Observed cost is still priced at the chosen model (`play.go:327-330`).
- **MORE (failure-path):** a failed or timed-out skill fails the lap; a cut-off call holds the reservation (`play.go:244-251,368-378`).
- **MORE:** " · no usage reported" is appended to a stage line (`driver.go:47,211-218`).

### E12 Agent skills

- **DIFFERENT:** a stage line is `<skill> <summary>`, where the summary is the agent's last output line cut to 120 characters (`play.go:256-260`). The simulation's lines carry structured results: "3 open questions settled", "0011-SPEC revised · 3 acceptance criteria", "3 tickets · blocking edges recorded", "tdd · 5 components · 14 files · boundary only · stops at the reserved part", "2 findings · fixed".
- **DIFFERENT:** each skill of a stage is its own paid call and line; the simulation shows one `implement` line for "implement · tdd" (`play.go:186-200`).
- **DIFFERENT:** code-review runs where the workflow puts it. In the sample that is inside Implement, before any dev ticket; the simulation runs it after the dev ticket, under Review, without waiting.
- **LESS:** no "boundary only" and no stop at the reserved part.

### E13 Hooks

- **DIFFERENT:** a hook runs a committed `sh` script and prints "hook <script> · exit N". Effort is not in its environment, and nothing records effort or unblocks implement (`play.go:216-233`, `infrastructure/stages/driver.go:97-124,198-206`). The simulation shows "hook effort recorded · implement unblocked".

### E14 Safe Seat

- **LESS:** `lib/cli` has no dev-ticket stage, no "waiting dev-ticket" or "closed dev-ticket" line, and no Safe Seat line, check or reveal (`finish.go:37-49`, `presentation/cli/session/reveals.go:46-57`). The reserved part exists only as a judge answer kind in `lib/core`.
- **Closed specs went further:** spec 0033 (SEAT-01/02/03, overrides, `harness safe-seat set`, the reserved part on the Plan card) added what the simulation does not show.

### E15 Local cluster

- **DIFFERENT (order):** the runtime is brought up right after the gate, before the first stage, and kept up through human waits (`next-phase.go:17-18`, `play.go:62-65`, `finish.go:57-89`; ticket 0225, commit c6fd7a3). The simulation brings it up after code-review.
- **DIFFERENT:** "nothing left" is printed while the kind cluster is always kept, `ClusterKept: true` (`infrastructure/lifecycle/native.go:158-163`, `finish.go:341-346`).
- **MORE (failure-path):**
  - a 20-minute bring-up timeout, and a failed bring-up ends the lap (`finish.go:53,62-78,472-517`);
  - an unconfirmed cleanup prints "runtime unknown … the reservation is held" (`finish.go:332-339`).
- **MORE:** components are built from their Dockerfiles (ADR 0004). The simulation inspects manifests and "resolves" images (`planning.go:471-545`).
- **MATCH:** the ready line's format (cat:195-197).

### E16 Delivery observation

- **DIFFERENT (who):** the simulation's run pushes and opens the pull request. The code only observes with `git` and `gh` (`infrastructure/delivery/observer.go:1-5,30-85`). "Pushed" requires an upstream 0 commits ahead; otherwise "pushed unknown" or "pr unknown" (`finish.go:181-196`).
- **LESS:** the `pushed` event has no hash reveal (`reveals.go:46-57`).
- **MATCH:** the pushed, pr and ready line formats.

### E17 Client CI checks

- **MORE (failure-path):** checks are followed up to `--checks-wait` (10 min), then "checks unknown"; a "✗ checks … failed" line exists (`finish.go:211-271`).
- **MATCH:** "checks … running" and "✓ checks N of N passed".

### E18 Evidence and ledger

- **MATCH:** "✓ completed implementation · $X observed of $Y reserved · usage record" (`finish.go:154-162`).
- **MORE:** printed nowhere, but written:
  - the planning spend entry;
  - the run record;
  - laps;
  - FG-06 billing;
  - the proposed ticket status;
  - the episode memory (`records.go:131-166,354-370`, `finish.go:360-403`).

  Each spec decides whether that record is required to produce a shown field.

### E19 Client components

- **DIFFERENT (source):** components come from `--component name=,path=…` and print "named in the CLI", not "named in the spec" (`application/plan/service.go:151-152`, `planning.go:405`).
- **DIFFERENT:** the port comes only from `port=` or the Dockerfile's `EXPOSE`; the port found in the manifest is ignored (`planning.go:494-498`, `domain/manifest/manifest.go:84-96`).
- **DIFFERENT:** `next.config.ts` is not a known manifest. Next.js is detected from `package.json` and prints as "package.json" and "Node.js+Next.js" (`domain/manifest/manifest.go:44,68-74,160-168`).
- **MORE:** Dockerfile, `HEALTHCHECK`, env-file, readiness and within-ceiling checks, each blocking (`planning.go:471-545`). A NATS or Redis folder without a Dockerfile blocks the plan.
- **DIFFERENT:** "5 images resolved" is the component count; images are derived tags, and nothing is resolved (`planning.go:418`, `records.go:250`).
- **MATCH:** the inspect and component line formats and the sizing table:
  - Go api 250m · 256 MiB;
  - Node.js api 300m · 512 MiB;
  - broker and cache 200m · 256 MiB;
  - web 500m · 1024 MiB (`domain/manifest/manifest.go:173-190`).

  The ports-distinct check also matches.

## 4. Documents that disagreed

ADR 0005 closed all of these on 2026-11-05. The list says where each one went wrong.

| Document | Where it went wrong |
| --- | --- |
| spec 0007, microruntimes hero | Its framing agreed, but its script is stale. |
| spec 0010 | "Guardrails come only from input… Nothing collected… sets one"; the planner's call is paid work. |
| spec 0016, the session | It holds most of the matching lines. It also has: the `agent ›` actor; the CLI-input envelope and budget; gated planning; the equal split; `run <runtime>`; FG-01 to FG-06; `status`, `reindex`, `--json`. |
| spec 0032, components by Dockerfile | Bring-up first; Dockerfile checks; `--component`; the port from `EXPOSE`. |
| ADR 0004 | Superseded where the simulation shows manifests and resolved images. |
| spec 0033, Safe Seat | Had more than the simulation: SEAT rules, overrides, the reserved part on the Plan card. |
| spec 0045, the judge as a platform service | Bring-up first; a planning reservation; a local RAG with `agent ›` lines. |
| spec 0024, init asks FinOps | A monthly budget, a per-run cap and a margin at init, not from the design doc. |
| `lib/cli/README.md`, `lib/core/contracts/go/README.md` | "Executing the runtime authorizes the plan"; bring-up before any stage; the stale workflow example; "a fact … never sets one". |
| tickets 0225, 0224, 0226 | Implemented: bring-up first, planning reservation, `dft agent stop`. |
