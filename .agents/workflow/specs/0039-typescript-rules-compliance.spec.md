# The TypeScript code complies with the TypeScript rules

**Status:** ready-for-agent

> **Note, 2026-11-05.** The owner closed 0277, 0282, 0283 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

## Tickets

- [ ] 0277: drift measures the TypeScript rules
- [ ] 0278: every import past a sibling uses an alias
- [ ] 0279: the messages barrel only re-exports
- [ ] 0280: every translator is named translate
- [ ] 0281: the website files carry the suffix of their responsibility
- [ ] 0282: the console files carry the suffix of their responsibility
- [ ] 0283: the packages files carry the suffix of their responsibility
- [ ] 0284: every object type in a signature is a named interface
- [ ] 0285: every await statement has its blank lines, in tests too

## Problem Statement

A read-only audit on 2026-11-04 measured `.claude/rules/04-typescript.rule.md` against every TypeScript file under `apps/` and `packages/` and found it applied at 91.4%. Biome's formatter and linter are clean, but the rules Biome does not carry, and one rule Biome carries with a wrong configuration, are not:

- **Import groups (Imports 4 and 6).** Spec 0044 renamed the internal packages to `web`, `tokens` and `streaming` and put them in the alias group of the import-group matcher in `biome.json`, so the tree already sorts them there. Nine files still break the groups by hand: a blank line inside a group (`apps/microservices/demo-requests/src/main.ts`, `apps/website/src/services/crypto.service.ts`) or a stylesheet side-effect import of the alias group placed among the third-party imports (`apps/console/src/app/[locale]/layout.tsx:8-9`, `apps/website/src/app/[locale]/layout.tsx:13-14`). `drift check biome` verifies the matcher with its fixtures.
- **Deep relative imports (Imports 1).** 26 imports reach past a sibling with `./dir/...`: the two proxies import their implementation folders (`apps/console/src/proxy.ts:7-10`, `apps/website/src/proxy.ts:3`), the book-a-demo route imports its interceptors folder, the website header imports its own folder, and the i18n request configurations and a catalog test import `./messages/...`. Neither app has an alias for `src/proxy/`.
- **Barrels (Imports 3).** The website messages barrel (`apps/website/src/app/template/i18n/messages/index.ts:1`) opens with a side-effect import that runs the catalog-shape assertions at load.
- **File suffixes.** 90 of 315 non-test, non-component, non-framework TypeScript files carry no suffix from the table or one outside it: 32 `page.metadata.ts`, `layout.metadata.ts` and `not-found.metadata.ts` companions that export a route's `generateMetadata`, one `layout.static.ts`, 7 files with `.validator`, `.response`, `.submission` or `.body`, and 51 unsuffixed files (`apps/console/src/app/lib/data/console-error.ts`, `apps/website/src/proxy/locale.ts`, the `packages/web` hero-shader helpers spec 0030 created). Some suffixed files do not match their suffix: `.compute` files that await a network read (`apps/console/src/app/lib/data/entity-title/*-title.compute.ts`, `apps/website/src/app/metadata/social-images.compute.ts`), `.types` files that export constant arrays (`apps/console/src/app/lib/views/tickets/ticket-phases.types.ts`, `apps/website/src/app/lib/components/book-a-demo/form.types.ts`) and a `.data` file that exports a function (`apps/website/src/app/template/i18n/loaders/website-messages.data.ts`). Spec 0027 placed the metadata companions beside their routes and spec 0030 split the hero-shader helpers without a suffix; the rule outranks both.
- **Translator names (Bindings 2).** 7 translators are named `t` (six console product pages such as `apps/console/src/app/[locale]/(menu)/products/[id]/analysis/page.tsx:20`, and `apps/console/src/app/[locale]/not-found.tsx:8`) and 6 carry another name (`messages`, `status`, `causes`, `kinds`, `ui`). The existing detector only sees `const t = `.
- **Inline object types (Types and contracts 2).** 29 signatures, generic arguments and assertions spell an object type inline: console read envelopes such as `useResource<{ products: Product[] }>` (17 sites), a private return type in `apps/website/src/services/cookie.service.ts:265`, the locale prefix table in `apps/website/src/proxy/locale.ts:12`, test helpers and the three `vitest.config.mts` files that assert the paths configuration.
- **Await spacing (Statement spacing).** 91 `await` statements in 15 files lack a blank line before or after them, nearly all in tests (`apps/website/src/app/lib/components/early-access/form.unit.tsx` alone has 23), plus one `await` that shares its statement with optional chaining (`form.unit.tsx:66`). The await-spacing detector skips tests, so `drift quality check` never saw them.

The try/catch blocks around pure code, the environment reader modules and the invented development secret belong to specs 0034 and 0035 and the Nest compliance spec, not yet written; the `packages/tokens/build.js` script belongs to spec 0037.

## Solution

Every TypeScript file under `apps/` and `packages/` follows the TypeScript rules, and `drift quality check 04-typescript` reports zero findings and measures every class of violation it can. The Biome matcher names the workspace packages and the tree is re-sorted in one mechanical pass. Every import past a sibling uses an alias, including a new `@proxy/*` alias in both Next apps. The messages barrel only re-exports, and the catalog-shape guard runs where the catalogs are assembled. Every file carries the suffix of its responsibility, chosen from the table; the calls the table does not settle are recorded below. Every translator is `translate`. Every object type in a signature, generic argument or assertion is a named interface. Every `await` statement in source and tests has its blank lines. Behaviour does not change: the unit tests, `typecheck`, `pnpm turbo run check` and `drift check localization` stay green at every ticket.

## User Stories

1. As the platform owner, I want `drift quality check 04-typescript` to report zero findings, so that compliance with the TypeScript rules is measured rather than asserted.
2. As the platform owner, I want a detector for every class of violation a pattern or awk program can catch reliably, so that a regression fails the check instead of waiting for a review.
3. As a developer, I want the workspace packages and every alias sorted in the alias group, with no blank line inside a group, so that the three import groups mean packages, workspace, and siblings everywhere.
4. As a developer, I want `drift check biome` to prove that the matcher keeps the workspace packages in the alias group, so that a Biome upgrade or a config edit cannot silently undo it.
5. As a developer, I want a relative import to reach only a sibling, so that a path tells me whether a dependency is local.
6. As a developer, I want a `@proxy/*` alias in each Next app, so that the proxy imports its implementation the way every other module imports across a boundary.
7. As a developer, I want the messages barrel to hold only named re-exports, so that importing a message never runs a hidden assertion.
8. As a developer, I want the catalog-shape guard to run where the website assembles its catalogs, so that the guard stays explicit and the runtime behaviour stays.
9. As a developer, I want every file's suffix to name its responsibility, so that I know what a file holds before I open it.
10. As a developer, I want a route's metadata companion to be named after the framework file it serves and suffixed as a mapping, so that a page folder reads `page.tsx` beside `page-metadata.map.ts`.
11. As a developer, I want a `.compute` file to be pure, so that the suffix promises what it names.
12. As a developer, I want a `.types` file to hold only types and a `.data` file only data, so that constants and loaders live under their own suffix.
13. As a developer, I want a re-export-only module to be an `index.ts` barrel or to disappear, so that a barrel is recognisable by its name.
14. As a developer, I want every translator to be `translate`, so that a message lookup reads the same in every component.
15. As a developer, I want a component with two translators to name each after its namespace with the `translate` prefix, so that the rule holds without hiding which catalog a key comes from.
16. As a developer, I want every object type in a signature, generic argument or assertion to be a named interface above its use, so that a shape has a name I can search for.
17. As a developer, I want the console read envelopes to be named interfaces owned by the data layer, so that a response shape is declared once.
18. As a developer, I want every `await` statement in a test to have its blank lines, so that tests follow the same spacing as the source.
19. As a developer, I want an awaited value consumed by optional chaining to land in a binding first, so that the failure mode is explicit.
20. As a reviewer, I want each ticket to be a batch that lands green on its own, so that a wide rename never leaves the tree broken.
21. As the owner, I want every reading of the suffix table that the table does not settle recorded, so that I can overturn it at pull request review.

## Implementation Decisions

- **Biome matcher.** The import-group matcher already places the workspace packages in the alias group (spec 0044); it gains `@proxy/*` with the alias ticket below. The nine hand-placed breaks are fixed by hand: the stylesheet side-effect imports of the alias group sit in the alias group with one blank line before them, and the blank line inside a group is removed. The `drift check biome` fixtures expect `web` in the alias group.
- **Detectors.** `drift quality check` gains, under `04-typescript`: a translator-name detector (awk) that reports any binding of `useTranslations` or `getTranslations`, including destructured bindings from `Promise.all`, whose name is neither `translate` nor `translate<Namespace>`; a deep-relative-import detector (pattern) for `./dir/...` sources; an inline-object-type detector (pattern) for an object literal type as a generic argument, after `as`, or as a return type annotation; a file-suffix detector (awk over tracked paths) that reports a `.ts` file under an app or package `src/` whose suffix is not in the TypeScript table, the Nest table or the framework's file conventions, and is not a test, an `index.ts` barrel, a configuration file or a declaration file. The await-spacing and import-group detectors move to the `all` scope so tests are measured. Inline object types in parameter annotations stay review-enforced, since a pattern cannot tell a type annotation from a destructured default. The existing `const t =` detector is replaced by the awk one.
- **Aliases.** Both Next apps declare `@proxy/*` for `src/proxy/*` in their paths configuration; the Biome matcher and the `drift check biome` fixtures name it. The proxies import their implementation through it. The book-a-demo route imports its interceptors through `@app/api/...`, the website header its folder through the alias spec 0041 gives the locale layout, and the i18n request configurations and tests their catalogs through `@i18n/messages/...`. No other alias is added; the file placement itself is spec 0041's.
- **Messages barrel.** The barrel loses its side-effect import. The catalog-shape assertion function keeps its behaviour and becomes a pure computation module; the website messages loader calls it for each catalog it assembles, so the guard still runs at startup, explicitly. The parity test of the catalogs keeps enforcing the rule (localization rules, item 7).
- **Suffix calls, decided from the table.** A route's `generateMetadata` companion maps route params and the parent metadata to a `Metadata` object and is not pure, so it is `.map`; its stem keeps the framework file it serves, because Next fixes the export name and `[locale]/` holds the companions of both `layout` and `not-found`: `page-metadata.map.ts`, `layout-metadata.map.ts`, `not-found-metadata.map.ts`. `generateStaticParams` is pure and becomes `layout-static-params.compute.ts`. A function that builds an HTTP response from an outcome is `.map`. A function that posts a form to the API is a command and is `.action`, as the console's `mutation.action` already is. A function whose work is a call into an external system or a platform API (the console API, console-stubs, a `fetch`, WebGL compilation, a DOM read) is `.adapter`, and `.compute` is reserved for pure functions; the entity-title readers, `socialImages`, `readConsole`, `consoleRequest`, `consoleStubsRequest`, `readResult`, `readPalette`, `compileShader` and `createProgram` move accordingly. A function that returns an imperative handle (`mountSurface`) is `.controller`. A function that builds a fresh instance (`createTerminalMachine`) is `.factory` even outside a `factories/` folder. An error class is a type with its methods and is `.type`. A module-level mutable object is `.state`. A GLSL source constant and a class-string constant are `.data`. A validator that is not a Zod schema is `.compute`. A one-interface props file is `.type`. A re-export-only module that is not `index.ts` is dissolved, its consumers importing the modules, unless it is a package export entry, which becomes the `index.ts` of a folder and the export map follows.
- **Suffix mismatches.** A `.compute` file that awaits an external read is renamed `.adapter`. A `.types` file that exports constant arrays gives them to a `.data` file. The website messages loader is `.map` since it selects catalogs by locale. `createEarlyAccessFormState` in a `.data` file is the state rules' concern (spec 0042) and is not renamed here. Files that specs 0034, 0035 or 0041 split or move receive their suffix at that split; the suffix tickets rename what remains and wait on those tickets. The console locale configuration module is replaced by the localization compliance spec, not yet written, and is not renamed here. When a suffix rename also fixes the stem (General, single responsibility 4), the rename does both at once and spec 0034's stem ticket excludes those files.
- **Translators.** Every single translator binding is `translate`, including the destructured `[{ id }, translate]` form. A scope with two translators names them `translate` plus the namespace in PascalCase (`translateCause`, `translateEvent`); message keys do not change.
- **Named object types.** Each object type in a signature, generic argument or assertion becomes an `interface` directly above its use in the same module, or, for the console read envelopes, in the data layer's types by concern (once spec 0034 has split them), imported with `import type`. An object type with one field counts: "non-trivial" is read strictly as "has a field". `Record<...>` and index signatures are not object literals. Unions of object shapes stay `type` aliases over named interfaces.
- **Await spacing.** Every `await` statement in a test gets its blank line before and after, two consecutive `await` statements are separated, and the one `await` shared with optional chaining lands in a binding first. The rule holds in tests as in source; the detectors measure both.

## Testing Decisions

- The seam is `drift quality check 04-typescript`: every ticket ends with its class of violation at zero and the whole spec ends with the rule at zero findings. `drift quality coverage` shows `04-typescript` applied at 100%.
- `drift check biome` is the seam for the matcher: its fixtures prove that `web` and `@proxy/*` sort in the alias group.
- Behaviour tests do not change in intent: the existing unit tests of each app and package pass unchanged apart from renamed imports and spacing; the website catalogs test keeps proving parity after the guard moves.
- A new detector is verified against the tree at its commit: it prints exactly the known violations of its class and nothing else, and zero once the class is fixed. Each new awk program also ships with fixtures that spec 0035's `drift quality detectors` runs.
- Prior art: spec 0030 for a rules-compliance pass over a package, `drift check biome` for configuration fixtures, `await-spacing.awk` for an awk detector over listed files and `atom-names.awk` for one over tracked paths.

## Out of Scope

- `try/catch` around pure code, empty catches, environment reader modules, the invented development secret, one-function-per-file splits, stem names, compound predicates and return spacing (specs 0034 and 0035 and the Nest compliance spec, not yet written).
- The `lib/` tree move, the proxies' implementation and the route handlers (spec 0041).
- `packages/tokens/build.js` (spec 0037).
- Test description names (the testing compliance spec, not yet written).
- Console locale configuration and locale literals (the localization compliance spec, not yet written).
- `*-state.data` defaults and form files by responsibility (spec 0042).
- Adding a suffix to the table: rules change only when the owner asks.

## Further Notes

Readings the owner can overturn at pull request review:

1. **Metadata companions are `.map` with the framework stem** (`page-metadata.map.ts`). The alternative, `generate-metadata.map.ts`, names the export but collides in `[locale]/`. A `.metadata` suffix would need a table change.
2. **I/O functions are `.adapter`, `.compute` is pure.** Readers of the console API, console-stubs and the DOM move to `.adapter`. The owner may prefer a `.query` or `.read` suffix in the table.
3. **Error classes are `.type`.** The table has no error suffix; `.service` is reserved for services.
4. **A factory outside `src/factories/` keeps `.factory`** for `createTerminalMachine` in `packages/web`.
5. **Two translators in one scope** are `translate<Namespace>`. The alternative, one un-namespaced translator with full keys, changes every key in those components.
6. **One-field object types are named.** The rule's "non-trivial" is read as "has a field".
7. **Pseudo-barrels are dissolved** rather than moved into folders, except package export entries.

Open question: the deep-relative-import rule and the `@proxy/*` alias touch the console proxy that spec 0041 rewrites; the alias ticket waits on that rewrite so the proxy is settled once.
