# SEO rules

Rules for metadata on the website. Loaded in every session.

1. **One public origin.** `https://langdrift.md`, from `metadata/site.json`. Canonicals are never derived from the request host, and a localhost social-image URL during local verification is no reason for a root document shell.

2. **`[locale]/layout.tsx` owns the document**: `<html>`, `<body>`, `metadataBase` and `applicationName`. `/` is locale negotiation in the proxy; it has no page and no document shell of its own.

3. **Indexable routes carry full metadata**: a unique localized title and description, a self canonical, hreflang, a server-rendered `lang`, Open Graph, Twitter metadata and a `sitemap.xml` entry. Child-page metadata belongs to the child page; shared layout metadata is a foundation, not a replacement. Transient pages are `noindex`.

4. **Social images are files.** `opengraph-image.*` and `twitter-image.*` own the image files and the metadata Next derives from them; width, height, MIME type and alt are never duplicated by hand.

5. **Metadata JSON stores no route paths.** Titles, descriptions, the site name, the public origin, Open Graph locales and organization fields live in `metadata/*.json`; locale ids and prefixes live in `routing.ts`.
