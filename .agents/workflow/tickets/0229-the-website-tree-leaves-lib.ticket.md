# 0229: the website's lib folder dissolves into its route owners

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** Every module under the website's `src/app/lib/` moves to the folder of the nearest route segment whose subtree holds every file that imports it, files of one concern together in a folder named after the concern, as the spec's website owners list places them. The JSON body modules both API routes share move into one concern folder under `api/`. The website's aliases follow the tree (`@layout/*`, `@routes/*`, `@api/*` in; `@components/*`, `@data/*`, `@event-testdata/*` out) and the Biome import matcher lists exactly the aliases the app declares. Tests move with their module and keep their file names. Nothing a visitor sees changes.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `apps/website/src/app/lib/` no longer exists and no file imports through a removed alias.
- [ ] Every moved module sits under the segment the spec's website owners list names, and the early-access environment computes move unchanged with their route.
- [ ] The import matcher in `biome.json` names exactly the aliases the two apps declare.
- [ ] `pnpm turbo run check typecheck test build --filter=website` and `platform/cli/drift check localization` pass; a move changed imports and nothing else.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift check localization`
- `platform/cli/drift check biome`
