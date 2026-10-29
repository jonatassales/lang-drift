# 0113: The console sends a strict nonce CSP and a full header set on every document

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** The console's proxy builds a per-request nonce policy, using a pure builder, and does two things with it:

- It hands the policy to next-intl's middleware on a new request, so Next stamps the nonce on its scripts.
- It sets the policy on the response.

The matcher covers every path except `/api`, `/_next` and the root convention files, so a path with a dot under a locale, and the 404, carry the policy and `no-store`.

The config adds, on every path:

- HSTS with includeSubDomains and preload;
- `X-Frame-Options: DENY`;
- the Permissions-Policy;
- `Cross-Origin-Opener-Policy: same-origin` and `Cross-Origin-Resource-Policy: same-origin`;

next to `nosniff`, `Referrer-Policy: same-origin` and `X-Robots-Tag`. API paths add `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`.

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] Every console document carries exactly one CSP:
  - its `script-src` holds this response's nonce and `'strict-dynamic'`, and no `'unsafe-inline'`;
  - Next's scripts carry the nonce;
  - two requests get different nonces.
- [ ] A dotted path under a locale, and an unknown path, get the policy and `Cache-Control: private, no-store`.
- [ ] API responses carry the API policy, and static files run no proxy.
- [ ] Every response carries the config's header set, and no response carries two policies.
- [ ] Every console route renders with no CSP violation. The one exception is zod's caught probe, where the runtime graph renders.
- [ ] The builder's development and preview variants match the spec.

### Verification

- The console's Node tests (the builder); `curl -sI` against the console's production server; a browser pass over every console route.
