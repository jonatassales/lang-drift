# The website complies with the repository rules

**Status:** ready-for-agent

## Problem Statement

The owner wants `apps/website` to apply every rule in `.claude/rules/00-general`, `01-filesystem`, `03-tailwind`, `04-typescript`, `05-react`, `06-nextjs` and `07-jotai`. An audit found the website far from it: about 700 lines of hand-written CSS, 32 source files with more than one function, files whose names do not name their export or carry no suffix from the table, the translator named `t`, object types written as `type` or inline in signatures, environment variables read through a helper with a fallback, compound predicates without a name, nested ternaries, empty `catch` blocks, inline event handlers, props that are not a named interface, effects that derive state, client form logic that is not split by responsibility, and statement spacing that misses the rules, in the source and in the tests.

The rules say existing code does not legalize a prohibited pattern, so every one of these is recorded for remediation. Until it is remediated, every new change to the website starts from a tree the rules reject, and the detectors of `drift quality` cannot report the website as compliant.

## Solution

The website is brought to the rules without changing what a visitor sees or what the site does, except where the owner approved a deviation:

- **Breakpoints** snap to Tailwind's defaults: 760px becomes `md`, 1000px becomes `lg`, 600px becomes `sm` and 900px becomes `md`. Between each pair of values the layout changes.
- **Two colors without a token** are approximated: the navigation drawer's scrim becomes `bg-foreground/38` and the trial plan's border `border-primary/65`.
- **The global focus ring, `white-space: nowrap` on buttons and `font: inherit` on buttons and inputs** stop being element rules and become classes on each element that relied on them, including `LanguageSwitcher`, `AnimatedTabs` and `BasicAccordion` in `packages/web`.
- **SmoothUI classes stay ahead.** Where an element carries a SmoothUI `smooth-*` class, the unlayered SmoothUI stylesheet still wins over Tailwind's layered utilities, so the utilities that must beat it take the `!` modifier, marked as transitional until the SmoothUI styles are remediated.

Every stylesheet the website authored becomes Tailwind utilities in the JSX that uses it and is deleted; the Tailwind entry keeps only its imports. Shared state moves to `src/app/state`. Every source file declares one function, is named after it and carries the suffix of its responsibility. TypeScript, control flow, React components and client forms follow their rules, and the statement spacing follows the general and TypeScript rules in the source and in the tests.

## User Stories

1. As the owner, I want the website to apply every rule in the seven rule files, so that new work starts from a compliant tree.
2. As the owner, I want `drift quality coverage` to report 100% for the measured rules in the website, so that compliance is measured rather than asserted.
3. As the owner, I want any detector finding that remains to be a documented false positive, so that the detectors are not worked around.
4. As a visitor, I want every page to look as it did, so that the refactor costs me nothing.
5. As a visitor on a viewport between 600 and 640px, 760 and 768px, 900 and 768px or 1000 and 1024px, I accept that the layout switches at Tailwind's breakpoint instead of the old one, as the owner approved.
6. As a visitor opening the navigation drawer, I see a scrim drawn from the foreground role at 38%, as the owner approved.
7. As a visitor on the pricing page, I see the trial plan's border drawn from the primary role at 65%, as the owner approved.
8. As a keyboard user, I want every link, button and input to keep the focus ring it has today, so that focus stays visible.
9. As a visitor, I want buttons to keep their text on one line, so that controls do not wrap.
10. As a visitor, I want buttons and inputs to keep their inherited font, so that controls read like the surrounding text.
11. As a visitor, I want smooth scrolling to anchors, with the sticky header's offset, so that a section is not hidden under the header.
12. As a visitor who prefers reduced motion, I want smooth scrolling and the theme transition to stay off, so that the page does not animate.
13. As a visitor switching themes, I want the page background and text to keep their 280ms transition, so that the switch reads as it did.
14. As a visitor, I want the header to keep its sticky glass bar, its border once the page scrolls, its icons and its drawer, so that navigation is unchanged.
15. As a visitor on a narrow viewport, I want the header to hide the demo button and tighten its controls below `md`, so that the bar fits.
16. As a visitor, I want the footer to keep its brand, links, credit line and social links, so that the footer is unchanged.
17. As a visitor on the pricing page, I want the plans section, the trial card, its badge, value, features and call to action, and the availability line to look as they did.
18. As the owner, I want the dead branches of the stylesheets removed: the secondary plan and the footer's `data-compact` state, so that no style exists for markup that is never rendered.
19. As the owner, I want the early-access form's styles moved to utilities even though the form is paused, so that restoring it brings no stylesheet back.
20. As the owner, I want font utilities that today have no effect on buttons and inputs removed, so that the look does not change when the global `font: inherit` rule leaves.
21. As a developer, I want `DemoButton` to accept a `className`, so that the drawer can give it its 44px minimum height.
22. As a developer, I want shared state imported through `@state` from `src/app/state`, so that the tree matches the Next.js rules.
23. As a developer, I want every file to declare one function, named after it, with the suffix of its responsibility, so that a file's name tells me what it holds.
24. As a developer, I want the messages barrel to hold only re-exports, so that locale logic lives in files of its own.
25. As a developer, I want imports that reach into a subfolder to use an alias, so that `./` only ever reaches a sibling.
26. As a developer, I want the translator named `translate`, so that the code reads as the rule says.
27. As a developer, I want object shapes declared as named interfaces, so that signatures stay readable.
28. As a developer, I want generic type parameters named descriptively, so that `T` never appears.
29. As a developer, I want every environment variable read with `process.env.NAME` where it is used, with an explicit check for a missing value, so that no helper or fallback stands in for it.
30. As a developer, I want nested ternaries and deep conditionals replaced with guards, a `switch` or a lookup, so that control flow stays at most two levels deep.
31. As a developer, I want compound predicates named, so that a condition says which rule it applies.
32. As a developer, I want every `catch` to state its fallback, so that no error is swallowed in silence.
33. As a developer, I want every component's props declared as an interface right above it and destructured on its first line, so that component contracts read alike.
34. As a developer, I want component scopes in the order the React rules set, so that every component reads in the same order.
35. As a developer, I want event handlers named and declared in the component, so that JSX holds no inline function.
36. As a developer, I want `'use client'` only where a component needs the browser, so that the rest renders on the server.
37. As a developer, I want state that can be derived during render or handled in an event handler kept out of effects, and every remaining effect to say why it exists, so that effects stay a last resort.
38. As a developer, I want the early-access form split into one handler per event, one map per state mutation and its Zod schema in `form.validation.ts`, so that the form follows the state rules.
39. As a developer, I want the book-a-demo form to have its `form.types`, `form-state.data` and `form.validation` files, so that both forms read alike.
40. As a developer, I want the statement spacing of the general and TypeScript rules in the source and in the tests, so that every file reads alike.
41. As a developer, I want the tests to use readable names instead of `t` or `(a, b)`, so that they follow the naming rules too.
42. As the owner, I want test helpers to stay inside their test, so that one function per file does not scatter the tests.
43. As the owner, I want `lib/` kept, except for `state`, so that the move stays small.
44. As the owner, I want the `error.tsx` `retry` and `reset` mismatch and the commented-out early-access code left as they are and noted in the pull request, so that I decide them.

## Implementation Decisions

- **One pull request, one commit per ticket.** The work lands on one branch with a commit per ticket, each with its `Ticket:` trailer.
- **CSS to Tailwind in five tickets:** global and base; header; footer; plans; early-access. Each stylesheet becomes utilities in the JSX of the component that uses it and is then deleted. The Tailwind entry keeps only `@import 'tailwindcss'`, the tokens stylesheet and the approved utilities file.
- **Element rules move to elements.** `<html>` takes `scroll-smooth`, the 96px scroll padding and `motion-reduce:` for the scroll behavior; `<body>` takes the background, the foreground, the minimum height and the 280ms background and color transition, with `motion-reduce:` turning it off. Every link, button and input that relied on the global focus ring takes the same ring as classes; every button that relied on the global `nowrap` takes `whitespace-nowrap`; Tailwind's preflight already gives buttons and inputs `font: inherit` in the base layer, so the font utilities the global rule was overriding are removed instead.
- **SmoothUI precedence.** SmoothUI's stylesheet is unlayered and beats Tailwind's layered utilities. Where a website stylesheet overrode a SmoothUI declaration on an element that carries a `smooth-*` class, the replacing utility takes `!`, with a comment that marks it transitional. SmoothUI's CSS is not touched.
- **Breakpoints.** Every `max-width` media query becomes the matching Tailwind breakpoint, written mobile-first or with `max-*:`: 760px → `md`, 1000px → `lg`, 600px → `sm`, 900px → `md`.
- **Colors without a token.** The drawer's `rgb(0 0 0 / 38%)` scrim becomes `bg-foreground/38`; the trial plan's mix of 65% primary into the border becomes `border-primary/65`. Every other color maps to an existing role.
- **The early-access form's error colors** map to the status danger role; the form is not rendered while sign-up is paused, so no visitor sees the change.
- **Dead branches.** The secondary plan's styles and the footer's `compact` prop and `data-compact` styles are removed, because nothing renders them.
- **`@state`.** The barrel moves from `src/app/lib/state` to `src/app/state`, and the alias is repointed in `tsconfig.paths.json`. Every importer keeps importing `@state`.
- **One function per file.** Each multi-function file becomes single-function files grouped in a folder named after their concern. A route's `generateMetadata` and `generateStaticParams` live in their own files beside the route and are re-exported by it, since a framework file may re-export what Next reads. Test helpers stay in their test.
- **Names and suffixes.** A file is named after its one export, in kebab-case, with the suffix of the table in the TypeScript rules where one names its responsibility. The messages barrel holds only named re-exports. Relative imports reach only siblings; anything deeper uses an alias.
- **TypeScript.** The translator is `translate`; object shapes are interfaces; non-trivial signature types are named interfaces declared above their function; type parameters are words. `process.env.NAME` is read at the operation, with an explicit check that treats a missing value as "not configured", as today, and no `?? ''`.
- **Control flow.** Nested ternaries and conditionals deeper than two levels become guards, a `switch` or a lookup; compound predicates get a name; every empty `catch` states its fallback and keeps today's effect.
- **React.** Props interfaces sit right above their component; scopes follow the rule's order; handlers are named functions in the component; `'use client'` stays only where a component needs it; effects that derive state move to render or to the handler, and each remaining effect carries a comment with its reason.
- **Forms.** The early-access form gets one `<event>.handler.ts` per event, one `<mutation>.map.ts` per state mutation and its schema in `form.validation.ts`; the book-a-demo form gets `form.types.ts`, `form-state.data.ts` and `form.validation.ts`.
- **Spacing.** `return`, `await`, props destructuring and block boundaries follow the general and TypeScript rules in the source and the tests; tests drop one-letter names.
- **`packages/web`.** Only `LanguageSwitcher`, `AnimatedTabs` and `BasicAccordion` change, and only to carry the focus ring, `nowrap` and font classes the website's global rules used to give them, without changing the console.

## Testing Decisions

- **Visual parity is the main seam.** Before any change, each website route is captured in the built-in browser at 1280px and 375px, in light and dark, in `en`, with the navigation drawer open: every element's computed styles and box. After each CSS or UI ticket, the same capture is compared with the baseline; only the approved deviations may differ. Nothing of it enters the repository, and no end-to-end test is written.
- **The existing test suite is the behavior seam.** `pnpm turbo run check typecheck test build --filter=website` passes after every ticket, and `--filter=@langdrift/web` too when `packages/web` changes. Good tests here assert external behavior: the rendered markup, the form results, the route handlers' responses, the proxy's redirects and headers. A test that only mirrors a file's internals is moved with the code it covers, not rewritten.
- **The detectors are the rules seam.** `platform/cli/drift quality check` reports no finding in `apps/website` for each ticket's rules, and `drift quality coverage` reaches 100% for the measured rules in the website at the end. A finding that remains is reviewed and, if it is a false positive, documented in the pull request.
- **The two-axis `code-review`** runs before every commit, against this spec and every rule file.
- Prior art: the website's unit tests beside their modules (`*.unit.ts`, `*.unit.tsx`), the catalog test for localization and the SEO contract test.

## Out of Scope

- `apps/console` and its use of `packages/web`, which a parallel session changes, and the root lockfile.
- Every `packages/web` component other than `LanguageSwitcher`, `AnimatedTabs` and `BasicAccordion`, and SmoothUI's stylesheets, whose remediation the Tailwind rules already record.
- New design tokens: the two colors without one are approximated instead.
- The rules in `08-localization`, `09-seo` and `10-testing`, except the spacing and the names in the tests.
- The `error.tsx` page, which receives `retry` where Next passes `reset`, and the commented-out early-access code. Both stay as they are and are noted in the pull request for the owner to decide.
- End-to-end tests.
- The website's own `cn()` formatter. The Tailwind rules forbid an app copy, but `dead-code.allow` records that the owner keeps the formatters, and the plan does not list it, so it stays and is noted in the pull request for the owner to decide.

## Further Notes

The owner settled the open questions before this spec: `lib/` stays except for `state`; tests are in scope except for one function per file; the visual proof is the built-in browser's capture, before and after; the breakpoints, the two approximated colors and the per-element focus, `nowrap` and font classes are approved deviations.
