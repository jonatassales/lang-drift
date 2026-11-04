# Components by Dockerfile path, and this repository's microruntime

**Status:** ready-for-agent

## Tickets

- [x] 0194: contracts v2: component by path, readiness from the Dockerfile, env by file
- [x] 0195: ex plans Dockerfile components
- [x] 0196: ex runs Dockerfile components
- [ ] 0197: release `@langdrift/ex` 0.3.0 (owner)
- [ ] 0198: this repository's Book a Demo microruntime

## Problem Statement

A microruntime's configuration, the `langdrift` fence of its `runtime.md`, describes each component by an image reference, a port, an HTTP readiness path and a build recipe, all authored inline. Four things follow that stop the repository from running its own services (spec 0031) under `ex`:

- The fence defines what the component's own files already say. A component has a Dockerfile; the fence should point at it and consume it, not repeat an image name and a build command beside it. The owner's direction is to decouple the fence from the dominant technologies as far as possible: it points at paths, it does not define.
- Readiness is HTTP only, rendered as an `httpGet` probe. Redis has no HTTP endpoint, so it cannot be a component, and a service that is not ready until its module graph is up cannot say so.
- Workloads are named by a digest of their path, so one component cannot address another by name: demo-requests cannot reach `redis`.
- No environment reaches a container. A service that needs the address of the log or a secret has no way to receive one.

Besides, `ex` is not a dependency of this repository, nothing here is planned or run by it, and ticket 0024 says the repository builds no image. The research note on Microruntimes left open whether that decision gets an exception for microruntime images.

## Solution

A component becomes a reference to files. The fence names the component and its folder; the folder holds the Dockerfile. `ex` derives the image tag, builds with `docker build` from the project root against that Dockerfile, takes the port from the Dockerfile's `EXPOSE` unless the fence says otherwise, and takes readiness from the Dockerfile's `HEALTHCHECK`, rendered as a Kubernetes `exec` probe, with `tcp` and `http` kept as explicit escapes. A component may name a committed env file without secrets and a list of variable names `ex` reads from its own environment at start, as the driver configuration already does; `ex` writes both into one Kubernetes Secret per component and records it in the execution ledger. Workloads and Services take the component's name, so components address each other by it.

The fence version becomes 2. The change ships as `@langdrift/ex` 0.3.0, which the owner releases. This repository then pins that release at its root, keeps the committed harness and finops inputs and a plan script under `platform/microruntime/`, and plans its microruntime of four components: `demo-requests`, `mailer`, `redis` and `minio`. `.langdrift/` stays machine-local and ignored. ADR 0004 records that component builds are Dockerfiles, supersedes the recipe allowlist of ADR 0002 while keeping its minimal-environment rule, and records the exception to ticket 0024 for microruntime components.

## User Stories

1. As an owner, I want a component in the fence to be its name and its folder, so that the fence points at files instead of defining what they already say.
2. As an owner, I want the image name derived by `ex` and never authored, so that two plans of the same revision build the same tag.
3. As an owner, I want the build fixed to `docker build` against the component's Dockerfile with the project root as context, so that no recipe arrives from outside the checkout.
4. As an owner, I want readiness read from the component's own `HEALTHCHECK`, so that the component owns its health semantics and the fence stays technology-neutral.
5. As an owner, I want a Dockerfile without its own `HEALTHCHECK` refused at plan time with the file named, so that a component never runs without a readiness signal.
6. As an owner, I want `tcp` and `http` readiness kept as explicit kinds, so that an image that cannot run a command or a page that proves readiness still works.
7. As an owner, I want the port taken from the Dockerfile's first `EXPOSE` when the fence omits it, so that the fence repeats nothing.
8. As an owner, I want a component to name a committed env file and a list of variables passed from the agent's environment, so that configuration and secrets reach a container through a declaration the review can see.
9. As an owner, I want the env file refused when it is outside the checkout or holds a name that is not an uppercase variable, so that the declaration stays auditable.
10. As an owner, I want the Secret `ex` creates recorded in the execution ledger, so that cleanup removes it and the evidence shows it existed.
11. As a developer, I want workloads and Services named after their component, so that `redis://redis:6379` resolves inside the microruntime.
12. As a developer, I want `ex agent plan` to take a component as `name=,path=` plus the optional keys, so that planning this repository's runtime is one command per component.
13. As a developer, I want the plan's checks to say which Dockerfile is missing, which has no `HEALTHCHECK`, which env file is absent and which ports collide, so that a draft fails with a reason.
14. As a developer, I want the fence version 2 refused by an older `ex` and version 1 refused by the new one with "re-plan the runtime", so that a mismatch is a clear message, never a silent misread.
15. As a developer, I want the golden fixtures, the README fences and the site's examples to show the new shape, so that the documented contract is the executed one.
16. As a developer, I want `ex runtime init` to ask for the Dockerfile and the readiness kind instead of an image and a recipe, so that the guided path matches the fence.
17. As a maintainer, I want `@langdrift/ex` pinned at the repository root from the npm release, so that the repository consumes the CLI the way a customer does.
18. As a maintainer, I want the harness and finops inputs committed under `platform/microruntime/` with a POSIX plan script, so that anyone can plan the runtime from a fresh checkout.
19. As a maintainer, I want `.langdrift/` to stay ignored, so that machine-local records never reach a commit.
20. As a maintainer, I want a root `.dockerignore`, so that a root-context build sends only what the services need.
21. As the owner, I want ADR 0004 to record that component builds are Dockerfiles and that microruntime components are the exception to "no container image", so that the decision outranks the ticket that said otherwise.
22. As a developer, I want `ex agent run` on this repository's runtime to build four images, load them into kind, apply the workloads and Secrets, wait for the exec probes and pause at the first human stage, so that the pause is the proof the runtime came up.

## Implementation Decisions

**Contract.** In the Microruntimes contracts, a component is:

```json
{
  "name": "demo-requests",
  "path": "apps/microservices/demo-requests",
  "dockerfile": "Dockerfile",
  "port": 3402,
  "readiness": { "kind": "tcp" },
  "envFile": "apps/microservices/demo-requests/runtime.env",
  "passEnv": ["MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "DEMO_REQUESTS_TOKEN"]
}
```

`name` is a unique DNS label; `path` is the component folder as today; `dockerfile` is a file name inside the folder, defaulting to `Dockerfile`, and must exist; `port` defaults to the first `EXPOSE`; `readiness` defaults to the Dockerfile's `HEALTHCHECK` and otherwise is `tcp` or `http` with a path; `envFile` is a clean relative path inside the checkout whose keys are uppercase variable names and whose values are plain text within a bound; `passEnv` holds names of the same shape that do not overlap the file. `id`, `kind` and the technologies stay as they are. `image`, `build` and `readinessPath` leave the fence; the build-tool allowlist and the recipe-refused error go with them. The fence version is 2; validation names the re-plan on any other version.

**Readiness.** The planner parses the component's own Dockerfile `HEALTHCHECK` (shell form to `sh -c`, exec form as the array; interval, timeout, retries and start period to the probe's period, timeout, failure threshold and initial delay, with `ex`'s current defaults where an option is absent) and renders an `exec` probe. A `HEALTHCHECK` inherited from a base image is not seen: the component's own Dockerfile declares it. Parsing text keeps the plan deterministic and testable before any image exists, which the lifecycle needs because manifests are rendered before the build. An optional post-build check reads the image's healthcheck and refuses an image that lost it.

**Build and image.** The image tag is derived from the runtime id, the component name and the revision. The build is the fixed `docker build --file <folder>/<dockerfile> --tag <image> .` from the project root with the minimal environment ADR 0002 prescribes. The build-recipe allowlist is retired: the committed Dockerfile is what code review sees.

**Environment.** At start, `ex` reads the env file and the pass-through names from its own environment, writes one Kubernetes Secret per component, references it from the container with `envFrom`, and records the Secret in the execution ledger so cleanup owns it. A pass-through name absent from the agent's environment refuses the start with the name.

**Blueprint.** Workloads and Services take the component's name; the container port keeps its named port; the probe renders as `httpGet`, `tcpSocket` or `exec` by kind. Materialization preserves a native probe the owner edited, as it does today.

**Planner input.** `ex agent plan --component name=,path=[,dockerfile=][,port=][,readiness=tcp|http:/path][,env-file=][,pass-env=A;B][,kind=][,tech=id:Label;…]`. The plan's cross-checks become: ports distinct, Dockerfiles exist, readiness known, env files exist, within the ceiling.

**Fixtures and documentation.** The golden fixtures (sample and expected runtimes, plans, runs), the planner's test project (which gains a Dockerfile with a `HEALTHCHECK`), the CLI README's fence and flag examples, the root README and its translations, the glossary ("Build recipe" becomes "Component Dockerfile") and ADR 0004 change together.

**This repository.** `@langdrift/ex` is a root devDependency pinned to the exact release; with the workspace's default linking, a plain version resolves from the registry even though the CLI's package is in the workspace. `platform/microruntime/` holds the workflow, driver, prices and forecasts inputs and an idempotent POSIX `plan.sh` that initializes the project when needed, records the harness and finops inputs, and plans the runtime from the microruntime ticket with four components: `redis` and `minio` from their Dockerfiles, `demo-requests` and `mailer` with their committed `runtime.env` (log address, MinIO address, port) and their pass-through secrets. The first workflow marks every stage as requiring a human, so the run brings the components up and pauses.

**Out of the contract.** The stub API and the console's runtime form still speak the old fields; they do not import the contracts module and follow in their own spec.

## Testing Decisions

A good test feeds a fence, a Dockerfile or a flag string and asserts the parsed configuration, the rendered manifest or the named refusal; it never reaches Docker, kind or the network, in line with spec 0016.

Seams:

1. **Configuration parsing and validation** in the contracts module: version 2 accepted, version 1 and unknown fields refused with the re-plan message, each field rule exercised with one refusing case, the old fields refused as unknown.
2. **Dockerfile inspection** in the planner: `HEALTHCHECK` in both forms with and without options, `HEALTHCHECK NONE`, no `HEALTHCHECK` (refused with the file named), the first `EXPOSE` as the default port.
3. **The plan session** through the existing golden fixtures: `ex agent plan` with the new flags reproduces the expected `runtime.md`, `plan.json` and `run.md`, and the findings name each cross-check.
4. **Blueprint and manifest rendering**: workloads named by component, the three probe kinds, the Secret reference, the fixed build command captured by the process fake, the ledger entry for the Secret.
5. **The guided `ex runtime init`** with the prompter fake asking for the Dockerfile and the readiness kind.

Prior art: the contracts' configuration tests, the planner's plan, findings and review tests, the connected native tests with their process fakes, the golden fixtures under the CLI's test data. Go tests with `-race`, descriptions as today.

The repository's own runtime is verified by hand on a machine with Docker, kind, kubectl and Terraform, and the pull request summarizes the run.

## Out of Scope

- A components-only operation that brings a runtime up without a ticket run; the first workflow's human stages give the pause.
- Hosting the microruntime anywhere but the local kind cluster.
- The stub API and the console's runtime form adopting the new fields.
- Secrets from anywhere but the agent's environment (no vault, no file of secrets).
- Multi-port components; one named port per component as today.
- The release itself, which is the owner's.

## Further Notes

- ADR 0004 supersedes ADR 0002's recipe allowlist only; its minimal-environment rule for every subprocess stays and applies to the fixed build.
- Ticket 0024 ("no container image") keeps standing for the website; microruntime components are the recorded exception.
- The first dogfood run of the platform's own engine is expected to surface more contract gaps; each one becomes a finding in the plan and a ticket, never a workaround in a Dockerfile.
