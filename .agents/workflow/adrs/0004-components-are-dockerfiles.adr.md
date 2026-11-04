---
status: accepted
---

# A microruntime component is built from its own Dockerfile

A component in a `runtime.md` fence used to carry an image name, a port, an HTTP readiness path and a build recipe, all authored inline beside the files that already said the same. The fence now points at files and does not define them: a component is its name and its folder, the folder holds a Dockerfile, and `ex` derives the rest. The image tag is `langdrift-<runtimeId>-<name>:r<revision>`, never authored. The build is the fixed `docker build --file <folder>/<dockerfile> --tag <image> .`, run from the project root. The port is the Dockerfile's first `EXPOSE` unless the fence names one. Readiness is the Dockerfile's own `HEALTHCHECK`, rendered as a Kubernetes `exec` probe, with `tcp` and `http` kept as explicit kinds for an image that cannot run a command or a page that proves readiness better. A Dockerfile with no `HEALTHCHECK` of its own, or with `HEALTHCHECK NONE`, and no explicit kind is refused at plan time. Workloads and Services take the component's name, so components address each other by it. The fence version is 2; an `ex` reading any other version asks to re-plan the runtime.

This supersedes the build-recipe allowlist of ADR 0002. The committed Dockerfile is what code review sees, and the only command `ex` runs to build is fixed, so no recipe arrives from the console or the fence. ADR 0002's other rule stands: every subprocess, the fixed build included, receives a minimal explicit environment, here the base plus the container engine's own `DOCKER_*` selection variables.

A component may name a committed env file without secrets and a list of variable names `ex` reads from its own environment at start. `ex` writes both into one Kubernetes Secret per component, `<name>-env`, referenced from the container with `envFrom`, and records the Secret in the execution ledger with the unit's other objects, so cleanup deletes it. The ledger records the Secret's identity, never its values, and the values travel in memory to the API server, never as command arguments. A pass-through name the agent's environment does not hold refuses the start with the name.

Ticket 0024 ("no container image") keeps standing for the website. Microruntime components are the recorded exception: building their images is what a microruntime is for.

## Considered Options

- **Keeping inline build recipes under the allowlist.** It works today, but the fence repeats what the Dockerfile says, and the allowlist has to keep growing with every tool a team builds with.
- **Reading readiness from the built image.** It would see a `HEALTHCHECK` inherited from a base image, but manifests are rendered before the build, so the plan would depend on an image that does not exist yet. Parsing the component's own Dockerfile keeps the plan deterministic and testable; a later check may read the built image and refuse one that lost its healthcheck.
- **Passing secrets through `kubectl create secret --from-literal`.** It is one command, but it puts every value in the process list; the Secret is built in memory and applied from standard input instead.
- **Naming workloads by a digest of their path.** It avoids collisions between folders of the same name, but no component could address another; component names are unique in the fence instead.

## Consequences

- An older fence, version 1, is refused with "re-plan the runtime"; there is no migration of inline images or recipes.
- A component whose image cannot carry its own `HEALTHCHECK` names `tcp` or `http` readiness in the fence.
- A secret reaches a container only through a pass-through name the fence declares; the declaration is the audit trail.
- The component's name is its folder's, so this repository's own components are `messaging` and `storage` under `infrastructure/` (spec 0044), and a fence that runs them names them so.
