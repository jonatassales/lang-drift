# 0019: The web package

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** `@langdrift/web` exists and exports, at `@pack/react`'s subpaths, exactly the components and utilities the website imports, ported from langdrift. Its colors come from `@langdrift/tokens` roles and groups; it carries no `--ld-*` variable.

**Blocked by:** 0018 (Website token theme in Exphere purple).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The exports map lists only: styles, utilities, product vision curve, brand, agent orb, server icons, language switcher, AI avatars, target distribution chart and the four SmoothUI components.
- [x] Its stylesheet imports the tokens stylesheet; the drift-sign rules read the drift colors.
- [x] No `--ld-*` variable or `ld-*` utility remains in the package.
- [x] Its ported tests pass under `turbo run test`, and `typecheck` passes.
