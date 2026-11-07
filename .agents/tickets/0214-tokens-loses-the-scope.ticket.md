# 0214: the design tokens package is `tokens`

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The design tokens package becomes the internal package `tokens`. The website, the console and the web package depend on it as `tokens` and import its stylesheet as `tokens/website.css`; the tokens test resolves it under its new name; the Tailwind rule names `tokens` where it named `@langdrift/tokens`. Both apps build with the same theme as before.

**Blocked by:** 0211.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `packages/tokens` is named `tokens`; the website, the console and `packages/web` depend on it through `workspace:*`.
- [ ] Every Tailwind entry and the web package's tokens stylesheet import `tokens/website.css`; both apps build and their pages render with the same colors.
- [ ] The tokens test passes under the new name.
- [ ] The Tailwind rule and the context document name `tokens`.
- [ ] `git grep -n '@langdrift/tokens'` finds only workflow history.

## Dev part

**Override:** No part qualifies: a package rename across its consumers (policy, Categories: never a rename). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile && pnpm turbo run build typecheck test check --filter=tokens --filter=@langdrift/web --filter=website --filter=console`
