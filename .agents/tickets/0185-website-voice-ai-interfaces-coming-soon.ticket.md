# 0185: Add Voice AI Interfaces to platform menu with coming soon badge

**What to build:** In the website navigation (both desktop `PlatformMenu` and mobile drawer in `WebsiteHeader`), add "Voice AI Interfaces" as a platform product with a discreet "Coming soon" badge. Add `VoiceAiInterfacesIcon` in `packages/web/src/ui/icons/platform`, update `platform-products.ts`, localize across all 5 locales in `header.json` and `messages.types.ts`, and add unit tests.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Create `VoiceAiInterfacesIcon` in `packages/web/src/ui/icons/platform/voice-ai-interfaces-icon.tsx` and export from `index.ts`.
- [x] Add unit test in `packages/web/src/ui/icons/platform/platform-icons.unit.tsx`.
- [x] Add `voiceAiInterfaces` and `voiceAiInterfacesDescription` to `apps/website/src/app/template/i18n/messages/header.json` across all 5 locales (`en`, `pt-BR`, `zh-Hant`, `ja`, `vi`).
- [x] Update `HeaderMessages` in `apps/website/src/app/template/i18n/messages/messages.types.ts`.
- [x] Add `voiceAiInterfaces` with `badge: 'comingSoon'` to `apps/website/src/app/lib/components/website-header/platform-products.ts`.
- [x] Update unit tests in `apps/website/src/app/lib/components/website-header/platform-menu.unit.tsx`.
- [x] Verify full check, typecheck, test, and build pass cleanly across the monorepo.
