# 0321: the hand-off records

**Spec:** .agents/workflow/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** What the console receives has a typed shape in the Core contracts, replacing the protocol's channel: the plan's card data (`plan.json`), the recorded intent (`intent.json`) and the run's record (`run.json`). Every figure is in its unit, integers for millicores, MiB, cents and tokens and two decimals for similarities, and a value not known is `null`, never zero. The fields inside each section join as the element specs that fill them land; this ticket gives the records their sections and the intent its full shape.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The plan's card data has the sections `ticket`, `analysis.documents`, `analysis.episodes`, `analysis.suggestion`, `plan.components`, `plan.sizing`, `plan.price`, `plan.authorization` and `plan.evidence` (`file`, `state`, `revision`, `proposedTo`, for EVIDENCE `runtime.md draft · revision 1 · proposed to the console`) and `workflow`; a section absent from the record has not been revealed yet, as a card reads `Not set yet`.
- [ ] The recorded intent holds `ticket`, `runtimeId`, `revision`, `model`, `effort`, `human` (a boolean for each of `planning`, `spec`, `tickets`, `implement` and `review`) and `status` (`proposed` or `authorized`), for `console › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1` and the HUMAN IN THE LOOP ticks.
- [ ] The run's record has the sections for the applied choice, each stage's state, the effort record, the dev ticket, the stage reports, the pull request and its checks, the runtime's state and the run's id and outcome, so that the Workflow card's RESERVED `$11.40`, OBSERVED `$9.80`, RELEASED `$1.60`, PULL REQUEST `#19 · 3f9c2e1 · ready for review` and SAFE SEAT `met` have a place.
- [ ] Unknown values decode and encode as `null`, and a record with an unknown field is refused.
- [ ] Go types only; no JSON Schema is published for the hand-off.

**Note (2026-11-05):** specs 0050 and 0065 keep RESERVED, OBSERVED and RELEASED out of the run's record: they come from the run's entry in the money ledger, `finops/ledger/<run>.ledger.json` (spec 0056), which the record's `run` id names. The record's sections take spec 0065's keys: `run`, `intent`, `stages` (each with its `report`, the dev ticket among them), `effort`, `safeSeat`, `runtime`, `delivery` and `outcome`.

## Dev part

**Category:** `core-domain`. A recorded intent is what the run applies as the team's authorization; which intents are well formed decides what can ever be run.

**Paths:**
- `lib/core/contracts/go/handoff/validate-intent.go`
- `lib/core/contracts/go/handoff/validate-intent_test.go`

**Boundary:** `func (intent Intent) Validate() []Problem`, where `Intent`, its `Human` checkpoints and `Problem` are declared by the agent in their own files.

**Goal:** Accept a recorded intent only when its ticket and runtime id are present, its revision is at least 1, its effort is `low`, `medium` or `high`, its model is not empty, `planning` is always a human checkpoint (the simulation's Planning has no unchecked state), and its status is `proposed` or `authorized`; return every problem.

- [ ] The simulation's default intent (Sonnet 5.5, medium, Planning and Review checked, authorized) passes.
- [ ] Planning unchecked, an unknown status and a revision of 0 are each a problem, and an intent with all three returns three.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `validate-intent.go` and its test.
