# 0202: ex stops at the dev ticket and the person closes it

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** A workflow stage gains `drivenBy`, either `agent` (the default) or `person`, and a workflow may declare at most one person stage. The sample workflow adds `dev-ticket`, labeled "Dev ticket", after `implement` and before `review`.

The person stage depends on the policy:
- **Safe Seat on:** the play loop reaches the person stage, sets the stage and the run to `waiting`, and prints the reserved part's paths, goal and acceptance criteria. It never calls the stage driver or runs skills there; the stage's hooks still run.
- **Safe Seat off:** the person stage is skipped.
- **Safe Seat on, no person stage:** the plan fails a check that names the missing stage.

Every agent call before the dev ticket receives the reserved paths through a `{{reserved}}` placeholder and a `LANGDRIFT_RESERVED_PATHS` variable. The driver's fixed prompt tells the agent to write only the boundary outside those paths.

`ex agent confirm <runtime-id> dev-ticket` closes the person stage only with an attestation, recorded with `--by` on the lap:
- `--without-agent` or `--with-agent`;
- or, at a terminal, an answer to the prompt asking "Written without an agent?".

`--confirm` on `ex agent run` never closes a person stage. A close without a terminal and without an attestation flag is refused with the reason. The review stage that follows runs `code-review` through the agent as usual.

**Blocked by:** 0201.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] With Safe Seat on, a run reaches the dev ticket, waits, prints the reserved part, and the scripted driver records no call for that stage.
- [ ] A scripted driver call before the dev ticket receives the reserved paths in its argument and its environment.
- [ ] `ex agent confirm … dev-ticket --without-agent` closes the stage, records the attestation, the `--by` and the time in the lap's private state, never in the rendered lap or the committed run record, and the run moves to `review`. `--with-agent` closes it too, recording that answer.
- [ ] The run's estimate and reservation count no paid call for the person stage, and the reservation stays held while the dev ticket waits.
- [ ] At a terminal, closing asks the attestation question. `--confirm` never closes the dev ticket. Without a terminal and without a flag, the close is refused with the reason.
- [ ] With Safe Seat off, the person stage is skipped. With Safe Seat on and no person stage, the plan fails its check with guidance. A workflow with two person stages is refused when recorded.
- [ ] The workflow mirror carries `drivenBy` for every stage.
- [ ] The session prints the website's dev-ticket run lines (spec 0043: `devTicketWaiting`, `devTicketClosed`, `safeSeatMet`, and the implement and code-review lines it moved), and the CLI's session catalog and its test mirror them.
- [ ] The `ex` README documents `drivenBy`, the dev ticket and the attestation flags.

## Dev part

**Category:** core-domain. When a person's attestation is valid.

**Paths:**
- `lib/cli/src/application/execution/close-dev-ticket.go`
- `lib/cli/src/application/execution/close-dev-ticket_test.go`

**Boundary:** The confirm use case calls `closeDevTicket(request devTicketClose) (attestation, error)`. The request carries:
- the attestation flag, if any;
- the answer to the terminal prompt, if a terminal asked it;
- whether the close came from `--confirm`;
- the stage state;
- the `--by`.

The confirm use case asks the prompt when a terminal is attached and no flag was given, and passes the answer. The agent declares `devTicketClose` and `attestation` in their own files and writes the callers.

**Goal:** Decide whether and how the dev ticket closes:
- `--confirm` never closes;
- a flag decides;
- otherwise, the terminal's answer decides;
- with neither, the close is refused with the reason;
- only a waiting stage closes.

- [ ] Table tests cover every combination of flag, terminal answer, `--confirm` and stage state.

Until the owner writes these files, the build fails on `closeDevTicket` alone.

## Verification

- `pnpm turbo run test typecheck --filter @langdrift/cli`
- `pnpm turbo run check`
