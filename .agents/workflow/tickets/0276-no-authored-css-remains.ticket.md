# 0276: no hand-written stylesheet remains

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** `smoothui.css` and its import go; the package stylesheet keeps the Tailwind import, the tokens import and the `@source` directive. `drift quality check tailwind` gains rows for a `smooth-` class in TypeScript, a stylesheet import in TS or TSX, an arbitrary color value containing `color-mix(`, `var(--`, `white` or `black`, a static value in a `style` object, and `clsx` or `tailwind-merge` in an app manifest.

**Blocked by:** 0273, 0275.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] No `.css` file holds a selector or element rule.
- [ ] `platform/cli/drift quality check tailwind` reports nothing.

## Dev part

**Override:** No part qualifies: detector rows are patterns in the existing table, configuration rather than a problem with its own goal (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check tailwind`
- `pnpm turbo run check typecheck test`
