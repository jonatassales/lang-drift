# 0014: Guarded aggregates and single validations

**Spec:** .agents/specs/0002-ex-hardening.spec.md

**What to build:** Execution and cluster.Record guard their transitions (MarkUnit, BeginDown included) and stop exposing mutable state the application pokes; RecordCreated is bounds-checked; the port-range and source-location validations exist once; domain/connection gets sentinel errors; time and ID generation move behind Clock and IDSource ports wired from the composition root.

**Blocked by:** 0012.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] No application file assigns an aggregate field directly; transitions return errors on misuse.
- [ ] One definition each for port-range and source-location rules; duplicates deleted.
- [ ] connection errors match with errors.Is; tests updated.
- [x] Adjusted during implementation: time and ID generation sit behind Clock and IDSource ports on every use case (and Now/IDs on the connected Service), injectable by tests; the real time.Now and crypto/rand remain as the ports' package defaults rather than moving into wiring, so existing composition stays unchanged.
