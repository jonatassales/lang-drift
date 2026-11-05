# 0048: The lifecycle steps move from the hero to the evidence section

**Spec:** .agents/workflow/specs/0007-microruntimes-hero.spec.md

**What to build:** a visitor still selects a lifecycle step and reads what it does, the evidence it leaves and the record it produces, but in the evidence section, next to that record, instead of in the hero. The hero's old lifecycle card and its runtime summary are gone, so the hero is ready to receive the ticket card. Nothing else on the page changes.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Evidence section

- [ ] `LifecycleSteps` (client): the six `ChoiceButton` steps with their numbered badges and the "What it does" and "Evidence it leaves" panel, extracted from the old hero card, mounted in the evidence section above `StageEvidence`, still writing `selectedLifecycleStageAtom`.
- [ ] `StageEvidence` and the atom are unchanged; selecting a step still moves the evidence record.

### Hero

- [ ] `LifecycleDemo` is removed from the hero and deleted with its test; the hero keeps its heading, lead, calls to action and the terminal until ticket 0049 replaces it.
- [ ] The catalog keys only the old card used (its label, the runtime summary terms, the components plural, the requested and observed states and the observed detail) are removed in all five locales; the stepper label, the stages, "What it does" and "Evidence it leaves" stay, and the CPU and memory value formats stay for ticket 0049 to reuse.

### Tests

- [ ] The old card's test moves to `LifecycleSteps`: clicking "Prepare" moves `aria-current`, puts `prepare` in the store and changes the evidence record's text; the runtime summary assertions go with the summary.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check localization`, `drift check dead-code`.
- `/en/platform/microruntimes/` renders the steps in the evidence section; Tab reaches every step, Enter and Space select, the panel and the record announce.
