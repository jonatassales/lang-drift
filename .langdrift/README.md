# .langdrift

The LangDrift folder of this project, versioned with it except `agent/state/`. It is self-contained: it reads the project and writes only inside itself, and a console renders every screen from it.

> **Illustrative sample.** One run of the ticket `0001-TICKET-console-microruntime`, planned and executed as the microruntime `console-local` in two laps. Nothing here came from a real execution, and the CLI does not read this layout yet. Spec `.agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md` records the decisions behind it.

| Path | Scope | Written by |
| --- | --- | --- |
| `config.json` | general | `ex init`: the project identity |
| `agent/state/` | general, ignored | the connected agent: locks, cluster record, build copies, outbox, the derived RAG index |
| `harness/` | general | the team: guardrails, the hardware envelope, the delivery stages, where artifacts live |
| `vision/` | Vision | Vision: the artifacts a ticket links |
| `microruntimes/` | Microruntimes | Exphere: execution rules, one `runtime.md` per microruntime, one folder per run |
| `finops/` | FinOps | FinOps: budgets, gates, prices, forecasts, one ledger per run |

A microruntime is the car, a run is a race, a lap is one pass of it.
