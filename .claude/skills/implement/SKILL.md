---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
metadata:
  source: https://github.com/mattpocock/skills/tree/5beb35b/skills/engineering/implement
  adapted: Adjusted to this repository's layout; delivery follows the workflow policy.
---

Implement the work described by the user in the spec or tickets.

If a ticket reference is passed, fetch it from the issue tracker and state its title before starting. Without one, take the spec's tickets in the order their blocking edges allow.

Call the Skill tool with "tdd" where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, call the Skill tool with "code-review" to review the work.

Then deliver the work as `.agents/workflow/human-in-the-loop.policy.md` says: it sets whether you commit, push and open the pull request yourself, and which guardrails gate you. Without that file, do not commit; report the work as ready for review.
