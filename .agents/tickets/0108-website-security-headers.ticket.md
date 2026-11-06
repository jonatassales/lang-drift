# 0108: The website answers with a complete security header set

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** Every website response carries the spec's static Content-Security-Policy, together with:

- HSTS with includeSubDomains and preload;
- `nosniff` and `Referrer-Policy: strict-origin-when-cross-origin`;
- `X-Frame-Options: DENY`;
- the Permissions-Policy;
- `Cross-Origin-Opener-Policy: same-origin`.

The powered-by header is off. Every page stays statically prerendered. Development adds only what the development server needs, and preview builds add only the Vercel Toolbar's origin. The README's steps for reopening sign-up name the Turnstile origins the policy will then need.

**Blocked by:** 0107.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A page, an API route and a static file each carry every header with the spec's value, and no response carries `X-Powered-By`.
- [ ] The policy has no `upgrade-insecure-requests`.
- [ ] Development adds only `'unsafe-eval'` and the development websocket; preview adds only `https://vercel.live`; production adds nothing.
- [ ] The build still prerenders every locale page statically.
- [ ] No page logs a CSP violation in the browser, the Microruntimes editor and the hero shader included.

### Verification

- The website unit suite, with a Next config test through `unstable_getResponseFromNextConfig`; the `next build` route table; `curl -sI` against `next start`; the browser console on every page.
