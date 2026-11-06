# 0324: the fence v3 is the only fence

**Spec:** .agents/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** `runtime.md` carries only what the simulation shows. Once the plan writes the v3 fence and the run and the cluster read it, the v2 `Configuration` leaves with everything only it served, and an older draft is refused with "re-plan the runtime" instead of being misread.

**Blocked by:** 0322, and the tickets of specs 0055 (the plan writes the v3 fence), 0066 (the components resolve into it) and 0062 (the cluster applies it).

**Status:** ready-for-agent

**Alignment:** remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The v2 fields leave: `runtimeUid`, `projectId`, the runtime's `name`, `llmBudgetUsd`, `native` and the migration rule, and per component `name`, `id`, `kind`, `dockerfile`, `readiness`, `envFile` and `passEnv`; so do the env-file parser, the variable-name rule, the Dockerfile path helpers, the derived image tag and the Markdown wrapper prose.
- [ ] Parsing a version-2 fence fails with "re-plan the runtime".
- [ ] The v3 types take the names the v2 ones held, and the fence's schema is the only one published.
- [ ] The figures the simulation prints from the fence stay exact: `✓ plan       runtime.md draft, revision 1 · proposed to the console`, `runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB`.
- [ ] The module and the CLI build and pass, in the workspace and outside it.

## Dev part

**Override:** No part qualifies: removals and a rename once no caller remains (policy, Categories: never a rename or a mechanical removal). 2026-11-05.
