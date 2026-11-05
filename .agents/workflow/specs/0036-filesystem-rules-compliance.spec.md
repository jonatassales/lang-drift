# The tree complies with the filesystem rules

**Status:** ready-for-agent

## Tickets

- [x] 0264: the handbook and the skill references carry the guide suffix
- [ ] 0265: no path climbs with ../
- [ ] 0266: the Go files take kebab-case

## Problem Statement

The audit of 2026-11-04 measured `.claude/rules/01-filesystem.rule.md` at 97.8% and `drift quality check filesystem` prints 87 findings. Three classes remain, and the detectors miss a fourth:

- **Snake_case Go files (rule 1).** 78 Go files carry underscores the toolchain does not require: 48 sources such as `lib/cli/src/application/cluster_create.go`, `lib/cli/src/presentation/cli/catalog/gates_markdown.go`, `lib/cli/src/infrastructure/scaffold/ci_gates.go` and `apps/microservices/console-stubs/src/application/load_tester.go`, and 30 tests such as `lib/cli/src/infrastructure/connected/local_source_test.go`. Only `_test.go` and the GOOS/GOARCH suffixes are the toolchain's; the rest is a convention the rule replaces with kebab-case.
- **Markdown without a type suffix (rule 2).** Seven documents name no kind: `.agents/context.md`, `.claude/skills/tdd/mocking.md`, `.claude/skills/tdd/tests.md`, `.claude/skills/domain-modeling/adr-format.md`, `.claude/skills/domain-modeling/glossary-format.md`, `.claude/skills/codebase-design/deepening.md` and `.claude/skills/codebase-design/design-it-twice.md`. `.agents/context.md` is linked from `AGENTS.md`, named by six skills and two platform scripts, so a rename has to carry every live reference with it.
- **Paths that climb with `../` (rule 3).** `platform/cli/src/drift.sh` resolves the project root as `$DRIFT_CLI_DIR/../../..`; `platform/cli/src/commands/git-pre-commit.sh` calls `$DRIFT_CLI_DIR/../drift`; `apps/console/src/proxy/proxy-matcher.test.ts` reads `new URL('../proxy.ts', import.meta.url)`; `apps/microservices/demo-requests/README.md` and `apps/microservices/mailer/README.md` link the Nest rules as `../../../.claude/rules/11-nestjs.rule.md`.
- **Detector gaps.** The `../` detectors look only at imports in code (so the console test and the Markdown links pass unseen) and only at shell under `platform/` and `.github/` (so `.claude/hooks/` and `lib/cli/build.sh` are never looked at). The Markdown detector's suffix list is the de facto approved list, and it cannot accept a new kind until it is extended.

The bucket findings are judgement calls that other specs resolve: `apps/{console,website}/src/app/lib/` dissolves under spec 0041, and the multi-function files `helpers.go` and `common.sh` are split by spec 0035, whose single-function files are named after their export.

## Solution

Bring every tracked path to the filesystem rules and make `drift quality check filesystem` report zero findings: rename the unsuffixed documents with a suffix that names their kind and move every live reference with them, remove every `../` from code, shell and Markdown while `drift` keeps running from any working directory, rename the Go files to kebab-case once spec 0035 has finished splitting them, and widen the detectors so that each of these classes stays caught.

## User Stories

1. As the platform owner, I want `drift quality check filesystem` to print zero findings, so that the pre-commit hook and CI stop listing known violations as noise.
2. As a maintainer, I want every Markdown document to say what kind of document it is in its name, so that a glance at a path tells me whether I am reading a rule, a spec, a policy or a guide.
3. As an agent reading `AGENTS.md`, I want its link to the mechanics document to resolve after the rename, so that the handbook stays one hop away.
4. As a skill author, I want each skill's supporting documents to keep resolving from `SKILL.md` and from each other, so that a skill never loads a missing file.
5. As a developer running `drift` from a subdirectory, a symlink in my bin directory or a CI runner, I want it to find the repository root, so that removing `../` costs no capability.
6. As a developer committing, I want the pre-commit hook to reach `drift quality check` through an absolute path, so that the hook works regardless of where Git invokes it from.
7. As a console maintainer, I want the proxy matcher test to read `proxy.ts` as a sibling, so that the test obeys the path rule and still guards the matcher.
8. As a reader of a Nest service README on GitHub, I want the link to the Nest rules to open, so that root-absolute links replace climbing ones without losing navigation.
9. As a Go maintainer, I want every Go file in kebab-case with only `_test.go` and the GOOS/GOARCH suffixes kept, so that the Go tree follows the same naming as the rest of the repository.
10. As an implementer of spec 0035, I want the Go rename to wait for my splits, so that I create new files in kebab-case once and the rename covers only what remains.
11. As a reviewer, I want a `../` in any TypeScript or JavaScript string, in any shell script of the repository and in any Markdown link to be flagged by a detector, so that the rule is checked mechanically rather than by eye.
12. As a maintainer adding a new kind of document, I want the approved suffix list to live in one place, so that extending it is one line.
13. As an `ex` user, I want the record format of my `.langdrift/` folder (`runtime.md`, `run.md`, `*.lap.md`, `*.gr.md`) unchanged, so that a repository-internal naming rule never alters the product's contract.
14. As an author of another compliance spec, I want the renames here marked as prefactors, so that my tickets reference the new names and never race the moves.

## Implementation Decisions

- **Document suffix: `.guide.md`.** The seven unsuffixed documents are guides a person or an agent follows: the repository handbook (`context.guide.md`) and the skills' supporting references (`mocking.guide.md`, `tests.guide.md`, `adr-format.guide.md`, `glossary-format.guide.md`, `deepening.guide.md`, `design-it-twice.guide.md`). One suffix covers all seven; `.reference.md` was the alternative for the skill files and was rejected to keep one kind.
- **References move with the files.** `AGENTS.md`, every `SKILL.md` and skill document that names a renamed file, the comments in the platform CLI that cite the handbook, and the memory of any hook follow the rename in the same change. Specs, tickets and research notes already published under `.agents/workflow/` are history and keep their original wording.
- **The approved suffix list gains `guide`.** The Markdown detector in the quality runner is the enforcement of rule 2; its list is extended, and the rule's own examples may be extended by the owner (Further Notes).
- **`drift` resolves its root without `../`.** The dispatcher already knows its own absolute directory; the project root is derived from it by `dirname`, three levels up, so the path never contains `..` and the result is the same absolute directory from any working directory and through the symlink `pnpm install` creates. The pre-commit command calls the dispatcher through its absolute `DRIFT_CLI_DIR` instead of climbing to the launcher.
- **The console proxy test moves beside its owner.** The matcher test becomes a sibling of `proxy.ts` and reads `./proxy.ts`; the console's test glob already covers that location. Its behaviour (the matcher accepts application routes and excludes API, `_next` and static assets) is unchanged.
- **Markdown links are root-absolute.** The two Nest READMEs link `/.claude/rules/11-nestjs.rule.md`. GitHub resolves a leading slash against the repository root, and the `ex` collection code treats a leading slash the same way, so the link works in both readers.
- **Go files take kebab-case, keeping the toolchain's suffixes.** `cluster_create.go` becomes `cluster-create.go`, `local_source_test.go` becomes `local-source_test.go`; `_test.go`, `_linux.go` and the like stay. The Go toolchain accepts hyphens in file names, and no `go:embed`, comment or document names a Go file, so the change is `git mv` plus a green build. The rename covers what remains after spec 0035's splits, which create their new files in kebab-case directly.
- **Detectors.** The `../` detector for TypeScript and JavaScript looks at every quoted `../` string, tests included, not only at imports; the shell detector's roots widen to `.claude`, `apps`, `lib` and `packages`; a new detector flags a `../` inside a Markdown link; the Markdown-suffix detector accepts `guide`. The kebab-case and bucket detectors stay as they are.
- **No product change.** The `ex` record files under `lib/cli/src/testdata/` keep their names: they are the product's record format, and testdata is outside the detectors' scope by design.

## Testing Decisions

- **Seam: `drift quality check filesystem`.** The rule's detectors are the verification; the target is zero findings once the three tickets land, and each ticket leaves only the findings of the tickets after it. A widened or new detector is proven by running it against the tree before and after the fix it guards.
- **Behaviour stays green.** `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test`, `drift check localization` and `drift check biome` pass after every ticket. The Go rename is verified by `go vet` and `go test -race` in `lib/cli`, `lib/core/contracts/go` and `apps/microservices/console-stubs`, which the Turborepo tasks run.
- **`drift` from any directory.** After the path fix, `drift --help` and `drift quality check filesystem` run from the repository root, from a nested directory, and through the symlink `pnpm install` places in the bin directory, with the same result.
- **Links resolve.** After the rename, every path cited in `AGENTS.md` and in `.claude/skills/**` exists; a grep for the old names over the live tree (outside `.agents/workflow/`) returns nothing.
- **Review-enforced.** Whether a chosen suffix names the document's kind, and whether a link target is the right document, stay with the review.

## Out of Scope

- The `lib/` directories of the console and the website (spec 0041), and the multi-function `helpers.go` and `common.sh` (spec 0035).
- Splitting Go files: spec 0035 does it, and this spec's Go rename waits for it.
- The names of `ex`'s record files and of anything under `testdata/`.
- Rewriting the wording of published specs, tickets and research notes.
- Renaming the `drift` launcher and the `ex` launcher to `.sh` (spec 0037).
- Any change to `.claude/rules/`.

## Further Notes

- **Suffix call.** `.guide.md` is the chosen kind for the handbook and the skill references. The owner may prefer `.reference.md` for the skill files, or add `.guide.md` to the examples in rule 2; the detector list is the enforcement either way.
- **`dirname` instead of `../`.** Deriving the root with `dirname` three times is read as compliant: the rule bans paths that climb, and the derived value is an absolute path. If the owner reads `dirname` climbing as the same thing, the alternative is `git rev-parse --show-toplevel`, which fails in a checkout that is not yet a Git repository, a case `drift setup` supports today.
- **Historical references.** Published specs and tickets that name `.agents/context.md` are left as written; only live documents, skills and scripts move.
- **Go hyphens.** Hyphenated Go file names are accepted by `go build`, `go vet` and `go test`; only `_test.go` and the build-constraint suffixes carry meaning, and the kebab-case detector already exempts them.
- **Proper names in paths.** Nothing in this spec touches a Next.js route segment (`[locale]`, `[...path]`), which the detector already exempts.
- **`packages/tokens/shared/`** is not a bucket: "shared token" is a glossary term, and the folder holds the tokens every theme shares.

> **Note, 2026-11-05.** `drift` moved from `platform/cli` to `packages/drift` (ticket 0309) at the same depth, so the `../../..` root climb this ticket removes is unchanged; read its paths under `packages/drift/`.
