# The packages and the drift CLI comply with the general rules

**Status:** ready-for-agent

## Tickets

- [x] 0245: the general rules name the Go test idioms they allow
- [ ] 0246: drift measures one unit per Go file, nesting, Go tests and shell functions
- [ ] 0247: the Microruntimes contracts declare one unit per file
- [ ] 0248: the microruntime domain declares one unit per file
- [ ] 0249: the other domain packages declare one unit per file
- [ ] 0250: the application root declares one unit per file
- [ ] 0251: the execution package declares one unit per file
- [ ] 0252: the planner declares one unit per file
- [ ] 0253: inputs, session, setup, gates and the wiring declare one unit per file
- [ ] 0254: the connected agent's supervision loop is a type
- [ ] 0255: the native runtime splits by concern
- [ ] 0256: the infrastructure packages files, process, toolchain, kind, kubectl and terraform declare one unit per file
- [ ] 0257: the infrastructure packages configuration, manifest, collection, inputs and scaffold declare one unit per file
- [ ] 0258: the infrastructure packages records, history, storage, rag and stages declare one unit per file
- [ ] 0259: the infrastructure packages lifecycle, delivery, hooks and judge declare one unit per file
- [ ] 0260: the CLI's presentation declares one unit per file and documents its environment
- [ ] 0261: the streaming transport and the web components apply the general rules
- [ ] 0262: drift's core shell declares one function per file
- [ ] 0263: drift's commands declare one function per file

## Problem Statement

`.claude/rules/00-general.rule.md` applies to every file in the repository, in every language. The read-only audit of 2026-11-04 measured `lib/`, `packages/` and `platform/cli` at 83.8% against it (about 86% without the Go tests), against 94% for the apps. The violations concentrate in the Go code of `lib/cli` and `lib/core/contracts/go`, in the shell of `platform/cli` and, to a smaller degree, in the TypeScript of `packages/streaming` and `packages/web`:

- **Single responsibility.** 109 of the 193 Go source files and 103 of the 124 Go test files declare more than one top-level function or more than one type with its methods, for instance `lib/core/contracts/go/configuration/lexical.go` (7 functions), `lib/cli/src/application/inputs/service.go` (`Service` and 12 free functions), `lib/cli/src/application/plan/helpers.go` (22 functions, under a bucket name), `lib/cli/src/infrastructure/records/format.go` (17) and `lib/cli/src/wiring.go` (7). `lib/cli/src/infrastructure/connected/native.go` runs 916 lines and mixes materialization, lifecycle, migration and the build tree; `lib/cli/src/application/identity.go` holds `IDSource` and `Clock`. 29 files are not named after what they declare: `application/execution/confirm.go` declares `Service`, `domain/finops/ledger.go` declares `SettleTotals`, `packages/streaming/src/transport/redis-streams.strategy.ts` declares `RedisStreamsServer`. In the shell, `platform/cli/src/core/output.sh` declares 18 functions, `core/common.sh` 12, `commands/check-biome.sh` 17 and `commands/quality-dead-code.sh` 11.
- **Control flow.** `runOne` in `lib/cli/src/application/connected/service.go` is 368 lines at nesting depth 5, built from closures over a dozen local variables. `preserveNative` in `native.go` reaches depth 3 and carries an `else` after a `return`; so do `infrastructure/files/json.go`, `infrastructure/scaffold/publish.go` and `infrastructure/storage/projects.go`. 70 of the 248 Go conditions with three or more operands are unnamed: `domain/microruntime/cleanup.go` combines 13 operands in one `if`, `lib/core/contracts/go/configuration/configuration.go` has nine such conditions in `Validate`, and `packages/web/src/ui/hero-shader/read-palette.ts` has one.
- **Statement spacing.** 101 Go source returns, 3 shell returns and 7 in `packages/tokens/build.js` follow another statement without a blank line; 61 Go functions of ten lines or more have no blank line at all between their phases, `runOne` and the 131-line function in `application/connected/project.go` among them.
- **Multi-field literals.** 370 of the 569 multi-field literals in Go source sit on one line, such as the 11-field `Event` literal inside `runOne`, the struct in `application/cluster.go` and the map in `application/execution/finish.go`. The Go tests hold 1,056 more, most of them the rows of test tables.
- **Readable names.** Go source uses `c`, `e`, `m`, `i` in `native.go` and `b` in `connected/service.go`; the tests use `t` 719 times, as the toolchain's convention wants, and other one-letter names such as the receiver `w` in `application/cluster_create_test.go`. 13 results are discarded with `_` before being needed: `cpu, _ := ParseQuantity(...)` in `native.go`, `data, _ := json.MarshalIndent(...)` in `lib/core/contracts/go/configuration/configuration.go`, `strconv.Atoi` in `infrastructure/records/mirrors.go`.
- **External input and errors.** `dft` reads `RUNTIME_SECRET`, `PROJECT_SECRET`, `BASE_SHA`, `HEAD_SHA`, `GITHUB_STEP_SUMMARY` and `ANTHROPIC_BASE_URL`, and `lib/cli/README.md` documents only `ANTHROPIC_API_KEY`. 39 errors are discarded: `_ = command.MarkFlagRequired(...)` seven times in `presentation/cli/catalog/inputs.go`, the unlock and close in `infrastructure/files/lock.go`, `_ = save()` in `runOne`; `packages/web/src/ui/agent-orb/agent-orb.tsx` and `hero-shader/hero-shader.tsx` end a dynamic import with `.catch(() => {})`. `packages/streaming/src/transport/redis-streams.strategy.ts` wraps the pure `parseStreamEntry` in `try/catch`.

The rule file itself has one gap the owner closed on 2026-11-04: idiomatic Go tests. The `t *testing.T` parameter and the one-line rows of a test table are what every Go reader expects, and the rule's Readable names and Multi-field literals sections forbid them as written. The owner approved an amendment that exempts exactly those two forms in `_test.go` files; every other rule keeps applying to Go tests, and spec 0027's decision to keep tests out of the one-function-per-file rule is withdrawn, since rules outrank specs.

`drift quality check general` measures part of this today. Its Go detector for one function per file counts only free functions, so a free function beside a type's methods, or two receiver types in one file, pass; nothing measures nesting depth, shell functions per file or the Go tests at all.

## Solution

Bring every file under `lib/` (the Go of `dft` in `lib/cli` and of the contracts in `lib/core/contracts/go`), `packages/` (the TypeScript of `streaming`, `web` and `tokens`) and `platform/cli` to full compliance with `00-general`, with no change in behavior: every Go test, vitest suite, `go vet`, Biome check and `drift` command stays green, and the LangDrift CLI, the events transport and the web components behave as before.

First, the rule gets its approved amendment and the detectors learn the same exemption and the violation classes they miss today. Then the Go code is split one unit per file, layer by layer, with `runOne` and `native.go` decomposed by hand where their difficulty sits; the Go tests follow their layers; the TypeScript packages close their few findings; the drift CLI's shell is laid out one function per file and sourced by name. At the end `drift quality check general` reports zero findings under `lib/`, `packages/` and `platform/`, and the classes no detector can catch are listed for review.

## User Stories

1. As the owner, I want `drift quality check general` to report zero findings under `lib/`, `packages/` and `platform/`, so that the general rules are measured, not assumed, outside the apps.
2. As the owner, I want `00-general.rule.md` to say in its own words that `t *testing.T` and table-test rows are idiomatic Go, so that a reviewer and a detector read the same exemption and nothing else is exempted by analogy.
3. As a maintainer, I want every Go source file to declare one function or one type with its methods, named after it, so that I find `SettleTotals` in `settle-totals.go` and nothing else there.
4. As a maintainer, I want every Go test file to hold one test with its subtests, and every fake or helper in its own file, so that a failing test names the file that holds it.
5. As a maintainer, I want the connected agent's supervision loop to be a type whose methods I can read one at a time, so that I can debug a lost lease or a revoked intent without holding 368 lines in my head.
6. As a maintainer, I want the native runtime's materialization, start, observation, migration and build tree in separate files, so that a change to migration cannot touch observation by accident.
7. As a maintainer, I want every condition with several business rules to carry a name, so that `isOwnedSelectedPod` reads as a list of named facts instead of 13 operands.
8. As a reader, I want a blank line before every return and between the phases of a function, in Go, TypeScript and shell alike, so that the shape of a function is visible before its text.
9. As a reader, I want every multi-field literal on several lines, so that a diff to one field shows one line.
10. As a reader, I want no one-letter name outside the Go test parameter, so that `cluster`, `event` and `manifest` say what `c`, `e` and `m` did not.
11. As a `dft` user, I want `lib/cli/README.md` to list every environment variable the CLI reads, by key, so that I know what to set before `dft agent connect` or `dft gates report`.
12. As an operator, I want no error discarded in `dft`, so that a lock that fails to release or a flag that fails to register is reported instead of silent.
13. As a microservice maintainer, I want the events transport to park an unreadable entry through a result it inspects, not an exception it catches, so that the boundary between the stream and the parser is explicit.
14. As a website visitor, I want the agent orb and the hero shader to keep their fallback when their runtime fails to load, and as a maintainer I want that fallback to be a named decision in the component, not an empty catch.
15. As a drift contributor, I want each shell function in its own file, sourced by name, so that I can read `drift_print_error` without scrolling past seventeen others.
16. As the owner, I want every file this spec creates to be kebab-case, so that spec 0036's rename of the existing Go files is the only rename left.
17. As a reviewer, I want the spec to say which violation classes no detector catches, so that the two-axis review knows what it alone enforces.
18. As an agent implementing a ticket, I want each ticket to fit one fresh context window and land green, so that the refactoring never leaves the tree broken between tickets.

## Implementation Decisions

### The rule amendment

`00-general.rule.md` gains, under Readable names, a second sentence in item 1: "In a Go test, the `t *testing.T` parameter keeps the name the toolchain's convention gives it." Under Multi-field literals, item 1 gains: "In a Go test, each case of a test table may stay on one line, so the table reads as rows; a literal outside a table follows the rule." Nothing else changes: receivers, loop variables and helpers in tests keep needing full names, test files keep one function each, and the spacing rules apply to tests as to source. The amendment is the first ticket, so every later ticket and detector reads the rule as amended.

### The detectors

`drift quality check` stays the verification seam, through `rule-detectors.tsv` and awk programs under `platform/cli/src/quality/`, POSIX shell and awk only. The spec adds or changes these detectors for `00-general`:

- **One unit per Go file** becomes an awk program: a file is reported when it declares more than one free function, a free function beside a method, or methods on more than one receiver type. It replaces the per-file count of `^func`.
- **Conditional nesting in Go**: an awk program reports a file where an `if`, `for`, `switch` or `select` opens at three or more tab stops inside a top-level function, which is a third conditional level under gofmt.
- **More than one function in a shell file**: a per-file detector over `*.sh` under `platform` and `packages`, on the `name() {` declaration form.
- **One-letter names in Go tests**: the existing Go one-letter detector gains a tests-scoped twin whose pattern excludes `t`.
- **Multi-field literals in Go tests**: an awk program reports a one-line keyed literal in a `_test.go` file unless it is an element of an enclosing slice or map literal, which is a table row.
- **One test per Go test file**: a per-file detector counting `^func Test`.
- **Discarded error in Go** keeps its advisory wording and extends to the tests.

Each awk detector ships with fixtures under `platform/cli/src/quality/testdata/<detector>/`, a file set it must report and one it must not, and a `drift quality detectors` self-check runs every awk program over its fixtures and fails when the reported set differs from the expected one, in the manner of `check-biome.sh`'s fixture-driven self-check. `testdata/` is already outside every code-scoped detector's paths, so the fixtures never count as findings. Detection stays review-enforced, and the spec says so, for: named compound predicates, phase spacing inside a function, the file name matching its declaration, nesting that gofmt does not express in tab stops, and the handling of an error at a boundary.

### Go: one unit per file

- The strict reading of Single responsibility 2 applies: a file declares exactly one free function, or exactly one type with its methods. A constructor is a function and moves to its own file beside its type (`new-service.go` declares `NewService`). A second receiver type moves with its methods to its own file. Interfaces and plain types with no methods may share a `ports`-style file, as the rule allows types that exist for one concern.
- A package keeps its import path and its exported identifiers; the splits move declarations between files of the same package, so no caller changes. No new Go package is created to house a split, and `helpers.go` dissolves into files named after each function's concern, never into a `helpers` folder.
- A new file is named after its declaration in kebab-case: `settle-totals.go`, `markdown-hash.go`, `fake-projects_test.go`. The Go toolchain requires only the `_test.go` suffix and build-constraint suffixes; the rest of the name follows the filesystem rules. Existing snake_case files that keep their content are renamed by spec 0036 after these splits.
- The split goes layer by layer: `contracts`; `domain/microruntime`; the other domain packages; the application root; `application/execution`; `application/plan`; `application/inputs`, `session`, `setup`, `gates` and the wiring; `application/connected`; `infrastructure/connected`; the remaining infrastructure packages in four batches; `presentation/cli`. Each ticket leaves `go build`, `go vet` and `go test -race` green on its own.

### Go: the two decompositions

- **The supervision loop.** `runOne` becomes a supervision type in `application/connected`, one file per method: the state it carried in local variables (the active operation, its cancel, the lease, the last contact, the stop cause, the intent) becomes the type's fields; `save`, `appendEvent`, `launch`, `shutdown` become methods; the `select` loop dispatches to one method per channel (a worker report, a worker result, a console tick, the context's end), and the console tick's reconciliation of the sync response against the record is its own method. The owner writes the reconciliation. Behavior is pinned by the existing supervisor scenarios and service tests before the split.
- **The native runtime.** `native.go` splits by concern: verification of the repository, materialization, the snapshot, start, stop, observation, migration with its source lookup and intent files, the build tree and its copy, and the managed-state check, each in its own file with `Native` keeping its methods' names and signatures. `preserveNative` is flattened to guards and named steps (one unit per workload object, a service port rewrite, a workload rewrite) with no `else` after an exit and at most two conditional levels; the owner writes it.

### Go: control flow, spacing, literals, names

- Every `else` after `return`, `continue` or `break` is removed; every function at three or more conditional levels is restructured with guards, named functions or a `switch`.
- Every condition with three or more operands that combines business or runtime rules becomes a named predicate: a boolean method or function beside the code that uses it, with its own test when the rule is a domain one, such as whether a pod is an owned, selected pod of an execution, or whether a configuration's envelope allocates enough per component.
- A blank line precedes every `return` that follows a statement, separates the phases of a function and never touches a block boundary. The two spacing awk programs already cover Go.
- Every multi-field struct and map literal in Go source is written one field per line, including the `Event` built inside the supervision loop.
- One-letter names are replaced by the thing they hold, in source and tests, except `t`.
- A discarded result is named and handled: a parse that cannot fail by construction is replaced by a constructor that does not return an error, or its error is returned; a discarded `MarkFlagRequired` error is an internal invariant and panics with the flag's name through one function; the lock's release returns its error and callers join it into theirs; a test that discards a marshal error fails through `t.Fatal`.

### Go: external input

- `lib/cli/README.md` documents every environment variable `dft` reads, by key and purpose, never by value: `ANTHROPIC_API_KEY`, `ANTHROPIC_BASE_URL`, `RUNTIME_SECRET`, `PROJECT_SECRET`, `BASE_SHA`, `HEAD_SHA`, `GITHUB_STEP_SUMMARY`. Each read stays at the operation that uses it, and a missing value fails there.

### The TypeScript packages

- `packages/streaming`: `parseStreamEntry` returns a parse result that is either the entry or the rejection with its reason, and `deliver` branches on it; no `try/catch` wraps the parser. The strategy file is renamed after the class it declares, `RedisStreamsServer`, and the transport barrel follows. Spacing and literal findings in the package close.
- `packages/web`: the agent orb and the hero shader replace `.catch(() => {})` with a named handler that applies the product-approved fallback explicitly, recording the failure in component state that the markup exposes, so a test can assert the fallback. The unnamed predicate in the palette reader gets a name. Spec 0030's remaining spacing findings close.
- `packages/tokens`: `build.js` is replaced by shell under spec 0037; this spec does not edit it, and its spacing findings vanish with it.

### The drift CLI's shell

- One function per file. The functions of `core/common.sh` and `core/output.sh` move to `platform/cli/src/core/<function>.sh`, one each, named after the function without its `drift_` prefix (`core/print-error.sh` declares `drift_print_error`); `core/commit-rules.sh` likewise. `core/common.sh` and `core/output.sh` disappear; a `core/load.sh` with no function of its own sources every sibling, and every command sources it where it sourced `common.sh`. Shell `source` of per-function files is the shell analog of the Go split.
- A command's local functions move to a folder named after the command (`commands/check-biome/<function>.sh`), one each, and the command sources them by name. `drift.sh`'s `drift_usage_error` moves to `core/usage-error.sh`.
- `drift quality dead-code` keeps tracking shell functions across the new layout, and `exceptions/dead-code.allow` is updated only when a function is really unused.
- A `return` or `exit` that follows a statement gets its blank line. The `|| true` idiom on a probe (a command run for its exit status) is a handled result, not a swallowed error, and stays; an `|| true` that hides a failure is removed.
- `platform/cli/drift`'s missing `.sh`, `lib/cli/bin/ex`, `lib/cli/build.sh` and the `../` in `drift.sh` belong to specs 0037 and 0036.

## Testing Decisions

- **The detectors are the rules seam.** `platform/cli/drift quality check general` reports zero findings under `lib/`, `packages/` and `platform/` at the end, and each ticket runs it on its own paths. A new awk detector ships with fixtures under `platform/cli/src/quality/testdata/<detector>/` that it must report and must not report, and `platform/cli/drift quality detectors` holds them, in the manner of `check-biome.sh`'s fixture-driven self-check; a detector is reliable when its fixtures hold, and a class it cannot catch is listed as review-enforced above.
- **The existing suites are the behavior seam.** `pnpm turbo run test typecheck check --filter=@langdrift/cli --filter=@langdrift/core --filter=streaming --filter=web --filter=tokens` passes after every ticket (`go test -race ./...` and `go vet ./...` for Go, vitest for events and web, `node --test` for tokens), and `drift doctor`, `drift check biome`, `drift check localization`, `drift quality check`, `drift quality coverage` and `drift quality dead-code` run as before after each shell ticket. A good test asserts external behavior: what `dft` writes and prints, what the supervisor does to a record for a given sync response, what the transport parks; a test that only mirrors a file's internals moves with the code it covers.
- **The two decompositions are pinned before they move.** The supervisor scenarios and service tests in `application/connected`, and the native tests in `infrastructure/connected`, are run and kept green at every step of the split; where a scenario is missing for a branch the split touches (a lease lost with an active stop, a purge epoch advance, a migration with a Job unit), it is added before the split, at the service's public methods, not on the new internal methods.
- **The two-axis `code-review`** runs before every commit, against this spec and every rule file, and checks the review-enforced classes on every change.
- **Prior art:** the Go table tests across `lib/cli`, the awk detectors with `drift quality coverage`, the vitest suites of `packages/streaming` and `packages/web`, and the fixture-driven self-checks of `check-biome.sh`.

## Out of Scope

- `apps/`: spec 0034 owns the general rules there, including every `route.ts` through spec 0041.
- Renaming the existing snake_case Go files and the `../` paths in `drift.sh` and `git-pre-commit.sh`: spec 0036, which waits on this spec's splits.
- `platform/cli/drift` and `lib/cli/bin/ex` without `.sh`, `lib/cli/build.sh`, `packages/tokens/build.js`, the commit-message hook and the pre-commit's localization check: spec 0037.
- Test descriptions and bare `test()`: the testing compliance spec, not yet written. The Nest microservices: the Nest compliance spec, not yet written.
- Whether `drift quality check` fails on findings: today it prints and exits 0; see Further Notes.
- Any change to what `dft`, the transport or the components do.

## Further Notes

**Coordination with spec 0045.** Spec 0045's tickets (0222 to 0227) rewrite parts of `lib/cli` and `lib/core/contracts/go`: the judge adapter, the planner's judge calls, the run loop and its phases, `dft agent stop`, and the removal of `ex agent init` and its scaffold. The Go tickets of this spec wait for them, so no file is split while another branch rewrites it, and each Go ticket re-measures its layer before it starts. The judge transport's `ANTHROPIC_BASE_URL` and `ANTHROPIC_API_KEY` reads leave `dft` with spec 0045, so the README lists what `dft` reads when the ticket runs.


Strict readings taken, for the owner to overturn at pull request review:

1. **Constructors.** `NewX` is a function, so it lives in `new-x.go` beside `x.go`. The owner may permit a constructor beside its type; the detector would then allow one free function whose name starts with `New` beside the methods.
2. **Go tests, one function per file.** A `_test.go` file declares one `Test` function; cases that exercise one behavior become `t.Run` subtests of it, a fake is one type with its methods per file, and a helper is one function per file. The decision-2 exemption covers only `t` and table rows. Spec 0027 kept tests out of this rule; rules outrank specs, so that is withdrawn.
3. **Kebab-case test files.** `fake-projects_test.go` is the shape the two rules produce together; the owner may prefer that spec 0036 settles a convention before the test tickets run.
4. **`MarkFlagRequired`.** Its error means the flag was never registered, an internal invariant, so one function panics with the flag's name instead of discarding it.
5. **`ANTHROPIC_BASE_URL`.** The judge transport falls back to Anthropic's public endpoint when the variable is unset. The spec keeps that as the API's own default, not an invented value, and documents the key; the owner may require the variable instead.
6. **Shell `|| true`.** Kept on probes, removed where it hides a failure, decided per site in the shell tickets.
7. **Interfaces together.** A `ports.go` holding several interfaces of one concern stays one file; the rule's "one type with its methods" is read as applying to types that have methods.
8. **One receiver per file, not one file per receiver.** The methods of one type may span several files of its package when each file holds methods of that type only, grouped by concern (`Execution`'s cleanup methods stay in their own file beside its other methods). A file never mixes two receiver types or a receiver with a free function; the one-unit detector reads the rule this way. The owner may require every method beside its type's declaration instead.

Open questions:

- `drift quality check` never fails (it prints the findings and exits 0), so zero findings are a target the review verifies, not a gate. Making it fail on findings belongs to the mechanical checks of spec 0037; the owner should confirm that placement.
- The Go one-letter detector cannot tell a `t` that is not `*testing.T` in a test file; the tests-scoped twin exempts every `t` in `_test.go`, which is the owner's chosen scope.
