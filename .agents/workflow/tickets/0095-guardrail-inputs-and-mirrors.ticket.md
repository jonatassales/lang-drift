# 0095: Guardrail inputs are CLI input, held privately and mirrored

**Spec:** .agents/workflow/specs/0015-ex-performs-the-microruntimes-session.spec.md

**What to build:** the owner records every value that can block or allow a run through the CLI: `ex harness envelope set`, `ex harness workflow set`, `ex harness model set` (the default model and effort, and the planner's model) and `ex harness driver set` (the coding-agent command), and `ex finops budget set`, `ex finops gates set` (the sample's six gates by default, or a file), `ex finops prices set` and `ex finops forecasts set`. Each input lands in a private store under `agent/state/` with who set it, when and the previous value, locked against concurrent writes, and is mirrored to `harness/` or `finops/` in the sample's shape, naming the CLI as its source and the record that set it. Invalid input is refused as usage; a mirror is never read back by the CLI.

**Blocked by:** 0092, 0094.

**Status:** ready-for-agent

- [ ] Running the eight commands with the sample's values produces mirrors equal to the sample's `harness/` and `finops/` files except for the source fields, the record ids and the timestamps, which are injected in the test.
- [ ] Setting a value twice keeps the previous value in the store and the mirror shows only the latest.
- [ ] A workflow with a duplicate stage id, an unknown human marker or a default human stage that is not a stage is refused with exit code 2; a price list in another unit is refused.
- [ ] The budget mirror computes the unallocated cents from the allocations.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; the eight commands in a scratch project, then `cat .langdrift/finops/budgets.json`.
