# 0020: The website app

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** `apps/website` serves langdrift's website, with the same routes, locales, content, proxy behavior, early-access flow and Docs rewrites, in Exphere purple. Only its token layer and configuration differ from the original.

**Blocked by:** 0019 (The web package).

**Status:** ready-for-agent

- [x] The website's tree is copied from langdrift `e518dac` and reformatted to this repository's Biome style.
- [x] `@pack/react` imports resolve to `@langdrift/web` through a pnpm workspace alias; the source's imports are unchanged.
- [x] The website's own theme files are deleted; its colors come only from token roles and groups, and `.dark` replaces `data-theme`.
- [x] The test fixture alias points inside the app; no path climbs with `../`.
- [x] Biome groups the website's aliases as aliases; `turbo run check` passes.
- [x] Its 15 unit tests, `typecheck` and `build` pass.
- [x] Side by side with langdrift's website, pages match in both color schemes apart from the brand color and the four accepted light-mode deltas.

_Superseded: ticket 0026 replaced the `@pack/react` workspace alias with direct `@langdrift/web` imports._
