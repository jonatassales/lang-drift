# Human in the loop

Where the owner takes part in this repository's spec-driven workflow. It varies between projects and teams, so the workflow skills read it from here instead of assuming it.

## Checkpoints

The owner takes part at three points only:

1. **Grilling.** It ends when the owner confirms the shared understanding.
2. **Dev part.** Every ticket reserves a part the owner writes by hand, as [`safe-seat.policy.md`](./safe-seat.policy.md) sets out. The owner writes it, or writes an override with its reason and date in the ticket.
3. **Pull request review.** The owner reviews the pull request and decides whether to merge it.

## Automatic stages

Once grilling ends, the agent runs these without stopping for approval, except at each ticket's dev part:

1. `to-spec`, with the test seams recorded in the spec.
2. `to-tickets`, with the breakdown reviewed by the agent itself and a dev part in every ticket.
3. `implement`, using `tdd`, until only the ticket's dev part remains. The agent stops there and waits for the owner.
4. `code-review`, once the owner has written the dev part.
5. Delivery. The agent commits each ticket on a branch, never on `main`, with its `Ticket: <NNNN>` trailer. It then pushes the branch and opens a pull request into `main` that names the spec and tickets.

## Gates

Only guardrails block the automatic stages:

- the effort hook before `implement`;
- the Git hooks;
- CI.

An agent never bypasses one: no `--no-verify` and no skipped checks. When a guardrail blocks, the agent fixes the cause. If it cannot, it stops and reports why.

## Always the owner's

- Merging and releasing.
- Invoking `grill-me` and `grill-with-docs`.
