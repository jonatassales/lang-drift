# Next.js rules

Rules for the Next.js apps under `apps/`. Loaded in every session.

1. **App Router under `src/app`.** Framework files keep their framework meaning: `layout.tsx`, `page.tsx`, `template.tsx`, `loading.tsx`, `error.tsx`, `not-found.tsx`, `route.ts` and `proxy.ts`. No custom wrapper recreates one of those jobs, and no Pages Router tree exists.

2. **The tree.** Next discovers `proxy.ts` in the parent of `app/`, so it sits in `src/`, with its implementation beside it.

   ```text
   apps/<app>/src/
     proxy.ts          the framework proxy, sibling of app/
     proxy/            its implementation
     app/
       template/       shared App Router presentation infrastructure
         tailwind/     the Tailwind entry: directives and the tokens import
         formatters/   e.g. cn.fmt.ts
         i18n/         routing, navigation, request config and messages/
       metadata/       SEO source: site.json, pages.json, locales.json, organization.json, page-metadata.ts
       api/            route handlers
       atoms/, state/  shared state (07-jotai)
       …               routes, layouts and the UI and data they own
     infra/            third-party integrations
     services/         cookie and crypto services
     factories/        factories; ports/ and adapters/ only when a real boundary needs them
   ```

   No `utils/`, `helpers/`, `common/` or `misc/`.

3. **Localization lives in `src/app/template/i18n/` only.** It owns routing, navigation, request configuration, messages and locale helpers; no second localization directory appears.

4. **Metadata.** `src/app/metadata/` is the application's SEO source. The convention files (`favicon.ico`, `icon.png`, `apple-icon.png`, `opengraph-image.*`, `twitter-image.*`, `robots.txt`, `sitemap.xml`, `manifest.webmanifest`) sit directly under `src/app/` as static files, never as TypeScript generators.

5. **`Link` for internal navigation.** An anchor stays where its semantics are needed: external destinations, fragment links, downloads and deliberate full-document navigation.

6. **Server-only configuration stays on the server.** Client code receives a value only when it has a real need for it; being a URL makes a value neither secret nor client-side. No centralized environment wrapper (see `04-code.rule.md`).

7. **Aliases live in `tsconfig.paths.json` only**, extended by `tsconfig.json`: no `baseUrl` and no duplicated `paths`. A private app manifest carries no version. Scripts call the framework directly (`next build`), never through a wrapper.
