# Microruntimes contracts

The contracts of LangDrift Microruntimes, and only of Microruntimes. The LangDrift CLI imports this module, and the console reads its published JSON Schemas; it depends on nothing but the Go standard library. Money formats belong to FinOps and artifact formats to the team's tracker; neither lives here. It is the Go surface of the `@langdrift/core` package, at `lib/core/contracts/go`.

The home simulation at <https://langdrift.md/en/> defines what these shapes carry (ADR 0005): every field of the new shapes feeds a line `dft` prints or a field a card shows, and nothing else.

The module is published by git tag, `lib/core/contracts/go/vX.Y.Z`, the form Go expects of a module in a subdirectory. `v0.3.0` is the first to carry the simulation's shapes: the fence version 3 and the `judgment`, `stages` and `handoff` packages. The CLI's `go.mod` requires a tagged version, so the CLI builds outside this repository too, and the root `go.work` resolves the module from source for local work. A `require` of a version that is not tagged makes Go fetch it, even inside the workspace, so a tag comes before the `require` that names it.

| Package | Contract |
| --- | --- |
| `configuration` | The `runtime.md` fence, version 3: `ParseFence` reads it and `Validate` returns every rule it breaks. `schema/` holds its JSON Schema and the simulation's fence as a sample. The version 2 `Configuration` stays beside it until ticket 0324 (see [The older shapes](#the-older-shapes)). |
| `judgment` | The three answers the plan asks its judge for: the guardrail's `criteria`, the verdict on each retrieved `episodes` candidate, and the `suggestion` of a model and effort. `ParseAnswer` reads an answer and refuses anything outside its kind. `schema/` holds a JSON Schema and a sample per kind. |
| `stages` | The report each stage's skill writes, one shape per stage. `CheckReport` refuses a report outside its stage's shape. `schema/` holds a JSON Schema and a sample per stage. |
| `handoff` | What the console receives, as files: the plan's card data (`plan.json`), the recorded intent (`intent.json`) and the run's record (`run.json`). Go types only; no JSON Schema is published for them. |
| `judge` | The older judge contract, until ticket 0325. |
| `protocol` | The connected agent's channel, until ticket 0323. |

```sh
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
```

Each published JSON Schema and its Go check are tested against the same samples, so a change on one side is a change on both.

## The fence

`runtime.md` holds exactly one `langdrift` JSON fence. `ParseFence` reads it as version 3 and refuses any other version, and any field version 3 does not define, with "re-plan the runtime": an older fence is never misread as a newer one. It reads the shape only; `Validate` returns every problem, never stopping at the first, and the plan prints how many there are (`validate runtime.md draft against the schema · 0 errors`).

```json
{
	"version": 3,
	"runtimeId": "book-a-demo-runtime",
	"revision": 1,
	"cpuMillis": 2000,
	"memoryMiB": 3072,
	"model": "claude-sonnet-5-5",
	"effort": "medium",
	"components": [
		{
			"path": "services/signup-api",
			"stack": "Go",
			"role": "api",
			"port": 8080,
			"image": "book-a-demo-runtime/signup-api:r1",
			"cpuMillis": 250,
			"memoryMiB": 256,
			"replicas": 1
		}
	]
}
```

| Field | Rule |
| --- | --- |
| `version` | `3`. |
| `runtimeId` | A lowercase DNS label, the ticket's slug. |
| `revision` | At least 1. |
| `cpuMillis`, `memoryMiB` | The bound the team set, in millicores and MiB, within the per-execution limits. |
| `model`, `effort` | The suggested model, and `low`, `medium` or `high`. |
| `components` | 1 to 16 components; their paths are unique, their ports and images distinct, and their shares sum within the bound. |
| `components[].path` | The component's folder: a clean relative path outside `.langdrift`. |
| `components[].stack`, `role` | One of the stacks (`Go`, `Node.js`, `NATS`, `Redis`, `Next.js`) and roles (`api`, `broker`, `cache`, `web`) the Plan card's COMPONENTS rows name. |
| `components[].port`, `image` | 1 to 65535, and an image reference. |
| `components[].cpuMillis`, `memoryMiB` | The component's proposed share, positive. |
| `components[].replicas` | `1`. |

The JSON Schema checks the shape. `Validate` also enforces what the schema does not express: a clean relative path outside `.langdrift`, an image reference, unique paths, distinct ports and images, and shares within the bound; the schema says so in its description.

## The judge's answers

The plan asks its judge three questions and pays for nothing doing so, so no answer carries a price, usage or provenance.

| Kind | Answer | Feeds |
| --- | --- | --- |
| `criteria` | The guardrail criteria that apply to the ticket, each with its id and statement | `FIN-02 estimate before a paid call`, Analysis DOCUMENTS, the guardrail row |
| `episodes` | Per candidate, by its id `ep-NNNN`, the verdict `kept` with the reason `shares the stack and scale`, or `dropped` with one of `other runtime`, `no broker` or `docs only` | `judging · 3 share the stack and scale · 3 differ: …`, Analysis EPISODES |
| `suggestion` | A model and an effort, with the reason for each | `weighing effort · …`, `weighing model · …`, `✓ suggests   claude-sonnet-5-5 · effort medium`, SUGGESTED MODEL, SUGGESTED EFFORT |

## The stage reports

Each stage's skill writes one report, and the run's stage line prints its figures. Every field is required, a count is never negative, and a value not known is refused, never read as zero.

| Stage | Fields | Line |
| --- | --- | --- |
| `planning` | `openQuestionsSettled` | `console › grilling   3 open questions settled` |
| `spec` | `spec` (the spec's project-relative path), `revised`, `acceptanceCriteria` | `console › to-spec    0011-SPEC revised · 3 acceptance criteria` |
| `tickets` | `tickets`, `blockingEdgesRecorded` | `console › to-tickets 3 tickets · blocking edges recorded` |
| `implement` | `method` (`tdd`), `components`, `files`, `boundaryOnly`, `stoppedAtReservedPart` | `console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part` |
| `review` | `findings`, `fixed` | `console › code-review 2 findings · fixed` |

The to-spec line prints the spec's short id, the first two dash-separated parts of its file name; `SpecReport.ShortID` takes it from the path.

## The hand-off

Three records under the runtime's folder in the client's `.langdrift/`, `microruntimes/runtimes/<runtime id>/` (spec 0050), which a console only receives. Every figure is in its unit, integers for millicores, MiB, cents and tokens and two decimals for similarities, and a value not known is `null`, never zero. `Decode` reads a record and refuses a field its shape does not define.

| Record | Holds | Card |
| --- | --- | --- |
| `plan.json` | `ticket`, `analysis.documents`, `analysis.episodes`, `analysis.suggestion`, `plan.components`, `plan.sizing`, `plan.price`, `plan.authorization` (`pending`), `plan.evidence` (`file`, `state`, `revision`, `proposedTo`) and `workflow`; a section absent has not been revealed yet | Ticket, Analysis, Plan, Workflow |
| `intent.json` | `ticket`, `runtimeId`, `revision`, `model`, `effort`, `human` (one tick per stage, `planning` always set) and `status` (`proposed` or `authorized`); `Intent.Validate` returns every problem | `console › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1`, HUMAN IN THE LOOP |
| `run.json` | `run` (the run's id), `intent`, `stages` with each one's `report`, `effort`, `safeSeat`, `runtime`, `delivery` with the `pullRequest`, and `outcome` (`running`, `paused`, `stopped` or `completed`) | Workflow, PULL REQUEST, SAFE SEAT |

The fields inside each section join these types as the element specs that fill them land. RESERVED, OBSERVED and RELEASED are not in the run's record: they come from the run's entry in the money ledger, which the run's id names.

## The older shapes

They remain until the code that still reads them moves to the shapes above; each leaves in its own contract ticket.

| Shape | Leaves with | Once |
| --- | --- | --- |
| `protocol`: the MCP version, the five tools (`langdrift.project.sync`, `langdrift.project.bind`, `langdrift.bind`, `langdrift.sync`, `langdrift.claim`) and their request and response types | Ticket 0323 | The connected agent leaves the CLI (spec 0048) |
| The version 2 `Configuration` and `Parse`: `runtimeUid`, `projectId`, the runtime's `name`, `llmBudgetUsd`, `native` and the migration rule; per component `name`, `id`, `kind`, `dockerfile`, `readiness`, `envFile` and `passEnv`; the env-file parser, the variable-name rule, the Dockerfile path helpers, the derived `ImageTag` and the Markdown wrapper prose. Parsing a version 2 fence then fails with "re-plan the runtime", and the version 3 types take the names | Ticket 0324 | The plan writes the version 3 fence, the components resolve into it and the cluster applies it (specs 0055, 0066, 0062) |
| The whole `judge` package: the request and response envelope with usage and cost, the `sources` and `reserved-part` kinds, its free `reasons` list and its run-keyed episodes, with its schemas, samples and the `./judge/*` npm export | Ticket 0325 | The plan asks through `judgment`, reads its documents deterministically and plans Safe Seat without a judge (specs 0051, 0053, 0061) |
