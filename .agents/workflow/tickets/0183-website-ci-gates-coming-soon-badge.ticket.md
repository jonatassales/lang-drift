# 0183: Add Coming Soon badge to website CI Gates menu item

**What to build:** Add a discreet "Coming soon" badge to the CI Gates menu item in the website navigation (both desktop `PlatformMenu` and mobile drawer in `WebsiteHeader`). Localize the badge copy across all 5 locales in `apps/website/src/app/template/i18n/messages/header.json` (`en`: "Coming soon", `pt-BR`: "Em breve", `zh-Hant`: "即將推出", `ja`: "近日公開", `vi`: "Sắp ra mắt") and update `HeaderMessages` in `messages.types.ts`. Create a dedicated component `apps/website/src/app/lib/components/website-header/product-badge.tsx` with unit tests in `product-badge.unit.tsx`. Update `platform-products.ts`, `platform-menu.tsx`, and `website-header.tsx`.

**Blocked by:** none

**Status:** ready-for-agent

- [x] Add `comingSoon` translation keys to `header.json` for all 5 locales (`en`, `pt-BR`, `zh-Hant`, `ja`, `vi`).
- [x] Update `HeaderMessages` interface in `apps/website/src/app/template/i18n/messages/messages.types.ts`.
- [x] Add `badge: 'comingSoon'` to `ciGates` in `apps/website/src/app/lib/components/website-header/platform-products.ts`.
- [x] Create `apps/website/src/app/lib/components/website-header/product-badge.tsx` with discreet badge styling.
- [x] Create unit tests in `apps/website/src/app/lib/components/website-header/product-badge.unit.tsx`.
- [x] Render the badge in `PlatformMenu` (`platform-menu.tsx`) and mobile drawer (`website-header.tsx`).
- [x] Verify test suite, typecheck, lint, localization checks, and build pass cleanly.
