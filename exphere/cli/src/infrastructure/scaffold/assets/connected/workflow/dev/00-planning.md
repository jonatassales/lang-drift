# Planning

**Input:** the user's request, available project context, and any existing PRD, spec, tickets, or evidence. Apply the [cycle contract](CYCLE.md); discovery depth depends on material uncertainty rather than a mandatory interview.

## Work

1. Read AGENTS.md (consumer-owned reference: `../../../AGENTS.md`), context (consumer-owned reference: `../../../.agents/CONTEXT.md`), domain configuration (consumer-owned reference: `../../artifacts/design/0007-DD-domain-model.dd.md`), and the tracker (consumer-owned reference: `../../artifacts/design/0006-DD-artifact-tracker.dd.md`). Inspect the current tree and preserve unrelated changes.
2. Identify the problem, intended users, outcome, owner, scope, and evidence gaps. Use [research](../skills/research/SKILL.md) for needed primary-source investigation. Use [grill-with-docs](../skills/grill-with-docs/SKILL.md) when a discovery discussion is explicitly requested; do not make it a prerequisite for every entry point.
3. Record consequential domain decisions with [domain-modeling](../skills/domain-modeling/SKILL.md). Use the available domain context; an absent glossary is not an instruction to recreate a deleted file.
5. When ticket preparation is authorized, use [to-tickets](../skills/to-tickets/SKILL.md). Give each ticket a complete, bounded deliverable, criteria, and genuine blocking edges. Keep stable IDs and the configured artifact paths and naming.
6. Identify the highest useful verification surface and the applicable guardrails. Respect scoped test authorization and restrictions on unrelated tests; documenting a future verification option does not authorize executing it. Record scope acceptance and execution authorization separately from artifact creation.

**Output:** only the requested typed artifacts, with source context, decisions, acceptance criteria, dependencies, verification approach, and remaining uncertainty. New proposals remain drafts unless actual authorization covers their scope and revision.

**Completion:** the requested planning deliverable is reviewable without guessing its purpose or next decision. Stop at that deliverable if implementation has not been requested. If implementation is already authorized and its prerequisites are met, [execution](01-execution.md) may proceed without asking again.

**Return condition:** unresolved facts that materially change the scope remain explicit. Continue independent research or drafting; do not invent an approval to advance.
