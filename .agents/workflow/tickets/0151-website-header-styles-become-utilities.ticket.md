# 0151: The website header's styles become utilities

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** `website-header.css` becomes utilities in the header, its drawer and its controls, and is deleted. The sticky bar, its border once the page scrolls, its icons, the desktop navigation, the drawer and its access area look as they did. The 760px breakpoint becomes `md` and the drawer's scrim becomes `bg-foreground/38`, as the owner approved. `DemoButton` accepts a `className` so the drawer can give it its 44px minimum height.

**Blocked by:** 0150

**Status:** ready-for-agent

- [ ] `website-header.css` is deleted and nothing imports it.
- [ ] The header, the drawer and the Platform menu keep their look and behavior, reduced motion included.
- [ ] Below `md` the demo button is hidden and the bar tightens, as it did below 760px.
- [ ] The drawer's scrim is `bg-foreground/38` with the same blur.
- [ ] `DemoButton` takes a `className`, and the drawer passes the 44px minimum height through it.
- [ ] The captures match the baseline, the open drawer included, except for the scrim.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check tailwind`
- Built-in browser captures compared with the baseline.
