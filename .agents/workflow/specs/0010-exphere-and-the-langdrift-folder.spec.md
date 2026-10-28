# Exphere and the `.langdrift` folder

**Status:** ready-for-agent

## Problem Statement

The Microruntimes page shows a simulated session: `ex agent plan` reads a ticket and its documents, recommends a model and an effort, sizes the components inside an envelope, and `ex agent run` plays a gated, human-in-the-loop workflow to a pull request. That session is what the console has to show for real, live, and what sells the product. The real `ex` reproduces only part of it (research note `ex-microruntimes-capability.research.md`): the lifecycle on kind, Terraform and Kubernetes exists; planning from a ticket, retrieval, pricing, gates, workflow stages and delivery do not, and nothing in this repository serves the protocol the agent speaks. The CLI also sits alone in `packages/ex`, with its contracts mixed into its own domain, and nothing defines the files a run leaves behind for a console to render.

## Solution

LangDrift becomes a FinOps project with Microruntimes as its engine. Exphere becomes a workspace, the engine of Microruntimes and of nothing else: `exphere/cli` holds the `ex` CLI, `exphere/contracts` the Microruntimes contracts, and a server that answers the agent follows. A project's `.langdrift/` folder is versioned and self-contained, and a run leaves records a console renders without reading prose. Every value that can block or allow a run comes from explicit input in the console or the CLI, never from a file. This change moves the CLI, extracts the contracts without changing behavior, and commits an illustrative `.langdrift/` whose one run covers every field of the four cards of the page.

## Decisions

Taken with the owner on 2026-10-28.

- **The project pivots to FinOps.** Its value comes from data the platform produces itself: runs, laps, ledgers. Vision stops as a product: it would need ever more of the client's data, and producing artifacts for it would break the other products. It may return as a connector, an MCP to tools such as GitHub, Linear or Obsidian, that brings their data into the console. The website does not change for now; the brand and the site are a separate decision.
- **Exphere is the Microruntimes engine.** A workspace with the CLI, the contracts and, next, the server. Its contracts cover Microruntimes only; money formats belong to FinOps. `apps/extub` is the console's page stub and unrelated.
- **Microruntimes ships with FinOps.** Its Play makes paid model calls in parallel, and the platform's principles require a team budget and a deterministic gate before paid work. Every plan with Microruntimes includes FinOps; FinOps sells without Microruntimes. Microruntimes keeps its own page.
- **The console is the primary path.** A person creates a microruntime in the console, sets its budget and limits, and plays it; the `ex` CLI makes the same request for automation. Both show the same run, live. The agent already supervises up to 128 runtimes in parallel, one Go worker each.
- **Guardrails come only from input.** Hardware limits, budgets, gates, the price and forecast sources, model and effort are explicit input in the console or the CLI, recorded with who, when and the previous value. The Exphere server is their source of truth. Nothing collected and nothing in a file sets one: a person or an agent could write a file with a huge limit. Files in `.langdrift/` are read-only mirrors; a mirror that differs from the server is shown as an unauthorized change.
- **The client's harness first.** LangDrift prescribes no workflow and no artifact. The person names the ticket in the console or the CLI, and the planner follows its links recursively, whatever file types they reach.
- **Collection is agentic inside a deterministic shell.** The traversal, its limits, the refusal of secrets and paths outside the project, and the hashes are deterministic. An agent decides what is input for the plan and what is discarded, with a reason and a confidence, and every node it considered stays in the record. Collected text is context only: it can inform a suggestion or flag a discrepancy, never set a value. The agent's call is paid work under the FinOps gates.
- **The racing metaphor.** A microruntime is the car, a run is a race, a lap is one pass of it: an attempt paused at a gate, an execution, a load step. A run's money is one FinOps ledger per run, one line per lap; there is no dated report.
- **Drift survives as FinOps metrics,** measured from the platform's own records: plan drift (proposed against executed), intent drift (a run planned on a source that has since changed, by its hash), cost drift (spent against reserved and budgeted), and decisions against the suggestion.
- **Memory is the run's Markdown; the index is derived.** A run's memory file and laps carry JSON fences; the RAG index is rebuilt from them into `agent/state/rag/`. Retrieval starts structured; embeddings come only when that measurably misses runs.

## The `.langdrift` folder

```text
.langdrift/
  config.json                         project identity
  .gitignore                          agent/state/
  harness/                            mirrors: envelope.json, workflow.json; guardrails/ the team's policy text
  microruntimes/
    AGENTS.md                         execution rules, each marked code, planned or guidance
    runtimes/<id>/runtime.md          the car: one langdrift fence
    runs/NNNN-<slug>/                 the race: plan.json, run.json, NNNN-<slug>.run.md, laps/, evidence/
  finops/                             mirrors: budgets.json, gates.json, forecasts.json, prices/, ledger/
```

The formats of a run's records are described in `exphere/contracts/README.md`.

## Console coverage

Every field of the page's four cards has a source in the sample; none is read from prose except a ticket's title and status line.

| Card | Fields | Source |
| --- | --- | --- |
| Ticket | number, slug, title, status, what it links | the ticket in the client's tracker; `plan.json` › `ticket`, `sources` |
| Analysis | sources reached and kept or discarded, with reason and confidence; facts by topic with citations; runs considered; sizing basis; suggestion | `plan.json` › `sources`, `facts`, `retrieval`, `sizing`, `suggestion` |
| Plan | proposed totals and state, limits set in the console and their records, components with technologies and shares, checks, materialization, price basis, the `runtime.md` | `plan.json` › `inputs`, `totals`, `components`, `checks`, `materialization`, `price`, `runtime` |
| Workflow | stages and where a person takes part, model and price, effort, estimate and reservation, laps, stages, delivery | `harness/workflow.json`; `run.json`; `finops/prices/`, `finops/forecasts.json`, the run's ledger |

## Implementation Decisions

- `packages/ex` moves to `exphere/cli`; the Go module becomes `github.com/jonatassales/tao/exphere/cli` and the package stays `@langdrift/ex`. `pnpm-workspace.yaml` adds `exphere/*`.
- `exphere/contracts`, package `@langdrift/exphere-contracts`, is a Go module with no dependency beyond the standard library: `configuration` holds the `runtime.md` fence and the lexical rules it uses; `protocol` holds the MCP version, the tool names and the wire types.
- The CLI keeps its names through type aliases and delegation, so no behavior changes: `domain/connection` aliases the fence, `application/connected` aliases the wire types, `domain/microruntime` delegates its identifier, path, port, platform-path and image rules and shares the error values, and the MCP client uses the protocol constants. The agent state converts to the wire type at the one place it is sent.
- A root `go.work` lists `apps/extub`, `exphere/cli` and `exphere/contracts`, so the CLI imports the contracts without a `replace` that climbs with `../`. The CI action reads the Go version from `go.work`.
- The `.langdrift/` sample holds data only. Its run consumes ticket 0070 of this repository's own tracker, a real follow-up, and the files it links. Every JSON file in `harness/` and `finops/` declares itself a mirror of a server record. Field descriptions live in this spec and in the contracts README; the sample's hashes are byte hashes, checked against the files, and its `runtime.md` passes the CLI's parser.

## Testing Decisions

- The CLI's full suite, with the race detector, passes unchanged after the move and the extraction: it is the evidence that behavior did not change.
- The fence's tests move with it to the contracts module, and the lexical rules gain table tests there, where they are owned.
- `pnpm turbo run build test typecheck` for the three Go packages, `pnpm turbo run check`, and `drift check dead-code`.

## Out of Scope

Follow-up work, each for its own ticket once grilled:

- The Exphere server: the protocol tools, each runtime's intent, credentials, the record of every guardrail input, and the operator surface the console and the CLI use to request a run.
- Planning from a ticket (`ex agent plan`): the recursive collection with its limits, the agent that keeps or discards each source, structured retrieval over earlier runs, and the suggestion.
- Records written by the CLI and the server: `plan.json`, `run.json`, the memory, laps and evidence inside the run folder, mirrors written on every change of a guardrail, and the CLI layout moving `runtimes/` under `microruntimes/` and private state under `agent/state/`.
- Progress events in the protocol, so the console shows the cards live.
- Contract changes: canonical hashing of JSON records, per-component resources and technologies in the fence, a run-to-completion kind, and whether the fence becomes `runtime.json`.
- The executor of the workflow stages, which orchestrates a coding agent, passes the money gates on each paid call and waits for a person's confirmation from the console.
- The FinOps gate contract: a reservation request and its decision, held atomically by the server, never in a committed file.
- Limits for parallel runs: a project-wide ceiling, run identity in runtime names and image tags, and lap records of what else was running.
- The console image and microruntime of `apps/console` with `apps/extub`: ticket 0070.

## Further Notes

Open decisions for the owner: whether only an owner in the console may raise a cap while the CLI may use and lower limits (recommended); how runs are numbered without collisions across branches; who writes the general parts of `.langdrift/`; what of a run leaves the machine for the console; the website and the brand after the pivot; and the embedding provider, if retrieval ever needs one.
