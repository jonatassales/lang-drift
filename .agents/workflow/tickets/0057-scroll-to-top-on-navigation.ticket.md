# 0057: Navigating between pages always starts at the top

**What to build:** moving between pages keeps the scroll position of the previous page. Each page renders its own header and main, so Next treats the new page's first element as already in view and does not scroll. A `ScrollToTop` Client Component in the locale layout resets the scroll to the top on every pathname change.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Scroll

- [ ] `ScrollToTop` scrolls to the top, without smooth animation, whenever the pathname changes.
- [ ] It does nothing on the initial load and when the URL has a fragment.
- [ ] The locale layout renders it once.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`.
- Scroll down on `/en/platform/vision/`, follow the header to `/en/platform/finops/` and the logo to `/en/`: each page opens at the top. A `#fragment` link still scrolls to its section.
