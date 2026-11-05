# 0036: FinOps ledger, receipt and chain demos

**Spec:** .agents/workflow/specs/0005-finops-and-microruntimes-pages.spec.md

**What to build:** the FinOps hero shows the illustrative team budget ledger; selecting a work item updates the delivery receipt in the "Know what it cost" section; the chain section walks through three scenarios in four pillars. The item with unknown spend says so in every place it appears.

**Blocked by:** 0035.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### State

- [ ] `selectedWorkItemAtom` in its own atom module, re-exported by name from the state barrel, initialized from the content module.

### Demos

- [ ] `BudgetLedgerDemo` (client): the envelope as a stacked bar (`role="img"` with a label, four segments in role colours, a visible legend) and the work items as toggle buttons with `aria-pressed`; the unknown item renders the warning copy, never a number; amounts formatted with `Intl.NumberFormat` for the locale; transitions `motion-safe` only.
- [ ] `DeliveryReceipt` (client): a live region with the selected item's title and a `DetailList` of estimate, reserved, observed or unknown, model tier, effort, retries (plural), verification, cleanup and the evidence reference.
- [ ] `PillarScenarios` (client): scenario tabs over `AnimatedTabs` with one tab panel per scenario carrying the ids the tabs emit, and a four-column pillar grid; the FinOps page maps its `chain` copy into it with the pillars intent, delivery, evidence and outcome.
- [ ] The hero slot, the chain section and the receipt section join the page; the chain outro links to the authorization section.

### Tests

- [ ] jsdom tests with a fresh Jotai store: clicking the second work item puts its id in the store and the receipt shows it; the unknown item shows the warning copy and no `0`.

### Verification

- The checks of 0035.
- Keyboard: Tab reaches every work item, Enter and Space select, the receipt announces; the tabs follow the arrow keys; reduced motion shows no transition.
