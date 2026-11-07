# 0128: @langdrift/ex packages prebuilt binaries and runs through npx

**Spec:** .agents/specs/0021-ex-cli-published-on-npm.spec.md

**What to build:** `@langdrift/ex` becomes a publishable public npm package (`0.1.0`, MIT license) with a POSIX `sh` launcher at `bin/ex` and prebuilt static binaries for macOS and Linux (`darwin`/`linux` on `arm64`/`x64`). Running `pnpm --filter @langdrift/ex build` compiles the host binary `dist/ex` with `main.version` injected from `package.json`; running `pnpm pack` or `pnpm publish` invokes `prepack` (`sh build.sh --release`), which cross-compiles all four static `CGO_ENABLED=0` binaries (`dist/ex-darwin-arm64`, `dist/ex-darwin-amd64`, `dist/ex-linux-amd64`, `dist/ex-linux-arm64`). The launcher resolves symlinks from `node_modules/.bin/ex` or the `npx` cache, selects the matching platform binary (or host `dist/ex` in local workspace builds), and `exec`s it with all arguments, stdio, and exit status preserved, exiting `4` on unsupported platforms. `@langdrift/contracts` moves from `dependencies` to `devDependencies` so consumers install zero runtime dependencies.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/ex/package.json` is version `0.1.0`, not private, with `"publishConfig": { "access": "public" }`, `"license": "MIT"`, `"bin": { "ex": "./bin/ex" }`, `"os": ["darwin", "linux"]`, `"cpu": ["x64", "arm64"]`, `"homepage"`, `"repository"`, `"prepack": "sh build.sh --release"`, `"@langdrift/contracts": "workspace:*"` in `"devDependencies"`, and `"files"` limited to `bin/ex`, the four `dist/ex-*` binaries, `LICENSE`, and `README.md`.
- [x] `LICENSE` and `packages/ex/LICENSE` carry the MIT license (`Copyright (c) 2026 LangDrift`).
- [x] `packages/ex/build.sh` injects `package.json` version into `main.version` via `-ldflags="-s -w -X main.version=..."` with `CGO_ENABLED=0`, builds host `dist/ex` by default, and builds all four `dist/ex-<os>-<arch>` binaries when passed `--release`.
- [x] `packages/ex/bin/ex` resolves symlinks, maps `uname -s` and `uname -m` to `darwin`/`linux` and `arm64`/`amd64`, `exec`s the matching binary, falls back to `dist/ex` in a host build, and exits `4` with a diagnostic on stderr when the platform is unsupported or no binary exists.
- [x] Automated tests in `packages/ex/src` verify `bin/ex` symlink resolution, platform dispatch, exit code `4` on unsupported platforms, and version injection.
- [x] `pnpm pack` in `packages/ex` produces a tarball containing only the declared files, and running it via `npx` in a clean directory outside the workspace executes `ex version`, `ex --help`, and `ex init`.

### Verification

- `pnpm --filter @langdrift/ex test`, `pnpm --filter @langdrift/ex typecheck`, `pnpm --filter @langdrift/ex build`, `pnpm turbo run check`, `platform/cli/drift check dead-code`.
- Clean-room `pnpm pack` + `npx` execution in `/tmp`.
