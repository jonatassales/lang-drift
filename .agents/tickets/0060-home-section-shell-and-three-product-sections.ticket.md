# 0060: The home section shell and the Vision, Microruntimes and FinOps sections

**Spec:** .agents/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** under the hero, three sections in one skeleton: the product's icon and name on a row at the right, a question as the heading and a description at the left, the product's demo, a link to the product page at the end. Vision's is the Vision page's hero as it is, the curve beside the graph; Microruntimes' is the planning session into its cards under a question of its own; FinOps' is the ledger with the docked receipt under the envelope statement.

**Blocked by:** 0065 (Products in the platform order; the features grid goes).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A `HomeProductSection` Server Component renders the skeleton: a landmark labelled by its heading and described by its optional caption and summary; the icon chip and the product name from the header catalog aligned right; a level-two heading in the Vision hero's heading style; a lead in its lead style; the demo; an optional text link with the arrow icon. Role utilities only, no stylesheet.
- [ ] The Vision section reuses the home catalog's hero question (with `ChromaText`, still under reduced motion), description, caption and summary, and renders the curve and the graph in a Tailwind twin of the Vision hero's demo grid; selecting a point or a node on the home moves the other, as on the Vision page. The Vision page and its stylesheet are untouched.
- [ ] The Microruntimes section renders the session terminal and the session cards with the page's props and grid under a new question and description from the platform catalog; its summary is the page's.
- [ ] The FinOps section renders the envelope statement as its heading, a new description, the ledger and the docked receipt with the page's props and grid, the page's caption and summary; unknown spend reads unknown.
- [ ] Each of the three ends with a link to its product page, from the platform catalog's new `sections` copy, drafted in five locales for the owner's review.
- [ ] A jsdom test of the shell asserts the product name, the heading relation, the caption and summary relations when given, and the link only when given.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check localization`.
- `/en/` at desktop width and 375px, light and dark: the three sections in order, labels at the right, links last, no horizontal overflow; the curve and the graph follow each other; the receipt follows the selected work item.
