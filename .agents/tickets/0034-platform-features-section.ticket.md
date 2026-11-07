# 0034: Platform features section

**Spec:** .agents/specs/0004-platform-landing-hero.spec.md

**What to build:** the home page's second section, after the hero: the four platform products as cards under a heading, adapted from SmoothUI's Features Grid block (smoothui.dev/docs/blocks/features).

**Blocked by:** 0031.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Component

- [ ] `packages/web/src/vendors/smoothui/features-grid.tsx`, exported as `@langdrift/web/vendors/smoothui/features-grid`: a heading, a lead and a grid of cards (`sm:grid-cols-2 lg:grid-cols-4`) that rise into view one after another and lift on hover on pointer devices; reduced motion shows them still. The caller passes the content; colors are role utilities only.

### Section

- [ ] `apps/website/src/app/lib/components/platform-features.tsx` (server) composes it with the four products from the header's nav model: their icons, names and one-line descriptions from `header.json`.
- [ ] `platform.json` gains `features.title` ("Four products. One platform.") and `features.lead` in five locales; `PlatformMessages` gains the type.
- [ ] The home page renders the section after the hero.

### Verification

- `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `pnpm turbo run check`, `drift check localization`, `drift check dead-code`.
- The cards fade in on scroll in both color schemes and stack at 375px.
