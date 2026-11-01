# 0163: Add the @langdrift/ex npm badge to the website footer

**What to build:** In the website footer (`WebsiteFooter`), add the `@langdrift/ex` npm badge on the right of the GitHub icon in the social links section, linking to `https://www.npmjs.com/package/@langdrift/ex`. Allow `https://img.shields.io` in the website's `img-src` Content-Security-Policy directive so the badge renders securely.

**Blocked by:** none

**Status:** ready-for-agent

- [x] Add the npm badge link on the right of the GitHub icon in `apps/website/src/app/lib/components/website-footer.tsx`.
- [x] Allow `https://img.shields.io` in `img-src` in `apps/website/next.config.ts`, `apps/website/next.config.unit.ts`, and `apps/website/README.md`.
- [x] Verification with Biome check, test suite, and drift quality checks passes.
