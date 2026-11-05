# 0038: Microruntimes lifecycle demos

**Spec:** .agents/workflow/specs/0005-finops-and-microruntimes-pages.spec.md

**What to build:** the Microruntimes hero shows the six-step lifecycle stepper and the runtime summary; selecting a step shows what it does and the evidence it leaves, and the evidence section mirrors the selected step's record; the lifecycle section walks through three scenarios from compose to clean up. The observed state reads unknown wherever it is unknown.

**Blocked by:** 0036, 0037.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### State

- [ ] `selectedLifecycleStageAtom` in its own atom module, re-exported by name from the state barrel.

### Demos

- [ ] `LifecycleDemo` (client): the stepper as an ordered list of buttons with numbered badges, a connector line and `aria-current="step"` on the selected step; the runtime summary as a `DetailList` with the observed state in the warning tone; a live stage panel with "what it does" and "evidence it leaves".
- [ ] `StageEvidence` (client): the selected step's record (kind, title, provenance) in the evidence section, as a live region.
- [ ] The lifecycle section reuses `PillarScenarios` with the phases compose, run, inspect and clean up and the scenarios add a component, preparation fails and the agent loses its lease; its outro links to the state section.

### Tests

- [ ] A jsdom test with a fresh Jotai store: clicking the "prepare" step moves `aria-current` and puts the stage in the store; the summary shows desired `running` and observed `unknown`.

### Verification

- The checks of 0035.
- Keyboard: Tab reaches every step, Enter and Space select, the stage panel and the evidence record announce; the lease-lost scenario says cleanup was attempted and its outcome stays unknown.
