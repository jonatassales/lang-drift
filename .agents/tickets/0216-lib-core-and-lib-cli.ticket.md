# 0216: the contracts and the CLI move to `lib/`

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The published code moves to the `lib/` workspace in one step, because Go resolves the two modules through the root workspace. The contracts become the Go surface of the core package, at `lib/core/contracts/go`, with the core package's manifest named `@langdrift/core` and private, and a README that names its surfaces. The CLI becomes `lib/cli` with its Go module at that path, every import rewritten, its manifest pointing at the new directory and depending on the core package, and its Go module still resolving the core module through the root workspace: a `require` of an untagged version makes Go fetch it even inside the workspace, so the `require` lands with the first tag (0217). CI builds, vets, tests and packages the binaries from the new paths, and the microruntimes widget reads the CLI from `lib/cli`. The npm name stays `@langdrift/ex` in this ticket; 0217 renames it.

**Blocked by:** 0211.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `packages/ex` and `packages/contracts` no longer exist; `lib/cli` and `lib/core/contracts/go` are Go modules named after their paths, both listed in the root Go workspace.
- [ ] The root Go workspace resolves the core module for `lib/cli`, whose `go.mod` names no version of it until the module is tagged; `go vet ./...` and `go test -race ./...` pass in both modules, and `sh lib/cli/build.sh && lib/cli/dist/ex version` prints the manifest version.
- [ ] `lib/core/package.json` is `@langdrift/core`, private, and its README names the Go surface at `contracts/go/`; `lib/cli/package.json` depends on it and names `lib/cli` as its directory.
- [ ] The validate and main workflows, the setup-workspace action and the microruntimes widget use the new paths; `drift check workflows` passes.
- [ ] `AGENTS.md` and the context document name `lib/cli` and `lib/core`.
- [ ] `git grep -n 'packages/ex\|packages/contracts\|tao/packages'` finds only workflow history and CLI fixtures.

## Dev part

**Override:** No part qualifies: a module move with mechanical import rewrites (policy, Categories: never a rename or a migration run by a tool). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile && pnpm turbo run build typecheck test --filter=@langdrift/ex --filter=@langdrift/core`
- `cd lib/cli && GOTOOLCHAIN=local go test -race ./... && sh build.sh && dist/ex version`
- `platform/cli/drift check workflows`
