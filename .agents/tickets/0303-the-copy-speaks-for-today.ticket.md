# 0303: the copy speaks for today's business

**Spec:** .agents/specs/0046-the-website-focus-pass.spec.md

**What to build:** The footer says "Agent-assisted delivery, accountable to your team's budget." in every locale. "Illustrative" and its translations leave every catalog and comment of the website and `packages/web`. No Vision product content remains; the logo's ring, the brand's vision, stays. "The ex CLI" and "ex connected agent" become the LangDrift CLI in prose, while terminal lines keep the `ex` command.

**Blocked by:** 0302.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `git grep -n -i -E 'illustrativ|ilustrativ|示意|イメージ上|イメージです|minh họa'` over `apps/website/src` and `packages/web/src` matches nothing; イメージ alone also means a container image and stays.
- [x] No catalog says "vision" in the product sense, and the footer line is the new tagline.
- [x] No prose names "the ex CLI"; terminal lines still run `ex`.
- [x] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
