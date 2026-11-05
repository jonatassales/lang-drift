# Safe Seat

How LangDrift practices Safe Seat on its own work, before the `dft` runtime measures it. Spec 0033 defines the capability. This policy applies it to this repository's spec-driven workflow, beside [`human-in-the-loop.policy.md`](./human-in-the-loop.policy.md), and the skills read it instead of assuming.

## Floor

A ticket's dev part is at least 5% of the lines the ticket's diff adds. The 5% is LangDrift's principle, not a threshold derived from a study (research note `safe-seat.research.md`). It confirms that the dev part was really written; it is not a dose of practice. The practice comes from the problem the dev part holds.

## Dev part

Every ticket from 0199 on reserves one dev part: a problem with its own goal and acceptance criteria, which the owner solves by hand, writing both the solution and its tests. Tickets up to 0198 predate this policy and have none.

- The ticket's `Dev part` section names its category, its paths, its boundary, its goal and its acceptance criteria.
- Only the owner writes in its paths. The agent writes the boundary the rest of the code calls, a signature or an interface, outside those paths, and nothing inside them.
- `implement` stops when only the dev part remains and waits for the owner. Until the owner writes it, the build fails on the dev part's missing symbols and on nothing else.
- `code-review` reviews the dev part once the owner has written it. A finding inside its paths is the owner's to fix.

## Categories

A dev part is one of:

- `core-domain`: domain logic the owner must be able to explain and debug during an incident;
- `hard-part`: the place where the ticket's difficulty sits;
- `new-stack`: first contact with a tool, library or stack.

It is never boilerplate, configuration, generated code, a rename or a migration run by a tool.

## Decisions

`to-tickets` proposes each ticket's dev part. The proposal stands as the decision unless the owner changes it, replaces it or overrides it, at the latest when `implement` stops at it. The ticket's `Dev part` section, as committed, records the decision.

## Override

An override is never silent: it records its reason and its date in the ticket's `Dev part` section. It takes one of two forms:

- **No part qualifies.** The override stands in place of the other fields, and the ticket has no dev part.
- **The owner does not write a part that qualifies.** The override sits beside the fields, and the agent may write its paths.

A ticket from 0199 on with neither a dev part nor an override stops `implement` until the owner gives one.

## Measurement

Nothing measures the share yet. Commits carry no AI signature, so the repository cannot tell who wrote a line. The ticket's dev part and the review hold the practice until the `dft` runtime of spec 0033 measures it. From then on, this policy follows the runtime's floor, attribution and gates.
