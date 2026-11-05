# 0094: Private state moves under `agent/`, runtimes under `microruntimes/`, evidence into the runtime folder

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** every file the CLI writes lands where the sample's layout says: the private state (locks, cluster record, kubeconfig, execution ledgers, Terraform state, build logs, connection records, migrations, build copies) under `.langdrift/agent/state/`, the runtime credentials there too instead of a per-runtime `.env` with its own ignore file, the authored and materialized runtime files under `.langdrift/microruntimes/runtimes/<id>/`, and the evidence of every operation and inspection inside that runtime's `evidence/` keyed by execution, instead of `.audits/`. An execution the connected agent starts without a plan gets a runtime folder with an execution-kind record on first start, so inspections have a home. The cluster deletion check tolerates a materialized runtime that never started. Every existing test follows the new paths.

**Blocked by:** None (can start immediately; parallel with 0093).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] After `ex init`, `ex cluster create` and a connected start in a disposable project, no file is written under `.langdrift/state`, `.langdrift/runtimes` or `.audits`; the paths above hold them.
- [ ] An inspection with no planned run creates `microruntimes/runtimes/<id>/run.json` of kind `execution` and writes its evidence under that folder.
- [ ] `ex cluster delete` proceeds past a runtime directory that holds only materialization intent and refuses while any execution ledger owns resources.
- [ ] The full suite passes under the race detector.

### Verification

- `GOTOOLCHAIN=local go test -race ./...` and `go vet ./...` in the CLI package.
