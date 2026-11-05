# 0275: the website uses the shared primitives and one cn

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** The demo form submit, the trial plan's call to action and the six text links consume the pill `Button` and `TextLink`; the `!` overrides and "Transitional" comments go. The footer's measurements become utilities. The website's `<html>` takes `scheme-light dark:scheme-dark` and the header stops writing `colorScheme`. The website's `cn.fmt.ts` and its `clsx` and `tailwind-merge` dependencies are removed.

**Blocked by:** 0229, 0272, 0274.

**Status:** ready-for-agent

- [x] No `smooth-` class and no "Transitional" comment remains in the website.
- [x] The website manifest lists neither `clsx` nor `tailwind-merge`.
- [x] Captures of the affected pages match, recorded under `.audits/`.

> **Note, 2026-11-05.** The captures under `.audits/0275/` compare computed styles, not pixels. Every element the ticket moved was read before and after, at 1280px in the dark theme, on the home, FinOps, Microruntimes and pricing pages: the six text links, the three plan calls to action, the footer and `<html>`.
> - Everything matches except two computed values with no visible effect. The CTA's radius reads `rounded-full` instead of 999px. Its gap is the 12px of `gap-3`, which SmoothUI's 7px used to beat and `justify-between` hides.
> - The Book a Demo form is not rendered while sign-up is paused, so its unit tests cover the pill submit.
> - `.audits/` is missing from `.gitignore`, although rule 10 says it is ignored, so it was kept out of the commit by hand.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
