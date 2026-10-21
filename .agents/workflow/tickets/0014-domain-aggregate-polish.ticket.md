# 0014: Guarded aggregates and single validations

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** Execution and cluster.Record guard their transitions (MarkUnit, BeginDown included) and stop exposing mutable state the application pokes; RecordCreated is bounds-checked; the port-range and source-location validations exist once; domain/connection gets sentinel errors; time and ID generation move behind Clock and IDSource ports wired from the composition root.

**Blocked by:** 0012.

**Status:** ready-for-agent

- [ ] No application file assigns an aggregate field directly; transitions return errors on misuse.
- [ ] One definition each for port-range and source-location rules; duplicates deleted.
- [ ] connection errors match with errors.Is; tests updated.
- [ ] time.Now/crypto-rand appear only in adapters and the composition root.
