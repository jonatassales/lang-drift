# 0055: Definitions, intended sources and a band that speaks for FinOps

**Spec:** .agents/specs/0008-finops-opening.spec.md

**What to build:** the budget section explains envelope, allocation and reservation as a definition list instead of three cards; the sources section lists the four sources under an "Intended sources" label, the Microruntimes one keeping its link, with the intended-model sentence closing the list; and the early-access band closes the page with a FinOps sentence instead of Vision's sentence about product decisions. The Vision and Microruntimes bands read as before.

**Blocked by:** 0052.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Definitions and sources

- [ ] The budget section renders a `DetailList` whose terms are the former card eyebrows and whose details set the former title stronger over its description; `InfoCards` leaves the page and stays for the Microruntimes authoring flow.
- [ ] `SourceList` (Server Component): a plain list, two columns from the `md` breakpoint, under a mono "Intended sources" label; each item's title set stronger over its description; the Microruntimes item keeps its link; the intended-model sentence closes the list.

### Band

- [ ] `EarlyAccessSection` accepts an optional description that the landing band shows instead of the shared sentence; without one it renders as before.
- [ ] The FinOps page passes the description from its catalog.

### Copy

- [ ] `budget.cards` becomes `budget.terms` with `term`, `title` and `description`; `sources` gains `intended`; a new `earlyAccess.description` invites the reader to see what a team envelope, a reservation and a receipt would look like for their own work and what FinOps does not promise. Five locales, same key order as `en`.

### Tests

- [ ] The early-access section test gains a case: with a description the landing band shows it and not the shared sentence; without one, the shared sentence.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `pnpm turbo run build --filter=website`, `drift check localization`, `drift check dead-code`.
- `/en/platform/finops/` and `/pt-br/platform/finops/` at desktop width and 375px, light and dark: the definitions, the two-column source list with its link, the FAQ and the band's FinOps sentence; `/en/platform/vision/` and `/en/platform/microruntimes/` unchanged.
