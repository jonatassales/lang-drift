# 0307: the Platform menu lists five products in two columns

**Spec:** .agents/specs/0046-the-website-focus-pass.spec.md

**What to build:** The Platform menu lists Microruntimes, FinOps, LangDrift CLI, CI Gates and Safe Seat. The CLI and CI Gates link to their pages; CI Gates loses its Coming soon badge and Safe Seat keeps its badge. On desktop the panel is a two-column grid; the mobile drawer lists the same products in one column. The header catalog gains the CLI's name and description in five locales.

**Blocked by:** 0305, 0306.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The menu test asserts the order, the five links and one Coming soon badge.
- [x] The desktop panel shows two columns at 1280px; the drawer one column at 375px.
- [x] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
