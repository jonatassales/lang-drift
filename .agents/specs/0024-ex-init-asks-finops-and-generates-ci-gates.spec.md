# ex init asks for the FinOps inputs and can generate the CI Gates workflow

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

`ex init` records only the project identity and the Microruntimes execution rules. A team that has just initialized a project cannot plan or run anything yet. It still has to discover, one by one, the commands that record the team budget and the default gates. Nothing records where its tickets and specs live. And nothing offers CI Gates, the pull request report LangDrift wants to run in the team's own CI, so CI never shows what LangDrift knows, or does not know, about a change.

## Solution

When a person runs `ex init` at a terminal, `init` asks a short set of questions, then records each answer exactly as its own command would:

- the team budget for the current month, in dollars, with the per-run cap and the reserve margin;
- whether to apply the default gates;
- the folder that holds the team's artifacts, its tickets and specs;
- whether to generate the CI Gates workflow, defaulting to no.

Every answer also has a flag. Without a terminal, `init` uses the flags it was given and asks nothing, so a script that runs `ex init` behaves exactly as it does today. Running `init` again asks only for what is not yet recorded and never overwrites a recorded input. At the end, `init` lists the commands still to run: the hardware envelope, the model and the driver.

With a yes, `init` writes one file outside `.langdrift/`: `.github/workflows/langdrift-gates.yml`. It is a thin workflow pinned to the `ex` version that generated it, and it runs the new `ex gates report` on every pull request. The report writes a Markdown summary of the pull request's commit range to the job summary. It marks every FinOps reading as unknown, because `.langdrift/` is gitignored and CI has no FinOps source. It needs no project, posts nothing and fails only when git cannot be read. An existing workflow file is kept byte for byte.

## User Stories

1. As a team lead initializing a project, I want `ex init` to ask for the monthly budget, so that the team envelope is recorded without learning `ex finops budget set` first.
2. As a team lead, I want to type the budget in dollars, such as `1500` or `1500.50`, so that I never convert to cents by hand.
3. As a team lead, I want the prompt to name the month it records, so that I know which period the budget covers.
4. As a team lead, I want `init` to ask for the per-run cap, so that no single run can spend the whole envelope.
5. As a team lead, I want `init` to ask for the reserve margin in percent, with no default offered, so that the limit stays my decision.
6. As a team lead not ready to set a budget, I want an empty answer to the budget question to skip the whole budget, so that `init` still finishes and records the rest.
7. As a team lead, I want a malformed amount, such as `abc` or `-5`, to be refused with the reason before anything is recorded, so that a typo never becomes a budget.
8. As a team lead, I want `init` to offer the default gates, so that the first run has gates without a separate command.
9. As a team lead, I want to name the folder that holds our tickets and specs, so that LangDrift knows where our artifacts live.
10. As a team lead, I want the artifacts folder to be refused when it is absolute, climbs out of the project or does not exist, so that `config.json` never names a path outside the project.
11. As a team lead, I want an empty answer to skip the artifacts folder, so that a team without one can still initialize.
12. As a team lead, I want `init` to ask whether to generate CI Gates, defaulting to no, so that nothing lands in my CI without my yes.
13. As a team lead who answers yes, I want `.github/workflows/langdrift-gates.yml` written in the project, so that the next pull request runs the gates.
14. As a team lead whose project already has that workflow, I want `init` to keep it byte for byte and say so, so that my edits survive.
15. As a team lead, I want the generated file to open with a header naming the generator and its version and saying that the team owns the file, so that a reader knows where it came from and that it is theirs to change.
16. As a security reviewer, I want the generated workflow to grant only `contents: read`, so that CI Gates cannot write to the repository or the pull request.
17. As a security reviewer, I want `actions/checkout` and `actions/setup-node` pinned to commit SHAs with their version in a comment, so that a moved tag cannot change what runs.
18. As a security reviewer, I want checkout to keep no credentials, so that later steps cannot reuse the token.
19. As a team lead, I want the workflow to run `npx --yes @langdrift/ex@<version> gates report` at the exact version that generated it, so that a new `ex` release never changes CI without a change to the file.
20. As a team lead, I want the workflow to have a timeout and one concurrency group per pull request, so that a stuck or superseded run does not hold a runner.
21. As a team lead, I want the workflow to trigger on every pull request, whatever its base branch, so that every change gets its report.
22. As a reviewer of a pull request, I want the job summary to show the pull request's commit range, so that I see what the change contains.
23. As a reviewer, I want every FinOps reading in the report marked unknown, with the reason that CI has no committed FinOps source, so that a missing reading is never read as zero spend.
24. As a team lead, I want `ex gates report` to pass on a repository without `.langdrift/`, so that CI Gates works in a fresh clone.
25. As a team lead, I want `ex gates report` to fail with the reason when git cannot read the range, so that a broken checkout is visible rather than reported as empty.
26. As a team lead, I want `ex gates report` to print the report to standard output when no job summary file is given, so that I can preview it locally.
27. As a script author, I want `ex init` without a terminal to ask nothing, so that piped or CI use never blocks on a prompt.
28. As a script author, I want flags for every answer, so that I can initialize a project fully in one non-interactive command.
29. As a script author, I want `ex init` with no flags and no terminal to behave exactly as today, so that existing automation keeps working.
30. As a script author, I want `ex init --json` to ask nothing, so that its result stays one JSON document.
31. As a test author, I want a flag that makes `init` ask its questions on any input, so that the questions can be tested without a terminal.
32. As a team lead running `init` a second time, I want it to ask only for inputs not yet recorded, so that I am not asked twice.
33. As a team lead, I want a second `init` never to overwrite a recorded budget, gates or artifacts folder, so that changing one stays a deliberate act through its own command.
34. As an auditor, I want every input `init` records to carry who set it, through one `--by`, so that its record matches one written by its own command.
35. As an auditor, I want the budget and gates recorded through the same use cases as `ex finops budget set` and `ex finops gates set`, so that their mirrors and history are identical whichever command wrote them.
36. As a team lead, I want `init` to end by listing the commands still to run, the envelope, the model and the driver, so that I know what blocks the first run.
37. As a team lead, I want the result to say what `init` recorded, skipped or kept, so that I can check what happened.
38. As a maintainer, I want the README's boundary to say that `init` may write the one CI Gates workflow file after an explicit yes, so that the documented boundary matches the behavior.

## Implementation Decisions

- **The questions.** `init` asks, in this order:
  1. the monthly budget for the current month, in dollars;
  2. the per-run cap, in dollars;
  3. the reserve margin, in percent;
  4. whether to apply the default gates, `[y/N]`;
  5. the artifacts folder;
  6. whether to generate CI Gates, `[y/N]`.

  The questions go to the diagnostics stream through the existing prompter, as in `runtime init --guided`, so they never mix with the result.
- **Optional answers.** The current prompter treats an empty answer as a cancellation. `init` needs an empty answer to skip a question instead, so the prompter gains an optional text question that returns an empty answer as "skipped". An empty first budget answer skips questions 2 and 3.
- **When `init` asks.** It asks when its input is a terminal (`support.Interactive`) and `--json` is off. A boolean `--ask` overrides the terminal check either way, following `agent run --confirm`. With `--json`, `init` asks nothing.
- **Flags.** `--budget <dollars>`, `--per-run <dollars>`, `--reserve-margin <percent>`, `--default-gates`, `--artifacts <path>`, `--ci-gates` and `--by <who>`. `--by` uses the shared by-flag helper and its `DefaultBy`. Each flag answers its question, so that question is never asked. The three budget flags go together: one without the others is a usage error.
- **Dollars.** A dollar amount is a positive number with at most two decimals, converted to whole cents. The conversion and its validation form one pure function in the FinOps domain, beside the existing dollar formatter. The budget's other bounds keep the values `ex finops budget set` already applies. The currency stays the existing default, USD.
- **Recording.** `InitializeProject` keeps its job: creating the identity and reporting `initialized` or `already initialized`. A new application use case, run after it, records the answers:
  - the budget through `inputs.Service.SetBudget`, with the period set to the current month (`YYYY-MM`) from the injectable clock;
  - the default gates through the same service call `ex finops gates set --default` makes;
  - the artifacts folder in `config.json`.

  "Already recorded" means the service already holds a budget for the period, gates, or an artifacts folder in the config; that input is skipped and not asked.
- **`config.json`.** It gains an optional `artifacts` field: a slash-separated path relative to the project root. The file is omitted when empty, so a project without it keeps today's exact layout, and the version stays 1. The path must be relative, must not climb with `..` and must name an existing directory inside the project, checked through the existing `files.Within` boundary.
- **The CI Gates workflow.** A new scaffold adapter writes `.github/workflows/langdrift-gates.yml` from an embedded template. It fills in only the `ex` version, using the same exclusive-create, keep-if-present pattern as the execution rules. A version of `0.0.0-dev` (an unreleased build) is refused, because `npx` could not resolve it; `init` says so and writes nothing. The template:
  - opens with a header comment naming `ex init`, the version, and that the team owns the file and may edit it;
  - triggers on `pull_request` with no branch filter;
  - grants top-level `permissions: contents: read`;
  - sets `concurrency` per pull request, with `cancel-in-progress`;
  - has one job with `timeout-minutes`;
  - runs `actions/checkout` at a pinned SHA with a version comment, `persist-credentials: false` and `fetch-depth: 0`;
  - runs `actions/setup-node` at a pinned SHA with a version comment;
  - runs one step, `npx --yes @langdrift/ex@<version> gates report`. `BASE_SHA` and `HEAD_SHA` come from the pull request event through `env`, never inside `run`.
- **`ex gates report`.** A new top-level `gates` command group with a `report` subcommand:
  - It needs no project and never opens `.langdrift/`.
  - It reads `BASE_SHA` and `HEAD_SHA` at the operation. A missing one is a usage error that names the variable.
  - It lists the commits of the range through `git` as its own process.
  - It renders Markdown, written to the file `GITHUB_STEP_SUMMARY` names, or to standard output when that variable is unset. The Markdown has a title, the range with its commit count, the commits as short SHA and subject, and a FinOps section marking the budget, the spend and the ledger as unknown, with the reason that no committed FinOps source exists.
  - It exits 0 on success. It exits with the existing external-failure code when git fails or the summary cannot be written.
- **The result.** `ProjectResult` keeps `project`, `id` and `status`. It gains a `recorded` list (`budget`, `gates`, `artifacts`, `ci-gates`), a `kept` list for an existing workflow file, and a `next` list of the commands still to run (`ex harness envelope set`, `ex harness model set`, `ex harness driver set`). Lists left empty are omitted, so an `init` without answers prints today's JSON.
- **The boundary.** The README's boundary says the CLI writes only `.langdrift/`, except for the one CI Gates workflow file `init` writes after an explicit yes. The scaffold comment "LangDrift prescribes no workflow" stays true of the `.langdrift/` scaffold. The existing init test's "no workflow" promise holds for an `init` without the yes.

## Testing Decisions

- A good test drives `ex` as a command, with arguments, piped answers and a temporary project. It checks the exit status, the result, the files written and the recorded inputs, never the internal types.
- The seam is the existing command-level harness: `newProject`, `p.ex` and `p.exInput` in the module's root tests, as in `langdrift_init_test.go`, `inputs_test.go` and the guided answers of `runtime init --guided`. `--ask` makes the questions testable on piped input. No new seam is needed.
- Tests for `init`:
  - all six answers, checking the budget mirror (cents, period, `by`), the gates mirror, `config.json`'s `artifacts`, and the workflow file against a golden copy kept under testdata;
  - an empty first budget answer skips the budget;
  - a malformed amount and a climbing artifacts path are refused before anything is recorded;
  - flags without a terminal record without asking;
  - with no flags and no terminal, the existing identity-only test keeps passing unchanged;
  - a second `init` asks only for what is missing;
  - an existing workflow file is kept byte for byte;
  - an unreleased version refuses the workflow.
- Tests for `ex gates report` run in a temporary git repository with two commits:
  - the report lists the range and marks FinOps unknown;
  - `GITHUB_STEP_SUMMARY` receives it;
  - a missing `BASE_SHA` is a usage error;
  - an unknown SHA fails with the external-failure status.
- The dollars-to-cents conversion gets a table test in the FinOps domain, as for the existing formatter.
- The golden workflow is also checked for the repository's own workflow conventions where they apply to a consumer file: pinned SHAs with version comments, top-level permissions, no expression inside `run`.

## Out of Scope

- Asking for the hardware envelope, the model or the driver: `init` only lists their commands.
- `ex agent plan` resolving a ticket by number through the artifacts folder.
- A committed FinOps source, or the console API, for CI; and any FinOps reading other than unknown in the report.
- Sticky comments, check runs, widgets and a LangDrift avatar.
- The Marketplace action and billing through a LangDrift API key.
- `ex gates init --upgrade`, and any rewrite of an existing workflow file.
- Changing a recorded input through `init`.

## Further Notes

- Owner decisions, from this spec's grilling (2026-10-29): Q1 budget, default gates, artifacts and CI Gates; Q2 terminal only, a flag for every answer; Q3 a second run asks only for what is missing; Q4 one `--by`, the existing use cases; Q5 dollars, the current month, no reserve-margin default, an empty answer skips; Q6 artifacts means the tickets and specs folder, stored only; Q7 one file outside `.langdrift/` after an explicit yes; Q8 a report in the job summary only; Q9 FinOps unknown in CI; Q10 `npx` pinned to the generating version, actions pinned to SHAs; Q11 `pull_request`, `contents: read`, timeout, concurrency, checkout without credentials and with full history; Q12 a header and team ownership; Q13 `gates report` needs no project.
- CI Gates ships as a workflow generated into the team's repository, not as a GitHub App, so it runs with the team's own token.
