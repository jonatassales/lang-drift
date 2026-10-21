# LangDrift CLI

This Go module, adopted from the Exphere repository, bootstraps projects and runs the independent agent connected to the LangDrift console. The console owns microruntime creation, composition, configuration and requested lifecycle; the agent builds explicit project recipes and reuses the bounded local kind/Terraform/Kubernetes lifecycle. The SDK remains postponed and applications need no LangDrift dependency. This is a local development candidate, not a published release or a claim of operational readiness.

## Build and checks

Use an existing Go 1.23 or newer toolchain on macOS or Linux:

```sh
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex --help
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go vet ./...
```

The module retains its standalone identity, `github.com/jonatassales/tao/packages/ex`, and is
integrated here as the private workspace package `@langdrift/ex`. This does not
imply a published release. Native Go commands own the build and checks. Building does not install infrastructure tools, alter shell
profiles, or change PATH. Use `./dist/ex` to avoid invoking the unrelated system
editor also named `ex`.

## Workspace package

The package holds only the Go module, its source, and `build.sh`. Installing
the workspace runs `build.sh`, which builds `dist/ex` with an existing Go
toolchain, never installs Go or tools, and never changes shell profiles. Set
`LANGDRIFT_SKIP_EX_BUILD=1` to skip it. From the workspace root,
`pnpm --filter @langdrift/ex build` rebuilds the executable.

This package declares no npm bin. Run `packages/ex/dist/ex` directly from this
repository. The separate SDK package and its bridge were removed on October 21, 2026. The former shell-setup command,
package bridge, historical schemas, and consumer-harness fixtures were removed
on October 17, 2026; Git history retains them.

## Commands

In the commands below, `ex` means this tree's compiled `./dist/ex` executable.

| Command | Behavior |
| --- | --- |
| `ex version` | Print the version; `--version` and `-V` are equivalent. |
| `ex init` | Create or reuse project identity; add the agent instructions, 26 installed Pocock skills with attribution, SDD workflow, policies and templates. Existing files remain unchanged. |
| `ex agent project <console-url>` | Read `PROJECT_SECRET` and supervise all microruntimes of one project, including those created after connection. The console normally launches this process. |
| `ex agent connect <console-url>` | Read `RUNTIME_SECRET` from the environment, verify the selected GitHub checkout, and save a private environment binding. |
| `ex agent run` | Synchronize all bound microruntimes, apply stopped configurations, build/start/observe/stop from authenticated console intent. |
| `ex agent migrate <runtime-id> [--recover]` | Write an explicit import guide for a stopped v2 runtime; `--recover` explicitly invokes ownership-aware legacy cleanup first. Preserve original native files, source and evidence. |
| `ex doctor` | Inspect existing native tools without installing anything. |
| `ex cluster create/inspect/delete` | Explicit local cluster administration, preserving existing ownership checks. |
| `ex cluster load-image <image>` | Explicitly load an existing image into the private project cluster. |
| `ex agent init <name>` | Retain the separate editable finite LangGraph Job scaffold; it is not the connected agent. |

All previous `ex runtime ...` creation, configuration and operation commands now refuse with console/migration guidance. Their native use cases remain internal to the connected agent and test fixtures.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success. |
| 1 | The operation failed. |
| 2 | Invalid usage: unknown command, flag, or argument shape. |
| 3 | Refused by a guard or policy: adoption of unowned resources, a changed configuration, a build recipe outside the allowlist, revoked console access. |
| 4 | A prerequisite is missing: a required tool is not installed, or a required project, cluster or runtime record does not exist. |
| 130 | Interrupted. |

Global flags are `--project <directory>`, `--json`, and `--timeout <duration>`.
Without `--project`, the nearest parent project is discovered. Finite commands default to two minutes. `ex agent run` owns a foreground process until interrupted; each native operation is separately bounded to two minutes and cleanup uses its own bounded context.

The standalone binary can be built directly with Go and placed in a user-selected executable directory as `ex`; Node, pnpm and the SDK are unnecessary in consumer projects. Builds do not modify a shell profile or PATH. No binary download, package publication or personal installation was performed in this delivery.

The local console normally handles bootstrap and agent launch. One project has one `.langdrift` and one agent; a process lock prevents a second project agent or legacy `agent run` from concurrently owning the same checkout. The agent admits at most 128 runtime workers, polls membership every five seconds, and joins its workers on cancellation. A console-owned stdin pipe also cancels the agent if its parent disappears. This process does not install dependencies or start workloads merely by connecting.

Local project bindings validate the canonical folder instead of requiring a GitHub origin, clean Git tree, or commit. Builds snapshot current tracked and nonignored untracked files, or bounded plain-folder contents without Git. Private `.langdrift` state, dotenv files, dependency/cache directories, and symlinks are excluded. Each local build records a content fingerprint. Legacy GitHub bindings retain their original source checks until the console explicitly associates a stopped project with a local folder.

The following runtime-specific connection flow remains for compatibility:

```sh
./dist/ex init
# Read the one-time secret without putting its value in shell history.
read -r -s RUNTIME_SECRET
export RUNTIME_SECRET
./dist/ex agent connect https://your-console.example
unset RUNTIME_SECRET
./dist/ex agent run
```

Use the console's newly generated runtime credential. The agent writes it and the restricted environment recovery credential to `.langdrift/runtimes/<runtime-id>/.env` with private permissions and a local ignore rule. A preexisting `.env` is preserved and reported as a conflict. Never pass credentials as command arguments.

Exit codes are 0 for success, 1 for operational errors, 2 for invalid usage,
and 130 for interruption. Results use stdout; diagnostics use stderr.

## Source boundaries

```text
src/main.go, src/wiring.go   Process entry and concrete dependency wiring
src/presentation/cli/        Commands, flags, prompts, output, exit codes
src/application/             Use cases and the interfaces they consume
src/domain/                  Identity, configuration, ownership, lifecycle rules
src/infrastructure/          Files, storage, native tools, embedded scaffolds
```

The domain does not import the CLI, filesystem, or native process adapters.
Application use cases own orchestration; the entry point wires concrete
adapters. These boundaries are a convention: without an `internal` directory,
the Go compiler does not stop other modules from importing the packages. The
module is private and unpublished, so only this repository builds it.

## Configuration and recovery

Each runtime has a stable console UID, native runtime ID, environment identity and distinct execution identities. Binding preserves the existing local project identity. MCP uses HTTPS (HTTP only for explicit loopback development), bounded JSON-RPC Streamable HTTP responses, initialization and tool calls `langdrift.bind`, `langdrift.sync`, and `langdrift.claim`. Redirects are refused. The recovery credential can report cleanup after revocation; it cannot claim another start.

Authored configuration is `.langdrift/runtimes/<id>/runtime.md`, containing exactly one `langdrift` JSON fence. Free prose is retained but never executed or sent to an LLM by this agent. The fence declares identity, revision, component paths, existing build recipes and images, ports/readiness, aggregate CPU/memory per execution, and the LangDrift agent model/effort. The console recommendation service uses that model setting; this deterministic interpreter makes no paid model call. Between one and sixteen components, at most 64 CPU cores, and at most 1,048,576 MiB are accepted. Generated Deployments use one replica per component and explicit resource shares whose sum equals the envelope; resource quotas do not establish monetary limits or observed capacity.

The agent records the authored Markdown hash, interpreted revision, source commit and concrete native fingerprint. Repository origin and HEAD must match the console's imported commit. Dirty source is an explicit conflict; generated `.langdrift` files are excluded from that check. Root components and distinct paths with equal basenames are supported. New unit identities derive from source paths.

Builds receive a disposable copy of Git-tracked source, excluding `.langdrift` and dotenv credential files. Existing `.dockerignore` and executable modes are preserved. The supported context is bounded to 8 MiB per source file and 128 MiB total; symlinks and external build contexts/credential mounts are refused. Explicit commands run with argv, never an implicit shell; project programs remain trusted consumer code, not sandboxed code. Inherited runtime, agent and common provider-secret variables are removed. Recipes must produce the configured local images. Missing tools, Dockerfiles, incompatible recipes and unsuccessful builds are explicit failures; no fallback build or tool installation occurs.

Every five seconds the agent synchronizes up to 100 durable events and current observations. A 15-second communication lease operates independently of finite native work. Confirmed lease loss cancels startup and attempts ownership-aware cleanup; cleanup failure remains unknown. The console cannot infer a stopped execution from a disconnected transport. Requests/results and user start/stop, connection loss, reconnect, deletion and revocation have distinct causes. Pending events keep stable IDs, increasing per-environment sequences, occurrence times and purge epochs. They are acknowledged before a new start claim. After reconnect, a new claim requires current intent and the acknowledged configuration revision; user stop, deletion or revocation prevents resume. Availability is reported only after native readiness. Missing metrics stay unavailable. Resource aggregates require complete samples for every running pod and container; native measurement timestamps and windows are retained, and repeated stale samples do not become new observations. Private connection state has a shared 4 MiB read/write limit with reserved capacity for cleanup and shutdown outcomes. A lost start-claim response is reconciled using the console execution identity before another start can be attempted.

Only confirmed stopped executions may materialize another revision. Existing authored/generated modifications block replacement. Native configuration is materialized as `schema.json`, Kubernetes files and Terraform configuration for the existing ownership-aware executor. Stopping/deleting cleans only recorded owned resources, preserving the shared project cluster, sources and evidence. A failed start is cleaned up and does not retry repeatedly under unchanged user intent. An interrupted agent restarts with unknown state and reconciles cleanup before another start.

For v2 migration, run `ex agent migrate <id>` (or explicitly `ex agent migrate <id> --recover` for an active/unsettled legacy execution), select **Existing JSON v2 runtime ID** in the console, and configure the same project components and their explicit build recipes. The active descriptor is archived privately; original native files remain unchanged and discoverable. The new descriptor points to generated `console-units` and `console-terraform` paths, retaining existing matching unit names and namespace/context. Console resource changes apply to these new outputs. Consumer workload settings outside image, replica and resource fields are retained; ambiguous multi-workload or multi-container units refuse automatic migration. Historical execution records and evidence remain intact, and migration creates no usage measurements.

Private kubeconfig, operation ledgers, connection records and pending events live under `.langdrift/state`. Evidence lives in `.audits/evidence/runs`. Preserve operational ownership until cleanup succeeds. Permanent console history deletion advances a purge fence: buffered old events may reconcile cleanup but cannot restore erased history. Already acknowledged local event copies are removed; operational state remains separate from history.

## Tests and verification limits

`go test ./...` runs the suites kept beside the code they cover: the domain
rules; every use case through in-memory fakes of its ports, including a
Kubernetes fake that imitates the controllers a runtime relies on; the file,
storage, configuration, and scaffold adapters in disposable directories; the
kubectl, kind, Terraform, toolchain, and hook adapters against recorded
command invocations; and the composed CLI. Connected-agent tests additionally cover durable replay, conditional resume, source/configuration conflicts, migration, private credentials and bounded MCP transport using fakes and disposable files. The process runner starts only local `sh` processes. These suites never run Docker, kind, kubectl,
Terraform, or an agent application; only the opt-in infrastructure E2E below
does.

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
and the durable state write. The recorded baseline and its limits are in
[SPEC-EX-HARDENING](../../.langdrift/artifacts/specs/0039-SPEC-ex-hardening.spec.md).

```sh
GOTOOLCHAIN=local go test -run '^$' -bench . -count=10 ./...
```

Tests, static checks, and source reviews do not establish successful live
cluster startup, readiness, cleanup, or production integration. Those outcomes
require separately authorized execution and recorded observations.

## Infrastructure E2E (opt-in)

`src/main_e2e_test.go` drives a test-only native lifecycle fixture against a disposable kind cluster. After console cutover this fixture preserves existing native ownership, cancellation, recovery, logs and forwarding assertions through the retained catalogs; it does not test the retired public runtime commands or establish end-to-end console integration. The fixture is compiled as a Go test executable and has no production command or environment switch. It compiles only with the `e2e` build tag and runs only when
`EX_E2E=1`; ordinary `go test ./...` never builds it, and
`go vet -tags e2e ./...` (the package `typecheck` task) keeps it compiling.
Running it creates real local infrastructure, so run it only where you are
authorized to, preferably on a disposable machine:

```sh
EX_E2E=1 GOTOOLCHAIN=local go test -tags e2e -run '^TestE2ELifecycle$' -count=1 -timeout 70m -v ./src
```

The test reserves 20 minutes of `-timeout` for cleanup and refuses to start
if fewer than 30 minutes would remain for the scenario, because Go skips
cleanup when the timeout expires. The first interrupt of the test process
stops the scenario and runs cleanup.

Prerequisites, checked before anything is created: a reachable Docker daemon;
`go`, `kind`, `kubectl`, and `terraform` on `PATH`; and a local base image
providing `sh`, `echo`, `sleep`, and BusyBox `httpd` (`busybox:1.36` by
default, or `EX_E2E_BASE_IMAGE`). A missing prerequisite fails the run with
instructions. The harness installs nothing and pulls no image, but kind
downloads its node image and Terraform the pinned Kubernetes provider when
they are not cached.

The scenario initializes a temporary project, creates its cluster, builds two
disposable images, generates a runtime with a Deployment and Service unit and
a Job unit through the test-only historical guided catalog, and loads the images. Operation names below identify that internal fixture, not available public commands. `runtime up`
must refuse a foreign namespace and record nothing. An `up` is then
interrupted once it has recorded its last object and waits for readiness;
`inspect` must report the interruption, a new `up` and a plain `down` must
refuse without changing the record or the live objects, and `down --recover`
must remove everything. A full `up` follows, with `inspect`, bounded `logs`
for both units, and a `forward` that answers HTTP and exits with code 130
when interrupted. Finally `runtime down` runs, `cluster delete` must refuse
while a foreign namespace exists and leave the cluster unchanged, and the
cluster is deleted with its private kubeconfig.

Every finite command gets an explicit `--timeout`, and every native tool sees a
disposable `KUBECONFIG`, so your own kubeconfig is never read or rewritten.
Stopping a CLI process sends an interrupt first, so it settles its records
and stops the tools it started. Cleanup runs even when the scenario fails. It
removes only what the run identified: its labeled foreign namespaces, the
runtime through `runtime down --recover`, the cluster through `cluster
delete`, and images whose tags still name the IDs recorded at build time. It
never deletes a cluster or namespace by name alone. Each removal is confirmed
by observation; anything still present, or whose absence cannot be observed,
fails the run with the reason, and the temporary workspace, including the
private kubeconfig, is kept for inspection. Docker's build cache and kind's
shared `kind` network are left as they are, because other work may use them.

Each run writes a sanitized report to
`.audits/evidence/runs/ex-e2e-<nonce>/report.json` at the repository root,
which Git ignores: tool and platform versions, each step's arguments, exit
code, duration, the native tools the CLI invoked directly, and bounded output
with paths replaced and credential-like values redacted, plus cleanup
results. It never contains a kubeconfig, Terraform state, plan, variable
file, or ownership token.

The manual `Infrastructure E2E` workflow in `.github/workflows/ex-e2e.yml`
runs the same test on a disposable GitHub-hosted runner and uploads the
evidence for seven days. It runs only when dispatched by hand, never on pull
requests or pushes.
