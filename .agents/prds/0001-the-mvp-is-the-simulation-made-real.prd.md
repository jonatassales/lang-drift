# The MVP is the simulation made real

**Status:** accepted

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Specs:** 0048 to 0066 in `.agents/workflow/specs/`. **Evidence:** `.agents/workflow/research/simulation-conformance.research.md`.

## Goal

LangDrift's MVP is the session on the website's home, made real. A team points `dft` at one of its own tickets, and what happens is exactly what the simulation shows. Nothing more, nothing less.

## Target

The MVP is won when a real run on a client's own artifacts reproduces the simulation exactly. The artifacts are the ticket, the documents it links, the guardrail, the project's earlier runs and the components the spec names. Exactly means:

- **The same steps, in the same order.**
  - `dft agent plan <ticket>` reads the ticket and its links, judges the guardrail, retrieves the episodes, suggests a model and effort, resolves and sizes the components, and proposes `runtime.md`.
  - `dft agent run` opens the gate, waits where a person must confirm, plays the stages, waits for the dev ticket, reviews, brings the microruntime up, records the usage, pushes, opens the pull request, follows the checks, meets Safe Seat, stops the runtime and releases what was not spent.
- **The same information.** Every terminal line and every field of the Ticket, Analysis, Plan and Workflow cards is present, with the same wording and units. The values come from the client's artifacts and the real run, not from the simulation's sample.
- **The same motion and effects.** Lines appear as the work produces them. Each card fills when the terminal reaches its step. Waits hold until the person acts. The run reveals the Workflow card's result the way the simulation does.
- **Possibly the same components.** The surface that shows the run may reuse the website's session components, since they already render this data. Whether to reuse them or build new ones is open.

## Why

The simulation is the product's promise: the pitch a visitor sees and the definition the owner has settled (ADR 0005). Every gap between it and the engine is a broken promise. Behavior the simulation does not show makes the product harder to explain and to trust. Behavior it shows and the engine lacks is a claim the product cannot keep. Closing both gaps on real artifacts proves the whole chain: product intent, delivery context, execution evidence and economic outcome, accountable to the team's budget.

## Users

- **The team's owner** writes the ticket and its documents, chooses the human checkpoints, model and effort, confirms planning, writes the dev ticket and decides the merge.
- **The coding agent** is driven by `dft` through the stages with the chosen model and effort.
- **A viewer** of the run, in a terminal or on the surface that receives its data, sees the same lines and cards the simulation shows.

## Scope

Each element of the feature has a spec, and each spec says what the element does, what it produces under the client's `.langdrift/` and which lines and card fields it owns.

| Spec | Element |
| --- | --- |
| 0048 | The LangDrift CLI: `init`, `plan` and `run` only, and the shared output format |
| 0049 | The LangDrift Core contracts |
| 0050 | The hand-off of the card data |
| 0051 | The judge |
| 0052 | The RAG over project episodes |
| 0053 | The ticket and its linked documents |
| 0054 | The harness: guardrails and the envelope |
| 0055 | `runtime.md` and the sized plan |
| 0056 | The money gates |
| 0057 | The workflow and its human checkpoints |
| 0058 | The stage driver |
| 0059 | The agent skills |
| 0060 | The effort hook |
| 0061 | Safe Seat in the run |
| 0062 | The local cluster |
| 0063 | Push, pull request and ready |
| 0064 | The client's CI checks |
| 0065 | Evidence and the usage record |
| 0066 | The client's components |

The specs are executed one at a time, in this order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. A spec's tickets are written when an agent starts that spec, not in advance, and they are implemented before the next spec starts.

## Success criteria

1. The simulation's own ticket 0031, set up as a client project, produces the simulation's transcript exactly and every field of the four cards.
2. A real client ticket, run end to end, prints the same lines in the same order. Every card field carries that run's values, with nothing added and nothing missing.
3. Each Workflow choice changes the run as the simulation does:
   - a checked stage waits;
   - an unchecked Review ends "to merge";
   - another model or effort changes the estimate, the reservation and the observed spend;
   - a reservation over the budget pauses at the gate with nothing paid for.
4. The surface that receives the run shows it with the simulation's motion and effects.
5. `dft` has no command but `init`, `plan` and `run`, and no line or field the simulation does not show.

## Non-goals

- Anything the simulation does not show: other commands, gates, records printed nowhere, console areas beyond the run.
- The console's own screens and its stub API, except as the surface that receives and shows the run's data.
- Pricing, plans and limits beyond what the simulation states.

## Open questions

1. The Plan card still reads "pending team authorization" after the run (research note §1.6).
2. The terminal's and the card's "4 documents" count different sets.
3. The `reads … FIN-02 · FIN-04` line omits FIN-03.
4. Where the judge's calls go, and who pays for planning, given "nothing … paid for in this step".
5. Whether the surface that shows a real run reuses the website's session components or rebuilds them with the same motion.
