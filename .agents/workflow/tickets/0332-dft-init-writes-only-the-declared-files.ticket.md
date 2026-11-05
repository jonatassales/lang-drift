# 0332: `dft init` writes only the declared files, and the input commands leave

**Spec:** .agents/workflow/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** `dft init` silently creates `.langdrift/` with only what plan and run need, asks nothing and prints nothing. The envelope, the budget and the choices come from the team's documents and the recorded intent, so the commands that set them as inputs leave.

**Blocked by:** the tickets of specs:
- 0053: the templates and the design doc's fields;
- 0054: the guardrail, and the envelope from the design doc;
- 0056: the price list, the forecasts, and the budget from the design doc;
- 0057: the workflow file;
- 0058: the driver the team writes;
- 0059: the skills.

**Status:** ready-for-agent

**Alignment:** change and remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] `dft init` on an empty folder writes exactly the files spec 0048's Produces table declares, and nothing else:
  - `.langdrift/skills/`;
  - `.langdrift/finops/prices.json` and `.langdrift/finops/forecasts.json`;
  - `.langdrift/harness/workflow.json`;
  - `.langdrift/harness/guardrails/0001-finops.gr.md`;
  - `.langdrift/templates/`, with `ticket.md`, `prd.md`, `dd.md` and `spec.md`.

  It prints nothing and exits 0.
- [ ] An existing file is kept byte for byte, and a second `dft init` changes nothing.
- [ ] `dft init` doesn't write `harness/driver.json`.
- [ ] Its questions and flags leave: `--ask`, `--budget`, `--per-run`, `--reserve-margin`, `--default-gates`, `--artifacts`, `--ci-gates` and `--by`, with the month's budget, the per-run cap, the reserve margin, the default gates, the artifacts folder and the CI Gates workflow.
- [ ] These commands leave: `harness envelope set`, `workflow set`, `model set` and `driver set`, and `finops budget set`, `gates set`, `prices set` and `forecasts set`, with their `--by`, `--status`, `--owner` and `--file` flags.
- [ ] `gates report` leaves with the CI Gates workflow that called it, with the gates report and its history.
- [ ] What `dft init` wrote beyond the table leaves, with the code only it used:
  - the project identity in `config.json`;
  - the input records and their lock;
  - the mirrors;
  - `microruntimes/AGENTS.md`;
  - the CI Gates workflow.
- [ ] The golden test runs `dft init` on its copy, and asserts that it writes exactly the declared files and prints nothing.

## Dev part

**Override:** No part qualifies: the files' contents belong to specs 0053 to 0059, and what remains is scaffolding and removal (policy, Categories: never boilerplate, configuration or a migration). 2026-11-05.
