# 0144: Every console file declares one function, named after it

**Spec:** .agents/workflow/specs/0026-console-rules-compliance.spec.md

**What to build:** Multi-function files are split into folders by concern. Every file is named after its export, with the suffix from the 04 table. entity-title leaves metadata/, and the extub route handler delegates to a .route file. No visual or behavioural change.

**Blocked by:** 0143.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).
