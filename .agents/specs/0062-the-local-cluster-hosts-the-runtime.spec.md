# The local cluster hosts the runtime

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E15.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 18 of 19, after spec 0063 and before spec 0064. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation the microruntime appears twice, late in the run. After the agent's code review, `runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB`: five Deployments, one per component, at the shares the plan proposed. After the pull request's checks and Safe Seat, `runtime    stopped · 5 Deployments removed · nothing left`.

The LangDrift CLI does it differently. It brings the runtime up right after the gate, before the first stage, and keeps it up through every human wait (ticket 0225). It sizes each workload with an equal split of the whole envelope (400m and 614 MiB each) instead of the proposed shares. It builds images from Dockerfiles, reads readiness from `HEALTHCHECK`, writes env-file Secrets, prepares a namespace and a quota with Terraform, and keeps the kind cluster when it is done while printing "nothing left". It adds lines the simulation does not have: a failed bring-up that ends a lap after a 20-minute timeout, and "runtime unknown … the reservation is held" when a cleanup is not confirmed.

## Solution

After `code-review`, `dft agent run` creates a kind cluster for this run, starts one one-replica Deployment per component of the applied `runtime.md` revision, at exactly the proposed share, and prints the ready line once all five are available. After `✓ Safe Seat  met`, it deletes the Deployments and the cluster, confirms that nothing of the run remains, and prints the stopped line. "Nothing left" is true: no cluster, no Deployment, no container of the run survives it. Any stop between the two takes everything down without a line.

## Responsibility

This element brings the runtime up after the code review, keeps it up until Safe Seat, takes it down and proves that nothing is left, and keeps the execution ledger of what it created. It decides how workloads are named, what "ready" means and whether anything is built, which spec 0066 leaves to it.

It does not:

- resolve images, ports, stacks or shares (specs 0066, 0055);
- run the stages or the code review (specs 0058, 0059);
- push, open the pull request or follow its checks (specs 0063, 0064);
- decide Safe Seat (spec 0061);
- measure CPU or memory for an episode (spec 0052 records `measured: null`);
- settle money (spec 0056).

## Produces

- **The cluster**, outside `.langdrift/`: one kind cluster named `langdrift-<runtime id>` (`langdrift-book-a-demo-runtime`), and in it one Deployment per component of the applied revision:
  - named after the component's path, with `/` replaced by `-` (`services-signup-api`, `services-notify-api`, `platform-nats`, `platform-redis`, `apps-signup`);
  - one replica (the fence's `replicas`);
  - one container running the component's resolved `image` on its `port`;
  - CPU and memory, both requested and limited, at the component's proposed share (`250m` and `256Mi` for `services/signup-api`).

  Nothing else is created: no namespace, quota, Service, Secret or volume.
- **The execution ledger**, private, at `.langdrift/agent/state/runs/<run>/execution.json`, beside the run's stage reports, where `<run>` is the run's id (spec 0065): the cluster's name and the five Deployments' names, written before each is created, and settled with what was removed and the cluster's confirmed absence. It feeds the removed count and "nothing left", and lets a later run remove what a stopped one left.
- **The runtime's state, in the run's record (spec 0065)**: `revision`, `deployments` (`5`), the applied totals (`cpuMillis` `1450`, `memoryMiB` `2304`), `state` (`ready`, then `stopped`), `removed` (`5`) and `left` (`0`). Token-free.

## Shows

Phase 2 only, never in the over-budget variant. Both are output lines, and both are the same with Review ticked or not and with any stages checked.

After `console › code-review 2 findings · fixed` (spec 0059) and before `✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record` (spec 0065):

```
runtime    book-a-demo-runtime · revision 1 · 5 Deployments ready · 1450m · 2304 MiB
```

The runtime id and the revision are the applied intent's (spec 0050); the count is the number of Deployments started; the figures are the sum of their shares, in plain integers.

After `✓ Safe Seat  met` (spec 0061) and before `✓ ready      #19 for review · $1.60 released to the envelope` or `✓ ready      #19 to merge · $1.60 released to the envelope` (spec 0063):

```
runtime    stopped · 5 Deployments removed · nothing left
```

The count is the number of Deployments the ledger records as removed. "nothing left" prints only when the cluster's absence is confirmed.

## User Stories

1. As a developer, I want the runtime to come up only after the code review, so that nothing runs while the stages and the person work.
2. As a developer, I want no runtime up while the run waits for a person, so that a human wait costs no local resources.
3. As a team lead, I want each component started at the share the plan proposed, so that what runs is what the team authorized.
4. As a team lead, I want the ready line to show the runtime, its revision, the number of Deployments and their total CPU and memory, so that I can check it against the Plan card.
5. As a developer, I want one one-replica Deployment per component, so that "5 one-replica Deployments" is what the cluster holds.
6. As a developer, I want each Deployment named after its component's path, so that I can find a component in the cluster.
7. As a developer, I want the ready line only once every Deployment is available, so that "ready" is a fact.
8. As a developer, I want the runtime to stay up until Safe Seat is met, as the simulation orders it, so that the run's steps keep the simulation's order.
9. As a team lead, I want the stopped line to say how many Deployments were removed, so that the cleanup is accounted for.
10. As a team lead, I want "nothing left" to be true, cluster included, so that the run leaves no local resource behind.
11. As a developer, I want "nothing left" printed only when the cluster's absence is confirmed, so that the line never claims what was not checked.
12. As a developer, I want every stop after the bring-up to take the runtime down without a line, so that a failure never leaves a cluster running.
13. As a developer, I want a later run to remove what a stopped run left, from the ledger, so that leftovers do not accumulate.
14. As the owner, I want no Terraform, namespace, quota, Secret or Service the simulation does not show, so that the runtime is exactly five Deployments.
15. As the owner, I want no bring-up line, timeout line or "runtime unknown" line, so that the terminal shows only the simulation's two runtime lines.
16. As a developer, I want the tools the runtime needs limited to Docker, kind and kubectl, so that I install only what the lines require.
17. As an auditor, I want the execution ledger written before each creation, so that cleanup acts only on what the run created.
18. As an auditor, I want the runtime's state in the run's record, so that the Workflow card and the evidence agree with the lines.

## Implementation Decisions

### Remove

- The bring-up right after the gate and the runtime kept up through the stages and human waits (ticket 0225). Serves the ready line's position after `code-review`.
- The 20-minute bring-up timeout as a line, the failed-lap ending of a bring-up, and taking the runtime down to "fail the lap". Serve nothing shown; a failure is a stop (spec 0048).
- "runtime    unknown · … · the reservation is held". Serves nothing shown; an unconfirmed cleanup is a stop.
- Keeping the kind cluster (`ClusterKept`). Serves "nothing left".
- The equal split of the envelope as the workloads' resources. Serves `1450m · 2304 MiB` and the per-component shares.
- Building images from Dockerfiles (ADR 0004), the build command, readiness probes read from `HEALTHCHECK`, env-file and passed-variable Secrets, and the cleanup of an earlier lap's open execution. Serve nothing shown; images are resolved (spec 0066).
- The Terraform preparation of a namespace and a quota, and the Terraform dependency. Serves nothing shown.
- The connected agent's native engine and its operator record, the claim of an execution, and the sampling of resource peaks. Serve nothing shown.
- The `cluster` and `runtime` commands (spec 0048 removes them).

### Build

- The bring-up step after the `code-review` line: record the cluster in the ledger, create the kind cluster `langdrift-<runtime id>`, make each resolved image available to it (pulled by the cluster, or loaded from the local Docker when it is only there), record and apply each Deployment, then wait until every Deployment reports itself available. No line prints until the ready line.
- "Ready" is the Deployment's own available condition with its one replica running. No probe is added.
- The ready line, from the applied revision and the started Deployments' shares.
- The take-down step after `✓ Safe Seat  met`: delete each Deployment, delete the cluster, confirm with kind that the cluster no longer exists, settle the ledger, then print the stopped line.
- The silent take-down on any stop after the bring-up started (spec 0048's rule): the same steps, without the stopped line.
- Before a bring-up, the execution ledgers of earlier runs of the same runtime that are not settled are read, and every cluster they record as left is deleted silently and their ledgers settled.
- The bring-up keeps an internal bound on the wait for availability, the existing 20 minutes; reaching it is a stop with no line of its own.
- The tools are run with the minimal subprocess environment the CLI already uses for native tools (ADR 0002).

### Change

- The workloads' resources change from the envelope's equal split to the fence's per-component shares, requested and limited.
- Workload names change from component names to path-derived names, since the fence no longer has component names (spec 0049).
- The kind cluster changes from a kept, shared cluster to one created and deleted by each run.

### Interfaces

- From spec 0050: the applied intent's runtime id and revision. From spec 0055: the `runtime.md` of that revision, with each component's path, image, port, share and replicas.
- From spec 0059: the moment `code-review` has printed. From spec 0061: the moment `✓ Safe Seat  met` has printed.
- To spec 0065: the runtime's state for the run's record, and the moment the ready line has printed, after which `✓ completed` prints.
- A cluster port with four operations (create, apply a Deployment and wait for availability, delete, confirm absence), so that the golden test scripts it.

### Failure paths

Every one follows spec 0048's stop rule, with no line of its own:

- kind, kubectl or Docker missing, the cluster not created, an image not available to it, a Deployment never available within the bound: the run stops before the ready line, takes down what it created and settles the ledger;
- a stop after the ready line and before Safe Seat (a failed push, a failed check, an interrupt): the runtime is taken down silently;
- the cluster's absence not confirmed: the stopped line does not print, the ledger keeps the cluster as left, and the run stops.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test**, with the cluster port scripted to report five Deployments available and then removed and the cluster absent. This spec adds:
  - the ready line and the stopped line, byte for byte, at their positions;
  - no runtime operation before the `code-review` line, in every variant, and none at all in the over-budget variant;
  - the five Deployments requested from the scripted cluster with the fixture's paths turned into names, one replica each, and the shares 250m/256Mi, 300m/512Mi, 200m/256Mi, 200m/256Mi and 500m/1024Mi;
  - the execution ledger settled with five removed and the cluster absent;
  - the runtime's state in the run's record.
- Golden stops: a Deployment that never becomes available, and a cluster whose absence is not confirmed. Each leaves standard output a prefix of the golden, takes down what was created, writes one sentence to standard error and exits 1.
- Unit tests only for the rules: the path-derived name, and the sum of shares for the ready line.
- Prior art: the existing lifecycle and kubectl tests with a scripted process runner.

## Out of Scope

- Image resolution, ports and shares (specs 0066, 0055).
- Delivery and checks (specs 0063, 0064), Safe Seat (spec 0061), money (spec 0056).
- The console and its stub API.

## Open Questions

- **What uses the runtime between ready and stopped.** The simulation brings it up after the code review and stops it after the pull request's checks and Safe Seat; the website's data calls it "the microruntime the checks use". No line shows anything reaching it, and the checks run on the forge (spec 0064). Until the owner decides, nothing in `dft` uses it: it is up, available and removed, in the simulation's order, with no Service between components.
- **Cleanup and waiting.** Deleting clusters that earlier stopped runs left behind before a bring-up, and waiting up to 20 minutes for readiness, are this spec's choices; the simulation shows neither.
- **"Nothing left".** This spec reads `nothing left` as deleting the run's whole kind cluster. The Microruntimes page's lifecycle copy says the shared cluster stays for the next run. Which reading holds is the owner's decision.

## Further Notes

- The ready line already matches the simulation's format in the CLI's catalog; only its moment and its figures' source change.
- ADR 0004's component Dockerfile no longer reaches the cluster; spec 0066 says which parts of ADR 0004 survive.
