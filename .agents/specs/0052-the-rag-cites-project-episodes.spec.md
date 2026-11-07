# The RAG cites the project's episodes

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E5.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 8 of 19, after spec 0056 and before spec 0051. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

To suggest a model and an effort, the simulation's plan runs a RAG over the project's episodes:
- it embeds the ticket intent into one query vector;
- it takes the top six episodes by cosine similarity;
- it drops those below 0.55 and cites the rest.

Each candidate appears on the Analysis card as `ep-NNNN` with a slug, a similarity, and kept or dropped with a reason. The card also says whether any cited episode measured resources.

The LangDrift CLI keeps one episode per runtime, identified by the runtime's id rather than `ep-NNNN` and a slug. It adds structured scoring (same component path, same stack, measured resources, measured tokens) and keeps an episode below 0.55 when two criteria match. When a cited episode measured resources, it replaces every component's share with an equal split of the peak × 1.5. None of this is shown. The episode a run leaves behind is written as a per-runtime memory fence that a new run of the same runtime overwrites. A team reading "kept" on the card cannot rely on it meaning "at or above 0.55".

## Solution

Every finished run leaves one episode under `.langdrift/`. It carries an `ep-NNNN` id, the slug of the ticket it delivered, the intent that was embedded, its components with their stack and role, and its measured resources, which stay null when nothing was measured. The plan:
1. embeds the ticket intent into one query vector with the local embedder, which makes no paid call;
2. takes the six most similar episodes by cosine from the project's index;
3. keeps exactly those at or above 0.55 and drops the rest;
4. hands the candidates to the judge (spec 0051) for their reasons;
5. reports how many cited episodes measured resources.

The six lines and the EPISODES and MEASURED RESOURCES fields carry exactly that. The structured scoring and the sizing from history are removed.

## Responsibility

This element:
- writes the episode of each finished run, and keeps the project's episode index with its vectors;
- embeds the ticket intent, searches the index, applies the strict 0.55 threshold and counts the candidates, the kept, the dropped and the cited episodes with measured resources;
- prints the six RAG lines;
- owns the card's EPISODES section (threshold and rows) and the MEASURED RESOURCES field.

It does not:
- decide any episode's reason, or print the `judging · 3 share …` line (spec 0051);
- state the sizing basis or print the `reasoning …` line (spec 0051), although its measured count is that line's input and the reasoning's last words appear in its MEASURED RESOURCES field;
- suggest a model or an effort (spec 0051);
- size components (spec 0066);
- decide when a run has finished, or write the run's other records (specs 0063 and 0065).

## Produces

- **Episode record**, one per finished run, at `.langdrift/microruntimes/episodes/<id>-<slug>.json`. Fields:
  - `id`: `ep-NNNN`, the next free number in the project;
  - `slug`: the delivered ticket's slug, for example `book-a-demo-runtime`;
  - `ticket`: the ticket's id;
  - `intent`: the exact text the episode is embedded from, the delivered ticket's title and body;
  - `components`: one entry per component the run ran, each with `path`, `stack` (`Go`, `Node.js`, `NATS`, `Redis`, `Next.js`) and `role` (`api`, `broker`, `cache`, `web`), in the run's order;
  - `resource`: the workload kind (`Deployment`), the count and the replicas each, as the plan's SUGGESTED RESOURCE stated it;
  - `measured`: `cpuMillis` and `memoryMiB` peaks, or `null` when the run measured neither.

  An episode with no components is one whose run changed no runtime. The judge reads `components`, `resource` and `intent` to give "share the stack and scale", "other runtime", "no broker" or "docs only". The `measured` field feeds the count in `3 episodes cited · none with measured resources`. A run that paused at the gate (spec 0056) leaves no episode.
- **Episode index**, at `.langdrift/agent/state/rag/index.json` (unchanged location). Fields: `version`, the embedder's `id` and `dimensions`, and one entry per episode with its `id`, `slug` and `vector`. Vectors are computed once, when an episode is written. The index is rebuilt from the episode records only when it is missing or was built by another embedder.
- **Retrieval** (into the plan's card data, spec 0050):
  - `threshold` (`0.55`);
  - `candidates`, in descending similarity, each with `id`, `slug`, `similarity` (two decimals), `verdict` and `reason` (the reason from spec 0051);
  - `kept`, `dropped` and `total` counts;
  - `cited` and `citedMeasured` counts.

  The kept and total counts also feed the Analysis header's second state (spec 0051). The candidates feed spec 0051's episodes question. `citedMeasured` feeds spec 0051's sizing conclusion.

## Shows

Phase 1, `dft agent plan`. After `budget     $25.00 per run for model calls · set by the team in the design doc` (spec 0056):

```
console › RAG over project episodes · a model and effort for this ticket
console › embedding the ticket intent · 1 query vector
console › vector search · project episodes · top 6 by cosine similarity
console › 6 candidates · similarity 0.81 to 0.42
```

then the judge's `console › judging · 3 share the stack and scale · 3 differ: other runtime, no broker, docs only` (spec 0051), then:

```
console › dropped 3 below 0.55 · kept 3
console › 3 episodes cited · none with measured resources
```

then the judge's `console › reasoning · …` (spec 0051).

Card fields, Analysis card:

- **EPISODES header**: `threshold 0.55`, shown with the card.
- **EPISODES rows**:
  - not set until the `vector search` line, which shows six placeholder rows;
  - from `6 candidates`, each row carries its id, slug and similarity, with a bar of the similarity;
  - from `dropped 3 below 0.55 · kept 3`, each row also carries its verdict:
    - `ep-0412 checkout-preview-sizing 0.81 kept`
    - `ep-0388 notify-api-rollout 0.74 kept`
    - `ep-0351 signup-redesign-preview 0.66 kept`
    - `ep-0297 docs-search-index 0.52 dropped · docs only`
    - `ep-0240 batch-export-worker 0.47 dropped · other runtime`
    - `ep-0188 storefront-cache-warmup 0.42 dropped · no broker`
- **MEASURED RESOURCES**: `none cited` from `3 episodes cited · none with measured resources`. `sizing follows the spec, not history` is added beside it from the judge's `reasoning` line (text from spec 0051).

## User Stories

1. As a team lead, I want the plan to say it runs a RAG over the project's episodes for a model and an effort, so that I know what the retrieval is for.
2. As a developer, I want the plan to embed the ticket intent into exactly one query vector, so that one ticket is one query.
3. As a budget owner, I want the embedding to run locally with no paid call, so that the plan step pays for nothing.
4. As a developer, I want the search to take the top six episodes by cosine similarity, so that the candidate set is small and predictable.
5. As a reviewer, I want the candidates line to give the highest and lowest similarity, so that I see the spread at a glance.
6. As a reviewer, I want every candidate at or above 0.55 kept and every one below dropped, so that "kept" always means the same thing.
7. As a reviewer, I want no structured score able to keep an episode below the threshold, so that the threshold the card shows is the only rule.
8. As a console user, I want each candidate shown with its `ep-NNNN` id and a slug, so that I can tell episodes apart and look one up.
9. As a console user, I want each row to show its similarity as a figure and a bar, so that I can compare candidates visually.
10. As a console user, I want dropped rows to say why, in the same short vocabulary as the terminal, so that the card and the terminal agree.
11. As a team lead, I want the plan to say how many cited episodes measured resources, and "none" when none did, so that missing measurements never read as zero.
12. As a team lead, I want MEASURED RESOURCES to say "none cited" with the reason sizing follows the spec, so that I know history did not set the figures.
13. As a developer, I want every finished run to leave one episode with its own id, so that the next plan can learn from every run, not just the last of each runtime.
14. As a developer, I want an episode to keep the components it ran with their stack and role, so that the judge can tell a matching precedent from one with no broker or another runtime.
15. As a developer, I want a run that changed no runtime to leave an episode without components, so that the judge can recognize a docs-only precedent.
16. As an auditor, I want an episode's measured resources stored as null when nothing was measured, so that unknown stays unknown.
17. As a developer, I want a run paused at the gate to leave no episode, so that only work that happened becomes precedent.
18. As a developer, I want the index to keep each episode's vector from when it was written, and to be rebuilt only when it is missing or the embedder changed, so that similarities do not drift between plans.
19. As a developer, I want no `reindex` step to remember, so that the index is always right without a command.
20. As a test author, I want the fixture's six episodes and index to reproduce 0.81, 0.74, 0.66, 0.52, 0.47 and 0.42 exactly, so that the golden test compares every figure byte for byte.
21. As a team lead, I want the episodes ordered by descending similarity on the card and in the counts, so that the strongest precedent reads first.

## Implementation Decisions

**Remove**
- Structured retrieval: the same-component-path, same-stack, measured-resources and measured-tokens criteria, the per-candidate criteria list and score, the method name "vector+structured", and the query's stacks and component paths. Serves `6 candidates · …` and `dropped 3 below 0.55 · kept 3`, which rank and keep by similarity alone.
- Keeping a candidate with two matched criteria below the threshold. Serves the strict `dropped 3 below 0.55`.
- Sizing from history: the strongest measured peak × 1.5, rounded up and split equally across components, the "measured-runs" sizing basis, and the note "N cited episodes measured cpu and memory · sizing follows their peaks with headroom". The simulation shows only "none with measured resources", so no history-based sizing is built.
- The episode as a per-runtime memory fence (`langdrift-run` in each runtime's `run.md`), read as the episode source, together with the runtime id as episode identity. Serves the `ep-NNNN` rows.
- Episode fields no line or field uses: measured tokens, outcome, facts and the run's title as a separate field.
- The threshold's reason strings ("similarity X · above the threshold", "… below 0.55 · no criterion matched") and the judge's free-text reasons. Reasons are spec 0051's closed set.
- The plan record's retrieval fields no field shows: method, criteria, keep-at-matches, score, and the note "first run of this project; no earlier run to cite".
- The `Reindex` service behind the removed `reindex` command (the command itself is removed by spec 0048).

**Build**
- The episode record described under Produces, written once per finished run, numbered `ep-NNNN` from the highest existing number plus one, and added with its vector to the index at the same moment. It prints nothing. Serves the next plan's rows.
- The strict threshold: a candidate is kept exactly when its similarity, rounded to two decimals, is at least 0.55.
- Candidate selection: the six episodes with the highest cosine similarity to the query vector, in descending similarity, ties broken by ascending id.
- The six lines, from the counts: `6 candidates · similarity <highest> to <lowest>`, `dropped <dropped> below 0.55 · kept <kept>`, `<cited> episodes cited · none with measured resources` (shown when `citedMeasured` is 0).
- The retrieval card data, with each candidate's reason as spec 0051 returns it.

**Change**
- The query text is the ticket intent: the ticket's title and body, the same shape an episode keeps as its `intent`. The embedder stays the local deterministic one, so "1 query vector" makes no paid call.
- The index stores each episode's vector from when the episode was written, instead of re-embedding every runtime's memory on each load.
- The episode is written at the end of every finished run (after `✓ ready`, spec 0063) instead of being the runtime's memory overwritten by its next run.

**Interfaces with neighbours**
- To spec 0051: the candidates with their records, similarities and threshold verdicts, plus `cited`, `citedMeasured`, `kept` and `total`. From spec 0051: one reason per candidate.
- From spec 0053: the ticket's id, slug, title and body for the query and, at the end of the run, for the episode.
- From specs 0055 and 0066, at the end of the run: the components with stack and role, and the suggested resource, for the episode.
- From the run (specs 0063 and 0065): the moment the run is finished, and that it was not paused at the gate.
- To spec 0050: the EPISODES and MEASURED RESOURCES data.

## Testing Decisions

- The seam is the golden test of spec 0048 over the shared fixture. This element adds its six lines, the EPISODES header and rows, and MEASURED RESOURCES, byte for byte.
- The fixture holds the six episode records of the simulation. Each record's components and resource let the recorded judge's reasons be plausible:
  - ep-0412, ep-0388 and ep-0351 share the stacks and component roles;
  - ep-0297 has no components;
  - ep-0240 runs another workload;
  - ep-0188 has no broker.

  All six have `measured: null`. The fixture's index holds the six vectors, built by the test from the fixture ticket's query vector so that their cosines are exactly 0.81, 0.74, 0.66, 0.52, 0.47 and 0.42. Its embedder id matches the local embedder's, so no rebuild happens.
- After `dft agent run` completes in the golden test, the test asserts that one new episode, `ep-0413-book-a-demo-runtime.json`, exists with the run's components and `measured: null`, and that the index holds seven entries.
- After the over-budget variant of `dft agent run`, the test asserts that no episode was written.
- Unit tests only for:
  - the threshold at its edge (0.55 kept, 0.54 dropped);
  - the selection order with ties;
  - the next `ep-NNNN` number.

## Out of Scope

- The reasons for each candidate and the `judging`, `reasoning`, `weighing` and `suggests` lines: spec 0051.
- Sizing components: spec 0066. Recording the sizing: spec 0055.
- Deciding that a run finished, and the run's evidence and usage record: specs 0063 and 0065.
- The `reindex` command and every other removed command: spec 0048.
- The console and its stub API.

## Open Questions

- **Fewer than six episodes.** A new project has no episodes, and a young one fewer than six. The simulation only shows six. Two things are undecided:
  - what the lines `6 candidates · similarity 0.81 to 0.42` and `3 episodes cited …` become with fewer candidates, or with none;
  - what the EPISODES section shows.

  Until the owner adds that case to the simulation, the lines keep their templates with the real counts. With zero candidates, the plan prints `0 candidates` without the similarity part. That is a wording the simulation does not yet show, and the owner may replace it.
- **A cited episode with measured resources.** Runs as the simulation shows them measure no cpu or memory, so every episode this product writes has `measured: null`, and "none with measured resources" always holds. If episodes gain measurements, the `3 episodes cited …` line, MEASURED RESOURCES and the sizing basis need wording and a rule that the simulation does not show. They are out of scope until it does.

## Further Notes

- The slugs in the simulation (`checkout-preview-sizing`, `docs-search-index`, …) are illustrative. In a real project they are the slugs of the tickets those runs delivered.
- The glossary's **Episode** entry ("one earlier run of the project as the RAG remembers it: the memory's fence, with its stacks, component paths, measured peaks and tokens") describes the per-runtime memory this spec replaces. Updating it is a `domain-modeling` task for the owner.
