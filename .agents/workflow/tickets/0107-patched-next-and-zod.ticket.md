# 0107: Both apps run patched Next.js and zod releases

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** The website and the console both run `next` 16.4.3. sinapsi resolves the patched zod 4.6.7 through a workspace override, and the website's tests run on `vitest` 5.0.5. The orbo-voice release-age exclusion, whose stated expiry has passed, is gone. npm reports no high or critical advisory for the production dependencies. This ticket ships alone in the first pull request.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `pnpm audit --prod` reports no high or critical advisory.
- [ ] The lockfile resolves `next` 16.4.3 for both apps, and zod 4.6.7 under sinapsi.
- [ ] No new `minimumReleaseAgeExclude` entry appears, and the orbo-voice entry is gone.
- [ ] Build, check, typecheck and test pass.
- [ ] Both apps render their pages in development.

### Verification

- `pnpm audit --prod`; `pnpm turbo run build check typecheck test`; `pnpm why zod --filter console`; both apps opened in the browser pane.
