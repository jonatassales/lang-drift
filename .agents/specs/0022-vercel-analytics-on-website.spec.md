# Vercel Analytics on the website

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The platform website (`apps/website`) deployed at `https://langdrift.md` on Vercel does not currently track web traffic, visitors, or page views through Vercel Analytics. Vercel provides a first-party `@vercel/analytics` package for Next.js App Router applications, which injects an `<Analytics />` component into the root layout. Furthermore, the website enforces a strict Content Security Policy (`Content-Security-Policy`), which must explicitly allow Vercel Analytics script and telemetry endpoints so that tracking requests are not blocked by the browser.

## Solution

1. Add `@vercel/analytics` to `apps/website` dependencies.
2. Render `<Analytics />` from `@vercel/analytics/next` in the root locale layout (`apps/website/src/app/[locale]/layout.tsx`).
3. Update the production Content Security Policy in `apps/website/next.config.ts`:
   - `script-src`: permit `https://va.vercel-scripts.com`
   - `connect-src`: permit `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com`
4. Update `apps/website/next.config.unit.ts` and `apps/website/README.md` to match the updated CSP contract.

## User Stories

1. As a platform team member, I want visitor and page view metrics collected on Vercel Analytics when users browse `https://langdrift.md`.
2. As a visitor browsing the site under a strict Content Security Policy, I want the browser not to report CSP violations for Vercel Analytics script loading or telemetry dispatch.
3. As a developer running unit tests, I want `next.config.unit.ts` and `tsc --noEmit` to verify the new CSP directives and component types.

## Implementation Decisions

- **Package installation:** Install `@vercel/analytics` as a runtime dependency in `apps/website`.
- **Layout placement:** Place `<Analytics />` inside `<body>` in `apps/website/src/app/[locale]/layout.tsx`. In Next.js App Router, `@vercel/analytics/next` automatically tracks route changes and pageviews across client navigation.
- **CSP contract:**
  - `script-src`: add `https://va.vercel-scripts.com`
  - `connect-src`: add `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com`
- **Unit test assertions:** Update `productionPolicy` in `apps/website/next.config.unit.ts` to assert that the CSP includes these endpoints on production routes.

## Testing Decisions

- **Seam 1 — CSP Header Verification (`apps/website/next.config.unit.ts`):** Vitest tests assert that `Content-Security-Policy` contains the updated `script-src` and `connect-src` on production routes (`/en/`, `/api/book-a-demo/`, `/robots.txt`).
- **Seam 2 — Build & Typecheck:** `pnpm --filter website typecheck` and `pnpm --filter website test` pass cleanly.
