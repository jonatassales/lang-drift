# 0020: The website app

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** `apps/website` serves langdrift's website, with the same routes, locales, content, proxy behavior, early-access flow and Docs rewrites, in Exphere purple. Only its token layer and configuration differ from the original.

**Blocked by:** 0019 (The web package).

**Status:** ready-for-agent

- [ ] The website's tree is copied from langdrift `e518dac` and reformatted to this repository's Biome style.
- [ ] `@pack/react` imports resolve to `@langdrift/web` through a pnpm workspace alias; the source's imports are unchanged.
- [ ] The website's own theme files are deleted; its colors come only from token roles and groups, and `.dark` replaces `data-theme`.
- [ ] The test fixture alias points inside the app; no path climbs with `../`.
- [ ] Biome groups the website's aliases as aliases; `turbo run check` passes.
- [ ] Its 15 unit tests, `typecheck` and `build` pass.
- [ ] Side by side with langdrift's website, pages match in both color schemes apart from the brand color and the four accepted light-mode deltas.
