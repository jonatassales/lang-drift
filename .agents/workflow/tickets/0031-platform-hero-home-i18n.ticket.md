# 0031: Platform hero and home page

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** the home page becomes header, a platform hero over the `HeroShader`, and footer. The hero's copy is a new `platform` catalog in five locales, drafted for the owner's review.

**Blocked by:** 0029, 0030.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Copy

- [ ] `template/i18n/messages/platform.json` in five locales: `skip`, `hero.titleStart` ("Agentic"), `hero.words` (DevEx, FinOps, Runtime, Voice AI), `hero.titleEnd` ("Platform"), `hero.primaryCta` ("Explore the products"), `hero.secondaryCta` ("Book a Demo"); `PlatformMessages` in `messages.types.ts`; `platformMessages` exported from `messages/index.ts` behind the shape assertion; the family joins `catalogs.unit.ts`.

### Hero

- [ ] `lib/components/platform-hero.tsx` (server): `HeroShader` filling the section, a bottom fade to `bg-background`, the h1 on the left in two lines, "Agentic" and then a kinetic word followed by "Platform" that never wraps (its size follows the column, so "Runtime Platform", the widest, fits from 320px up), the word cycling through DevEx, FinOps, Runtime and Voice AI with SmoothUI's Text Morph (vendored as `@langdrift/web/vendors/smoothui/text-morph`, driven by `@langdrift/web/ui/rotating-text`), a black hole sun in the middle, a black-and-white pill to `/platform/vision` (foreground on background, like Vercel's primary button) and a text link to `/platform/vision#early-access`; no lead and no claims. Tailwind role utilities plus the existing `smooth-text-link` class only.
- [ ] The sun's disc holds still; the smoke behind it is a `data-magnetic` participant of SmoothUI's Magnetic Field, vendored as `@langdrift/web/vendors/smoothui/magnetic-field` with the `motion` dependency, so it leans and turns toward the pointer and brightens while the pointer is over the sun. Reduced motion and touch leave it at rest.
- [ ] Inside the disc, and only there, SmoothUI's Gravity Stars (vendored as `@langdrift/web/vendors/smoothui/gravity-stars`) drifts a starfield that falls toward the pointer. The disc is dark in both color schemes and the stars take its text color, so they stay light; the palette is re-read when `<html class>` changes. The Magnetic Field's loop sleeps once everything is at rest.
- [ ] On the right, the four products in a vertical list named "Platform products": dimmed at rest but readable, each lighting up on hover or focus with a purple dot and a soft glow; Vision links to its page.

### Home

- [ ] `(pages)/page.tsx` renders `WebsiteStructuredData`, the skip link to `#top`, `WebsiteHeader`, `<main id="top">` with the hero, and `WebsiteFooter`; the unused imports are gone.

### Verification

- `drift check localization`, the catalogs test, `drift check dead-code` (every `.website-hero*` class is still carried by the Vision page).
- `/en/` shows header, hero and footer only; both CTAs reach `/en/platform/vision/`.
