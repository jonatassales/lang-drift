# 0027: Remove the code that supports missing features

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** nothing in the repository supports a feature, component, integration or service that tao does not have. The website copy brought code for a confirmation email the website no longer sends, a toast, the old header, pricing, account and comment components and unused dependencies; it all goes, and the documents describe only what exists. The website renders and behaves exactly as before. The owner keeps the `.gitkeep` placeholders and the `.fmt.ts` formatters.

**Blocked by:** none.

**Status:** ready-for-agent

### Confirmation email

The registry sends the confirmation; the website never sends mail (`apps/website/README.md`).

- [x] `early-access-email.json` is deleted, with its import, `assertCatalogShape` call and `earlyAccessEmailMessages` export in `template/i18n/messages/index.ts`, the `EarlyAccessEmail` type in `messages.types.ts` and its entry in `catalogs.unit.ts`.
- [x] `turbo.json` no longer passes `RESEND_*` through.
- [x] `early-access/env.server.unit.ts` is renamed after what it tests, the provider configuration of the route and of Siteverify, instead of an `env.server` module that does not exist.

### Feature flag and build environment

`WEBSITE_PREVIEW_SIGNUP_DISABLED` was a feature flag without a proper home for flags. The owner removes it and keeps sign-up paused, the state production and CI already ran in.

- [x] The variable is read nowhere: not in `next.config.ts`, the early-access section, the privacy page, the submission route, Siteverify, the tests, `turbo.json`, CI or the docs.
- [x] The paused behavior is live code: the section shows the return-soon notice, the privacy page states the pause and `/api/early-access/` answers 503 without calling a provider, which `api/early-access/route.unit.ts` proves.
- [x] The code that opens sign-up stays commented where it lives: the form branch of `section.tsx` and its test, the active privacy notice and the real `POST` export. The submission route and Siteverify stay live and tested, so opening sign-up is restoring those comments.
- [x] `next build` now reads no environment variable, so `website#build` and `globalEnv` go. CI sets no website values.

### CI

- [x] The Harness job leaves the pull-request workflow, and AGENTS.md no longer lists it.

### Toast

- [x] `atoms/toast.atom.ts` and its re-export from `lib/state/index.ts` are deleted.
- [x] The `.smooth-toast` rules in `smoothui.css` are deleted.

### CSS for components that do not exist

Each file is on the ADR 0003 exception list, which lets it shrink. Only rules whose selectors match no markup go.

- [x] `packages/web/src/ui/product-experience.css` (only `header-four*`) is deleted, with its import in `packages/web/src/styles.css` and its line in the exception list of `02-frontend.rule.md`.
- [x] `smoothui.css` keeps only `smooth-primary-button`, `smooth-text-link`, `smooth-footer`, `chroma-text` and `avatar-group`: the `header-four*`, `pricing-two`, `pricing-card`, `pricing-toggle`, `context-meter*`, `context-breakdown`, `figma-comment*`, `account-avatar*`, `ai-message*`, `smooth-scrubber*`, `glow-card*`, `smooth-icon-button`, `smooth-footer-main`, `smooth-footer-bottom` and `compact-chart` rules are gone.
- [x] `ui-consistency.css` is deleted, with its import and its line in the exception list. Its last live rules, `.avatar-chip` here and in `smoothui.css`, were touched, so they move to Tailwind classes in `animated-avatar-group.tsx`, their only consumer (`02-frontend.rule.md`, rule 5).
- [x] `globals.css` loses `header-four*`, `hero-product-card`, `public-auth-actions`, `public-sign-in`, `packaging-section`, `reason-*` and `voice-meter-card`.
- [x] `hero.css` loses `reason-impact`, `website-early-access` and `smooth-footer-bottom`.
- [x] `responsive.css` loses its `header-four*`, `hero-product-card`, `packaging-section` and `reason-*` rules.
- [x] `website-header.css` loses `website-locale-pill` and `website-locale-link`, including the drawer and reduced-motion selectors; the header uses `LanguageSwitcher` from `@langdrift/web`.

### Placeholders for missing integrations

- [x] `src/infra/better-stack/.gitkeep` stays: where there is a `.gitkeep`, the owner keeps the folder. The `@infra/*` alias goes with its Biome import groups, since nothing imports through it.
- [x] The `@metadata/*` alias goes with its Biome import groups: every metadata import uses `@app/metadata/`.
- [x] The empty `src/app/state/state.data.ts` is deleted; the state barrel stays `src/app/lib/state/index.ts`.

### Formatters and unused dependencies

- [x] The four files in `src/app/template/formatters/` stay: the owner keeps the `.fmt.ts` formatters.
- [x] `axios` and `tw-animate-css` leave the website's `package.json` and the lockfile; `clsx` and `tailwind-merge` stay, because `cn.fmt.ts` imports them.

### Found by the dead-code check

- [x] The empty `src/factories/index.ts` and the `src/services/index.ts` barrel go: every import names its module directly.
- [x] `afterPaint` leaves `@langdrift/web/utilities`: the web port brought it, and nothing calls it.

### Documents

- [x] `apps/website/.env.development` is tracked as the shared development defaults, empty for now and never holding a secret; `.gitignore` keeps ignoring every other `.env*` file, and AGENTS.md and the README say so.
- [x] AGENTS.md describes `packages/shell` as planned, not as where the switcher and the shell live.

### Owner

Agents do not edit `.env` files unless the owner asks.

- [x] At the owner's request, the local `apps/website/.env.development` is empty: none of the dashboard, SSO, Turnstile or Resend values it held exist in tao yet, and sign-up is paused in code.

### Verification

- [x] `pnpm turbo run build`, `typecheck`, `test` and `check` pass, and `drift check localization` passes after the catalog change.
- [x] Every page in every locale looks the same before and after, in light and dark: no removed rule was live.
- [x] No reference to a removed file, alias, export, variable or class remains.

Parity was proved by comparing every computed property, elements and `::before`/`::after`, with main's CSS served and with the pruned CSS: `/en/` at 1280 and 375 px in dark and at 1280 px in light, `/en/privacy/` and `/en/early-access-sent/` at 1280 px. Nothing differed outside dev-injected scripts. The other locales render the same markup, and no removed selector names a theme or locale.

Intended survivors: the renamed test still stubs `RESEND_*` to prove the route needs none, and `earlyAccessEmailSchema` validates the form's email field.
