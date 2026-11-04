# 0188: `@langdrift/events/transport`: Redis Streams strategy and client

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** A Nest application can consume a Redis Stream as a microservice and append to it as a client. The package's `transport` entry ships a custom transport strategy that creates the consumer group with `MKSTREAM`, reads with `XREADGROUP … BLOCK` on a dedicated connection, dispatches each entry to the `@EventPattern` handler registered for its type or to the wildcard handler, acknowledges only after the handler's result completes, leaves the entry pending on a throw, claims idle pending entries with `XAUTOCLAIM` at startup and on an interval, parks an entry past the maximum deliveries on the dead-letter stream, acknowledges and skips an entry no handler wants, and closes gracefully; and a client that appends an envelope with `XADD … MAXLEN ~` and refuses request-response.

**Blocked by:** 0187.

**Status:** ready-for-agent

- [x] The strategy implements Nest's custom transport contract and is constructed with the connection URL, the stream, the group, the consumer name and bounded options (block, batch, claim idle, claim interval, maximum deliveries).
- [x] An entry whose handler completes is acknowledged; an entry whose handler throws stays pending and is claimed on the next claim pass.
- [x] An entry delivered more times than the maximum is appended to the dead-letter stream and acknowledged.
- [x] An entry with no handler for its type and no wildcard handler is acknowledged and skipped.
- [x] Closing the strategy breaks the blocking read by disconnecting its connection, awaits the batch in flight and quits the command connection.
- [x] The client resolves `emit` with the stream id of the appended entry and rejects `send`.
- [x] Every behavior above is covered by unit tests against a recording fake of the `ioredis` subset in use; no test reaches Redis.

### Verification

- `pnpm turbo run build typecheck test --filter=@langdrift/events`
- `pnpm turbo run check`
- `platform/cli/drift quality check 00-general 04-typescript 10-testing`
