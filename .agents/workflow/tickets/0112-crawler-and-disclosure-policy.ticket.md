# 0112: robots.txt and security.txt state the crawler and disclosure policy

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** robots.txt keeps search and AI answer crawlers allowed (everything but `/api/`) and adds one group that disallows everything for these training crawlers:

- GPTBot
- ClaudeBot
- Google-Extended
- CCBot
- Applebot-Extended
- Bytespider

The sitemap line stays. `/.well-known/security.txt` follows RFC 9116, with these fields:

- `Contact: mailto:security@langdrift.md`
- `Expires`, one year after the commit
- `Preferred-Languages: en, pt`
- `Canonical` on langdrift.md

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] robots.txt parses into the two groups and the sitemap line.
- [ ] `/.well-known/security.txt` answers 200 as `text/plain`, with the four fields.
- [ ] The SEO contract test pins both files.

### Verification

- The website unit suite; `curl` on both files against `next start`.
