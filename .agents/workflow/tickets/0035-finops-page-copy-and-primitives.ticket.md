# 0035: FinOps page, copy and product-page primitives

**Spec:** .agents/workflow/specs/0005-finops-and-microruntimes-pages.spec.md

**What to build:** `/[locale]/platform/finops/` renders in the five locales with the Vision page's layout: header, a hero with the question headline, lead and two calls to action, the budget, authorization, report and sources sections, the FAQ, the demo band and the footer, with its own metadata and sitemap entries. The hero's demo slot and the chain and receipt sections arrive with ticket 0036.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

### Primitives in `@langdrift/web`

- [ ] `SectionHeading` (kicker, heading with id, optional lead), `DetailList` (`dl` rows with an optional status tone) and `InfoCards` (card grid with eyebrow, title, description, optional icon and link, ordered or not), each Tailwind-only, exported by name from its own entry in the package export map, with a `renderToString` unit test.

### Page shell in the website

- [ ] `ProductHero` (the Vision hero layout in Tailwind: centered intro, `h1` with the highlight in `ChromaText`, primary and text-link calls to action with a 44px tap target, an optional demo slot described by a caption and a screen-reader summary), `ProductSection` (section shell with scroll margin, top border and `SectionHeading`) and `ProductFaqSection` (Tailwind-only FAQ over `BasicAccordion`).

### Copy and data

- [ ] `finops.json` in the five locales, same key order everywhere: `skip`, `hero` (start, highlight, end, lead, two calls to action, caption, summary), `budget`, `authorization`, `report`, `sources`, `faq`, plus the `ledger`, `chain` and `receipt` keys that 0036 consumes, so the catalog's TypeScript shape is complete in this ticket.
- [ ] The FinOps content module: work-item, scenario and pillar ids, the envelope amounts, the weekly report numbers, the initial selection and the catalog type; a unit test for the envelope sum and the single unknown item.
- [ ] The catalog index asserts the shape and exports the typed catalog; the parity test's family list gains it.

### Route and metadata

- [ ] The page with the Vision boilerplate, the skip link, the hero, the budget, authorization, report and sources sections, the FAQ, the demo band and the footer.
- [ ] `pages.json` gains `finops` in five locales; the nested-path map gains `platform/finops/`; `sitemap.xml` gains the five FinOps URLs with their alternates (25 in total).

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/finops/` renders with the FinOps title and canonical; `/pt-br/platform/finops/` renders translated copy; nothing overflows at 375px.
