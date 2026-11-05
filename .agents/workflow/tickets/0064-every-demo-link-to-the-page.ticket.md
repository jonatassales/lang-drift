# 0064: Every "Book a demo" leads to the page

**Spec:** .agents/workflow/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** the header's rainbow button, on the bar and in the drawer, the hero's link, a band that closes the home after the FAQ and the band that closes the Vision, Microruntimes and FinOps pages all read "Book a demo" and open the Book a demo page; "Ask for a demo" appears nowhere; the pricing page keeps its early-access band.

**Blocked by:** 0062 (One FAQ with a tab per product), 0063 (The Book a demo page, its form and its route).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A `BookDemoBand` Server Component with the fragment id `book-a-demo` shows "Book a demo", its own description or the sentence a page passes, and a pill link to the page; a jsdom test covers both descriptions and the link.
- [ ] The home renders the band after the FAQ; the Vision, Microruntimes and FinOps pages render it in place of the landing early-access band, FinOps with its own sentence.
- [ ] `EarlyAccessSection` keeps only its pricing branch; the landing branch, the description prop and the "Ask for a Demo" keys leave the code, the catalog and the test.
- [ ] The Vision page's three in-page links that pointed at the early-access fragment point at the band's.
- [ ] The header's demo button is a link to the page with its rainbow ring, and tapping it in the drawer closes the drawer; the hero's secondary call to action links to the page.
- [ ] The website README says demo calls to action lead to the page, whose form validates and answers that requests are paused.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check dead-code`, `drift check localization`.
- `/en/`, `/en/platform/vision/`, `/en/platform/microruntimes/`, `/en/platform/finops/`: the band last before the footer; every demo link lands on `/en/book-a-demo/`; the Vision page's in-page links reach its band; the drawer closes on tap; `/en/pricing/` keeps its early-access band.
