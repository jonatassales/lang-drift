# The run delivers: push, pull request and ready

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E16.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 17 of 19, after spec 0061 and before spec 0062. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

In the simulation, once the implementation is complete, the run delivers it: `✓ pushed     feat/0031-book-a-demo-runtime · 3f9c2e1`, then `✓ pr #19     feat(demo): run book a demo as a microruntime`, and at the very end `✓ ready      #19 for review · $1.60 released to the envelope`, or `… to merge …` when Review is not ticked. The Workflow card's header turns into `#19 · ready for review · $9.80`, and its PULL REQUEST field reads `#19 · 3f9c2e1 · ready for review`.

The LangDrift CLI only observes. It reads the checked-out branch with `git` and calls it pushed only when an upstream holds it with no commit ahead; otherwise it prints `pushed     unknown · …`. It asks `gh` for an existing pull request and prints `pr         unknown · …` when there is none. It chooses "for review" over "to merge" by whether delivery looked clean, not by the Review choice, and writes a proposed ticket status that nothing shows. Nothing creates the branch, pushes it, opens the pull request or marks it ready, so the simulation's delivery never happens unless someone else does it.

## Solution

The run delivers the work itself, with the repository's own tools, `git` and `gh`:

- Before the first stage, it starts the branch `feat/0031-book-a-demo-runtime` from the checked-out commit; the stages and the person work on it.
- After `✓ completed …`, it pushes the branch and prints its name and the pushed commit's short hash.
- It opens a draft pull request titled with the branch's first commit subject and prints its number and title.
- At the end, after the runtime is stopped, it marks the pull request ready and prints `for review` or `to merge`, as the Review choice says, with the released figure.

Who delivers is the owner's to confirm (see Open Questions); this is the minimal mechanism that produces the shown lines.

## Responsibility

This element creates the run's branch, pushes it, opens the pull request, marks it ready, prints the three delivery lines and fills the Workflow card's phase-2 header and PULL REQUEST.

It does not:

- write or commit the code: the stages (specs 0058, 0059) and the person (spec 0061) do;
- follow the checks (spec 0064) or decide Safe Seat (spec 0061);
- compute the observed or released figures (spec 0056);
- decide the Review choice or print the `# next` lines (spec 0057);
- merge: the merge is always the owner's, and nothing in `dft` merges.

## Produces

- **The branch**, in the client's repository: `feat/<ticket number>-<ticket slug>` (`feat/0031-book-a-demo-runtime`, from spec 0053's number and slug), started from the checked-out commit before the first stage, and pushed to the `origin` remote with its upstream set.
- **The pull request**, on the forge: a draft from that branch to the repository's default branch, titled with the subject of the branch's first commit, marked ready at the end. Never merged.
- **The delivery, in the run's record (spec 0065)**:
  - `branch` (`feat/0031-book-a-demo-runtime`), feeding the pushed line;
  - `commit`, the pushed head's seven-character short hash (`3f9c2e1`), feeding the pushed line and PULL REQUEST;
  - `pullRequest.number` (`19`), feeding the pr line, the ready line, the header and PULL REQUEST, and the `# next` line (spec 0057);
  - `pullRequest.title` (`feat(demo): run book a demo as a microruntime`), feeding the pr line;
  - `pullRequest.state`: `draft`, then `ready for review` or `ready to merge`, feeding the header and PULL REQUEST. While the pull request is a draft, PULL REQUEST's state shows the checks' state that spec 0064 records beside it (`checks running`, then `4 of 4 checks passed`).

## Shows

Phase 2 only; none of these print in the over-budget variant.

After `✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record` (spec 0065):

```
✓ pushed     feat/0031-book-a-demo-runtime · 3f9c2e1
✓ pr #19     feat(demo): run book a demo as a microruntime
```

Then come `checks …` and `✓ checks …` (spec 0064), `✓ Safe Seat  met` (spec 0061) and `runtime    stopped · …` (spec 0062). After the stopped line, with Review ticked:

```
✓ ready      #19 for review · $1.60 released to the envelope
```

With Review unticked:

```
✓ ready      #19 to merge · $1.60 released to the envelope
```

All three are success lines. The released figure is spec 0056's; it is the run's reserved minus observed, so it changes with the model and effort. The `# next` line follows (spec 0057).

On the Workflow card, PULL REQUEST is a dash until the pr line, and then fills in four steps, as the simulation shows:

- from `✓ pr #19 …`: `#19 · 3f9c2e1`, with the pull request's mark;
- from `checks     … running` (spec 0064): its state reads `checks running`;
- from `✓ checks     4 of 4 passed` (spec 0064): its state reads `4 of 4 checks passed`, in the success tone;
- from the ready line: its state reads `ready for review`, or `ready to merge`.

At the ready line, the header becomes `#19 · ready for review · $9.80`, or `#19 · ready to merge · $9.80` (the observed figure is spec 0056's).

## User Stories

1. As a developer, I want the run to work on its own branch named after the ticket, so that the delivery never lands on the branch I was on.
2. As a developer, I want the branch named `feat/<number>-<slug>`, so that I can tell from the name which ticket it delivers.
3. As a developer, I want the branch pushed once the implementation is complete, so that the work leaves my machine without a step of mine.
4. As a team lead, I want the pushed line to show the branch and the commit's short hash, so that I know exactly what was delivered.
5. As a team lead, I want the run to open the pull request, so that the review starts where the team reviews.
6. As a team lead, I want the pr line to show the number and the title, so that I can find it.
7. As a developer, I want the pull request's title to follow the project's commit convention, so that it reads like the team's own.
8. As a team lead, I want the pull request opened as a draft and marked ready only at the end, so that nobody reviews it before the checks pass and Safe Seat is met.
9. As the owner, I want the ready line to say "for review" when Review is ticked and "to merge" when it is not, so that the wording follows the choice the team made.
10. As the owner, I want the ready line to show what the reservation released, so that the money closes on the same line as the delivery.
11. As the owner, I want nothing in `dft` ever to merge, so that the merge stays mine.
12. As a console developer, I want the branch, the commit, the number, the title and the state in the run's record, so that the header and PULL REQUEST fill from data.
13. As a developer, I want a push or a pull request that fails to stop the run without a made-up line, so that the terminal never says `unknown`.
14. As a developer, I want a stopped run to leave its pull request as a draft, so that an unfinished delivery is never presented as ready.
15. As the owner, I want no proposed ticket status written, so that the run keeps only what the simulation shows.
16. As an auditor, I want the pushed hash recorded, so that the delivered commit can be traced to the run.

## Implementation Decisions

### Remove

- The read-only observer as the source of the delivery lines: reading the checked-out branch, the "upstream holds it with no commit ahead" rule, and looking up an existing pull request by branch. Serves the pushed and pr lines, which now report what the run did.
- `pushed     unknown · …` and `pr         unknown · …`, and the "the branch" stand-in for a missing number in the ready and `# next` lines. Serve nothing shown.
- Choosing "for review" when delivery is not clean (the delivered-cleanly rule). Serves the ready line, which follows the Review choice (spec 0057).
- The proposed ticket status (`in-review`, `ready-to-merge`) written to the run's record. Serves nothing shown.
- The delivery words of a lap's steps ("… at …, not pushed, pull request #…"). Serve nothing shown; laps are removed (spec 0065).

### Build

- Starting the branch `feat/<number>-<slug>` from the checked-out commit, after the gate opens and before the first stage, with no line. When the branch already exists from an earlier run of the same ticket, the run checks it out and works on it instead. Serves the pushed line's branch.
- After `✓ completed …`: `git push` of the branch to `origin` with its upstream set, then the short hash of the pushed head, seven characters. Serves the pushed line.
- Opening a draft pull request with `gh` from the branch, titled with the subject of the branch's first commit, or, when the branch already has an open pull request from an earlier run, using that one. Serves the pr line and PULL REQUEST's first step.
- After `runtime    stopped · …`: marking the pull request ready with `gh`, then the ready line with the number, the wording from spec 0057's `reviewed`, and spec 0056's released figure. Serves the ready line.
- Recording the delivery in the run's record at each of these moments. Serves the header and PULL REQUEST.
- A forge port with four operations (start the branch, push, open a draft pull request, mark it ready), so that the golden test scripts them.

### Change

- The delivery changes from observed to performed by the run (see Open Questions).
- The ready line's wording depends only on the Review choice.
- The pushed line now always carries the hash.

### Interfaces

- From spec 0053, through the plan's card data: the ticket's number and slug.
- From spec 0065: the moment `✓ completed …` has printed. To spec 0064: the pull request's number, once the pr line has printed.
- From spec 0062: the moment the stopped line has printed.
- From spec 0064: the checks' state recorded beside the pull request, which PULL REQUEST shows while the pull request is a draft.
- From spec 0057: `reviewed`. From spec 0056: `releasedCents` for the ready line, `observedCents` for the header.
- To spec 0057: the pull request's number for `# next: the owner reviews #19 and decides whether to merge`.
- To spec 0065: the delivery fields for the run's record.

### Failure paths

Every one follows spec 0048's stop rule, with no line of its own, and the runtime, when up, is taken down silently (spec 0062):

- the branch can be neither started nor checked out: the run stops before the first stage;
- the push is refused or `origin` is missing: the run stops before the pushed line;
- `gh` is missing or not signed in, or the pull request cannot be opened: the run stops before the pr line, with the branch pushed;
- the pull request cannot be marked ready: the run stops before the ready line, and the pull request stays a draft.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test**, with the forge port scripted: the branch started, the push reporting `3f9c2e1`, pull request #19 opened as a draft, titled with the subject of the branch's first commit, `feat(demo): run book a demo as a microruntime`, which the fixture's fake coding agent writes (spec 0048), and marked ready. This spec adds:
  - the pushed, pr and ready lines, byte for byte, at their positions, with `for review` in the default golden and `to merge` in the review-unticked golden;
  - the branch started before the first stage line, and no forge operation in the over-budget variant;
  - the delivery fields in the run's record, compared with the Workflow header and PULL REQUEST after the run, and PULL REQUEST's four steps at their lines;
  - no merge operation in any variant.
- Golden stops: a refused push and a pull request that cannot be opened. Each leaves standard output a prefix of the golden, takes the runtime down, writes one sentence to standard error and exits 1.
- Unit tests only for the rules: the branch name from a ticket's number and slug.
- Prior art: the existing delivery observer tests with a scripted process runner.

## Out of Scope

- The commits' content and authorship (specs 0058, 0059, 0061).
- The checks (spec 0064), Safe Seat (spec 0061), the runtime (spec 0062), money (spec 0056), the `# next` lines (spec 0057).
- Forges other than the one `gh` serves.
- The console and its stub API.

## Open Questions

- **Who delivers.** The simulation shows the branch pushed and the pull request opened as steps of the run, but not by whom. Until the owner decides, `dft` starts the branch, pushes it, opens the pull request and marks it ready; the alternative is that the coding agent does these and `dft` only reports them.
- **Where the pull request's title comes from.** `feat(demo): run book a demo as a microruntime` is not derivable from the ticket's title. Until the owner decides, it is the subject of the branch's first commit, which the stages write in the project's commit convention.
- **A second run of the same ticket.** The simulation shows one run on a fresh branch. A stopped run's intent stays runnable (spec 0048), and its branch, and perhaps its draft pull request, remain. This spec reuses them: the next run checks the branch out and uses the open pull request. The alternatives are a fresh branch name per run, or a stop.
- **Work the stages did not commit.** If the stages or the person leave changes uncommitted, the push carries only the committed head. Until the owner decides, `dft` commits nothing itself.
- **A draft pull request.** Opening the pull request as a draft and marking it ready at the `ready` line is this spec's choice; the simulation shows only `✓ pr #19 …` and then `✓ ready …`.

## Further Notes

- The website's status sentence ("pull request #19 is ready for review, $9.80 observed of $11.40 reserved") is an announcement for assistive technology, not a panel field; its data is the same as the header's.
