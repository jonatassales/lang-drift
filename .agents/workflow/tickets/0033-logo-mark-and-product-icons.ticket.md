# 0033: Logo mark and product icons

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** LangDrift gets a new logo mark and the four platform products get monochrome icons, both in `@langdrift/web`. The hero's centerpiece, in the middle like Vercel's triangle, is not the mark but a black hole sun: a disc in the page background with a thin bright rim and a corona in the brand purple.

**Blocked by:** 0028.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Mark

- [ ] `LogoMark` in `packages/web/src/ui/brand/brand.tsx` draws the drift mark: a ring with one segment carried outward, both strokes in the foreground color so it reads black on light and white on dark. It is square (`viewBox 0 0 64 64`), and `Brand`, the footer and the integrations section keep using it.
- [ ] The hero's middle column shows the black hole sun, built from Tailwind role utilities only; its corona breathes with `animate-pulse` and stays still under reduced motion.

- [ ] The raster icons carry the new mark instead of the orange swoosh: `favicon.ico` (16 to 256), `icon.png`, `apple-icon.png`, and `public/metadata/icon-192.png`, `icon-512.png` and `logo-512.png`. The mark is white (#F5F5F5) on the dark background (#0B0B0C), the dark scheme's foreground and background, so it reads on light and dark browser tabs; the Apple icon is full-bleed because iOS rounds it.

- [ ] The social share images (`opengraph-image.png` and `twitter-image.png`, 1200×630) show the drift mark and wordmark, the headline "Agentic Platform", the black hole sun with its purple corona, and the four products, in the dark scheme's tokens.

### Icons

- [ ] `packages/web/src/ui/icons/platform/`: `VisionIcon`, `FinopsIcon`, `MicroruntimesIcon` and `AiVoiceIcon`, 24-unit stroke icons in `currentColor` with a `className` prop, exported as `@langdrift/web/ui/icons/platform`.
- [ ] The Platform menu and the drawer show each product's icon beside its name.

### Verification

- `pnpm turbo run typecheck test --filter=@langdrift/web --filter=website`, `pnpm turbo run check`, `drift check dead-code`.
- The header, footer and hero show the new mark in both color schemes; the dropdown shows the four icons.
