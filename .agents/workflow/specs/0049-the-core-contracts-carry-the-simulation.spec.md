# The Core contracts carry the simulation

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E2.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 1 of 19, the first, before spec 0048. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

- [ ] 0318: the runtime.md fence v3, beside v2
- [ ] 0319: the judge's three answers the simulation shows
- [ ] 0320: the five stage reports
- [ ] 0321: the hand-off records
- [ ] 0322: `@langdrift/core` 0.3.0 publishes the new shapes
- [ ] 0323: the protocol package leaves
- [ ] 0324: the fence v3 is the only fence
- [ ] 0325: the old judge contract leaves

> **Note, 2026-11-05.** The simulation was played in production before these tickets were written, in both phases and the over-budget variant, and matches this spec's tables. One nuance is carried into 0320: the to-spec line prints the spec's short id, `0011-SPEC`, while the report's `spec` field holds its path. The tickets expand then contract: 0318 to 0322 add the new shapes beside the old ones, so the CLI keeps building in the workspace, and 0323 to 0325 remove the old ones once the specs that retire their last consumers land. The judge's answers are a new `judgment` package for that reason.

## Problem Statement

The LangDrift Core contracts (`lib/core/contracts/go`, published as `@langdrift/core`) are the shapes behind every line and card field of the simulation: the `runtime.md` fence the plan drafts and the run applies, the judge's answers, and what the CLI hands to whoever shows its data. Today they carry the wrong things.

- The fence has no per-component share, stack, role or replica count. It checks unique paths and names, not unique ports, and nothing compares the shares with the bound. Its validation stops at the first error, so the plan cannot say `0 errors` against a schema that does not exist. Its budget is a float in dollars described as "never a billing gate", while the simulation gates on it.
- The fence carries what nothing shows: a runtime uid, a project id, a required name, env files, passed variables, Dockerfiles, readiness, native migration and a derived image tag.
- The judge's contract answers a sources question that document reading no longer needs, carries free-text episode reasons keyed by run, prices every answer with a cost, usage and provenance against "nothing … paid for in this step", and keeps a reserved-part answer the simulation never asks the judge for.
- The protocol package describes a connected agent talking MCP to a console: tools, bindings, claims, leases, events, measurement windows, sync. None of it carries a card field, and the hand-off the simulation needs (the cards' data, the recorded intent, the run's record) has no shape at all.

## Solution

The contracts carry exactly the data the simulation shows, and nothing else:

- **The fence**, version 3: the runtime's id and revision, the bound, the suggestion, and per component its path, stack, role, port, image, share and replicas. It has a JSON Schema, and validation reports every error, so the plan can count them.
- **The judge's answers**: the three questions the plan asks (criteria, episodes, suggestion), each answer typed and checked against a closed contract.
- **The stage reports**: the five shapes a skill writes so its stage line can print structured results.
- **The hand-off**: the plan's card data, the recorded intent and the run's record, replacing the protocol package.

Each field's value belongs to the element spec that produces it. This spec decides that the field exists, its name, its unit and its validation.

## Responsibility

This element defines the shared shapes, their JSON Schemas where a schema is validated or published, and their parsing and validation, in the Go module and its npm surface. It does not:

- produce any value (specs 0050 to 0066);
- read or write any file of a client's project; the CLI does;
- decide where the judge runs (spec 0051) or how the console reads the hand-off (out of scope);
- hold records internal to the CLI: the episode record and index (spec 0052), the money ledger (spec 0056), the usage record (spec 0065) and the execution ledger (spec 0062) stay with their owners in the CLI.

## Produces

No file in a client's project. The module `lib/core/contracts/go` and the npm package `@langdrift/core` publish:

### `configuration`: the `runtime.md` fence, version 3

One `langdrift` JSON block per `runtime.md`. Its fields:

| Field | Type and unit | Example | Feeds |
| --- | --- | --- | --- |
| `version` | integer, `3` | `3` | `Parse` refuses any other version with "re-plan the runtime" |
| `runtimeId` | identifier | `book-a-demo-runtime` | `runtime    book-a-demo-runtime · …` (spec 0062) |
| `revision` | integer ≥ 1 | `1` | `✓ plan … revision 1`, EVIDENCE, `intent … runtime.md revision 1`, the ready line |
| `cpuMillis`, `memoryMiB` | integers, millicores and MiB, the bound | `2000`, `3072` | `✓ size … of 2000m · … of 3072 MiB`, BOUNDED, HEADROOM, the bars |
| `model` | price-list id, the suggestion | `claude-sonnet-5-5` | the draft's recommendation, beside the Analysis card's SUGGESTED MODEL (spec 0051); the applied choice is the recorded intent's |
| `effort` | `low`, `medium` or `high`, the suggestion | `medium` | the draft's recommendation, beside SUGGESTED EFFORT |
| `components[].path` | clean relative path, outside `.langdrift` | `services/signup-api` | COMPONENTS rows, inspect and component lines |
| `components[].stack` | one of `Go`, `Node.js`, `NATS`, `Redis`, `Next.js` | `Go` | COMPONENTS rows, component lines, the episode's stacks |
| `components[].role` | one of `api`, `broker`, `cache`, `web` | `api` | COMPONENTS rows, component lines |
| `components[].port` | integer 1 to 65535 | `8080` | the ports line, inspect lines |
| `components[].image` | image reference, resolved by spec 0066 | the fixture resolver's | `5 images resolved`, the cluster (spec 0062) |
| `components[].cpuMillis`, `components[].memoryMiB` | integers, the proposed share | `250`, `256` | component lines, COMPONENTS, totals, `5 Deployments ready · 1450m · 2304 MiB` |
| `components[].replicas` | integer, `1` | `1` | `5 one-replica Deployments`, SUGGESTED RESOURCE |

Validation collects every error instead of stopping at the first:

- the version, the identifier and the revision;
- the bound positive and within the existing per-execution limits (1 to 64,000 millicores, 1 to 1,048,576 MiB);
- the effort one of the three; the model not empty (whether it is in the price list is spec 0056's check);
- 1 to 16 components;
- each path clean, relative, outside `.langdrift` and unique;
- each stack and role in its closed set;
- each port valid and distinct;
- each image a valid reference and distinct;
- each share positive, and the sum of the shares at most the bound;
- each `replicas` equal to 1;
- no field the version does not define.

`Validate` returns the list of errors; its length is the count `tool call: validate runtime.md draft against the schema · 0 errors` prints (spec 0055). The JSON Schema of the fence is published beside the judge's schemas, and the Go validation and the schema are tested against the same samples. The sizing report that spec 0055 writes beside the fence is not part of the fence.

### `judge`: the three answers

The plan asks the judge three questions (spec 0051). Each answer is typed, and `ParseAnswer` refuses anything outside the contract:

| Kind | Answer | Feeds |
| --- | --- | --- |
| `criteria` | `criteria`: one entry per criterion, `id` (`FIN-02`) and `statement` (`estimate before a paid call`) | `judging · FIN-02 estimate before a paid call · …`, the FIN line |
| `episodes` | `episodes`: one entry per candidate, `id` (`ep-0412`), `verdict` (`kept` or `dropped`) and `reason` (`shares the stack and scale` for kept; `other runtime`, `no broker` or `docs only` for dropped) | `judging · 3 share the stack and scale · 3 differ: …`, EPISODES verdicts and reasons |
| `suggestion` | `model` (`claude-sonnet-5-5`), `effort` (`medium`), `effortReason`, `modelReason` | `weighing effort · …`, `weighing model · …`, `✓ suggests …`, SUGGESTED MODEL, SUGGESTED EFFORT |

The reasons are closed sets; the statements and the two weighing reasons are text. The JSON Schemas and one sample per kind stay under the judge's schema folder and on npm. The sizing basis (`reasoning · …`) needs no question (spec 0051) and has no answer kind.

### `stages`: the five stage reports

The report a skill writes at `LANGDRIFT_REPORT` (spec 0059), one shape per stage, each with a JSON Schema published on npm so a skill author can validate it:

| Stage | Fields | Feeds |
| --- | --- | --- |
| `planning` | `openQuestionsSettled` (integer) | `grilling   3 open questions settled` |
| `spec` | `spec` (path), `revised` (boolean), `acceptanceCriteria` (integer) | `to-spec    0011-SPEC revised · 3 acceptance criteria` |
| `tickets` | `tickets` (integer), `blockingEdgesRecorded` (boolean) | `to-tickets 3 tickets · blocking edges recorded` |
| `implement` | `method` (`tdd`), `components`, `files` (integers), `boundaryOnly`, `stoppedAtReservedPart` (booleans) | `implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part` |
| `review` | `findings` (integer), `fixed` (boolean) | `code-review 2 findings · fixed` |

### `handoff`: what the console receives

It replaces the `protocol` package. Three records, structured as spec 0050 says, in units (integers for millicores, MiB, cents and tokens; numbers with two decimals for similarities), with every value that is not known `null`, never zero:

- **The plan's card data** (`plan.json`): `ticket`, `analysis.documents`, `analysis.episodes`, `analysis.suggestion`, `plan.components`, `plan.sizing`, `plan.price`, `plan.authorization` (`pending`), `plan.evidence` (`file`, `state`, `revision`, `proposedTo`) and `workflow`. A field absent from the record has not been revealed yet.
- **The recorded intent** (`intent.json`): `ticket`, `runtimeId`, `revision`, `model`, `effort`, `human` (a boolean for each of `planning`, `spec`, `tickets`, `implement` and `review`, with `planning` always true) and `status` (`proposed` or `authorized`).
- **The run's record** (`run.json`, kept by spec 0065): the applied choice and per stage its state (spec 0057), the effort record (spec 0060), the dev ticket's state with `paidCalls` and the Safe Seat result (spec 0061), each stage's report figures (spec 0059), the pull request with its checks' state (specs 0063, 0064), the runtime's state (spec 0062), and the run's id and outcome (spec 0065).

The fields inside each section are named by the element that fills them, in its own spec; they join these types when that spec's tickets land. Go types only: no JSON Schema is published for the hand-off until a receiver needs one.

## Shows

Nothing of its own. Every terminal line and card field is printed or filled by another element from these shapes; the tables under Produces name which line or field each field feeds.

## User Stories

1. As an implementer of the plan, I want the fence to carry each component's share, stack, role, port, image and replicas, so that the draft holds what the Plan card shows.
2. As an implementer of the run, I want the fence to hold the proposed shares, not an equal split, so that the cluster applies exactly what the team saw.
3. As an implementer of the plan, I want validation to return every error, so that the validate line can print a count.
4. As an implementer of the plan, I want a JSON Schema for the fence, so that "against the schema" names a real schema.
5. As a team lead, I want the fence to refuse two components on one port, so that "5 ports distinct" is a checked fact.
6. As a team lead, I want the fence to refuse shares that exceed the bound, so that a draft can never propose more than the design doc allows.
7. As a team lead, I want the fence to refuse a replica count other than one, so that "one-replica Deployments" holds.
8. As a developer, I want an older fence refused with "re-plan the runtime", so that an old draft is never misread.
9. As the owner, I want the fence free of uid, project id, name, env files, passed variables, Dockerfile, readiness and migration fields, so that it carries only what the simulation shows.
10. As an implementer of the judge, I want three typed answers matching the three questions, so that every judged line has one place its words come from.
11. As an implementer of the judge, I want episode reasons from a closed set, so that the card's reasons read the same on every surface.
12. As an implementer of the judge, I want episodes identified as `ep-NNNN`, so that the judge's answer joins the EPISODES rows.
13. As the owner, I want no cost, usage or provenance on a judge answer, so that the contract does not contradict "nothing … paid for in this step".
14. As the owner, I want the reserved-part answer removed, so that the judge decides nothing the simulation does not show it deciding.
15. As a skill author, I want the five stage-report schemas published, so that I can validate my skill's report before a run.
16. As a console developer, I want the hand-off records typed in the Core contracts, so that the console and the CLI read the same shapes.
17. As a console developer, I want every hand-off value in units and unknown values `null`, so that I can format them per locale and never show a missing value as zero.
18. As the owner, I want the protocol package removed, so that no contract describes a channel the product does not have.
19. As a maintainer, I want the Go types and the JSON Schemas tested against the same samples, so that a change on one side is a change on both.
20. As a maintainer of the CLI, I want the module to keep depending on nothing but the Go standard library, so that the CLI builds anywhere.

## Implementation Decisions

### Remove

- From the fence: `runtimeUid`, `projectId` and the runtime's `name`. Serve no shown field; the runtime is named by `runtimeId`.
- From the fence: `llmBudgetUsd`. The budget is the design doc's, carried in the plan's hand-off and the money ledger in integer cents (spec 0056). Serves `budget     $25.00 per run …`.
- From the fence: `native` and the migration rule. Serve no shown line.
- From a component: `name`, `id`, `kind`, `dockerfile`, `readiness`, `envFile`, `passEnv`. Serve no shown field; `role` replaces `kind`, and the cluster's naming and readiness are spec 0062's.
- The env-file parser, the variable-name rule, the Dockerfile path helpers and the derived image tag `langdrift-<runtimeId>-<name>:r<revision>`. Serve no shown line; the image is resolved (spec 0066).
- The fence's Markdown wrapper prose ("Configuration is interpreted by the connected project agent …"); the file around the fence is spec 0055's.
- The judge's `sources` kind: per-source verdict, reason, confidence, role and summary, facts and figures. Document reading is deterministic (spec 0053).
- The judge's `reserved-part` kind and its fields (spec 0061).
- The judge's response envelope fields `usage`, `cost` (cents, provenance, illustrative), `refused`, `reason` and `judgeVersion`, and the request's `system`, `prompt`, `maxTokens` and `promptVersion`. Serve no shown line. Where the judge runs is open (spec 0051); if the answer needs a request across a process boundary, its shape returns here with only what that boundary needs.
- The suggestion's free `reasons` list: the two weighing reasons are named fields.
- The episode answer's `run` key: episodes are `ep-NNNN`.
- The whole `protocol` package: the MCP version, the five tools (`langdrift.project.sync`, `langdrift.project.bind`, `langdrift.bind`, `langdrift.sync`, `langdrift.claim`), `Binding`, `Claim`, `Event`, `Intent` with its desired state, lease, revocation, deletion, purge epoch and accepted events, `Operation`, `UnresolvedExecution`, `Observation` with its float `llmCostUsd`, `MeasurementWindow`, `ProjectMembership`, `Repository`, `Sync` and `AgentState`. Serve no shown field; the hand-off is files (spec 0050).
- The README's statement that "a fact collected from a source never sets one": the design doc sets the bound and the budget (specs 0053, 0054).

### Build

- The fence's `stack`, `role`, `cpuMillis`, `memoryMiB` and `replicas` per component. Serve the component lines, COMPONENTS, the totals and SUGGESTED RESOURCE.
- Validation that collects every error, with the distinct-port, distinct-image, sum-within-bound and one-replica rules. Serves `… against the schema · 0 errors`, `5 ports distinct`, `no collisions`, `1450m of 2000m`.
- The fence's JSON Schema, published on npm and tested with the Go validation against shared samples. Serves "against the schema".
- The `criteria`, `episodes` and `suggestion` answer kinds with their schemas and samples, and the closed reason set. Serve the judged lines and fields.
- The `stages` package: five report shapes, their schemas on npm, and a check that refuses a report outside them. Serves the five stage lines.
- The `handoff` package: the plan's card data, the recorded intent and the run's record, with `status` `proposed` or `authorized`. Serves every card field, the run's header and the `intent` line.

### Change

- The fence's version becomes 3. `Parse` refuses version 2 with "re-plan the runtime".
- The fence's `model` and `effort` are the suggestion, recorded as a recommendation; the applied choice lives in the recorded intent.
- The judge's `Criterion` field `applies` becomes `statement`.
- The README of the module and of `@langdrift/core` describe these four packages and nothing else.

### Interfaces

- To spec 0055: the fence type, its schema and `Validate`. To specs 0062 and 0066: the component fields.
- To spec 0051: the three answer kinds and `ParseAnswer`.
- To specs 0058 and 0059: the stage-report shapes and their check.
- To spec 0050: the hand-off types; to spec 0065: the run's record type.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** This spec adds: the fixture's written `runtime.md` fence parses as version 3 and validates with 0 errors; `plan.json` and `intent.json` decode into the hand-off types without unknown fields; each fixture stage report passes the stage check; the recorded judge's answers parse.
- Unit tests in the module, prior art being its existing schema-and-sample tests of the judge:
  - a fence with two equal ports and shares over the bound returns exactly two errors;
  - a fence with an unknown field, a version-2 fence, a stack or role outside its set, and `replicas: 2` are each refused;
  - each answer kind's sample parses, and a dropped episode with a reason outside the set is refused;
  - each stage-report sample passes, and one with a missing field is refused;
  - every published schema and its Go validation agree on the same samples.

## Out of Scope

- The values: specs 0050 to 0066.
- The CLI's internal records: episodes and index (spec 0052), money ledger (spec 0056), execution ledger (spec 0062), usage record (spec 0065).
- The console, its stub API and any schema it may later need.

## Further Notes

- The module's version tag moves with this change, and the CLI's `go.mod` requires the new tag, as the module's README already describes.
- The protocol's removal leaves the glossary's **Connected agent** entry without a mechanism; spec 0048 notes the stale glossary entries.
