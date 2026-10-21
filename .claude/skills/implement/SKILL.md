---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
metadata:
  source: https://github.com/mattpocock/skills/tree/5beb35b/skills/engineering/implement
  adapted: Adjusted to this repository's layout; never commits.
---

Implement the work described by the user in the spec or tickets.

If the user passes a ticket reference, fetch it from the issue tracker and state its title before starting. If the reference is ambiguous, ask.

Call the Skill tool with "tdd" where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, call the Skill tool with "code-review" to review the work.

Do not commit. Report the work as ready for review; committing needs the owner's explicit authorization (see `AGENTS.md`).
