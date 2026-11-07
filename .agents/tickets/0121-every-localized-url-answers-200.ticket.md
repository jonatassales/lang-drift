# 0121: Every localized URL answers 200

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** Find why every `pt-br` and `zh-hant` URL looped on 308 under a local production server, and make sure every sitemap URL answers 200.

**Blocked by:** 0107.

**Status:** wontfix

**Outcome.** The loop is Next's, and only on a server bound with `--hostname 127.0.0.1`. NextURL renames a loopback host `localhost` (`server/web/next-url.js`), so next-intl's rewrite to the locale route names `localhost`. The router relativizes that rewrite against the URL it received, which keeps `127.0.0.1` (`resolve-routes.js`), and an absolute result counts as an external rewrite: the router proxies the request to itself. The second pass meets the trailing-slash redirect at `/pt-BR` and the proxy's case redirect at `/pt-BR/…`, so the browser lands where it started. `en`, `ja` and `vi` never rewrite. The website's own `next start`, unbound, answers 200 on all 30 sitemap URLs, through `localhost` and `127.0.0.1` alike, and Vercel routes rewrites itself. A first fix patched Next through pnpm; it was dropped, because carrying a framework patch for a local-only artifact costs more than it protects. The website README warns against binding its production server to `127.0.0.1`.

- [x] Each of the 30 sitemap URLs answers 200 from `next start` as the website's script runs it, with no redirect.
- [x] `/PT-BR/…` still answers 308 to `/pt-br/…`, `/zh-cn/…` still answers 307 to `/en/…?availability=zh-CN`, and `/` still answers 307 to a locale.
- [x] The cause is recorded here and in the spec.

### Verification

- `next build`, then `curl` against `next start --hostname localhost` on every sitemap URL and on the redirect cases above.
