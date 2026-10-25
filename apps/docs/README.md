# Docs

The documentation for LangDrift and every platform product: Getting Started, Platform, Vision, FinOps and Development. It was ported from Exphere's Nextra site under [spec 0006](/.agents/workflow/specs/0006-langdrift-docs.spec.md) with its homepage removed, so the documentation is the root. The design is Exphere's: Nextra's docs theme, dark only, the purple accent.

Nextra owns navigation, breadcrumbs, the article wrapper, the table of contents, pagination, the mobile menu and Pagefind search. Pages are route-local `page.mdx` files under `src/app/`, ordered by `_meta.ts`.

## Local development

Run commands from the repository root:

| Command | Purpose |
| --- | --- |
| `pnpm docs` | Build shared dependencies and start development at `http://localhost:3001/docs`. |
| `pnpm turbo run build --filter=docs...` | Build shared dependencies, compile the docs and generate search assets. |
| `pnpm --filter docs start` | Serve the completed production build on port 3001. |
| `pnpm --filter docs typecheck` | Generate Next.js route types and check TypeScript. |

No secrets, backing service or environment variables are required.

## Base path

The app is served under `basePath: '/docs'`, so its links, assets and search results already match `langdrift.md/docs`. Internal MDX links omit `/docs`; Next adds it. Serving the docs from the website's domain is a separate deploy decision; the website has no Docs zone.

## Search lifecycle

`next build` runs the app's own adapter, `next.adapter.mts`, which uses Next's build-complete output records to index every prerendered page except the framework's error pages. Pages are indexed under their URL without the base path, because Nextra's search adds it both when it loads Pagefind and when it opens a result. Indexing errors fail the build.

Generated `public/_pagefind/` files are ignored by Git and kept in the app's Turbo build outputs. The development server needs a previously generated index for search; verification checks result URLs and served assets, not just the input field.

## Presentation and metadata

Colors come from `@langdrift/tokens/docs.css`, whose app theme holds the same purple as the website; the app declares no theme value of its own and renders with the `.dark` class. The navbar and footer use the LangDrift mark from `@langdrift/web/ui/brand`.

The root layout owns `lang="en-US"`, the application name and `noindex`; pages own their titles and descriptions through front matter. Indexing, a sitemap and social images need the owner's decision. The docs are English only.

## Dependency compatibility

The app keeps Exphere's verified pins: Next.js 16.4.0, React 19.3.0, Nextra 4.6.1 and TypeScript 6.0.3, which Nextra's Twoslash dependency requires. The workspace overrides Nextra's Zod to 4.1.12, because newer Zod rejects the theme's omitted children field. `next.config.ts` maps Nextra's MDX import source to `src/mdx-components.tsx` for both bundlers; the build uses Turbopack.
