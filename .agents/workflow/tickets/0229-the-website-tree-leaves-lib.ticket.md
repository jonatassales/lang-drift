# 0229: the website's lib folder dissolves into its route owners

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** Every module under the website's `src/app/lib/` moves to the folder of the nearest route segment whose subtree holds every file that imports it, files of one concern together in a folder named after the concern, as the spec's website owners list places them. The JSON body modules both API routes share move into one concern folder under `api/`. The website's aliases follow the tree (`@layout/*`, `@routes/*`, `@api/*` in; `@components/*`, `@data/*`, `@event-testdata/*` out) and the Biome import matcher lists exactly the aliases the app declares. Tests move with their module and keep their file names. Nothing a visitor sees changes.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `apps/website/src/app/lib/` no longer exists and no file imports through a removed alias.
- [x] Every moved module sits under the segment the spec's website owners list names, and the early-access environment computes move unchanged with their route.
- [x] The import matcher in `biome.json` names exactly the aliases the two apps declare.
- [x] `pnpm turbo run check typecheck test build --filter=website` and `platform/cli/drift check localization` pass; a move changed imports and nothing else.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift check localization`
- `platform/cli/drift check biome`

> **Note, 2026-11-05.** Two items are obsolete: spec 0046 already removed `@event-testdata/*` with the early-access route, and since then only the book-a-demo route reads a JSON body, so its body modules stay in `api/book-a-demo/body/` rather than a shared folder. The early-access modules the owners list names no longer exist. `product-badge`, which the list does not place, goes to the locale layout, the nearest segment holding its importers: the header and its Platform menu, the pricing plans section and the Safe Seat page.
