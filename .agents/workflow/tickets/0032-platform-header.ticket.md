# 0032: Platform header

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** the header keeps its look and changes its content: logo, locale pill, `Platform` (a glassy dropdown: Vision, FinOps, Microruntimes, AI Voice) and `Pricing`, then Docs, GitHub, a rainbow-highlighted `Book a Demo` button and the theme toggle. The drawer mirrors the nav.

**Blocked by:** 0028, 0030.

**Status:** ready-for-agent

### Copy

- [ ] `header.json` (five locales) and `HeaderMessages`: `platform`, `platformMenu`, `vision`, `finops`, `microruntimes`, `aiVoice`, `visionDescription`, `finopsDescription`, `microruntimesDescription`, `aiVoiceDescription`, `pricing`, `docs`, `github`, `bookDemo`; `product`, `attribution`, `integrations`, `voice`, `plans` and `getStarted` removed.

### Desktop

- [ ] One nav model: the `Platform` menu (Vision linked to `/platform/vision`; FinOps, Microruntimes and AI Voice as text with a description) and `Pricing` to `/pricing`.
- [ ] `website-header/platform-menu.tsx`: a disclosure button with a caret and a glassy Tailwind panel that opens on hover, focus and click and closes on Escape, outside click and navigation.
- [ ] Right side: Docs and GitHub as text, the demo button, the hamburger, the theme toggle.

### Demo button

- [ ] `website-header/demo-button.tsx` replaces `access-links.tsx`: a `span` without href, `bg-foreground text-background`, a rotating conic ring of `rainbow-1` to `rainbow-6` behind it, static under reduced motion.
- [ ] The `.website-get-started` rules leave `website-header.css`.

### Drawer

- [ ] The drawer lists a Platform heading, the four items, Pricing, Docs, GitHub, the language pill and the demo button.

### Verification

- `pnpm turbo run check`, typecheck and tests for `website`, `drift check localization`, `drift check dead-code`.
- At 1440px and above the dropdown works with mouse and keyboard; below, the drawer carries the same items; at 375px nothing overflows.
