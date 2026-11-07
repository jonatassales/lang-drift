# 0026: The website imports @langdrift/web

**Spec:** .agents/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the website depends on and imports `@langdrift/web` by its own name, at the subpaths it used for `@pack/react`. No `@pack/*` name remains in the repository's code or configuration.

**Blocked by:** none.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `apps/website/package.json` depends on `@langdrift/web` (`workspace:*`); the `@pack/react` alias is gone.
- [x] Every import, the layout's `styles.css` import and the lazy chart import name `@langdrift/web`, and `transpilePackages` lists `@langdrift/web`.
- [x] The frontend rule, ADR 0003 and spec 0003 describe the direct dependency.
