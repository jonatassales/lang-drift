# Microruntimes contracts

The contracts of LangDrift Microruntimes, and only of Microruntimes. The LangDrift CLI and the planned server both depend on them; they depend on nothing but the Go standard library. Money formats belong to FinOps and artifact formats to the team's tracker; neither lives here. It is the Go surface of the `@langdrift/core` package, at `lib/core/contracts/go`.

The module is published by git tag, `lib/core/contracts/go/vX.Y.Z`, the form Go expects of a module in a subdirectory; `v0.2.0` is the first to carry the `judge` package. The CLI's `go.mod` requires a tagged version, so the CLI builds outside this repository too, and the root `go.work` resolves the module from source for local work. A `require` of a version that is not tagged makes Go fetch it, even inside the workspace, so a tag comes before the `require` that names it.

| Package | Contract |
| --- | --- |
| `configuration` | The `runtime.md` fence: exactly one `langdrift` JSON block, version 2, its fields and bounds, the component env file format, the derived image tag, and the lexical rules for identifiers, paths, ports, images and variable names. |
| `protocol` | What the connected agent and the server exchange at `<console-url>/api/mcp`: the MCP version, the five tools and their request and response types. |
| `judge` | What the planner and the platform's judge exchange at `<console-url>/api/judge` (spec 0045): the request, the response with usage and cost, and the typed answer of each kind (`sources`, `episodes`, `suggestion`, and `reserved-part` for the part a person writes under Safe Seat). `schema/` holds the JSON Schemas and one sample per kind; the console tests against the same files. `ParseAnswer` turns an answer's JSON into its typed answer and refuses anything outside the contract. |

```sh
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
```

## The fence

A fence is version 2. `Parse` refuses any other version, and any field this version does not define, with "re-plan the runtime": an older fence is never misread as a newer one.

A component is a reference to files: its name and its folder, with the Dockerfile inside the folder. The fence points at those files and does not repeat what they say.

```json
{
  "name": "demo-requests",
  "path": "apps/microservices/demo-requests",
  "dockerfile": "Dockerfile",
  "port": 3402,
  "readiness": { "kind": "tcp" },
  "envFile": "apps/microservices/demo-requests/runtime.env",
  "passEnv": ["MINIO_ACCESS_KEY", "MINIO_SECRET_KEY"]
}
```

| Field | Rule |
| --- | --- |
| `name` | Required. A DNS label, unique in the runtime; workloads and Services take it, so components address each other by it. |
| `path` | Required. The component folder: a clean relative path, unique, outside `.langdrift`; `.` is the project root. |
| `dockerfile` | A file name inside the folder, without `/`; `Dockerfile` when absent. |
| `port` | 1 to 65535; when absent, the Dockerfile's first `EXPOSE`. |
| `readiness` | `{ "kind": "tcp" }` or `{ "kind": "http", "path": "/…" }`; when absent, the Dockerfile's own `HEALTHCHECK`. |
| `envFile` | A committed file without secrets: a clean relative path inside the checkout, one `NAME=value` per line, uppercase names, plain-text values of at most 4 KiB. |
| `passEnv` | Uppercase names the agent reads from its own environment at start; none repeats, and none is also in the env file. |
| `id`, `kind` | Optional, as before. |

Readiness comes from the component: the planner reads the `HEALTHCHECK` of the component's own Dockerfile (one inherited from a base image is not seen) and renders it as an `exec` probe, the shell form as `/bin/sh -c`, the exec form as its array, with `--interval`, `--timeout`, `--retries` and `--start-period` as the probe's period, timeout, failure threshold and initial delay. A Dockerfile without its own `HEALTHCHECK`, or with `HEALTHCHECK NONE`, needs an explicit `tcp` or `http` readiness.

The image is never authored. `ImageTag` derives it from the runtime id, the component name and the revision, `langdrift-<runtimeId>-<name>:r<revision>`, and the CLI builds it with the fixed `docker build --file <path>/<dockerfile> --tag <image> .` from the project root. See ADR 0004.

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
