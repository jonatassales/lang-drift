# 0124: The ticket page reads as one document

**Spec:** .agents/workflow/specs/0019-console-ticket-widgets.spec.md

**What to build:** the ticket page loses its cards. A header carries the ticket's number, its title as the page heading and its slug. A strip of labels over values shows the status, the phase as "Plan · 3 of 4" and when the ticket was created. Below come the linked artifacts and one section per phase under a full-width hairline, each with a heading and its summary line and the phase's body as it is today; a phase not reached yet stays one muted line. An entry in the Analysis, Plans or Workflows listing opens the page at that phase's section, below the sticky header.

**Blocked by:** 0123, Ticket entries become a card on mobile and a row on desktop.

**Status:** ready-for-agent

- [ ] `summaryPhases`, `phaseStep(ticket)` and `ticketPath(ticket, listing)`, test first: the phases without `ticket`; 1 of 4 and 4 of 4; no fragment in Tickets and the phase key in the other three listings.
- [ ] The widget frame becomes a section with an `h2` and the summary line, offset below the sticky header; the Analysis, Plan and Workflow widgets wrap their unchanged bodies in it.
- [ ] The Ticket widget becomes the linked artifacts section, without the number, slug, status and title the header now carries.
- [ ] The ticket page renders the header, the strip, the artifacts and the phases with no card; status appears once, in the strip.
- [ ] Once the ticket has rendered, the page scrolls to the URL's fragment; the listings' entries link to it.
- [ ] A ticket detail skeleton replaces the runtime detail skeleton in the ticket page's loading state and Suspense fallback.
- [ ] The breadcrumb builds its ticket label with `ticketSummary`; the console README describes the listings and the ticket page as they now are.
- [ ] The strip's labels and phase progress in the four console catalogs, with the same keys.

### Verification

- `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`, and the four catalogs' leaf keys equal.
- In the browser: no card on the ticket page; a Plans entry lands on its section below the header, on navigation and on reload; polling never moves the page; the breadcrumb reads "Project / Tickets / number · slug"; `en` and `pt-br`, light and dark, 375px; Executions and History unchanged.
