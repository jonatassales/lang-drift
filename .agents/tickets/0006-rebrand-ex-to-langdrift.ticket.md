# 0006: Rebrand ex to LangDrift

**Spec:** .agents/specs/0002-ex-hardening.spec.md

**What to build:** a user who runs ex sees LangDrift: npm `@langdrift/ex`, Go module `github.com/jonatassales/tao/packages/ex`, project directory `.langdrift/`, `LANGDRIFT_SKIP_EX_BUILD`, LangDrift in messages and docs. Binary name `ex`, `PROJECT_SECRET`/`RUNTIME_SECRET` unchanged.

**Blocked by:** 0005.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Module path renamed in go.mod and every import; package.json renamed; build.sh env var renamed.
- [ ] `.exphere` appears nowhere in src (code, scaffold assets, tests); `.langdrift` replaces it.
- [ ] No user-facing "Exphere" remains (README, help, errors, scaffold content). Historical notes may name the origin once in the README.
- [ ] Full suite green.
