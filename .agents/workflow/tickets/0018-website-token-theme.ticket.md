# 0018: Website token theme in Exphere purple

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** `@langdrift/tokens` builds one stylesheet, `website.css`, that carries the website's look: the Exphere purple as primary, the website's neutrals as shared roles, classification and drift colors as fixed groups, the editorial and mono fonts, and Tailwind's `dark` variant on `.dark`. The Vision and FinOps themes are removed.

**Blocked by:** 0017 (Apps replace products).

**Status:** ready-for-agent

- [ ] The package exports only `website.css`, with the `--site-*` theme behind `primary` and `primary-hover` (`#6D0AD6`; hover `#5408A5` light, `#8C3BEF` dark) and a white primary foreground.
- [ ] Background, card, popover, foreground, muted foreground and border take the website's light and dark values.
- [ ] Classification colors (intentional, aligned, unexplained, review) and drift colors (negative, positive, zero) are exposed as Tailwind colors and never reference the brand.
- [ ] `--font-sans` is the editorial stack and `--font-mono` comes from a new mono token.
- [ ] The stylesheet declares the `dark` variant on the `.dark` class.
- [ ] The contract test covers all of the above; the cross-product test is gone.
- [ ] AGENTS.md's design section and the glossary describe the single theme and the two new groups.
