# 0022: Frontend harness

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the guardrails langdrift wrote for its website apply here: rules for code, React, Next.js, Jotai, localization and SEO; Tailwind conventions and a scoped CSS exception in the frontend and Tailwind rules; langdrift's Biome lint rules; and `drift check localization` in CI.

**Blocked by:** 0020 (The website app).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The new rules exist with this repository's paths and are listed in AGENTS.md.
- [x] The frontend and Tailwind rules carry `cn()`, role utilities only, no component classes and the scoped exception that names the existing CSS files.
- [x] Biome reports type aliases where an interface fits, non-kebab file names and unused imports as errors; `turbo run check` passes.
- [x] `drift check localization` passes on the website's catalogs, fails on a broken catalog, runs in CI's Harness job and appears in `drift --help`.
