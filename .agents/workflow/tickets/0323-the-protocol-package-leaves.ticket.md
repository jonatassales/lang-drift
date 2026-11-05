# 0323: the protocol package leaves

**Spec:** .agents/workflow/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** No contract describes a channel the product does not have. The `protocol` package (the MCP version, the five tools, bindings, claims, events, intents with leases, operations, observations, measurement windows, memberships, repositories, sync and agent states) leaves the Core contracts once the CLI's connected agent, its last consumer, is gone.

**Blocked by:** 0322, and spec 0048's ticket that removes the connected-agent commands (`agent connect`, `agent project`, `agent supervise`) and the connected agent.

**Status:** ready-for-agent

**Alignment:** remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The `protocol` package and its tests are gone, and nothing in the repository imports it.
- [ ] No terminal line or card field of the simulation is lost: none came from it.
- [ ] The module and the CLI build and pass, in the workspace and outside it.

## Dev part

**Override:** No part qualifies: the removal of a package with no remaining caller (policy, Categories: never a mechanical removal). 2026-11-05.
