# The LangDrift CLI answers to `dft`

**Status:** ready-for-agent

## Tickets

- [ ] 0313: the CLI answers to `dft`
- [ ] 0314: the site, the console and the docs say `dft`
- [ ] 0315: the repository's documents say `dft`

## Problem Statement

The CLI that plans, runs and records microruntimes is published as `@langdrift/cli`, but everyone meets it as `ex`: the command people type, the launcher, the release binaries, the ledger's provenance, the website's terminal demos, the console's connection instructions, the docs in five languages and this repository's own guides. `ex` was a code name. It names neither the brand nor what the tool does, and the product it ships with already carries the brand in its name (LangDrift FinOps, Microruntimes). On 2026-11-05 the owner decided that the CLI is called **LangDrift CLI** and that its command is **`dft`**.

## Solution

From the release that carries this change, a person installs `@langdrift/cli` as before and types `dft`: `dft init`, `dft agent plan`, `dft agent run`, `dft runtime up`. Every surface that shows the command shows `dft`, and every sentence that names the tool calls it the LangDrift CLI. `ex` stops working: there is no alias. The release notes say what changed. Records written by earlier versions stay readable. Closed specs and tickets keep the name they were written with, and the glossary records that `ex` is the former name of `dft`.

## User Stories

1. As a developer installing LangDrift, I want to type `dft`, so that the command I learn carries the brand's name.
2. As a developer, I want `npx @langdrift/cli init` to keep working, so that the install line I already know does not change.
3. As a developer reading `dft --help`, I want the help, the usage lines and the error messages to say `dft`, so that what the tool prints matches what I typed.
4. As a developer on macOS or Linux, I want the release binary for my platform to be named `dft-<os>-<arch>`, so that the downloads say what they are.
5. As a developer whose project was set up by 0.4.0, I want the release notes to tell me that `ex` became `dft`, so that a script calling `ex` fails with a reason I can find.
6. As a developer, I want ledger records written by 0.4.0 with `ex` as their provenance to stay readable, so that upgrading does not break my history.
7. As a reader of the ledger, I want new records to name `dft` as the tool that set them, so that provenance names the tool I ran.
8. As a website visitor, I want the CLI, Microruntimes and CI Gates pages to show `dft` in their commands and terminal demos, in every locale, so that what I copy works.
9. As a website visitor reading in Portuguese, Traditional Chinese, Japanese or Vietnamese, I want the same commands as the English page, so that the meaning survives translation.
10. As a console user connecting a project, I want the connection instructions to read `dft init`, `dft agent connect` and `dft agent run`, so that I can paste them.
11. As a docs reader in any of the five languages, I want the README to call the tool the LangDrift CLI and to show `dft`, so that the docs agree with the binary.
12. As a docs reader, I want the terminal screenshots to show `dft` running, so that the pictures match the command.
13. As the owner, I want the dashboard and Microruntimes widgets to say "LangDrift CLI" and show the published version, so that the repository's status pages use the product's name.
14. As an agent working in this repository, I want AGENTS.md, the context guide and the rules to say `dft`, so that I never bring `ex` back.
15. As an agent picking up an open ticket, I want its text to say `dft`, so that what I implement uses the current name.
16. As a reader of a closed spec or ticket, I want it to keep the name it was written with, and the glossary to say that `ex` is the former name of `dft`, so that history stays true and still reads.
17. As a contributor, I want the glossary to tell `dft`, the product's command, apart from `drift`, this repository's own tooling, so that the two short names are not confused.
18. As the owner, I want the rename to ship with 0.5.0, together with spec 0045, so that the breaking change comes once.

## Implementation Decisions

- **The name.** In prose the tool is the **LangDrift CLI**. The command, and every identifier named after the command, is `dft`. The npm package stays `@langdrift/cli`. The Go module path does not change: it names the folder, not the command.
- **No alias.** The package's `bin` maps only `dft`, so `ex` stops resolving with the release that carries the change. No shim, deprecation command or second bin entry is added (General, Authority 6).
- **What takes `dft`:**
  - the root command's `Use`, with every help, usage and error string that names the command;
  - the launcher, renamed with its internal variables, which `bin` and `files` list;
  - the launcher test;
  - the release binaries and their checksum file, built by the main workflow, with its job and artifact names;
  - the widgets' variables and labels ("LangDrift-X (CLI)" becomes "LangDrift CLI");
  - the comments and READMEs of `lib/cli` and `lib/core`.
- **Generated files.** Whatever the CLI writes into a consumer's project that calls the command (guardrails, harness files, microruntime scripts, record notes) calls `dft`. Fixtures under `testdata/` that hold expected output follow the generator. Fixtures that stand for a project written by an earlier version keep `ex`, and a test reads them.
- **Cluster names keep `ex-`.** A project's kind cluster is named `ex-<id>`, and its ownership record names that cluster. Renaming the prefix would leave every cluster 0.4.0 created behind, so the prefix stays and a comment says why.
- **Provenance.** New ledger and harness records set by the tool carry `dft` where they carried `ex`. Readers accept any value there, as they do today; nothing validates the old value, so records from 0.4.0 stay readable.
- **Surfaces.**
  - The website's terminal demos, command lines and catalogs say `dft` in all five locales, with catalog parity kept.
  - The console's connection instructions say `dft`.
  - The five docs READMEs call the tool the LangDrift CLI and show `dft`.
  - Their screenshots are retaken from the renamed binary and named after the command they show. The capture is the owner's own terminal session, so the owner retakes it.
- **Documents.**
  - The current documents say `dft`: AGENTS.md, the context guide, the Safe Seat policy, `SECURITY.md` and the infrastructure READMEs. The root README is the public landing page, so it ships with the surfaces.
  - The context guide's row for the root's pinned CLI keeps `ex` while the pin is 0.4.0; ticket 0228 moves the pin.
  - So do the rules that name the command: `10-testing` (the microruntime the CLI generates) and `11-nestjs` (the CLI derives readiness from the `HEALTHCHECK`). This updates a name, not a rule.
  - The glossary gains `dft` (the LangDrift CLI's command), `drift` (this repository's tooling, never shipped) and `ex` (the former name of `dft`, kept in closed specs and tickets).
- **Open versus closed.**
  - A ticket is open while its spec's Tickets list leaves it unchecked and no commit on `main` carries its `Ticket:` trailer, because the lists lag the merges. A spec is open while any of its tickets is.
  - Open specs and tickets are rewritten to `dft` where they plan work. Where they record history (what 0.4.0 did, a command that only ever existed as `ex`, a ticket already merged), they keep `ex`, and a dated note at the top says so. Closed ones are not touched.
  - Ticket 0228 (release 0.5.0) gains blocking edges on 0313, 0314 and 0315.
- **The internal `drift` keeps its name.** It is repository tooling, and the glossary separates it from `dft`.
- **Release.** The change ships with `@langdrift/cli` 0.5.0. Publishing stays the owner's, through ticket 0228.
- **Timing.** A merge to `main` deploys the site and the console, and npm serves 0.4.0 until 0.5.0 is published. The CLI and the documents (0313, 0315) merge when ready, because nobody types the tree's build. The surfaces (0314) go in their own pull request, which merges with the release, so that no visitor reads `dft` while npm still ships `ex`.

## Testing Decisions

- **Behavior is the seam.** A good test runs the command or reads what the build ships: what `dft` prints and writes, what the package lists, what a page renders. It never asserts the name of a variable.
- **The CLI.**
  - The launcher test runs the renamed launcher.
  - The root command's tests assert that help and usage say `dft`.
  - The golden fixtures under `lib/cli/src/testdata/langdrift/expected/` hold the new provenance.
  - A test reads a 0.4.0 record that names `ex`.
  - `pnpm pack --dry-run` in `lib/cli` lists `bin/dft.sh` and no `bin/ex.sh`.
  - `go vet` and `go test -race` pass in `lib/cli` and `lib/core`.
  - Prior art: the launcher test and the golden-record tests in `lib/cli/src`.
- **The surfaces.**
  - The website's catalog test and `drift check localization` keep parity and valid ICU in every locale.
  - The content unit tests that assert commands read `dft`.
  - `pnpm turbo run check typecheck test build` passes for `website` and `console`.
  - Prior art: `apps/website/src/app/template/i18n/catalogs.unit.ts` and `microruntimes-content.unit.ts`.
- **The residue check.** A `git grep` for `ex` used as the command finds no match outside closed specs and tickets and the 0.4.0 fixtures. The check runs once per ticket, by hand. It is not a new detector.

## Out of Scope

- An `ex` alias, a deprecation shim or a migration command.
- Renaming the npm package, the Go module or the repository's `drift` tooling.
- Rewriting closed specs, tickets, ADRs or research notes.
- Publishing the release and deploying it, which ticket 0228 owns.
- Any change to what the CLI does.

## Further Notes

The owner settled every decision on 2026-11-05 by accepting the recommended answers:
- everything named after the command takes `dft`;
- no alias;
- the package name stays;
- new provenance says `dft`;
- closed history stays and open work is rewritten;
- every locale and the screenshots change;
- `drift` keeps its name;
- the change ships with 0.5.0.
