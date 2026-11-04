# 0197: release `@langdrift/ex` 0.3.0

**Spec:** .agents/workflow/specs/0032-components-by-dockerfile-path.spec.md

**What to build:** The owner publishes the release that carries the version 2 fence (0194–0196): bump the package version to 0.3.0, publish to npm, tag `packages/ex/v0.3.0` and write the GitHub Release, as the earlier releases were done. Nothing in this repository can plan or run its microruntime with the pinned CLI until this lands.

**Blocked by:** 0196.

**Status:** ready-for-human

- [ ] `@langdrift/ex@0.3.0` is on npm with the four binaries and reports its version.
- [ ] The tag `packages/ex/v0.3.0` and its GitHub Release exist.
- [ ] `pnpm dlx @langdrift/ex@0.3.0 --version` prints `0.3.0`.

### Verification

- `pnpm dlx @langdrift/ex@0.3.0 --version`
