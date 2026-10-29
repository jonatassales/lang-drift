# 0125: The Executions overview scrolls with the page

**Spec:** .agents/workflow/specs/0020-console-activity-scrolls-with-the-page.spec.md

**What to build:** the Executions overview stops drawing its executions inside a fixed-height scroll box. The executions form an ordered list in the page's flow, each in a bordered item of its own with a gap between items and no card around the list, and the page is the only thing that scrolls. Older executions load when the end of the list comes into view or when Load more is pressed. A list that cannot reach extub says so once and stops asking until the user retries. An execution keeps its body: the runtime link, the state, the environment and execution id, and the CPU, RAM, revision and observation figures, with long names and ids wrapping at 375px. The loading placeholders take the shape of the new items.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] The Executions overview renders an ordered list named by the page, with no scroll box, fixed height, list-wide border or virtualization; Refresh stays above it, and the loading row, error notice and empty line below.
- [ ] The overview ends with `ListEnd` fed by the paged resource, and its own Load more button is gone.
- [ ] `ExecutionEntry` renders its body in an item of its own with a border, rounded corners and the card colour, no hover colour and no link over the whole entry; the runtime link stays the only link, without prefetch.
- [ ] A figure that was not observed still reads as unavailable, never as zero.
- [ ] Long runtime names and ids wrap, so 375px has no horizontal overflow.
- [ ] `ActivityEntrySkeleton` (one bordered item shaped like an entry) and `ActivityOverviewSkeleton` (a Refresh placeholder over six entry skeletons, busy and labelled as loading) exist; the Executions page uses the overview skeleton for its loading state and Suspense fallback, and the loading-more row uses the entry skeleton in a status region.
- [ ] `ActivityListSkeleton` and the per-runtime History page are unchanged.
- [ ] No copy changes; the four console catalogs keep the same keys.

### Verification

- `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`.
- In the browser, with more than one page of executions in extub: no inner scroll box and the first entry fully visible; scrolling requests each cursor once and the screen keeps filling on a tall window; with extub stopped, one error and no request loop, and Retry and Load more resume; no horizontal overflow at 375px; the runtime link opens the runtime and an id can be selected; no hydration warning in `en` and `pt-br`, light and dark.
