# 0087: The console lists the project's tickets and opens one

**Spec:** .agents/workflow/specs/0014-console-ticket-phases.spec.md

**What to build:** the project's sidebar group gains Tickets, Analysis, Plans and Workflows. Tickets opens an infinite list of ticket cards, each built on the website's Ticket panel and opening the ticket's detail page. The other three entries arrive with tickets 0088 to 0090.

**Blocked by:** 0085.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Console types mirror extub's ticket.
- [ ] `products/[id]/tickets` and `products/[id]/tickets/[ticketId]` with page, loading, error and template.
- [ ] `TicketListing` pages with `usePagedResource` and `VirtualList` and keeps the load-more fallback.
- [ ] `TicketWidget`: eyebrow, summary line, current phase as a status label, number, slug, status, title and linked artifacts.
- [ ] `TicketDetail`: the ticket head, the reached widgets in order and one muted line per phase not reached.
- [ ] Sidebar entries and the `tickets` catalog namespace in the four console locales.
- [ ] Node tests for the phase-reached predicate and the Ticket summary line.

### Verification

- `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`; the listing and the detail in the browser.
