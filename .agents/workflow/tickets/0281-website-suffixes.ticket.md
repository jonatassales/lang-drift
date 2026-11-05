# 0281: the website files carry the suffix of their responsibility

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Every TypeScript file in the website that the earlier splits did not already rename takes the suffix of its responsibility, by the spec's calls: metadata companions `.map` with the framework stem, I/O functions `.adapter`, pure ones `.compute`, commands `.action`, handles `.controller`, fresh instances `.factory`, error classes `.type`, module-level mutable objects `.state`, constants `.data`; a re-export-only module that is not `index.ts` is dissolved unless it is a package export entry.

**Blocked by:** 0239, 0240, 0279, 0277.

**Status:** ready-for-agent

- [x] The file-suffix detector reports nothing in the website.
- [x] No import path breaks and no export is renamed beyond the spec's calls.

> **Note, 2026-11-05.** The file-suffix detector is ticket 0277's. Until it exists, a script applies its definition: a tracked `.ts` file under `src/` that is not a test, an `index.ts`, a declaration or a framework file must carry a suffix from the TypeScript or Nest table. Under that check, the website keeps three files without one: `template/i18n/routing.ts`, `navigation.ts` and `request.ts`. They are next-intl's file conventions, and rule 06's tree names them. The calls this ticket made:
> - The metadata companions are `page-metadata.map.ts`, `layout-metadata.map.ts` and `not-found-metadata.map.ts`, and `generateStaticParams` is `layout-static-params.compute.ts`.
> - The response builder is `book-demo-reply.map.ts`.
> - The four guards that are not Zod schemas, `triple` among them, are `.compute`, and `socialImages`, which awaits the parent metadata, is `.adapter`.
> - `drawerFocus` is `.state` and `ticketCardFormat` is `.factory`.
> - `positions`, the three content modules and the book-a-demo error codes are `.data`. The error codes leave `form.types.ts`.
> - `platform-products` splits into its data and a `platformProduct` lookup.
> - The content modules and `page-metadata` lose their re-exports, and their consumers import each helper from its own module.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check 04-typescript`
