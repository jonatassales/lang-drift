# 0316: the pricing plans sit in one row

**Spec:** None. The owner's report of 2026-11-05: on a desktop screen the pricing page shows Starter and Plus side by side and Pro alone on a second row.

**What to build:** From the `lg` breakpoint, Starter, Plus and Pro sit in one row of three equal columns, at every width up to the page's 1240px frame. Below `lg` they stack in one column as before.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] At 1024, 1120 and 1440px the three cards share one row, with equal widths and heights and no horizontal scroll.
- [x] Below `lg` the cards stack as before.

> **Note, 2026-11-05.** The grid was `repeat(auto-fit, minmax(280px, 390px))`. With a definite maximum, `auto-fit` counts its tracks by that 390px, so any content width under 1206px (about 1262px of viewport) held two tracks and Pro wrapped. Three explicit columns hold the row: each card is 311px wide at 1024px and 383px at 1440px.

## Dev part

**Override:** No part qualifies: a one-class layout fix holds no core-domain logic, hard part or new stack (policy, Categories). 2026-11-05.
