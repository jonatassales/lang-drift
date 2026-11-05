# 0278: every import past a sibling uses an alias

**Spec:** .agents/workflow/specs/0039-typescript-rules-compliance.spec.md

**What to build:** Both Next apps declare `@proxy/*` for `src/proxy/*`, and the Biome matcher and the `drift check biome` fixtures name it. The proxies, the book-a-demo route's interceptors, the website header and the i18n request configurations and tests import past their siblings through an alias. The nine hand-placed import breaks are fixed.

**Blocked by:** 0229, 0230, 0233.

**Status:** ready-for-agent

- [ ] No `./dir/...` import remains in an app.
- [ ] `drift check biome` proves `web` and `@proxy/*` sort in the alias group.
- [ ] `pnpm turbo run check typecheck test build --filter=website` and `pnpm turbo run check typecheck test build --filter=console` pass.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift check biome`
- `pnpm turbo run check typecheck test build --filter=website`
- `pnpm turbo run check typecheck test build --filter=console`
