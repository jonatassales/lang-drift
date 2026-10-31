# 0153: The pricing plans' styles become utilities

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** `conversion.css` becomes utilities in the plans section and the pricing page, and is deleted. The introduction, the trial card, its badge, value, features and call to action, and the availability line look as they did. The 1000px breakpoint becomes `lg` and the 600px one `sm`, and the trial plan's border becomes `border-primary/65`, as the owner approved. The secondary plan's styles, which nothing renders, are removed. Where an element carries `smooth-primary-button` or `smooth-text-link`, the utilities that must beat SmoothUI take `!`, marked transitional.

**Blocked by:** 0150

**Status:** ready-for-agent

- [ ] `conversion.css` is deleted and nothing imports it.
- [ ] The plans section keeps its look at 1280px and 375px.
- [ ] The trial plan's border is `border-primary/65`.
- [ ] The secondary plan's styles are gone.
- [ ] Every `!` that beats SmoothUI carries a transitional note.
- [ ] The captures match the baseline, except for the trial border.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check tailwind`
- Built-in browser captures compared with the baseline.
