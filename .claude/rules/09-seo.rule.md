# SEO rules

Metadata on the website. Loaded in every session.

1. **One public origin.** `https://langdrift.md`, from `src/app/metadata/site.json`. Canonicals are never derived from the request host. A localhost social-image URL during local verification is no reason for a root document shell; absolute production social-image origins are verified against `https://langdrift.md` after deployment.
2. **`[locale]/layout.tsx` owns the document**: `<html>`, `<body>`, `metadataBase` and `applicationName`. `/` is locale negotiation in the proxy; it has no page and no document shell of its own.
3. **Indexable routes carry full metadata**: a unique localized title and description, a self canonical, hreflang alternates, a server-rendered `lang`, Open Graph, Twitter metadata and a `sitemap.xml` entry. Child-page metadata belongs to the child page; shared layout metadata is a foundation, not a replacement. Transient pages are `noindex`.
4. **Social images are files.** `opengraph-image.*` and `twitter-image.*`, with their `.alt.txt`, live in the segment whose layout sets `metadataBase`: `src/app/[locale]/`. There Next resolves them against the public origin and serves each once, at `/-/opengraph-image.*` and `/-/twitter-image.*`; at the root they would resolve against the deployment environment's origin. The files own the metadata Next derives from them: width, height, MIME type, alt and absolute image descriptors are never duplicated by hand.
5. **Metadata JSON stores no route paths.** Titles, descriptions, the site name, the public origin, Open Graph locales and organization fields live in `metadata/*.json`; route paths, locale ids and prefixes live in code (`routing.ts`).
