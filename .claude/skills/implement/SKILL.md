---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
metadata:
  source: https://github.com/mattpocock/skills/tree/5beb35b/skills/engineering/implement
  adapted: Adjusted to this repository's layout; delivery follows the workflow policy; every ticket's Safe Seat dev part is left to the owner.
---

Implement the work described by the user in the spec or tickets.

If a ticket reference is passed, fetch it from the issue tracker and state its title before starting. Without one, take the spec's tickets in the order their blocking edges allow.

Every ticket reserves a dev part for the owner, as the Dev part section of `.agents/policy/safe-seat.policy.md` sets out. Read the ticket's `Dev part` section before starting:

- Never create, edit or delete anything in its paths, and that includes its tests: the owner writes the solution and the tests. Write only the boundary the section names, outside those paths.
- When only the dev part remains, stop. Typechecking and the tests then fail on the dev part's missing symbols, as expected, and on nothing else. Report its paths, the boundary you wrote, its goal and acceptance criteria, and those expected failures, then wait for the owner.
- An override in place of the fields leaves no dev part. An override beside the fields lets you write its paths.
- A ticket from 0199 on with neither a `Dev part` section nor an override stops you before you start: ask the owner for one. Tickets up to 0198 predate the policy and run as before.

Call the Skill tool with "tdd" where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once the owner has written the dev part, call the Skill tool with "code-review" to review the work against all rules in `.claude/rules/*.rule.md` and the originating spec. Any findings must be addressed and resolved before delivery. A finding inside the dev part's paths is the owner's to resolve: report it and do not fix it.

Then deliver the work as `.agents/policy/human-in-the-loop.policy.md` says: it sets whether you commit, push and open the pull request yourself, and which guardrails gate you. Without that file, do not commit; report the work as ready for review.
