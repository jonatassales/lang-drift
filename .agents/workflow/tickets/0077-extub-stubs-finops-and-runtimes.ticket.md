# 0077: The stub API serves FinOps and runtimes only

**Spec:** .agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md

**What to build:** `apps/microservices/console-stubs` seeds one product, Microruntimes, with its four runtimes and their FinOps readings. The Vision fixture types, the product field that carries them, the seeded Vision product and its local folder go. Every route stays. The console shows one product and keeps its product switcher disabled.

**Blocked by:** 0076

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The console home reads one product from the stub, and every route answers as before.
- [ ] The local folder list offers no Vision folder.
- [ ] Widget and secret isolation are proven against a product the tests create, and the stub's suite passes under the race detector.

### Context

- [The spec](/.agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md)
