# 0142: Install @vercel/analytics and add the Analytics component to the console layout

**Spec:** .agents/specs/0025-vercel-analytics-on-console.spec.md

**What to build:** Add `@vercel/analytics` to `apps/console` and render `<Analytics />` from `@vercel/analytics/next` inside `<body>` in `apps/console/src/app/[locale]/layout.tsx`. The console's document policy allows `https://va.vercel-scripts.com` in `script-src`, and `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com` in `connect-src`, in every environment; its test and the README say so.

**Blocked by:** None (can start immediately).

**Status:** completed

- [x] `apps/console/package.json` adds `@vercel/analytics` as a runtime dependency.
- [x] `apps/console/src/app/[locale]/layout.tsx` renders `<Analytics />` within `<body>`.
- [x] `documentPolicy` adds the analytics sources to `script-src` and `connect-src`.
- [x] `content-security-policy.test.ts` asserts the new sources.
- [x] `apps/console/README.md` names the analytics sources in the nonce policy.
