# `runtime.md` records the sized plan

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E8.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 11 of 19, after spec 0066 and before spec 0057. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

Once the five components are resolved, the simulation's plan sizes them inside the envelope: 1450m of 2000m and 2304 of 3072 MiB, five one-replica Deployments, the components' proposed shares. It says gpu and disk are not `runtime.md` fields and none was requested. It validates the `runtime.md` draft against the schema with 0 errors, as revision 1. The Plan card shows those figures, the CPU and memory bars, the components' total, BOUNDED, HEADROOM, SUGGESTED RESOURCE, GPU · DISK and VALIDATION.

The LangDrift CLI prints the same size line from the proposed shares, but then records an equal split of the whole envelope (400m and 614 MiB each) as what executes. It clamps headroom at zero and never fails a proposal over the envelope. It validates with a Go function that stops at the first error, with no schema and no error count. It writes a blocked plan anyway, refuses to plan while a run is open, and takes the model and effort from a pre-recorded choice rather than the plan's suggestion. Its plan fence holds many records that no line or field shows. A team authorizing "1,450m · 2,304 MiB" would get a runtime sized differently.

## Solution

The plan writes `runtime.md` for the ticket's runtime as a draft whose revision is 1 on the first plan. Its `langdrift` fence carries:
- the bound;
- each component's proposed share, port and image;
- one replica each;
- the judge's suggested model and effort, recorded as a recommendation.

Its sizing report carries the proposed totals, the headroom, the suggested resource, the limits statement and the validation result. The proposal is the sum of the components' shares (spec 0066). It must fit the bound (spec 0054), or the plan stops. The draft is validated against the `runtime.md` JSON Schema (spec 0049), and every error is counted. Three lines print, and the Plan card's figures and report fields fill from the same data. What `dft agent run` later applies is what this draft records.

## Responsibility

This element:
- sums the components' proposed shares, checks them against the bound and computes the headroom;
- states the suggested resource (one-replica Deployments, one per component) and the limits;
- builds the `runtime.md` draft, validates it against the schema and counts the errors;
- writes the draft with its revision;
- prints `✓ size`, `limits` and the `validate` tool call;
- owns the Plan card's figures, the CPU and memory bars, the components' total line, BOUNDED, HEADROOM, SUGGESTED RESOURCE, GPU · DISK and VALIDATION.

It does not:
- size each component or resolve its port and image (spec 0066);
- read the bound (spec 0054) or the suggestion (spec 0051);
- hand the draft to the team, state that it is proposed or pending authorization, or show EVIDENCE (spec 0050);
- state the price basis (`estimate runtime has no price source …`, PRICE) or the budget (spec 0056);
- define the fence's fields or the schema (spec 0049).

## Produces

- **`.langdrift/microruntimes/runtimes/<runtime id>/runtime.md`**, the runtime's file at its current location, where `<runtime id>` is the ticket's slug (`book-a-demo-runtime`, spec 0053). It holds:
  - The `langdrift` fence, whose fields are spec 0049's:
    - `runtimeId` (`book-a-demo-runtime`) and `revision` (`1` on the first plan, the previous revision plus one on a later plan of the same runtime);
    - `cpuMillis` (`2000`) and `memoryMiB` (`3072`), the bound;
    - `model` (`claude-sonnet-5-5`) and `effort` (`medium`), the suggestion;
    - `components`, in the spec's order, each with `path`, `stack`, `role`, `port`, `image`, `cpuMillis` and `memoryMiB` (the proposed share, for example `services/signup-api`, `Go`, `api`, `8080`, `250`, `256`) and `replicas` (`1`), the fields spec 0049 defines.
  - The sizing report:
    - `proposed` (`cpuMillis 1450`, `memoryMiB 2304`);
    - `bound` (`2000`, `3072`, set by the team in the design doc);
    - `headroom` (`550`, `768`);
    - `resource` (`Deployment`, count `5`, replicas each `1`);
    - `limits` (`gpu` and `disk` not fields, none requested);
    - `validation` (`schema`, errors `0`);
    - `suggestion` (`by: recommendation`).

  These feed `✓ size`, `limits`, the `validate` line and the Plan card fields below. The run's `intent … runtime.md revision 1` (spec 0050) and `runtime … revision 1 · 5 Deployments ready · 1450m · 2304 MiB` (spec 0062) read the revision and the proposed shares from this file. The draft's revision feeds `✓ plan runtime.md draft, revision 1 · …` and EVIDENCE (spec 0050).
- Into the plan's card data (spec 0050): the same figures, plus the running totals as components resolve.

## Shows

Phase 1, `dft agent plan`. After `console › cross-checking · 5 ports distinct · 5 images resolved · no collisions` (spec 0066):

```
✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares
limits     gpu and disk are not runtime.md fields · none requested
```

then `estimate   runtime has no price source · cpu and memory stay in millicores and MiB` (spec 0056), then:

```
console › tool call: validate runtime.md draft against the schema · 0 errors
```

then `✓ plan       runtime.md draft, revision 1 · proposed to the console` (spec 0050).

Card fields, Plan card, each at its moment:

- **Plan header figures** and **PROPOSED PLAN figures**, stepping up as each component's `✓` line prints (spec 0066): `250m · 256 MiB`, `550m · 768 MiB`, `750m · 1,024 MiB`, `950m · 1,280 MiB`, then `1,450m · 2,304 MiB`. They are not set before the first component, and show a placeholder from the `resolving` line. The words ` · pending team authorization` after them are spec 0050's.
- **CPU** `1,450m of 2,000m` and **Memory** `2,304 MiB of 3,072 MiB`, with bars at the share of the bound, stepping with the figures.
- **Components' total line**: `<n> of 5 resolved · <cpu>m · <memory> MiB` while components resolve (for example `1 of 5 resolved · 250m · 256 MiB`), then `5 components · 1,450m · 2,304 MiB`.
- **BOUNDED**: `2,000m · 3,072 MiB`, from the `bounded` line (spec 0054).
- **HEADROOM**: `550m · 768 MiB`, from `✓ size`.
- **SUGGESTED RESOURCE**: `5 one-replica Deployments`, from `✓ size`.
- **GPU · DISK**: `not runtime.md fields`, from `limits`.
- **VALIDATION**: `schema · 0 errors`, from the `validate` line.

The COMPONENTS rows and the ports line belong to spec 0066. PRICE belongs to spec 0056, and EVIDENCE and "pending team authorization" to spec 0050.

## User Stories

1. As a team lead, I want the plan to size the components inside the envelope our design doc sets, so that the runtime never asks for more than we allowed.
2. As a team lead, I want the size line to show the proposed total against the bound for cpu and memory, so that I see how much of the envelope the plan takes.
3. As a team lead, I want the proposal to be the sum of the components' own shares, not an equal split of the envelope, so that each component gets what its stack needs.
4. As a developer, I want what the run applies to be exactly the shares the plan showed, so that "1450m · 2304 MiB" on the card is what runs.
5. As a team lead, I want the headroom shown as what the bound leaves, so that I know what is spare.
6. As a team lead, I want a proposal over the bound to stop the plan, so that a plan never proposes more than the envelope.
7. As a platform owner, I want the plan to say it proposes one one-replica Deployment per component, so that I know the shape of what will run.
8. As a platform owner, I want the plan to state that gpu and disk are not `runtime.md` fields and that none was requested, so that nobody expects a limit the file cannot hold.
9. As a reviewer, I want the draft validated against the published schema with the number of errors, so that "0 errors" is a measured fact.
10. As a reviewer, I want a draft with errors to stop the plan with no draft written, so that an invalid `runtime.md` never reaches the team.
11. As a developer, I want the first plan of a runtime to be revision 1 and each later plan to raise the revision by one, so that the run can name the revision it applies.
12. As a team lead, I want the draft to record the judge's suggested model and effort as a recommendation, so that the intent the team authorizes starts from the plan's reasoning.
13. As a developer, I want to plan again while an earlier run of the runtime exists, so that no lock stops me from revising a draft.
14. As a console user, I want the Plan card's figures to step up as each component resolves, so that the card follows the terminal.
15. As a console user, I want the components' total to read "n of 5 resolved" while resolving and "5 components" when done, so that progress is visible.
16. As a console user, I want BOUNDED to appear as soon as the bound is stated, before sizing, so that the limit is on the card before the proposal.
17. As a console user, I want HEADROOM, SUGGESTED RESOURCE, GPU · DISK and VALIDATION to fill at their own lines, so that each fact appears when the plan establishes it.
18. As an auditor, I want `runtime.md` to keep only what the plan shows and what the run applies, so that the file is readable and every field has a reason.
19. As an implementer of the cluster, I want each component's share, port, image and replica count in the fence, so that bring-up needs nothing else from the plan.

## Implementation Decisions

**Remove**
- The equal split of the whole envelope across components (400m and 614 MiB each) as the shares recorded and applied. Serves `✓ size … · proposed shares` and the run's `5 Deployments ready · 1450m · 2304 MiB`.
- Clamping headroom at zero. A negative headroom stops the plan instead (see Build).
- Writing a blocked plan, and the `✗ plan … blocked` line with its "blocked" plan state.
- The lock that refuses to plan while a run of the runtime is open.
- Taking the model and effort from a pre-recorded choice input. Serves the draft's recommended model and effort.
- Validation that stops at the first error with a hard-coded count of 1. Serves `· 0 errors`, which is now a count.
- Scanning judge facts for "gpu" or "disk" to say "requested in …". The plan extracts no gpu or disk field (spec 0053), so the limits line states the schema's fact.
- Plan-fence records that feed no line or field and that no run step reads:
  - the collection limits, reach and context hash;
  - the per-source records;
  - the facts and discrepancies;
  - the input records with set-by, via and set-at;
  - the ceiling;
  - the retrieval's method, criteria and scores;
  - the sizing basis note;
  - the price record and mirror path;
  - the check list beyond the schema check;
  - the contract gaps;
  - the fence hash.

  The Ticket and Analysis data the console receives travel in the card data of spec 0050, not in `runtime.md`.

**Build**
- Totals: the sum of the components' shares (spec 0066), and headroom as the bound minus that sum, for cpu and memory.
- The size line: `✓ size       <cpu>m of <bound cpu>m · <memory> of <bound memory> MiB · <n> one-replica Deployments · proposed shares`.
- The resource: one Deployment per component, one replica each, recorded per component as `replicas: 1` and in the report as `Deployment`, count and replicas each.
- The limits line: `limits     gpu and disk are not runtime.md fields · none requested`. It is true because the schema has no gpu or disk field and no extracted fact requests one.
- Schema validation: the draft's fence is validated against the `runtime.md` JSON Schema of spec 0049, every error is counted, and the line `tool call: validate runtime.md draft against the schema · 0 errors` prints only when the count is 0.
- The running totals for the card while components resolve, with the "n of 5 resolved" and "5 components" total line.
- Safe stops:
  - A proposal over the bound in cpu or memory stops the plan before `✓ size`.
  - A validation count other than 0 stops the plan before the `validate` line, since the simulation shows only `· 0 errors`; the count goes into the one sentence on standard error.

  Either way the plan writes no draft and exits non-zero, with the reason through the CLI's error channel (spec 0048). No line the simulation lacks is printed.

**Change**
- The fence's per-component cpu and memory are the proposed shares, and the runtime-wide `cpuMillis` and `memoryMiB` are the bound.
- The fence's model and effort are the judge's suggestion, marked as a recommendation, and stay so: the team's applied choice lives in the recorded intent (specs 0049, 0050, 0057), and the run never rewrites `runtime.md`.
- The draft is written only when its validation count is 0. It is written after the `validate` line and before spec 0050's `✓ plan` line.

**Interfaces with neighbours**
- From spec 0066: the components in the spec's order, with share, port and image, as they resolve.
- From spec 0054: the bound with its provenance. From spec 0051: the suggestion.
- From spec 0049: the fence's fields and the JSON Schema.
- To spec 0050: the draft's path and revision, and the Plan card data. Spec 0050 records the draft's proposed and authorized state.
- To spec 0062: the shares, ports, images and replicas that the cluster brings up.
- To spec 0052: the components with stack and role, and the resource, for the episode the run leaves.

## Testing Decisions

- The seam is the golden test of spec 0048 over the shared fixture. This element adds:
  - its three lines, byte for byte;
  - the Plan card fields above, including every intermediate value of the figures and the total line;
  - the written `runtime.md`, whose fence is checked field by field against the values in Produces and validated against the schema of spec 0049.
- The golden run of `dft agent run` asserts that `runtime … revision 1 · 5 Deployments ready · 1450m · 2304 MiB` (spec 0062) uses the shares from this file.
- Golden cases for the safe stops:
  - a fixture variant whose spec names enough components to exceed the bound;
  - a fixture variant whose draft fails the schema.

  Each prints up to the stop, writes no `runtime.md` and exits non-zero.
- A unit test for the totals and headroom only.

## Out of Scope

- Sizing each component, resolving ports and images, and the COMPONENTS rows and ports line: spec 0066.
- The bound's source and the `bounded` line: spec 0054.
- The suggestion: spec 0051.
- The price basis line, PRICE and the budget: spec 0056.
- `✓ plan … proposed to the console`, "pending team authorization", EVIDENCE and the authorization state: spec 0050.
- The fence's field list and the schema itself: spec 0049.
- Bringing the runtime up and down: spec 0062.
- The console and its stub API.

## Further Notes

- The fence keeps the runtime-wide bound beside the per-component shares on purpose. The bound is the team's limit, and the shares are the plan's proposal inside it. HEADROOM is their difference.
- Re-planning the same ticket writes the next revision over the draft. The simulation shows only revision 1, so nothing about a second revision is printed differently.
