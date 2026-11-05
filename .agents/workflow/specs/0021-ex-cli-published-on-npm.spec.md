# The ex CLI is published on npm and runs through npx

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The `ex` CLI lives in `packages/ex` as a private workspace package (`"private": true`, `"version": "0.0.0"`) with no `"bin"` entry, no license file, and a runtime `"dependencies"` entry on the private `@langdrift/contracts` workspace package that only exists so Turborepo orders task builds. A team cannot run `npx @langdrift/ex` or install `@langdrift/ex` from npm: the package is unpublished, declares no executable entry point, would fail to resolve `@langdrift/contracts` from the public registry, and would require a local Go 1.26+ toolchain and the monorepo's `go.work` if it tried to compile from source on the consumer's machine. Meanwhile, the package and repository READMEs tell readers that no binary is published and that `ex` must be built from a clone of the repository.

## Solution

`@langdrift/ex` becomes a public, MIT-licensed npm package at version `0.1.0` under the `@langdrift` organization that ships prebuilt static binaries for macOS and Linux (`darwin`/`linux` on `arm64`/`x64`) alongside a POSIX `sh` launcher at `bin/ex`. Running `npx @langdrift/ex <command>` (or `pnpm dlx @langdrift/ex <command>`) anywhere, or `npx ex <command>` inside a project where `@langdrift/ex` is installed, resolves the host operating system and CPU architecture and `exec`s the matching compiled Go binary with its arguments, standard streams, and exit code intact—without requiring Go on the consumer's machine or running a `postinstall` script.

Local development stays fast: `pnpm --filter @langdrift/ex build` (and `sh build.sh` with no flags) still compiles only the host binary `dist/ex` with the version from `package.json` injected into `main.version`. Packing or publishing runs `prepack` (`sh build.sh --release`), which cross-compiles all four static `CGO_ENABLED=0` binaries (`dist/ex-darwin-arm64`, `dist/ex-darwin-amd64`, `dist/ex-linux-amd64`, `dist/ex-linux-arm64`) with stripped symbols and the package version injected. The published tarball contains only `bin/ex`, the four platform binaries, `LICENSE`, `README.md`, and `package.json`. `@langdrift/contracts` moves to `devDependencies` so Turborepo keeps its task ordering while consumers install zero external dependencies.

## User Stories

1. As a developer evaluating LangDrift, I want to run `npx @langdrift/ex version` without cloning the repository or installing Go, so that I can verify the CLI immediately.
2. As a developer evaluating LangDrift, I want `npx @langdrift/ex --help` and `npx @langdrift/ex init` to run directly in my project folder, so that I can initialize `.langdrift/` and plan a ticket with one command.
3. As a developer with `@langdrift/ex` installed in my project's `node_modules`, I want `npx ex` and `pnpm exec ex` to invoke the `ex` binary rather than the system Vi line editor, so that project scripts can call `ex` directly.
4. As a developer on macOS Apple Silicon (`darwin/arm64`), I want `npx @langdrift/ex` to execute the `darwin-arm64` binary, so that the CLI runs natively on my machine.
5. As a developer on macOS Intel (`darwin/amd64`), I want `npx @langdrift/ex` to execute the `darwin-amd64` binary, so that the CLI runs natively on my machine.
6. As a developer on Linux x86_64 (`linux/amd64`), I want `npx @langdrift/ex` to execute the `linux-amd64` binary, so that the CLI runs natively in Linux workstations and CI containers.
7. As a developer on Linux ARM64 (`linux/arm64`), I want `npx @langdrift/ex` to execute the `linux-arm64` binary, so that the CLI runs natively on ARM64 Linux hosts and containers.
8. As a developer on an unsupported operating system or CPU architecture, I want `bin/ex` to print a clear diagnostic to stderr and exit with code 4 (missing prerequisite), so that the failure obeys the CLI's exit code contract.
9. As a security-conscious consumer with `ignore-scripts=true` or pnpm `allowBuilds` restrictions, I want `@langdrift/ex` to work without any `install` or `postinstall` lifecycle script, so that installing or running it never executes arbitrary install hooks.
10. As a consumer installing `@langdrift/ex`, I want it to pull zero npm runtime dependencies (and never request the private `@langdrift/contracts` package from the registry), so that installation succeeds in a clean project.
11. As a consumer inspecting the npm tarball, I want only `bin/ex`, the four `dist/ex-*` platform binaries, `LICENSE`, `README.md`, and `package.json` included, so that Go source files, `src/testdata` fixtures, and `.build.lock` artifacts are never shipped.
12. As a user running `ex version`, I want it to print `ex 0.1.0` instead of `ex 0.0.0-dev` in packaged builds, so that I know the exact published release I am running.
13. As a contributor running `pnpm turbo run build`, I want `@langdrift/ex` to build only the host binary `dist/ex` after `@langdrift/contracts`, so that local workspace builds remain fast and ordered.
14. As a maintainer running `pnpm pack` or `pnpm publish` in `packages/ex`, I want `prepack` to cross-compile all four static binaries via `sh build.sh --release`, so that a packed tarball is always complete and self-contained.
15. As a consumer reading the package metadata and license on npm, I want `MIT` declared with a matching `LICENSE` file (`Copyright (c) 2026 LangDrift`), `homepage`, `repository`, `os`, `cpu`, and `"access": "public"`, so that licensing and platform support are explicit.
16. As a reader of `packages/ex/README.md`, the root `README.md`, and its translations (`pt-BR`, `zh-Hant`, `ja`, `vi`), I want the documentation to show how to run `@langdrift/ex` via `npx` as well as how to build from source, so that the docs reflect the published package.

## Implementation Decisions

- **Package manifest (`@langdrift/ex`).** Version becomes `0.1.0`, `"private"` is removed, and `"publishConfig": { "access": "public" }` is set. Metadata fields: `"description"`, `"license": "MIT"`, `"homepage": "https://langdrift.md"`, `"repository": { "type": "git", "url": "git+https://github.com/jonatassales/lang-drift.git", "directory": "packages/ex" }`, `"os": ["darwin", "linux"]`, `"cpu": ["x64", "arm64"]`, and `"bin": { "ex": "./bin/ex" }`.
- **Workspace dependency ordering.** `"@langdrift/contracts": "workspace:*"` moves from `"dependencies"` to `"devDependencies"` in `@langdrift/ex`. Turborepo still orders `@langdrift/contracts` before `@langdrift/ex`, while npm consumers do not attempt to install `@langdrift/contracts` from the public registry.
- **Published files list.** `"files"` in `@langdrift/ex` is restricted to `["bin/ex", "dist/ex-darwin-arm64", "dist/ex-darwin-amd64", "dist/ex-linux-amd64", "dist/ex-linux-arm64", "LICENSE", "README.md"]`. Host `dist/ex`, Go sources, test fixtures, and lock directories stay out of the tarball.
- **POSIX `sh` launcher (`bin/ex`).** Resolves its own symlinked path (using the same POSIX symlink-resolution loop as `platform/cli/drift` so `node_modules/.bin/ex` and `npx` cache symlinks resolve back to the package root), detects OS via `uname -s` (`Darwin` → `darwin`, `Linux` → `linux`) and architecture via `uname -m` (`arm64`/`aarch64` → `arm64`, `x86_64`/`amd64` → `amd64`), locates `dist/ex-<os>-<arch>` (falling back to `dist/ex` when running against a local host-only workspace build), restores executable permissions (`chmod 755`) if the package manager unpacked the `dist/ex-*` file with mode `0644`, and `exec`s the binary with `"$@"`. If the platform is unsupported or the binary is missing, it writes a diagnostic to stderr and exits `4`.
- **Build script (`build.sh`).** Reads the package version from `package.json` and passes `-ldflags="-s -w -X main.version=<version>"` with `CGO_ENABLED=0`. Without arguments (used by `scripts.build`), it compiles `dist/ex` for the host OS/arch. With `--release` (used by `scripts.prepack`), it cross-compiles `dist/ex-darwin-arm64`, `dist/ex-darwin-amd64`, `dist/ex-linux-amd64`, and `dist/ex-linux-arm64` atomically via temporary files under the build lock, plus `dist/ex` for the host.
- **License.** MIT `LICENSE` file (`Copyright (c) 2026 LangDrift`) placed at the repository root and in `packages/ex/LICENSE` so both the repository and the packed npm tarball carry the full license text.
- **Documentation.** `packages/ex/README.md`, `.agents/context.md`, root `README.md`, and the four translated READMEs under `platform/docs/` (`README.pt-br.md`, `README.zh-hant.md`, `README.ja.md`, `README.vi.md`) describe `@langdrift/ex` on npm and `npx @langdrift/ex` alongside building from source.

## Testing Decisions

A good test verifies observable behavior at a public boundary—exit status, stdout, stderr, and packed artifact execution—never internal shell variables.

- **Seam 1 — CLI launcher and build script contract (`packages/ex/src/main_test.go`):** Go integration tests invoke `bin/ex` and `build.sh` in disposable directories, verifying that `bin/ex` dispatches the matching `dist/ex-<os>-<arch>` binary with arguments and exit codes preserved, falls back to `dist/ex` in a host build, exits `4` with a clear stderr message on an unsupported `uname` output, and that `build.sh` injects the version from `package.json` into `ex version`.
- **Seam 2 — Clean-room tarball and `npx` verification:** Pack `@langdrift/ex` with `pnpm pack`, inspect the tarball entry list to confirm only the allowed files are present, and execute the tarball through `npx` in an isolated temporary directory outside the monorepo (`version`, `--help`, `init`).
- **Seam 3 — Repository guardrails:** `pnpm turbo run build test typecheck check`, `platform/cli/drift check dead-code`, and `platform/cli/drift check localization`.

## Out of Scope

- Automated GitHub Actions release workflow (`.github/workflows/publish.yml`).
- Publishing `@langdrift/contracts`, `@langdrift/tokens`, or `@langdrift/web` to npm.
- Windows (`win32`) support.
- Split per-platform optional npm packages (`@langdrift/ex-darwin-arm64`, etc.).

## Further Notes

- Spec 0002 hardened `@langdrift/ex` as a workspace package; this spec opens its public npm distribution under `@langdrift/ex`.
