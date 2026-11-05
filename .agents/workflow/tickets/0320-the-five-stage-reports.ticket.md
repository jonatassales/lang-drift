# 0320: the five stage reports

**Spec:** .agents/workflow/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** A skill's report has a contract a skill author can validate before a run: one shape per stage, each with a JSON Schema on npm and a check that refuses a report outside it. The stage lines of the run print their figures from these reports (spec 0059).

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] `planning` carries `openQuestionsSettled`, for `console › grilling   3 open questions settled`.
- [ ] `spec` carries `spec`, `revised` and `acceptanceCriteria`, for `console › to-spec    0011-SPEC revised · 3 acceptance criteria`. The simulation prints the spec's short id, `0011-SPEC`, while the field holds the spec's path (`.agents/0011-SPEC-book-a-demo-runtime.spec.md`); the contract keeps the path and says that the line's id is the first two dash-separated parts of the file name, which spec 0059's line takes from it.
- [ ] `tickets` carries `tickets` and `blockingEdgesRecorded`, for `console › to-tickets 3 tickets · blocking edges recorded`.
- [ ] `implement` carries `method` (`tdd`), `components`, `files`, `boundaryOnly` and `stoppedAtReservedPart`, for `console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part`.
- [ ] `review` carries `findings` and `fixed`, for `console › code-review 2 findings · fixed`.
- [ ] Each shape has a JSON Schema and a sample, published on npm; each sample passes the check, one with a missing field is refused, and the Go check and the schemas agree.

## Dev part

**Override:** No part qualifies: five data shapes, their schemas and a mechanical check, with no domain decision of their own (policy, Categories: never boilerplate or configuration). 2026-11-05.
