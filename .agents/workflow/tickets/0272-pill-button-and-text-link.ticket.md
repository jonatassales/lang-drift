# 0272: the pill button and the text link are shared primitives

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** `Button` in `packages/web` gains a `pill` variant that reproduces the SmoothUI primary button, composed through `cn()` so a consumer's overrides win without `!`. A `TextLink` primitive reproduces the SmoothUI text link, renders an anchor by default and lends its look to its one child through `asChild`, exporting its props interface.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] A `Button` test renders the pill variant with the plans section's overrides and asserts the override wins over the conflicting utility.
- [x] A `TextLink` test covers its anchor form and its `asChild` form.
- [x] `pnpm turbo run check typecheck test build --filter=web --filter=tokens` passes.

> **Note, 2026-11-05.** `TextLink` also carries the site's focus outline, as the pill does, so its callers do not repeat it.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
