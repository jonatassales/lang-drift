# .langdrift

The LangDrift folder of this project, versioned with it except `agent/state/`. It is self-contained: it reads the project and writes only inside itself, and a console renders every screen from it.

> **Illustrative sample.** One run of ticket `0070-console-microruntime` from this repository's own tracker, planned and executed as the microruntime `console-local` in two laps. Nothing here came from a real execution, and the CLI does not read this layout yet. Spec `.agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md` records the decisions behind it.

| Path | Scope | Written by |
| --- | --- | --- |
| `config.json` | general | `ex init`: the project identity |
| `agent/state/` | general, ignored | the connected agent: locks, cluster record, build copies, outbox, the derived RAG index |
| `harness/` | general | the server, as mirrors of console input: the hardware envelope, the delivery stages; `guardrails/` holds the team's policy text |
| `microruntimes/` | Microruntimes | The engine: execution rules, one `runtime.md` per microruntime, one folder per run |
| `finops/` | FinOps | the server, as mirrors: budgets, gates, prices, forecasts, one ledger per run |

Every value that can block or allow a run is console or CLI input held by the server. The JSON files here mirror it for reading; no gate reads them, and an edit to one is shown as an unauthorized change.

A microruntime is the car, a run is a race, a lap is one pass of it.
