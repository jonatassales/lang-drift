# 0088: The console's Analysis listing

**Spec:** .agents/specs/0014-console-ticket-phases.spec.md

**What to build:** the Analysis entry opens an infinite list of the tickets that have reached the analysi phase, each as a card built on the website's Analysi panel: the documents read with what each yielded, the episodes kept and dropped against the threshold with their similarity bars, the measured resources, the suggested model and effort, and that the team authorizes the recommendation. The ticket's detail page shows the same widget.

**Blocked by:** 0086, 0087.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `products/[id]/analysis` with page, loading, error and template, reading `reached=` the phase.
- [ ] `AnalysiWidget` with the panel's eyebrow and summary line over its full body; unknown values are a dash.
- [ ] The detail page renders the widget once the phase is reached.
- [ ] Copy in the four console locales; Node tests for the summary line and any figure the widget computes.

### Verification

- `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`; the listing in the browser.
