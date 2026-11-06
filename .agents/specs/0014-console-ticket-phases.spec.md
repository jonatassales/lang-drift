# Console ticket phases: Tickets, Analysis, Plans and Workflows

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The website's hero shows what Microruntimes does with one ticket as four panels beside a simulated `ex` session: the ticket and its linked artifacts, the analysis that reads them and cites past episodes, the plan that sizes the components inside an envelope, and the workflow that spends against the team's budget up to a pull request. The console, where a team would follow that work, shows none of it. Its Microruntimes area lists runtimes, executions and history, and extub serves no ticket at all.

## Solution

The project's sidebar group gains four entries, Tickets, Analysis, Plans and Workflows. Each opens an infinite list of cards built on the website's panels, one card per ticket that has reached that phase: Tickets lists every ticket, Analysis every ticket with an analysis, Plans every ticket with a plan, Workflows every ticket with a workflow. A phase is a status, not a disclosure: each card shows the panel's header line over the panel's full body, with no collapse. A ticket in the middle of a phase shows a dash where a value is not known yet, never a zero. Every card opens the ticket's detail page, which stacks the four widgets and says which phases are not reached yet. extub serves the tickets through two new routes, with cursor paging, from a deterministic stub seed.

## User Stories

1. As a team member, I want Tickets, Analysis, Plans and Workflows in the project's sidebar, so that I can follow the work by phase.
2. As a team member, I want each listing to show one card per ticket that has reached that phase, so that I see every analysis, every plan and every workflow in one place.
3. As a team member, I want each card to show the same header line and body as the website's panel, so that the console and the site describe a ticket the same way.
4. As a team member, I want a card whose phase is in progress to show a dash where a value is unknown, so that nothing reads as zero before it is measured.
5. As a team member, I want the listings to load more cards as I scroll, so that a long list never loads at once.
6. As a team member, I want to open a ticket's detail page from any card, so that I can read the whole ticket in one place.
7. As a team member, I want the detail page to say which phases are not reached yet, so that a missing widget never looks like an error.
8. As a team member, I want the Analysis card to show the documents read, the episodes kept and dropped against the threshold, and the suggested model and effort, so that I see why the recommendation was made.
9. As a team member, I want the Plan card to show the proposal, the envelope, the components with their marks and shares, and the report, so that I see what will run and inside which limits.
10. As a team member, I want the Workflow card to show the human stages with their progress, the estimated spend against the budget, and the result with reserved, observed and released amounts and the pull request, so that I see what the work cost and where it stands.
11. As a team member, I want every recommendation to say the team authorizes it, so that a suggestion never reads as a decision.
12. As a team member reading pt-BR, ja or zh-Hant, I want the listings and the detail in my language, so that nothing depends on English.
13. As a team member on a narrow screen, I want the cards to stack without horizontal overflow at 375px.
14. As a team member, I want the listings to say extub is unavailable when it is down, as every other screen does.
15. As a frontend developer, I want the share bar shared by the website and the console from `@langdrift/web`, so that both draw it the same way.
16. As a frontend developer, I want extub's ticket paging tested for duplicates, cursor end and cumulative phase filtering, so that the infinite list is reliable.

## Implementation Decisions

- **Phases.** `ticket`, `analysis`, `plan`, `workflow`, in that order. A ticket's `phase` is the last phase it has reached; a section of the ticket exists once its phase is reached and may hold unknown fields while the phase is current. Tests is not a phase.
- **extub data.** A deterministic generator in the domain builds about sixty tickets for the seeded project from a handful of use cases, spread across the four phases, every figure inside its envelope, every response labelled `source: stub`. The early-access ticket carries the website's figures. A project created at runtime has no tickets.
- **extub routes.** `GET console/products/{id}/tickets?reached={phase}&cursor=` returns `{source, items, nextCursor}`, newest first, twenty per page, the cursor a decimal position as in the activity routes; `reached` keeps the tickets whose phase is at or after it and defaults to `ticket`; an unknown phase or a malformed cursor is invalid input. `GET console/tickets/{id}` returns one ticket; an unknown id is not found. This supersedes, for tickets, the line of spec 0013 that kept new stub routes out of scope.
- **Console routes.** `products/[id]/tickets`, `products/[id]/analysis`, `products/[id]/plans`, `products/[id]/workflows` and `products/[id]/tickets/[ticketId]`, each with page, loading, error and template like the executions route. Server pages read through `readConsole` and hand the promise to a client view.
- **Listing.** One client view, `TicketListing`, takes the phase, pages with `usePagedResource`, renders with `VirtualList` and keeps the load-more fallback, exactly as the executions overview. The phase picks the widget.
- **Widgets.** `TicketWidget`, `AnalysisWidget`, `PlanWidget` and `WorkflowWidget`: a card with an eyebrow, the panel's summary line, the ticket's current phase as a status label, and the panel's body. The workflow widget is read-only.
- **Detail.** `TicketDetail` reads the ticket with `useResource`, shows the ticket head and the four widgets in order; a phase not reached yet is one muted line.
- **Shared pieces.** `ShareBar` moves to `@langdrift/web/ui/share-bar`; the website imports it from there. Technology icons, `Eyebrow` and `DetailList` already live in `@langdrift/web`. Rows for components, episodes and documents are console-local.
- **Copy.** A `tickets` namespace in the console's four catalogs with the same keys, ICU for counts. Model names, `runtime.md`, file names, paths and criteria ids stay untranslated. Figures are formatted with `Intl` for the locale: millicores, MiB, dollars from cents and token counts.

## Testing Decisions

A good test observes what a user sees or what a URL returns, never internal structure.

1. **extub** (Go, beside the activity test): the Tickets route pages twenty at a time with no duplicate id across pages and ends with a null cursor; `reached=analysis` returns only tickets at or past analysis and every one of them also appears under `reached=ticket`; an unknown phase and a malformed cursor are 400; a ticket's detail is 200 and an unknown id 404; a created project has no tickets.
2. **Console logic** (Node test runner, as the console's other tests): the phase-reached predicate, the summary line of each phase including the dash for an unknown value, money and token formatting, and share computation.
3. **Repository checks**: `pnpm turbo run check`, `pnpm turbo run typecheck test`, `drift check localization`, `drift check dead-code`.

Manual verification with `pnpm console`: the four listings scroll and load more without duplicates, each card matches its website panel, cards open the detail, `en` and `pt-br`, light and dark, 375px, and `extub_unavailable` with extub stopped.

## Out of Scope

- Creating, editing or advancing a ticket from the console; the widgets are read-only.
- A Tests listing.
- Real ticket data from `.langdrift/` or the `ex` CLI.
- Changes to the website's panels beyond importing the shared share bar.

## Further Notes

The panels and their figures come from the website's `microruntimes-content.ts` and `microruntimes-workflow.ts`. Money is in cents and token counts in thousands, as there.
