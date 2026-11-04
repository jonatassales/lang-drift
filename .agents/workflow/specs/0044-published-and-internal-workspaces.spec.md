# Published and internal workspaces

**Status:** ready-for-agent

## Tickets

- [ ] 0211: the workspace knows `lib/` and `infrastructure/`
- [ ] 0212: the event log and the archive become `infrastructure/messaging` and `infrastructure/storage`
- [ ] 0213: the event log package is `streaming`
- [ ] 0214: the design tokens package is `tokens`
- [ ] 0215: the shared primitives package is `web`
- [ ] 0216: the contracts and the CLI move to `lib/`
- [ ] 0217: the CLI is published as `@langdrift/cli` 0.4.0
- [ ] 0218: the stub API moves to `apps/microservices/console-stubs`
- [ ] 0219: console-stubs is its name on every visible surface
- [ ] 0220: open tickets and documents follow the new layout

## Problem Statement

The repository does not tell apart what LangDrift publishes to its users from what only this monorepo uses. Everything under `packages/` carries the `@langdrift/*` scope: the published CLI (`@langdrift/ex`), its Go contracts (`@langdrift/contracts`, private), and three packages nothing outside the repository will ever install (`@langdrift/web`, `@langdrift/tokens`, `@langdrift/events`). A reader cannot tell from a name whether a change ships to a user, and `SECURITY.md` lists private packages as published ones.

Two placements made under spec 0031 and spec 0032 also say the wrong thing. The Redis Streams event log and the MinIO archive sit under `apps/microservices/` beside the Nest services, as if they were microservices; they are infrastructure the products run on, the site and the console, and nothing to do with the CLI. And the Go stub API the console reads is called `extub`, a name carried over from Exphere that names neither what it serves nor whom.

## Solution

Four workspaces, each with one meaning, and a naming rule that follows it:

- **`lib/`** holds only what is published, under the `@langdrift/*` scope. `@langdrift/cli` is the CLI, today `@langdrift/ex`; the command stays `ex`. `@langdrift/core` holds the Microruntimes contracts; its Go module lives at `contracts/go/` inside it, leaving room for other surfaces beside it.
- **`packages/`** holds only packages internal to the monorepo, with bare names like the apps: `web`, `tokens` and `streaming` (today `events`; the bare name `events` would collide with Node's built-in module).
- **`infrastructure/`** holds the product infrastructure components, each a folder with a Dockerfile, a README and a minimal private manifest: `messaging` (Redis Streams, the event log) and `storage` (MinIO, the archive).
- **`apps/microservices/`** holds the services: the Nest services and `console-stubs`, the Go stub API the console reads, renamed from `extub` through everything that shows the name, the folder, the module, the environment variables, the console route, the catalog key, the realm and the labels.

The published CLI is released under its new name as `@langdrift/cli` 0.4.0, and `@langdrift/ex` on npm is deprecated with a pointer to it. `@langdrift/core`'s Go module is published the way `ex` consumes it, as a Go module by git tag; the npm manifest exists and stays private until a non-Go surface exists.

## User Stories

1. As a maintainer, I want the folder a package lives in to tell me whether it ships to users, so that I know what a change in it reaches.
2. As a maintainer, I want only published packages to carry the `@langdrift/*` scope, so that the scope means "published".
3. As a maintainer, I want internal packages named like the apps, with bare names, so that one convention covers everything the monorepo alone consumes.
4. As a CLI user, I want to install the CLI as `@langdrift/cli` and still run it as `ex`, so that the package name says what it is and my scripts keep working.
5. As a CLI user who installed `@langdrift/ex`, I want npm to tell me it was renamed and to what, so that I move without searching.
6. As a consumer of the CI gates scaffold, I want `ex init` to write the new package name into the generated workflow, so that a fresh project pins a package that keeps receiving releases.
7. As a maintainer, I want the contracts to live inside `@langdrift/core` under a Go surface folder, so that a TypeScript or other surface can sit beside it later without another move.
8. As a maintainer building the CLI outside the workspace, I want its module to require the contracts module by a published path, so that `go build` works once the module is tagged.
9. As a maintainer, I want the Redis Streams and MinIO components under `infrastructure/`, named `messaging` and `storage`, so that the tree says they are what the products run on, not services of their own.
10. As a maintainer, I want each infrastructure component listed as a workspace member, so that pnpm, Turborepo and the `drift` checks treat every folder uniformly.
11. As a developer running a service image locally, I want the build commands in the Dockerfiles and READMEs to name the new paths, so that copying them works.
12. As a developer of the Nest services, I want them to depend on `streaming`, so that the event log package is named for what it does.
13. As a developer of the console or the site, I want to import `web/...` and `tokens/...`, so that the import says the dependency is internal.
14. As a maintainer, I want the Tailwind, TypeScript, React and Nest rules to name the packages by their new names, so that the rules describe the tree as it is.
15. As a console developer, I want the stub API called `console-stubs`, so that its name says what it serves and for whom.
16. As a console operator, I want the environment variables named `CONSOLE_STUBS_URL`, `CONSOLE_STUBS_TOKEN` and `CONSOLE_STUBS_PORT`, so that the configuration names the service.
17. As a console developer, I want the proxy route, the data channel label, the auth realm, the service field, the load-test target and the catalog key to carry the new name, so that no visible surface still says `extub`.
18. As a console visitor in any served locale, I want the unavailable-stubs notice to keep its meaning after the key rename, so that the parity test and the copy both hold.
19. As the owner, I want the Vercel steps that only I can take listed in order, so that the hosted preview and the console keep working across the rename.
20. As a maintainer, I want the `Ticket:` trailer rule to cover `lib/` and `infrastructure/`, so that every change that reaches users or products names its ticket.
21. As a maintainer, I want the dead-code and localization checks to scan the new workspaces, so that nothing escapes review by moving.
22. As a maintainer, I want CI to build, vet and test the Go modules at their new paths and to publish the CLI binaries from `lib/cli`, so that the pipeline stays green through the move.
23. As a maintainer, I want historical specs and tickets left as they were written and only open tickets updated, so that history stays truthful and the backlog stays actionable.
24. As a maintainer, I want the `ex` test fixtures that mirror a consumer project left untouched, so that the CLI's behaviour against an existing project does not change.
25. As a reader of `SECURITY.md`, I want the published packages listed to be the ones that are published, so that a report reaches the right surface.
26. As a reader of `AGENTS.md` and the context document, I want the layout described with the four workspaces and the naming rule, so that the next agent places a new package correctly.

## Implementation Decisions

- **Four workspaces.** `pnpm-workspace.yaml` gains `lib/*` and `infrastructure/*` beside `apps/*`, `apps/microservices/*` and `packages/*`. The `Ticket:` trailer rule and the `drift` dead-code workspace loops gain the two new roots.
- **Published packages.** `lib/cli` is the CLI (Go module `.../lib/cli`, npm `@langdrift/cli` 0.4.0, bin `ex`). `lib/core` is the package `@langdrift/core`, private for now, whose Go module is `.../lib/core/contracts/go`. The root Go workspace ties them for local work; the CLI's module requires the core module by path once the module is tagged (a `require` of an untagged version makes Go fetch it even in workspace mode), and from then on the CLI builds outside the workspace. Tags become `lib/cli/vX.Y.Z` and, for the Go module in its subdirectory, `lib/core/contracts/go/vX.Y.Z`.
- **Internal packages.** `web`, `tokens` and `streaming` keep their folders under `packages/` (`events` moves to `streaming`) and lose the scope. Every dependency on them keeps the `workspace:*` protocol: all three bare names exist on npm, and the protocol is what pins the local package.
- **Infrastructure components.** `infrastructure/messaging` and `infrastructure/storage` hold the Dockerfile and README moved from `apps/microservices/redis` and `minio`, plus a manifest with a name and `private: true` and nothing else. ADR 0004 stands: the component is its folder and its Dockerfile, and a fence names it by the folder, so a future microruntime of this repository names `messaging` and `storage`.
- **Nest rule scope.** The Nest rule covers the Nest services only; an infrastructure component is not a microservice and lives under `infrastructure/`; `console-stubs` is a Go service under `apps/microservices/` outside that rule.
- **console-stubs.** Folder `apps/microservices/console-stubs`, workspace name `console-stubs`, Go module `.../apps/microservices/console-stubs`. The environment contract becomes `CONSOLE_STUBS_URL`, `CONSOLE_STUBS_TOKEN` and `CONSOLE_STUBS_PORT`; the console's proxy route becomes `/api/console-stubs/...`; the realm, the service field, the load-test target, the data channel label, the log prefixes and the catalog key follow. The console's dev task starts `console-stubs#dev`.
- **The CLI's own text.** The launcher's reinstall hint, the setup message and the CI gates scaffold name `@langdrift/cli`; the scaffold's fixture changes in lockstep. The consumer scaffold keeps `npx`: the tooling rule governs this repository's automation, not workflows generated into a consumer's repository.
- **Documentation and badges.** The root README, the CLI README, the translated READMEs, the website footer badge and the dashboard widget name `@langdrift/cli` and read the version from the CLI manifest rather than a literal. `SECURITY.md` lists `@langdrift/cli` and `@langdrift/core` as published and nothing else.
- **History.** Specs and closed tickets stay as written. Open tickets that cite moved paths or the old package name are updated: the console microruntime ticket, the Book a Demo microruntime ticket (which pins the CLI at the root and names the components), the Safe Seat CLI tickets and the four open stub tickets.
- **Order of work.** Each phase lands green on its own: scaffolding; infrastructure; `streaming`; `tokens` then `web`; `lib/core` with `lib/cli` in one move (Go resolution is through the workspace); the npm rename; the stub move then the stub rename; the ticket sweep. The release is the owner's.
- **Owner-only steps.** Publishing `@langdrift/cli@0.4.0`, deprecating `@langdrift/ex`, tagging, the GitHub Release, the Vercel Root Directory of the stubs project and the recreated `CONSOLE_STUBS_*` variables on the console and stubs projects. The agent never enters a credential or a value.

## Testing Decisions

A good test here checks that the repository behaves the same from the outside after each move: a package resolves under its new name, a binary prints its version, a service answers on its route, an image builds. No test asserts folder paths for their own sake.

Seams, existing ones first:

- **The workspace.** `pnpm install --frozen-lockfile` and `pnpm turbo run build typecheck test check` are the seam for every rename: a missed import, dependency or `transpilePackages` entry fails there.
- **The `drift` checks.** `drift quality check`, `drift quality dead-code`, `drift check localization` and `drift check biome` are the seam for the rules, the catalogs and the workspace loops; the console catalog parity test keeps the renamed key in every served locale.
- **Go.** `go vet` and `go test -race` per module, and the CLI's own tests, which already compare the scaffold with its fixture and read the version from the manifest, are the seam for the module moves and the CLI's text.
- **Docker.** `docker build` of the two infrastructure components and of the two Nest service images from the repository root is the seam for the Dockerfile paths and the turbo prune of `streaming`.
- **The console.** Its existing route, proxy-matcher, link and load-test unit tests are the seam for the stub rename; they change their expectations, not their shape.
- **Retired identifiers.** A `git grep` of each old name is part of each ticket's verification; the only allowed hits are the workflow history and the CLI's consumer fixtures.

Prior art: the CLI's `main_test.go` reads the manifest version and runs the launcher; the console's `extub-link.test.ts`, `proxy-matcher.test.ts` and `entity-id.test.ts`; the Nest services' controller tests with fakes; the website and console catalog parity tests.

## Out of Scope

- Publishing `@langdrift/core` to npm, which waits for a surface that is not Go.
- Renaming the `ex` command or its Go packages.
- Hosting the Nest services, the event log and the archive (spec 0031 left it out; the hosting conversation of 2026-11-04 continues separately).
- The console's missing Vietnamese catalog, a pre-existing gap against the localization rule.
- Changing how a consumer project references an infrastructure component from a package; the components are referenced by path inside this repository.
- Rewriting historical specs, tickets or the CLI's consumer fixtures.

## Further Notes

- The staged move that preceded this spec (stubs under `apps/microservices/extub`, components under `platform/infrastructure/`) was discarded; this spec supersedes it.
- The Vercel project of the stubs builds from Root Directory `apps/extub` today; it must point at the new folder when the move merges, or its preview stops building.
- Version 0.4.0 continues the `ex version` sequence; the rename is the reason for the minor.
