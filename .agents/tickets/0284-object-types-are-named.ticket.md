# 0284: every object type in a signature is a named interface

**Spec:** .agents/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Each object type in a signature, generic argument or assertion becomes an `interface` directly above its use, or, for the console read envelopes, in the data layer's types by concern, imported with `import type`. A one-field object counts. Unions of object shapes stay `type` aliases over named interfaces.

**Blocked by:** 0281, 0282, 0283.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The inline-object-type detector reports nothing.
- [x] `pnpm turbo run check typecheck test` passes.

> **Note, 2026-11-05.** The inline-object-type detector is ticket 0277's. Until it exists, a TypeScript-AST scan stands in for it, over every TypeScript file in `apps/` and `packages/`, tests and configurations included. It looks for an object type literal used as a type argument, an assertion (`as` or `satisfies`, through `readonly` arrays), a return type, a parameter type or a variable annotation, and it reports none after this ticket (97 before).
> - Route params are named per route: `LocaleRouteParams`, `ProductRouteParams`, `RuntimeRouteParams`, `TicketRouteParams` and `BreadcrumbRouteParams`.
> - The console's envelopes live in the data layer by concern: `product-envelopes.types.ts`, `runtime-envelopes.types.ts`, `history-envelopes.types.ts` and `console-envelopes.types.ts`.
> - The unions of object shapes (`ReadResult`, `AnalysisSummary`, `WorkflowSummary`) are type aliases over named interfaces.
> - The web package exports `LanguageSwitcherLinkOption`, and `LanguageSwitcherLinkProps` is an interface.
> - The catalog shapes, the cookie and crypto services, the locale prefixes, the test mocks and the three Vitest configurations name their shapes above their use.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test`
- `platform/cli/drift quality check 04-typescript`
