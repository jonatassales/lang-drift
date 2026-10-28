# Exphere

The engine of LangDrift Microruntimes. This workspace holds everything that plans, runs and records a microruntime, and nothing that belongs to another product: Vision's artifacts and FinOps's money formats stay outside, and Exphere only reads them.

| Member | Package | Holds |
| --- | --- | --- |
| [`cli/`](cli/README.md) | `@langdrift/ex` | The `ex` CLI and the connected agent: project setup, the bounded local lifecycle on kind, Terraform and Kubernetes. |
| [`contracts/`](contracts/README.md) | `@langdrift/exphere-contracts` | The Microruntimes contracts: the `runtime.md` fence, the agent-server protocol, and the formats of a run's records. |

A server that answers the agent's protocol and holds each runtime's intent is planned as the next member. Until it exists, nothing in this repository can drive `ex agent`.

## Boundary

Exphere writes only `.langdrift/microruntimes/` in a project. It reads `.langdrift/harness/`, `.langdrift/finops/` and the team's artifacts, and asks a FinOps gate before any paid work. Every plan with Microruntimes includes FinOps.

## Go workspace

The root `go.work` ties `cli`, `contracts` and `apps/extub` together, so the CLI imports the contracts without a `replace` directive. Building the CLI outside this repository needs the contracts module published first.

```sh
pnpm turbo run build test typecheck --filter=@langdrift/ex --filter=@langdrift/exphere-contracts
```
