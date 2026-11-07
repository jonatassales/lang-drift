# 0046: The terminal reports its progress and can be paused

**Spec:** .agents/specs/0007-microruntimes-hero.spec.md

**What to build:** a page can drive a sibling from the terminal's own clock. `MacTerminal` tells its caller how many lines it has completed as the session plays, freezes when asked to, and prints a warning line in the warning status color. The home page, which passes none of this, keeps playing exactly as before.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Terminal

- [ ] An optional `onProgress` callback receives the number of completed lines whenever that number changes, from an effect, so typed commands, printed output and the reduced-motion branch that completes every line at once all report through the same path. It never fires during server rendering.
- [ ] An optional `paused` flag joins the existing play condition beside viewport and page visibility: while it is set, no line advances and no character types.
- [ ] A `warning` line type prints in the warning status color.
- [ ] Every new prop is optional; the home page's terminal renders and plays unchanged.

### Tests

- [ ] A jsdom unit test for `MacTerminal`, with an `IntersectionObserver` stub that reports intersection, a `matchMedia` stub and fake timers: `onProgress` is called with 1, then 2, as lines complete; setting `paused` stops the count; under reduced motion the callback reports the full line count at once.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=@langdrift/web`.
- The home page hero terminal still types and prints its session in the built-in browser.
