# 0338: the expandable card moves to web

**Spec:** .agents/specs/0068-the-session-panel-primitives-move-to-web.spec.md

**What to build:** `packages/web` gains `ExpandableCard`, under its own `web/ui/expandable-card` entry. It is the card that expands: a header that names the card and carries its one-line summary, over a body that slides open and shut. Which card is open is decided from outside. The website's session stack composes the four cards (Ticket, Analysis, Plan, Workflow) from it in place of its local session card, which is deleted. The hero panel opens, collapses and reopens by hand exactly as before.

**Blocked by:** None (can start immediately). Like 0337 and 0339, it adds an export entry to the web package's manifest.

**Status:** ready-for-agent

**Alignment:** change, in code only: no simulated line or card field changes. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The headers "Ticket", "Analysis", "Plan" and "Workflow" each carry their summary line and a caret, drawn by `web/ui/expandable-card`.
- [ ] The card of the act the session reaches opens and the others collapse. A click opens a card by hand until the next act takes it back. The session cards tests prove both with no changed assertion.
- [ ] A collapsed card's body stays mounted, inert and hidden from assistive technology. Its header reports `aria-expanded` and controls the region that it labels.
- [ ] The website no longer has its session card component.
- [ ] `pnpm turbo run check`, `typecheck` and `test` pass.

## Dev part

**Category:** `hard-part`. This ticket's difficulty is a controlled disclosure that keeps its body mounted: the header and region ids, `aria-expanded` and `aria-controls`, the inert and `aria-hidden` shut body, and the grid-row slide with no motion under reduced motion.

**Paths:**

- `packages/web/src/ui/expandable-card/expandable-card.tsx`
- `packages/web/src/ui/expandable-card/expandable-card.unit.tsx`

**Boundary:** the agent writes the barrel, which re-exports `ExpandableCard` and `ExpandableCardProps` from `./expandable-card`, and the package's `./ui/expandable-card` export entry. The website's session stack calls `<ExpandableCard>` with `ExpandableCardProps`:

- `body: ReactNode`;
- `expanded: boolean`;
- `label: string`;
- `onToggle: (event: MouseEvent<HTMLButtonElement>) => void`;
- `summary: ReactNode`;
- `value?: string`, set on the header button as `data-value`, so one handler serves every card.

The stack's handler reads `event.currentTarget.dataset.value` in place of `dataset.phase`.

**Goal:** write the piece and its test, with the classes, ARIA and motion of the website's session card.

- [ ] The header button carries the label through `Eyebrow`, the truncated summary in mono 12px and a caret that turns 180° when expanded, with `aria-expanded`, `aria-controls` and `data-value`.
- [ ] The body region is labelled by the header. Shut, it is `inert` and `aria-hidden="true"` and stays in the DOM.
- [ ] The unit test, described as `web/ui/expandable-card/expandable-card`, covers both states, the ids that link header and region, and `onToggle` receiving a click whose target carries `data-value`.
