# 0304: platform headings carry the product icon

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** The icon tile Microruntimes draws before its title becomes one component in the platform folder, taking an icon, and Microruntimes uses it. FinOps and Safe Seat replace the "LangDrift FinOps" eyebrow with the tile holding the FinOps icon and the word "FinOps"; Safe Seat keeps its Coming soon badge on that line. Their sentence h1 is unchanged.

**Blocked by:** 0302.

**Status:** ready-for-agent

- [x] Microruntimes renders unchanged through the shared tile.
- [x] FinOps and Safe Seat show the FinOps icon and "FinOps" above their h1, in every locale.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`

> **Note, 2026-11-05.** FinOps and Safe Seat no longer show the FinOps icon and name as a label above a sentence title: every product page now opens with the same heading, its own icon and product name at the Microruntimes size, so Safe Seat takes its own mark, and the sentence follows as a tagline (ticket 0311).
