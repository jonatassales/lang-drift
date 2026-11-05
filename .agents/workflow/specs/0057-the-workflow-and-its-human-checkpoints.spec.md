# The workflow runs its stages and waits for a person where the team asked

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E10.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 12 of 19, after spec 0055 and before spec 0058. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The hero simulation's Workflow card is the run's input. The team ticks the stages where a person takes part (Planning always, Review by default), picks a model and an effort, and plays the run. `dft agent run` then waits for the owner before each ticked stage's work, inside the same command: "waiting planning · the owner confirms", then "confirmed planning · by the owner". It runs every other stage without waiting and ends "for review" when Review is ticked, "to merge" when it is not. The card follows each stage from pending to done and ends with Review "waiting for the owner" and Merge "always the owner's".

The LangDrift CLI does something else:

- It takes the waiting stages from a workflow mirror's defaults, with no per-run choice, and the model and effort from the plan's suggestion or from `--model` and `--effort` flags.
- Its sample workflow has no Tickets stage and no Dev ticket stage, and it runs `code-review` inside Implement.
- It asks "Confirm planning? [y/N]" on stderr. Without a terminal it exits while waiting, with the runtime up, and resumes after `dft agent confirm` with a second `run` that prints "lap … resumed".
- It prints "by owner" instead of "by the owner".
- It ends "for review" whenever delivery is not clean, not because of the Review checkbox.
- It writes a proposed ticket status that nothing shows.

## Solution

The workflow is the simulation's seven stages, in its order, with the skills the card names. One choice per run carries the model, the effort and the stages a person confirms, and the run applies it. Its defaults are the workflow's: Claude Sonnet 5.5 at medium effort, with Planning plus Review, the choice the Workflow card shows from the moment the panel appears, before the judge has suggested anything. Every ticked stage among Planning, Spec, Tickets and Implement prints its `waiting` line before its work and its `confirmed` line once the owner answers. The answer is given inside the same `dft agent run`: no other command, no second run and no prompt line. The Review checkbox alone decides whether the run ends for review or to merge, and the run's last comment says what the owner does next. The stage states the card shows are data the run records as it goes.

## Responsibility

This element defines the workflow's stages, the per-run choice of model, effort and human stages with its defaults, and the human checkpoints of `dft agent run`. It also records each stage's state for the card, decides from the Review checkbox how the run ends, and prints the run's closing `# next` comment.

It does not:

- run a stage's skill; spec 0058 (stage driver) and spec 0059 (agent skills) do;
- run the effort hook before Implement; spec 0060 does;
- define the dev ticket's own wait and close or Safe Seat; spec 0061 does, and this element only places the stage and shows it in HUMAN IN THE LOOP;
- price the choice or gate the run; spec 0056 does;
- print the `ready` line; spec 0063 does, from this element's decision;
- carry the choice between plan, card and run; spec 0050's hand-off and recorded intent do.

## Produces

Under the client's `.langdrift/`:

- **`harness/workflow.json`, the workflow.** It is LangDrift's opinionated spec-driven development workflow, set up by `dft init` (spec 0048), and the same workflow this repository runs on its own work, as `.agents/context.guide.md` and `.agents/policies/human-in-the-loop.policy.md` describe it, outside any `.langdrift/` folder; its stages run the skills `dft init` installs under `.langdrift/skills/` (spec 0059), and the tickets that take part in it are the client's, in `.agents/` or `.claude/` (spec 0053). It keeps its place and drops its mirror header, `status`, `owner`, `note`, `efforts`, the `human` levels and `defaults.humanAt`. It holds the seven stages in order, each with `id`, `skill` (the skill text the card shows, or none) and `person`:

  | `id` | `skill` | `person` |
  | --- | --- | --- |
  | `planning` | `grilling` | `required`: always confirmed |
  | `spec` | `to-spec` | `optional` |
  | `tickets` | `to-tickets` | `optional` |
  | `implement` | `implement · tdd` | `optional` |
  | `dev-ticket` | none | `writes`: the person's own work, never a checkbox |
  | `review` | `code-review` | `optional` |
  | `merge` | none | `owner`: always the owner's, never played |

  It also holds the default choice, `defaults`: `model` `claude-sonnet-5-5`, `effort` `medium` and `human` `planning` and `review`, the simulation's initial Workflow choices. It feeds HUMAN IN THE LOOP's rows, their skill texts, which rows carry a checkbox, their initial ticks, and the MODEL and EFFORT the card shows before any suggestion.
- **The choice, in the recorded intent that spec 0050 carries.** It holds `model` (a price-list id, spec 0056), `effort` (`low`, `medium` or `high`) and `human`, a flag for each of `planning`, `spec`, `tickets`, `implement` and `review`. `dft agent plan` records the workflow's default choice, `defaults`: `claude-sonnet-5-5` at `medium`, with Planning and Review. It writes it into the plan's card data before the plan's first line (spec 0050), so the Workflow card shows it from the moment the panel appears, and into the proposed intent at `✓ plan`. The judge's suggestion (spec 0051) is a recommendation on the Analysis card and does not change the default choice. The team may change the choice in the intent before it authorizes it (spec 0050), and the run applies the authorized choice. `planning` is always true. It feeds the MODEL and EFFORT choices, the ticks, the `intent` line's model and effort (spec 0050), spec 0056's estimate and spec 0058's call.
- **The stage states, in the run's record (spec 0065).** It holds the choice as applied, and per stage its `id`, `human` (ticked or not), `state` (`pending`, `waiting`, `active` or `done`) and `confirmedBy` (`owner` once confirmed). It also holds `reviewed`, the Review tick, which feeds spec 0063's choice between `for review` and `to merge` and this element's `# next` line. The states feed the card's stage states.

## Shows

Phase 2, `dft agent run`. For each ticked stage among `planning`, `spec`, `tickets` and `implement`, this pair prints before the stage's own work:

```
console › waiting    {stage} · the owner confirms
✓ confirmed  {stage} · by the owner
```

With the default choice only Planning waits, right after `gate … opened` (spec 0056) and before `grilling` (spec 0059):

```
console › waiting    planning · the owner confirms
✓ confirmed  planning · by the owner
```

Where each pair goes when its stage is ticked:

- **Spec:** after `grilling` and before `to-spec`:

  ```
  console › waiting    spec · the owner confirms
  ✓ confirmed  spec · by the owner
  ```

- **Tickets:** after `to-spec` and before `to-tickets`:

  ```
  console › waiting    tickets · the owner confirms
  ✓ confirmed  tickets · by the owner
  ```

- **Implement:** after `to-tickets` and before the `hook` line (spec 0060):

  ```
  console › waiting    implement · the owner confirms
  ✓ confirmed  implement · by the owner
  ```

An unticked stage prints no pair. Review never prints a pair: its `code-review` runs without waiting even when Review is ticked, and the tick only decides how the run ends.

The run's last line follows the `ready` line (spec 0063). With Review ticked:

```
# next: the owner reviews #19 and decides whether to merge
```

With Review unticked:

```
# next: the merge stays the owner's
```

The pull request's number comes from spec 0063.

Phase-1 card fields, from the plan's hand-off:

- **HUMAN IN THE LOOP**, one row per stage:
  - Planning `grilling`, ticked; the box is unavailable, with the note `Planning always has a person in this workflow`;
  - Spec `to-spec`, Tickets `to-tickets` and Implement `implement · tdd`, unticked;
  - Dev ticket `driven by the person`, with no box;
  - Review `code-review`, ticked;
  - Merge `always the owner's`, with no box and a mark.

  Each box's accessible name is `{stage} · human in the loop`.
- **MODEL:** `Claude Opus 5.5 | Claude Sonnet 5.5 | Claude Haiku 5.5`, with Claude Sonnet 5.5 chosen. Its price line is spec 0056's.
- **EFFORT:** `low | medium | high`, with medium chosen.

The choices are locked while a run plays.

Phase-2 card fields, from the run's stage states:

- A stage is pending until its first line prints, active while its lines print (in the primary tone) and done, with a mark, after its last.
- A ticked stage reads `waiting for the owner` while it is active.
- Review's lines run from `code-review` to `ready`. Once `ready` has printed, Review reads `waiting for the owner` when ticked; when unticked it is done.
- The Dev ticket row reads `the person writes the reserved part` while active (its states are spec 0061's).
- Merge always reads `always the owner's`.

After the default run, HUMAN IN THE LOOP reads: Planning ✓ · Spec ✓ · Tickets ✓ · Implement ✓ · Dev ticket ✓ · Review `waiting for the owner` · Merge `always the owner's`. A run paused at the gate leaves every stage pending.

## User Stories

1. As a team member, I want the workflow to list Planning, Spec, Tickets, Implement, Dev ticket, Review and Merge with the skills each runs, so that I know what the run will do before I play it.
2. As a team member, I want Planning always to have a person, its box ticked and unavailable with a note saying why, so that no run starts its work without the owner.
3. As a team member, I want Planning and Review ticked by default, so that the common case needs no change.
4. As a team member, I want to tick Spec, Tickets or Implement for one run, so that I can confirm the stages I care about this time.
5. As a team member, I want the Dev ticket shown as driven by the person, with no box, so that nobody can switch off the part a person writes.
6. As a team member, I want Merge shown as always the owner's, so that I know the run never merges.
7. As a team member, I want the model and the effort chosen per run, starting from the workflow's default and shown beside the judge's suggestion, so that I can follow or override it before any spend.
8. As the owner running `dft agent run`, I want the run to stop before each ticked stage's work and say "waiting {stage} · the owner confirms", so that I know it is my turn.
9. As the owner, I want to confirm in the same terminal and see "confirmed {stage} · by the owner", so that I need no second command and no second run.
10. As the owner, I want no prompt line added to the session, so that the terminal reads exactly as the simulation shows.
11. As the owner, I want unticked stages to run without waiting, so that the run only stops where the team asked.
12. As the owner, I want `code-review` to run without waiting even when Review is ticked, so that my review starts on a pull request the agent has already reviewed.
13. As the owner, I want a ticked Review to end the run "for review" and tell me "the owner reviews #19 and decides whether to merge", so that I know the pull request waits for me.
14. As the owner, I want an unticked Review to end the run "to merge" and say "the merge stays the owner's", so that I know only the merge is left to me.
15. As the owner, I want the Review tick alone to decide that ending, so that the outcome never depends on anything I did not choose.
16. As a team member watching the card, I want each stage to turn active when its first line prints and done after its last, so that the card follows the terminal.
17. As a team member, I want a ticked stage to read "waiting for the owner" while it is active, and Review to read it after the run when ticked, so that I see where a person is involved.
18. As a team member, I want the choices locked while the run plays, so that the card never describes a choice the run did not apply.
19. As an auditor, I want the run's record to keep the choice it applied and who confirmed each stage, so that I can tell which stages a person approved.
20. As the console, I want the stage states recorded by the run, so that I only show them.
21. As the owner, I want a run whose input ends while it waits to stop there with nothing more printed, so that an unattended run never proceeds without me.

## Implementation Decisions

**Remove**, each beyond what the simulation shows:

- The workflow mirror's `human: required | optional | owner-only` levels, `defaults.humanAt` including merge, and the run record's link to the workflow record and mirror. Serves HUMAN IN THE LOOP and its default ticks.
- `code-review` inside Implement, and `implement` and `tdd` as two skills. Serves the `implement · tdd` and `code-review` rows (spec 0059 makes them one stage line each).
- The `--model`, `--effort`, `--by` and `--confirm` flags, the choice's "overrode suggestion" record and its `By`, `Via` and `Record` fields. Serves the `intent` line's model and effort, taken from the applied choice.
- The `Confirm planning? [y/N]` prompt on stderr. Serves the waiting and confirmed pair, which has no line between its two lines.
- The exit in the waiting state when there is no terminal, with the runtime up, and with it `dft agent confirm`, the second `run` that resumes an open lap, the "lap … open · … resumed" line, the refusal of a model or effort while a lap is open, and the replay of a run from its first stage not done. Serves the rule that a checkpoint is answered inside `dft agent run`.
- "For review" whenever delivery is not clean (the delivered-cleanly rule), and the search for a person-only stage without skills. Serves `✓ ready … for review` or `… to merge` (spec 0063) and the `# next` lines.
- The proposed ticket status, `in-review` or `ready-to-merge`, written to the run record and never printed.
- The stage state `failed` and the per-stage `skillsDone` and `summary` used for resuming.

**Build**, each shown and missing:

- The seven-stage workflow above, with Tickets and Dev ticket and `code-review` under Review. Serves HUMAN IN THE LOOP and the run's line order.
- The per-run choice of model, effort and human stages, recorded with its defaults by the plan and applied by the run, with Planning forced on. Serves the MODEL, EFFORT and checkbox fields and every `waiting`/`confirmed` pair.
- A `waiting`/`confirmed` pair for each ticked stage among Spec, Tickets and Implement, placed before that stage's work and, for Implement, before the effort hook. Serves the pairs under Shows.
- The confirmation inside `dft agent run`. After its `waiting` line the run reads one line from standard input; in a terminal it reads with echo off, so no character reaches the session. The line confirms: the run records `confirmedBy: owner`, prints `confirmed` and plays the stage. End of input stops the run right after the `waiting` line (see failure paths). Serves `✓ confirmed  {stage} · by the owner`.
- The stage states, recorded at each line that changes them, and `reviewed`. Serves the card's stage states and the closing state of Review.

**Change**, each different:

- `confirmed` reads "by the owner", not the `--by` default "by owner".
- The Review tick alone decides the ending: `reviewed` true gives `for review` and `# next: the owner reviews #{number} and decides whether to merge`; false gives `to merge` and `# next: the merge stays the owner's`.
- Which stages wait comes from the applied choice, not the workflow's defaults; the defaults only seed the plan's default choice.
- Every `dft agent run` plays the stages from Planning, in order, after its gate; nothing resumes.

**Interfaces with neighbours:**

- To spec 0050: the default choice and the workflow, carried in the plan's hand-off. From spec 0050: the applied choice, as the recorded intent.
- To spec 0056: the applied model and effort. From spec 0056: the gate's outcome; nothing here runs after a pause.
- To spec 0058: for each stage with a skill, the stage id, its skill text and the applied model and effort, called in order once the stage's checkpoint, if any, is confirmed.
- To spec 0060: Implement's turn, after Implement's confirmation when it is ticked.
- To spec 0061: the Dev ticket's turn, after Implement and before Review; spec 0061 records its states.
- To spec 0063: `reviewed`, and the pull request's number for the `# next` line.

**Failure paths the simulation does not show:**

- When standard input ends while a stage waits, the run stops after that `waiting` line. It prints nothing more and leaves the stage `waiting` in the record; spec 0048's rule for every stop applies.
- A stage that fails stops the run after its last printed line, with its state left as it was.
- A later `dft agent run` is a new run from the gate; nothing resumes.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test** of `dft agent plan 0031-TICKET-book-a-demo-runtime` and `dft agent run` against the shared fixture. The test feeds one line on standard input for each expected confirmation. This spec adds:
  - the default run: only the Planning pair, between `gate … opened` and `grilling`; `# next: the owner reviews #19 and decides whether to merge` last; stage states ending Planning to Dev ticket done and Review waiting;
  - a variant with Spec, Tickets and Implement ticked: each pair at its place under Shows, Implement's before `hook`;
  - a variant with Review unticked: the run ends with `# next: the merge stays the owner's`, after spec 0063's `to merge` line, and Review is done;
  - a variant whose input is empty: the run prints through `console › waiting    planning · the owner confirms` and nothing after, calls no stage driver and records Planning as `waiting`;
  - the plan's hand-off: the seven stages with their skill texts and `person` kinds, and the default choice `claude-sonnet-5-5`, `medium`, Planning and Review.
- A unit test covers the choice rule: Planning stays ticked whatever the recorded choice says, and the ending follows `reviewed` alone.
- Prior art: the session catalog tests and the run tests that drive a fake stage driver.

## Out of Scope

- The skills and their lines (spec 0059) and the driver that runs them (spec 0058).
- The effort hook (spec 0060).
- The dev ticket's wait, its close and Safe Seat (spec 0061).
- The money and the gate (spec 0056), the hand-off (spec 0050), the `ready` line (spec 0063).
- The console and its stub API, which only receive the workflow, the choice and the states.

## Open Questions

- **The default choice when the suggestion differs.** The simulation's Workflow card starts from fixed choices, Claude Sonnet 5.5 at medium, which happen to equal the judge's suggestion. This spec takes the default from `harness/workflow.json`, so the card and the intent show it before any suggestion. Should the default instead follow the judge's suggestion when the two differ? How the team changes the choice and authorizes is spec 0050's open question.
- **Declining a checkpoint.** The simulation only shows a confirmation. What the owner types to decline, and what the run prints and records then, is open. Until then any line confirms, and only end of input or an interrupt stops the run.
- **Who "the owner" is.** `dft` cannot tell who answers at the terminal. It records the role `owner` for whoever confirms. Whether the confirmation must be tied to an identity is open.
- **Changing the workflow.** Whether a team may change the stages, their skills or their order is open, since the simulation shows one fixed workflow. Until the owner decides, `dft` plays the seven stages above.

## Further Notes

- `waiting` prints stage ids (`planning`, `spec`, `tickets`, `implement`), while the card shows labels. The labels are copy in the console's catalogs, not data.
- The paired lines' columns follow the CLI's layout, spec 0048.
