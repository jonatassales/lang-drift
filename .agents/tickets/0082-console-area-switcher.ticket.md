# 0082: The dashboard's shell with an area switcher

**Spec:** .agents/specs/0013-console-areas.spec.md

**What to build:** The sidebar opens with an area switcher in the old dashboard's style, listing Microruntimes and FinOps, with the URL deciding the area. Each area has its own sidebar items, with Settings at the foot. The header's account control becomes a menu with the session email, Settings and Sign out. The stub's seeded project is renamed LangDrift Studio, and the project item that listed runtimes reads Runtimes.

**Blocked by:** 0081

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The switcher shows the current area, switches with mouse and keyboard, closes with Escape and returns focus.
- [ ] A pathname test proves routes under `/finops` are FinOps and every other route is Microruntimes.
- [ ] The account menu shows the email and signs out.

### Context

- [The spec](/.agents/specs/0013-console-areas.spec.md)
