---
name: to-tickets
description: Break a plan, spec, or the current conversation into a set of tracer-bullet tickets, each declaring its blocking edges, published to the configured tracker (edges as text in one file per ticket locally, or native blocking links on a real tracker).
metadata:
  source: https://github.com/mattpocock/skills/tree/5beb35b/skills/engineering/to-tickets
  adapted: Adjusted to this repository's layout; checkpoints follow the workflow policy; every ticket carries a Safe Seat dev part.
---

# To Tickets

Break a plan, spec, or conversation into a set of **tickets**: tracer-bullet vertical slices, each declaring the tickets that **block** it.

The issue tracker and triage label vocabulary are defined in the `Spec-driven development` section of `.agents/context.md`.

## Process

### 1. Gather context

Work from whatever is already in the conversation context. If the user passes a reference (a spec path, an issue number or URL) as an argument, fetch it and read its full body and comments.

For a spec decided by ADR 0005 (Perfect Simulation Alignment), first play the simulation in production at https://langdrift.md/en/. Watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. Change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. Compare what it shows with the spec, line by line and field by field. Where they differ, the simulation in production wins: report the difference and write the tickets to the simulation. Write tickets only for the spec whose turn it is in the specs' Execution order, when an agent starts that spec, never for later specs in advance.

### 2. Explore the codebase (optional)

If you have not already explored the codebase, do so to understand the current state of the code. Ticket titles and descriptions should use the project's domain glossary vocabulary, and respect ADRs in the area you're touching.

Look for opportunities to prefactor the code to make the implementation easier. "Make the change easy, then make the easy change."

### 3. Draft vertical slices

Break the work into **tracer bullet** tickets.

<vertical-slice-rules>

- Each slice cuts a narrow but COMPLETE path through every layer (schema, API, UI, tests): vertical, NOT a horizontal slice of one layer
- A completed slice is demoable or verifiable on its own
- Each slice is sized to fit in a single fresh context window
- Any prefactoring should be done first

</vertical-slice-rules>

Give each ticket its **blocking edges**: the other tickets that must complete before it can start. A ticket with no blockers can start immediately.

Give each ticket its **simulation alignment**, as ADR 0005 and item 4 of `.claude/rules/12-console.rule.md` require. A ticket that removes what the code does beyond the simulation, builds what the simulation shows and the code lacks, or changes what differs says which of the three it is. Its acceptance criteria quote the exact terminal lines and card fields of the home simulation it must produce.

Give each ticket its **dev part**: the problem in it that the owner writes by hand, as `.agents/policy/safe-seat.policy.md` sets out.

- The problem has its own goal and acceptance criteria.
- Its category is one of those in the policy's Categories section: `core-domain`, `hard-part` or `new-stack`.
- It is never one of the kinds the policy excludes.
- It reaches the policy's Floor: at least 5% of the lines the ticket will add.
- It is located by paths that only the owner writes, including their tests, so that `implement` can leave them alone.
- The boundary the rest of the code calls (a signature or an interface) is named in the ticket and sits outside those paths.

Your proposal stands as the decision unless the owner changes it, as the policy's Decisions section says. When no part of a ticket qualifies, write an override with its reason and date instead, as the policy's Override section says.

**Wide refactors are the exception to vertical slicing.** A **wide refactor** is one mechanical change (rename a column, retype a shared symbol) whose **blast radius** fans across the whole codebase, so a single edit breaks thousands of call sites at once and no vertical slice can land green. Don't force it into a tracer bullet; sequence it as **expand–contract**. First expand: add the new form beside the old so nothing breaks. Then migrate the call sites over in batches sized by blast radius (per package, per directory), each batch its own ticket blocked by the expand, keeping CI green batch to batch because the old form still exists. Finally contract: delete the old form once no caller remains, in a ticket blocked by every migrate batch. When even the batches can't stay green alone, keep the sequence but let them share an integration branch that all block a final integrate-and-verify ticket; green is promised only there.

### 4. Review the breakdown

Draft the breakdown as a numbered list. For each ticket, note:

- **Title**: short descriptive name
- **Blocked by**: which other tickets (if any) must complete first
- **What it delivers**: the end-to-end behaviour this ticket makes work
- **Alignment**: remove, build or change, and the simulated lines and fields it produces
- **Dev part**: its category and paths, or its override

Check it yourself, or with the owner if `.agents/policy/human-in-the-loop.policy.md` puts a checkpoint here:

- Does the granularity feel right? (too coarse / too fine)
- Are the blocking edges correct: does each ticket only depend on tickets that genuinely gate it?
- Should any tickets be merged or split further?
- Is each dev part a problem worth reasoning about, and does it reach the Floor? Do its paths stay clear of everything the agent must write in that ticket?

Revise until every answer holds, then publish.

### 5. Publish the tickets to the configured tracker

Publish the approved tickets. **How** depends on the tracker configured in `.agents/context.md`; the tickets are the same either way, only the shape of the blocking edges changes:

- **Local files** → write one file per ticket at `.agents/tickets/<NNNN>-<slug>.ticket.md`, continuing the repository-wide ticket sequence in dependency order (blockers first). Each file's "Blocked by" lists the numbers/titles it depends on. Use the per-ticket file template below: one ticket per file, never a single combined file.
- **A real issue tracker (GitHub, Linear, …)** → publish one issue per ticket in dependency order (blockers first) so each ticket's blocking edges can reference real identifiers. Use the platform's native blocking relationship where it has one; otherwise set each ticket's "Blocked by" to the blocking issues. If the source was an existing issue, make each ticket its sub-issue (tracker doc's operation). Apply the `ready-for-agent` triage label unless instructed otherwise; the tickets are agent-grabbable by construction.

Work the **frontier**: any ticket whose blockers are all done. For a purely linear chain that means top to bottom.

Do NOT close or modify any parent issue.

<local-ticket-template>

# <NNNN>: <Ticket title>

**Spec:** the path of the spec this ticket implements, or "None".

**What to build:** the end-to-end behaviour this ticket makes work, from the user's perspective, not a layer-by-layer implementation list.

**Blocked by:** the numbers/titles of the tickets that gate this one, or "None (can start immediately)".

**Status:** ready-for-agent

**Alignment:** remove, build or change. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] Acceptance criterion 1, quoting the exact simulated line or card field
- [ ] Acceptance criterion 2

## Dev part

**Category:** `core-domain`, `hard-part` or `new-stack`, and why this problem.

**Paths:** the files the owner writes, tests included.

**Boundary:** the signature or interface the rest of the code calls, which the agent writes outside those paths.

**Goal:** the problem the owner solves.

- [ ] Dev part acceptance criterion

**Override:** the reason and the date, only when there is one. It stands in place of the fields above when no part qualifies, and beside them when the owner does not write a part that qualifies.

</local-ticket-template>

<issue-template>

## Parent

A reference to the parent issue on the tracker (if the source was an existing issue, otherwise omit this section).

## What to build

The end-to-end behaviour this ticket makes work, from the user's perspective, not layer-by-layer implementation.

## Alignment

Remove, build or change, with the canonical rule of ADR 0005: what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

## Acceptance criteria

- [ ] Criterion 1, quoting the exact simulated line or card field
- [ ] Criterion 2

## Dev part

The category, paths, boundary, goal and acceptance criteria of the part the owner writes, and the override with its reason and date when there is one, as in the local template.

## Blocked by

- A reference to each blocking ticket, or "None (can start immediately)". Omit this section when blockers were set as native edges.

</issue-template>

In either form, avoid specific file paths or code snippets: they go stale fast. The dev part's paths and boundary are an exception, because they are the decision itself. Another exception: if a prototype produced a snippet that encodes a decision more precisely than prose can (state machine, reducer, schema, type shape), inline it and note briefly that it came from a prototype. Trim to the decision-rich parts, not a working demo, just the important bits.
