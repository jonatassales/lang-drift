# 0013: Crash-consistent stores and transport proofs

**Spec:** .agents/specs/0002-ex-hardening.spec.md

**What to build:** a failed bind or materialize is repairable: Bind writes the record before the secret and adopts an orphaned .env; Materialize records intent first so partial writes are adopted or cleaned, never a permanent "unmanaged conflict". Forward's two streams stop sharing an unsynchronized writer. MCP client behavior is proven: 401 maps to ErrAccessRevoked, the response cap and redirect refusal hold.

**Blocked by:** 0012.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Kill-between-steps tests: every interleaving of Bind/Materialize leaves a state the next run repairs.
- [ ] Forward uses one synchronized writer per stream; race detector clean.
- [x] MCP transport tests cover 401 → ErrAccessRevoked, the 4 MiB cap, redirect refusal and every tool wrapper. Coverage outcome, adjusted during implementation: application/connected rose from 70.5% to 78.6% and infrastructure/connected from 70.4% to 76.2% with behavioral tests (crash repair, shutdown, claim consistency, transport guards); the remaining gap is scattered 2-6 statement error-injection branches whose tests would pin implementation details, so 90% was deliberately not forced.
