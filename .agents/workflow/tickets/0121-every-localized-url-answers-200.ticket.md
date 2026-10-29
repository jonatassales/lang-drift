# 0121: Every localized URL answers 200

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** Under the production server, every `pt-br` and `zh-hant` URL loops on 308. The proxy hands next-intl the lowercase prefix, next-intl rewrites to the locale route (`/pt-BR` for the home page, without the trailing slash Next is configured for), Next redirects, and the proxy lowercases again. Find the root cause in Next 16.4.3's routing and next-intl 4.14.5, then fix it where it starts. Every URL in the sitemap must then answer 200 with no redirect, for every locale. What must keep working:

- the uppercase-to-lowercase prefix redirect (308);
- the zh-CN fallback to English (307);
- the negotiating redirect at `/` (307);
- the locale cookie.

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] Each of the 30 sitemap URLs answers 200 from `next start`, with no redirect.
- [ ] `/PT-BR/…` still answers 308 to `/pt-br/…`, `/zh-cn/…` still answers 307 to `/en/…?availability=zh-CN`, and `/` still answers 307 to a locale.
- [ ] A test fails if the rewrite target for a mixed-case locale stops being one Next serves without a redirect, at the highest seam that reproduces the loop.
- [ ] The cause is explained in the commit body.

### Verification

- The website unit suite; `next build`, then `curl` against `next start` on every sitemap URL and on the redirect cases above.
