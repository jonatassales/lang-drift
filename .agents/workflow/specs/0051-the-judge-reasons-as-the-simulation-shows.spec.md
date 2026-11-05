# The judge reasons as the simulation shows

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E4.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 9 of 19, after spec 0052 and before spec 0066. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The simulation's plan reasons in six lines. It judges the guardrail's criteria, each in a few words. It judges why three episodes are kept ("share the stack and scale") and why three differ ("other runtime, no broker, docs only"). It concludes that sizing follows the spec, not history. It weighs the efforts and the models it did not pick, and suggests one model and one effort, which the Analysis card shows as a recommendation the team authorizes. All of it happens in a step that pays for nothing. The LangDrift CLI's judge reasons differently:
- Its planning calls are priced, reserved and gated through FG-01, FG-02 and FG-03, and can refuse the plan.
- Without a judge, the plan prints judge-unavailable or judge-stopped warnings, and a `--no-judge` mode drops the judging and weighing lines.
- The guardrail line is the judge's free "id applies" text.
- On episodes, the judge only rewrites reason text and never keeps or drops one, while "share the stack and scale" is fixed text.
- The sizing note reads "follows the inputs, not history".
- The judge's figure discrepancies print lines the simulation does not have.

A team cannot trust that the reasoning it reads is the reasoning the product does.

## Solution

The judge is the plan's reasoning step. It answers three questions:
- what each guardrail criterion asks, in a few words;
- for each episode candidate, whether it shares the ticket's stack and scale or why it differs, from a closed set of reasons;
- which model and effort the ticket warrants, with one reason against each alternative it did not pick.

It also states one conclusion that follows from the retrieval: with no cited episode that measured cpu or memory, sizing follows the spec. Its answers are checked against fixed rules before any line prints. They feed the six lines and the Analysis card's suggestion fields exactly as the simulation shows them. Nothing about the judge is priced, reserved or gated in the plan, and no line about it exists beyond the six. Where the judge runs, and who pays for its calls, is an owner decision recorded below. The judge sits behind a port, so that decision changes the port's implementation and nothing else.

## Responsibility

This element:
- asks the judge the three questions and checks each answer;
- prints the two `judging` lines, the `reasoning` line, the two `weighing` lines and `✓ suggests`;
- decides each episode candidate's reason, the kept reason or one of the closed drop reasons, consistent with the threshold of spec 0052;
- states the sizing basis;
- owns the Analysis card's header, SUGGESTED MODEL, SUGGESTED EFFORT and "recommendation · the team authorizes".

It does not:
- read documents or extract facts (spec 0053), or read the guardrail's criteria (spec 0054);
- embed, search, rank or apply the 0.55 threshold, print the candidate, drop or cited lines, or show the EPISODES rows and MEASURED RESOURCES (spec 0052);
- record the suggestion in `runtime.md` (spec 0055) or decide its authorization state (spec 0050);
- offer the model and effort choices or price them (specs 0057 and 0056), or set the Workflow card's default choice, which is the workflow's (spec 0057); the suggestion is a recommendation shown on the Analysis card and recorded in the draft;
- size components (spec 0066).

## Produces

No file of its own. Its values go into the plan's card data (spec 0050) and the `runtime.md` draft (spec 0055):

- **Criteria statements**: per guardrail criterion, `id` and `statement` (`FIN-02` "estimate before a paid call", `FIN-03` "a budget before a paid run", `FIN-04` "size cpu and memory to the load"). They feed the first `judging` line and, through spec 0054, the card's FIN line.
- **Episode judgments**: per candidate, `id`, `verdict` (`kept` or `dropped`, equal to the threshold's) and `reason`:
  - `shares the stack and scale` for a kept candidate;
  - one of `other runtime`, `no broker`, `docs only` for a dropped one.

  They feed the second `judging` line. Through spec 0052 they also feed each EPISODES row's verdict and reason.
- **Sizing basis**: `spec`, with the statement "no cited episode measured cpu or memory · sizing follows the spec, not history". It feeds the `reasoning` line. Through spec 0052 its last part, "sizing follows the spec, not history", feeds the MEASURED RESOURCES field.
- **Suggestion**: `model` (`claude-sonnet-5-5`), `effort` (`medium`), `effortReason` ("low skips the broker wiring · high is not warranted by 5 components"), `modelReason` ("Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome") and `by` (`recommendation`). It feeds the weighing lines, `✓ suggests`, SUGGESTED MODEL, SUGGESTED EFFORT, the header's last state and, through spec 0055, the draft's model and effort.

## Shows

Phase 1, `dft agent plan`, at three moments.

After `reads      harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04` (spec 0054), and before `bounded …`:

```
console › judging · FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load
```

After `console › 6 candidates · similarity 0.81 to 0.42`, and before `console › dropped 3 below 0.55 · kept 3` (spec 0052):

```
console › judging · 3 share the stack and scale · 3 differ: other runtime, no broker, docs only
```

After `console › 3 episodes cited · none with measured resources` (spec 0052), and before `# plan · resolving the components, then sizing them inside the envelope` (spec 0048):

```
console › reasoning · no cited episode measured cpu or memory · sizing follows the spec, not history
console › weighing effort · low skips the broker wiring · high is not warranted by 5 components
console › weighing model · Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome
✓ suggests   claude-sonnet-5-5 · effort medium
```

Card fields, Analysis card:

- **Header**, in three states:
  - `4 documents read` from the `✓ read` line (count from spec 0053);
  - `3 of 6 episodes kept` from the `dropped` line (counts from spec 0052);
  - `Claude Sonnet 5.5 · effort medium` from `✓ suggests` on.

  Before the first state it is not set.
- **SUGGESTED MODEL**: a placeholder from the `RAG over project episodes …` line, then `Claude Sonnet 5.5` with `claude-sonnet-5-5` beneath it, from `✓ suggests`.
- **SUGGESTED EFFORT**: a placeholder from the same line, then `medium` from `✓ suggests`.
- **`recommendation · the team authorizes`**, from `✓ suggests`.

## User Stories

1. As a team lead, I want the plan to state what each guardrail criterion asks in a few words, so that I can see the plan understood our rules.
2. As a reviewer, I want the criteria judged in the guardrail's order, all of them, so that none is silently skipped.
3. As a developer, I want the plan to say how many episodes share the ticket's stack and scale, so that I know what precedent the suggestion rests on.
4. As a developer, I want the plan to say why the others differ, in a small fixed vocabulary, so that the reasons read the same on every plan.
5. As a reviewer, I want the drop reasons listed once each in a fixed order, so that the line is stable and comparable between plans.
6. As a reviewer, I want the judge's keep or drop to agree with the similarity threshold, so that "dropped 3 below 0.55" is never contradicted by the line before it.
7. As a team lead, I want the plan to say that no cited episode measured cpu or memory and that sizing therefore follows the spec, so that I know where the component figures come from.
8. As a team lead, I want the plan to weigh the lower and the higher effort against the ticket, so that the suggested effort comes with its reasons.
9. As a team lead, I want the plan to weigh the cheaper and the dearer model against the ticket, so that the suggested model comes with its reasons.
10. As a developer, I want `✓ suggests` to print the model's id and the effort, so that I can copy them exactly.
11. As a budget owner, I want the suggestion limited to the three models the Workflow card offers and to low, medium or high, so that the plan never suggests what the run cannot use.
12. As a team member, I want the Analysis card to call the suggestion a recommendation the team authorizes, so that nobody reads it as a decision.
13. As a console user, I want the Analysis header to step from documents read, to episodes kept, to the suggested model and effort, so that the card's summary follows the act.
14. As a console user, I want SUGGESTED MODEL to show the model's name with its id, so that the name I recognize and the id I type sit together.
15. As a budget owner, I want nothing about the judge priced, reserved or gated in the plan, so that the plan step pays for nothing, as its scope comment says.
16. As a developer, I want no `--no-judge` mode and no judge warnings in the terminal, so that every plan prints the same lines.
17. As a developer, I want an answer that breaks the rules to stop the plan with nothing written, so that a plan never shows reasoning that was refused.
18. As a platform owner, I want the judge behind one port, so that deciding where it runs changes one implementation and no line.
19. As a test author, I want a recorded judge for the fixture, so that the golden test compares the six lines byte for byte.
20. As an auditor, I want the suggestion recorded as a recommendation pending authorization, so that the plan's evidence keeps it apart from what the team authorizes.

## Implementation Decisions

**Remove**
- Pricing, reserving and gating the planner's calls: FG-01, FG-02 and FG-03 on planning, the planning forecast, estimate and reservation, the planning spend line in the ledger, and the refusal of a plan at a planning gate. Serves `# plan only · nothing is built, started or paid for in this step` (spec 0048), which no judging line may contradict. The ledger side of this removal is specs 0056 and 0065's.
- The judge's cost on each judgment and the planner's usage and price record.
- `--no-judge`, the judge-unavailable and judge-stopped warnings, and the "deterministic shell" that skips the judging and weighing lines.
- The judge's sources question (verdicts, confidences, roles, summaries, facts, figures). Document reading is deterministic now (spec 0053).
- The figure comparison and its "discrepancy … the input holds" lines.
- Fixed reason text for kept episodes and the judge's free-text rewriting of reasons. Serves the second `judging` line: reasons come from the closed set.
- The note "sizing follows the inputs, not history". Serves the `reasoning` line.
- The suggestion's free `reasons` list and any suggestion field no line or field shows.
- Falling back to a recorded model and effort choice as the "suggestion" when no judge answers.

**Build**
- The criteria question: the judge gets the guardrail's criteria with their rule text (spec 0054) and returns one statement per criterion. The answer is refused if it skips a criterion or names one that is not there. The line is `judging · ` followed by `<id> <statement>` per criterion, joined by ` · `, in table order.
- The episodes question:
  - The judge gets, for each candidate, the episode record and similarity (spec 0052), together with the read documents and extracted facts (spec 0053).
  - It returns a verdict and a reason for every candidate. A kept candidate's reason is `shares the stack and scale`. A dropped candidate's is one of `other runtime`, `no broker`, `docs only`.
  - The answer is refused when a verdict differs from the threshold's, a reason is outside the set, or a candidate is missing.
  - The line is `judging · <kept> share the stack and scale · <dropped> differ: ` followed by the distinct drop reasons in the set's order (`other runtime`, `no broker`, `docs only`), joined by `, `.
- The sizing conclusion: when spec 0052 reports no cited episode with measured resources, the basis is `spec` and the line is `reasoning · no cited episode measured cpu or memory · sizing follows the spec, not history`. This needs no question to the model.
- The suggestion question:
  - The judge gets the documents, the facts (among them the five components), the cited episodes, the three models of the price list with their short names (spec 0056), the same three the Workflow card offers, and the efforts low, medium and high.
  - It returns a model, an effort, an effort reason and a model reason.
  - The answer is refused when the model is not one of the three, the effort is not one of the three, the effort reason does not name each effort not chosen in ascending order, or the model reason does not name each model not chosen by its short name in ascending order of list price.
  - The lines are `weighing effort · <effortReason>`, `weighing model · <modelReason>` and `✓ suggests   <model id> · effort <effort>`.
- The Analysis header's three states and the suggestion fields, in the plan's card data (spec 0050). The model's name comes from the price list (spec 0056).
- One judge port with three questions, whose implementation is open (see Open Questions), and a recorded judge for the fixture.
- Safe stops: the judge cannot be reached, declines, or gives an answer the checks refuse. The plan stops before the line the answer would feed, writes no draft and exits non-zero, with the reason through the CLI's error channel (spec 0048). An answer is not retried silently into a different line, and no warning line is printed.

**Change**
- The judge decides the episodes' reasons, and its verdict must agree with the threshold, instead of only rewriting the threshold's reason text.
- The suggestion always comes from the judge and is recorded as a recommendation pending the team's authorization. The draft's model and effort are the suggestion's (spec 0055), not a pre-recorded choice.
- The suggestion's models are the price list's three (spec 0056), the ones the Workflow card offers, not any model a price file might list.

**Interfaces with neighbours**
- From spec 0053: the documents' text, the extracted facts and the documents-read count for the header.
- From and to spec 0054: the criteria in, the statements out.
- From and to spec 0052: candidates with similarity and the threshold's verdict in, reasons out. The cited count, the measured count and the kept and total counts for the header come in.
- From spec 0056: the price list's three models with ids, names and short names. The efforts are the Core contracts' low, medium and high (spec 0049).
- To spec 0055: the suggested model and effort, recorded as a recommendation.
- To spec 0050: the header states and the suggestion fields.

## Testing Decisions

- The seam is the golden test of spec 0048 over the shared fixture. This element adds its six lines, the header's three states and the three suggestion fields, byte for byte.
- The fixture's judge is a recorded judge behind the port. It answers:
  - the criteria question with the three statements;
  - the episodes question with the six verdicts and reasons of the simulation (ep-0412, ep-0388 and ep-0351 kept; ep-0297 docs only, ep-0240 other runtime, ep-0188 no broker);
  - the suggestion question with Sonnet, medium and the two reasons as shown.
- Golden cases for the safe stops: an unreachable judge, and an episodes answer that keeps a candidate below 0.55. Each prints up to the stop, writes no draft and exits non-zero.
- Unit tests for the rules only:
  - the answer checks;
  - the drop-reason ordering and de-duplication in the second `judging` line;
  - the ascending order required in the weighing reasons.
- No test reaches a real model.

## Out of Scope

- Document reading and the documents count: spec 0053.
- The guardrail's criteria and the FIN line on the card: spec 0054.
- Retrieval, the threshold, the EPISODES rows and MEASURED RESOURCES: spec 0052.
- The draft and its authorization state: specs 0055 and 0050.
- Model prices, the estimate, the budget and the gates of the run: spec 0056.
- The Workflow card's choices: spec 0057.
- The console and its stub API.

## Open Questions

- **Where the judge runs, and who pays for planning.** The plan's scope says "nothing is built, started or paid for in this step". Yet the judge is a model, and a call to it costs someone. Three answers are possible:
  - the team's own model credential, which contradicts "nothing … paid for" unless the cost is absorbed outside the run's budget;
  - a LangDrift platform service that bears the cost (the closed spec 0045's direction, without its planning reservation);
  - a local model.

  This spec builds the port and the recorded fixture judge, and removes every planning gate and reservation. The production implementation waits for this answer.
- **An episode the judge finds unlike, above the threshold.** The simulation's judge and threshold always agree. May the judge drop a candidate whose similarity is at or above 0.55, or keep one below it? Until decided, a disagreement is a refused answer and stops the plan.

## Further Notes

- The glossary's **Judge** entry still says its calls are paid work under the money gates and that without a credential the judged fields stay null. Under ADR 0005 that entry is stale. Updating it is a `domain-modeling` task for the owner, not part of this spec.
- The simulation never shows a cited episode with measured resources. Neither sizing from history nor the wording of the `reasoning` line in that case is built here (see spec 0052).
