# 0284: every object type in a signature is a named interface

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Each object type in a signature, generic argument or assertion becomes an `interface` directly above its use, or, for the console read envelopes, in the data layer's types by concern, imported with `import type`. A one-field object counts. Unions of object shapes stay `type` aliases over named interfaces.

**Blocked by:** 0281, 0282, 0283.

**Status:** ready-for-agent

- [x] The inline-object-type detector reports nothing.
- [x] `pnpm turbo run check typecheck test` passes.

> **Note, 2026-11-05.** The inline-object-type detector is ticket 0277's. Until it exists, a TypeScript-AST scan stands in for it, over every TypeScript file in `apps/` and `packages/`, tests and configurations included. It looks for an object type literal used as a type argument, an assertion, a return type, a parameter type or a variable annotation, and it reports none after this ticket (97 before).
> - Route params are named per route: `LocaleRouteParams`, `ProductRouteParams`, `RuntimeRouteParams`, `TicketRouteParams` and `BreadcrumbRouteParams`.
> - The console's read envelopes live in `lib/data/read-envelopes.types.ts`, one concern, until spec 0034 splits the data layer's types.
> - The web package exports `LanguageSwitcherLinkOption`.
> - The catalog shapes, the cookie and crypto services, the locale prefixes, the test mocks and the three Vitest configurations name their shapes above their use.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test`
- `platform/cli/drift quality check 04-typescript`
