# 0109: Book a demo refuses cross-site and non-JSON requests

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** The book-a-demo route applies the early-access request contract before it reads the body:

- `Sec-Fetch-Site: cross-site`, or an `Origin` other than the request's own, answers 403;
- any `Content-Encoding`, or a content type other than JSON, answers 415;
- a declared `Content-Length` over 4 KiB answers 413.

A same-origin JSON request behaves exactly as it does today.

**Blocked by:** 0107.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Each refused request answers its status without the body being read.
- [ ] A same-origin JSON request, with the site's own `Origin` or none, still answers 400 for field errors and 503 when valid.
- [ ] Every response keeps `no-store` and `nosniff`.

### Verification

- The website unit suite (the book-a-demo route test); `curl` against `next start` with a foreign `Origin`.
