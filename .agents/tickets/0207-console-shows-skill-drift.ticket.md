# 0207: the console shows the team's skill drift

**Spec:** .agents/specs/0033-safe-seat.spec.md

**What to build:** The console's FinOps area shows the month's skill drift for the team: the applicable tickets (met plus overridden), the met tickets, the overridden tickets and the share overridden, with the drift colors. Tickets that are `not-applicable` or `unknown`, and tickets `ex` did not run, stay out of the share. The reading never shows a person or a ticket's share. The console computes the reading from the tickets that 0206 seeds, so console-stubs holds no skill drift of its own; the reading is labeled illustrative, and the copy goes into the five catalogs.

**Blocked by:** 0206.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] The FinOps area shows skill drift for the current month from the stub, with the drift colors and the illustrative label.
- [ ] The reading is computed from the seeded tickets' states. console-stubs serves no skill drift figure.
- [ ] A month with no applicable ticket shows the reading as unknown, never as zero drift.
- [ ] The copy exists in the five locales with the same leaf keys.
- [ ] No horizontal overflow at 375px.

## Dev part

**Category:** core-domain. The team's skill drift.

**Paths:**
- `apps/console/src/app/lib/views/finops/skill-drift.compute.ts`
- `apps/console/src/app/lib/views/finops/skill-drift.test.ts`

**Boundary:** The FinOps view calls `skillDrift(tickets: SafeSeatTicket[], period: string): SkillDrift | null`, where `period` is a month in the budget's `YYYY-MM` form. `SkillDrift` holds `applicable`, `met`, `overridden` and `overriddenShare`. The agent declares both types in `skill-drift.types.ts` and writes the view.

**Goal:** Count one period of tickets:
- only `met` and `overridden` tickets are applicable;
- `not-applicable`, `unknown` and `needs-override` tickets stay out, and so does a ticket with no Safe Seat record, which `ex` did not run;
- no applicable ticket gives `null`, never zero.

- [ ] Node tests cover a period mixing all five states, a period with none applicable, a ticket with no Safe Seat record, and tickets from another period left out.

Until the owner writes these files, the build fails on `skillDrift` alone.

## Verification

- `pnpm --filter console test`
- `pnpm --filter console typecheck`
- `pnpm turbo run test --filter console-stubs`
- `pnpm turbo run check`
- `platform/cli/drift check localization`
