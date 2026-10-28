# 0090: The console's Workflows listing

**Spec:** .agents/workflow/specs/0014-console-ticket-phases.spec.md

**What to build:** the Workflows entry opens an infinite list of the tickets that have reached the workflow phase, each as a card built on the website's Workflow panel: the model and effort, the human-in-the-loop stages with their progress, the estimated spend against the budget, and the result with reserved, observed and released amounts, the pull request, its checks and the outcome. The ticket's detail page shows the same widget.

**Blocked by:** 0086, 0087.

**Status:** ready-for-agent

- [ ] `products/[id]/workflows` with page, loading, error and template, reading `reached=` the phase.
- [ ] `WorkflowWidget` with the panel's eyebrow and summary line over its full body; unknown values are a dash.
- [ ] The detail page renders the widget once the phase is reached.
- [ ] Copy in the four console locales; Node tests for the summary line and any figure the widget computes.

### Verification

- `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`; the listing in the browser.
