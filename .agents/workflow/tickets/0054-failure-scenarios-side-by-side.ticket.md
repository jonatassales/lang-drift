# 0054: The two failure scenarios stand side by side and the weekly report folds into them

**Spec:** .agents/workflow/specs/0008-finops-opening.spec.md

**What to build:** after the gate, a section with the kicker "Intent → Delivery → Evidence → Outcome" shows pause and replan and unknown after a timeout side by side from the `lg` breakpoint, stacked below it, without tabs. Each scenario has its title, its context line and the four pillars in order, and ends in a status row: the pause scenario in "Needs attention" with the sentence that the reindexing job is paused and a decision is pending, in the danger tone; the unknown scenario in "Observed spend" with "Unknown · awaiting reconciliation", in the warning tone. The within-envelope scenario and the weekly report section are gone.

**Blocked by:** 0053.

**Status:** ready-for-agent

### Content

- [ ] The chain scenario ids become the two failure scenario ids, each with a status tone in the content module; the pillar ids stay.

### Scenarios

- [ ] `FailureScenarios` (Server Component): two articles with the scenario's title, its context and a `DetailList` of the four pillars closed by the toned status row.
- [ ] The section takes the fragment id `failures`, the old chain title as its kicker, a title that says the envelope holds when the plan does not, and a lead that names the two failures and points at the receipt above for the case inside the envelope. `PillarScenarios`, `SectionOutro` and `RecordCard` leave the page and stay for the Microruntimes page.
- [ ] The weekly report section and its computed rows leave the page.

### Copy

- [ ] `chain` becomes `failures` (kicker, title, lead, pillars, two scenarios with title, context, four pillars and a status term and detail); the tab labels become titles; the selector label, the outro and the within-envelope scenario are removed; `report` is removed and its attention sentence becomes the pause scenario's status detail. Five locales, same key order as `en`.

### Tests

- [ ] The content module test: the failure ids leave the within-envelope case out and every failure scenario has a status tone.
- [ ] A jsdom test renders `FailureScenarios` with the English catalog and asserts the two titles, the pause scenario ending in the needs-attention sentence, the unknown scenario ending in the unknown copy, and no `$0`.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/finops/` and `/pt-br/platform/finops/`: two columns from the `lg` breakpoint, one column at 375px; the status rows in their tones in light and dark.
