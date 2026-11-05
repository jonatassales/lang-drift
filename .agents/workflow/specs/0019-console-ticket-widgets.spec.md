# Console ticket widgets: a card on mobile, a row on desktop, a flat ticket page

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The console's four ticket listings, Tickets, Analysis, Plans and Workflows, draw every ticket inside one bordered box that scrolls on its own inside a page that also scrolls. Every entry is the website's full panel, between 350 and 640 pixels tall, so about one and a half tickets fit on a screen. The first entry is clipped at the top of the box. The number and slug are printed twice per entry. A Load more button sits under a box that also loads as it scrolls. The Analysis, Plans and Workflows listings never show a ticket's title. The ticket's own page puts each phase inside a card, so it reads as boxes inside boxes rather than as one ticket.

## Solution

Each listing becomes a page-scrolled list of one small widget per ticket. On a narrow screen the widget is a card: the number, the date the ticket was created and its status, the title, the slug, and a footer with the entry's figure, which wraps rather than being cut. From `lg` up it is a compact bordered row in the grammar of Vercel's deployment rows: number, title over slug, status over the figure, created date and a caret. The whole widget opens the ticket. The figure is the panel's existing one-line summary. In the Tickets listing it is the summary of the ticket's current phase; in the Analysis, Plans and Workflows listings it is that listing's phase, and the entry opens the ticket page at that phase's section. The panels' full bodies leave the listings and live on the ticket page only.

The ticket page becomes one document with no card. Its header carries the number and slug on one line over the title, which is the page heading. Below it comes a strip of labels over values (Status, Phase and Created), as Vercel's deployment page and Geist's Description do. Then come the linked artifacts as a definition list and one section per phase under a full-width hairline, each with its heading and summary line. Each phase's body stays as it is today. A phase not reached yet stays one muted line.

## User Stories

1. As a team member, I want a listing to scroll with the page instead of inside a box, so that nothing is clipped and there is one scrollbar.
2. As a team member, I want to see about ten tickets per screen on desktop, so that I can scan the project's work newest first.
3. As a team member on a phone, I want each ticket as its own card, so that every tap target is a whole ticket.
4. As a team member on desktop, I want each ticket as one compact row, so that numbers, titles, statuses and dates line up across tickets.
5. As a team member, I want every entry to show the ticket's number, title, slug and status once each, so that nothing reads twice.
6. As a team member, I want the Analysis, Plans and Workflows listings to show each ticket's title, so that I know which work a figure belongs to.
7. As a team member, I want each Tickets entry to show the latest figure of the ticket's current phase, so that I see where every ticket stands without opening it.
8. As a team member, I want each Analysis entry to show what the analysis settled last (documents read, episodes kept, or the suggested model and effort), so that I can compare analyses.
9. As a team member, I want each Plans entry to show the proposed totals and whether the team authorized them, so that I can find plans waiting for authorization.
10. As a team member, I want each Workflows entry to show the estimate against the budget, a pause at the gate, or how the run ended with its pull request and observed spend, so that I can find runs that need me.
11. As a team member, I want a figure that is not known yet to read as a dash, never a zero, so that an unmeasured value never looks measured.
12. As a team member, I want to see when each ticket was created, so that I can tell recent work from old work.
13. As a team member, I want the whole widget to open the ticket, by mouse, by touch and by keyboard, so that I never aim at a small link.
14. As a keyboard user, I want one tab stop per ticket with a focus ring around the whole widget, so that I can see where I am.
15. As a screen reader user, I want each entry's link to be named by the ticket's number and title, so that the list reads as tickets, not as repeated "open" links.
16. As a team member, I want older tickets to load when I reach the end of a listing, and a Load more button that does the same, so that scrolling and the keyboard both work.
17. As a team member, I want a listing that cannot reach extub to say so once and stop asking, so that a failure never turns into a stream of requests.
18. As a team member, I want an entry in the Analysis, Plans or Workflows listing to open the ticket page at that phase, below the console's header, so that I land on what I came for.
19. As a team member, I want the ticket page to read as one document, with no card around any part, so that it reads like a deployment page, not a dashboard of boxes.
20. As a team member, I want the ticket's title as the page heading with its number and slug above it, so that the page names the work first.
21. As a team member, I want the ticket's status, its phase as "Plan · 3 of 4" and when it was created as labels over values near the top, so that I take in where it stands at a glance.
22. As a team member, I want the linked artifacts listed by kind with their file names, so that I can find the PRD, design doc, spec and guardrail.
23. As a team member, I want each phase on the ticket page under a heading with its summary line on the right, so that the page's outline matches the four phases.
24. As a team member, I want each phase's body unchanged on the ticket page, so that the documents, episodes, components, report, stages, spend and pull request are all still there.
25. As a team member, I want the ticket page's polling not to move my scroll position, so that reading is never interrupted.
26. As a team member, I want loading placeholders shaped like the cards, rows and ticket page they precede, so that nothing jumps when the content arrives.
27. As a team member reading pt-BR, ja or zh-Hant, I want the new labels, dates and intros in my language, so that nothing depends on English.
28. As a team member on a 375px screen, I want no horizontal overflow in any listing or on the ticket page.
29. As a team member, I want the listing intros to describe what the listings now show, so that no intro promises content that only the ticket page has.

## Implementation Decisions

- **No scroll box for tickets.** The ticket listing renders an ordered list in the page's flow. A new console UI piece, `ListEnd`, renders the Load more button and observes it with an `IntersectionObserver` ahead of the viewport. It calls the paged resource's `next` while there are more pages, nothing is loading and the last request has not failed. Each page that lands re-arms it, so a short page keeps filling the screen; a failure disarms it until the user retries. `VirtualList` stays as it is for the Executions and History overviews.
- **Two widgets, one entry.** `TicketCard` (below `lg`) and `TicketRow` (from `lg`) render the same ticket. Each is hidden at the other's breakpoint, so only one is in the tab order and the accessibility tree. Both end with one link covering the widget, named by the number and title through `aria-labelledby`, without prefetch, as the activity entries.
- **Row grammar.** From `lg` the row has columns for the number, the title over the slug, the status over the created date, and the caret, with the figure as a third line under the slug: beside the 224px sidebar the content column is about 750px at `lg`. Below `lg` (about 500px at `md`) a row would cut titles to a few words, so tablets get the card. From `xl` the figure moves under the status and the date gets its own column, as in Vercel's deployment rows. The columns have fixed widths so rows align.
- **The figure.** `PhaseSummary` renders the one-line summary of the analysis, plan or workflow phase. The three `summaryLine` functions move into it from their widgets, which render it in their headers too. It renders nothing when that phase's data is null; unknown values inside a phase keep the dash. The plan's authorization label moves into the ticket format hook, shared by the plan body and its summary.
- **Which phase an entry summarizes.** `listedPhase(ticket, listing)`: the listing's phase in Analysis, Plans and Workflows; in Tickets, the ticket's current phase, or none while it is still in the Ticket phase. `SummaryPhase` is the phase type without `ticket`.
- **Where an entry opens.** `ticketPath(ticket, listing)`: the ticket page, with `#analysis`, `#plan` or `#workflow` in the three phase listings. The fragment is the phase key, not the route segment.
- **Dates.** `TicketCreated` renders `createdAt` as an absolute date in a `time` element: month, day and time in the listings, and the full date with its year on the ticket page. The console's request configuration fixes the time zone to UTC, so server and client render the same text; relative time would read two clocks and break hydration.
- **The ticket page.** `TicketDetail` renders a header (the number and slug as one mono line, as the breadcrumb names the ticket, over the title as `h1`), then `TicketDescription` (a `dl` strip of Status, Phase and Created, label over value, as the runtime detail's figures), then `TicketArtifacts`, then each phase's widget or its "not reached yet" line. Status appears once, in the strip.
- **Sections.** The widget frame becomes `TicketSection`: a `section` under a full-width hairline, with an `h2` and the summary line, scrolled into view below the 52px sticky header. The Analysis, Plan and Workflow widgets wrap their unchanged bodies in it. The Ticket widget becomes `TicketArtifacts`, a section around the existing definition list of the four files; its 40px number, slug, status and title leave because the page header carries them.
- **Landing on a fragment.** `useFragmentLanding` scrolls to the URL's fragment once the ticket has rendered. Next gives up on a fragment whose target is not in the DOM while the route's loading state shows, so without it a Plans entry would not land on `#plan`.
- **Phase progress.** `phaseProgress(ticket)` returns the phases reached and the total, for "Plan · 3 of 4".
- **Removed.** The `listed` prop of the phase widgets, the frame's link line, the Ticket widget's linked number and the `openTicket` copy key. The breadcrumb builds its ticket label with `ticketSummary`.
- **Skeletons.** A ticket entry skeleton (card below `lg`, row from `lg`), a ticket list skeleton for the four listings' loading states and Suspense fallbacks, and a ticket detail skeleton for the ticket page in place of the runtime detail skeleton.
- **Copy.** New `tickets.ticket.status`, `tickets.ticket.phase`, `tickets.ticket.created` and `tickets.ticket.phaseProgress` ("{phase} · {reached} of {total}"). The four listing intros are rewritten to describe the entries and point to the ticket page. Same keys in the console's four catalogs.
- **Supersedes, in spec 0014:**
  - In the Solution: "an infinite list of cards built on the website's panels" and "each card shows the panel's header line over the panel's full body, with no collapse". Entries show the header line only; the full body is on the ticket page.
  - User story 3 and stories 8 to 10. Their content is now on the ticket page.
  - In the Listing decision: "renders with `VirtualList` and keeps the load-more fallback, exactly as the executions overview".
  - In the Widgets decision: "a card with an eyebrow". A widget is now a section with a heading.
  - In the Detail decision: "the four widgets in order". The detail shows the artifacts and the three phase widgets under a header and strip.
  - In the manual check: "each card matches its website panel". Each phase section's body matches its panel.

## Testing Decisions

A good test observes what a user sees or what a function returns for a ticket, never markup structure.

1. **Console logic** (Node test runner, beside the existing ticket phase tests): `listedPhase` for each listing and phase, including none for a Ticket-phase ticket in Tickets and `analysis` for a workflow ticket in Analysis; `ticketPath` with no fragment in Tickets and the phase key in the others; `phaseProgress` at 1 of 4 and 4 of 4; `summaryPhases` is the phase order without `ticket`. The summary-line tests already in place keep covering the figures.
2. **Repository checks**: `pnpm --filter console test`, `pnpm --filter console typecheck`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`, and a diff of the leaf keys of the four console catalogs, since the console has no catalog parity test.

Manual verification with `pnpm console`:

- No inner scroll box in the four listings, and no horizontal overflow at 375px.
- Rows at 768, 1024 and 1280 wide.
- Each cursor requested once while scrolling, and the screen keeps filling on a tall window.
- With extub stopped, one error and no request loop; Load more and Retry resume.
- One focus ring around the whole widget, and Enter opens the ticket.
- A Plans entry lands on its section below the header, on navigation and on reload.
- No card on the ticket page, and polling never moves it.
- The breadcrumb reads "Project / Tickets / number · slug".
- `en` and `pt-br`, light and dark, and no hydration warning.
- Executions and History unchanged.

## Out of Scope

- The Executions and History overviews, which keep `VirtualList`.
- Rewriting the phase bodies as label-value grids; they stay the website panels' bodies.
- Status colours for ticket statuses, which stay muted.
- Filtering, sorting or searching the listings.
- A `vi` catalog for the console.

## Further Notes

References: Vercel's Geist Entity rows and Description, and shadcn/studio's project-card widget and list of separate small cards, for composition only. The colours, type and icons stay the console's.

`feat/search-coverage-and-hardening` also changes the console ticket pages' titles and the catalogs (its tickets 0114 and 0115); whichever lands second rebases.
