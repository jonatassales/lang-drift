# 0137: Release ex 0.2.0

**Spec:** `.agents/workflow/specs/0024-ex-init-asks-finops-and-generates-ci-gates.spec.md`

**What to build:** `@langdrift/ex` moves to 0.2.0, the minor version that carries the init questions, the artifacts folder, the CI Gates workflow and `ex gates report`. Once the pull request merges, the merge commit on `main` gets the tag `packages/ex/v0.2.0`, Go's tag form for a module in a subfolder, and a GitHub Release with notes. Publishing to npm stays a separate step for the owner.

**Blocked by:** 0133, 0134, 0135, 0136.

**Status:** ready-for-agent

- [x] `packages/ex/package.json` is at 0.2.0, so `build.sh` injects 0.2.0 into `ex version` and into the CI Gates workflow `ex init` writes.
- [ ] The tag `packages/ex/v0.2.0` points at the merge commit on `main`.
- [ ] A GitHub Release on that tag lists what 0.2.0 adds.
