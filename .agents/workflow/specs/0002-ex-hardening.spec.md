# Ex CLI hardening

**Status:** ready-for-agent

## Problem Statement

The ex CLI arrived from the Exphere repository as a pristine copy in `packages/ex`. Its lifecycle engineering is strong (ledger-first execution, UID-guarded cleanup, 93–100% coverage in the core layers), but it is not yet a LangDrift package: the branding, module identity and build wiring are Exphere's, `pnpm install` would fail on machines without Go, and the review found real gaps — the console can run arbitrary commands on the user's machine, subprocesses inherit `PROJECT_SECRET`, failures swallow their causes, the connected agent is the least-tested code while being the most privileged, aggregates accept unguarded mutation, and the e2e suite never runs in CI.

## Solution

Adopt ex as `@langdrift/ex`, then harden it in place, phase by phase: close the console trust boundary with an allowlist and minimal subprocess environments, make every failure explain itself (stderr tails, evidence files, a six-code exit contract, human output where humans read), bring the connected agent to the same rigor as the native lifecycle (typed states, locks, crash-consistent stores), encapsulate the domain aggregates, and run the e2e suite on a schedule in CI. The FinOps engine itself (metrics, workflow executor, RAG, episodic memory, LLM-effort) is explicitly out of scope.

## User Stories

1. As a frontend contributor, I want `pnpm install` to succeed without a Go toolchain, so that the monorepo does not force Go on me.
2. As a maintainer, I want ex to build and test through `pnpm turbo run build/test/typecheck`, so that CI covers it with no extra wiring.
3. As a maintainer, I want PR CI to also run `go vet -tags e2e`, the race detector and a pinned `govulncheck`, so that cheap classes of defect never reach review.
4. As a user, I want the CLI, its files and its messages to say LangDrift (`@langdrift/ex`, `.langdrift/`, `LANGDRIFT_SKIP_EX_BUILD`), so that the product is coherent; the binary stays `ex`.
5. As a user, I want a console that sends a build recipe outside the policy (committed project script, or `docker build`/`go`/`npm`/`pnpm`/`make` under the argument constraints) to be refused with a message naming the policy, so that the console cannot run arbitrary code on my machine.
6. As a user, I want every interpreter to refuse a flag-style first argument, so that `-c`/`-e`/`--eval` escapes are closed.
7. As a user, I want every subprocess to run with a minimal explicit environment, so that `PROJECT_SECRET` and my ambient credentials never leak into tools that did not declare them.
8. As a project author, I want hooks to receive exactly the variables my hook configuration declares, so that what a hook can see is reviewable.
9. As a maintainer, I want every git invocation to run with `GIT_*` unset, so that an ambient `GIT_DIR` cannot redirect repository checks.
10. As a user, I want a failed build to tell me the failing command and a bounded tail of its stderr, so that I can fix it without re-running by hand.
11. As a user, I want long build output streamed to an evidence file instead of failing the build at 1 MiB, so that verbose builds work and remain auditable.
12. As a user, I want exit codes to distinguish success (0), operation failure (1), usage (2), policy refusal (3), missing prerequisite (4) and interruption (130), so that scripts can react without parsing messages.
13. As a user, I want the agent commands' argument errors to exit 2 like every other command, so that the contract is uniform.
14. As a user, I want `doctor`, `cluster inspect`, `runtime inspect` and `version` to print human-readable output by default, so that I do not read indented JSON; `--json` keeps the machine form, and every other command keeps JSON output.
15. As a user, I want errors rendered as one line of cause plus a hint, so that the first line tells me what to do.
16. As a maintainer, I want the connected agent's run loop modeled as a typed state machine in the domain with table-driven tests, so that the most privileged code is the best-tested code.
17. As a user, I want `ex agent migrate` to take the project lock, so that it cannot race the supervisor.
18. As a user, I want a failed bind or materialize to be repairable on the next attempt, so that one crash does not wedge the runtime forever.
19. As a maintainer, I want `kubectl` port-forward streams to stop sharing one writer across goroutines without synchronization, so that diagnostics cannot race.
20. As a maintainer, I want MCP tests covering 401 (`ErrAccessRevoked`), the response size cap and redirect refusal, so that the console client's guarantees are proven.
21. As a maintainer, I want the `connected` application and infrastructure packages at ninety percent coverage or above, achieved rather than CI-gated.
22. As a maintainer, I want `Execution` and `cluster.Record` to guard their transitions and hide mutable state behind methods, so that a caller typo cannot corrupt a ledger.
23. As a maintainer, I want the port-range and source-location validations defined once, so that the three copies cannot drift.
24. As a maintainer, I want `domain/connection` to classify errors with sentinels like `microruntime` does, so that callers can match them.
25. As a maintainer, I want time and ID generation behind `Clock` and `IDSource` ports, so that tests control them and the hidden dependencies disappear.
26. As a user, I want `ex runtime …` to keep refusing with migration guidance while the dead command tree and its duplicated test harness are deleted, so that the live surface is the tested surface.
27. As a maintainer, I want a separate e2e workflow (manual dispatch plus a weekly schedule) running the `e2e`-tagged suite with a pinned kind on ubuntu, so that the lifecycle is exercised against a real cluster without slowing PR CI.
28. As the owner, I want every commit that touches `packages/` to carry its `Ticket:` trailer and no AI signature, so that the repository's guardrails hold.

## Implementation Decisions

- **Adoption order.** The pristine Exphere copy is committed first, unchanged, so every later diff is reviewable against it.
- **Identity.** npm `@langdrift/ex`; Go module `github.com/jonatassales/tao/packages/ex`; project directory `.langdrift/`; `LANGDRIFT_SKIP_EX_BUILD`; user-facing strings say LangDrift. `PROJECT_SECRET` and `RUNTIME_SECRET` keep their names. The binary stays `ex` (`drift` is the platform CLI).
- **Build wiring.** No npm `postinstall`. Turborepo tasks call the package scripts; the `test` script runs with the race detector. `build.sh` remains for standalone use, rebranded. PR CI gains `go vet -tags e2e` and pinned `govulncheck` in the Harness job.
- **Trust boundary.** Recipe validation lives in the domain (`connection`): shape one, a committed project script resolved inside the checkout; shape two, an allowlisted tool with the existing argument constraints; flag-style first arguments refused for all interpreters. Subprocess environments are built by one allowlist helper per tool family in `infrastructure/process`; hooks merge in only declared variables; git runs with `GIT_*` scrubbed everywhere.
- **Diagnosable failures.** `process` gains a spill-to-file capture mode used by build recipes; the evidence directory holds the full output; errors embed a bounded tail. Exit codes map from sentinels in one place in `presentation/cli`; `ErrAdoptionRefused` and policy refusals exit 3, missing tools and absent prerequisites exit 4. Human formatting is a presentation-layer renderer over the same data the JSON path emits.
- **Connected agent.** The `runOne` closure state machine becomes a typed domain model (states and causes as declared types with guarded transitions), the service keeps orchestration only. `Migrate` locks the project. `Store.Bind` writes the record before the secret file and repairs an orphaned `.env`; `Materialize` records intent first so a partial write is adopted or cleaned, never a permanent conflict. `Forward` gives each stream its own synchronized writer.
- **Domain polish.** Aggregate fields unexported where transitions exist; `RecordCreated` bounds-checked; shared validations consolidated into the owning domain type; `connection` sentinels; `Clock` and `IDSource` ports wired from the composition root.
- **Out of band.** No new third-party dependencies anywhere in the module.

## Testing Decisions

- The existing suites are the safety net; every change lands with its table-driven tests at the same seams the package already uses (domain and application tests against fakes, adapter tests against recorded runners, `httptest`-style transport fakes for MCP).
- New behavior contracts get tests at the highest existing seam: exit codes and human output through `cli.Execute` with injected streams; recipe policy through the domain validator plus one application-level refusal test; environment minimization through the recorded runner's observed env.
- The e2e-tagged suite is not extended in this work; it runs as-is in the new workflow.
- A good test asserts what a consumer observes (exit code, stderr text, saved ledger, refused recipe), never internal call order except where ordering is the contract (ledger-first writes).

## Out of Scope

- The FinOps engine: metrics ingestion, workflow executor, RAG, episodic memory, LLM-effort suggestion.
- Any console/MCP server implementation.
- Coverage ratchets or lint tooling beyond vet, race, govulncheck.
- Windows support; publishing the package or binary anywhere.

## Further Notes

- ADR 0002 records the trust-boundary policy. Glossary: microruntime, connected agent, build recipe, execution ledger, evidence.
- Evidence basis: full test suite, vet (both tag sets), gofmt, per-package coverage, and two layer reviews on 2026-10-21. E2E suite and real cluster operations were not executed during the review.
