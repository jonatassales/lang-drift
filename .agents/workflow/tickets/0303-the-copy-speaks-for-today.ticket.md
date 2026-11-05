# 0303: the copy speaks for today's business

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** The footer says "Agent-assisted delivery, accountable to your team's budget." in every locale. "Illustrative" and its translations leave every catalog and comment of the website and `packages/web`. No Vision content remains. "The ex CLI" and "ex connected agent" become the LangDrift CLI in prose, while terminal lines keep the `ex` command.

**Blocked by:** 0302.

**Status:** ready-for-agent

- [ ] `git grep -n -i -E 'illustrative|ilustrativ|示意|イメージ|minh họa'` over `apps/website/src` and `packages/web/src` matches nothing.
- [ ] No catalog says "vision" in the product sense, and the footer line is the new tagline.
- [ ] No prose names "the ex CLI"; terminal lines still run `ex`.
- [ ] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
