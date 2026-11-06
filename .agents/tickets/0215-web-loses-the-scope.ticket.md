# 0215: the shared primitives package is `web`

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The shared primitives package becomes the internal package `web`. The website and the console depend on it as `web`, import `web/ui/...`, `web/utilities`, `web/vendors/...` and `web/styles.css`, and list `web` in their transpiled packages; the package's own imports of its subpaths follow; the Biome check that names the package, and the Tailwind, TypeScript and React rules, name `web`. Biome keeps grouping the new specifiers as package imports.

**Blocked by:** 0214.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `packages/web` is named `web`; the website and the console depend on it through `workspace:*` and transpile `web`.
- [ ] No source imports `@langdrift/web`; both apps build, typecheck and pass their tests.
- [ ] `drift check biome` passes and the import groups of a file importing `web/ui` and a relative sibling are unchanged.
- [ ] The Tailwind, TypeScript and React rules and the `drift` Biome check name `web`.
- [ ] `git grep -n '@langdrift/web'` finds only workflow history.

## Dev part

**Override:** No part qualifies: a package rename across its consumers (policy, Categories: never a rename). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile && pnpm turbo run build typecheck test check --filter=web --filter=website --filter=console`
- `platform/cli/drift check biome`
