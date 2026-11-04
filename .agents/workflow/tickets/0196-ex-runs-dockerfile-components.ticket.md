# 0196: ex runs Dockerfile components

**Spec:** .agents/workflow/specs/0032-components-by-dockerfile-path.spec.md

**What to build:** `ex agent run` brings a version 2 runtime up: it builds each component with the fixed `docker build` against its Dockerfile from the project root under the minimal environment, loads the derived image into kind, names the workload and the Service after the component so components address each other by name, renders the probe as `httpGet`, `tcpSocket` or `exec` by kind, writes one Secret per component from its env file and the pass-through names read from the agent's environment, references it with `envFrom`, records the Secret in the execution ledger so cleanup removes it, and refuses the start with the name of a pass-through variable the agent does not have. `ex runtime init` asks for the Dockerfile and the readiness kind. ADR 0004 records that component builds are Dockerfiles, supersedes ADR 0002's recipe allowlist while keeping its minimal-environment rule, and records the exception to ticket 0024 for microruntime components; the glossary renames "Build recipe" to "Component Dockerfile". This is the integrate-and-verify ticket of the branch shared with 0194 and 0195: both Go modules are green here.

**Blocked by:** 0195.

**Status:** ready-for-agent

- [x] The blueprint names workloads and Services by component name and carries the derived image; the three probe kinds render as Kubernetes expects.
- [x] The process fake records `docker build --file <folder>/<dockerfile> --tag <image> .` with the minimal environment and nothing else as the build.
- [x] The Secret is created from the env file plus the pass-through values, referenced from the container, and listed in the ledger; cleanup deletes it; a missing pass-through variable refuses the start with its name.
- [x] A native probe edited by the owner is preserved by materialization as today.
- [x] `ex runtime init` with the prompter fake asks for the Dockerfile name and the readiness kind.
- [x] ADR 0004 is accepted and the glossary is updated; the root README and its translations show the new fence.
- [x] `go test -race ./...` and `go vet ./...` pass in the contracts and CLI modules; `govulncheck` is clean.

### Verification

- `pnpm turbo run build test typecheck --filter=@langdrift/ex --filter=@langdrift/contracts`
- `platform/cli/drift quality check`
- `platform/cli/drift check localization` (the README translations)
