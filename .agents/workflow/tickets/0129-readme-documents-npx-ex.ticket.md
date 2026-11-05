# 0129: Package and repository READMEs document npx @langdrift/ex

**Spec:** .agents/workflow/specs/0021-ex-cli-published-on-npm.spec.md

**What to build:** the `ex` package README, `.agents/context.md`, the root English `README.md`, and its four translations (`pt-BR`, `zh-Hant`, `ja`, `vi`) describe `@langdrift/ex` as a published npm package with prebuilt macOS and Linux binaries that runs directly via `npx @langdrift/ex` (or `npx ex` when installed in a project), alongside the existing instructions for building from source.

**Blocked by:** 0128 (`@langdrift/ex` packages prebuilt binaries and runs through `npx`).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/ex/README.md` documents running `@langdrift/ex` via `npx @langdrift/ex`, the `bin/ex` launcher and prebuilt `darwin`/`linux` (`arm64`/`amd64`) binaries, and `sh build.sh --release` alongside local source builds.
- [x] `README.md` and `platform/docs/README.{pt-br,zh-hant,ja,vi}.md` update Project status and Get started with `ex` so readers can run `npx @langdrift/ex` directly or build from source.
- [x] `.agents/context.md` reflects that `@langdrift/ex` ships a `"bin"` launcher for `npx @langdrift/ex`.

### Verification

- `pnpm turbo run check`, `platform/cli/drift check localization`, `platform/cli/drift check dead-code`.
