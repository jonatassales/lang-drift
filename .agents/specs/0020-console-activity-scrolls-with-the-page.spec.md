# Console activity overviews scroll with the page

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The console's Executions and History overviews draw every entry inside one bordered box with a fixed height of 65% of the screen, which scrolls on its own inside a page that also scrolls. The first entry is clipped at the top of the box, a team member scrolls twice to reach the end of a long list, and a Load more button sits under a box that also loads as it scrolls. The four ticket listings had the same problems and were fixed in spec 0019; these two overviews were left behind on purpose and now look unfinished beside them.

## Solution

Each overview becomes a page-scrolled list: an ordered list in the page's flow with one bordered item per entry, a gap between items and no card around the list. The page is the only thing that scrolls. Older entries load when the end of the list comes into view, and a Load more button does the same by keyboard. A list that cannot reach extub says so once and stops asking until the user retries.

The entries keep their bodies. An execution is still a runtime link and its state over the environment and execution id, with the CPU, RAM, revision and observation figures; an event is still a runtime link and its kind over the cause and actor, with the environment and the two dates, the ids and the detail. Each body moves into an item of its own. The runtime link stays the only link of an entry: an entry holds ids and evidence a team member copies, and no execution or event page exists to open.

Loading placeholders take the shape of the new items. `VirtualList` and the dependency it needed leave the console, since nothing uses them any more.

## User Stories

1. As a team member, I want the Executions overview to scroll with the page instead of inside a box, so that nothing is clipped and there is one scrollbar.
2. As a team member, I want the History overview to scroll with the page instead of inside a box, so that nothing is clipped and there is one scrollbar.
3. As a team member, I want each execution as its own bordered item, so that I can tell where one execution ends and the next begins.
4. As a team member, I want each event as its own bordered item, so that I can tell where one event ends and the next begins.
5. As a team member, I want no card around the whole list, so that the items read as a list, not as boxes inside a box.
6. As a team member, I want the first entry fully visible when the page opens, so that nothing starts clipped.
7. As a team member, I want older entries to load when I reach the end of the list, so that I never press a button to keep reading.
8. As a keyboard user, I want a Load more button after the last entry, so that I can load older entries without a pointer.
9. As a team member, I want each page of entries requested once, so that scrolling never repeats a request.
10. As a team member on a tall window, I want the list to keep loading until the screen is filled, so that a short page never leaves the end of the list in view with more to read.
11. As a team member, I want a list that cannot reach extub to say so once and stop asking, so that a failure never turns into a stream of requests.
12. As a team member, I want Retry and Load more to resume loading after a failure, so that I decide when to try again.
13. As a team member, I want Refresh above the list to reload the newest entries, so that I see current observations.
14. As a team member, I want to see a placeholder shaped like the entries while the first page loads, so that nothing jumps when the content arrives.
15. As a team member, I want a placeholder shaped like one entry while the next page loads, so that I know more is coming.
16. As a team member, I want to see the empty line only when there are no entries and nothing failed, so that an error never reads as an empty project.
17. As a team member, I want each execution to keep its runtime link, state, environment, execution id, CPU, RAM, revision and observation time, so that nothing the overview showed is lost.
18. As a team member, I want a figure that was not observed to keep reading as unavailable, never as zero, so that an unmeasured value never looks measured.
19. As a team member, I want each event to keep its runtime link, kind, cause, actor, environment, dates, ids and detail, so that nothing the overview showed is lost.
20. As a team member, I want to select and copy an execution id, an event id or a detail text without opening anything, so that I can paste evidence elsewhere.
21. As a team member, I want the runtime link in an execution to open the runtime and the one in an event to open the runtime's history, so that I reach the runtime in context.
22. As a screen reader user, I want each list to be a named list of entries, so that I can tell how many records there are and move through them.
23. As a keyboard user, I want one tab stop per entry, on its runtime link, so that I can reach every entry without crossing its text.
24. As a team member on a 375px screen, I want no horizontal overflow, even with a long runtime name, actor or id, so that nothing is cut off or needs a sideways scroll.
25. As a team member reading pt-BR, ja or zh-Hant, I want the pages to read the same in my language, so that nothing depends on English.
26. As a team member, I want the per-runtime History page, with its table, to load and look as it does today, so that this change does not touch it.
27. As a maintainer, I want `VirtualList` and the virtualization dependency gone once nothing uses them, so that no dead code or unused dependency is kept.

## Implementation Decisions

- **Page-scrolled list.** The Executions and History overviews render an ordered list in the page's flow, labelled with the page's name, one list item per entry, with a gap between items. There is no scroll box, no fixed height, no list-wide border and no virtualization. Refresh stays above the list; the loading row, the error notice and the empty line stay below it.
- **The end of the list is `ListEnd`.** Each overview ends with the `ListEnd` piece from spec 0019, fed with the paged resource's `hasMore`, `loading` and `next`, and `failed` from its error. Its behavior is already settled: it asks for the next page as the end comes within reach of the viewport, re-arms after each page, and stays quiet after a failure until Retry or Load more. The overviews' own Load more button is removed, since `ListEnd` renders it. The paged resource is unchanged.
- **Entries keep their bodies, each in its own item.** `ExecutionEntry` and `HistoryEntry` each render an `article` with a bordered, rounded, card-coloured frame and padding, and no hover colour, since the item itself is not a link. Their contents are unchanged except that long unbroken text wraps: runtime names, actors and ids break rather than widen the page. The list item wrapper replaces the virtual row. There is no card and row split, because the figures are labelled values and the entries are already compact; a card and row would need repeated labels and add no information.
- **The runtime link stays the only link.** An execution links to its runtime and an event to its runtime's history, without prefetch, as today. No link covers the whole entry.
- **Skeletons.** A new `ActivityEntrySkeleton` is one bordered item shaped like an entry. A new `ActivityOverviewSkeleton` is a Refresh placeholder over six entry skeletons, a `section` marked busy with the loading label, as the ticket list skeleton. The Executions and History `loading.tsx` and Suspense fallbacks use it, and the loading-more row uses the entry skeleton in a status region. `ActivityRowSkeleton` is deleted. `ActivityListSkeleton` is left exactly as it is, because it still loads the per-runtime History page, whose body is a card holding a table.
- **Dead code.** `VirtualList` is deleted. `@tanstack/react-virtual` is removed from the console's dependencies, through pnpm, so the manifest and the lockfile change. The copy of the dependency in the `ex` package's sample project fixture is left alone. `drift check dead-code` passes without a new allow entry.
- **Copy.** None changes. The shared intro (`activityIntro`, "Scroll for older records; refresh for current observations") is still true, and the list labels reuse the existing `executions` and `history` keys. The four console catalogs keep the same keys.
- **Supersedes, in spec 0019:** in the Implementation Decisions, "`VirtualList` stays as it is for the Executions and History overviews", and in Out of Scope, "The Executions and History overviews, which keep `VirtualList`".
- **Supersedes, in spec 0014:** in the Listing decision, "exactly as the executions overview", which no longer renders with `VirtualList`.

## Testing Decisions

A good test observes what a user sees or what a function returns, never markup structure.

1. **No new automated test.** This change moves markup and removes code; it adds no pure logic. `ListEnd`, the paged resource and the ticket logic keep the behavior and the tests they have, and the console has no React test harness to extend.
2. **Repository checks**: `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`, and a diff of the leaf keys of the four console catalogs, since the console has no catalog parity test.

Manual verification with `pnpm console`, with enough activity in extub to fill more than one page (a page holds 50; the stub seeds no records, so start and stop a runtime repeatedly):

- No inner scroll box on either page, and the first entry fully visible.
- Each cursor requested once while scrolling, and the screen keeps filling on a tall window.
- With extub stopped, one error and no request loop; Load more and Retry resume.
- No horizontal overflow at 375px, with a long runtime name, actor and id.
- Refresh reloads the newest entries.
- The runtime link of an execution opens the runtime, and an event's opens its history; an id can be selected and copied.
- `en` and `pt-br`, light and dark, and no hydration warning.
- The per-runtime History page loads and looks as before.

## Out of Scope

- The per-runtime History page and its table, and the runtime detail's execution table.
- Turning the entries into the ticket listings' card and row.
- A link over the whole entry, an execution page and an event page.
- Filtering, sorting or searching the overviews.
- A change to the paged resource, the page size or extub's activity endpoints.
- A `vi` catalog for the console.

## Further Notes

`ActivityListSkeleton` keeps a fixed-height box shape that fits no listing now. It is kept because renaming and reshaping it for the per-runtime History page is a separate change, which this spec does not make.

Spec 0019 and its tickets 0122 to 0124 are the prior art: the list, its end and its loading states are built the same way here.
