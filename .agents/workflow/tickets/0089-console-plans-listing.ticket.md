# 0089: The console's Plans listing

**Spec:** .agents/workflow/specs/0014-console-ticket-phases.spec.md

**What to build:** the Plans entry opens an infinite list of the tickets that have reached the plan phase, each as a card built on the website's Plan panel: the proposed plan and its authorization state, the envelope rows with share bars, the components with their technology marks, figures and shares, the ports, and the report (bounded, headroom, deployments, GPU and disk, price, validation, evidence). The ticket's detail page shows the same widget.

**Blocked by:** 0086, 0087.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `products/[id]/plans` with page, loading, error and template, reading `reached=` the phase.
- [ ] `PlanWidget` with the panel's eyebrow and summary line over its full body; unknown values are a dash.
- [ ] The detail page renders the widget once the phase is reached.
- [ ] Copy in the four console locales; Node tests for the summary line and any figure the widget computes.

### Verification

- `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`; the listing in the browser.
