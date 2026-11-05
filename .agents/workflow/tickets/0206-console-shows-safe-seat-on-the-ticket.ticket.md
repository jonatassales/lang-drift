# 0206: the console shows Safe Seat on each ticket

**Spec:** .agents/workflow/specs/0033-safe-seat.spec.md

**What to build:** console-stubs seeds LangDrift Studio's tickets with Safe Seat data, labeled illustrative like the rest of the stub:
- a person-driven `dev-ticket` stage in each workflow;
- a reserved part on each plan;
- a Safe Seat state per ticket, covering `met`, `needs-override`, `overridden` with a reason, `not-applicable` and `unknown`.

In the console:
- the Workflow card marks a person-driven stage apart from a human-approved one;
- the Plan card shows the reserved part (paths, goal, category) and who chose it;
- each ticket shows its Safe Seat state, with the override reason when there is one.

No type and no view carries a share. The copy goes into the five console catalogs with key parity, and the name Safe Seat stays untranslated.

> **Amended 2026-11-05.** The console serves four locales (`en`, `pt-BR`, `zh-Hant`, `ja`): the owner closed 0308 (`vi` for the console) as won't fix, so the copy goes into those four catalogs with key parity.

**Blocked by:** 0199.

**Status:** ready-for-agent

- [x] console-stubs serves the Safe Seat fields for LangDrift Studio, and its catalog test checks that the seeded tickets cover every state.
- [x] `TicketStage` carries `drivenBy`. The ticket carries its Safe Seat state and override reason. The plan carries the reserved part. No field holds a share.
- [x] The Workflow card shows a distinct marker for the person-driven stage, tested beside the ticket phase tests for the person-driven, human-approved and agent cases. The Plan card shows the reserved part and `chosenBy`. The ticket shows its state, labeled as illustrative stub data.
- [x] The copy exists in `en`, `pt-BR`, `zh-Hant` and `ja` (the console's four locales, amended 2026-11-05), with the same leaf keys, and "Safe Seat" is never translated or abbreviated.
- [x] No horizontal overflow at 375px on the ticket views.

## Dev part

**Category:** core-domain. How a ticket's Safe Seat record reads in the console.

**Paths:**
- `apps/console/src/app/lib/views/tickets/safe-seat-view.map.ts`
- `apps/console/src/app/lib/views/tickets/safe-seat-view.test.ts`

**Boundary:** The ticket views call `safeSeatView(ticket: Ticket): SafeSeatView`. `SafeSeatView` carries:
- the state's message key and its status tone;
- the override reason, when there is one;
- the reserved part's paths, goal and category, and who chose it;
- whether the record is illustrative stub data.

The agent declares `SafeSeatView` in `safe-seat-view.type.ts` and writes the cards.

**Goal:** Map a stub ticket's Safe Seat fields to what the cards show:
- every state value has its key and tone;
- a missing record reads as `unknown`, never as `met`;
- an override reason appears only with `overridden`;
- a missing reserved part reads as absent, not as an empty list.

- [x] Node tests cover each of the five states, a ticket with no Safe Seat record, an override with and without a reason, and a ticket with and without a reserved part.

Until the owner writes these files, the build fails on `safeSeatView` alone. (Superseded by the override below: the agent wrote them.)

**Override:** the owner will not write this part ("I won't write anything", 2026-11-05); the agent writes its paths.

## Verification

- `pnpm --filter console test`
- `pnpm --filter console typecheck`
- `pnpm turbo run test --filter console-stubs`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
