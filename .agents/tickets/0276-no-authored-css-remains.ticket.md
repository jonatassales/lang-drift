# 0276: no hand-written stylesheet remains

**Spec:** .agents/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** `smoothui.css` and its import go; the package stylesheet keeps the Tailwind import, the tokens import and the `@source` directive. `drift quality check tailwind` gains rows for a `smooth-` class in TypeScript, a stylesheet import in TS or TSX, an arbitrary color value containing `color-mix(`, `var(--`, `white` or `black`, a static value in a `style` object, and `clsx` or `tailwind-merge` in an app manifest.

**Blocked by:** 0273, 0275.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] No `.css` file holds a selector or element rule.
- [ ] `platform/cli/drift quality check tailwind` reports nothing.

> **Note, 2026-11-05.** `smoothui.css` and its import left anyway, with the CI fix of #146: once 0275 moved the website to the shared primitives, `drift quality dead-code` reported its three classes as unused. The drift rows this ticket adds stay won't fix.

## Dev part

**Override:** No part qualifies: detector rows are patterns in the existing table, configuration rather than a problem with its own goal (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check tailwind`
- `pnpm turbo run check typecheck test`
