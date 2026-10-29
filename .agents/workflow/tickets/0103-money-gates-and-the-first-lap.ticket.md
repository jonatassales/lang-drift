# 0103: `ex agent run` estimates, reserves and pauses at the gate

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** `ex agent run <runtime>` prints the run session's opening for real: running it authorizes the plan, the intent prints from the fence or the command line, the estimate at list prices from the token forecast, the per-run budget, and the gate: reserved and opened, or a reservation above the per-run cap, the ticket's allocation or the period's remainder paused with nothing paid and the decision the owner must take. Every attempt is a lap: a numbered lap file with its steps table and fence, a line of the runtime's ledger with the gate results, and the run record's lap entry. An unacknowledged overspend of an earlier lap requires the owner's acknowledgement; a price list that does not cover the model or a ticket without an allocation blocks.

**Blocked by:** 0100.

**Status:** ready-for-agent

- [ ] Against the fixture runtime, a run with Opus at high prints lines 1 to 6 of the pause variant and leaves lap `0001` and the ledger line with the sample's figures (3400, 4080, nothing paid), exit code 3.
- [ ] A run with the suggestion prints the opened gate with 950 reserved as 1140 and leaves lap `0002` open, with the run record's lap entry and the ledger line.
- [ ] The gate arithmetic is unit-tested over the website's three-by-three model and effort matrix, including the rounding of the Haiku figures.
- [ ] `ex agent run` is the hero's command; the connected agent's loop answers to `ex agent supervise` with its behavior unchanged.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; the two runs in a scratch copy of this repository.
