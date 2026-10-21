---
name: to-spec
description: "Turn the current conversation into a draft spec in the configured local tracker: no interview, just synthesis of what you've already discussed."
disable-model-invocation: true
---

This skill takes the current conversation context and codebase understanding and produces a spec. Do NOT interview the user; just synthesize what you already know.

Read `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md` for the configured tracker, artifact paths and states. If it is missing, report that configuration gap; do not infer a remote tracker.

## Process

1. Explore the repo to understand the current state of the codebase, if you haven't already. Use the project's domain glossary vocabulary throughout the spec, and respect any ADRs in the area you're touching.

2. Sketch out the seams at which you're going to test the feature. Existing seams should be preferred to new ones. Use the highest seam possible. If new seams are needed, propose them at the highest point you can. The fewer seams across the codebase, the better - the ideal number is one.

Reuse explicit agreement on these seams already present in the session. If no agreement exists, record them as proposed in the draft and ask for confirmation before implementation; drafting need not wait for that confirmation.

3. Write the spec in en-US under `.exphere/artifacts/specs/`, using `NNNN-SPEC-<lowercase-english-slug>.spec.md`. Allocate the next available four-digit SPEC number according to the configured tracker; preserve an existing spec and its stable ID when updating it. Follow `.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md`. When `.agents/templates/spec.md` exists, use it as the canonical format and preserve its fields, identifiers and state conventions; use the guidance below to fill its corresponding sections without replacing its structure. Use the inline template only when no local spec template exists. Use the existing initiative ID when present. New content is a draft by default. Record approval only when explicit user authorization covers this spec and scope; cite that authorization in the record, including prior authorization in the session when applicable. Do not invent an approval or grant `ready-for-agent` merely because a document was generated. Approval of a previous version does not approve a changed scope. External publication follows separately authorized tracker operations.

<spec-template>

**Status:** draft
**Approval:** Not recorded. Replace only with a reference to actual authorization covering this scope.

## Problem Statement

The problem that the user is facing, from the user's perspective.

## Solution

The solution to the problem, from the user's perspective.

## User Stories

A numbered list of user stories covering the agreed scope. Scale the number and detail to the initiative profile: a small correction may need one story; a larger feature may need several. Each user story should be in the format of:

1. As an <actor>, I want a <feature>, so that <benefit>

<user-story-example>
1. As a mobile bank customer, I want to see balance on my accounts, so that I can make better informed decisions about my spending
</user-story-example>

Cover the relevant users, critical behavior and material edge cases without inventing scope to make the document longer. Follow the configured workflow profile where one exists; proposed profiles remain proposals until adopted for this initiative.

## Implementation Decisions

A list of implementation decisions that were made. This can include:

- The modules that will be built/modified
- The interfaces of those modules that will be modified
- Technical clarifications from the developer
- Architectural decisions
- Schema changes
- API contracts
- Specific interactions

Do NOT include specific file paths or code snippets. They may end up being outdated very quickly.

Exception: if a prototype produced a snippet that encodes a decision more precisely than prose can (state machine, reducer, schema, type shape), inline it within the relevant decision and note briefly that it came from a prototype. Trim to the decision-rich parts, not a working demo, just the important bits.

## Testing Decisions

A list of testing decisions that were made. Include:

- A description of what makes a good test (only test external behavior, not implementation details)
- Which modules will be tested
- Prior art for the tests (i.e. similar types of tests in the codebase)

## Out of Scope

A description of the things that are out of scope for this spec.

## Further Notes

Any further notes about the feature.

</spec-template>
