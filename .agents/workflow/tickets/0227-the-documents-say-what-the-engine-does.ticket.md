# 0227: the documents say what the engine does

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** A reader of the CLI README, the context document, spec 0032 and this repository's microruntime folder finds the engine described as it now behaves. The CLI README's command table has `stop`, no `init` under `agent` and no `--no-judge`; its "Configuration" section says planning needs a connected project and names the planning cost the gate reserves; spec 0032's sentence about the run is amended to the new order with a note pointing here; `platform/microruntime/plan.sh` drops `--no-judge`, its README says the project must be connected and how (`ex agent connect` with `RUNTIME_SECRET`), and ticket 0198's run criteria are revisited against 0.5.0. The context document's commands table and glossary name the judge as the platform's service.

**Blocked by:** 0221, 0223, 0224, 0226, and the merge of ticket 0198's plan script (pull request #117).

**Status:** ready-for-agent

- [ ] The CLI README (its connect, plan and stop lines written in 0224 and 0226), the translated READMEs, the context document and the microruntime README name no `--no-judge`, no `ex agent init` and no provider key for planning, and describe connect, plan, run (runtime up first), confirm, stop and status.
- [ ] Spec 0032's "This repository" paragraph says the run brings the components up and pauses, consistent with spec 0045, with a note of the amendment and its date.
- [ ] `sh platform/microruntime/plan.sh` has no `--no-judge`; its header and README say the project must be connected first.
- [ ] `git grep -n -i 'no-judge\|agent init\|ANTHROPIC_API_KEY' -- lib README.md platform .agents/context.md apps/console/README.md` finds only the driver's pass-through and the console's own provider key.

## Dev part

**Override:** No part qualifies: documentation and a script flag follow decisions taken in 0224 to 0226; nothing here has a problem of its own to solve (policy, Dev part: a problem with its own goal and acceptance criteria). 2026-11-04.

## Verification

- The `git grep` above.
- `sh -n platform/microruntime/plan.sh`
- `platform/cli/drift quality check 01-filesystem`
