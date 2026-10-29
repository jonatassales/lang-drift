# 0110: Hreflang and social images agree with the sitemap

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:**

- Every indexable page's x-default is its own English URL, as the sitemap already says.
- next-intl stops sending its `Link` header alternates (`alternateLinks: false`, as the console already does), so the HTML and the sitemap are the only hreflang sources.
- The social images resolve on langdrift.md in the build output.
- The localized 404 has a title.
- An SEO contract test pins how the metadata JSON, the metadata helper and the sitemap fit together. It covers every page in the metadata JSON in every locale, and checks:
  - a unique title and description within the locale;
  - a self canonical on langdrift.md with a trailing slash;
  - hreflang for all five locales;
  - Open Graph and Twitter fields;
  - the index directive;
  - sitemap parity.

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] In every prerendered page, the HTML x-default equals that page's sitemap x-default.
- [ ] No website response carries a `Link` header with hreflang alternates.
- [ ] `og:image` and `twitter:image` in every prerendered page are `https://langdrift.md` URLs. If Next cannot resolve them there, the README's post-deploy checklist records the origin check, and the ticket says why.
- [ ] `/<locale>/missing/` answers 404 with a localized title and noindex.
- [ ] The contract test fails when a page enters the metadata JSON without a sitemap entry, or the sitemap lists a URL with no page.
- [ ] The early-access confirmation page stays noindex and out of the sitemap.

### Verification

- The website unit suite; `next build`, then reading the head of every prerendered page; `curl -sI` on a missing path.
