# 0164: Release ex 0.2.1

**What to build:** `@langdrift/ex` moves to 0.2.1, the patch version that includes the root banner image in `packages/ex/README.md` and the npm badge in the website footer.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/ex/package.json` is at 0.2.1, so `build.sh` injects 0.2.1 into `ex version`.
- [x] `@langdrift/ex` 0.2.1 is published to npm with the banner image in README.
