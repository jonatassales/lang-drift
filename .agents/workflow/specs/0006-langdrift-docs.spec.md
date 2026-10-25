# LangDrift documentation

**Status:** ready-for-agent

## Problem Statement

LangDrift has no documentation. The material that explains FinOps, its CLI and its architecture lives in Exphere's Nextra site (`jonatassales/exphere` `apps/website`, served at exphere.vercel.app/docs), behind a homepage that presents Exphere as its own product. Exphere is now LangDrift FinOps, Vision has no documentation at all, and nothing describes the platform both products share or how to work in this repository.

## Solution

A new app, `apps/docs`, copied from Exphere's website with its homepage removed, so the documentation is the root. It is served under `basePath: '/docs'`, so its links, assets and search results already match `langdrift.md/docs`. The design stays as it is: the Nextra docs theme, dark only, the purple accent, the same navbar, sidebar, table of contents and footer arrangement. Only the identity and the content change: the LangDrift mark and name replace the Exphere wordmark and Orbo, and the content is restructured into Getting Started, Platform, Vision, FinOps and Development.

## User Stories

1. As a reader, I want the documentation at `/docs/` to open on Getting Started, so that I land on the documentation, not a marketing page.
2. As a reader, I want Getting Started to explain that LangDrift is a platform of products on one domain and to point me to each product, so that I know where to go.
3. As a reader, I want a Platform section that explains the one domain, the console, the product switcher and the principles the products share, so that I understand what holds the products together.
4. As a reader, I want Platform → Architecture to show the repository's apps and packages and what each owns, so that I know where things live.
5. As a reader, I want a Vision page that explains the Product Vision score, its classifications and what exists today, so that I understand Vision without reading marketing copy.
6. As a reader, I want a FinOps section with an overview, concepts, architecture and the `ex` CLI, so that the Exphere material is still there under its new name.
7. As a contributor, I want a Development page that describes this repository's toolchain, commands, hooks and CI, so that I can work here.
8. As a reader, I want every page to say what is proposed, planned or unverified, so that I never mistake a design for a shipped feature.
9. As a reader, I want search to find a page and open it at the right URL, so that search is useful.
10. As a reader, I want the sidebar, table of contents, breadcrumbs, pagination and mobile menu to work as they did on Exphere's docs, so that the reading experience does not regress.
11. As the owner, I want the docs to look like Exphere's docs, so that the port does not become a redesign.
12. As a frontend developer, I want the docs' colors to come from `@langdrift/tokens` roles, so that the app declares no theme value of its own.

## Implementation Decisions

- **App.** `apps/docs`, Next.js App Router with route-local `page.mdx` under `src/app`, Nextra and `nextra-theme-docs` 4.6.1 on Exphere's verified pins (next 16.4.0, react 19.3.0, TypeScript 6.0.3 for Twoslash, pagefind 1.5.2). `basePath: '/docs'`; `pnpm docs` starts it on port 3001. Pointing the website at it is a later deploy decision; spec 0003 removed the old Docs zone on purpose.
- **Homepage removed.** The shader page, the copy-command preview, the Orbo preview, `@repo/web`, `orbo-voice` and the homepage-only metadata go. The root layout and the docs layout become one root layout that owns `<html>` and Nextra's `Layout`.
- **Identity.** `Brand` and `LogoMark` from `@langdrift/web/ui/brand` replace the Exphere wordmark and Orbo in the navbar and footer; the footer reads "© 2026 LangDrift". The website's LangDrift icons replace Exphere's.
- **Tokens.** `@langdrift/tokens` gains a `docs` app theme with the same purple as the website, and the docs import `@langdrift/tokens/docs.css`. The Exphere `site-*` palette is replaced by roles: `border-border`, `text-muted-foreground`, `font-sans`. `.dark` replaces `data-theme`.
- **Structure.** `_meta.ts` orders Getting Started, Platform (with Architecture), Vision, FinOps (Overview, Concepts, Architecture, CLI) and Development. Internal links omit `/docs`, which Next adds. The two FinOps diagrams sit beside their page.
- **Content.** FinOps pages keep Exphere's wording with "Exphere" replaced by "LangDrift FinOps"; the CLI page follows `packages/ex/README.md`, since Exphere's commands and paths are out of date. Platform, Vision and Development are written from AGENTS.md, `.agents/context.md` and the website's copy; nothing is invented. Links into Exphere's repository become links into this one or are dropped.
- **Search.** The Pagefind adapter indexes every prerendered page except the framework's error pages, under URLs without the basePath, because Nextra's search adds it when it loads Pagefind and when it navigates. The index is written to `public/_pagefind`, ignored by Git and kept in the Turbo build outputs. Vercel's adapter chaining goes; tao is not on Vercel.
- **Metadata.** `noindex` stays until the owner decides to index the docs. No `robots.txt`, which means nothing under a basePath.
- **Workspace.** `pnpm-workspace.yaml` carries Exphere's scoped Zod override for Nextra. The root gains the `docs` launcher; AGENTS.md and `.agents/context.md` describe the second app.
- **Known deviation.** The docs are English only, so the few chrome strings in the layout ("Search docs…", "On this page") stay in TSX, as in Exphere. Localizing the docs across the five locales is a later owner decision.

## Testing Decisions

A good test checks what a reader or a build sees. The docs carry no unit tests of their own: their behavior is Nextra's, and the build already fails when a page, a link target or the index breaks.

1. **Build:** `pnpm turbo run build --filter=docs...` compiles every page and indexes nine pages; an indexing error fails the build.
2. **Tokens:** the tokens test covers the new theme.
3. **Repository checks:** `pnpm turbo run check`, `typecheck` and `test`, `drift check dead-code` and `drift check localization`.
4. **Browser:** in the built docs, every page loads by direct URL; the sidebar, table of contents, breadcrumbs, pagination and mobile menu work; both diagrams render; a search for "microruntime" opens `/docs/finops/concepts#…`; a missing page answers 404; the console shows no errors; the layout matches exphere.vercel.app/docs.

## Out of Scope

- Serving the docs from the website's domain (rewrites or a zone).
- Localizing the docs, indexing them, a sitemap and social images.
- Porting langdrift's earlier SDK, MCP and schema documentation.
- Any design change beyond identity.
