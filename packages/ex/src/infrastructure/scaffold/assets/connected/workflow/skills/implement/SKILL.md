---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
---

Implement the work described by the user in the spec or tickets.

If the user passes a ticket reference, fetch it from the issue tracker and state its title before starting. If the reference is ambiguous, ask.

Use the `tdd` skill where possible, at pre-agreed seams. Invoke it through the Skill tool when available; otherwise read and follow `.exphere/workflow/skills/tdd/SKILL.md`.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, use `code-review` to review the work, including the scoped changes that are still uncommitted. Invoke it through the Skill tool when available; otherwise read and follow `.exphere/workflow/skills/code-review/SKILL.md`.

Commit only when the user has authorized committing this scope; existing authorization in the session counts. Otherwise leave the reviewed changes uncommitted and report their status. Invoking implementation or completing a review does not itself authorize a commit, push, pull request, or deployment.
