# 0111: Every indexable page describes itself in JSON-LD

**Spec:** .agents/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** Every indexable page, in every locale, renders one JSON-LD graph with a WebPage node:

- its `@id` and `url` come from the canonical;
- its `name` and `description` equal the meta tags;
- its `inLanguage` equals the document's `lang`;
- its `isPartOf` references the WebSite node.

The home page's graph keeps its Organization and WebSite nodes, sameAs included, exactly as they are. The graph is built from the metadata JSON the meta tags already read.

**Blocked by:** 0110.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Every indexable page holds exactly one `application/ld+json` script, and it parses.
- [ ] Its WebPage node carries the expected `@id`, url, name, description, `inLanguage` and `isPartOf`.
- [ ] The home page's Organization and WebSite nodes are unchanged.
- [ ] `<` is escaped in the serialized JSON.
- [ ] Noindex pages carry no WebPage node.

### Verification

- The website unit suite (a builder test); `next build`, then parsing every prerendered page's JSON-LD.
