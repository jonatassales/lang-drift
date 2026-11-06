# 0119: The contracts are documented

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** Each README documents the contracts this spec added:

- **Website README:**
  - the security header set and its environment variants;
  - the Turnstile origins to add when sign-up reopens;
  - the crawler policy;
  - the SEO contract and its test;
  - security.txt and its renewal date;
  - the owner's post-deploy checklist.
- **Console README:** the nonce policy, the header set, robots.txt and the entity titles.
- **extub README:** the response headers and the authentication details.

The checklist:

- confirm the production domain and the `og:image` origin;
- redirect the `vercel.app` alias;
- keep Deployment Protection on;
- submit HSTS preload;
- renew security.txt;
- keep the security mailbox monitored.

**Blocked by:** 0108, 0109, 0110, 0111, 0112, 0113, 0114, 0115, 0116, 0117, 0118, 0121.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Each README describes what its app now does, and nothing it does not.
- [ ] The checklist lists every step that code cannot verify.
- [ ] `pnpm turbo run check` and `drift check dead-code` pass.

### Verification

- `pnpm turbo run check`; `drift check dead-code`; the README claims compared with `curl` output.
