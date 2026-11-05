# 0062: One FAQ with a tab per product

**Spec:** .agents/workflow/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** a fifth home section, a FAQ for the whole platform: tabs for Vision, Microruntimes and FinOps switch one accordion between the questions each product page already answers, collapsed on every switch, and the tabs move with the arrow keys.

**Blocked by:** 0065 (Products in the platform order; the features grid goes).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The Vision FAQ's id list moves to a data module the Vision page and the home read; the Vision page's FAQ renders as before.
- [ ] A `PlatformFaqSection` Server Component builds the three groups from the three product catalogs, labelled with the header's product names, in the product FAQ's shell with the fragment id `faq`, and renders `FaqTabs`.
- [ ] `FaqTabs` is a Client Component: controlled Animated Tabs, a tab panel per group wired by the tabs' ids, the Basic Accordion mounted only inside the active panel so a switch shows it collapsed.
- [ ] The platform catalog gains the FAQ's title, description and tab-list label in five locales, drafted for the owner's review.
- [ ] A jsdom test asserts that selecting Microruntimes shows its first question and none of Vision's, that an item expanded under one tab is collapsed after switching away and back, and that arrow keys move the selection.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check localization`.
- `/en/`: the FAQ after AI Voice; tabs switch by pointer and keyboard; an open answer closes on switch.
