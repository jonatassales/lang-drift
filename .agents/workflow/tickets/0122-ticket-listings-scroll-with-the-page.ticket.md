# 0122: The ticket listings scroll with the page

**Spec:** .agents/workflow/specs/0019-console-ticket-widgets.spec.md

**What to build:** the Tickets, Analysis, Plans and Workflows listings stop drawing their tickets inside a scroll box of their own. The tickets form an ordered list in the page's flow, the page is the only thing that scrolls, and older tickets load when the end of the list comes into view or when Load more is pressed. A listing that cannot reach extub says so once and stops asking until the user retries. The entries themselves are unchanged until ticket 0123.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] A console UI piece renders Load more and, while there are more pages, nothing is loading and the last request has not failed, calls the listing's next page when it comes within reach of the viewport.
- [x] Each page that lands re-arms it, so a short page on a tall window keeps loading until the screen is filled.
- [x] A failed request disarms it; Retry and Load more resume.
- [x] The ticket listing renders an ordered list with no scroll box, keeps Refresh above it and the loading row, error notice and empty line below.
- [x] `VirtualList`, the Executions and the History overviews are unchanged.

### Verification

- `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run check`, `drift check dead-code`.
- In the browser: no inner scroll box in the four listings; scrolling to the end requests each cursor once; with extub stopped, one error and no request loop; Load more reachable by keyboard.
