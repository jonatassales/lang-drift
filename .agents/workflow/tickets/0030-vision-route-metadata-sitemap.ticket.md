# 0030: Vision page at /platform/vision

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** today's home page content lives at `/[locale]/platform/vision/`, exactly as it is, with its own metadata and sitemap entries, and the links that pointed at the home page's sections point there.

**Blocked by:** none.

**Status:** ready-for-agent

### Route

- [ ] `apps/website/src/app/[locale]/(pages)/platform/vision/page.tsx` is today's home page, with `pageMetadata('vision', …)`, named `VisionPage`, without `WebsiteStructuredData`.

- [ ] The voice preview section (the scripted orb, `#voice`) leaves the Vision page: AI Voice becomes its own product. The `VoicePreview` component stays for that product's page, allow-listed in `platform/quality/dead-code.allow`, and the `.voice-band` rules leave `globals.css` and `responsive.css`.

### Metadata

- [ ] `metadata/pages.json`: a `vision` entry per locale carrying today's home copy under the Vision name; `home` gets platform copy.
- [ ] `metadata/page-metadata.ts` maps `vision` to `platform/vision/`; the JSON stores no route path.
- [ ] `sitemap.xml` gains the five Vision URLs with their alternates.

### Links

- [ ] `website-footer.tsx` and `plans-section.tsx` FAQ links point at `/platform/vision#faq`.

### Verification

- `pnpm turbo run typecheck test --filter=website`
- `/en/platform/vision/` renders today's home; its title and canonical name Vision.
