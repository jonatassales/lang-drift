# 0021: Website build, CI and container

**Spec:** .agents/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** Turborepo, CI and the website's Dockerfile know the website: its build environment is declared, `dev` builds the tokens first, CI builds it with safe values, and the container builds from this repository's paths with pinned inputs.

**Blocked by:** 0020 (The website app).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The website's build declares its environment variables; `dev` depends on upstream builds.
- [x] The pull-request workflow sets the safe values (local build, early access paused, sign-up disabled, a local Docs origin).
- [x] The Dockerfile uses this repository's paths, a digest-pinned Node image and a checksum-pinned pnpm, and no Corepack.

_Superseded: ticket 0024 removed the Dockerfile, and ticket 0023 removed the local-build flag and the Docs origin from the safe values._
