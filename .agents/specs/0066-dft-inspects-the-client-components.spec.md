# `dft` inspects the client's components

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E19.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 10 of 19, after spec 0051 and before spec 0055. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The simulation's plan resolves the five components named in the spec, "paths, images, ports". For each one, it inspects the manifest that identifies its stack (`go.mod`, `package.json`, `nats.conf`, `redis.conf`, `next.config.ts`) with the port it serves on, and prints the component's stack, role and proposed share. It then cross-checks that the five ports are distinct, that five images are resolved and that nothing collides. The Plan card shows the five rows and the ports line.

The LangDrift CLI gets its components from `--component name=,path=…` flags and prints "named in the CLI". It takes the port only from `port=` or the Dockerfile's `EXPOSE`, ignoring the port the manifest declares. It does not know `next.config.ts`: a Next.js app prints as "package.json" and "Node.js+Next.js". It requires a Dockerfile with a `HEALTHCHECK`, and runs env-file, readiness and within-ceiling checks, each able to block the plan, so a NATS or Redis folder without a Dockerfile blocks it. Its "5 images resolved" is the component count, and its images are derived tags that nothing resolves. A team whose spec names five components gets a plan that ignores the spec and demands files the simulation never mentions.

## Solution

The components are exactly the paths the spec names (read by spec 0053), in the spec's order. For each path, `dft` reads the folder's manifest by a fixed precedence:
- `next.config.ts` → Next.js · web;
- `go.mod` → Go · api;
- `package.json` → Node.js · api;
- `nats.conf` → NATS · broker;
- `redis.conf` → Redis · cache.

That gives the stack and the role. It takes the port the manifest declares and resolves the component's image. The share comes from a fixed per-stack sizing table. After the five components, it cross-checks ports, images and paths for collisions. No Dockerfile, healthcheck, env file or ceiling takes part in the plan. The lines print as the simulation shows, and the card's COMPONENTS rows and ports line carry the same data. Which parts of ADR 0004 survive is stated below.

## Responsibility

This element:
- takes the component paths from the spec (handed over by spec 0053) and resolves each one: its manifest, stack, role, port, image and proposed share;
- cross-checks the resolved components for distinct ports, resolved images and no collisions;
- prints the `resolving` line, the five `inspect` tool calls, the five `✓ <path>` lines and the `cross-checking` line;
- owns the Plan card's five COMPONENTS rows and the ports line.

It does not:
- read the spec or count its components (spec 0053);
- sum the shares, check them against the bound, compute headroom or write `runtime.md` (spec 0055);
- build, load, start, probe or remove anything; images are resolved here, never built or run (spec 0062);
- print the column padding of its lines (spec 0048).

## Produces

No file of its own. Each resolved component goes to the `runtime.md` draft (spec 0055) and the plan's card data (spec 0050):

- `path` (`services/signup-api`), in the spec's order;
- `manifest` (`go.mod`, `package.json`, `nats.conf`, `redis.conf` or `next.config.ts`), the file that decided the stack;
- `stack` (`Go`, `Node.js`, `NATS`, `Redis`, `Next.js`) and `role` (`api`, `broker`, `cache`, `web`);
- `port` (integer: `8080`, `8081`, `4222`, `6379`, `3000`);
- `image`, the resolved image reference, recorded in the draft for spec 0062;
- `cpuMillis` and `memoryMiB` (integers, millicores and MiB), the proposed share from the sizing table;
- the **cross-check** result: `portsDistinct` (`5`), `imagesResolved` (`5`) and `collisions` (`0`).

They feed this element's lines and fields. Through spec 0055 they also feed the size line, the Plan figures, the bars, the components' total and the runtime's `5 Deployments ready · 1450m · 2304 MiB` (spec 0062). Through spec 0052 they feed the stacks and roles of the episode the run leaves.

## Shows

Phase 1, `dft agent plan`. After `# plan · resolving the components, then sizing them inside the envelope` (spec 0048):

```
console › resolving 5 components named in the spec · paths, images, ports
console › tool call: inspect services/signup-api · go.mod · :8080
✓ services/signup-api   Go · api        250m · 256 MiB
console › tool call: inspect services/notify-api · package.json · :8081
✓ services/notify-api   Node.js · api   300m · 512 MiB
console › tool call: inspect platform/nats · nats.conf · :4222
✓ platform/nats         NATS · broker   200m · 256 MiB
console › tool call: inspect platform/redis · redis.conf · :6379
✓ platform/redis        Redis · cache   200m · 256 MiB
console › tool call: inspect apps/signup · next.config.ts · :3000
✓ apps/signup           Next.js · web   500m · 1024 MiB
console › cross-checking · 5 ports distinct · 5 images resolved · no collisions
```

The line after them is `✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares` (spec 0055).

Card fields, Plan card:

- **COMPONENTS rows**: five placeholder rows from the `resolving` line. Each row is set by its own `✓` line, with the stack's mark, the path, the stack, the role, the cpu and the memory, and a bar of the component's cpu as a share of the bound's 2,000m:
  - `services/signup-api  Go       api     250m  256 MiB`
  - `services/notify-api  Node.js  api     300m  512 MiB`
  - `platform/nats        NATS     broker  200m  256 MiB`
  - `platform/redis       Redis    cache   200m  256 MiB`
  - `apps/signup          Next.js  web     500m  1,024 MiB`
- **Ports line**: `ports 8080 · 8081 · 4222 · 6379 · 3000 · distinct · no collisions`, from the `cross-checking` line.

The components' total line under the rows belongs to spec 0055.

## User Stories

1. As a developer, I want the plan to resolve exactly the components my spec names, so that the runtime's scope is the spec's.
2. As a developer, I want the plan to say the components are "named in the spec", so that I know where the list came from.
3. As a developer, I want the components resolved in the spec's order, so that the terminal and the card read in the order I wrote them.
4. As a developer, I want each component's inspection to name the manifest it read and the port, so that I can check the plan looked at the right file.
5. As a developer with a Next.js app, I want `next.config.ts` recognized as Next.js, even beside a `package.json`, so that my web app is not sized as a Node.js API.
6. As a developer with a Go service, I want `go.mod` recognized as a Go API, so that the stack reads correctly.
7. As a developer running NATS or Redis from a config file, I want `nats.conf` and `redis.conf` recognized as a broker and a cache, so that infrastructure components need no Dockerfile.
8. As a developer, I want the port taken from the manifest when the manifest declares one, so that the plan uses the port my component really serves on.
9. As a team lead, I want each component's share taken from a fixed table by stack and role, so that the same component is always sized the same.
10. As a team lead, I want each `✓` line to show the stack, the role and the share, so that I can review the sizing component by component.
11. As a developer, I want the plan to check that no two components share a port, so that the runtime cannot start with a port conflict.
12. As a developer, I want the plan to check that every component's image resolves, so that the run does not fail later on a missing image.
13. As a developer, I want the plan to check that no two components collide on path, port or image, so that "no collisions" is a checked fact.
14. As a developer, I want no Dockerfile, healthcheck, env file or ceiling required at plan time, so that the plan asks only for what the simulation shows.
15. As a developer, I want a spec path that does not exist, has none of the five manifests, or gives no port, to stop the plan with the path named, so that I know which component to fix.
16. As a developer, I want a port or image collision to stop the plan with nothing written, so that a colliding draft never reaches the team.
17. As a console user, I want five placeholder rows while the components resolve and each row filled by its own line, so that the card follows the terminal.
18. As a console user, I want each row to show the stack's mark and a bar of its share of the cpu bound, so that I see which component takes the most.
19. As a console user, I want the ports line to list the five ports and say they are distinct, so that the cross-check is visible on the card.
20. As an implementer of the cluster, I want each component's port and resolved image in the draft, so that bring-up has what it needs without a Dockerfile.
21. As a maintainer, I want the parts of ADR 0004 that survive stated, so that nobody re-adds what ADR 0005 removed.

## Implementation Decisions

**Remove**
- Components from `--component name=,path=,port=,kind=,technology=` input, and the text "named in the CLI". Serves `resolving 5 components named in the spec · …`. The flags themselves are removed by spec 0048.
- The Dockerfile requirement at plan time: reading the component's Dockerfile, the port from its first `EXPOSE`, readiness from its `HEALTHCHECK` (and the `tcp`/`http` explicit kinds), the missing-Dockerfile and readiness checks, and the readiness record. The simulation's plan never mentions a Dockerfile.
- The env-file check, the env-file parser at plan time, and the pass-through names.
- The within-ceiling check (the ceiling itself is removed by spec 0054).
- `package.json` deciding Next.js, and the combined label "Node.js+Next.js". Serves `inspect apps/signup · next.config.ts · :3000` and `Next.js · web`.
- Manifests and stacks the simulation does not show: `pyproject.toml` (Python), `Dockerfile` as a manifest (Docker), the `Unknown` stack, and the `worker` role.
- Sizing from measured runs in place of the table (removed by spec 0052).
- "5 images resolved" as a count of components with derived tags. Serves `5 images resolved`, which counts images actually resolved.
- The derived image tag `langdrift-<runtimeId>-<name>:r<revision>` as the image the plan records, superseded by the resolved image (see Open Questions).

**Build**
- Manifest precedence for each component folder: `next.config.ts`, then `go.mod`, then `package.json`, then `nats.conf`, then `redis.conf`. The first file present decides the manifest, the stack and the role:
  - Next.js · web;
  - Go · api;
  - Node.js · api;
  - NATS · broker;
  - Redis · cache.
- The per-stack sizing table, in millicores and MiB:

  | Stack | Role | cpu | memory |
  | --- | --- | --- | --- |
  | Go | api | 250m | 256 MiB |
  | Node.js | api | 300m | 512 MiB |
  | NATS | broker | 200m | 256 MiB |
  | Redis | cache | 200m | 256 MiB |
  | Next.js | web | 500m | 1024 MiB |

- The port rule:
  1. the port the manifest declares (`port` in `nats.conf`; `port` in `redis.conf`; `PORT=`, `--port` or `-p` in a `package.json` script; a port set in `next.config.ts`);
  2. otherwise the stack's conventional port (NATS 4222, Redis 6379, Next.js 3000);
  3. otherwise the port the spec names beside the component's path (spec 0053).

  The source of the port for a manifest that cannot declare one is an open question below.
- An image resolver behind a port. It returns one image reference per component, or a refusal. Its production rule is an open question below. The fixture uses a resolver that resolves all five.
- The cross-check: ports distinct, every image resolved, and no two components sharing a path, a port or an image. The line is `cross-checking · <n> ports distinct · <n> images resolved · no collisions`.
- The lines and card rows above, one `inspect` and one `✓` line per component, in the spec's order.
- Safe stops:
  - A path that is missing, has none of the five manifests, or gets no port stops the plan before its `inspect` line.
  - An image that does not resolve, or any collision, stops the plan before the `cross-checking` line.

  In every case the plan writes no draft and exits non-zero, with the reason through the CLI's error channel (spec 0048). No line the simulation lacks is printed.

**Change**
- The component list's source moves from CLI flags to the spec's Components section.
- The port's source moves from `port=` or `EXPOSE` to the manifest, with the fallbacks above.
- `next.config.ts` becomes a known manifest with precedence over `package.json`.
- "Images resolved" becomes a resolution with a result per component, not a count of derived tags.

**What survives of ADR 0004**
- A component is its folder in the client's project, named by its path. The fence points at that path and never carries an authored build recipe. The image is resolved by `dft`, never typed into the fence or sent by the console.
- Components are unique by path, a rule the cross-check extends to ports and images.
- How workloads are named in the cluster, readiness, the build command if one remains, and the minimal subprocess environment of ADR 0002 are spec 0062's to keep or drop. The plan no longer reads a Dockerfile for any of them.
- Superseded: the Dockerfile as the component's definition, the derived tag, `EXPOSE` as the port, `HEALTHCHECK` as readiness at plan time, the env file and pass-through Secret, and the fence version 2 refusal (spec 0049 owns the fence's version).

**Interfaces with neighbours**
- From spec 0053: the spec's component paths in order, each with the port the spec names beside it, if any.
- From spec 0054: the bound's cpu, the denominator of each row's bar.
- To spec 0055: each component with share, port, image, stack and role, as it resolves.
- To spec 0062: the resolved image and port of each component, through the draft.
- To spec 0050: the COMPONENTS rows and the ports line.

## Testing Decisions

- The seam is the golden test of spec 0048 over the shared fixture. This element adds its twelve lines, the five COMPONENTS rows and the ports line, byte for byte. The fixture project holds:
  - `services/signup-api/go.mod`, with the spec naming `:8080` beside the path while the port rule stays open;
  - `services/notify-api/package.json`, whose start script sets `PORT=8081`;
  - `platform/nats/nats.conf` with `port: 4222`;
  - `platform/redis/redis.conf` with `port 6379`;
  - `apps/signup/next.config.ts` beside a `package.json`.

  No Dockerfile exists in the fixture. Its image resolver resolves all five.
- Golden cases for the safe stops: a spec path with no known manifest, two components on the same port, and an image the resolver refuses. Each prints up to the stop, writes no draft and exits non-zero.
- Unit tests for the rules only:
  - manifest precedence (`next.config.ts` beside `package.json`);
  - the port rule's order;
  - the sizing table;
  - the cross-check.

## Out of Scope

- Reading the spec and counting its components: spec 0053.
- The totals, headroom, size line and `runtime.md`: spec 0055.
- Building, loading and starting images, readiness, workload names and cleanup: spec 0062.
- The fence's fields and version: spec 0049.
- The column padding and the `console ›` prefix: spec 0048.
- The console and its stub API.

## Open Questions

- **How an image is "resolved" without a Dockerfile.** The simulation says `5 images resolved` and never mentions a Dockerfile or a build. For NATS and Redis, a published image of the stack would do, but at which version and from which registry? For the Go, Node.js and Next.js components, which come from source, the plan needs some way to resolve or produce an image:
  - a stack builder;
  - a prebuilt image the project names;
  - a Dockerfile kept only for the build, outside the plan.

  This spec builds the resolver port and the fixture resolver. The production resolver waits for this answer.
- **Where the port comes from when the manifest cannot declare one.** `go.mod` never declares a port, and `next.config.ts` rarely does. The simulation still prints `go.mod · :8080`. The provisional rule takes the manifest's port, then the stack's conventional port, then the port the spec names beside the path. Is that the rule, or should every component's port come from the spec?
- **Where the shares come from.** The fixed per-stack sizing table, and the fixed precedence over the five manifests, are this spec's choices. The simulation shows the shares and says "sizing follows the spec, not history", but not how a share is derived.

## Further Notes

- The simulation's `✓` lines are column-aligned (`✓ services/signup-api   Go · api        250m · 256 MiB`). The values here are path, `<stack> · <role>` and `<cpu>m · <memory> MiB`. The padding between them is spec 0048's.
- Real components often need environment variables. The simulation shows none, so none are read at plan time. A need for them enters the simulation first (ADR 0005).
