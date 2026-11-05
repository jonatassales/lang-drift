# 0274: every color and static style reads from a role or a utility

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** The switcher's translucent surfaces take the role with Tailwind's opacity modifier. The app theme gains `primary-highlight` and `primary-shade`, built by the tokens build in both schemes, and the orb's fallback gradient becomes a radial utility through them. The brand's `tone` prop is removed and the wordmark takes the surrounding text color; its layout and typography become utilities, as do the text morph's static styles. A `style` prop keeps only a runtime measurement.

**Blocked by:** 0267.

**Status:** ready-for-agent

- [x] The tokens test asserts the two new roles in both schemes.
- [ ] No `color-mix(`, `var(--` color or static `style` value remains in `packages/web`.
- [x] Captures differ only in the orb's gradient and the wordmark's color, as the spec accepts.

> **Note, 2026-11-05.** The captures compare the computed styles of the header wordmark, the text links, the pill CTAs, the footer and the footer orb. They cover the home, FinOps, Microruntimes and pricing pages at 1280px and 375px, in dark and light. The only colour that changes is the wordmark in light, which moves to the foreground role. The second criterion stays open on `hero-shader` and `chroma-text`, which wait for the owner.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`

> **Note, 2026-11-05.** Two components stay outside this ticket: `chroma-text` keeps its five-stop rainbow gradient, which Tailwind's from/via/to utilities cannot express without dropping stops, and `hero-shader` keeps the CSS variables its WebGL program reads and its two-layer fallback gradient. Both still hold `var(--` or `color-mix(`, so the second criterion stays open until the owner decides how they should be expressed; it also gates 0276's detector. Everything else the criterion names is a role or a utility, and `elevated-frame` went with it.
