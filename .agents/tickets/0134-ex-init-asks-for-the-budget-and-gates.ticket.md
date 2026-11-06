# 0134: ex init asks for the budget and the default gates

**Spec:** `.agents/specs/0024-ex-init-asks-finops-and-generates-ci-gates.spec.md`

**What to build:** At a terminal, or with `--ask`, `ex init` asks for the current month's budget in dollars, the per-run cap in dollars, the reserve margin in percent, and whether to apply the default gates. It records the answers through the same use cases as `ex finops budget set` and `ex finops gates set --default`, under one `--by`.
- An empty first budget answer skips the budget.
- Flags answer the questions without a terminal: `--budget`, `--per-run` and `--reserve-margin`, which go together, and `--default-gates`.
- A second run asks only for what is not yet recorded.
- With `--json`, or with no terminal and no flags, `init` asks nothing and behaves as today.
- The result lists what was recorded and the commands still to run.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Piped answers with `--ask` record a budget whose mirror holds the current month, the amounts in cents and the `--by` value, and the default gates.
- [x] An amount with more than two decimals, a negative amount or text is refused before anything is recorded. A table test covers the dollars-to-cents conversion.
- [x] An empty first budget answer skips the per-run and reserve-margin questions and records no budget.
- [x] The budget flags without `--ask` and without a terminal record the budget without a prompt. One budget flag without the other two is a usage error.
- [x] A second `init` asks nothing for a budget or gates already recorded, and changes neither.
- [x] The existing identity-only init test passes unchanged, and `init --json` asks nothing.
- [x] The result carries `recorded` and `next` (`ex harness envelope set`, `ex harness model set`, `ex harness driver set`). Both are omitted when empty.
- [x] The README documents the questions and the flags.
