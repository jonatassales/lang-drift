# 0198: this repository's Book a Demo microruntime

**Spec:** .agents/workflow/specs/0032-components-by-dockerfile-path.spec.md

**What to build:** From a fresh checkout on a machine with Docker, kind, kubectl and Terraform, a developer plans and runs the Book a Demo services as a microruntime with the platform's own CLI. `@langdrift/cli` is pinned at the repository root to the 0.4.0 release; `platform/microruntime/` holds the workflow (every stage requiring a human), the driver, the prices and the forecasts as committed inputs and an idempotent POSIX `plan.sh` that initializes the project when needed, records the harness and finops inputs and plans the runtime from this ticket with four components: `messaging` and `storage` from their Dockerfiles under `infrastructure/`, `demo-requests` and `mailer` with their committed `runtime.env` (log address, MinIO address, port) and their pass-through secrets. A root `.dockerignore` bounds the root-context builds. `.langdrift/` stays ignored. `ex agent run` builds the four images, loads them into kind, applies the workloads and Secrets, waits for the exec probes and pauses at the first human stage; a port forward to `demo-requests` lets the website's curl path run against the microruntime.

**Blocked by:** 0197, 0189, 0190, 0191, 0192.

**Status:** ready-for-agent

- [ ] The root `package.json` pins `@langdrift/cli` to `0.4.0` in `devDependencies` and `pnpm exec ex --version` prints `0.4.0` (resolved from the registry, not linked from the workspace).
- [ ] `sh platform/microruntime/plan.sh` is idempotent and ends with a `runtime.md` under `.langdrift/` whose fence lists the four components with derived images, ports and probes.
- [ ] `pnpm exec ex agent run <runtime-id>` with the pass-through secrets exported shows four Ready pods in the runtime's namespace and pauses at the first human stage.
- [ ] Through a port forward, posting a request to `demo-requests` appends the received event, the mailer's sent events follow, and the MinIO objects exist.
- [ ] `platform/microruntime/README.md` documents the prerequisites, the commands and each file; the context document names the folder.
- [ ] `.langdrift/` is still ignored and nothing under it is committed.

### Verification

- `pnpm install`
- `pnpm exec ex doctor`
- `sh platform/microruntime/plan.sh`
- `pnpm exec ex agent run <runtime-id>` and `kubectl --kubeconfig .langdrift/agent/state/cluster.kubeconfig get pods -n <runtime-id>`
- `pnpm exec ex agent status <runtime-id>`
- `platform/cli/drift quality check 02-tooling`
