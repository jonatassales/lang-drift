# Exphere and the `.langdrift` folder

**Status:** ready-for-agent

## Problem Statement

The Microruntimes page shows a simulated session: `ex agent plan` reads a ticket and its documents, recommends a model and an effort, sizes the components inside an envelope, and `ex agent run` plays a gated, human-in-the-loop workflow to a pull request. That session is what the console of LangDrift Microruntimes has to show for real, live, and what sells the product. The real `ex` reproduces only part of it (research note `ex-microruntimes-capability.research.md`): the lifecycle on kind, Terraform and Kubernetes exists; planning from a ticket, retrieval, pricing, gates, workflow stages and delivery do not, and nothing in this repository serves the protocol the agent speaks. The CLI also sits alone in `packages/ex`, with its contracts mixed into its own domain, and nothing defines the files a run leaves behind for a console to render.

## Solution

Exphere becomes a workspace, the engine of Microruntimes and of nothing else: `exphere/cli` holds the `ex` CLI, `exphere/contracts` the Microruntimes contracts, and a server that answers the agent follows. A project's `.langdrift/` folder is versioned, self-contained and split by product, and a run leaves records a console renders without reading prose. This change moves the CLI, extracts the contracts without changing behavior, and commits an illustrative `.langdrift/` that covers every field of the four cards of the page.

## Decisions

Taken with the owner on 2026-10-28.

- **Exphere is the Microruntimes engine.** A workspace with the CLI, the contracts and, next, the server. Its contracts cover Microruntimes only; Vision's artifacts and FinOps's money formats stay outside. `apps/extub` is the console's page stub and unrelated.
- **Microruntimes ships with FinOps.** Its Play makes paid model calls in parallel, and the platform's principles require a team budget and a deterministic gate before paid work. Every plan with Microruntimes includes FinOps; FinOps sells without Microruntimes. The website keeps Microruntimes's own page and does not change.
- **The console shows what the page simulates.** The same cards, fed by real data and updated live; Play from the console or from the CLI makes the same request, and both show the same run.
- **Runs in parallel.** The connected agent already supervises up to 128 runtimes, one Go worker each.
- **`.langdrift/` is versioned and self-contained.** It reads the project and writes only inside itself; only `agent/state/` is ignored. General folders: `config.json`, `agent/`, `harness/`. Product folders: `vision/`, `microruntimes/`, `finops/`, each written only by its product.
- **The racing metaphor.** A microruntime is the car, a run is a race, a lap is one pass of it: an attempt paused at a gate, an execution, a load step. A run's money is one FinOps ledger per run, one line per lap; there is no dated report.
- **FinOps is money only.** Bounded budgets, gates that block paid work, dated prices, token forecasts, the ledger of what was reserved, spent, recovered and overspent. Hardware limits in units stay in `harness/envelope.json`.
- **Memory is the run's Markdown; the index is derived.** A run's memory file and laps carry JSON fences; the RAG index is rebuilt from them into `agent/state/rag/`. No vector database in `ex init`: retrieval starts structured, and embeddings come only when that measurably misses runs.
- **The workflow is the team's.** `harness/workflow.json` describes the stages and where a person takes part; LangDrift reads it and never prescribes it.

## The `.langdrift` folder

```text
.langdrift/
  config.json                         project identity
  .gitignore                          agent/state/
  harness/                            envelope.json, workflow.json, artifacts.json, guardrails/
  vision/artifacts/                   prds/, designs/, specs/, tickets/
  microruntimes/
    AGENTS.md                         execution rules, each marked code, planned or guidance
    runtimes/<id>/runtime.md          the car: one langdrift fence
    runs/NNNN-<slug>/                 the race: plan.json, run.json, NNNN-<slug>.run.md, laps/, evidence/
  finops/                             budgets.json, gates.json, forecasts.json, prices/, ledger/
```

The formats of a run's records are described in `exphere/contracts/README.md`.

## Console coverage

Every field of the page's four cards has a source in the sample; none is read from prose except a ticket's title and status line.

| Card | Fields | Source |
| --- | --- | --- |
| Ticket | number, slug, title, status, linked artifacts by role | the ticket; `plan.json` › `ticket`, `context.sources` |
| Analysis | documents read and their facts, guardrail criteria, provenance of each bound, runs considered, sizing basis, suggestion | `plan.json` › `documents`, `bounds`, `retrieval`, `sizing`, `suggestion` |
| Plan | proposed totals and state, bounded and headroom, components with technologies and shares, checks, materialization, price basis, the `runtime.md` | `plan.json` › `totals`, `components`, `checks`, `materialization`, `price`, `runtime`; `harness/envelope.json` |
| Workflow | stages and where a person takes part, model and price, effort, estimate and reservation, laps, stages, delivery | `harness/workflow.json`; `run.json`; `finops/prices/`, `finops/forecasts.json`, the run's ledger |

## Implementation Decisions

- `packages/ex` moves to `exphere/cli`; the Go module becomes `github.com/jonatassales/tao/exphere/cli` and the package stays `@langdrift/ex`. `pnpm-workspace.yaml` adds `exphere/*`.
- `exphere/contracts`, package `@langdrift/exphere-contracts`, is a Go module with no dependency beyond the standard library: `configuration` holds the `runtime.md` fence and the lexical rules it uses; `protocol` holds the MCP version, the tool names and the wire types.
- The CLI keeps its names through type aliases and delegation, so no behavior changes: `domain/connection` aliases the fence, `application/connected` aliases the wire types, `domain/microruntime` delegates its identifier, path, port, platform-path and image rules and shares the error values, and the MCP client uses the protocol constants. The agent state converts to the wire type at the one place it is sent.
- A root `go.work` lists `apps/extub`, `exphere/cli` and `exphere/contracts`, so the CLI imports the contracts without a `replace` that climbs with `../`. The CI action reads the Go version from `go.work`.
- The `.langdrift/` sample holds data only. Field descriptions live in this spec and in the contracts README; the sample's hashes are byte hashes, checked against the files, and its `runtime.md` passes the CLI's parser.

## Testing Decisions

- The CLI's full suite, with the race detector, passes unchanged after the move and the extraction: it is the evidence that behavior did not change.
- The fence's tests move with it to the contracts module, and the lexical rules gain table tests there, where they are owned.
- `pnpm turbo run build test typecheck` for the three Go packages, `pnpm turbo run check`, and `drift check dead-code`.

## Out of Scope

Follow-up work, each for its own ticket once grilled:

- The Exphere server: the protocol tools, each runtime's intent, credentials, and the operator surface the CLI and the console use to request a run.
- Records written by the CLI: `plan.json`, `run.json`, the memory, laps and evidence inside the run folder, and the CLI layout moving `runtimes/` under `microruntimes/` and private state under `agent/state/`.
- Progress events in the protocol, so the console shows the cards live.
- Contract changes: canonical hashing of JSON records, per-component resources and technologies in the fence, a run-to-completion kind, and whether the fence becomes `runtime.json`.
- Planning from a ticket (`ex agent plan`): resolving artifacts by role through `harness/artifacts.json`, structured retrieval over earlier runs, and the suggestion.
- The executor of the workflow stages, which orchestrates a coding agent, passes the money gates on each paid call and waits for a person's confirmation from the console.
- The FinOps gate contract: a reservation request and its decision, held atomically by a server, never in a committed file.
- Limits for parallel runs: a project-wide ceiling, run identity in runtime names and image tags, and lap records of what else was running.
- The console image and microruntime of `apps/console` with `apps/extub`.

## Further Notes

Open decisions for the owner: how runs are numbered without collisions across branches (a server-issued number, or an independent id); who writes the general parts of `.langdrift/` now that Exphere writes only `microruntimes/`; what of a run leaves the machine for the console (facts and hashes, or documents); and the embedding provider, if retrieval ever needs one.
