# 0052: The FinOps page opens with the statement, the ledger and the docked receipt

**Spec:** .agents/workflow/specs/0008-finops-opening.spec.md

**What to build:** `/[locale]/platform/finops/` opens with a mono eyebrow "LangDrift FinOps" and the left-aligned statement "One team envelope. Every paid call reserved before it runs." as its `h1`, with no gradient, no lead and no buttons. Under it the budget ledger fills three fifths of the width and the delivery receipt the remaining two fifths from the `lg` breakpoint, stacked below it, so selecting a work item fills its receipt at once. The illustrative caption stays under the demo. The receipt section two sections below is gone, and `ProductHero` is deleted because nothing renders it any more. The sections after the hero stay as they are for now.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

### Opening

- [ ] The page composes its own hero section labelled by the `h1`: the eyebrow, the statement in the product pages' heading classes with no `ChromaText`, then a grid wired by `aria-describedby` to the visible caption and the screen-reader summary, with the ledger first and the receipt second; the caption sits under the grid, left-aligned. No lead and no calls to action.
- [ ] The skip link keeps its text and points at the gate section, the first section after the hero.
- [ ] `ProductHero` and its copy interface are deleted; the Microruntimes and Vision pages are untouched.

### Ledger and receipt

- [ ] `BudgetLedgerDemo` loses its own centering and maximum width and fills its column; nothing else changes.
- [ ] `DeliveryReceipt` loses its section top margin, its heading becomes an `h2`, and its head stacks the eyebrow, title and owner over the observed spend at every width. The live region, the rows and the unknown-spend rendering stay.
- [ ] The receipt section with its kicker, title and lead leaves the page.

### Copy

- [ ] `hero` in the FinOps catalog becomes `eyebrow`, `title`, `caption` and `summary` in five locales, in the same key order as `en`: the eyebrow is the product name untranslated, the title is the statement, the caption keeps its text and the summary describes the ledger, the receipt beside it and the simulated amounts. `receipt` loses `kicker`, `title` and `lead`. The catalog's TypeScript shape follows.

### Tests

- [ ] The ledger and receipt jsdom test keeps its cases and finds the selected item's title as a level-two heading.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/finops/` and `/pt-br/platform/finops/`: the eyebrow and the statement alone above the ledger, the receipt beside it from the `lg` breakpoint and under it at 375px, following the selected item; no horizontal overflow.
