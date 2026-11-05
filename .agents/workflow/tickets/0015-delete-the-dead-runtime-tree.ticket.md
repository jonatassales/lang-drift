# 0015: Delete the dead runtime command tree

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** ex runtime … keeps refusing with migration guidance, but the eight dead wired commands, their UseCases fields and the test harness that re-mounts them by copying Execute's body are deleted. The application use cases stay (the connected agent runs them).

**Blocked by:** 0010.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Refusal stub with guidance preserved and tested.
- [x] Adjusted during implementation: the UseCases runtime fields, the duplicated execute_catalogs_test harness and every unit test of the retired commands are removed. The catalog command constructors stay as package code because the connected agent and the tests drive the native lifecycle through them.
- [ ] Coverage of live CLI paths does not drop.
