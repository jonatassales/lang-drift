# 0180: Make Code Quality CI gate required

**What to build:** Update `.github/widgets/dashboard.sh` to publish the "Code Quality" check run as an active passing CI gate with conclusion `success` and title `Rules Verified` (instead of `neutral`), aligning it with the other four CI gates (Workflow, FinOps, Microruntimes, and Agentic). Update documentation in `.github/widgets/code-quality.sh`.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Update `.github/widgets/dashboard.sh` to publish `Code Quality` check run with conclusion `success` and title `Rules Verified`.
- [x] Update `.github/widgets/code-quality.sh` to document that the Code Quality check run is an active, required CI gate.
- [x] Verify `platform/cli/drift widget dashboard` executes cleanly.
