# 0226: `dft agent stop` abandons a waiting run

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** An operator whose run waits at a stage can let it go: `dft agent stop <runtime>` closes the open lap with the outcome `abandoned`, who and when, releases the reservation in the ledger, stops and cleans up the runtime through the execution ledger as the end of a run does, keeps the evidence and the memory, and leaves the runtime plannable: `dft agent plan` is accepted again. A run that is not waiting is refused with its state.

**Blocked by:** 0225.

**Status:** ready-for-agent

- [ ] After `dft agent stop` on a waiting run, the lap's outcome is `abandoned` with the person and the time, the ledger line shows the reservation released and nothing observed, the lifecycle's stop and cleanup ran for every owned resource, and the evidence folder is intact.
- [ ] `dft agent plan` on that runtime is accepted afterwards and raises the revision.
- [ ] `dft agent stop` on a runtime whose run is `completed`, `failed` or has no open lap is refused naming the state; on an unknown runtime, with "no such runtime".
- [ ] `dft agent status` shows the abandoned lap; the console's run record reads it as abandoned, not failed. (2026-11-05: `run.json` records the outcome `abandoned`; no console code reads `run.json` yet, and the console is outside this phase, so its reading lands with the console's run view.)
- [ ] The CLI README documents the command beside `confirm`.

## Dev part

**Category:** `core-domain`. What an abandoned lap settles to in the ledger is the money rule the owner must be able to explain during an incident.

**Paths:**
- `lib/cli/src/domain/finops/abandon.go`
- `lib/cli/src/domain/finops/abandon_test.go`

**Boundary:** `finops.Abandon(lap Lap, at time.Time, by string) Abandonment`, called by the stop command after the lifecycle stopped. `finops.Settlement` already names what `Settle` returns, so the result is `Abandonment` (released cents, observed cents, the overspend and the lap's closing line). The agent declares it with `Lap` (its reservation, its observed spend and whether a call's spend is unknown) and `LapClosing` (the outcome, the person and the RFC 3339 time, with `LapAbandoned`), one type per file as the rules require: `abandonment.go`, `lap.go` and `lap-closing.go` in `domain/finops`. The ledger's `LapLine` gains `closed`, the closing line.

**Goal:** Settle an abandoned lap: release the whole reservation when nothing was observed, release the difference when paid calls were observed before the wait, never release more than was reserved, keep the observed spend, and write the closing line with the outcome, the person and the time.

- [ ] A lap with a reservation and no observed spend releases it all; one with observed spend releases the remainder and keeps the observed figure; one with observed spend above the reservation releases nothing and keeps the overspend visible.
- [ ] The closing line carries `abandoned`, the person and the RFC 3339 time, and the ledger's totals add up after it.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `Abandon` and its test.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go test -race ./...`
- On this repository's microruntime: `dft agent run`, then `dft agent stop book-a-demo-microruntime`, then `dft agent plan` accepted.
