# 0032: Platform header

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** the header keeps its look and changes its content: logo and locale pill on the left, `Platform` (a glassy dropdown: Vision, FinOps, Microruntimes, AI Voice), Docs, GitHub and `Pricing` together in the middle, then a rainbow-highlighted `Book a Demo` button and the theme toggle on the right. The bar is glassier, the locale pill shows at every width, and the drawer mirrors the nav.

**Blocked by:** 0028, 0030.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Copy

- [ ] `header.json` (five locales) and `HeaderMessages`: `platform`, `platformMenu`, `vision`, `finops`, `microruntimes`, `aiVoice`, `visionDescription`, `finopsDescription`, `microruntimesDescription`, `aiVoiceDescription`, `pricing`, `docs`, `github`, `bookDemo`; `product`, `attribution`, `integrations`, `voice`, `plans` and `getStarted` removed.

### Desktop

- [ ] One nav model: the `Platform` menu (Vision linked to `/platform/vision`; FinOps, Microruntimes and AI Voice as text with a description and an icon), then Docs and GitHub as text, then `Pricing` to `/pricing`, all in the center nav.
- [ ] The desktop nav shows from Tailwind's `lg` breakpoint (64rem) and the hamburger below it; the breakpoint lives in Tailwind classes, not in `website-header.css`, and the drawer's auto-close uses the same `64rem`.
- [ ] The Platform dropdown and the drawer share one glass surface (`website-header/surface.ts`): radius, transparency, blur, border and shadow.
- [ ] `website-header/platform-menu.tsx`: a disclosure button with a caret and a glassy Tailwind panel that opens on hover, focus and click and closes on Escape, outside click and navigation.
- [ ] Right side: the demo button, the hamburger, the theme toggle. The glass moves from `website-header.css` to Tailwind on the bar, with more blur and less opacity, and the rule hiding the locale pill below 760px leaves the CSS.

### Demo button

- [ ] `website-header/demo-button.tsx` replaces `access-links.tsx`: a `span` without href, `bg-foreground text-background`, a rotating conic ring of `rainbow-1` to `rainbow-6` behind it, static under reduced motion.
- [ ] The `.website-get-started` rules leave `website-header.css`.

### Drawer

- [ ] The drawer lists a Platform heading, the four items with their icons, Docs, GitHub, Pricing, the language pill and the demo button.

### Verification

- `pnpm turbo run check`, typecheck and tests for `website`, `drift check localization`, `drift check dead-code`.
- At 1440px and above the dropdown works with mouse and keyboard; below, the drawer carries the same items; at 375px nothing overflows.
