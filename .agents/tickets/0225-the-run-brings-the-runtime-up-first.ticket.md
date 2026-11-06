# 0225: the run brings the runtime up before the first stage

**Spec:** .agents/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** An operator runs a planned runtime and sees it come up before anyone is asked anything: `ex agent run` opens the gate and reserves, then creates the cluster, builds the images, claims an execution, starts, waits for readiness and prints the runtime-ready line; only then it plays the stages, and a stage that waits for a person pauses with the runtime running. The end of the run observes delivery, stops and cleans up by the ledger, releases the unused reservation, settles the lap and writes the memory, as today, and brings nothing up. The bring-up has its own timeout, long enough for four image builds, named in the README.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] In the composed test with the fake lifecycle, the transcript prints the runtime-ready line before the first `waiting` line, and the written tree holds the execution's evidence before any stage record.
- [ ] A run whose every stage waits for a person ends its first lap `waiting` at the first stage with the execution started and the ledger's reservation held; `ex agent status` shows the runtime up and the stage waiting.
- [ ] A bring-up that fails (build, readiness, cluster) ends the lap with the failure, releases the reservation and leaves the evidence; no stage is played.
- [ ] The end of the run no longer builds or starts anything; the delivery observation, stop, cleanup, settlement and memory are unchanged and their tests still pass.
- [ ] The CLI README describes the order and the bring-up timeout.

## Dev part

**Category:** `hard-part`. The order of a run's phases, and what each failure leaves behind, is where this ticket's difficulty sits.

**Paths:**
- `lib/cli/src/application/execution/next-phase.go` (named `run-phases.go` when written; renamed on 2026-11-05 after the file-naming rule)
- `lib/cli/src/application/execution/next-phase_test.go`

**Boundary:** `execution.NextPhase(run RunState) Phase`, called by the run loop after each phase, where `RunState` (gate, bring-up, stage and finish outcomes so far) and `Phase` (`reserve`, `bringUp`, `play`, `wait`, `finish`, `fail`, `done`) are declared in `phases.go` beside it by the agent, which also writes the loop that executes each phase through the existing ports.

**Goal:** Decide the next phase from what has happened: reserve before bring-up, bring-up before any stage, wait at a human stage with the runtime up, finish after the last stage, fail on a bring-up or stage failure with the reservation released, and never play a stage without a started execution.

- [ ] Every transition of the table is covered by one test, including a bring-up failure, a stage failure, a human wait and a confirmed resume.
- [ ] No sequence of outcomes reaches `play` without `bringUp` having succeeded.

**Override:** The owner does not write this part. On 2026-11-05 the owner said "I won't write anything", so the agent writes these paths. The category, boundary, goal and criteria above stand.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go test -race ./...`
- `lib/cli/dist/ex agent run <runtime> --timeout 30m` on this repository's microruntime, then `kubectl --kubeconfig .langdrift/agent/state/cluster.kubeconfig get pods -n <runtime>` showing four Ready pods while the run waits at `planning`.
