# Every React component complies with the React rules

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

> **Note, 2026-11-05.** The owner closed 0286, 0287, 0288, 0289, 0291 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

## Tickets

- [ ] 0286: the web package's components apply the React rules
- [ ] 0287: the icon family has one owner
- [ ] 0288: the console's components apply the React rules
- [ ] 0289: the recommendation request runs in one transition
- [ ] 0290: the website's components apply the React rules
- [ ] 0291: drift measures the React rules

## Problem Statement

`.claude/rules/05-react.rule.md` governs the components of `apps/console`, `apps/website` and `packages/web`. Specs 0026, 0027 and 0030 brought each surface closer to it; the audit of 2026-11-04 measured the rule at 90.0% and found what they left:

- Nineteen `template.tsx` files under `apps/console/src/app/[locale]/` that return their children unchanged. The console import carried a `PageTemplate` wrapper whose comment said the template "owns fresh form, paging and interaction state on entry"; spec 0026 replaced the wrapper with the pass-through files. Rule 9 says `template.tsx` exists only when a route must remount on navigation, and no file names such a navigation.
- `apps/console/src/app/lib/components/route-state/heading-skeleton.tsx` opens with `'use client'` and uses no hook, handler or browser API; `views/products/product-detail.tsx` and `views/runtimes/runtime-connection.tsx` show the same signature to a heuristic.
- Eight arrow functions in component scope in five files: `apps/console/src/app/lib/views/finops/cost-overview.tsx` (`usd`, `plain`, `reported`), `views/runtimes/runtime-graph.tsx` (`updatePalette`, `pause`), `views/runtimes/runtime-node-cards.tsx` (`place`, `update`), `apps/website/src/app/lib/components/website-header.tsx` (`closeOnDesktop`) and `website-header/language-pill.tsx` (`syncUrlSuffix`).
- Forty-nine props contracts off the rule: 13 with no blank line between the interface and the component (`apps/console/src/app/lib/components/page-heading/page-heading.tsx`), 34 with a declaration between them (`status-label/status-label.tsx` puts `statusTones` between; the card and dialog parts, `skeleton.tsx`, `chroma-text.tsx`, `elevated-frame.tsx` and `rotating-text.tsx` in `packages/web`), and route pages whose props interface is not named after the component.
- About 26 components whose scope leaves the rule 5 order: a custom hook after a third-party hook (`views/executions/execution-entry.tsx`), a third-party hook after a derived constant (`apps/website/src/app/lib/components/finops/delivery-receipt.tsx`), a handler after derived constants (`views/runtimes/runtime-playback.tsx`), a React hook after a third-party hook (`views/runtimes/runtime-graph.tsx`, `views/products/secret-dialog.tsx`, `early-access/form.tsx`).
- Twenty-four of 51 `useEffect` calls carry no comment saying why an effect is necessary (`navigation/mobile-navigation.tsx`, `runtime-graph.tsx`, `runtime-node-cards.tsx`, `packages/web/src/ui/scroll-hint/scroll-hint.tsx`, six in `packages/web/src/vendors/smoothui/mac-terminal.tsx`, two each in `hero-shader.tsx` and `vendors/vidro/code-editor.tsx`), and one derives state: `mobile-navigation.tsx` resets `openPath` from the pathname in an effect although the state already records the path it opened on.
- A loading flag, `reading`, in `views/runtimes/recommendation-panel.tsx` beside a `useTransition` that already models pending work.
- `apps/console/src/app/[locale]/(menu)/products/page.tsx` wraps its whole content in `<Suspense>` although `loading.tsx` beside it renders the same skeleton.
- Local icons: `apps/console/src/app/lib/ui/stack-icon.tsx` draws eight brand paths from `stack-icons.json`, four of which `packages/web/src/ui/icons/technology/` also draws; the website header draws its menu, sun, moon and close icons inline and the platform menu its caret, which `packages/web/src/ui/icons/arrow/caret-down-icon.tsx` already provides.

Each is a pattern the rule prohibits, and the existing detectors see only the arrow component, the destructured parameter, the inline handler, four loading-flag names, `createContext` and the blank line after destructuring.

## Solution

The console's nineteen pass-through templates are deleted; a `template.tsx` returns only for a route that must remount on a navigation the file names. `'use client'` sits only where a component needs it. Every function in component scope is a declaration. Every props interface sits directly above its component with exactly one blank line between, named after the component except where Next reserves the name. Every component's scope follows the rule 5 order. An effect that derives state or reports a state change moves to render or to the place that changes the state; the subscriptions the terminal, the website header and the console's recommendation panel keep move into custom hooks; every effect that stays says why. The file read in the recommendation panel joins its transition. The products page relies on its `loading.tsx`. The icon family has one owner: the console's stack icons and the website's header icons become `packages/web` icons, with one artwork per technology. `drift quality check react` reports zero findings, with new detectors for each pattern a program can catch.

## User Stories

1. As the owner, I want `drift quality check react` to report zero findings, so that the React rule is verified by the seam the repository already has.
2. As a console developer, I want no `template.tsx` that returns its children, so that a template exists only where a route must remount and says which navigation needs it.
3. As a console user, I want every page to keep its state exactly as today across navigation, so that removing the templates changes nothing I can see.
4. As a maintainer, I want `'use client'` only on components that need it, so that a static skeleton renders on the server.
5. As a reader, I want every function inside a component to be a declaration, so that a handler, a formatter and a subscription callback read alike.
6. As a reader, I want the props interface directly above its component, so that the contract and the component are one unit on the page.
7. As a reader, I want each component's scope in the rule's order, so that hooks, handlers and derived values are found where expected.
8. As a reader, I want every remaining `useEffect` to carry its reason, so that an effect is never a habit.
9. As a console user, I want the mobile navigation to close on route change as it does today, derived during render instead of through an effect.
10. As a console user, I want the recommendation panel to show its busy state while it reads my ticket file and asks for a recommendation, modeled by one transition.
11. As a website visitor, I want the terminal to play when in view and the page is visible, pause when I ask and report its progress, with that subscription owned by a hook.
12. As a website visitor, I want the navigation drawer to open as a modal, return focus when it closes and close when the viewport grows, with that subscription owned by a hook.
13. As a console user, I want the products page to show its skeleton through the route's `loading.tsx` while the list loads.
14. As a maintainer, I want every icon to come from `packages/web`, so that the console's stacks and the website's header share one icon identity.
15. As a console user, I want a runtime's stack to show the shared technology icon, accepting that four of them take the shared artwork.
16. As a website visitor, I want the header's menu, theme and close controls and the platform menu's caret to look as they do today.
17. As the owner, I want detectors for an arrow function in component scope, an undocumented effect, a props interface away from its component, a scope out of order, a needless `'use client'`, a pass-through template, an inline `<svg>` in an app and more loading-flag names, so that the patterns found here are reported from now on.
18. As the owner, I want each strict reading recorded, so that I can overturn it at PR review.

## Implementation Decisions

- **Templates (rule 9).** The nineteen console templates are deleted. A page below a dynamic segment already remounts when the segment's value changes, and a page under a static segment unmounts when the visitor leaves it; the console has no in-route navigation whose state the product requires to reset, and its forms are leaf client components that mount with their page. The ticket verifies, route by route, that no state survives a navigation it did not survive before. Should a route need a reset on a navigation that keeps its segment, its `template.tsx` returns with a comment naming that navigation; an effect or a key never imitates it.
- **Suspense (rule 9).** A page whose `loading.tsx` renders the same fallback drops its page-level `<Suspense>`. A `<Suspense>` stays around a section that renders on its own inside a page that shows other content meanwhile.
- **Client boundary (rule 1).** `'use client'` is removed from every component that uses no hook, handler, ref, browser API or client-only library, starting with the heading skeleton; a route's `error.tsx` keeps it, as Next requires.
- **Functions (rules 3, 6).** Every arrow function assigned in component scope becomes a function declaration, including those declared inside an effect; a formatter the component builds from a hook result (currency, plain number, reported value) is a declaration in the component.
- **Props contract (rule 4).** The props interface sits directly above the component with exactly one blank line between. A comment that introduces the component sits above the interface, so the interface and the component stay adjacent. A declaration that sat between them (a tone map, a class string, a constant) moves above the interface. The interface is `<Component>Props`; a route page whose component is `Page` keeps a route-named interface, because Next generates a global `PageProps`, and a shared icon props interface keeps its shared owner.
- **Scope order (rule 5).** Each component follows: props destructuring; React hooks; custom hooks; third-party hooks; named handlers; derived constants and memoized values; the return. A value a hook needs that would be derived first is computed inside the hook's argument or by a `.compute` function called there. A guard that returns early on a derived value moves below every hook, as React requires.
- **Effects (rule 7).**
  - An effect that derives state moves to render: the mobile navigation's open state is `openPath === pathname`, with no effect.
  - An effect that reports a state change moves to the place that changes the state: the terminal reports its progress where it completes a line.
  - A subscription moves into a custom hook with its cleanup: the terminal's in-view, page-visibility and animation loop become a terminal playback hook; the website header's drawer modal, body scroll lock, focus return and desktop close become a navigation drawer hook; the recommendation panel's abortable file read and request become a recommendation request hook.
  - Every effect that stays carries a comment with its reason on the line above it: a DOM measurement, an observer, an event subscription, a lazy import of a client-only library, a cleanup on unmount.
- **Transitions (rule 8).** The recommendation panel's `reading` flag is removed; reading the ticket file and requesting the recommendation run inside the existing transition, so `pending` models both, and a newer submission aborts the earlier one through the request controller.
- **Icons (rule 11).** `packages/web` gains the technology icons the console draws that it lacks (React, Python, PostgreSQL, Kafka) and the header icons the website draws inline (menu, sun, moon, close), each one component in the icon family with the family's props. The console's `StackIcon` becomes a lookup from a stack id to the shared technology icon, falling back to the neutral cube; `stack-icons.json` is deleted. Where both drew a technology (Node.js, Next.js, Go, Redis), the shared artwork wins. The platform menu takes the shared caret.
- **Detectors (02-tooling rule 4).** `drift quality check` gains, for `05-react`: an arrow function assigned to a constant inside a TSX function body; a `useEffect` whose preceding line is not a comment (awk); a component whose props interface does not end two lines above it with one blank line between (awk); a scope whose categories appear out of order, classifying React hooks, custom hooks imported from `.hook` modules, third-party hooks, handlers, derived constants and the return (awk); a `'use client'` file with no hook call, handler prop, ref, browser global or client-only import (awk); a `template.tsx` whose body returns its children (awk); an inline `<svg>` in an app; and the loading-flag detector extended to `reading`, `fetching` and `busy`. The scope-order and client-boundary programs are heuristic; a false positive is documented in the pull request, never worked around.
- **Order.** `packages/web` lands first (its own components, the new icons, the terminal hook), then the console and website once spec 0041 has moved their trees, then the detectors, so that `drift quality check react` is zero when they land.

## Testing Decisions

- **Behavior is the seam.** `pnpm turbo run check typecheck test build` for `web`, `console` and `website` passes after every ticket. A good test renders the component and asserts its markup, attributes and what a user event changes, never the order of hooks or the presence of an effect. The hooks that replace effects get their own tests: the terminal playback hook with a fake observer, visibility and clock; the navigation drawer hook with a fake dialog element; the recommendation request hook with a fake file and fetch, asserting that a newer submission aborts the earlier and that `pending` covers the read. Prior art: `packages/web/src/vendors/smoothui/mac-terminal.unit.tsx`, `scroll-hint.unit.tsx`, the console's view tests under `apps/console/src/app/lib/views/`, the website's `early-access/form.unit.tsx`.
- **Navigation parity for the templates.** Before deleting the templates, each console route is visited in the built-in browser and its state after a navigation away and back is recorded under `.audits/` (form fields, paging position, open dialogs); after the deletion, the same walk is compared. No browser runner is added.
- **Visual parity for the icons.** The console's runtime views and the website's header are captured before and after in the built-in browser, at 1280px and 375px, in light and dark; the four technologies that take the shared artwork are the only accepted difference.
- **The detectors are the rules seam.** `platform/cli/drift quality check react` reports zero findings at the end; each new awk program is exercised against a fixture that holds one compliant and one violating component.
- **Review-enforced.** A `<Suspense>` wider than its section, an effect whose comment does not state a real reason, a `'use client'` the heuristic cannot judge, and a template that returns with a reason are checked by the review.

## Out of Scope

- Rules other than `05-react`: the SmoothUI stylesheets (spec 0038), client form structure and atoms (spec 0042), file suffixes and names (spec 0039), the tree move (spec 0041), route handlers (spec 0041).
- The behavior of any component; only its structure changes, except the four technology icons' artwork, which the owner can overturn.
- Pages already wider than the viewport on mobile.
- Framework files that require `'use client'` (`error.tsx`, `global-error.tsx`).

## Further Notes

- **Strict readings taken.** (1) A comment introducing a component moves above its props interface, so that "exactly one blank line" separates interface and component; the owner may prefer the comment between them. (2) Arrow functions assigned inside an effect body count as component scope and become declarations. (3) A pass-through template is treated as a key hack the rule forbids, overriding spec 0026's "return their children directly", since rules outrank specs. (4) A subscription with cleanup moves into a custom hook even when a single component uses it, because rule 7 names the hook as the place for a subscription.
- **Framework constraint recorded.** Next.js generates global `PageProps` and `LayoutProps` helpers, so a route component named `Page` keeps a route-named props interface; the detector allows an interface ending in `Props` whose name the signature uses.
- **Visual change to confirm.** Four console stack icons (Node.js, Next.js, Go, Redis) take the shared `packages/web` artwork; the alternative, keeping the console's paths, would leave two drawings of one technology in the icon family.
- **Spec 0026's decision** that `StackIcon` and the technology icons stay with the console's domain is overturned by rule 11, as the general rules' precedence requires; the icons move, the stack-to-icon mapping stays with the console.
