# 0015: Delete the dead runtime command tree

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** ex runtime … keeps refusing with migration guidance, but the eight dead wired commands, their UseCases fields and the test harness that re-mounts them by copying Execute's body are deleted. The application use cases stay (the connected agent runs them).

**Blocked by:** 0010.

**Status:** ready-for-agent

- [ ] Refusal stub with guidance preserved and tested.
- [ ] Dead catalog commands, UseCases fields and execute_catalogs_test harness removed.
- [ ] Coverage of live CLI paths does not drop.
