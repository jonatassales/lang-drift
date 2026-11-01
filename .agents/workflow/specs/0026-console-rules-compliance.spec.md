# The console complies with the repository rules

**Status:** ready-for-agent

## Problem Statement

`apps/console` breaks the rules in `.claude/rules/` that apply to it: `00-general`, `01-filesystem`, `04-typescript`, `05-react`, `06-nextjs` and `07-jotai`. An audit on 2026-10-31 counted:
- 39 files with more than one function, and 32 files whose name or suffix does not match their export;
- the translator named `t` in 75 places;
- `type` aliases and inline object types where interfaces belong, and one-letter type parameters;
- unnamed compound predicates, chained ternaries, empty or discarding catches, and an environment fallback;
- props interfaces away from their component, scope order breaks, 39 unnamed handlers, 2 loading flags, derivable effects, and `'use client'` wider than needed;
- a login/signup page discriminator, and a pass-through `PageTemplate` wrapper behind 19 `template.tsx` files;
- app-local UI primitives and an app-local icon dependency, where rule 05.11 puts shared primitives and the icon family in `packages/web`.

## Solution

The console complies with those rules and keeps its behaviour. The owner approved one visual change: its button, card, dialog and skeleton become shadcn primitives owned by `packages/web`, themed by `@langdrift/tokens`. Everything else looks and behaves as before.

## User Stories

1. As the owner, I want `drift quality check` to report nothing in `apps/console` for rules 00, 01, 04, 05, 06 and 07, so the console is the example the rules describe.
2. As a reader of the console, I want one function per file, named after it, so I find code by its file name.
3. As a console user, I want every page to look and behave as before, except the approved primitive restyle.
4. As a maintainer, I want the console's generic primitives and its icons to come from `packages/web`, so every app shares one button, card, dialog, skeleton and icon family.

## Implementation Decisions

- **Primitives (05.11):**
  - `packages/web` gains shadcn Button, Card, Dialog and Skeleton as owned source under `src/ui/`, styled with token role utilities and exported by name.
  - `ThemeToggle` and `ListEnd` are generic and move to `packages/web`.
  - `StackIcon` and the technology icons stay with the console's domain.
  - The Phosphor icons the console uses come from a `packages/web` export, and the console drops `@phosphor-icons/react`.
- **Structure (06):**
  - `lib/` stays as is.
  - `metadata/entity-title.ts` leaves `metadata/`, which holds SEO sources only.
  - `api/extub/[...path]/route.ts` delegates its logic to a `.route.ts` file.
- **One function per file (00):**
  - Multi-function files are split into folders named after their concern.
  - Each file is named after its export, with the suffix from the 04 table.
- **TypeScript (04):**
  - The translator is `translate`.
  - Object shapes are interfaces, and signatures use named interfaces.
  - Type parameters have descriptive names.
- **Control flow and external input (00):**
  - Chained ternaries and `else if` chains become guards, a `switch` or a lookup.
  - Compound predicates get names.
  - A catch translates or reports the failure instead of discarding it, keeping today's outcome.
  - `EXTUB_URL` is read where it is used, without a hardcoded fallback.
- **React (05):**
  - The props interface sits directly above its component.
  - Scope follows the 05.5 order.
  - Handlers are named function declarations.
  - Loading flags become transitions.
  - `'use client'` sits at the narrowest boundary.
  - Derivable effects move into render or handlers, and the effects that stay say why.
  - Login and signup compose their own views.
  - The `template.tsx` files return their children directly.
  - The body-level focus ring becomes classes on each interactive element.
- **Tests:** the rules apply to test files too. The exception is that a test's helpers stay in its file.

## Testing Decisions

- Each ticket keeps `pnpm turbo run check typecheck test build --filter=console` green, plus `--filter=@langdrift/web` when `packages/web` changes.
- `drift quality check` shows no console finding for the ticket's rules.
- A layout signature per route compares every element's box and computed style before and after: 20 routes, at 1280 and 375px, in dark and light. Tickets C2 to C6 must match it. C1 and the focus ring differ only where the primitives and the focus classes apply.

## Out of Scope

- The website, which has its own spec, 0027.
- The SmoothUI stylesheet in `packages/web`.
- Rules 03, 08, 09 and 10, except naming and spacing in tests.
- The mobile pages that are already wider than the viewport.
