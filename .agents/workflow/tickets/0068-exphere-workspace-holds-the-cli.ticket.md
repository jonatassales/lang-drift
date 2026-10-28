# 0068: The Exphere workspace holds the `ex` CLI and the Microruntimes contracts

**Spec:** .agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md

**What to build:** `packages/ex` moves to `exphere/cli`, and a new Go module, `exphere/contracts`, takes the Microruntimes contracts out of it. The CLI's module path becomes `github.com/jonatassales/tao/exphere/cli` and every import follows; the package stays `@langdrift/ex` and the binary `ex`. `exphere/contracts`, package `@langdrift/exphere-contracts`, depends on nothing beyond the standard library: `configuration` takes the `runtime.md` fence from the CLI's `domain/connection` with its tests, and the lexical rules for identifiers, paths, ports, the platform path and image references; `protocol` takes the wire types from the CLI's `application/connected` and names the MCP version and the five tools. The CLI keeps its names through aliases and delegation, so nothing it does changes. A root `go.work` ties `exphere/cli`, `exphere/contracts` and `apps/extub` together; `pnpm-workspace.yaml` adds `exphere/*`, the lockfile follows, and the CI action reads the Go version from `go.work`. AGENTS.md, `.agents/context.md`, the docs pages and the READMEs name the new places.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] `packages/ex` no longer exists; `exphere/cli` builds `dist/ex` and its full suite passes unchanged with the race detector.
- [x] The fence's tests run in the contracts module, and the lexical rules have table tests there.
- [x] The MCP client uses the protocol's version and tool constants; no literal remains.
- [x] The contracts module imports nothing outside the standard library.
- [x] No live reference to `packages/ex` remains outside specs, tickets, dated research and the CLI's own scaffold assets.
- [x] `pnpm install --frozen-lockfile` passes with the updated lockfile.

### Verification

- `pnpm turbo run build test typecheck --filter=@langdrift/ex --filter=@langdrift/exphere-contracts --filter=@langdrift/extub`, `pnpm turbo run check typecheck`, `pnpm turbo run build --filter=docs`, `drift check dead-code`.
