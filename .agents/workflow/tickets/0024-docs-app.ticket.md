# 0024: The docs app

**Spec:** .agents/workflow/specs/0006-langdrift-docs.spec.md

**What to build:** `apps/docs` serves Exphere's Nextra documentation without its homepage, under `/docs`, in LangDrift's identity, with its content restructured into Getting Started, Platform, Vision, FinOps and Development. The design is unchanged.

**Blocked by:** none.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `@langdrift/tokens` builds `docs.css` from a `docs` app theme, and its test covers it.
- [x] `apps/docs` is copied from Exphere `bf0133e`, without the homepage, reformatted to this repository's Biome style, with colors only from token roles and `.dark` instead of `data-theme`.
- [x] The navbar and footer show the LangDrift mark and name; the icons are LangDrift's.
- [x] Nine pages, in the order of the spec, with no Exphere-only facts left and Exphere named only as FinOps's origin.
- [x] The build indexes nine pages, and a search result opens its page under `/docs/`, not `/docs/docs/`.
- [x] `pnpm docs` starts the app; AGENTS.md, `.agents/context.md` and `apps/docs/README.md` describe it.
- [x] `check`, `typecheck`, `test`, `build`, `drift check dead-code` and `drift check localization` pass.
- [x] Side by side with exphere.vercel.app/docs, the layout matches apart from identity and content.
