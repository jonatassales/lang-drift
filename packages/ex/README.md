# LangDrift CLI

This Go module, adopted from the Exphere repository, performs the Microruntimes session inside a project. `ex agent plan` reads a ticket and the artifacts it links, judges them through a model, retrieves the project's earlier runs, sizes the components inside the team's envelope and writes the runtime file that plans the work. `ex agent run` runs that plan as laps under the team's money gates: it drives the owner's coding agent through the delivery stages, brings the runtime up on the bounded local kind/Terraform/Kubernetes lifecycle, observes what the harness delivered, settles the ledger and remembers what it measured. The module also keeps the connected agent that supervises microruntimes from a LangDrift console, which needs a server this repository does not have. Applications need no LangDrift dependency. This is a local development candidate, not a published release or a claim of operational readiness.

## Build and checks

Use an existing Go 1.26 or newer toolchain on macOS or Linux:

```sh
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex --help
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
```

The module, `github.com/jonatassales/tao/packages/ex`, is the private package
`@langdrift/ex`. It imports the [Microruntimes contracts](/packages/contracts/README.md),
`github.com/jonatassales/tao/packages/contracts`, for the `runtime.md` fence,
the lexical rules and the agent-server protocol, and the official Anthropic Go
SDK for the planner's model calls. The root `go.work` ties the two modules and
`apps/extub` together, so the CLI imports the contracts without a `replace`
directive; building it outside this repository needs the contracts module
published. This does not imply a published release. Native Go commands own the
build and checks. Building does not install infrastructure tools, alter shell
profiles, or change PATH. Use `./dist/ex` to avoid invoking the unrelated system
editor also named `ex`.

## Workspace package

The package holds only the Go module, its source, and `build.sh`. `build.sh`
builds `dist/ex` with an existing Go toolchain, never installs Go or tools, and
never changes shell profiles; `LANGDRIFT_SKIP_EX_BUILD=1` skips it. Installing
the workspace does not run it. From the workspace root,
`pnpm --filter @langdrift/ex build` rebuilds the executable.

This package declares no npm bin. Run `packages/ex/dist/ex` directly from this
repository. The separate SDK package and its bridge were removed on October 21, 2026. The former shell-setup command,
package bridge, historical schemas, and consumer-harness fixtures were removed
on October 17, 2026; Git history retains them.

## Boundary

In a project, the CLI writes only `.langdrift/`: the project identity, the mirrors of its inputs under `harness/` and `finops/`, one folder per runtime under `microruntimes/runtimes/`, and its private state under `agent/state/`. It reads the team's artifacts where the ticket links them, never outside the project, and asks a FinOps gate before any paid work: the planner's own model calls and every call the stages make through the coding agent. It never pushes, opens a pull request or merges: delivery belongs to the harness, and the CLI only observes it through `git` and `gh`. Anything it cannot observe stays unknown, with the reason; unknown is never zero.

The connected agent (`ex agent connect`, `ex agent project`, `ex agent supervise` and `ex agent migrate`) needs a server that answers its protocol and holds each runtime's intent. None exists in this repository yet, so nothing here can drive it.

## The `.langdrift` folder

The folder is machine-local and gitignored: `ex` regenerates it, and nothing a team authors lives inside it. Spec 0016 defines the layout below; the reference copy is the sample under `src/testdata/langdrift/sample`, which the test suite pins, beside a frozen copy of the ticket its run planned and the sources that ticket links.

| Path | Scope | Written by |
| --- | --- | --- |
| `config.json` | general | `ex init`: the project identity |
| `agent/state/` | general, private | the CLI: input records, the RAG index, the runs' laps, stage logs and operator records, execution ledgers, Terraform state, build logs, connection records, credentials, locks |
| `harness/` | general | the CLI, as mirrors of its inputs: the hardware envelope and the delivery stages |
| `microruntimes/` | Microruntimes | the engine: the execution rules, one folder per runtime with its `runtime.md` (the plan), `run.json` (the execution), `run.md` (the memory), `laps/` and `evidence/` |
| `finops/` | FinOps | the CLI, as mirrors of its inputs: budgets, gates, prices, forecasts, one ledger per runtime |

Every value that can block or allow a run is CLI input held in the private store. The JSON files mirror it for reading; no gate reads them. A microruntime is the car, a run is the race, and a lap is one pass of the run: an attempt paused at a gate, or the execution that finished.

Existing private state is not migrated: a checkout that still has `.langdrift/state` from an earlier version should delete its cluster with `ex cluster delete` before upgrading, or remove that directory once nothing runs.

## Inputs and mirrors

Each input is recorded with who set it, when, through what, its value and the record it replaced. A mirror opens with a note saying it was written from the record; an edit to a mirror changes nothing a gate reads.

| Input | Command | Mirror |
| --- | --- | --- |
| Envelope | `ex harness envelope set --per-execution-cpu-millis … --per-execution-memory-mib … --default-cpu-millis … --default-memory-mib … [--not-runtime-field gpu]` | `harness/envelope.json` |
| Workflow | `ex harness workflow set --file workflow.json` | `harness/workflow.json` |
| Model | `ex harness model set <model> --effort <low\|medium\|high> [--planner <model>]` | none |
| Driver | `ex harness driver set --file driver.json` | none |
| Budget | `ex finops budget set --period YYYY-MM --period-cents … --per-run-cents … --reserve-margin-percent … [--allocate <ticket>=<cents>]` | `finops/budgets.json` |
| Gates | `ex finops gates set --default` or `--file gates.json` | `finops/gates.json` |
| Prices | `ex finops prices set --file prices.json --slug <slug>` | `finops/prices/NNNN-<slug>.prices.json` |
| Forecasts | `ex finops forecasts set --file forecasts.json` | `finops/forecasts.json` |

The workflow names the stages, their skills, where a person takes part and the committed scripts a stage runs first (`hooks`). The forecasts carry the tokens each effort is expected to use and a `planning` entry for the planner's own calls. Every command takes `--by` for who sets the input.

## The session

`ex agent plan <ticket>` prints the plan session: the ticket and its links, each document read, the guardrail criteria judged, the bounds and the budget, the retrieval over the project's episodes with the candidates kept and dropped at the 0.55 threshold, the model and effort suggested, each component inspected and sized, the cross-checks, the schema validation and the draft. It writes `runtime.md` with the executable `langdrift` fence and a `langdrift-plan` fence that holds the analysis and the plan, a `run.json` with the stages pending, and the runtime's ledger with the planner's own line. Planning the same ticket again, while the runtime is stopped, raises the fence's revision.

`ex agent run <runtime>` prints the run session. Running the runtime authorizes its plan. The intent, the estimate at list prices and the per-run budget print, then the gate: a reservation inside the cap, the ticket's allocation and the period's remainder opens a lap; one above them pauses with nothing paid. The stages then play in order: a stage where a person takes part waits for the owner (the run asks at a terminal, or with `--confirm`; otherwise it exits as waiting and `ex agent confirm <runtime> <stage>` records the yes), each skill runs through the coding agent and prints its one-line summary, FG-04 runs before every paid call, and the stage with `human: owner-only` never runs. After the stages the runtime comes up from its fence and prints ready, the implementation is complete with what the lap observed against its reservation, the branch, the pull request and its checks are observed (`--checks-wait` bounds the wait), the runtime goes down, the unused reservation returns to the budget, and the run says whether it is ready for review or to merge and what comes next. Every change lands in the lap file, `run.json` and the ledger; the memory `run.md` is written at the end and the episode index rebuilt.

Lines print with `agent ›` as the actor. With `--json`, each line is one NDJSON event: `{at, session, kind, id, text, card, field, data}`, where `id` is the website's line id and `card` and `field` name the card the line fills, so a console or a TUI can render the same session.

`ex agent status <runtime>` shows the run record and the open lap. `ex agent reindex` rebuilds the episode index from every runtime's memory.

### The judge

The planner's model is the judge: it keeps or discards each source with its facts, weighs the episodes, and suggests a model and effort. Its calls go through the official Anthropic Go SDK, which resolves the credential each time the judge builds its client, from the first source that yields one: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, the profile `ANTHROPIC_PROFILE` names, workload identity federation (`ANTHROPIC_FEDERATION_RULE_ID`, `ANTHROPIC_ORGANIZATION_ID` and `ANTHROPIC_IDENTITY_TOKEN_FILE`), then the active or `default` profile that `ant auth login` writes under `~/.config/anthropic` (`ANTHROPIC_CONFIG_DIR` or `XDG_CONFIG_HOME` moves it). The base URL is `ANTHROPIC_BASE_URL`, else the profile's `base_url`, else `https://api.anthropic.com`. The model is the one `ex harness model set --planner` names, or the default model. The calls are paid work under the same gates as a lap, against the forecasts' `planning` entry, and land on the ledger's planning line. Answers are structured JSON, checked and never asked again; only the SDK retries a call that failed on the way, at most twice. With `--no-judge`, or without a credential, the plan keeps the deterministic shell: the line says so once, and the judged fields stay `null`.

The judge's HTTP client follows no redirect: the SDK sends an API key as `X-Api-Key`, a header Go carries to the host a redirect names. It sends nothing to a base URL that is not HTTPS, HTTP to a loopback host aside: an `ANTHROPIC_BASE_URL` that breaks the rule, or does not parse, stops the judge before the client is built, and a base URL from any other source is refused before the request leaves. It reads at most 4 MiB of a response, as the console client does. An error status reaches the session's warning and the judge's stop reason as its HTTP code and request id, never as the provider's body, and a call that fails on the way (a redirect, the bound, the connection) as a fixed phrase without the URL. A credential the API refuses (401 or 403) and a refused base URL both make the model unavailable.

### The driver

The stages' skills run through the coding agent the owner configured with `ex harness driver set`:

```json
{
	"executable": "claude",
	"args": ["-p", "{{prompt}}", "--output-format", "json"],
	"passEnv": ["ANTHROPIC_API_KEY"],
	"timeoutSeconds": 1800,
	"usage": { "inputTokens": "usage.input_tokens", "outputTokens": "usage.output_tokens" },
	"summary": "result",
	"unmeasured": "pause"
}
```

The executable runs as a trusted process from the project root, never through a shell. Each argument may carry `{{skill}}`, `{{ticket}}`, `{{run}}`, `{{stage}}`, `{{lap}}` or `{{prompt}}`, filled per argument. The environment is the minimal base, the variables `passEnv` names, and `LANGDRIFT_PROJECT_ROOT`, `LANGDRIFT_RUNTIME_ID`, `LANGDRIFT_LAP`, `LANGDRIFT_STAGE`, `LANGDRIFT_SKILL` and `LANGDRIFT_TICKET`. The output streams to a private log under `agent/state/runs/<runtime>/stages/`; `usage` names the dotted paths of the tokens in the JSON output and `summary` the path of the one-line summary, else the output's last line. A call that reports no usage leaves the spend unknown: FG-04 pauses the next paid call unless `unmeasured` is `allow`. A stage's `hooks` are committed scripts run first through `sh`, inside the project. The evidence of every call, token-free, lands under the runtime's `evidence/laps/<lap>/stages/`.

## Commands

In the commands below, `ex` means this tree's compiled `./dist/ex` executable.

| Command | Behavior |
| --- | --- |
| `ex version` | Print the version; `--version` and `-V` are equivalent. |
| `ex init` | Create or reuse the project identity in `.langdrift/config.json` and add the Microruntimes execution rules at `.langdrift/microruntimes/AGENTS.md`. Existing files remain unchanged; no template, skill or policy text is written. |
| `ex harness envelope\|workflow\|model\|driver set …` | Record one input of the delivery harness and write its mirror. |
| `ex finops budget\|gates\|prices\|forecasts set …` | Record one money rule and write its mirror. |
| `ex agent plan <ticket> [--runtime <id>] [--name …] [--component …]… [--cpu-millis … --memory-mib …] [--max-depth … --max-files … --max-bytes …] [--no-judge] [--by …]` | Read the ticket and the files it links, judge, retrieve, size and write the runtime file that plans it. |
| `ex agent run <runtime> [--model …] [--effort …] [--acknowledge-overspend] [--confirm] [--checks-wait <duration>] [--by …]` | Open a lap, or play on with the open one: the gate, the stages, the runtime, delivery, the settlement and the memory. |
| `ex agent confirm <runtime> <stage> [--by …]` | Confirm a waiting stage as the owner; the next run plays on from it. |
| `ex agent status <runtime>` | Show the run record and the open lap. |
| `ex agent reindex` | Rebuild the episode index from every runtime's memory. |
| `ex agent project <console-url>` | Read `PROJECT_SECRET` and supervise all microruntimes of one project, including those created after connection. The console normally launches this process. |
| `ex agent connect <console-url>` | Read `RUNTIME_SECRET` from the environment, verify the selected GitHub checkout, and save a private environment binding. |
| `ex agent supervise` | Synchronize all bound microruntimes, apply stopped configurations, build/start/observe/stop from authenticated console intent. |
| `ex agent migrate <runtime-id> [--recover]` | Write an explicit import guide for a stopped v2 runtime; `--recover` explicitly invokes ownership-aware legacy cleanup first. Preserve original native files, source and evidence. |
| `ex doctor` | Inspect existing native tools without installing anything. |
| `ex cluster create/inspect/delete` | Explicit local cluster administration, preserving existing ownership checks. |
| `ex cluster load-image <image>` | Explicitly load an existing image into the private project cluster. |
| `ex agent init <name>` | Retain the separate editable finite LangGraph Job scaffold; it is not the connected agent. |

All previous `ex runtime ...` creation, configuration and operation commands now refuse with console/migration guidance. Their native use cases remain internal to the run, the connected agent and test fixtures.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success, or a run that waits for a person. |
| 1 | The operation failed: a skill, a hook, a build or a start did not succeed. |
| 2 | Invalid usage: unknown command, flag, argument shape or input value. |
| 3 | Refused by a guard or policy: a money gate paused or blocked, a plan check failed, a stage is not waiting, a run finished, adoption of unowned resources, a changed configuration, a build recipe outside the allowlist, revoked console access. |
| 4 | A prerequisite is missing: a required tool is not installed, or a required project, input, cluster or runtime record does not exist. |
| 130 | Interrupted. |

Global flags are `--project <directory>`, `--json`, and `--timeout <duration>`.
Without `--project`, the nearest parent project is discovered. Finite commands default to two minutes; `ex agent run` is bounded by the driver's timeout per call and by `--checks-wait`. `ex agent supervise` owns a foreground process until interrupted; each native operation is separately bounded to two minutes and cleanup uses its own bounded context.

The standalone binary can be built directly with Go and placed in a user-selected executable directory as `ex`; Node, pnpm and the SDK are unnecessary in consumer projects. Builds do not modify a shell profile or PATH. No binary download, package publication or personal installation was performed in this delivery.

The local console normally handles bootstrap and agent launch. One project has one `.langdrift` and one agent; a process lock prevents a second project agent or a legacy supervision loop from concurrently owning the same checkout. The agent admits at most 128 runtime workers, polls membership every five seconds, and joins its workers on cancellation. A console-owned stdin pipe also cancels the agent if its parent disappears. This process does not install dependencies or start workloads merely by connecting.

Local project bindings validate the canonical folder instead of requiring a GitHub origin, clean Git tree, or commit. Builds snapshot current tracked and nonignored untracked files, or bounded plain-folder contents without Git. Private `.langdrift` state, dotenv files, dependency/cache directories, and symlinks are excluded. Each local build records a content fingerprint. Legacy GitHub bindings retain their original source checks until the console explicitly associates a stopped project with a local folder.

The following runtime-specific connection flow remains for compatibility:

```sh
./dist/ex init
# Read the one-time secret without putting its value in shell history.
read -r -s RUNTIME_SECRET
export RUNTIME_SECRET
./dist/ex agent connect https://your-console.example
unset RUNTIME_SECRET
./dist/ex agent supervise
```

Use the console's newly generated runtime credential. The agent writes it and the restricted environment recovery credential to `.langdrift/agent/state/credentials/<runtime-id>.env` with private permissions, inside the ignored private state, and writes the connection record under `.langdrift/agent/state/connections/` only after that file. A credential file whose record exists belongs to a connected runtime: it is never replaced, and the connection is refused. A credential file without a record is what an interrupted connection left behind, and the next attempt replaces it. Never pass credentials as command arguments.

Results use stdout; diagnostics use stderr.

## Source boundaries

```text
src/main.go, src/wiring.go   Process entry and concrete dependency wiring
src/presentation/cli/        Commands, flags, prompts, the session's lines, exit codes
src/application/             Use cases and the interfaces they consume: inputs, plan, execution, connected
src/domain/                  Identity, configuration, ownership, lifecycle rules, the records, money and the RAG
src/infrastructure/          Files, storage, records, native tools, the judge, the stage driver, the lifecycle, delivery
```

The `runtime.md` fence, the lexical rules and the protocol types live in the
Microruntimes contracts, `packages/contracts`; `src/domain/connection` and `src/application/connected`
name them for the CLI, and `src/domain/microruntime` delegates its identifier
and path rules to them. The domain does not import the CLI, filesystem, or
native process adapters.
Application use cases own orchestration; the entry point wires concrete
adapters. These boundaries are a convention: without an `internal` directory,
the Go compiler does not stop other modules from importing the packages. The
module is private and unpublished, so only this repository builds it.

## Configuration and recovery

Each runtime has a stable console UID, native runtime ID, environment identity and distinct execution identities. Binding preserves the existing local project identity. MCP uses HTTPS (HTTP only for explicit loopback development), bounded JSON-RPC Streamable HTTP responses, initialization and tool calls `langdrift.bind`, `langdrift.sync`, and `langdrift.claim`. Redirects are refused. A console URL that does not parse is refused without being quoted, so a credential pasted into it never reaches the terminal. The recovery credential can report cleanup after revocation; it cannot claim another start.

Authored configuration is `.langdrift/microruntimes/runtimes/<id>/runtime.md`, containing exactly one `langdrift` JSON fence, and, when `ex agent plan` wrote it, one `langdrift-plan` fence. Free prose is retained but never executed; the planner's judge is the one paid model call this CLI makes, and the interpreter of the fence stays deterministic. The fence declares identity, revision, component paths, existing build recipes and images, ports/readiness, aggregate CPU/memory per execution, and the LangDrift agent model/effort. Between one and sixteen components, at most 64 CPU cores, and at most 1,048,576 MiB are accepted. Generated Deployments use one replica per component and explicit resource shares whose sum equals the envelope; resource quotas do not establish monetary limits or observed capacity.

The agent records the authored Markdown hash, interpreted revision, source commit and concrete native fingerprint. For a console binding, repository origin and HEAD must match the console's imported commit, and dirty source is an explicit conflict; generated `.langdrift` files are excluded from that check. A run materializes the runtime as its own operator on the project's folder, with the same ownership rules. Root components and distinct paths with equal basenames are supported. New unit identities derive from source paths.

Builds receive a disposable copy of Git-tracked source, excluding `.langdrift` and dotenv credential files. Existing `.dockerignore` and executable modes are preserved. The supported context is bounded to 8 MiB per source file and 128 MiB total; symlinks and external build contexts/credential mounts are refused. Explicit commands run with argv, never an implicit shell; project programs remain trusted consumer code, not sandboxed code. Inherited runtime, agent and common provider-secret variables are removed. Recipes must produce the configured local images. Missing tools, Dockerfiles, incompatible recipes and unsuccessful builds are explicit failures; no fallback build or tool installation occurs.

Every five seconds the connected agent synchronizes up to 100 durable events and current observations. A 15-second communication lease operates independently of finite native work. Confirmed lease loss cancels startup and attempts ownership-aware cleanup; cleanup failure remains unknown. The console cannot infer a stopped execution from a disconnected transport. Requests/results and user start/stop, connection loss, reconnect, deletion and revocation have distinct causes. Pending events keep stable IDs, increasing per-environment sequences, occurrence times and purge epochs. They are acknowledged before a new start claim. After reconnect, a new claim requires current intent and the acknowledged configuration revision; user stop, deletion or revocation prevents resume. Availability is reported only after native readiness. Missing metrics stay unavailable. Resource aggregates require complete samples for every running pod and container; native measurement timestamps and windows are retained, and repeated stale samples do not become new observations. Private connection state has a shared 4 MiB read/write limit with reserved capacity for cleanup and shutdown outcomes. A lost start-claim response is reconciled using the console execution identity before another start can be attempted.

Only confirmed stopped executions may materialize another revision. Existing authored/generated modifications block replacement. Native configuration is materialized as `schema.json`, Kubernetes files and Terraform configuration for the existing ownership-aware executor. Stopping/deleting cleans only recorded owned resources, preserving the shared project cluster, sources and evidence. A failed start is cleaned up and does not retry repeatedly under unchanged user intent. An interrupted agent restarts with unknown state and reconciles cleanup before another start. A run whose cleanup was not confirmed ends with its outcome unknown and its reservation held; the next lap cleans up before materializing again.

For v2 migration, run `ex agent migrate <id>` (or explicitly `ex agent migrate <id> --recover` for an active/unsettled legacy execution), select **Existing JSON v2 runtime ID** in the console, and configure the same project components and their explicit build recipes. The active descriptor is archived privately; original native files remain unchanged and discoverable. The new descriptor points to generated `console-units` and `console-terraform` paths, retaining existing matching unit names and namespace/context. Console resource changes apply to these new outputs. Consumer workload settings outside image, replica and resource fields are retained; ambiguous multi-workload or multi-container units refuse automatic migration. Historical execution records and evidence remain intact, and migration creates no usage measurements.

Private kubeconfig, operation ledgers, connection records, credentials and pending events live under `.langdrift/agent/state`. Evidence lives in each runtime's `evidence/` folder, keyed by execution, and a runtime that was started without a plan gets a `run.json` of kind `execution`. Preserve operational ownership until cleanup succeeds. Permanent console history deletion advances a purge fence: buffered old events may reconcile cleanup but cannot restore erased history. Already acknowledged local event copies are removed; operational state remains separate from history.

## Tests and verification limits

`go test ./...` runs the suites kept beside the code they cover: the domain
rules; every use case through in-memory fakes of its ports, including a
Kubernetes fake that imitates the controllers a runtime relies on; the file,
storage, configuration, and scaffold adapters in disposable directories; the
kubectl, kind, Terraform, toolchain, hook, stage driver and delivery adapters
against recorded command invocations; the judge against local HTTP servers,
redirects, oversized answers, provider error bodies and plain HTTP base URLs
included;
and the composed CLI. The session's tests drive the composed CLI through a
project seeded from the fixtures with a scripted model, a scripted driver, a
scripted lifecycle and a scripted delivery observer, and reproduce the figures
the website's session shows. Connected-agent tests additionally cover durable
replay, conditional resume, source/configuration conflicts, migration, private
credentials and bounded MCP transport using fakes and disposable files. The
process runner starts only local `sh` processes. These suites never run Docker,
kind, kubectl, Terraform, a model, a forge or an agent application.

Measure coverage across package boundaries, and use the race detector where a
C toolchain is available:

```sh
GOTOOLCHAIN=local go test -coverpkg=./... -coverprofile=coverage.out ./...
GOTOOLCHAIN=local go tool cover -func=coverage.out
GOTOOLCHAIN=local go test -race ./...
```

Benchmarks need no infrastructure and run only when requested. They cover
process startup, command dispatch, cluster inventory parsing, cluster
deletion reconciliation, process spawn and cancellation, the project lock,
and the durable state write. The recorded baseline and its limits were kept
in the upstream SPEC-EX-HARDENING, which did not come over with the CLI.

```sh
GOTOOLCHAIN=local go test -run '^$' -bench . -count=10 ./...
```

Tests, static checks, and source reviews do not establish successful live
cluster startup, readiness, cleanup, or production integration. Those outcomes
require separately authorized execution and recorded observations.
