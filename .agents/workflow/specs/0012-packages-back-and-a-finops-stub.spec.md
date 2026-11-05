# The CLI and its contracts return to packages, extub stubs FinOps alone, and every app runs on its own port

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

Spec 0010 grouped the `ex` CLI and the Microruntimes contracts in a workspace of their own, `exphere/`, as the engine of a product line. The owner has since pivoted the project to one agentic FinOps platform, so a separate workspace with a separate name adds a concept no one needs: the CLI and its contracts are packages like any other. The console's stub API, `apps/extub`, still seeds a Vision product with its own fixture, AI Voice block and drift workspace, although Vision is no longer a product. Vision also remains in the docs app, in AGENTS.md and in the context file's glossary and navigation. A developer cannot run the docs and the console at the same time, because both bind port 3001, and no command starts the console together with the stub it reads.

## Solution

The CLI moves back to `packages/ex` and the contracts to `packages/contracts`, published as `@langdrift/contracts`; their Go module paths follow the folders, and the Exphere workspace and name go away. The stub API serves only FinOps and runtimes: one seeded product, Microruntimes, with its runtimes, executions, history, secrets and FinOps readings. Vision leaves the docs app and the repository documents, which describe one FinOps product with Microruntimes as its engine and the console that lives here. Three root commands start the apps on their own ports: `pnpm website` on 3000, `pnpm console` on 3001 together with the stub on 3401, and `pnpm docs` on 3002.

## User Stories

1. As a developer, I want the `ex` CLI under `packages/ex`, so that it sits beside the other packages.
2. As a developer, I want the Microruntimes contracts under `packages/contracts` as `@langdrift/contracts`, so that the folder and the package share one name.
3. As a developer, I want Go module paths that match the folders, so that an import tells me where the code lives.
4. As a developer, I want the workspace files, the lockfile and the CI cache to follow the move, so that installs and builds work without manual steps.
5. As a developer, I want the history of every moved file to follow it, so that `git log --follow` still reaches its origin.
6. As a contributor, I want no document to describe an Exphere workspace, so that the vocabulary matches the repository.
7. As a console user in development, I want the stub to seed only the Microruntimes product, its runtimes and their FinOps readings, so that the console shows the product the platform sells.
8. As a console user, I want the local folder picker to offer no Vision folder, so that I cannot recreate a product that no longer exists.
9. As a console user, I want the product switcher to stay quiet with one product, so that the interface does not offer a choice it does not have.
10. As a developer, I want the stub's tests to prove product isolation with a product they create, so that the tests do not depend on a seeded product that left.
11. As a docs reader, I want no Vision page or navigation entry, so that the docs describe the platform as it is.
12. As a docs reader, I want the platform pages to describe the console and the stub that live in this repository, so that nothing claims they live upstream.
13. As a contributor, I want AGENTS.md to describe one FinOps product with Microruntimes as its engine, so that I start from the current idea.
14. As a contributor, I want the glossary to define the console, the platform product and the classification color without Vision, so that terms match the product.
15. As a developer, I want `pnpm console` to start the console and the stub together, so that the console has data the moment it opens.
16. As a developer, I want the docs on port 3002, so that I can run the docs, the console and the website at the same time.
17. As a developer, I want the commands documented with their ports, so that I find each app without reading scripts.
18. As an agent using the preview tools, I want launch configurations for the website, the console and the docs on their ports, so that each can be previewed by name.

## Implementation Decisions

- **Move.** The two folders move with `git mv` in the same commit that rewrites their module paths, the workspace entries and the import lines, so rename detection holds. Imports are reformatted afterwards because the two module paths now sort in the opposite order.
- **Names.** `@langdrift/ex` keeps its name and binary. The contracts package becomes `@langdrift/contracts`; the CLI keeps its workspace dependency on it under the new name, because Turborepo orders builds by it. Go comments that said "Exphere contracts" say "the Microruntimes contracts".
- **Workspace.** The pnpm workspace drops its `exphere/*` entry, which `packages/*` now covers. The Go workspace lists the two packages and the stub. The CI setup action caches the CLI's `go.sum` at its new path.
- **Documentation of the move.** The workspace README goes; its boundary and Go workspace notes move into the two package READMEs. AGENTS.md, the context file, the docs pages that cite the CLI's path and the `.langdrift` sample stop naming Exphere; the sample's mirror records name the console as their source. Historical specs, tickets and the research note stay as written, and so do the `.langdrift` sources whose hashes the sample run records.
- **Stub API.** The Vision fixture types, the product field that carries them, the seeded Vision product and its local folder go. The API keeps every route; it seeds one product and four runtimes. Tests that used the seeded Vision product for isolation create a second product through the API instead.
- **Console.** No change: it never read the Vision fixture, and its product switcher already disables itself with fewer than two products.
- **Docs app.** The Vision page and its navigation entry go. The home and the platform pages describe one FinOps product with Microruntimes as its engine and the console and stub in this repository; the architecture table lists them.
- **Repository documents.** AGENTS.md describes the FinOps product, the pivot in what remains true, and the console and stub that live here. The context file's navigation and glossary drop Vision and the product switcher between products. The editor's folder-icon mapping for Vision goes.
- **Rules.** The frontend rule stops listing stylesheets that no longer exist. The localization rule no longer keeps "Product Vision" in English and keeps its rule that meaning survives translation.
- **Commands and ports.** The docs app binds 3002 in development and in production. The console app's Turborepo configuration runs the stub's development task alongside its own with `with`, and a root script `console` filters `turbo run dev` to the console. The context file's rule allowing only two root launchers allows three. The preview launch configuration names the website, the console and the docs with their ports; the docs entry runs development mode rather than a production server that needs a build first.

## Testing Decisions

- A good test observes what a user or a build sees: the stub's HTTP responses, the Go packages' behavior, and the commands answering on their ports. It does not assert file layout.
- **Seams.** The existing ones suffice: the Go test suites of the CLI, the contracts and the stub; the repository's check, typecheck, test and build tasks; the dead-code, localization and history checks; and a manual run of the three commands at once. No new seam.
- **Stub tests.** The home lists one product. Widget isolation and secret isolation are proven against a product the test creates. The Vision fixture test goes. The other tests keep passing unchanged.
- **Move.** The CLI's suite passes under the race detector after the move, and its build script still writes the binary next to it.
- **Prior art.** The stub's handler and activity tests, and the CLI's suite run in ticket 0068.

## Out of Scope

- The console's own Exphere branding: its logo, navigation text, messages and cookie name.
- The planned server that authors runtime definitions, and any change to the protocol or the fence.
- Historical specs, tickets and research notes, which keep the paths and names of their time.
- The website, its copy and its brand line.

## Further Notes

- Spec 0010 created the workspace; this spec reverses that part of it and leaves the rest of 0010 standing.
- The sample run under `.langdrift/` hashes the console's manifest and the stub's Go module file; neither changes here.
