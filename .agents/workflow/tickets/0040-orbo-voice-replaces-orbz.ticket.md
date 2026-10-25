# 0040: orbo-voice replaces the unpublished orbz package

**Spec:** None

**What to build:** `@langdrift/web` depends on `orbo-voice` instead of `@neongate-ai/orbz`, which no longer exists on the registry, so a fresh install works again and the agent orb renders as before.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `packages/web` depends on `orbo-voice` 1.1.2, pinned exactly like its predecessor, and the lockfile differs from main only by that dependency.
- [ ] `AgentOrb` imports the package's types and browser runtime and renders the `orb-o` element with the same preset, state, size, speed, paused and reduced-motion attributes; `AgentOrbState` keeps its name.
- [ ] The orb and voice-preview tests, including the old website's, query and mock the new names.
- [ ] No reference to `@neongate-ai/orbz` or `orb-z` remains outside the lockfile history.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check dead-code`, `drift check localization`.
- `pnpm install --frozen-lockfile` succeeds once the release is older than the repo's `minimumReleaseAge` (the repo's policy rejects `orbo-voice@1.1.2` until 2026-10-26 17:10 UTC).
- The built pages contain `<orb-o` and no `<orb-z`.
