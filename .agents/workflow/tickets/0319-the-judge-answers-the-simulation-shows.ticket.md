# 0319: the judge's three answers the simulation shows

**Spec:** .agents/workflow/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** The Core contracts carry the three answers the plan asks its judge for, and nothing else: the guardrail `criteria`, the `episodes` verdicts and the `suggestion`. Each is typed, has a JSON Schema and a sample, and `ParseAnswer` refuses anything outside it. They live in a new `judgment` package, because the old `judge` package keeps its sources, reserved-part and priced envelope for the CLI's current planner and the console until spec 0051 moves them (ticket 0325 removes it).

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The `criteria` answer holds one `id` and `statement` per criterion, so that `console › judging · FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load` and the guardrail row `FIN-02: estimate before a paid call · FIN-03: a budget before a paid run · FIN-04: size cpu and memory to the load` come from it.
- [ ] The `episodes` answer holds per candidate an `ep-NNNN` id, a `kept` or `dropped` verdict and a reason from a closed set (`shares the stack and scale` when kept; `other runtime`, `no broker` or `docs only` when dropped), so that `console › judging · 3 share the stack and scale · 3 differ: other runtime, no broker, docs only` and the EPISODES rows (`ep-0412 checkout-preview-sizing 0.81 kept`, `ep-0297 docs-search-index 0.52 dropped · docs only`) come from it.
- [ ] The `suggestion` answer holds `model`, `effort`, `effortReason` and `modelReason`, so that `console › weighing effort · low skips the broker wiring · high is not warranted by 5 components`, `console › weighing model · Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome`, `✓ suggests   claude-sonnet-5-5 · effort medium`, SUGGESTED MODEL `Claude Sonnet 5.5` and SUGGESTED EFFORT `medium` come from it.
- [ ] No answer carries a cost, usage, provenance or version: the plan prints `# plan only · nothing is built, started or paid for in this step`.
- [ ] Each kind has a JSON Schema and a sample, and the samples parse with Go and agree with the schemas; ticket 0322 publishes them on npm with the other new shapes.
- [ ] The old `judge` package is untouched, and the CLI and the console keep building and passing.

## Dev part

**Category:** `core-domain`. Whether a judge's verdict on an episode is acceptable decides what the Analysis card shows as kept or dropped and why; a reason outside the set would read differently on every surface.

**Paths:**
- `lib/core/contracts/go/judgment/check-episodes.go`
- `lib/core/contracts/go/judgment/check-episodes_test.go`

**Boundary:** `func checkEpisodes(answer EpisodesAnswer) error`, called by `ParseAnswer` for the `episodes` kind; `EpisodesAnswer` and its entry type are declared by the agent in their own files.

**Goal:** Accept an episodes answer only when every entry has an `ep-NNNN` id, no id repeats, a kept entry's reason is `shares the stack and scale`, and a dropped entry's reason is one of `other runtime`, `no broker` or `docs only`; refuse naming the field and the rule otherwise.

- [ ] The simulation's six verdicts pass.
- [ ] A malformed id, a repeated id, a kept reason on a dropped entry and a reason outside the set are each refused naming the field.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `checkEpisodes` and its test.
