---
status: accepted
---

# Perfect Simulation Alignment

The session that plays on the website's home is the definition of LangDrift's main feature: `dft agent plan <ticket>`, then `dft agent run`, in a terminal beside the Ticket, Analysis, Plan and Workflow panel. The product does exactly what that simulation shows, in the same order, with the same information. Doing more than it shows is wrong and is removed. Doing less is wrong, and what is missing becomes a ticket that builds it. Doing it differently is both: the difference is removed and the simulated behavior is built. Specs, tickets, ADRs and READMEs that disagree are stale. The simulation's script and data live in `apps/website/src/app/[locale]/(pages)/platform/microruntimes/` and its copy in `microruntimes.json`.

The data matters more than the screen that shows it. The console is an interface that receives what the mechanisms produce. What has to be right is what `dft` and the client's `.langdrift/` folder produce and reason: every line of the terminal and every field of the four cards is data they produce, and the card is the specification of that data. The `console ›` prefix on the terminal's lines is part of the output, not a statement that the console does the reasoning.

The LangDrift CLI keeps three commands: `dft init`, which creates the `.langdrift/` folder the other two need, `dft agent plan` and `dft agent run`. Every other command is removed. A human checkpoint that the simulation shows ("waiting planning · the owner confirms") is answered inside `dft agent run`.

On 2026-11-05 the owner closed every spec and ticket still to implement as won't fix under this ADR, whatever it concerned, so the backlog restarts from zero; a ticket whose acceptance criteria were all checked was delivered, and is marked completed instead. New specs are written one per element of the feature, from the LangDrift CLI and Core to the client's own components as `dft` reads them, with the console and its stub API left out. Each spec states what its element is responsible for, what it produces under `.langdrift/` and which terminal lines and card fields it owns. Each ticket's acceptance criteria quote the exact lines and fields it must produce, and every ticket carries the canonical rule: what is done beyond the simulation is removed, what is missing becomes a new ticket. Item 4 of `.claude/rules/12-console.rule.md` holds the same rule.

## Considered Options

- **Treating the specs as the truth and the simulation as marketing.** That is how the code drifted: the specs added gates, commands and checks the product never shows, and left out what it shows, such as the dev ticket and Safe Seat lines.
- **Allowing extra behavior that does not contradict the simulation.** Every such allowance is a place where the product and its definition part ways. Information the product needs and the simulation lacks enters the simulation first.
- **Revising the open specs one by one.** Most predate the simulation's current script, and their status lines no longer say what was implemented. Closing them all and writing per-element specs from the simulation is shorter and leaves no stale decision standing.

## Consequences

- The backlog is empty until the new specs and their tickets are written.
- ADR 0004's derivation of a component from its Dockerfile is superseded wherever the simulation shows something else: it inspects each component's manifest (`go.mod`, `package.json`, `nats.conf`, `redis.conf`, `next.config.ts`) and resolves images. The spec for the client's components says which parts of ADR 0004 survive.
- How a failure the simulation does not show is handled is decided spec by spec, under the canonical rule.
- Where the simulation contradicts itself, such as the Plan card still reading "pending team authorization" after the run, the specs follow it literally and record the contradiction as an open question for the owner.
