# 0274: every color and static style reads from a role or a utility

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** The switcher's translucent surfaces take the role with Tailwind's opacity modifier. The app theme gains `primary-highlight` and `primary-shade`, built by the tokens build in both schemes, and the orb's fallback gradient becomes a radial utility through them. The brand's `tone` prop is removed and the wordmark takes the surrounding text color; its layout and typography become utilities, as do the text morph's static styles. A `style` prop keeps only a runtime measurement.

**Blocked by:** 0267.

**Status:** ready-for-agent

- [ ] The tokens test asserts the two new roles in both schemes.
- [ ] No `color-mix(`, `var(--` color or static `style` value remains in `packages/web`.
- [ ] Captures differ only in the orb's gradient and the wordmark's color, as the spec accepts.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
