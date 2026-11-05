# The stage driver runs each stage on the chosen model and effort

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E11.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 13 of 19, after spec 0057 and before spec 0059. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The hero simulation's run prints "intent claude-sonnet-5-5 · effort medium", prices its estimate at Sonnet's list price, and then runs grilling, to-spec, to-tickets, implement and code-review. Everything the card and the terminal say about money assumes that those stages ran on the model and at the effort the team chose. The stage driver owns no line of its own, yet every skill line depends on it.

In the LangDrift CLI the chosen model and effort never reach the coding agent. A stage call carries the skill, the ticket, the run, the stage and the lap, but no model and no effort, and the driver's placeholders offer none. The observed spend is then priced at a model the agent may not have used. The driver also adds what the simulation does not show:

- " · no usage reported" appended to a stage line;
- an "exit N · no summary" fallback;
- a summary cut from the agent's last output line;
- a gate before every call, with an "allow unmeasured" setting;
- a timeout, and one paid call per skill rather than per stage;
- a configuration held as an input record set through a command.

## Solution

The stage driver is the team's coding agent, as `dft agent run` starts it once for each stage that has a skill: Planning, Spec, Tickets, Implement and Review. The applied model and effort reach every call, both as placeholders in the agent's arguments and as environment variables. A configuration that cannot pass them is refused before the run prints anything. The driver reads each call's input and output tokens for the usage record and gives the skill the path where it reports its result. It prints nothing itself; each stage's one line belongs to the skills (spec 0059). Its configuration is a plain file in `.langdrift/harness/`.

## Responsibility

This element starts the coding agent for a stage. It passes the run, the ticket, the stage, the skill text, the applied model and effort, and the path of the stage's report. It keeps the call's output in a private log, reads the call's input and output tokens, and returns the exit status, the tokens and the elapsed time.

It does not:

- choose the model, the effort or which stages run; spec 0057 does;
- print a stage line or read a stage's figures; spec 0059 does;
- price the tokens; spec 0056 does;
- keep the usage record or the evidence; spec 0065 does;
- run the effort hook; spec 0060 does;
- make any call during the dev ticket; spec 0061 makes none.

## Produces

Under the client's `.langdrift/`:

- **`harness/driver.json`, the driver configuration.** It replaces the driver input record and holds:
  - `executable`, the coding agent's command;
  - `args`, its arguments, where the placeholders `{{prompt}}`, `{{model}}`, `{{effort}}`, `{{report}}`, `{{skill}}`, `{{stage}}`, `{{ticket}}` and `{{run}}` are filled per argument and never split an argument;
  - `passEnv`, the variables passed through from `dft`'s environment;
  - `usage.inputTokens` and `usage.outputTokens`, the dotted paths of the tokens in the agent's JSON output.

  `args` must carry `{{model}}` and `{{effort}}`.
- **The call's environment.** It is the minimal base plus `LANGDRIFT_PROJECT_ROOT`, `LANGDRIFT_RUN`, `LANGDRIFT_TICKET`, `LANGDRIFT_STAGE`, `LANGDRIFT_SKILL`, `LANGDRIFT_MODEL`, `LANGDRIFT_EFFORT` and `LANGDRIFT_REPORT`, and the `passEnv` variables. `LANGDRIFT_MODEL` is the price-list id, such as `claude-sonnet-5-5`; `LANGDRIFT_EFFORT` is `low`, `medium` or `high`. `LANGDRIFT_RUN` and `{{run}}` carry the run's id (spec 0065).
- **The prompt.** `{{prompt}}` expands to the stage's skill text as the card names it (`grilling`, `to-spec`, `to-tickets`, `implement · tdd`, `code-review`), the ticket, the effort, and the instruction to write the stage's report at the report path (spec 0059 defines the report).
- **The call's log, `agent/state/runs/<run>/stages/<stage>.log`**, where `<run>` is the run's id (spec 0065). It is private and ignored, and is never evidence.
- **The call's result, handed to spec 0065.** It holds the stage, the exit status, the input and output tokens and the elapsed milliseconds. The tokens feed the usage record, from which spec 0056 computes the observed spend: with the fixture at medium effort, the five calls total 3,700K input and 240K output tokens.

## Shows

No terminal line and no card field. The driver runs the five calls behind spec 0059's lines, in this order, each after its stage's checkpoint when the stage is ticked (spec 0057):

| Stage | Call's skill text | Line it enables (spec 0059) |
| --- | --- | --- |
| `planning` | `grilling` | `console › grilling   3 open questions settled` |
| `spec` | `to-spec` | `console › to-spec    0011-SPEC revised · 3 acceptance criteria` |
| `tickets` | `to-tickets` | `console › to-tickets 3 tickets · blocking edges recorded` |
| `implement` | `implement · tdd` | `console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part` |
| `review` | `code-review` | `console › code-review 2 findings · fixed` |

The `implement` call follows spec 0060's `hook` line. The `review` call follows spec 0061's `✓ closed     dev-ticket …` line.

## User Stories

1. As the owner, I want the coding agent to run on the model I chose, so that the work matches the intent line and the price I was shown.
2. As the owner, I want the coding agent to run at the effort I chose, so that "effort medium" in the intent line is true of the work.
3. As a team member, I want the model id and the effort passed both as arguments and as environment variables, so that any coding agent can be configured to take them.
4. As a team member, I want a driver configuration that does not pass the model and the effort refused before the run prints anything, so that no run can silently use another model.
5. As a team member, I want to configure our coding agent in one file under `.langdrift/harness/`, so that I can read and version how stages are run.
6. As the owner, I want one call per stage, so that `implement · tdd` is one piece of work and one line.
7. As the owner, I want the review call to run after the dev ticket, so that the agent reviews the part the person wrote.
8. As the owner, I want no call made while the person writes the dev ticket, so that "no paid call" is true.
9. As a FinOps reader, I want each call's input and output tokens read from the agent's own output, so that the observed spend is measured, not guessed.
10. As a FinOps reader, I want the observed tokens priced at the list price of the model the agent actually ran, so that estimate and observation compare like for like.
11. As the owner, I want the driver to print nothing of its own, so that the session shows only the simulation's lines.
12. As a skill author, I want the path where my stage's report goes passed to the call, so that my skill can report its figures to `dft`.
13. As a security-minded team member, I want the agent started directly, with arguments that are never re-split and a minimal environment plus what we pass through, so that no shell interprets a ticket's text.
14. As an auditor, I want each call's output kept in a private log and out of evidence, so that tokens of text never leak into the record.
15. As the owner, I want a call that fails to stop the run without inventing a line, so that the session never claims work that did not happen.

## Implementation Decisions

**Remove**, each beyond what the simulation shows:

- " · no usage reported" appended to a stage line, and the "exit N · no summary" fallback. Serves each stage line, which is exactly spec 0059's form.
- The summary read from the agent's last output line, or from a configured path, and cut to 120 characters. Spec 0059's report replaces it.
- The gate before every call (FG-04) and the driver's `unmeasured` setting (spec 0056 removes the gate).
- `usage.costCents`. The observed spend is the list price applied to the tokens (spec 0056).
- `timeoutSeconds`, with its "did not finish within N seconds" failure. No time limit is shown (see Open Questions).
- One paid call per skill, which split `implement` from `tdd` and ran `code-review` inside Implement. Serves one line per stage.
- The `{{lap}}` placeholder and `LANGDRIFT_LAP`, with the lap flow spec 0057 removes. Also `LANGDRIFT_EXECUTION_ID`: no microruntime is up while the stages run (spec 0062 brings it up after `code-review`).
- The driver as an input record and the command that sets it (spec 0048 removes the commands).
- The driver's `Hook` operation, which ran committed `sh` scripts (spec 0060).
- The `summary` field of a stage's evidence. The evidence keeps the report's figures (spec 0065).

**Build**, each shown and missing:

- The applied model and effort on every stage call, as `{{model}}`, `{{effort}}`, `LANGDRIFT_MODEL` and `LANGDRIFT_EFFORT`. Serves the `intent` line's claim and every stage line after it.
- The check that `harness/driver.json` exists and that its `args` carry `{{model}}` and `{{effort}}`, made when `dft agent run` starts, before its first line (`# run · …`) prints. A missing or refused configuration is a stop with nothing printed and nothing reserved (spec 0048).
- The report path, as `{{report}}` and `LANGDRIFT_REPORT`. Serves every stage line through spec 0059.
- One call per stage that has a skill, in the order spec 0057 plays them. Serves the five stage lines.

**Change**, each different:

- The configuration is the file `harness/driver.json`, read when the run starts.
- The prompt names the stage's skill text, the ticket, the effort and the report path, instead of asking the agent to end with one summary line.
- A call's result carries tokens as required: a call whose output does not yield both token counts at the configured paths is a failed call.

**Interfaces with neighbours:**

- From spec 0057: the stage id, its skill text, the applied model and effort, and when to call.
- To spec 0059: the report path given to the call, and the call's exit status. Spec 0059 reads the report and prints the line.
- To spec 0065: the call's result. Through spec 0065, to spec 0056: the run's summed tokens.
- From spec 0048: the shared stop rule for a failed call.

**Failure paths the simulation does not show:**

- A call that does not start, exits with a non-zero status, is interrupted, or reports no usage stops the run after the last printed line. Its stage prints no line.
- The log is kept, and the call's result is handed to spec 0065 with tokens unknown, never zero, when they were not reported.
- Nothing is retried and nothing resumes (spec 0057). The stop follows spec 0048's single rule.

## Testing Decisions

- **The seam is spec 0048's end-to-end golden test.** The shared fixture's stage driver is a small fake coding agent: a POSIX shell script committed in the fixture. It records the arguments and the `LANGDRIFT_*` variables of each call, writes the stage's report (spec 0059), and prints JSON usage that sums to the simulation's observed tokens for the applied effort. This spec adds:
  - exactly five calls, in the order Planning, Spec, Tickets, Implement and Review; none between `waiting    dev-ticket` and `closed     dev-ticket`, and none in an over-budget run;
  - every call carrying `claude-sonnet-5-5` and `medium` in its arguments and its environment for the default run, and the applied pair in a variant that records Haiku at low effort;
  - the summed usage feeding the money record that spec 0056 asserts.
- **Unit tests of the driver adapter**, with a fake process runner:
  - substitution keeps an argument with spaces whole;
  - a configuration without `{{model}}` or `{{effort}}` is refused before any output;
  - output without the token paths is a failed call with unknown usage.

  Prior art: the stage driver's existing tests.

## Out of Scope

- Which stages run and when (spec 0057), and the effort hook (spec 0060).
- The skills' reports and their lines (spec 0059), and the dev ticket (spec 0061).
- Pricing (spec 0056) and the usage record and evidence (spec 0065).
- The judge's calls during planning (spec 0051).
- The console and its stub API.

## Open Questions

- **When a run without a usable driver configuration stops.** The simulation never shows a run without one. This spec checks the configuration before the run prints anything, so no gate opens and nothing is reserved for a run that cannot start a stage; an over-budget run without a configuration therefore stops instead of pausing. The alternative is to check it only before the first stage call.
- **A ready configuration for one coding agent.** The simulation runs Claude models. Whether `dft init` writes a working `harness/driver.json` for a particular coding agent, and which agent, arguments and usage paths, is the owner's decision. Until then the file is the team's to fill.
- **A call that fails.** The simulation never shows a stage call that fails, hangs or reports no usage. What the run prints and records then, and whether a time limit is wanted, is open. Until the owner decides, such a call stops the run silently as described above.

## Further Notes

- The driver remains a trusted process started directly, never through a shell command string. Its output never enters evidence.
- How a coding agent maps `{{effort}}` onto its own reasoning setting is the agent's configuration. `dft` guarantees only that the applied value reaches it.
