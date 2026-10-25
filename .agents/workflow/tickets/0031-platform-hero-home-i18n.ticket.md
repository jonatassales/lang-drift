# 0031: Platform hero and home page

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** the home page becomes header, a platform hero over the `HeroShader`, and footer. The hero's copy is a new `platform` catalog in five locales, drafted for the owner's review.

**Blocked by:** 0029, 0030.

**Status:** ready-for-agent

### Copy

- [ ] `template/i18n/messages/platform.json` in five locales: `skip`, `hero.title` ("Agentic Platform"), `hero.lead`, `hero.primaryCta`, `hero.secondaryCta`, `hero.claims.{agents,evidence,budget}`; `PlatformMessages` in `messages.types.ts`; `platformMessages` exported from `messages/index.ts` behind the shape assertion; the family joins `catalogs.unit.ts`.

### Hero

- [ ] `lib/components/platform-hero.tsx` (server): `HeroShader` filling the section, a bottom fade to `bg-background`, the h1 "Agentic Platform", lead, a primary CTA to `/platform/vision`, a secondary CTA to `/platform/vision#early-access`, and three claims. Tailwind role utilities plus the existing `section-kicker`, `smooth-primary-button` and `smooth-text-link` classes only.

### Home

- [ ] `(pages)/page.tsx` renders `WebsiteStructuredData`, the skip link to `#top`, `WebsiteHeader`, `<main id="top">` with the hero, and `WebsiteFooter`; the unused imports are gone.

### Verification

- `drift check localization`, the catalogs test, `drift check dead-code` (every `.website-hero*` class is still carried by the Vision page).
- `/en/` shows header, hero and footer only; both CTAs reach `/en/platform/vision/`.
