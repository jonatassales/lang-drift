# 0328: the commands the simulation doesn't show leave

**Spec:** .agents/workflow/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** `dft` no longer answers to the commands and flags that nothing in the simulation shows and that `dft agent plan` and `dft agent run` do not need. Each is refused as unknown, and the code only it used is gone. The commands plan and run still depend on stay until their own tickets (0330 to 0333), and so does `gates report`, which the CI Gates workflow that `dft init --ci-gates` writes still calls until 0332.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] These commands are gone and refused as unknown:
  - `doctor`, `agent status`, `agent reindex`, `version` and `help`;
  - the generated `completion` command;
  - the `runtime` group and its retirement message;
  - the `cluster` group (`create`, `delete`, `inspect`, `load-image`);
  - `agent stop`, `agent migrate`, `agent project` and `agent supervise`.
- [ ] These global flags are gone: `--json` with its NDJSON records, `--project`, `--timeout`, `--version` and `-V`.
- [ ] `dft agent plan` and `dft agent run` find the project root as the nearest folder, from the working directory up, that holds `.langdrift/`. None is a stop. `dft init` uses the working directory.
- [ ] The code only these commands used is gone:
  - the doctor;
  - the runtime catalogs;
  - the cluster's delete and inspect;
  - the connected loop's run, project and migration;
  - the console's MCP client;
  - `agent stop`'s abandonment.
- [ ] Kept, because the run or the plan still uses it:
  - the cluster's create and load, and the runtime's start, stop and inspect, inside the run;
  - the judge's connection read, inside the plan.
- [ ] The tests that read state through `agent status` read the records instead. The tests that read kinds through `--json` read the injected session sink, or the records. The tests that ran `agent reindex` build the index through the RAG service.
- [ ] `dft agent plan` and `dft agent run` behave as before, and their tests pass.

## Dev part

**Override:** No part qualifies: this ticket deletes commands, flags and the code only they use, with no domain decision of its own (policy, Categories: never a migration or a rename). 2026-11-05.
