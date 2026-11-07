# Platform landing hero, Vision page and platform header

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

LangDrift is becoming a platform of products under one domain (AGENTS.md, "The idea"), but the website still presents a single product: its home page is LangDrift Vision's landing page, its header links to that page's sections and its metadata describes Vision. The platform site needs a home that presents the platform, a place where today's Vision landing page lives unchanged, and a header whose content names the platform's products, while the website's design, its glassy header and its hamburger drawer stay exactly as they are.

## Solution

- Today's home page content moves, as it is, to `/[locale]/platform/vision/`, with its own metadata and sitemap entries.
- The home page becomes header, a new platform hero and footer. The hero is in a Vercel-like register (a dark field, a two-word headline, two buttons, the logo mark in the middle) in the brand purple, over a WebGL shader background adapted from Exphere's `exphere-shader`: more purple shades, a light theme, and colors read from the token roles at runtime.
- The header keeps its look and changes its content: logo and locale pill, then `Platform` (a glassy dropdown listing Vision, FinOps, Microruntimes and AI Voice with monochrome icons), Docs, GitHub and `Pricing` together in the middle, then a rainbow-highlighted `Book a Demo` button and the theme toggle. The drawer mirrors it, and the locale pill shows at every width.
- The rainbow is a new fixed color group in `@langdrift/tokens`, like the status, classification and drift colors.

## User Stories

1. As a visitor, I want the home page to present the LangDrift platform, so that I understand it is more than one product.
2. As a visitor, I want today's Vision landing page at `/platform/vision/`, exactly as it was, so that nothing I knew is lost.
3. As a visitor, I want the hero background to move subtly in my color scheme and stay still when I ask for reduced motion, so that it decorates without distracting.
4. As a visitor on a browser without WebGL2, I want a purple gradient in the hero instead of an empty space, so that the page still reads as designed.
5. As a visitor, I want the Platform menu to list the platform's products with a line each and open Vision's page, so that I can find the product I came for.
6. As a visitor on a narrow screen, I want the hamburger drawer to carry the same items, so that nothing depends on screen size.
7. As a visitor, I want the header's Book a Demo button to carry a rainbow highlight, so that the one call to action stands out.
8. As a search engine, I want the home and the Vision page to carry their own titles, descriptions, canonicals, hreflang and sitemap entries, so that both index correctly.
9. As a frontend developer, I want the shader as a `@langdrift/web` component that reads colors only from token roles, so that one package owns every color and the component serves any app.
10. As a frontend developer, I want the rainbow as a fixed token group exposed to Tailwind, so that no palette literal appears in a component.
11. As a frontend developer, I want every new string in JSON catalogs with parity across the five locales, so that the localization checks stay green.
12. As the owner, I want the new hero copy drafted for review, so that I edit the catalogs rather than start from nothing.

## Implementation Decisions

- **Routes.** `[locale]/(pages)/platform/vision/page.tsx` is today's home page file with the metadata id `vision` and without the site-level JSON-LD, which stays on the home page. The home page keeps its skip link, header and footer and renders only the new hero in `<main>`.
- **Metadata.** `metadata/pages.json` gains a `vision` entry per locale and new platform copy for `home`. `page-metadata.ts` maps the id `vision` to the path `platform/vision/`, so the JSON stores no route path. `sitemap.xml` gains five `<url>` blocks.
- **Hero shader.** `packages/web/src/ui/hero-shader/`: raw WebGL2 without a library, a fullscreen triangle and Exphere's petal scene with new uniforms `shadeDeep`, `shadeSoft`, `glow`, `violet` and `scheme`. Colors come from `--background`, `--card`, `--foreground`, `--primary` and `--primary-hover`, read from the canvas's computed style and normalized through a 2D canvas; the derived shades are computed in JS; a `MutationObserver` on `<html class>` recolors the scene when the theme changes. Reduced motion, visibility, the pixel caps, context loss and cleanup stay as in Exphere. The module loads lazily after first paint behind a token-only Tailwind gradient that also serves as the fallback.
- **Hero.** `lib/components/platform-hero.tsx` is a server component; its copy comes from a new `platform` catalog read server-side, so the client bundle gains no catalog. Its centerpiece is a black hole sun built from role utilities. A new monochrome logo mark, a ring with one segment carried outward, replaces the swoosh in the header, footer and integrations section. The headline is "Agentic Platform", the owner's fixed positioning (Vercel's hero says "Agentic Infrastructure"); it is translated in the other locales.
- **Header.** One nav model feeds the desktop nav and the drawer. The Platform dropdown is a disclosure (a button with `aria-expanded` and `aria-controls`) over a Tailwind panel with the header's glass treatment. Items without a page render as text, as do Docs and GitHub. The demo button is a `span` without href; its rainbow ring is a rotating conic gradient of the six rainbow stops, static under reduced motion. The old `.website-get-started` rules leave `website-header.css`, which may only shrink.
- **Tokens.** `shared/rainbow.json` declares `lg.rainbow.1` to `lg.rainbow.6`, the same in both schemes; `build.js` adds `rainbow` to the fixed groups so `--color-rainbow-<n>` reaches Tailwind.
- **Workflow.** One commit per ticket on `website-product-redesign`, each with its `Ticket:` trailer.

## Testing Decisions

- `packages/tokens/tokens.test.js`: the rainbow stops exist in light and are absent from dark.
- `packages/web/src/ui/hero-shader/hero-shader.unit.tsx`: the SSR markup (aria-hidden, absolute inset, pointer events, class merge) and the pure palette math (`mix`, `rotateHue`, `derivePalette`: light shades lighter than dark, `scheme` 0 or 1).
- `apps/website/src/app/template/i18n/catalogs.unit.ts`: the `platform` family joins the parity check.
- `drift check localization`, `drift check dead-code`, `pnpm turbo run check`, and typecheck and tests for `website`, `@langdrift/web` and `@langdrift/tokens`.
- Manual: `/en/` and `/en/platform/vision/` in the dev server; the theme toggle recolors the shader; reduced motion gives a static frame; WebGL disabled shows the gradient; the Platform dropdown opens on hover, focus and click and closes on Escape and outside click; the drawer mirrors the nav; `/en/sitemap.xml` lists 20 URLs.

## Out of Scope

- Pages for FinOps, Microruntimes and AI Voice, and destinations for Docs, GitHub and Book a Demo.
- Sections below the hero on the home page.
- Growing AGENTS.md's product table beyond Vision and FinOps.
- Running the shader behind the sticky header.
