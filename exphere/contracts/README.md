# Exphere contracts

The contracts of LangDrift Microruntimes, and only of Microruntimes. The `ex` CLI and the planned server both depend on them; they depend on nothing but the Go standard library. Money formats belong to FinOps and artifact formats to Vision or the team's tracker; neither lives here.

| Package | Contract |
| --- | --- |
| `configuration` | The `runtime.md` fence: exactly one `langdrift` JSON block, its fields and bounds, the build recipe allowlist, and the lexical rules for identifiers, paths, ports and images. |
| `protocol` | What the connected agent and the server exchange at `<console-url>/api/mcp`: the MCP version, the five tools and their request and response types. |

```sh
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
```

## The protocol

The agent calls `initialize`, then `notifications/initialized`, then `tools/call` with:

| Tool | Request | Response |
| --- | --- | --- |
| `langdrift.project.sync` | `{}` | `ProjectMembership`: the project's local folder and runtime ids |
| `langdrift.project.bind` | `runtimeId`, `environmentId` | `Binding`: an environment credential and the current `Intent` |
| `langdrift.bind` | `environmentId` | `Binding`, for a runtime secret |
| `langdrift.sync` | `Sync`: events, state, observation | `Intent`: desired state, revision, the authored fence, lease |
| `langdrift.claim` | `environmentId`, `intentVersion`, `revision` | `Claim`: the execution id of one start |

## The records of a run

These formats are the target of the next contract changes; the sample in `.langdrift/microruntimes/` follows them, and no code reads them yet. A run lives in `.langdrift/microruntimes/runs/NNNN-<slug>/`. A microruntime is the car, a run is a race, a lap is one pass of it.

| File | Holds | The console shows it in |
| --- | --- | --- |
| `plan.json` | Sources with role and hash, facts extracted per document, bounds with their sources, runs considered, sizing basis, suggestion, components proposed and executed, checks, the authorized `runtime.md` | Analysis and Plan |
| `run.json` | Kind, state, outcome, proposed ticket status, laps with each choice and outcome, stages with where a person took part, delivery: branch, commit, pull request, checks | Workflow and the run list |
| `NNNN-<slug>.run.md` | The memory: what the run teaches, and a `langdrift-run` JSON fence | The RAG |
| `laps/NNNN.lap.md` | One pass: steps, resources proposed, executed and observed, and a `langdrift-lap` JSON fence | The run detail and the RAG |
| `evidence/` | The operational record of each execution | Audits |

A run the plan considered is one entry of `plan.json` › `retrieval.considered`: the run, the criteria it matched, the share it matched as `score`, `kept` or `dropped`, and the reason.

Rules for every record: the console reads fields, never prose; anything not measured is `null`, never zero; a record keeps what happened, never the content it read; hashes of JSON records are taken after canonical serialization, so a formatter cannot change them.
