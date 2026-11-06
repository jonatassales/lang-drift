# 0218: the stub API moves to `apps/microservices/console-stubs`

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The Go stub API the console reads moves under the services as `apps/microservices/console-stubs`, workspace name `console-stubs`, Go module named after the path, with no change in what it serves or how it is configured. The console's dev task starts `console-stubs#dev`, so `pnpm console` still brings both up; CI vets and tests the module at its new path; the Nest services' Docker build context keeps the Go service out. Its Vercel manifest travels with the folder, and the owner points the stubs project's Root Directory at the new folder when this merges.

**Blocked by:** 0211.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `apps/extub` no longer exists; `apps/microservices/console-stubs` is the workspace `console-stubs` and the Go module named after its path, listed in the root Go workspace.
- [ ] `pnpm console` starts the console on 3001 and the stubs on 3401; `go vet ./...` and `go test -race ./...` pass in the module.
- [ ] The validate workflow and the setup-workspace action use the new path; `docker build -f apps/microservices/demo-requests/Dockerfile .` does not copy the Go service into its context.
- [ ] `AGENTS.md`, `SECURITY.md`, the context document and the console README link the new folder.
- [ ] Owner, when merging: the Vercel project `langdrift-stubs` has Root Directory `apps/microservices/console-stubs` and its preview builds.

## Dev part

**Override:** No part qualifies: a folder and module move (policy, Categories: never a rename or a migration run by a tool). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile && pnpm turbo run build typecheck test --filter=console-stubs --filter=console`
- `platform/cli/drift check workflows`
