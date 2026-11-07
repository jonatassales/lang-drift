# 0325: the old judge contract leaves

**Spec:** .agents/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** The judge's contract carries only the three answers the plan asks for. Once the plan asks through the `judgment` answers, reads its documents deterministically and plans Safe Seat without a judge, the old `judge` package leaves: its `sources` and `reserved-part` kinds, its free `reasons` list and its run-keyed episodes, and the request and response envelope with `usage`, `cost`, `refused`, `reason`, `judgeVersion`, `system`, `prompt`, `maxTokens` and `promptVersion`. Its npm path goes with it.

**Blocked by:** 0322, and the tickets of specs 0051 (the judge answers through `judgment`), 0053 (document reading needs no `sources` question) and 0061 (Safe Seat needs no `reserved-part` answer). The console's `/api/judge` route, outside this phase, still reads the old schemas; where the judge runs is spec 0051's decision, and that decision says what becomes of the route first.

**Status:** ready-for-agent

**Alignment:** remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The old `judge` package, its schemas and samples, and the `./judge/*` export are gone, and nothing in the repository reads them.
- [ ] No judge answer carries a price, against `# plan only · nothing is built, started or paid for in this step`.
- [ ] The module and the CLI build and pass, in the workspace and outside it.

## Dev part

**Override:** No part qualifies: the removal of a contract with no remaining caller (policy, Categories: never a mechanical removal). 2026-11-05.
