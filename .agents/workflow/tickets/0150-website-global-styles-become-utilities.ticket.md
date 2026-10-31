# 0150: The website's global and base styles become utilities

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** The website's global stylesheets leave: `globals.css`, `hero.css`'s scroll padding and `template/tailwind/base.css`. `<html>` and `<body>` in the locale layout carry what were element rules (smooth scrolling, the 96px scroll padding, the colors, the minimum height and the 280ms theme transition, with `motion-reduce:` turning motion off). The global focus ring, `white-space: nowrap` on buttons and `font: inherit` on buttons and inputs become classes on every element that relied on them, including `LanguageSwitcher`, `AnimatedTabs` and `BasicAccordion` in `packages/web`. The `section-kicker` class becomes utilities where it is used. Font utilities that the global `font: inherit` made ineffective on buttons and inputs are removed. The Tailwind entry keeps only its imports.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `globals.css` and `template/tailwind/base.css` are deleted, and `hero.css` keeps no `html` rule.
- [ ] `app.css` holds only the Tailwind import, the tokens stylesheet and the utilities file.
- [ ] `<html>` and `<body>` carry the former element rules as utilities, reduced motion included.
- [ ] Every link, button and input that showed the global focus ring shows the same ring through its own classes.
- [ ] Every button that relied on the global `nowrap` keeps it through `whitespace-nowrap`.
- [ ] No font utility on a button or an input starts to apply because the global `font: inherit` is gone.
- [ ] The three `packages/web` components keep their look on the website and in the console.
- [ ] The captures at 1280px and 375px, light and dark, match the baseline.

### Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=@langdrift/web`
- `platform/cli/drift quality check tailwind`
- Built-in browser captures compared with the baseline.
