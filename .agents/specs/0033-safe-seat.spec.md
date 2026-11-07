# Safe Seat: the person writes a reserved part of every ticket

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

> **Note, 2026-11-05.** The owner closed 0207 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

## Problem Statement

A team that delivers with coding agents through `ex` can hand an agent the whole of every ticket. The work ships, but the people who own it stop writing it. They stop practicing the reasoning that lets them read, debug and explain the code at the moment an incident needs them. Bainbridge described this irony of automation in 1983. A 2026 randomized trial found that developers who learned a library with an AI assistant scored lower on a quiz of concepts, code reading and debugging, with no significant gain in speed. Developers also misjudge how much AI tools help them, so their own impression cannot settle the question. The evidence, and its limits, is in the research note `safe-seat.research.md`.

Nothing in `ex` or the console keeps a person in the work. A `human` stage means a person approves and the agent then runs the stage. No stage is driven by a person. `ex` records no authorship, so nobody can tell how much of a delivered ticket a person wrote. A team lead who wants the team to keep its skill has no mechanism, no reading and no record.

## Solution

Safe Seat is a capability of the LangDrift FinOps harness. When it is on, every ticket reserves a part that a person writes. The part is a problem with its own goal and acceptance criteria, in one of three categories: core domain logic, the hard part, or first contact with a new stack. Boilerplate, configuration and generated code never qualify.

The flow, ticket by ticket:

1. **Plan.** `ex agent plan` suggests the reserved part, located by paths, and the lead authorizes it with the plan.
2. **Agent stages.** `ex agent run` drives the coding agent through the workflow. The agent writes only the boundary the reserved part needs (its signature or interface), outside the reserved paths.
3. **Dev ticket.** At the workflow's person stage, `ex` stops calling the agent. The person writes the solution and its tests, then closes the dev ticket and says whether they wrote it without an agent.
4. **Review.** After the dev ticket closes, the agent may review the work.

`ex` snapshots the working tree around every agent call. When the ticket is delivered, every line of the final diff is attributed to the agent, the person or unknown. The person's share must reach the floor, at least 5% of the final diff's lines, or a higher floor the team sets. Below the floor, the ticket needs an override from the tech lead, with a reason.

Each layer of the policy has its own owner:

- **LangDrift:** the 5% floor, a product principle that no team can lower.
- **Executive:** turns Safe Seat on or off.
- **Tech lead:** raises the floor, picks each ticket's reserved part and decides overrides.
- **Person:** closes their own dev ticket.

Visibility follows the same split. The person sees their own numbers. The rest of the team sees only each ticket's state and the team's skill drift, the share of tickets that needed an override, never a ranking.

The name is Safe Seat, never abbreviated. The terms dev ticket, reserved part and skill drift are working names that stay out of the glossary until the owner settles them.

## User Stories

### Policy

1. As an executive, I want to turn Safe Seat on for a project, so that every ticket planned there reserves a part a person writes.
2. As an executive, I want turning Safe Seat on to start at the 5% floor, so that the team has a working policy without choosing a number first.
3. As an executive, I want to be the one who turns Safe Seat off, recorded with my name, so that switching off the practice is a deliberate act of leadership and not a developer's shortcut.
4. As a tech lead, I want to raise the team's percentage above 5%, so that a team that wants more practice can commit to it.
5. As a tech lead, I want a percentage below 5% refused with the reason, so that no team can go below LangDrift's principle.
6. As a tech lead, I want to add exclusion patterns for files that do not count, so that a team's own generated code does not dilute the measurement.
7. As a tech lead, I want the product's exclusions (lockfiles and files marked as generated) to stay in force whatever the team adds, so that nobody inflates a share by excluding the code that matters.
8. As an auditor, I want every change to the policy recorded with who set it, when, through what and the record it replaced, so that the policy has the same history as the budget.
9. As a team member, I want the policy mirrored to a read-only file, so that I and the console can read it, knowing that no gate reads the mirror.
10. As a team lead, I want a policy change to apply to plans authorized after it, so that a ticket already running keeps the rules it was authorized under.
11. As an auditor, I want each run to record which policy record it was authorized under, so that its state can be checked against the rule in force.
12. As a team lead initializing a project at a terminal, I want `ex init` to offer Safe Seat, defaulting to on, so that a new project starts with the practice.
13. As a script author, I want `ex init` without a terminal and without the Safe Seat flag to record nothing about Safe Seat, so that existing automation behaves as today.
14. As a team member, I want a project with no Safe Seat record to say so in the plan, so that a missing policy reads as missing, never as off by decision.

### The reserved part

15. As a tech lead, I want `ex agent plan` to suggest a reserved part with its paths, goal, acceptance criteria, category, reason and confidence, so that I start from a reasoned proposal.
16. As a tech lead, I want the suggestion to cover a problem in one of the three categories, so that the person exercises reasoning, not typing.
17. As a tech lead, I want boilerplate, configuration, generated code, renames and tool-run migrations never suggested, so that the reserved part is never the dull leftover.
18. As a tech lead, I want the suggestion's planned size to reach the team's percentage of the ticket's planned size, so that the plan is on track for the floor before any code is written.
19. As a tech lead, I want the plan to say plainly when no part of the ticket meets the criteria, so that I choose one myself or override with a reason.
20. As a tech lead, I want to replace the suggestion with paths, a goal, acceptance criteria and a category of my own, so that the decision stays mine.
21. As a tech lead, I want my own pick checked against the planned size and shown as met, below or unknown, so that I see where it stands, knowing that delivery decides.
22. As a tech lead, I want authorizing the plan to authorize its reserved part, recorded with my name, so that the part is an explicit input and not something inferred from prose.
23. As a tech lead running a plan without a judge, I want the reserved part left empty with the reason, so that a missing credential never invents a suggestion.
24. As a team member, I want the Plan card to show the reserved part and who chose it, so that everyone sees what the person will write.
25. As a tech lead, I want a run with Safe Seat on and neither a reserved part nor an override refused at authorization with nothing paid, so that no ticket starts without its part.
26. As a tech lead, I want to override the reserved part with a reason when the ticket has none worth writing by hand, so that a trivial ticket does not block the team.

### The run

27. As a team lead, I want the workflow to declare one stage driven by the person, so that the team places the dev ticket where it fits its harness.
28. As a team lead, I want a plan with Safe Seat on and no person stage in the workflow to fail its checks with guidance, so that the gap shows before the run.
29. As a team member, I want the person stage skipped when Safe Seat is off, so that the workflow runs unchanged without the policy.
30. As a developer, I want every agent call before the dev ticket to receive the reserved paths with the instruction to leave them alone, so that the agent writes only the boundary.
31. As a developer, I want `ex` never to call the agent while the dev ticket is open, so that the reasoning in it is mine.
32. As a developer, I want the run to pause right after an agent call that changes a reserved path, so that an agent that ignores the boundary is caught at once.
33. As a developer, I want to restore the reserved paths from the snapshot taken before that call, recorded, so that I get my problem back intact.
34. As a developer, I want to accept the agent's change instead, recorded, so that the run continues and those lines count as the agent's.
35. As a developer, I want the run to stop at the dev ticket and show its paths, goal and acceptance criteria, so that I know exactly what to write.
36. As a developer, I want to write the solution and its tests myself, so that I frame and solve the problem, not just make given tests pass.
37. As a developer, I want to close the dev ticket with `ex agent confirm` and say whether I wrote it without an agent, so that my answer decides whether my lines count as mine.
38. As a developer, I want a terminal prompt to ask that question when I close the dev ticket, so that I never attest by accident.
39. As a script author, I want `--confirm` never to close a dev ticket, so that no automation attests for a person.
40. As a developer, I want closing without a terminal and without an attestation flag refused with the reason, so that a missing answer is never read as yes.
41. As a developer, I want my answer that I used an agent to make those lines unknown, so that honesty costs me only the share, never an accusation.
42. As a developer, I want lines I wrote by hand earlier in the same run to count once I attest them, so that my work outside the dev ticket is not lost.
43. As a tech lead, I want the dev ticket never skipped because of earlier lines, so that hand-written boilerplate cannot replace the reasoning.
44. As a developer, I want the agent allowed to review after I close the dev ticket, so that I get feedback once my reasoning is done.
45. As a developer, I want any change the agent applies to my lines to count as the agent's, so that the share reflects the delivered code.
46. As a FinOps owner, I want the dev ticket to make no paid call and the estimate to exclude the reserved part, so that the money reading matches the work the agent does.
47. As a FinOps owner, I want the reservation held while the dev ticket is open, as for any waiting stage, so that the money gates behave as they do today.

### Measurement

48. As a developer, I want `ex` to snapshot the working tree before and after every agent call without touching my index, worktree, branches or refs, so that measuring never disturbs my work.
49. As a team member, I want the share measured on the lines the ticket's final diff adds, from the run's base commit to the delivered head, so that the number describes what ships.
50. As a team member, I want each counted line attributed to the agent when an agent call introduced it, to the person when an attested interval did, and to unknown otherwise, so that nothing unobserved is credited to anyone.
51. As a team member, I want unknown lines kept in the denominator, so that unknown never inflates the person's share.
52. As a team member, I want excluded files left out of the count, so that lockfiles and generated code do not move the share.
53. As a team member, I want a ticket whose final diff counts no line to read as not applicable, never as 0%, so that a documentation ticket is not a failure.
54. As a developer, I want `ex agent status` to show my agent, person, unknown and excluded lines, my share and the floor, so that I see where I stand.
55. As a developer, I want those numbers kept in the project's private state, outside the committed records, so that my share stays mine.

### Gate and override

56. As a tech lead, I want a ticket below the floor at delivery to be marked as needing an override before the owner's stage, so that the gap is visible before anyone merges.
57. As a tech lead, I want to override it with a reason, recorded with my name, so that a justified exception leaves a trace.
58. As a tech lead, I want the measured share shown to me only when I decide an override, so that the number reaches another person only when a decision needs it.
59. As an auditor, I want the override's mirror to carry the state and the reason but not the share, so that the committed record does not expose a person's number.
60. As a tech lead, I want it said plainly that `ex` records a merge without an override but cannot prevent it, so that nobody mistakes the gate for a lock on the repository.
61. As a reviewer of a pull request, I want the CI Gates report to show Safe Seat as unknown, with the reason that CI has no Safe Seat source, so that a missing reading is never read as met.

### Console

62. As a team member, I want the Workflow card to mark the person stage as driven by a person, distinct from a stage a person approves, so that the two meanings never blur.
63. As a team member, I want each ticket's Safe Seat state shown as met, needs override, overridden with its reason, not applicable or unknown, so that I read the practice ticket by ticket.
64. As a team member, I want no ticket's share shown in the console, so that the console never becomes a ranking.
65. As a tech lead, I want the FinOps area to show the month's skill drift for the team (applicable tickets, met, overridden and the share overridden), with the drift colors, so that the practice reads beside plan, intent and cost drift.
66. As a team member, I want tickets `ex` did not run left out of skill drift, so that the reading covers only what the platform observed.
67. As a team member, I want the console's stub data for Safe Seat labeled as illustrative, as the rest of the stub is, so that a seeded reading is never taken for a measured one.
68. As a visitor using any of the five locales, I want the Safe Seat copy in my language, with the name Safe Seat untranslated, so that the capability reads the same everywhere.

### Dogfooding

69. As the owner of this repository, I want a Safe Seat policy for LangDrift's own work, so that the project is the first to practice it.
70. As the owner, I want the human-in-the-loop policy to name the dev part as a checkpoint of mine, so that the automatic stages stop for it.
71. As the owner, I want `to-tickets` to write a dev part into every ticket, with its paths, goal and acceptance criteria, so that each ticket carries my part.
72. As the owner, I want `implement` to leave the dev part's paths alone and stop when only they remain, so that I write them.
73. As the owner, I want `code-review` to review my dev part after I write it, so that I get feedback after the reasoning.
74. As the owner, I want an override of a ticket's dev part written in the ticket with its reason, so that skipping one is visible until the runtime measures it.

## Implementation Decisions

### Policy

- **The floor is a product constant.** The Safe Seat floor of 5% is a constant of the harness domain. Validation refuses any percentage below it or above 100.
- **One input record kind, `safe-seat`.** Its value is `enabled`, `percent` (an integer from the floor to 100) and `exclusions` (the team's path patterns).
  - It is set through `ex harness safe-seat set` with `--by`, recorded and mirrored like every other input record.
  - Turning Safe Seat on without a percentage records the floor.
  - Patterns are project-relative globs. An absolute pattern, an empty pattern or one that climbs out of the project is refused.
  - The roles (executive, tech lead, person) are recorded through `--by` only. Nothing authenticates them until the account exists.
- **Product exclusions.** These always apply, whatever the team adds:
  - lockfiles from a fixed list;
  - files with the `linguist-generated` attribute in the project's `.gitattributes`;
  - files whose header carries the `Code generated … DO NOT EDIT` marker.
- **Read at authorization.** `ex agent run` reads the `safe-seat` record when it authorizes the plan. The run record names the record it used. A later change applies only to plans authorized after it.
- **`ex init`.** At a terminal, `ex init` offers Safe Seat, defaulting to on, and records it as `ex harness safe-seat set` would. A flag gives the same answer without a terminal. Without a terminal and without the flag, `init` records nothing about Safe Seat (spec 0024 keeps non-interactive `init` unchanged). A project with no `safe-seat` record shows "Safe Seat not recorded" in the plan and runs without it.

### The reserved part

- **Shape.** A reserved part has `paths`, a `goal`, `acceptance` criteria, a `category` (`core-domain`, `hard-part` or `new-stack`), `plannedLines`, the ticket's `totalPlannedLines`, an `excludedKind`, a `reason`, a `confidence` and `chosenBy` (the judge's suggestion, or the lead). `excludedKind` is the judge's assessment that the part is boilerplate, configuration, generated code, a rename or a tool-run migration; it is null otherwise, and always null for the lead's pick.
- **The judge suggests.** The judge suggests a part only if its planned share reaches the team's percentage and its `excludedKind` is null. When no part qualifies, the plan records the reason. Without a judge, the part stays null with the reason, as other judged fields do.
- **The lead can replace it.** `ex agent plan` takes the lead's own pick (paths, goal, acceptance criteria, category), which replaces the suggestion. The plan check on planned share reads met, below or unknown for the lead's pick and never blocks it.
- **Where it is stored.** The reserved part lives in the plan record and in the `langdrift-plan` fence, and the Plan card reads it there. Authorizing the plan authorizes it.

### Gates

Three new gates, with ids that use the word, never the initials:

| Gate | When | What it does |
| --- | --- | --- |
| `SEAT-01` reserved-part-required | At authorization | With Safe Seat on, the run is refused, with nothing paid, unless the plan has a reserved part or the ticket has an override. |
| `SEAT-02` reserved-path-untouched | After every agent call | If the call changed a reserved path, the run pauses. |
| `SEAT-03` floor-met | At delivery | Below the floor, the ticket's state becomes `needs-override`. |

- **Overrides.** An override is an input record of kind `safe-seat-override`, set with `ex agent safe-seat override <runtime-id> --reason … --by …`. Its value holds the ticket, the run, the lap, the gate it answers (`SEAT-01` or `SEAT-03`) and the reason. For `SEAT-03` the private record also keeps the measured share. Its mirror carries the state and the reason, never the share.
- **After a `SEAT-02` pause:**
  - `ex agent safe-seat restore <runtime-id>` writes the reserved paths back from the snapshot taken before the offending call, and removes reserved files that call created.
  - `ex agent safe-seat accept <runtime-id>` keeps the change.
  - Both record who and when in the lap's evidence and resume the run.
- **`SEAT-03` records, it does not prevent.** It never prevents a merge on the hosting service. `ex` cannot, and the spec does not pretend otherwise. The ticket's latest lap decides its state: an override answers the lap it names, a later lap that meets the floor makes the ticket `met`, and a latest lap that was not observed makes it `unknown`, whatever an earlier lap showed.

### Workflow

- **A person-driven stage.** A workflow stage gains `drivenBy`, either `agent` (the default) or `person`. A workflow may declare at most one person stage.
- **Behavior at the person stage.**
  - With Safe Seat on, the play loop reaches the person stage, takes the opening snapshot, sets the stage and the run to `waiting`, and prints the reserved part. It never calls the stage driver or runs skills there; the stage's hooks still run.
  - With Safe Seat off, the person stage is skipped.
  - With Safe Seat on and no person stage, the plan fails a check that names the missing stage.
- **The sample workflow** places the person stage, `dev-ticket` labeled "Dev ticket", after `implement` and before `review`. Its review stage carries `code-review`, so the agent reviews after the person.
- **Reserved paths reach the agent.** Every agent call receives the reserved paths through a new `{{reserved}}` placeholder and a `LANGDRIFT_RESERVED_PATHS` variable. The driver's fixed prompt tells the agent to write only the boundary outside those paths.
- **Closing the dev ticket.**
  - `ex agent confirm <runtime-id> dev-ticket` closes the person stage only with an attestation: `--without-agent` or `--with-agent`, or the terminal prompt asking "Written without an agent?". The answer and the `--by` are recorded in the lap's private state, never in the committed records, and the closing snapshot is taken.
  - `--confirm` on `ex agent run` never closes a person stage.
  - Closing without a terminal and without an attestation flag is refused with the reason.

### Measurement

- **A new port: tree snapshots.**
  - A snapshot is a git tree of the working tree, tracked and untracked files minus ignored ones, written through a private index. The person's index, worktree, branches and refs are never touched; the only side effect is unreferenced objects in the repository's object store.
  - The run records its base commit at start.
  - Snapshots are taken before and after every agent call, and when the dev ticket opens and closes.
  - Each agent call's stage evidence gains the paths it changed, with lines added and removed, and the two tree ids.
- **The intervals.** The snapshots split the run into intervals. An interval inside an agent call belongs to the agent. Every other interval, including the dev ticket and the time between stages, belongs to the person if the dev ticket's attestation is "without an agent", and to unknown otherwise.
- **Attribution.**
  - At delivery, every line the final diff adds is attributed to the last interval that introduced it. The final diff runs from the base commit to the delivered head, minus exclusions, and a modified line counts once.
  - The attribution chains the snapshot trees into private, unreferenced commits and asks git which commit last introduced each line. Lines changed after the last snapshot belong to a final interval treated like any other interval outside an agent call.
  - A line the chain cannot trace is unknown.
- **Share and state.**
  - The person's share is person lines over all counted lines, with unknown lines in the denominator.
  - A diff that counts no line makes the ticket `not-applicable`.
  - The ticket's state is one of `met`, `needs-override`, `overridden`, `not-applicable` or `unknown`. It is `unknown` when delivery was not observed or the base commit cannot be read.
- **Where the numbers live.**
  - The counts (agent, person, unknown, excluded), the share and the attestation live in the project's private state, outside the committed records. `ex agent status` prints them.
  - The committed run record carries only the Safe Seat state, the policy record it used, and the override record when there is one.

### Records, documentation and CI

- **Contracts README.** It documents the new fields of the plan fence, `run.json`, the lap evidence and the two input record kinds.
- **`ex` README.** It documents `drivenBy`, the person stage, the `safe-seat` commands and the three gates.
- **CI Gates.** `ex gates report` adds a Safe Seat line marked unknown, with the reason that CI has no Safe Seat source, as it does for the FinOps readings.

### Console and extub

- **The console reads the stub.** The console still reads extub (the connector of spec 0017 is not built), so extub seeds LangDrift Studio's tickets with Safe Seat data, labeled illustrative like the rest of the stub:
  - a person stage in each workflow;
  - a reserved part on each plan;
  - a Safe Seat state per ticket, covering every state value, from which the console computes the month's skill drift.
- **Ticket types.**
  - `TicketStage` gains `drivenBy`, and the Workflow card marks a person-driven stage apart from a human-approved one.
  - The ticket gains its Safe Seat state, and the override reason when there is one.
  - The plan gains the reserved part.
  - No type carries a share.
- **Skill drift.** A pure compute counts, per month, the applicable tickets (`met` plus `overridden`), the met and the overridden ones, and the share overridden. The FinOps area shows it with the drift colors and leaves `not-applicable` and `unknown` tickets out of the share.
- **Localization.** The console copy goes into all five catalogs (`en`, `pt-BR`, `zh-Hant`, `ja`, `vi`) with key parity. The name Safe Seat stays untranslated. *(Amended 2026-11-05: the console serves four locales, `en`, `pt-BR`, `zh-Hant` and `ja`; the owner closed 0308, `vi` for the console, as won't fix.)*

### Dogfooding

- **The policy.** A Safe Seat policy for this repository goes into `.agents/policy/`, beside the human-in-the-loop policy:
  - every ticket has a dev part of at least 5% of its diff, written by the owner;
  - `to-tickets` proposes the dev part, and the owner decides it;
  - an override is written in the ticket with its reason;
  - nothing measures the share until the runtime exists.
- **The human-in-the-loop policy** names the dev part as a third owner checkpoint. The automatic stages stop for it and resume with `code-review` and delivery.
- **The skills.**
  - `to-tickets` writes a dev part section into every ticket, with its paths, goal, acceptance criteria and category, or an override with its reason.
  - `implement` leaves those paths alone and stops when only they remain.
  - `code-review` reviews the dev part once the owner has written it.

## Testing Decisions

A good test drives `ex` the way a person does, through its commands, and checks what the person sees, what the records hold and what the working tree contains. It never asserts internal structure.

1. **The seam is the `ex` command in process, with injected dependencies.** This is the prior art of the run, stage and delivery tests: `exWith`, the scripted driver, the scripted delivery and the seeded project. The seeded project becomes a git repository, so snapshots and attribution run against real git with no fake. The scripted driver simulates the agent by writing lines into the project. The test edits files between commands to simulate the person.
2. **Policy.**
   - `ex harness safe-seat set` records and mirrors the policy.
   - It refuses a percentage of 4 or 101, and absolute, empty and climbing patterns.
   - Turning Safe Seat on without a percentage records 5.
   - A second record keeps the first as its previous.
   - `ex init` at a terminal offers Safe Seat; without a terminal and without the flag it records nothing.
   - Prior art: the inputs and init tests.
3. **Plan.**
   - A scripted judge's suggestion lands in the plan fence.
   - A suggestion below the team's percentage is not offered, and the plan says why.
   - The lead's pick replaces the suggestion and shows its planned-share check.
   - `--no-judge` leaves the part null with the reason.
   - With Safe Seat on, a workflow without a person stage fails its check.
   - Prior art: the plan and judged-plan tests.
4. **Run.**
   - `SEAT-01` refuses a run with nothing paid, and an override lets it start.
   - The play loop never calls the scripted driver in the person stage.
   - A driver call that writes a reserved path pauses at `SEAT-02`. `restore` brings back the before-snapshot content and removes created files; `accept` keeps the change.
   - `--confirm` never closes the dev ticket, and a non-terminal close without an attestation flag is refused.
   - The dev ticket closes with either attestation.
   - With Safe Seat off, the person stage is skipped.
   - Prior art: the stage tests.
5. **Measurement.** One scenario covers all four sources of a line: the agent writes the boundary and other files, the person writes the reserved paths and attests "without an agent", the agent's review rewrites one person line, and an excluded lockfile changes. At delivery:
   - `ex agent status` prints the expected counts, share and state;
   - the committed run record holds only the state;
   - the person's index, worktree and refs are unchanged.

   Variants:
   - "with an agent" makes the person lines unknown;
   - a diff of only excluded files is `not-applicable`;
   - below the floor, the state is `needs-override`, and an override makes it `overridden`, with a mirror that holds no share.

   Prior art: the delivery tests.
6. **CI Gates.** The report prints the Safe Seat line as unknown. Prior art: the gates report test.
7. **Console logic** (Node test runner, beside the cost totals and ticket phase tests):
   - the skill drift compute for a month that mixes all five states;
   - the mapping of a stub ticket's Safe Seat state and reserved part;
   - the Workflow card's stage marker logic, for person-driven versus human-approved stages.
8. **extub.** The catalog test checks that the seeded tickets cover every Safe Seat state.
9. **Repository checks:**
   - `pnpm turbo run build test typecheck` for `ex`, the contracts, extub and the console, with Go's race detector;
   - `pnpm turbo run check`, `drift quality check` and `drift check localization`;
   - a diff of the leaf keys of the console catalogs (four since 2026-11-05).

## Out of Scope

Each item becomes its own work once grilled:

- The website's Safe Seat section and copy, now covered by spec 0043.
- The account, authenticated roles and an organization above the project. Until then, who acts is recorded but not enforced.
- The Exphere server as the policy's source of truth, and a required CI check that reads it and truly blocks a merge below the floor.
- The connector of spec 0017. The console reads Safe Seat from extub until it lands.
- A mode of `ex` in which the agent only explains, and any observation of IDEs or agents running outside `ex`.
- Reserving a part by symbol instead of by path.
- Measuring the dogfooding policy before the runtime exists.
- Glossary entries for dev ticket, reserved part and skill drift, which wait for the owner to settle their names.

## Further Notes

- The owner settled these decisions in a grilling on 2026-11-03. The research note `safe-seat.research.md` holds the evidence and its limits. The evidence does not show that agent-assisted work causes depression or anxiety, so Safe Seat makes no clinical claim and its copy avoids health vocabulary.
- **On by default.** Only an executive turns Safe Seat off, so an interactive `ex init` offers it on by default; the owner confirmed this on 2026-11-03. A project without a record runs without it, which keeps the rule that every value that can block a run comes from explicit input.
- The 5% floor is a product principle, not a threshold derived from a study, and the FAA advises against practice quotas without evidence. The floor confirms that a reserved part was really written; it is not a dose of practice. What makes the practice is the reserved part, defined by its problem rather than by its line count, and the lead can override it with a reason.
- Not calling the agent at all in the dev ticket is stricter than the evidence requires: in the skill-formation trial, asking an assistant conceptual questions preserved learning. Safe Seat keeps the strict rule because it is what `ex` can enforce and measure cleanly. A mode in which the agent only explains stays out of scope.
- The gate ids use `SEAT`, and the name is never shortened to its initials anywhere, in any locale.
