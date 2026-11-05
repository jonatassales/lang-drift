# 0331: `dft agent run` applies the authorized intent

**Spec:** .agents/workflow/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** A person types `dft agent run` alone, as the simulation does, and the run applies what the team authorized. No runtime is named, and no model, effort or confirmation is passed as a flag.

**Blocked by:** 0329, and the tickets of specs:
- 0050, which records intents;
- 0057, which answers the confirmations inside the run;
- 0064, which waits for the checks;
- 0065, which records a completed run.

**Status:** ready-for-agent

**Alignment:** change and remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] `dft agent run` takes no argument. Its first line is `# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team` in every variant.
- [ ] It applies the one recorded intent whose status is `authorized` and that has no completed run. None, or more than one, is a stop.
- [ ] A completed intent is not runnable again. A stopped run's intent stays runnable, and starts again from the gate.
- [ ] The `RUNTIME` argument leaves.
- [ ] These flags leave: `--model`, `--effort`, `--acknowledge-overspend`, `--confirm`, `--checks-wait` and `--by`.
- [ ] The golden test plays the run with no argument in every variant, and its pending list drops no line it had matched.

## Dev part

**Category:** `core-domain`. Which intent the run applies decides what is ever executed under the team's authorization.

**Paths:**
- `lib/cli/src/application/execution/applicable-intent.go`
- `lib/cli/src/application/execution/applicable-intent_test.go`

**Boundary:** `func applicableIntent(intents []RecordedIntent) (handoff.Intent, error)`. `RecordedIntent` carries an intent and whether it has a completed run. The agent declares it in its own file, and the run's start calls the function.

**Goal:** Pick the one authorized intent with no completed run, and refuse every other case with a reason a stop can print.

- [ ] One authorized intent without a completed run is applied.
- [ ] No intent, only proposed intents, or only completed ones are each refused.
- [ ] Two authorized intents without a completed run are refused.
- [ ] A stopped run's intent stays applicable.
