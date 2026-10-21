# Twelve-Factor App — guardrails

**Local policy v1, adopted October 10, 2026.** Application contracts with their execution environment, based on the [Twelve-Factor App](https://12factor.net/). Each row combines the original factor with its application in this monorepo. Evidence mechanisms and adaptations are local decisions; they do not constitute a new official edition of the manifesto.

## Scope

Apply to modified services and processes and their configuration, builds, and execution interfaces. Libraries, protobuf contracts, and documents are not deployable services: apply only relevant factors. Stateful services, jobs, and managed platforms require contextual assessment. Adopting these criteria does not claim that the current scaffold already meets them.

## The 12 criteria

| ID / factor | Local rule | Review evidence |
| --- | --- | --- |
| 12F-01 — [Codebase](https://12factor.net/codebase) | Identify the source and version of each application. A monorepo may contain several services; shared dependencies have explicit boundaries. | Workspace path, change reference, and artifact link when available. Do not require a separate repository per service. |
| 12F-02 — [Dependencies](https://12factor.net/dependencies) | Declare dependencies and tools required for build/execution; record versions as appropriate for the technology. | Manifests, lockfiles when supported, and tool requirements. Distinguish file inspection from a completed build. |
| 12F-03 — [Config](https://12factor.net/config) | Separate values that vary by environment from code. Declare inputs, required fields, and safe local defaults. Secrets are references, not versioned values. | Documented configuration, handling of missing inputs, and injection mechanism. The local adaptation accepts environment variables, mounted files, or secret managers as appropriate to the environment. |
| 12F-04 — [Backing services](https://12factor.net/backing-services) | Declare contracts and references for backing services. Consumers select compatible implementations. | Connection configuration without credentials, dependencies, and compatibility limits. Switching a database to another product is not automatically transparent. |
| 12F-05 — [Build, release, run](https://12factor.net/build-release-run) | Distinguish artifact construction, release composition, and execution. Preserve the identity of what was delivered. | Source, version/digest when available, and release configuration. Identify configuration embedded by the framework at build time and record when rebuilding is required. |
| 12F-06 — [Processes](https://12factor.net/processes) | Disposable application processes keep durable state in appropriate dependencies. Temporary local state has an explicit lifecycle. | Destinations for sessions, files, and persistent data. Assess databases and other stateful components against their own contracts; do not require every component to be stateless. |
| 12F-07 — [Port binding](https://12factor.net/port-binding) | State the entry interface of network services and how the environment exposes them. | Ports/network configuration or an equivalent managed-service contract. Libraries and jobs without a server may justify non-applicability. |
| 12F-08 — [Concurrency](https://12factor.net/concurrency) | Separate process types when the case requires it and declare concurrency/scaling limits. | Execution model, intended capacity, and coordination behavior. Configured replicas or autoscaling do not prove measured capacity. |
| 12F-09 — [Disposability](https://12factor.net/disposability) | Define startup, readiness, and safe shutdown, including in-flight work. | Configuration/code and service-specific time and integrity criteria. Observed behavior requires authorized execution; without it, behavior remains unverified. |
| 12F-10 — [Dev/prod parity](https://12factor.net/dev-prod-parity) | State relevant differences in versions, protocols, and dependencies across environments. | A matrix or record of differences and their effects. Stubs have declared limitations. Parity does not require equal cost/capacity or demonstrate production availability. |
| 12F-11 — [Logs](https://12factor.net/logs) | Emit events for collection by the environment, preserving diagnostic information without exposing secrets. | Event destination, service identity, and correlation when applicable. Use stdout/stderr when appropriate; retention and collection belong to the integration. Metrics and tracing are complementary requirements. |
| 12F-12 — [Admin processes](https://12factor.net/admin-processes) | Administrative tasks use versioned code and application-compatible context, with explicit permissions and effects. | Command, version, target, and execution record when authorized. Native tools own repository tasks; permitted CLI helper scripts use `.sh`. Creating a helper does not authorize migrations or data changes. |

## Completion criteria

Record the assessment and evidence in the [review format](../../templates/review.md), citing affected IDs. The former shared states and exceptions index is unavailable; this format does not replace it or grant exception authority. Inspection can demonstrate declared configuration; behavior, capacity, and recovery require execution evidence. The scoped hook initiative permits its contract suite and build/test push gate; unrelated tests remain subject to the authorization rules in root AGENTS.

New services assess all 12 factors. Incremental changes assess affected factors; do not expand a task into rewriting the entire repository. Preexisting differences outside the scope are separate findings. An adapter or skill that does not perform verification must state that limitation.

## Updates and sources

Factor sources consulted October 10, 2026. The [official modernization](https://github.com/twelve-factor/twelve-factor/blob/next/UPDATE_FAQ.md) retains the focus on the application/platform interface. External changes require local review before changing these criteria. FinOps, access control, and agent governance complement the manifesto; they are not additional original factors.
