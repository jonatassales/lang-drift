# 0152: The website footer's styles become utilities

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** `website-footer.css` and the footer rules left in `hero.css` become utilities in the footer, and both files are deleted. The brand, the description, the links, the credit line and the social links look as they did. The 760px breakpoint becomes `md`. The `compact` prop and its `data-compact` styles, which nothing renders, are removed. Where the footer's elements carry `smooth-footer`, the utilities that must beat SmoothUI take `!`, marked transitional.

**Blocked by:** 0150

**Status:** ready-for-agent

- [ ] `website-footer.css` and `hero.css` are deleted and nothing imports them.
- [ ] The footer keeps its look and its hover and focus states.
- [ ] The `compact` prop and its styles are gone.
- [ ] Every `!` that beats SmoothUI carries a transitional note.
- [ ] The captures match the baseline.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check tailwind`
- Built-in browser captures compared with the baseline.
