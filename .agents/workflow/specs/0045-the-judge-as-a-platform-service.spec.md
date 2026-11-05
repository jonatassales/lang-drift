# The judge as a platform service, and a run that brings the runtime up first

**Status:** ready-for-agent

## Tickets

- [x] 0221: `ex agent init` and the LangGraph scaffold leave the CLI
- [ ] 0222: one judge contract for the CLI and the console
- [ ] 0223: the console judges at `/api/judge`
- [ ] 0224: the planner asks the console, and no plan exists without a judge
- [ ] 0225: the run brings the runtime up before the first stage
- [ ] 0226: `ex agent stop` abandons a waiting run
- [ ] 0227: the documents say what the engine does
- [ ] 0228: release `@langdrift/cli` 0.5.0 and deploy the judge (owner)

## Problem Statement

The owner ran this repository's own microruntime with the published CLI (ticket 0198) and found that the engine did not do what the platform says it does.

- **Nothing came up.** `ex agent run` opened the money gate, reserved the run and stopped at the first stage waiting for a person, with no image built, no cluster and no pod. Since spec 0016, the runtime comes up at the end of the run, after every stage has been played; spec 0032 wrote the opposite ("the run brings the components up and pauses") and the two were never reconciled. A workflow whose stages all wait for a person, the workflow this repository committed, therefore pauses with nothing running.
- **No judge judged.** The plan was made with `--no-judge`, so the suggestion of model and effort came from an input, not from reasoning, while the session printed `agent ›` lines as if an agent had reasoned. The judge that exists is three model calls embedded in the CLI, and its adapter assumes a provider key in the consumer's environment, which no spec ever decided: the platform never asked a consumer to provide a model.
- **A placeholder posed as the agent.** `ex agent init NAME` still generates a two-node LangGraph scaffold into `.langdrift/agents/`, a Kubernetes Job of a retired flow that nothing plans, runs or supervises. Spec 0016 removed the embedded scaffold and kept this command without saying why.
- **A waiting run has no exit.** A run paused at a stage blocks any new plan of its runtime ("the runtime is running") and no command abandons it.

## Solution

The agent of `ex` is its judge: always present, part of every plan, and a service of the platform rather than a call the consumer pays a provider for.

- **The judge runs on the platform.** `ex agent plan` sends its three structured questions, the facts and verdicts over the sources it read, the relevance of the episodes the RAG retrieved, and the suggestion of model and effort, to the console at `/api/judge`. The console calls the model with LangDrift's own key, validates the answer against the question's schema and returns it with what the call used and what it cost. The consumer configures no provider, exports no key and does not choose the judge's model.
- **Planning requires a connected project, and a plan without a judge does not exist.** `--no-judge` is removed. A project that is not connected, a console that cannot be reached, or a judge that refuses, refuses the plan with the reason; no suggestion "from the inputs" stands in for a judgment.
- **The judge's spend passes the gates, with no way around them.** Before the first question the planner reserves the planning forecast against the ticket's allocation and the per-run cap, as today; the ledger's planning line records what the console charged, with its provenance. Until the owner prices a plan, the console charges the list cost of the tokens it spent, marked illustrative.
- **The judge is the platform's.** Its prompts are versioned in the CLI; each plan records the prompt version and hash it was judged with. The consumer configures only the planner model through the harness, as today, within what the platform offers.
- **The RAG stays as it is.** Deterministic, local, rebuilt by `ex` from the runs' memories; the judge reasons over what it retrieves.
- **The run brings the runtime up first.** `ex agent run` creates the cluster, builds the images, applies the workloads and Secrets and waits for readiness before the first stage; a stage that waits for a person pauses with the runtime running; the stages work against it. The end of the run observes delivery, stops, cleans up and remembers, as today, but no longer brings anything up.
- **A waiting run can be abandoned.** `ex agent stop <runtime>` closes the open lap as abandoned, releases the reservation, tears the runtime down through the execution ledger and keeps the evidence; a new plan is allowed afterwards.
- **The scaffold leaves.** `ex agent init` and the LangGraph templates are removed; the `agent` word stays in the commands and in the session's `agent ›` lines, meaning the judge and the deterministic shell around it.
- **Connection until the account exists.** The console's deployment holds one `RUNTIME_SECRET` the owner sets; `ex agent connect` records the console URL and that secret in the project's private state, and `/api/judge` accepts that secret. The account, when it comes, issues one secret per project behind the same route.
- **Stages stay the consumer's.** The workflow's stages keep running through the driver the consumer configured, their own coding tool, with the environment it names.

This is `@langdrift/cli` 0.5.0 and the console's first route that is not a stub. Ticket 0198 stays open until both exist; its `plan.sh` then drops `--no-judge` and requires the project to be connected.

## User Stories

1. As an operator running a microruntime, I want `ex agent run` to bring the components up before the first stage, so that a workflow that waits for a person pauses with the runtime running.
2. As an operator, I want a stage that waits for a person to find the runtime ready, so that confirming it plays the stage against live software.
3. As an operator, I want `ex agent stop <runtime>` to abandon a run that waits, so that I can plan again without deleting local state by hand.
4. As an operator, I want an abandoned run recorded as such, with who and when, its reservation released and its evidence kept, so that nothing disappears in silence.
5. As a team member planning a ticket, I want every plan judged, so that the suggestion of model and effort is reasoning over my ticket and my history, never an input repeated back.
6. As a team member, I want to configure no model provider and export no key to plan, so that planning asks nothing of me beyond a connected project.
7. As a team member, I want a plan refused with the reason when the project is not connected or the judge cannot be reached, so that I never mistake a plan without judgment for a plan.
8. As a team member, I want the judge's calls paid against my ticket's allocation and within the per-run cap, with the gate pausing before the first call when they do not fit, so that planning is as accountable as running.
9. As a team member, I want the plan to record which judge judged it, by prompt version and hash, so that two plans can be compared and a changed judge is visible.
10. As a team member, I want the ledger's planning line to show what the platform charged and where that figure came from, so that the plan's cost is evidence, not an estimate.
11. As a team member, I want the RAG to keep retrieving the same episodes for the same project state, so that the judge's inputs are reproducible even when its answer is not.
12. As the platform, I want the judge's model and key to be mine, so that no consumer's provider, quota or account shapes a judgment.
13. As the platform, I want one contract for the three questions, shared by the CLI's adapter and the console's route, so that a change on one side is a change on both.
14. As the platform, I want the console to validate every answer against the question's schema before returning it, so that the CLI never parses free text.
15. As the platform, I want the route to answer only a request carrying the deployment's secret, so that the judge is not an open endpoint.
16. As the platform, I want the route to record each judgment's usage and cost without storing the consumer's sources beyond the request, so that the platform keeps evidence and not the consumer's code.
17. As the owner, I want the price of a plan to stay my decision, so that the spec does not fix pricing.
18. As the owner, I want the console to charge the list cost of the tokens it spent, marked illustrative, until I price a plan, so that the gate and the ledger work from the first release.
19. As the owner, I want one `RUNTIME_SECRET` per console deployment until the account exists, set by me on the Vercel project, so that connecting a project needs no feature that is not built.
20. As a consumer who connects a project, I want `ex agent connect` to record the console URL and the secret in the private state, so that `plan` finds them without flags.
21. As a consumer, I want the stages of my workflow to keep running through my own coding tool via the driver, so that the platform judges and measures without choosing how I deliver.
22. As a maintainer, I want `ex agent init` and the LangGraph scaffold gone, so that no retired flow ships as a feature.
23. As a maintainer, I want `--no-judge` gone from the CLI, its tests and its documentation, so that no path plans without a judge.
24. As a maintainer, I want the session's `agent ›` lines to be true, so that a line about reasoning appears only when the judge reasoned.
25. As a maintainer, I want the run's order of operations tested at the composed CLI, so that a change to bring-up or stop is caught by the transcript and the written tree.
26. As a maintainer, I want the console's judge route tested with a fake model, so that its contract, its authentication and its refusals are covered without a paid call.
27. As a maintainer of this repository, I want `platform/microruntime/plan.sh` to drop `--no-judge` and to say the project must be connected, so that ticket 0198 can close on 0.5.0.
28. As a reader of spec 0032, I want its sentence about the run brought in line with this spec, so that the two no longer contradict each other.

## Implementation Decisions

- **One judge port, one contract.** The planner keeps asking its judge one operation, a question with a kind (`sources`, `episodes`, `suggestion`), the planner model, a system text, a prompt and a token bound, and receiving an answer with its text, the tokens used and whether the model refused. The console's route is the other side of that same contract: a request carries the question, the answer carries the text, the usage, the cost in cents with its provenance and the judge's version. The answer's text is JSON that the console validated against the kind's schema; the CLI validates it again.
- **The remote adapter replaces the provider adapter.** The CLI's judge adapter posts to `<console-url>/api/judge` with the deployment secret as a bearer token, a bounded timeout and no retry of a paid call; the provider SDK leaves the CLI. The console route holds the provider SDK, reads the platform's key from its environment at the call and refuses to start a call when it is missing.
- **The judge's model is the platform's.** The request carries the planner model the consumer configured through the harness; the console accepts only the models it offers, named in its own environment, prices the call at its list prices for that model and refuses any other with 400 naming the offered set. The consumer picks within what the platform offers; the platform holds the key.
- **Prompts live in the CLI.** The three prompts stay versioned in the CLI (today `judge-1`); the request carries the system text and the prompt, the console does not hold prompts. Each plan records the prompt version and the hash of the prompts in its plan fence, beside the judge's model.
- **No plan without a judge.** The `--no-judge` flag, the "inputs' own" suggestion path and the `by: input` suggestion are removed from the planner. A project without a recorded connection, an unreachable console, a route that refuses the secret or a model that refuses all refuse the plan with a message naming the cause; a blocked plan is not written as proposed.
- **Money as today, closed.** The planner reserves the planning forecast at the planner model's list price before the first question (FG gates), watches spend during, and writes the planning line of the ledger with the cost the console returned, its provenance (`console`, illustrative or priced) and the usage; the ledger's planning line gains those provenance fields. The CLI parses and validates every answer against the contract again before using it; the planner's existing errors for an unavailable model and a gate pause stay, and the new refusals join them. The console computes the illustrative cost from the list prices it holds for its own model until the owner prices a plan; that price is an open decision recorded in the spec, not in code.
- **The route.** `apps/console` gains `/api/judge` as a route handler that owns its methods, with one interceptor per concern beside it (method, bearer secret in constant time, content type, body size, schema of the question kind) and a provider adapter under the console's infrastructure folder. It answers 401 without the secret, 400 for a question outside the contract, 503 when the model is unavailable, 200 with the validated answer. It stores no source text; it logs usage and cost through the console's observability path.
- **Connection.** `ex agent connect` records the console URL and the `RUNTIME_SECRET` it reads from the environment in the project's private state and verifies them with one authenticated connection check against the judge route (a `GET` that answers 204 with the secret and 401 without), instead of the console's MCP route, which the console does not serve; `plan` reads the connection from there and refuses without it. The console's deployment-level secret is set by the owner on the Vercel project; when the account exists it issues one secret per project behind the same route and the CLI does not change.
- **Bring-up before the first stage.** The run's order becomes: open the gate and reserve; materialize the runtime from the fence through the native lifecycle, build the images, claim an execution, start, wait for readiness and print the runtime-ready line; then play the stages, pausing at the first that waits for a person with the runtime running; at the end observe delivery, stop and clean up by the ledger, release the unused reservation, settle the lap and write the memory. The bring-up timeout is the run's own, long enough for four image builds, and the README names it.
- **`ex agent stop`.** A new command on a runtime whose run is waiting: it closes the open lap with the outcome `abandoned`, who and when, releases the reservation in the ledger, stops and cleans up the runtime through the execution ledger as the end of a run does, keeps the evidence and the memory, and leaves the runtime plannable again. A run that is not waiting is refused with its state.
- **Removals.** `ex agent init`, the agent scaffolder port, the LangGraph templates and their tests leave the CLI; the README and the context document stop naming them. The session catalog keeps `agent ›` as the actor, and the lines about reasoning print only when the judge answered.
- **This repository.** `platform/microruntime/plan.sh` drops `--no-judge`; its README says the project must be connected and names the planning cost the gate reserves; ticket 0198 is revisited once 0.5.0 and the route are deployed. Spec 0032's sentence about the run is amended to this order.
- **Release.** `@langdrift/cli` 0.5.0, published by the owner as 0.4.0 was, with the tag and the GitHub Release; the core module gains a tag for its judge surface, which the CLI requires by version; this repository's root pin moves to 0.5.0 so its `plan.sh` runs the new engine; the console route ships with the console's deployment and the owner sets `RUNTIME_SECRET` and the platform's provider key on the Vercel project.

## Testing Decisions

A good test here observes what a user sees and what a file holds: the printed session lines, the exit code, the `.langdrift/` tree, the route's status and body; never a type or a function.

- **The composed CLI** stays the highest seam, as spec 0016 set it: the composition root accepts the injected ports (clock, id source, judge, lifecycle, stage driver, delivery observer) and one test reproduces a whole session in a disposable project seeded from fixtures. The judge port is injected with a fake that answers the three kinds; the run's new order is proven by the transcript (runtime-ready before the first `waiting` line) and by the written tree (evidence of the execution before any stage record); `stop` is proven by the lap's outcome, the ledger's released reservation, the lifecycle fake's stop call and a plan accepted afterwards; a plan without a connection, with an unreachable console or with a refusing judge is proven by the message and the absence of a proposed plan.
- **The remote judge adapter** is tested against an in-process HTTP test server that plays the console: the request it sends, the bearer it carries, the answer it parses, the refusals it maps to the planner's errors, and that a paid call is never retried.
- **The console route** is tested as the other route handlers are, with its interceptors and a fake provider adapter: 401 without the secret, 400 for a malformed question, 503 when the model is unavailable, 200 with a validated answer, the cost and usage it reports, and that the body holds no source text afterwards.
- **The contract** has one schema used by both sides' tests, so a drift between the Go adapter and the TypeScript route fails a test on the side that moved.
- **This repository's inputs** are proven by `sh platform/microruntime/plan.sh` against a connected project once the route is deployed; until then, the CLI's composed test covers the same path with the fake.
- Prior art: the composed-CLI tests at the CLI package root, the in-memory port fakes of the application package, the website's route-handler tests with interceptors, the console's proxy and route tests.

## Out of Scope

- The console account: who issues per-project secrets, sign-in and billing. Until it exists, one deployment-level secret.
- The price of a plan, which is the owner's; the console charges the illustrative list cost until then.
- The console's MCP route (`/api/mcp`) and the connected agent's supervision; `connect` records the connection, nothing else in that protocol changes here.
- Running the judge inside the microruntime, a local model, or embeddings by model in the RAG.
- Any change to the driver model: the stages keep running through the consumer's tool.
- Closing ticket 0198's run criteria, which wait for this release and the deployed route.
- The website's simulated session: whether its lines still describe the engine after this change is checked in its own ticket, not here.

## Further Notes

- The grilling of 2026-11-04 reversed an earlier direction within the same conversation: an editable LangGraph agent in `.langdrift/` was considered and rejected; the owner decided the agent is the judge, deterministic RAG stays, and the platform never asks a consumer for a model.
- `ex agent plan` on 0.4.0 printed `agent › reasoning` lines with `--no-judge`; the session catalog mirrors the website's terminal and said "agent" when nothing had reasoned. This spec makes those lines conditional on a judgment.
- The open decisions of this spec are the price of a plan and the account; both are recorded as open, and nothing in code presumes their answer.
