# 0033: Logo mark and product icons

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** LangDrift gets a new logo mark and the four platform products get monochrome icons, both in `@langdrift/web`. The hero's centerpiece, in the middle like Vercel's triangle, is not the mark but a black hole sun: a disc in the page background with a thin bright rim and a corona in the brand purple.

**Blocked by:** 0028.

**Status:** ready-for-agent

### Mark

- [ ] `LogoMark` in `packages/web/src/ui/brand/brand.tsx` draws the drift mark: a ring with one segment carried outward, both strokes in the foreground color so it reads black on light and white on dark. It is square (`viewBox 0 0 64 64`), and `Brand`, the footer and the integrations section keep using it.
- [ ] The hero's middle column shows the black hole sun, built from Tailwind role utilities only; its corona breathes with `animate-pulse` and stays still under reduced motion.

### Icons

- [ ] `packages/web/src/ui/icons/platform/`: `VisionIcon`, `FinopsIcon`, `MicroruntimesIcon` and `AiVoiceIcon`, 24-unit stroke icons in `currentColor` with a `className` prop, exported as `@langdrift/web/ui/icons/platform`.
- [ ] The Platform menu and the drawer show each product's icon beside its name.

### Verification

- `pnpm turbo run typecheck test --filter=@langdrift/web --filter=website`, `pnpm turbo run check`, `drift check dead-code`.
- The header, footer and hero show the new mark in both color schemes; the dropdown shows the four icons.
