# 0217: the CLI is published as `@langdrift/cli` 0.4.0

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The CLI package takes its new npm name and version, `@langdrift/cli` 0.4.0, and every place that tells a user how to install or run it follows: the launcher's reinstall hint, the setup message, the CI gates scaffold and its fixture, the root and CLI READMEs, the translated READMEs, the website footer badge and the dashboard widget, which reads the version from the manifest instead of a literal. `SECURITY.md` lists `@langdrift/cli` and `@langdrift/core` as the published packages and nothing else. The command stays `ex`. The release itself is the owner's.

**Blocked by:** 0216.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `lib/cli/package.json` is `@langdrift/cli` 0.4.0 with bin `ex`; `pnpm pack --dry-run` in `lib/cli` lists the launcher, the four binaries, LICENSE and README under the new name, and `dist/ex version` prints `ex 0.4.0`.
- [x] `ex init` writes a CI gates workflow that runs `@langdrift/cli@0.4.0`, and the CLI's scaffold test still matches its fixture.
- [x] No README, translated README, website badge, widget or CLI message names `@langdrift/ex`; the dashboard widget reads the version from the manifest.
- [x] `SECURITY.md` lists `@langdrift/cli` and `@langdrift/core` as published and no private package.
- [x] `git grep -n '@langdrift/ex\b'` finds only workflow history and CLI fixtures.
- [x] Once the owner has tagged `lib/core/contracts/go/v0.1.0`, `lib/cli/go.mod` requires that version (deferred from 0216: an untagged `require` makes Go fetch the module even inside the workspace).
- [x] Owner, after merge and a fresh `pnpm logout && pnpm login`: `@langdrift/cli@0.4.0` is on npm, `@langdrift/ex` is deprecated with a pointer to it, the tags `lib/cli/v0.4.0` and `lib/core/contracts/go/v0.1.0` and the GitHub Release exist, and `pnpm dlx @langdrift/cli@0.4.0 --version` prints `ex 0.4.0`.

## Dev part

**Override:** No part qualifies: a package rename in manifests, messages and documentation (policy, Categories: never a rename). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go test -race ./... && pnpm pack --dry-run`
- `pnpm turbo run build --filter=website`
- `platform/cli/drift quality check`
