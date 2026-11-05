# 0044: SmoothUI motion on the home and Microruntimes heroes

**Spec:** None

**What to build:** the owner's review of the home page asked for three SmoothUI components. The home product cards glow under the pointer (Glow Hover Card), the hero's product list drops in line by line (Short Slide Down) and the Microruntimes hero types an illustrative `ex` session (Mac Terminal). Each component is vendored in `@langdrift/web` and drawn with the tokens' roles.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The Features Grid's cards glow in the brand purple inside a circle that follows the pointer; the glowing copy of the grid is inert, so no link is focusable or announced twice. Reduced motion shows no glow.
- [ ] The hero's products drop in one after another and stay. Each name has a one-line description beneath it, translated in every locale, in solid foreground (white in the dark scheme, black in the light one) with no transparency and no hover effect.
- [ ] The descriptions fit on one line beside the sun from 1024px up in every locale.
- [ ] The Microruntimes hero shows a terminal above the lifecycle demo that types `ex doctor` and `ex agent run` for `checkout-preview`, matching the demo's illustrative values. The line kinds live in the content module and the text in the catalog; reduced motion shows the whole session at once.

### Verification

- The home page and `/platform/microruntimes/` render in both color schemes without console errors.
- `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test` and `drift check localization`.
