# 0056: Shared, shorter scroll hint on platform pages

**What to build:** the blurred bottom band (`ScrollHint`) is half as tall and is rendered by a shared `platform/layout.tsx`, so every `platform/*` product page shows it without wiring it page by page.

**Status:** ready-for-agent

### Scroll hint

- [ ] `ScrollHint` height is `h-8` and `sm:h-12`.
- [ ] `platform/layout.tsx` renders `ScrollHint` after its children; the Microruntimes page no longer renders it itself.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`.
- `/en/platform/vision/`, `/en/platform/finops/` and `/en/platform/microruntimes/`: the band shows near the bottom and fades at the end of the page.
