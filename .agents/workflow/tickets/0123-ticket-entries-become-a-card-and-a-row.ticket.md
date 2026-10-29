# 0123: Ticket entries become a card on mobile and a row on desktop

**Spec:** .agents/workflow/specs/0019-console-ticket-widgets.spec.md

**What to build:** every entry in the four ticket listings becomes one small widget. Below `lg` it is a card with the number, the date the ticket was created, the status, the title, the slug and a footer with the entry's figure. From `lg` it is a compact bordered row with the number, the title over the slug, the status, the figure and the created date. The figure is the panel's one-line summary: in Tickets, the ticket's current phase; in Analysis, Plans and Workflows, that listing's phase. The whole widget opens the ticket page. The panels' full bodies leave the listings, and the loading placeholders take the shape of the new entries.

**Blocked by:** 0122, The ticket listings scroll with the page.

**Status:** ready-for-agent

- [x] `SummaryPhase` and `listedPhase(ticket, listing)`, test first: none for a Ticket-phase ticket in Tickets, the current phase for the others, and the listing's phase in Analysis, Plans and Workflows.
- [x] One `PhaseSummary` renders the analysis, plan and workflow summary lines, moved out of their widgets; the widgets' headers render it; the plan's authorization label is shared through the ticket format hook.
- [x] `TicketCreated` renders `createdAt` as an absolute date in a `time` element.
- [x] `TicketCard` below `lg` and `TicketRow` from `lg`, each with one link covering the widget, named by the number and title, without prefetch; the row's figure moves under the status from `xl`.
- [x] The listings render the two widgets per ticket instead of the phase widgets; the `listed` prop, the frame's link line and the Ticket widget's linked number are removed, and so is the `openTicket` copy key.
- [x] Entry and list skeletons replace the activity skeletons in the four listings' loading states, Suspense fallbacks and loading-more row.
- [x] The four listing intros describe the entries in the four console catalogs, with the same keys.

### Verification

- `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`, and the four catalogs' leaf keys equal.
- In the browser: cards at 375px with no horizontal overflow; cards at 768, rows at 1024 and 1280; each listing's figure follows the rule above; click and Enter open the ticket; one focus ring around the whole widget; no hydration warning in `en` and `pt-br`, light and dark.
