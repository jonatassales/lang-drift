# 0053: The gate is one flow carrying the receipt's figures

**Spec:** .agents/specs/0008-finops-opening.spec.md

**What to build:** right after the hero, the authorization section shows the gate as one horizontal flow of three steps instead of three ordered cards: the agent recommends with 1,800 dollars proposed, the gate reserves 2,600 before the first paid call, the evidence closes with 2,460 observed and 140 released to the envelope. A line under the flow says the figures come from the Checkout retries receipt above and are illustrative. The kicker still reads "Recommendation is not authorization" and the lead still says the gate is the intended design, not a shipped control.

**Blocked by:** 0052.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Content

- [ ] The content module exports the gate example, the work item the ledger selects first, and a pure function deriving the gate figures from a work item with observed spend: proposed is the estimate, reserved the reservation, observed the observed spend, released the difference between the two. The authorization card ids become the gate step ids.

### Flow

- [ ] `GateFlow` (Server Component): an ordered list of three steps in one row from the `lg` breakpoint with an arrow between steps, stacked below it with the arrows turned; each step keeps its eyebrow, title and description and carries a mono figure with its label, the third step a second figure for the released amount; amounts formatted with the module's dollar formatter for the locale, never typed into copy or markup.
- [ ] The source line under the flow names the receipt's work item through an ICU message and says the figures are illustrative.
- [ ] The authorization section keeps its kicker, title, lead and fragment id and moves to sit right after the hero.

### Copy

- [ ] `authorization.cards` becomes `authorization.steps` with a figure label on each step and a released label on the third, plus `authorization.source`, in five locales in the same key order as `en`.

### Tests

- [ ] The content module test: the gate example is the initial work item and has observed spend; its figures are 1,800, 2,600, 2,460 and 140.
- [ ] A jsdom test renders `GateFlow` with the English catalog and asserts the three eyebrows in order, the four formatted amounts and the source line, with no `$0`.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/finops/`: the flow reads left to right with the four amounts; at 375px the steps stack with the arrows turned.
