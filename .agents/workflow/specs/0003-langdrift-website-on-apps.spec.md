# LangDrift website on apps/

**Status:** ready-for-agent

## Problem Statement

The repository was laid out for products (`products/<product>/<kind>/<workspace>`), but no product workspace was ever created and Vision and FinOps are no longer planned. Meanwhile the LangDrift marketing website already exists in the earlier `langdrift` monorepo, built with Next.js, next-intl, Jotai and Tailwind, together with the guardrails written for it. The owner wants that website here, unchanged in architecture and look, recolored to the Exphere purple, and wants this repository's structure (workspace layout, design tokens, rules, tooling) to stay the source of truth around it.

## Solution

The repository moves to `apps/*` and the website becomes its first and only app, copied from langdrift. Everything the website needs from outside its folder is rebuilt on this repository's terms:

- `@langdrift/tokens` carries the website's look in shadcn's vocabulary: a single `website` theme in Exphere purple, the website's neutrals as shared roles, and two new fixed color groups, classification and drift colors.
- A new `@langdrift/web` package holds the components the website imports from langdrift's `@pack/react`, and nothing else.
- The website adapts only its token layer and its configuration; its routes, components, state, content and hand-written CSS stay as they were.
- The guardrails langdrift wrote for React, Next.js, Jotai, Tailwind, localization and SEO become rules here, with a Biome lint upgrade and an executable localization check.

## User Stories

1. As the owner, I want workspaces to live under `apps/` and `packages/`, so that the layout matches what actually exists.
2. As the owner, I want the empty `products/` folders gone, so that nobody builds on an abandoned layout.
3. As the owner, I want a commit that changes `apps/` to need a ticket trailer, so that the guardrail covering `products/` keeps its meaning.
4. As the owner, I want `drift doctor` to check for `apps` instead of `products`, so that CI keeps passing after the move.
5. As the owner, I want AGENTS.md and the glossary to describe apps instead of products, so that agents use the current vocabulary.
6. As a visitor, I want the website's pages, copy and demos exactly as on langdrift, so that the site I know is the site I get.
7. As a visitor, I want the brand color to be the Exphere purple, so that the website carries the new identity.
8. As a visitor, I want dark to stay the default color scheme and my choice of scheme to persist, so that the site behaves as before.
9. As a visitor, I want the five locales (en, pt-BR, zh-Hant, ja, vi) with their locale prefixes, so that I can read the site in my language.
10. As a visitor whose browser asks for Simplified Chinese, I want to be sent to English with an availability notice, so that I am never shown Traditional Chinese as a substitute.
11. As a visitor, I want a mixed-case locale prefix to redirect to lowercase, so that every page has one URL.
12. As a visitor, I want my locale preference remembered, so that I return to the language I chose.
13. As a visitor, I want the early-access form to behave as before, including its paused state, so that sign-up keeps its rules.
14. As a search engine, I want the same metadata, hreflang, sitemap, robots and social images, so that the site's indexing does not regress.
15. As a frontend developer, I want the website's components to read colors only from `@langdrift/tokens` roles, so that one package owns every color.
16. As a frontend developer, I want classification and drift colors as tokens, so that Product Vision charts never hard-code them.
17. As a frontend developer, I want `dark:` utilities to follow the `.dark` class, so that shadcn and SmoothUI components switch with the site.
18. As a frontend developer, I want the components the website uses in `@langdrift/web`, so that shared UI lives in a package, not in the app.
19. As a frontend developer, I want `@langdrift/web` to export only what the website imports, so that no unused component is maintained.
20. As a frontend developer, I want the website to import `@langdrift/web` at the subpaths it used for `@pack/react`, so that the website names packages the way the rest of the repository does.
21. As a frontend developer, I want the website formatted in this repository's Biome style, so that one style applies everywhere.
22. As a frontend developer, I want the website's path aliases grouped as aliases by Biome, so that import order stays stable.
23. As a frontend developer, I want the website's unit tests and `@langdrift/web`'s tests to run under `turbo run test`, so that they gate pull requests.
24. As a frontend developer, I want `turbo run dev --filter=website` to build the tokens first, so that dev starts without a manual step.
25. As a frontend developer, I want rules for React, Next.js and Jotai, so that new code follows the same conventions as the copied code.
26. As a frontend developer, I want the Tailwind rules extended with `cn()`, role utilities only and no component classes, so that styling stays consistent.
27. As a frontend developer, I want the website's existing hand-written CSS allowed by a scoped exception that forbids new CSS, so that the copy is legal without opening the door to more.
28. As a frontend developer, I want rules for localization and SEO, so that every new string and page meets the bar the website set.
29. As a frontend developer, I want code rules for environment access and control flow, so that the copied code's style is the house style.
30. As a frontend developer, I want Biome to require interfaces, kebab-case file names and no unused imports, so that these rules are enforced, not just written.
31. As a frontend developer, I want `drift check localization` to fail on catalog drift, so that a missing translation breaks CI instead of production.
32. As the owner, I want CI to build the website with safe values (early access paused, sign-up disabled), so that CI never touches real services.
33. As the owner, I want the decisions behind this move recorded in an ADR, so that the CSS exception and the token additions can be revisited knowingly.

## Implementation Decisions

- **Layout.** The pnpm workspace globs are `apps/*` and `packages/*`. `products/` is deleted. The ticket-trailer guardrail matches `apps/` and `packages/`. `drift doctor` requires `apps`.
- **One theme.** `@langdrift/tokens` builds a single stylesheet, `website.css`, from the shared tokens plus a `website` theme with its own prefix (`--site-*`). The Vision and FinOps themes and stylesheets are removed. The theme mechanism stays, so a second app can add a theme later.
- **Palette.** Primary is the Exphere purple: `#6D0AD6` in both color schemes; hover `#5408A5` in light and `#8C3BEF` in dark. Primary foreground is white in both schemes.
- **Shared roles take the website's neutrals.** Background, card, popover, foreground, muted foreground and border take the values langdrift's website declared for itself; the other roles keep their reference chains.
- **New fixed color groups.** Classification colors (intentional, aligned, unexplained, review), their text tones (intentional, review, unexplained; darker in light, lighter in dark, because the palette reset removes the Tailwind tones the website used for small text) and drift colors (negative, positive, zero; light and dark) join status colors as groups that never take the brand. Each is exposed as a Tailwind color.
- **Fonts.** The body font becomes the website's editorial system stack, and a mono font token is added, each mapped to Tailwind's font variable.
- **Dark variant.** The tokens stylesheet declares Tailwind's `dark` variant on the `.dark` class, so every app gets it from the bridge.
- **`@langdrift/web`.** Ships TypeScript source with an exports map limited to what the website imports from `@pack/react`, at the same subpaths. Its components read colors through the token roles. Its stylesheet keeps its own Tailwind entry and its existing CSS files, and imports the tokens stylesheet.
- **Website adaptation boundary.** In the website's source only the token layer changes: `--ld-*` variables and utilities are remapped onto roles and the new groups, `data-theme` becomes the `.dark` class, and the website's own theme files are deleted. Its imports name `@langdrift/web` instead of `@pack/react`, at the same subpaths. Configuration adapts to the repository: the manifest depends on `@langdrift/web`, the test-fixture alias points inside the app, and the Next.js transpile list names `@langdrift/web`.
- **Accepted look deltas.** Four light-mode values collapse onto the nearest role: chart grid onto border, chart line and point onto foreground, langdrift's text color onto foreground, and the product font onto the sans font.
- **Integrations unchanged.** The early-access flow (Turnstile, Resend, registry), the locale proxy and the environment contract are copied as they are.
- **A new project.** tao keeps nothing for langdrift's sake: the website has no Docs zone (no `DOCS_ORIGIN`, no Docs rewrites or redirect, no Docs entries in the sitemap), no container image, and its README describes this repository only.
- **Formatting.** The copy is reformatted once to this repository's Biome style; object literals the multi-property plugin reports are expanded by hand.
- **Harness.** New rules cover code (environment access, control flow), React, Next.js, Jotai, localization and SEO. The frontend and Tailwind rules gain langdrift's Tailwind conventions and a scoped exception for the CSS files the copy brings. Biome gains langdrift's lint rules. `drift` gains a `check localization` command that CI's Harness job runs.
- **Build wiring.** Turborepo declares the website's build environment, and `dev` depends on upstream builds. CI sets the safe values langdrift used.

## Testing Decisions

- A good test drives a module through its public surface: the built stylesheet a consumer imports, the component a page renders, the response the proxy returns. It never asserts on internals.
- **Seam 1: the tokens contract test.** The existing `tokens.test.js` builds the package in isolation and resolves its stylesheet through the package exports. It gains the website theme, the classification and drift groups (exposed and never brand), the mono font and the dark variant, and loses the cross-product isolation test.
- **Seam 2: the copied unit tests.** The website's 14 unit tests (proxy, i18n catalogs, early-access flow, voice preview) and `@pack/react`'s tests move unchanged apart from formatting, and run under `turbo run test`. They prove the copy behaves like the original.
- **Seam 3: `drift check localization`** runs against the copied catalogs in CI.
- Visual parity is checked by hand in the browser pane against langdrift's running website; it is not automated.

## Out of Scope

- Converting the website's hand-written CSS to Tailwind.
- Removing code the copy brings but does not use (formatters, the toast atom, the `axios` and `tw-animate-css` dependencies).
- Other langdrift apps (dashboard, docs, sso) and packages beyond what the website imports.
- langdrift's test-ownership rule and its third-party agent skills.
- Vercel project configuration and production secrets.
- A Docs zone and a container image.

## Further Notes

- The source is langdrift `main` at `e518dac`.
- ADR 0003 records why products were dropped, why the website's values live in the tokens and why its CSS is excepted.
