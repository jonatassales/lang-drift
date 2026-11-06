# 0184: Add proxy matcher to apps/console to exclude static assets

**What to build:** In `apps/console/src/proxy.ts`, `export const config = { matcher: [...] }` was missing. As documented in `apps/console/README.md`, the proxy must run on every path except `/api`, `/_next`, and root static files (`favicon.ico`, `icon.svg`, `apple-icon.png`, `robots.txt`). Without this matcher, Next.js routed all `/_next/static/...` CSS and JS asset requests through `proxy.ts`, where `namesNoScreen` flagged `.css` extensions as non-entity IDs and rewrote asset requests to the 404 HTML page, breaking all styling on `https://console.langdrift.md`. Add `config.matcher` to `apps/console/src/proxy.ts` and add a test suite in `apps/console/src/proxy/matcher.test.ts`.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Add `config.matcher` to `apps/console/src/proxy.ts` excluding `/api`, `/_next`, root assets, and dotted files.
- [x] Add test suite in `apps/console/src/proxy/proxy-matcher.test.ts` verifying path matching and exclusion contracts.
- [x] Verify console test suite passes with `pnpm --filter console test`.
- [x] Verify full repository check, typecheck, test, and build pass with `pnpm turbo run check typecheck test build`.
