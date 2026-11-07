# 0333: `agent confirm` and `agent connect` leave, and `dft` has three commands

**Spec:** .agents/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** The last commands beyond the three leave. `dft` answers to `dft init`, `dft agent plan <ticket id>` and `dft agent run`, and nothing else, and its README says so.

**Blocked by:** 0330, 0331, 0332, and the tickets of specs:
- 0051: the judge needs no console connection;
- 0057: the confirmations are answered inside the run.

**Status:** ready-for-agent

**Alignment:** remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] `agent confirm` leaves. The run's standard input answers `console › waiting    planning · the owner confirms` with `✓ confirmed  planning · by the owner` (spec 0057).
- [ ] `agent connect` leaves, with the connection store, the credentials and the private state only it used.
- [ ] `dft` answers to exactly `init`, `agent plan` and `agent run`. Every other command and flag is refused as unknown.
- [ ] `lib/cli/README.md` describes the three commands, the line format and the stop rule, and nothing that was removed.
- [ ] The golden test seeds nothing through a removed command.

## Dev part

**Override:** No part qualifies: removal and documentation, with no domain decision of its own (policy, Categories: never a migration or configuration). 2026-11-05.
