# 0019: The web package

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** `@langdrift/web` exists and exports, at `@pack/react`'s subpaths, exactly the components and utilities the website imports, ported from langdrift. Its colors come from `@langdrift/tokens` roles and groups; it carries no `--ld-*` variable.

**Blocked by:** 0018 (Website token theme in Exphere purple).

**Status:** ready-for-agent

- [ ] The exports map lists only: styles, utilities, product vision curve, brand, agent orb, server icons, language switcher, AI avatars, target distribution chart and the four SmoothUI components.
- [ ] Its stylesheet imports the tokens stylesheet; the drift-sign rules read the drift colors.
- [ ] No `--ld-*` variable or `ld-*` utility remains in the package.
- [ ] Its ported tests pass under `turbo run test`, and `typecheck` passes.
