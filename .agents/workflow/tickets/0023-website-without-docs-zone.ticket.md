# 0023: The website has no Docs zone

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the website builds and runs without `DOCS_ORIGIN`. It no longer forwards or redirects Docs paths, a Docs path gets the localized 404 like any missing page, and nothing that existed for the Docs zone or for the next.config test remains.

**Blocked by:** none.

**Status:** ready-for-agent

- [ ] `next.config.ts` has no Docs redirect or rewrites; `next build` passes with no `DOCS_ORIGIN`.
- [ ] `next.config.unit.ts` and its vitest `include` entry are gone.
- [ ] The proxy matcher, the locale proxy and their tests know nothing of Docs paths.
- [ ] The unused `docs` header copy and `.website-docs-link` rules are gone, and the sitemap lists only the website's 15 URLs.
- [ ] `DOCS_ORIGIN`, `NEXT_PUBLIC_DOCS_URL` and `LANGDRIFT_LOCAL_BUILD` appear nowhere: not in `.env.development`, `turbo.json`, CI or AGENTS.md. `turbo.json` also drops the unused `VERCEL`, `VERCEL_URL` and `VERCEL_PROJECT_PRODUCTION_URL`.
