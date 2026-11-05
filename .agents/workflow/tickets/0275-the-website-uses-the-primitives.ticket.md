# 0275: the website uses the shared primitives and one cn

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** The demo form submit, the trial plan's call to action and the six text links consume the pill `Button` and `TextLink`; the `!` overrides and "Transitional" comments go. The footer's measurements become utilities. The website's `<html>` takes `scheme-light dark:scheme-dark` and the header stops writing `colorScheme`. The website's `cn.fmt.ts` and its `clsx` and `tailwind-merge` dependencies are removed.

**Blocked by:** 0229, 0272, 0274.

**Status:** ready-for-agent

- [x] No `smooth-` class and no "Transitional" comment remains in the website.
- [x] The website manifest lists neither `clsx` nor `tailwind-merge`.
- [x] Captures of the affected pages match, recorded under `.audits/`.

> **Note, 2026-11-05.** The captures compare computed styles, not pixels, before and after, in a local audit. They cover every element this ticket moved, the header wordmark and `<html>`, on the home, FinOps, Microruntimes and pricing pages at 1280px and 375px, in dark and light.
> - Every width and height matches.
> - Only notations differ: the underline colour as `oklab` instead of `srgb`, the CTA radius as `rounded-full` instead of 999px, and a zero-width border style. The pill CTAs keep SmoothUI's 7px gap and wrapping, and the submit does not dim while submitting.
> - Hover and focus states and the Book a Demo form, paused, are not measured.
> - `.audits/` is missing from `.gitignore`, although rule 10 says it is ignored.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
