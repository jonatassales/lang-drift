# Next.js rules

The Next.js apps under `apps/`. Loaded in every session.

1. **App Router under `src/app`.** The App Router's `app/` sits directly under `src/`. No Pages Router tree exists. Framework files keep their framework meaning (`layout.tsx`, `page.tsx`, `template.tsx`, `loading.tsx`, `error.tsx`, `not-found.tsx`, `route.ts`, `proxy.ts`); no custom wrapper recreates one of those jobs.
2. **The proxy.** Next discovers `proxy.ts` in the parent of `app/`, so it is `src/proxy.ts`; a `proxy.ts` beside `package.json` is not loaded. Its implementation stays in `src/proxy/`.
3. **The tree.** Every Next app follows it, and new files do not invent a second structure. `src/app` holds the routes and the UI, i18n, data and state they use; `template/` is shared presentation infrastructure, and not every file in it is JSX. `ports/` and `adapters/` are never created just to match the diagram.

   ```text
   apps/<app>/src/
     proxy.ts              the framework proxy, sibling of app/
     proxy/                its implementation
     app/
       template/           shared App Router presentation infrastructure
         tailwind/         the Tailwind entry: directives, tokens import, approved utilities
         formatters/       one formatter per file, e.g. to-currency.fmt.ts
         i18n/             routing, navigation, request config, locale helpers, messages/
       metadata/           SEO source: site.json, pages.json, locales.json, organization.json, page-metadata.ts
       api/                route handlers
       atoms/              shared state
       state/              the @state barrel and state contracts
       …                   routes, layouts, pages, slots and the UI and data they own
     infra/                concrete third-party SDKs and service integrations
     services/             cookie, crypto and domain service classes
     factories/            factories
     ports/                only when a real interface boundary requires one
     adapters/             only when a real external or framework adaptation requires one
   ```

4. **Localization lives in `src/app/template/i18n/` only.** It owns routing, navigation, request configuration, messages and locale helpers; no second localization directory exists.
5. **Metadata.** `src/app/metadata/` is the application's SEO source, never `src/metadata`. Next file-convention metadata stays directly under `src/app/` (`favicon.ico`, `icon.png`, `apple-icon.png`, `robots.txt`, `sitemap.xml`, `manifest.webmanifest`): framework-owned endpoints and assets, kept as static files and never replaced by TypeScript generators. Social images follow the SEO rules.
6. **`Link` for internal navigation.** An anchor stays where its semantics require it: external destinations, fragment links (`#section`), downloads and deliberate full-document navigation or recovery.
7. **Server-only configuration stays on the server.** Client code receives a value only when it has a real need for it. Being a URL makes a value neither secret nor client-side. No centralized environment wrapper.
8. **App configuration.** A private app manifest carries no package version unless publication becomes an explicit requirement; publishable packages keep their versioning. An app-local Turbo configuration holds only meaningful app-specific overrides. Scripts call the framework directly (`next build`). An app README explains the surface's product role and behavior; stack details support that explanation rather than replace it.
