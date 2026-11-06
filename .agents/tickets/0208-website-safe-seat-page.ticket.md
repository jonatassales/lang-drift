# 0208: the Safe Seat page

**Spec:** .agents/specs/0043-safe-seat-on-the-website.spec.md

**What to build:** A visitor opens `/platform/safe-seat` in any of the five locales and learns what Safe Seat is and why LangDrift FinOps keeps a person writing part of every ticket. The page states that the capability is coming soon.

The hero has an eyebrow reading "LangDrift FinOps", the "Coming soon" badge, a two-part heading and a lead. The sections follow in this order:

- **Why:** writing code is a way of thinking, and the irony of automation.
- **Evidence:** six cards, each giving a finding, its limit and a link to the primary source.
- **What stays with the person.**
- **How it will work:** five planned steps.
- **Principles.**
- **What the evidence does not show.**
- **Sources.**
- **FAQ.**

The page carries its own title, description, canonical, hreflang, structured data and sitemap entries.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `/en/platform/safe-seat/` and the four other locale prefixes render every section in order. The hero shows the "Coming soon" badge, and the page has no call to action.
- [ ] Every evidence card states its finding, its limit and a link to its primary source. Burnout reads as an association. No clinical or health vocabulary appears in any locale.
- [ ] The colleague's post is paraphrased, never quoted or attributed. Any quote from a source is a few words long.
- [ ] "Safe Seat" is never translated and never abbreviated in any locale.
- [ ] The `safe-seat` catalog exists in five locales with `en`'s key order. It is registered with the loaders, the shape validator, the messages barrel and the catalog parity test, and the parity test passes.
- [ ] A content module holds the id lists, the source links and the `SafeSeatMessages` type. Its unit test checks that the ids match the catalog and that every link is an absolute `https` URL.
- [ ] `pages.json`, the nested page path and the sitemap carry the page (35 sitemap entries in total). The SEO contract and structured-data tests pass.
- [ ] No horizontal overflow at 375px, in light and in dark.

## Dev part

**Override:** no part of this ticket qualifies. It is copy in five locales, page composition and metadata registration, with no domain logic, hard part or new stack. Proposed by `to-tickets` on 2026-11-04; the owner may replace it.

## Verification

- `pnpm --filter website test`
- `pnpm --filter website typecheck`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
