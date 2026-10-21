# 0015: Delete the dead runtime command tree

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** ex runtime … keeps refusing with migration guidance, but the eight dead wired commands, their UseCases fields and the test harness that re-mounts them by copying Execute's body are deleted. The application use cases stay (the connected agent runs them).

**Blocked by:** 0010.

**Status:** ready-for-agent

- [ ] Refusal stub with guidance preserved and tested.
- [x] Adjusted during implementation: the UseCases runtime fields, the duplicated execute_catalogs_test harness and every unit test of the retired commands are removed. The catalog command constructors stay as package code because the e2e suite drives the native lifecycle through them; the composition root exposes those use cases to that harness only.
- [ ] Coverage of live CLI paths does not drop.
