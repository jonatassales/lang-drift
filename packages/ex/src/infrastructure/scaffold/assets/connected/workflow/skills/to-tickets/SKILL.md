---
name: to-tickets
description: Break a plan, spec, or the current conversation into draft tracer-bullet tickets in the configured tracker, each declaring its blocking edges; record readiness only when approval exists.
disable-model-invocation: true
---

# To Tickets

Break a plan, spec, or conversation into a set of **tickets**: tracer-bullet vertical slices, each declaring the tickets that **block** it.

Read `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md` for the configured tracker, artifact paths and states. If it is missing, report that configuration gap; do not infer a remote tracker.

## Process

### 1. Gather context

Work from whatever is already in the conversation context. Read the supplied spec or the matching existing spec under `.exphere/artifacts/specs/` when present and follow `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md`. If the user passes a reference (a spec path, an issue number or URL) as an argument, fetch it and read its full body and comments.

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

Give each ticket its **blocking edges**: the other tickets that must complete before it can start. An approved ticket with no blockers is dependency-ready; starting implementation still requires authorization for that scope.

**Wide refactors are the exception to vertical slicing.** A **wide refactor** is one mechanical change (rename a column, retype a shared symbol) whose **blast radius** fans across the whole codebase, so a single edit breaks thousands of call sites at once and no vertical slice can land green. Don't force it into a tracer bullet; sequence it as **expand–contract**. First expand: add the new form beside the old so nothing breaks. Then migrate the call sites over in batches sized by blast radius (per package, per directory), each batch its own ticket blocked by the expand, keeping CI green batch to batch because the old form still exists. Finally contract: delete the old form once no caller remains, in a ticket blocked by every migrate batch. When even the batches can't stay green alone, keep the sequence but let them share an integration branch that all block a final integrate-and-verify ticket; green is promised only there.

### 4. Prepare the breakdown for review

Present the proposed breakdown as a numbered list. For each ticket, show:

- **Title**: short descriptive name
- **Blocked by**: which other tickets (if any) must complete first
- **What it delivers**: the end-to-end behavior this ticket makes work

Ask the user:

- Does the granularity feel right? (too coarse / too fine)
- Are the blocking edges correct: does each ticket only depend on tickets that genuinely gate it?
- Should any tickets be merged or split further?

Save the proposed breakdown as local draft tickets so the user can review a concrete result. Reuse explicit approval of this same breakdown already present in the session. Otherwise ask for the unresolved decisions and retain draft status until approval is recorded; do not infer approval from silence or document creation.

### 5. Save tickets and record their actual state

Follow `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md`:

- **Local files** → write one en-US file per ticket in `.exphere/artifacts/tickets/`, using `NNNN-TICKET-<lowercase-english-slug>.ticket.md`. Continue the global TICKET sequence with four-digit numbers in dependency order (blockers first); do not restart numbering for each initiative. Preserve existing IDs and ticket edits. Each file's "Blocked by" lists the numbers/titles it depends on. Use the per-ticket file template below: one ticket per file, never a single combined file.
- **A remote issue tracker**, only when explicitly configured and publication authorized → publish one issue per ticket in dependency order. Use native blocking relationships where supported, or list the blocking identifiers. If the source was an existing issue, follow the tracker configuration for parent relationships. Do not apply readiness labels automatically.

New tickets remain `draft` until approval of their scope and breakdown is recorded with a reference to actual authorization. Existing authorization in the session counts. Use the configured readiness state only after the required approvals and prerequisites hold; this skill does not grant its own approval or start implementation. A scope change requires reconsidering any earlier approval.

For later authorized implementation, the **frontier** contains approved tickets whose blockers are done. For a purely linear chain that means top to bottom.

Do NOT close or modify any parent issue.

<local-ticket-template>

# <NN>: <Ticket title>

**What to build:** the end-to-end behavior this ticket makes work, from the user's perspective, not a layer-by-layer implementation list.

**Blocked by:** the numbers/titles of the tickets that gate this one, or "None".

**Status:** draft
**Approval:** Not recorded. Replace only with a reference to actual authorization covering this scope.

- [ ] Acceptance criterion 1
- [ ] Acceptance criterion 2

</local-ticket-template>

<issue-template>

## Parent

A reference to the parent issue on the tracker (if the source was an existing issue, otherwise omit this section).

## What to build

The end-to-end behavior this ticket makes work, from the user's perspective, not layer-by-layer implementation.

## Acceptance criteria

- [ ] Criterion 1
- [ ] Criterion 2

## Blocked by

- A reference to each blocking ticket, or "None". Omit this section when blockers were set as native edges.

</issue-template>

In either form, avoid specific file paths or code snippets: they go stale fast. Exception: if a prototype produced a snippet that encodes a decision more precisely than prose can (state machine, reducer, schema, type shape), inline it and note briefly that it came from a prototype. Trim to the decision-rich parts, not a working demo, just the important bits.
