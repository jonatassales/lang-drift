# 0114: Console server reads trust only what they should

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:**

- Server-rendered reads reach extub with a request built from the incoming host alone. They carry no browser `Origin`, `Sec-Fetch-*` or cookie, so a page reached from another site, or rendered over HTTPS on the hosted preview, reads normally. Browser calls to the route handlers keep the cross-site rule.
- One entity-id predicate guards every id that reaches an extub path:
  - A route id outside the alphabet answers 404 without calling extub.
  - Client calls check ids the same way.
  - Nothing is encoded twice.
- The extub bridge imports `server-only`.
- The login and sign-up form declares `method="post"`.

**Blocked by:** 0113.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A navigation that carries `Sec-Fetch-Site: cross-site`, and a render with an HTTPS `Origin`, both read extub normally.
- [ ] A cross-site browser call to `/api/console/*` still answers 403.
- [ ] A route id such as `..`, or one with `%` or `/`, answers 404 with no extub request.
- [ ] Seed ids and generated ids still work everywhere.
- [ ] Importing the extub bridge from a client component fails the build.
- [ ] The login form's markup carries `method="post"`, so a pre-hydration submit carries no credential in the URL.

### Verification

- The console's Node tests (the id predicate and the cross-site rule); `curl` with fetch-metadata headers against the production server; `next build`.
