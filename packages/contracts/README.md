# Microruntimes contracts

The contracts of LangDrift Microruntimes, and only of Microruntimes. The `ex` CLI and the planned server both depend on them; they depend on nothing but the Go standard library. Money formats belong to FinOps and artifact formats to the team's tracker; neither lives here. The root `go.work` resolves this module for the CLI; building the CLI outside this repository needs it published first.

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
| `plan.json` | The collection: the agent, its limits and the context hash; every source the ticket's links reached, of any file type, with its type, hash, depth, the link that reached it, its state, and the agent's verdict, reason and confidence; facts by topic with citations; the input records that set each limit; discrepancies between facts and input; runs considered; sizing basis; suggestion; components inferred, confirmed, proposed and executed; checks; the authorized `runtime.md` | Ticket, Analysis and Plan |
| `run.json` | Kind, state, outcome, proposed ticket status, laps with each choice and outcome, stages with where a person took part, delivery: branch, commit, pull request, checks | Workflow and the run list |
| `NNNN-<slug>.run.md` | The memory: what the run teaches, and a `langdrift-run` JSON fence | The RAG |
| `laps/NNNN.lap.md` | One pass: steps, resources proposed, executed and observed, and a `langdrift-lap` JSON fence | The run detail and the RAG |
| `evidence/` | The operational record of each execution | Audits |

A run the plan considered is one entry of `plan.json` › `retrieval.considered`: the run, the criteria it matched, the share it matched as `score`, `kept` or `dropped`, and the reason.

A source is `read`, `binary` when it is recorded by type, size and hash without its content, `unresolved` when the link leads nowhere, or `refused` when it is a secret, a symlink, or outside the project. Its role is inferred and optional; nothing depends on it.

## Input and mirrors

Every value that can block or allow a run, hardware limits, budgets, gates, the price and forecast sources, model and effort, is explicit input in the console or the CLI, held by the server with who set it, when, and the value before. A fact collected from a source never sets one. The server writes each record to `.langdrift/harness/` or `.langdrift/finops/` as a mirror, with a `mirror` block naming the record; no gate reads a mirror, and a mirror that differs from the server is an unauthorized change.

Rules for every record: the console reads fields, never prose; anything not measured is `null`, never zero; a record keeps what happened, never the content it read; hashes of JSON records are taken after canonical serialization, so a formatter cannot change them.
