# Vercel Analytics on the console

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The website (`apps/website`) reports visitors and page views to Vercel Analytics since spec 0022, but the console (`apps/console`), deployed on Vercel as `langdrift-console`, reports nothing. The console's document policy is a per-request nonce policy (spec 0018): a script runs only when it carries the response's nonce or when a nonced script loads it (`'strict-dynamic'`), and `connect-src` allows only the console origin, so the analytics endpoints must be allowed explicitly.

## Solution

The console renders `<Analytics />` the way the website does, and its document policy allows the same Vercel Analytics hosts.

## User Stories

1. As the owner, I want console visits and page views in Vercel Analytics, as the website already reports them.
2. As a console user under the nonce policy, I want the browser to load and send analytics without a CSP violation.
3. As a developer, I want the policy test to state the analytics sources, so a later change cannot drop them silently.

## Implementation Decisions

- **Package:** `@vercel/analytics` as a runtime dependency of `apps/console`, at the website's range, `^2.0.1`.
- **Layout:** `apps/console/src/app/[locale]/layout.tsx` imports `Analytics` from `@vercel/analytics/next` and renders `<Analytics />` inside `<body>`, after the `NextIntlClientProvider`, as the website does.
- **Document policy (`apps/console/src/proxy/content-security-policy.ts`), in every environment:**
  - `script-src` adds `https://va.vercel-scripts.com`, and `connect-src` adds `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com`: the website's hosts (spec 0022).
  - The package loads `/_vercel/insights/script.js` from the console origin in production and `https://va.vercel-scripts.com/v1/script.debug.js` in development. A nonced script injects it, so `'strict-dynamic'` lets it run.
- **README:** the console README's nonce policy bullet names the analytics sources.

## Testing Decisions

- **Seam:** `apps/console/src/proxy/content-security-policy.test.ts`, which already asserts the policy per environment, expects the new sources in `script-src` and `connect-src`.
- `pnpm turbo run check typecheck test build --filter=console`, and the console in the browser with no CSP violation reported.

## Out of Scope

- Speed Insights, custom events and any analytics configuration beyond the component's defaults.
