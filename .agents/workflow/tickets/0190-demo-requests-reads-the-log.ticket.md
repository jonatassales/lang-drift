# 0190: demo-requests reads the log

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** A caller watches the stream over Server-Sent Events and reads one request's state. `GET /v1/demo-requests/events` is a Nest `@Sse()` route that starts at the `Last-Event-ID` header when present, at the live end otherwise, replays from a given entry on request, emits each entry with the stream id as the event id and the type as the event name, sends a ping on every empty block and tears the blocking read down when the client disconnects. `GET /v1/demo-requests/:id` folds the aggregate's events into a view (state `received`, `confirmed` or `failed`, the submission, the emails sent or failed) and answers 404 when none exist.

**Blocked by:** 0189.

**Status:** ready-for-agent

- [x] `curl -N` with the bearer token and `Last-Event-ID: 0-0` prints every entry of the stream as an SSE event with `id`, `event` and `data`, then keeps the connection open with pings.
- [x] Without `Last-Event-ID` the stream shows only entries appended after the connection opened.
- [x] Disconnecting the client ends the reader's blocking loop and closes its duplicated connection.
- [x] `GET /v1/demo-requests/:id` answers the folded view; a request with a sent confirmation reads `confirmed`, one with a failed confirmation reads `failed`, an unknown id answers 404.
- [x] The fold, the SSE mapping and the reader are unit-tested with fakes; the SSE route is declared before the id route.

### Verification

- `pnpm turbo run build typecheck test --filter=demo-requests`
- `pnpm turbo run check`
- `platform/cli/drift quality check 00-general 04-typescript 10-testing`
- Manual: post a request (0189), watch it arrive on the SSE stream, read its view.
