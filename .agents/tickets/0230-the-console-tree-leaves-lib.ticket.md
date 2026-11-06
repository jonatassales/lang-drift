# 0230: the console's lib folder dissolves into its route owners

**Spec:** .agents/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** Every module under the console's `src/app/lib/` moves beside the route segment that owns it, as the spec's console owners list places them: the locale layout's API client and domain types, the `(menu)` layout's navigation and shared views, the `runtimes`, `products`, `tickets` and `executions` views beside their segments. The bridge to console-stubs (the request, the base resolution, the host and cross-site checks, the hosted check and the server read builder) moves to `src/infra/console-stubs/`. The console's aliases follow the tree (`@layout/*`, `@auth/*`, `@menu/*`, `@product/*`, `@api/*`, `@infra/*` in; `@views/*`, `@components/*`, `@ui/*`, `@data/*` out) and the Biome matcher follows them. Nothing a console user sees changes.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] `apps/console/src/app/lib/` no longer exists and no file imports through a removed alias.
- [ ] The console-stubs bridge lives under `apps/console/src/infra/console-stubs/`.
- [ ] The import matcher in `biome.json` names exactly the aliases the two apps declare.
- [ ] `pnpm turbo run check typecheck test build --filter=console` and `platform/cli/drift check localization` pass; a move changed imports and nothing else.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift check localization`
- `platform/cli/drift check biome`
