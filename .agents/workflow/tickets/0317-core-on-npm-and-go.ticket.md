# 0317: `@langdrift/core` 0.2.0 is published on npm and as a Go tag

**Spec:** .agents/workflow/specs/0044-published-and-internal-workspaces.spec.md, which keeps the npm manifest private "until a non-Go surface exists". The judge's JSON Schemas and samples (0222) are that surface. The owner asked on 2026-11-05 to publish it on npm, then delegated the Go tag ("Add the Go tag too"), which spec 0044 keeps among the owner-only steps and which leaves 0228. As with `v0.1.0`, the tag goes on the merge commit on `main` after the owner merges this ticket's first pull request, and the CLI's `require` follows in a second one.

**What to build:** Anyone running the engine outside this repository installs `@langdrift/core` 0.2.0 from npm and imports the judge's contract as `@langdrift/core/judge/<file>.json`, the same files the console's tests and the contracts' own tests read. The package carries those JSON Schemas and samples, its README and its license, and nothing else. The same release is the Go module's tag `lib/core/contracts/go/v0.2.0`, the first with the `judge` package, and the CLI's `go.mod` requires it, so the CLI builds outside this repository against the published contracts. The owner runs the npm publish with their account.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] The manifest is public at 0.2.0, with a description, the Apache-2.0 license, the repository and `publishConfig.access: public`, and installs on every platform.
- [x] `pnpm pack --dry-run` lists the twelve files of `contracts/go/judge/schema/`, `README.md`, `LICENSE` and `package.json`, and no Go source.
- [x] The tag `lib/core/contracts/go/v0.2.0` exists on GitHub and `proxy.golang.org` serves the module at `v0.2.0` with its `judge` package.
- [x] The CLI's `go.mod` requires the contracts at `v0.2.0`, and `GOWORK=off go build ./src` and `go test ./...` pass in `lib/cli`.
- [x] The core README and the module README say what each surface carries and how to get it; `SECURITY.md` names both surfaces.
- [ ] `pnpm view @langdrift/core version` prints `0.2.0` once the owner publishes, and the package installs and loads a schema in a fresh directory.

## Dev part

**Override:** No part qualifies: release configuration (policy, Categories: never configuration). 2026-11-05.

## Verification

- `pnpm --dir lib/core pack --dry-run` and `pnpm --dir lib/core publish --dry-run --no-git-checks`
- `pnpm turbo run test typecheck --filter=@langdrift/core --filter=console`
- `curl https://proxy.golang.org/github.com/jonatassales/tao/lib/core/contracts/go/@v/v0.2.0.info`, then `GOWORK=off go -C lib/cli build ./src` and `GOWORK=off go -C lib/cli test ./...`
- After the owner's `pnpm --dir lib/core publish --access public`: `pnpm view @langdrift/core version files`, then in a fresh directory `pnpm add @langdrift/core@0.2.0` and `node --input-type=module -e "import schema from '@langdrift/core/judge/request.schema.json' with { type: 'json' }; console.log(schema.title)"`
