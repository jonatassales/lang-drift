# 0165: Sync ex package readme with root readme and release 0.2.2

**What to build:** `@langdrift/ex` moves to 0.2.2 with `packages/ex/README.md` matching root `README.md` exactly, updating image URLs to raw GitHub URLs so the dark/light hero banner and terminal screenshot resolve on npmjs.com.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/ex/README.md` is synced to match root `README.md` with raw GitHub image URLs.
- [x] `packages/ex/package.json` is at 0.2.2, so `build.sh` injects 0.2.2 into `ex version`.
- [x] `@langdrift/ex` 0.2.2 is published to npm.
