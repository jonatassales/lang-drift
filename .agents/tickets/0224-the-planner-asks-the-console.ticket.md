# 0224: the planner asks the console, and no plan exists without a judge

**Spec:** .agents/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** A team member plans a ticket with no provider configured and no key exported: `dft agent plan` sends its three questions to the console the project is connected to and writes the plan from the judged answers. `dft agent connect` records the console URL and the `RUNTIME_SECRET` it reads from the environment in the private state and verifies them with the judge route's connection check, no longer binding through the console's MCP route. The remote judge adapter posts each question with the secret as a bearer, within a bounded timeout, never retrying a paid call, and maps the route's refusals to the planner's errors. `--no-judge` and the plan "from the inputs" are gone: a project without a connection, an unreachable console, a refused secret or a model that refuses all refuse the plan naming the cause, and no proposed plan is written. The plan records the judge's prompt version and hash, and the ledger's planning line records the cost the console returned with its provenance. The provider SDK leaves the CLI. The session prints its reasoning lines only when the judge answered.

**Blocked by:** 0221, 0222.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] In a connected project, `dft agent plan` with the composed test's fake console answers writes a proposed plan whose fence records `judge.version`, `judge.promptHash` and a suggestion `by: judge`; the ledger's planning line, which gains these fields, carries the console's cost, its provenance `console` and whether it is illustrative.
- [ ] Before the first question the planner reserves the planning forecast at the planner model's list price against the ticket's allocation and the per-run cap; when it does not fit, the gate pauses before any question and nothing is paid, as today. Every answer is parsed against the contract before use; an answer outside it is a refusal.
- [ ] `dft agent connect` with `RUNTIME_SECRET` exported records the connection after a 204 from the console and refuses it after a 401, naming the console URL.
- [ ] Without a connection, the plan is refused with "connect the project"; with a console answering 401, 503 or a timeout, with the route's reason; with a model refusal, with the refusal; in every case no plan is written as proposed. A refusal before the model runs pays nothing; the model's own refusal is charged as the console prices it (spec 0045: the console charges the list cost of the tokens it spent), and a call cut off holds its reservation, since unknown is never zero.
- [ ] `--no-judge` is unknown to the CLI, `suggestion.by` is never `input`, and the CLI's module no longer depends on the provider SDK.
- [ ] The remote adapter's tests run against an in-process server: the request body matches the contract's request schema, the bearer is the recorded secret, each refusal maps to the planner's error, and a paid call is sent once.
- [ ] The session transcript prints `agent › reasoning` lines only after a judged answer; the composed test's expected transcript shows it.
- [ ] The CLI README documents connecting before planning and the planning cost the gate reserves.

## Dev part

**Category:** `core-domain`. When a plan may not exist is the rule that keeps a judgment from being faked; the owner decides it line by line.

**Paths:**
- `lib/cli/src/application/plan/judge-refusal.go`
- `lib/cli/src/application/plan/judge-refusal_test.go`

**Boundary:** `plan.JudgeRefusal(connection Connection, outcome AskOutcome) error`, called by the planner before each question with the money gates' verdict (`AskCleared`, or `AskPaused` with the gate's detail, a missing planning forecast among them) and after it with the console's outcome (`AskAnswered`, `AskUnauthorized`, `AskUnavailable`, `AskTimedOut`, `AskRefused` for the model's refusal, `AskMalformed` for an answer or a question outside the contract). `Connection` (recorded or absent, with the console URL and the secret), `AskOutcome` (the question's kind, its `AskStatus` and the reason in the gate's, the console's or the model's words) and `AskStatus` are declared by the agent in `connection.go`, `ask-outcome.go` and `ask-status.go` beside it, one type per file as the rules require; `ErrNotConnected` and `ErrJudgeRefused` join `ErrModelUnavailable` and `ErrPaused` in `judge.go`. A `nil` return lets the plan continue, and the planner's existing errors for an unavailable model and a gate pause are among the ones it returns; `ErrNotConnected` exits 4, `ErrPaused` exits 3.

**Goal:** Decide, for every state of the connection and every outcome of a question, whether the plan continues or is refused, and with which message naming the cause, so that no path reaches a proposed plan without three judged answers.

- [ ] The table of connection states and outcomes is covered one case per test, and the only continuing cases are a recorded connection with a question the gates cleared, before it, and a recorded connection with an answered question, after it.
- [ ] Every refusal message names what the person does next: connect, retry later, or read the model's refusal; a gate pause keeps the planner's existing message.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `JudgeRefusal` and its test.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go test -race ./...`
- `lib/cli/dist/ex agent plan <ticket>` against the composed test's fake console and against no connection; once 0223 exists, against a console started locally with a fake key.
