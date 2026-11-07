# 0337: the reveal value, check mark and card section move to web

**Spec:** .agents/specs/0068-the-session-panel-primitives-move-to-web.spec.md

**What to build:** `packages/web` gains three dumb pieces, each with its own `web/ui/<name>` entry:

- `RevealValue`, the value that reveals itself as a dash, then a pulsing skeleton, then the value fading in;
- `CheckMark`, the ✓ in the success role with its screen-reader label;
- `CardSection`, the block under a hairline rule.

The website's session cards, rows, summaries and workflow parts compose these three in place of their local card value, read mark and card section, and the local copies are deleted. The hero panel looks and behaves exactly as before.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Delivered:** commit 8bdfbba, merged in PR #166, under this ticket's first number. Its `Ticket: 0334` trailer means this ticket, 0337, not ticket 0334 (the agentic word).

**Alignment:** change, in code only: no simulated line or card field changes. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] Before its line plays, each card value reads as `—`, with "Not set yet" for screen readers. While it resolves, it pulses, with "Resolving" for screen readers. Once its line plays, it fades in. All three states are drawn by `web/ui/reveal-value`.
- [ ] A read document and a finished workflow stage show `✓`, with "Read" or the stage's name for screen readers, drawn by `web/ui/check-mark`.
- [ ] Every hairline section of the Ticket, Analysis, Plan and Workflow cards is drawn by `web/ui/card-section`, with the same spacing as before.
- [ ] The website no longer has its card value, read mark or card section components. The reveal states' type is owned by `web/ui/reveal-value`, and the website's `revealState` compute returns it.
- [ ] The website's session cards and session terminal tests pass with no changed assertion.
- [ ] `web/ui/check-mark` and `web/ui/card-section` have unit tests described by their paths.
- [ ] `pnpm turbo run check`, `typecheck` and `test` pass.

## Dev part

**Category:** `core-domain`. `RevealValue` is where "unknown is never zero" becomes visible: what a not-yet-known value shows, what a resolving one shows, and what assistive technology hears for each. The owner must be able to explain and debug it whenever a card shows a figure it should not.

**Paths:**

- `packages/web/src/ui/reveal-value/reveal-value.tsx`
- `packages/web/src/ui/reveal-value/reveal-value.unit.tsx`

**Boundary:** the agent writes these, outside the paths:

- `packages/web/src/ui/reveal-value/reveal-value-state.type.ts`, which exports `type RevealValueState = 'unset' | 'resolving' | 'set'`;
- `packages/web/src/ui/reveal-value/reveal-value-labels.type.ts`, which exports `interface RevealValueLabels { notSet: string; resolving: string }`;
- the barrel, which re-exports `RevealValue` and `RevealValueProps` from `./reveal-value` with both types, and the package's `./ui/reveal-value` export entry.

The website imports `RevealValue` from `web/ui/reveal-value`. The owner's component takes `RevealValueProps`:

- `state: RevealValueState`;
- `labels: RevealValueLabels`;
- `children?: ReactNode`;
- `className?: string`;
- `width?: string`, the skeleton's width class, defaulting to `w-24`.

**Goal:** write the piece and its test. It draws each state exactly as the website's card value drew it, with the same classes, ARIA and motion.

- [ ] `unset` renders a decorative `—` and `labels.notSet` for screen readers, inside a span carrying `className`.
- [ ] `resolving` renders a block with a decorative pulsing bar of `width` (no pulse under reduced motion) and `labels.resolving` for screen readers.
- [ ] `set` renders the children in a span that fades in, merged with `className` through `cn()`.
- [ ] The unit test, described as `web/ui/reveal-value/reveal-value`, covers the three states, both labels and the default and given widths.

**Override:** the owner does not write this part; the agent writes its paths. Reason: the owner's files did not reach the worktree after two rechecks, and the owner then chose to override the part rather than write it, in the session, 2026-11-05.
