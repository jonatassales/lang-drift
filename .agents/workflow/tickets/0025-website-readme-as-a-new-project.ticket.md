# 0025: The website README as a new project

**Spec:** .agents/workflow/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the website README describes the website as it lives in this repository: its role, behavior, configuration, commands and release gates, with no langdrift history and no links to files this repository doesn't have.

**Blocked by:** 0023 (The website has no Docs zone).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] No mention of the Docs zone, `COOKIE_ENCRYPTION_KEY`, older cookie formats, historical packages, format2 rollback or end-to-end suites elsewhere.
- [ ] Every link resolves inside this repository, and every command is this repository's own (`pnpm turbo run … --filter=website`, `drift check localization`).
