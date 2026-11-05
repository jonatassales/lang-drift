# 0130: Install @vercel/analytics and add Analytics component to website layout

**Spec:** .agents/workflow/specs/0022-vercel-analytics-on-website.spec.md

**What to build:** Add `@vercel/analytics` to `apps/website` dependencies and render `<Analytics />` from `@vercel/analytics/next` inside `<body>` in `apps/website/src/app/[locale]/layout.tsx`. Update the Content Security Policy in `apps/website/next.config.ts` so `script-src` allows `https://va.vercel-scripts.com` and `connect-src` allows `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com`. Update `apps/website/next.config.unit.ts` and `apps/website/README.md` to reflect the updated CSP contract.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `apps/website/package.json` adds `@vercel/analytics` as a runtime dependency.
- [x] `apps/website/src/app/[locale]/layout.tsx` imports `Analytics` from `@vercel/analytics/next` and renders `<Analytics />` within `<body>`.
- [x] `apps/website/next.config.ts` updates `contentSecurityPolicy()` with `https://va.vercel-scripts.com` in `script-src` and `connect-src`, and `https://vitals.vercel-insights.com` in `connect-src`.
- [x] `apps/website/next.config.unit.ts` updates `productionPolicy` to match the new CSP contract.
- [x] `apps/website/README.md` documents the updated `Content-Security-Policy`.

### Verification

- `pnpm --filter website test`
- `pnpm --filter website typecheck`
- `pnpm turbo run check`
- `platform/cli/drift check dead-code`
