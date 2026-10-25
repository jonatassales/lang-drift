# 0051: Hardware, the closing report and the pause and replay control

**Spec:** .agents/workflow/specs/0007-microruntimes-hero.spec.md

**What to build:** the card completes. A hardware table shows CPU and memory bounded against suggested with headroom bars, GPU and disk without a number and with the note that they are not `runtime.md` fields, and a line saying no price source is configured so the figures stay in units. The closing report fills with the bounded figures, the suggested figures, the suggested materialization and a final row pending team authorization, and the screen reader hears it once. A text button pauses the session and, once it ends, replays it from the start.

**Blocked by:** 0050.

**Status:** ready-for-agent

### Hardware

- [ ] A real table with resource, bounded and suggested columns: CPU and memory bounded after the bounded line and suggested after the size line, each with a headroom bar in the suggested cell; GPU `none | none` and disk `— | —` after the limits warning, both with the note `not a runtime.md field`; the line `price basis: none configured · units only` after the estimate line. No currency appears anywhere on the card.

### Report

- [ ] A detail list with bounded (after the bounded line), suggested (stepping with the running sums of ticket 0050), suggested resource (`5 one-replica Deployments`) and final (`1,450m · 2,304 MiB · pending team authorization`, set stronger), the last two after the plan line, followed by the footer line naming the `runtime.md` draft proposed to the console.
- [ ] The report footer is the card's only live region, polite and atomic, so the screen reader hears the report once when it fills and nothing per terminal line.

### Control

- [ ] `SessionControl` (client): a text button in the card footer that reads pause while the session plays and replay once the last line prints; pause freezes the terminal and the card together; replay resets the count, clears the pause and increments the run so the terminal remounts and retypes while the card returns to dashes in the same frame. It is reachable with Tab, operable with Enter and Space, has a visible focus ring, and is not rendered under reduced motion.

### Tests

- [ ] The `TicketCard` test grows: after the bounded line the table and the report show the envelope; after the limits line GPU and disk show their notes; at the full count the report shows the suggested totals and the pending authorization text and the control reads replay; clicking replay resets the store to count zero with the run incremented.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `pnpm turbo run build`, `drift check localization`, `drift check dead-code`.
- In the built-in browser: the hero in `/en`, `/pt-br`, `/zh-hant`, `/ja` and `/vi` with the same structure; light and dark; 375px without horizontal overflow and with the ticket head in the first screen; reduced motion shows the final card and the full transcript with no control; Tab reaches pause and replay with a visible focus ring and Enter toggles them; the first frame without JavaScript shows the heading, the terminal frame and the reserved card without a layout shift; the screen reader announces the report once at the end; the Vision and FinOps pages are unchanged.
