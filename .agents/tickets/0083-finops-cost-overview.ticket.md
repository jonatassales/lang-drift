# 0083: FinOps opens with a cost overview

**Spec:** .agents/specs/0013-console-areas.spec.md

**What to build:** `/finops` shows the summed known LLM cost across runtimes, the count of executions whose cost is unknown, and per runtime its project, CPU and memory against the envelope, its LLM budget when configured, its FinOps reading and its recommendation action, all from the stub.

**Blocked by:** 0082

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A test proves the totals add only measured costs and count the unknown ones.
- [ ] The page renders in the four console locales with catalog parity, and an unknown cost never reads as zero.

### Context

- [The spec](/.agents/specs/0013-console-areas.spec.md)
