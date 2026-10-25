# 0031: Platform hero and home page

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** the home page becomes header, a platform hero over the `HeroShader`, and footer. The hero's copy is a new `platform` catalog in five locales, drafted for the owner's review.

**Blocked by:** 0029, 0030.

**Status:** ready-for-agent

### Copy

- [ ] `template/i18n/messages/platform.json` in five locales: `skip`, `hero.title` ("Agentic Platform"), `hero.primaryCta` ("Explore the products"), `hero.secondaryCta` ("Book a Demo"); `PlatformMessages` in `messages.types.ts`; `platformMessages` exported from `messages/index.ts` behind the shape assertion; the family joins `catalogs.unit.ts`.

### Hero

- [ ] `lib/components/platform-hero.tsx` (server): `HeroShader` filling the section, a bottom fade to `bg-background`, the h1 "Agentic Platform" on the left, a black hole sun in the middle, a black-and-white pill to `/platform/vision` (foreground on background, like Vercel's primary button) and a text link to `/platform/vision#early-access`; no lead and no claims. Tailwind role utilities plus the existing `smooth-text-link` class only.
- [ ] The sun's disc holds still; the smoke behind it is a `data-magnetic` participant of SmoothUI's Magnetic Field, vendored as `@langdrift/web/vendors/smoothui/magnetic-field` with the `motion` dependency, so it leans and turns toward the pointer and brightens while the pointer is over the sun. Reduced motion and touch leave it at rest.

### Home

- [ ] `(pages)/page.tsx` renders `WebsiteStructuredData`, the skip link to `#top`, `WebsiteHeader`, `<main id="top">` with the hero, and `WebsiteFooter`; the unused imports are gone.

### Verification

- `drift check localization`, the catalogs test, `drift check dead-code` (every `.website-hero*` class is still carried by the Vision page).
- `/en/` shows header, hero and footer only; both CTAs reach `/en/platform/vision/`.
