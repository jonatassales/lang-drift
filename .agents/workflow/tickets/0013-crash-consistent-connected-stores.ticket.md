# 0013: Crash-consistent stores and transport proofs

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** a failed bind or materialize is repairable: Bind writes the record before the secret and adopts an orphaned .env; Materialize records intent first so partial writes are adopted or cleaned, never a permanent "unmanaged conflict". Forward's two streams stop sharing an unsynchronized writer. MCP client behavior is proven: 401 maps to ErrAccessRevoked, the response cap and redirect refusal hold.

**Blocked by:** 0012.

**Status:** ready-for-agent

- [ ] Kill-between-steps tests: every interleaving of Bind/Materialize leaves a state the next run repairs.
- [ ] Forward uses one synchronized writer per stream; race detector clean.
- [ ] MCP transport tests for 401, size cap, redirect; connected app+infra coverage ≥ 90%.
