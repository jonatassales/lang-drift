# 0199: LangDrift practices Safe Seat in writing

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** This repository adopts Safe Seat before the runtime exists. A Safe Seat policy sits beside the human-in-the-loop policy and states LangDrift's own practice. Every ticket has a dev part of at least 5% of its diff, which the owner writes by hand. `to-tickets` proposes the dev part, and the owner decides it. An override is written in the ticket with its reason, and nothing measures the share yet. The human-in-the-loop policy names the dev part as the owner's third checkpoint, so the automatic stages stop for it and resume with `code-review` and delivery.

The workflow skills follow the policy:
- `to-tickets` writes a dev part section into every ticket, with its paths, the boundary the rest of the code calls, a goal, acceptance criteria and a category, or else an override with its reason;
- `implement` leaves the dev part's paths alone and stops when only they remain;
- `code-review` reviews the dev part once the owner has written it.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The Safe Seat policy exists beside the human-in-the-loop policy and covers:
  - the 5% floor;
  - one dev part per ticket;
  - the three categories and the excluded kinds;
  - who decides;
  - the override written in the ticket;
  - the absence of measurement until the runtime exists.
- [ ] The human-in-the-loop policy lists the dev part as a third owner checkpoint and points to the Safe Seat policy.
- [ ] `to-tickets` writes a dev part section into every ticket, with category, paths, boundary, goal and acceptance criteria, or else an override with its reason.
- [ ] `implement` reads the dev part section, never edits its paths, and stops with a report naming them when only they remain.
- [ ] `code-review` reviews the dev part after the owner has written it, and reports a dev part written by an agent as a finding.
- [ ] The spec-driven section of the context document names the Safe Seat policy beside the human-in-the-loop policy.
- [ ] The name is never shortened to its initials in any of these files.

## Dev part

**Category:** core-domain. The rule of the practice itself, in the owner's words.

**Paths:** `.agents/policies/safe-seat.policy.md`

**Boundary:** The skills and the human-in-the-loop policy link to this file by its path and name its sections: Floor, Dev part, Categories, Decisions, Override, Measurement. The agent writes those references and nothing inside the file.

**Goal:** State how LangDrift practices Safe Seat on its own work.

- [ ] Each of the six sections exists and agrees with spec 0033's dogfooding decisions.
- [ ] The policy says what happens when a ticket has no part worth writing by hand.

**Override:** on 2026-11-03 the owner asked the agent to write this dev part ("escreva vc"), so the agent wrote the policy text.

## Verification

- `platform/cli/drift quality check`
- `pnpm turbo run check`
