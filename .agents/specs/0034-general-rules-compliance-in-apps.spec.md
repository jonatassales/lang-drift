# The apps apply the general rules

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

> **Note, 2026-11-05.** The owner closed 0236, 0237, 0238, 0241, 0242, 0244 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

## Tickets

- [ ] 0236: drift sees the general-rule violations in the apps
- [ ] 0237: console-stubs' Go applies the general rules
- [ ] 0238: demo-requests and mailer apply the general rules
- [ ] 0239: the website's files hold one responsibility each
- [ ] 0240: the website reads its environment where it uses it
- [ ] 0241: the console's types and stub frame split by concern
- [ ] 0242: the recommendation panel, theme storage and policy environment apply the general rules
- [ ] 0243: the website's statements are spaced by phase
- [ ] 0244: the console's statements are spaced by phase

## Problem Statement

The TypeScript under `apps/` (the website, the console and the two Nest microservices) breaks `.claude/rules/00-general.rule.md` in ways `drift quality check general` does not see. An audit on 2026-11-04 measured the 627 TypeScript files under `apps/` at 94% applied and `drift quality check general` printed no TypeScript finding under `apps/`: every remaining violation sits in a class its detectors do not cover. The same audit counted:

- **Single responsibility.** `apps/console/src/app/lib/data/types.ts` holds 48 types across every console domain; `apps/console/src/app/lib/views/runtimes/recommendation-panel.tsx` uploads a ticket, posts for a recommendation and renders the panel; `apps/console/src/app/lib/views/runtimes/stub-frame-validator.ts` mixes generic type guards with the parsing of a stub widget frame; `apps/website/src/app/lib/components/early-access/request.body.ts` declares a class and a function. Eleven test files declare two to five helper functions beside their suites (`apps/website/src/app/metadata/seo-contract.unit.ts` five, `apps/microservices/demo-requests/src/demo-requests/demo-requests.controller.unit.ts` four). Nine files are not named after their export (`apps/console/src/app/lib/data/server.ts` declares `readConsole`, `apps/console/src/app/template/i18n/actions.ts` declares `setLocale`, `apps/website/src/app/lib/components/early-access/turnstile.validation.ts` declares `verifyEarlyAccessChallenge`).
- **Control flow.** 31 `if` conditions combine three or more operands without a name (`stub-frame-validator.ts` seven of them, `apps/website/src/services/cookie.service.ts` three, `apps/website/src/proxy/locale.ts` two). Five places exceed two conditional levels through chained ternaries or an `else if` chain (`apps/console/src/app/lib/views/runtimes/recommendation-accounting.tsx`, `runtime-list.tsx`, `apps/website/src/app/lib/components/early-access/form-submit.handler.ts`, `parse-submission-json.compute.ts`).
- **Statement spacing.** 121 declarations are directly followed by the guard that tests them, with no blank line between the computation and the check (`apps/console/src/app/lib/components/navigation/product-breadcrumb.tsx`, `apps/website/src/app/atoms/microruntime-session.atom.ts` six, `apps/microservices/demo-requests/src/archive/archive-bucket.factory.ts`): 32 in the website, 28 in the console, 16 in the microservices, and the rest in Go. Six test files have a `return` without a blank line before it.
- **Readable names.** Two callbacks name a parameter `_` (`runtime-node-cards.tsx`, `stub-frame.test.ts`). 28 test assertions consume an awaited result before naming it (`expect(await …)` in `apps/website/next.config.unit.ts`, `api/book-a-demo/route.unit.ts`, `early-access/turnstile.validation.unit.ts`).
- **External input.** The website reads its registry and Turnstile variables through two modules whose job is reading and checking the environment and returning a configuration object (`read-registry.compute.ts`, `read-turnstile.compute.ts`), fed by a third that reads `NODE_ENV` and `EARLY_ACCESS_MODE` (`early-access-mode.compute.ts`); the console proxy has the same shape in `policy-environment.compute.ts`. `apps/website/src/factories/cookie-service.factory.ts` invents a `CRYPTO_SECRET` in development and test, introduced by ticket 0181; the rule outranks that ticket.
- **Errors.** Four `catch` blocks hold only a comment (`theme-control.tsx`, `language-pill.tsx`), five promise rejections are discarded with `.catch(() => undefined)` (`request.body.ts`, `turnstile.validation.ts`, `early-access/form.tsx`, and two in the early-access route), and `early-access/form.tsx` wraps a pure environment check in the `try` that guards the Turnstile controller.
- **console-stubs.** `apps/microservices/console-stubs` is Go, not TypeScript, and carries the only 34 findings `drift quality check general` prints under `apps/` today: free functions beside a type (`domain/ticket.go` eleven, `presentation/httpapi/handler.go` five), one-line struct literals, one-letter names, `_` discards and returns without a blank line, in 13 files.

Spec 0026 let a test's helpers stay in its file and spec 0027 kept the early-access environment readers; the rules outrank both, so this spec fixes what they kept. Until it lands, every change to an app starts from a tree the rules reject, and `drift quality check general` cannot report the apps as applying the rule because it cannot see most of what is wrong.

## Solution

Every TypeScript file under `apps/` applies `00-general`, and `drift quality check general` can see it: the detectors gain one row per violation class a pattern or an awk program catches reliably, and the apps are brought to zero findings, batch by batch, with no change in what a visitor, a console user or a demo requester gets.

- The environment is read where it is used. The early-access registry forward, the Turnstile verification and the console's document policy read their variables at the operation, with the checks that protect production beside the read; the three website reader modules and the console's disappear. The cookie service requires `CRYPTO_SECRET` in every environment and invents nothing.
- No error is swallowed. A browser-storage failure is translated into the result of the operation that wrote or read it; a discarded rejection gets a named handler that reports it; a `try` wraps only the external call.
- One file, one responsibility. The console's types are split by domain, the recommendation panel into its upload, its request and its view, the stub-frame guards apart from the widget parser; `BodyError` leaves the bounded body reader; a test's helpers become fixtures beside the test; a file is named after its export.
- Control flow reads top down. A compound condition gets a name, a chained ternary becomes a lookup or a guard, a declaration and its guard are separated by a blank line.
- console-stubs's Go applies the same rules, with its free functions in files of their own.

## User Stories

1. As the owner, I want `drift quality check general` to print nothing under `apps/`, so that the apps are measured as applying the general rules instead of asserted to.
2. As the owner, I want a detector for every violation class the audit found that a pattern or an awk program can catch reliably, so that the next regression is caught in pre-commit and CI rather than in review.
3. As the owner, I want the violations no detector can see named in this spec as review-enforced, so that the review knows what to look for.
4. As the owner, I want `drift quality check general` to keep printing a list to review and never fail, as it does today, so that the detectors stay advisory until the owner says otherwise.
5. As a reviewer, I want each ticket of this spec to leave `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test` and `drift check localization` green, so that compliance costs no behaviour.
6. As a visitor, I want every page of the website to look and behave as it does, so that the refactor costs me nothing.
7. As a console user, I want every screen to look and behave as it does, so that the refactor costs me nothing.
8. As a demo requester, I want a demo request to be accepted, forwarded, emailed and recorded as it is today, so that the microservices' refactor costs me nothing.
9. As an operator, I want the website to fail at start when `CRYPTO_SECRET` is missing in any environment, so that no instance runs on a key nobody set.
10. As an operator, I want the Turnstile verification to keep refusing a non-official endpoint and test credentials in production, so that moving the reads to the operation loses no guarantee.
11. As an operator, I want the early-access registry forward to keep refusing a registry or public origin that carries credentials, so that the inlined reads keep the credential-free check.
12. As a developer, I want each environment variable read with `process.env.NAME` at the operation that needs it, so that no module exists whose job is reading the environment.
13. As a developer, I want each required variable documented by key in its app's README, so that the contract stays visible as the reads move.
14. As a developer, I want a browser-storage read or write to return whether it succeeded, so that the caller states its fallback instead of an empty `catch`.
15. As a developer, I want every discarded promise rejection to reach a named handler that reports it, so that no failure disappears.
16. As a developer, I want a `try` to wrap only the external call, so that an internal check reads as a guard clause.
17. As a developer, I want the console's domain types split by concern, so that a change to tickets does not touch the runtimes' types.
18. As a developer, I want the recommendation panel split into its upload handler, its request and its view, so that each has one reason to change.
19. As a developer, I want the stub frame's generic guards apart from its widget parser, so that the guards can serve another frame.
20. As a developer, I want `BodyError` in its own file, so that the bounded body reader declares one function.
21. As a developer, I want a test's helpers to live as fixtures beside the test, one per file, so that test files obey the same rule as source files.
22. As a developer, I want every file named after its export, so that I find code by its file name.
23. As a developer, I want every compound condition named, so that a guard says which rule it applies.
24. As a developer, I want no chained ternary and no third conditional level, so that an outcome is read from a lookup or a guard.
25. As a developer, I want a blank line between a declaration and the guard that tests it, so that the computation and the check read as two phases.
26. As a developer, I want no `_` parameter and no awaited result consumed before it is named, so that every value has a name.
27. As a developer, I want console-stubs's free functions in files of their own, its literals one field per line and its names descriptive, so that the stub API reads like the rest of the Go.
28. As a developer, I want the tickets that touch the website or the console to wait for spec 0041's tree move, so that no file is moved twice.
29. As a developer, I want the microservices' tickets to leave the Nest-specific items (environment factories, the mailer controller, the reader's `try`) to the Nest compliance spec, not yet written, so that each finding has one owner.

## Implementation Decisions

- **Detectors (the verification seam).** `drift quality check general` gains, in the detector table, one row per class:
  - a guard directly after the declaration it tests, through a new awk program for TypeScript: a `const`, `let` or `var` line whose binding is tested by the `if` on the very next line, with no blank line between; Go is excluded because the rule's own Go example puts `if err != nil` directly under the assignment;
  - a chained ternary: a `?` branch line directly under a `?` or `:` branch line, as Biome formats a nested ternary;
  - a compound condition with three or more operands in an `if` or `while`, across lines;
  - a swallowed rejection (`.catch(() => undefined)`, `null`, `{}`, `false`) and a comment-only `catch`, by widening the empty-catch row;
  - a `_` parameter, by widening the one-letter-name row;
  - an awaited result consumed before it is named, in tests (`expect(await …)`);
  - test files in scope for the TypeScript multi-function and return-spacing rows, through rows scoped to `*.unit.*`, `*.test.*` and `testdata/`.
  The detectors stay POSIX shell and awk in the tooling's existing table and programs, and `drift quality check` keeps listing without failing. The rows look at `apps` and `packages`; what they find under `packages` belongs to spec 0035.
- **Review-enforced classes.** Single responsibility of a file's concern (a types file across domains, a component that uploads, posts and renders), a `try` around a pure internal check, and the "one reason to change" judgement have no detector; the review checks them.
- **Website environment.** The early-access registry forward becomes an adapter under the website's infra, in the shape of the demo-request forward adapter: it reads `EARLY_ACCESS_REGISTRY_URL`, `EARLY_ACCESS_REGISTRY_KEY` and `EARLY_ACCESS_PUBLIC_ORIGIN` at the fetch, with the credential-free origin check and the key-length check beside the reads. The Turnstile verification reads `TURNSTILE_SECRET_KEY`, `TURNSTILE_VERIFY_URL` and `TURNSTILE_EXPECTED_HOSTNAME` at the fetch, with the hostname shape, the official-endpoint-in-production and the no-test-credentials-in-production checks beside them. The mode (development, test, production) is derived at each of those two operations from `NODE_ENV` and `EARLY_ACCESS_MODE`, as a named type and two reads; the three reader modules are deleted. The README keeps every key.
- **Cookie secret.** The cookie service factory reads `CRYPTO_SECRET` and throws when it is missing, in every environment; the in-memory generated secret and its `NODE_ENV` branch go. The unit tests set the variable. The README and the context note drop "generated on the fly"; where the development value comes from is an open question below.
- **Console environment.** The proxy's policy builder reads `NODE_ENV` and `VERCEL_ENV` where it chooses the document policy's environment; the policy-environment module is deleted. Spec 0041 owns the proxy's structure; this spec only moves the two reads.
- **Browser storage.** Each app that persists a preference in `localStorage` (the console's theme, the website's language) does it through one adapter per preference that returns whether the write took and `undefined` for a blocked read; the component names the fallback (keep the default, keep navigating). The adapter's `catch` translates the failure into that result.
- **Rejections.** A body or reader cancel and the Turnstile prepare get a named rejection handler that reports the failure through the app's existing observability path and keeps today's outcome. The early-access route's two discarded rejections are fixed by spec 0041's rewrite of that route under this contract.
- **Pure checks outside `try`.** The early-access form reads its site key and refuses a missing one as a guard; the `try` wraps only the Turnstile controller's creation.
- **Console types.** The console's data types are split into one `.types.ts` per concern (account, product, runtime, execution, recommendation, ticket, page), imported with `import type`; nothing is renamed.
- **Recommendation panel.** The ticket upload becomes a file-change handler file, the recommendation request a function of its own, and the panel keeps the view; the reading version and the abort controller stay with the handler that owns them.
- **Stub frame.** The generic guards (`isRecord`, `isText`, `isDate`, `isNumber`, `isNullable`, `hasFields`) become one validation file each; the component and widget guards keep their frame concern; the parser composes them. Each of its seven compound conditions gets a name.
- **Bounded body.** `BodyError` moves to its own file; the reader keeps its signature.
- **Test helpers.** A helper a test declares moves to a `testdata/` folder beside the test, one function per file: a fake or fixture builder as a `.mock.ts`, a pure computation as a `.compute.ts`, as the microservices already do. The test keeps only its suites. Spec 0026's exception for test helpers is withdrawn.
- **File names.** The seven files whose stem does not name their export are renamed after it, taking the suffix of their responsibility in the same rename, so spec 0039's suffix work does not touch them again; the two whose file name is the one the state rules prescribe (`form-state.data`) rename their export instead (`bookDemoFormState`).
- **Compound conditions.** Each gets a named boolean directly above the guard, named after the rule it applies (`hostAllowed`, `isCrossSite` already exist as the pattern). Chained ternaries become a lookup from code to outcome or a sequence of guards.
- **Phase spacing.** A blank line separates a declaration from the guard that tests it, batch by batch, after the structural tickets of each app so that no file is reformatted twice.
- **console-stubs.** Each free function of `apps/microservices/console-stubs` moves to a kebab-case file named after it in its package; methods of one type may span files grouped by concern as they do today. Literals take one field per line, one-letter names become words, a discarded value is named or its error handled, and every `return` that follows a statement gets its blank line. Spec 0036's Go rename waits for these new files.
- **Order and batches.** The detectors land first. Each app then takes two tickets, structure then spacing, sized to one context window; the website and console tickets wait for spec 0041's dissolution of `lib/` so a file is moved once.

## Testing Decisions

- The seam is `drift quality check general` (and `drift quality coverage general`): zero findings under `apps/` once every ticket lands, and each ticket's verification runs it on its paths. A detector's acceptance is the list it prints against the current tree: the known sites and no file in a named negative set.
- Behaviour must not change: `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test` and `drift check localization` stay green on every ticket; the website's SEO contract test and its catalog test, the console's node tests and the microservices' vitest suites are the regression net.
- A good test checks the outcome of an operation, not its structure: a cookie factory that throws without `CRYPTO_SECRET` in every `NODE_ENV`; a Turnstile verification that answers `unavailable` for a non-official endpoint in production and `verified` for the official one; a registry forward that refuses a registry URL with credentials; a storage adapter that answers `false` when `localStorage` throws; a widget guard that refuses each malformed frame the stub-frame test already pins.
- Prior art: the cookie factory's unit tests, the early-access provider-configuration and Turnstile validation tests, the console's stub-frame node test, the demo-request forward adapter's tests, and the microservices' `testdata/*.mock.ts` fakes.
- Moved test helpers keep their tests' assertions unchanged; a fixture gets no test of its own unless it computes something.

## Out of Scope

- `route.ts` files (spec 0041), including the early-access route's `try` and its two discarded rejections, which 0041 fixes under this spec's contract.
- Translator names (spec 0039), file suffixes beyond the renames above (0039), inline object types and await spacing (0039).
- The dissolution of `apps/{website,console}/src/app/lib/` and the console proxy's structure (spec 0041).
- Client form splits, atom granularity and `*-state.data` defaults (spec 0042).
- The console's locale configuration module, including `isLocale` (the localization compliance spec, not yet written).
- Test describe names and bare `test()`/`it()` (the testing compliance spec, not yet written).
- Nest-specific items: environment reads outside factories, the mailer controller's `deliver`, the event-log reader's `try` around a pure parse, the vitest decorator setup (the Nest compliance spec, not yet written). A Nest `.factory.ts` reading `process.env` at creation is allowed by the owner's decision 1 and is not a finding here.
- `lib/`, `packages/` and `platform/cli` (spec 0035), even where this spec's new detectors report them.

## Further Notes

- **Strict readings taken** (the owner may overturn each at PR review):
  - A test file obeys "one function per file": its helpers move to fixtures. This withdraws spec 0026's exception.
  - A module that reads two variables and returns a mode (`earlyAccessMode`, `policyEnvironment`) is an environment module and is dissolved, at the cost of two short mode derivations in the website.
  - A compound condition of three or more operands is unnamed even when its operands are named booleans; the detector flags it.
  - Methods of one Go type may span files grouped by concern (the console-stubs `Handler`); merging them into one file would hurt readers more than it serves the rule.
  - `expect(await …)` consumes a result before naming it; tests name the awaited value first.
- **Open questions for the owner:**
  - `apps/microservices/console-stubs` is Go under `apps/`; the brief's numbering gives this spec the TypeScript under `apps/` and spec 0035 the Go under `lib/` and `packages/`. This spec drafts the console-stubs ticket so the 34 findings have an owner; it moves to 0035 if the owner prefers.
  - With the generated development `CRYPTO_SECRET` removed, where does a fresh checkout get its value: a tracked development value in `.env.development` (which ticket 0181 removed on purpose) or the developer's untracked `.env`, checked by `drift env validate`? This spec assumes the untracked `.env` and documents it.
  - The suffix table has no entry for an error class (`BodyError`) nor for a test interaction driver (`mount`, `fill`, `submit`); this spec uses `testdata/*.mock.ts` for fixtures and drivers and leaves the error class unsuffixed for spec 0039 to settle.
- **Coordination.** The early-access route's reads of the registry configuration are replaced when 0041's rewrite of that route lands, so this spec's website environment ticket waits for it. The microservices' spacing ticket may touch files the Nest compliance spec rewrites; whichever lands last keeps the spacing, and the detector catches a regression.
