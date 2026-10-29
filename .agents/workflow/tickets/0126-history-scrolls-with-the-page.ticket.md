# 0126: The History overview scrolls with the page

**Spec:** .agents/workflow/specs/0020-console-activity-scrolls-with-the-page.spec.md

**What to build:** the History overview gets the same page-scrolled list as the Executions overview. The events form an ordered list in the page's flow, each in a bordered item of its own, older events load as the end of the list comes into view or when Load more is pressed, and a list that cannot reach extub says so once and stops asking until the user retries. An event keeps its body: the runtime link to the runtime's history, the kind, the cause and actor, the environment and the two dates, the ids and the detail, with long actors, ids and details wrapping at 375px. Its loading placeholders are the ones the Executions overview introduced.

**Blocked by:** 0125, The Executions overview scrolls with the page.

**Status:** ready-for-agent

- [x] The History overview renders an ordered list named by the page, with no scroll box, fixed height, list-wide border or virtualization; Refresh stays above it, and the loading row, error notice and empty line below.
- [x] The overview ends with `ListEnd` fed by the paged resource, and its own Load more button is gone.
- [x] `HistoryEntry` renders its body in an item of its own with a border, rounded corners and the card colour, no hover colour and no link over the whole entry; the runtime link stays the only link, without prefetch.
- [x] Long runtime names, actors, ids and details wrap, so 375px has no horizontal overflow.
- [x] The History page uses `ActivityOverviewSkeleton` for its loading state and Suspense fallback, and the loading-more row uses `ActivityEntrySkeleton` in a status region.
- [x] The per-runtime History page and `ActivityListSkeleton` are unchanged.
- [x] No copy changes; the four console catalogs keep the same keys.

### Verification

- `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`.
- In the browser, with more than one page of events in extub: no inner scroll box and the first entry fully visible; scrolling requests each cursor once and the screen keeps filling on a tall window; with extub stopped, one error and no request loop, and Retry and Load more resume; no horizontal overflow at 375px; the runtime link opens the runtime's history; no hydration warning in `en` and `pt-br`, light and dark; the per-runtime History page loads as before.
