# 0009: Failures carry their cause

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** when a build or lifecycle subprocess fails, the error names the command and includes a bounded tail of its stderr; long build output streams to an evidence file instead of failing at 1 MiB, and the error names that file. No secret values appear in errors.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] process gains a spill-to-file capture used by build recipes; cap failures disappear for builds.
- [ ] native.go stops discarding error causes and stderr (build, stop, report paths).
- [ ] relocateMigration's unchecked type assertions replaced with checked handling.
- [ ] Tests cover: failing build shows tail; >1 MiB build succeeds and evidence file holds full output.
