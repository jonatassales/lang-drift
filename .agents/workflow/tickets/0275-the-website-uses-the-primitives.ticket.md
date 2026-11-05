# 0275: the website uses the shared primitives and one cn

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** The demo form submit, the trial plan's call to action and the six text links consume the pill `Button` and `TextLink`; the `!` overrides and "Transitional" comments go. The footer's measurements become utilities. The website's `<html>` takes `scheme-light dark:scheme-dark` and the header stops writing `colorScheme`. The website's `cn.fmt.ts` and its `clsx` and `tailwind-merge` dependencies are removed.

**Blocked by:** 0229, 0272, 0274.

**Status:** ready-for-agent

- [ ] No `smooth-` class and no "Transitional" comment remains in the website.
- [ ] The website manifest lists neither `clsx` nor `tailwind-merge`.
- [ ] Captures of the affected pages match, recorded under `.audits/`.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
