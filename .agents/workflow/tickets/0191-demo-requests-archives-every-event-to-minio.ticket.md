# 0191: demo-requests archives every event to MinIO

**Spec:** .agents/workflow/specs/0031-book-a-demo-microservices.spec.md

**What to build:** Every event on the stream becomes an object in the MinIO bucket, so a trimmed stream or a restarted Redis never loses a request's history. The service becomes a hybrid application that also runs the `archive` consumer group on the Redis Streams transport; a wildcard handler stores each envelope as `demo-request/<aggregateId>/<streamId>.json` with the event type as object metadata, and the archive service ensures the bucket exists when the module starts. The MinIO client is a token-injected provider whose factory reads the endpoint, port, TLS flag, credentials and bucket from the environment at creation.

**Blocked by:** 0189.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] After a request is posted, the bucket holds `demo-request/<id>/<streamId>.json` with the envelope as its body and the type in its metadata.
- [x] A mailer event for the same request (0192) lands beside it under the same prefix.
- [x] The bucket is created once on module init when absent and left alone when present.
- [x] The archive handler acknowledges only after the object is written; a failed write leaves the entry pending.
- [x] The service shuts down gracefully on SIGTERM, finishing the batch in flight.
- [x] Unit tests cover the object key, the bucket bootstrap and the store call with a fake MinIO client; none reaches MinIO.

### Verification

- `pnpm turbo run build typecheck test --filter=demo-requests`
- `pnpm turbo run check`
- `platform/cli/drift quality check 00-general 04-typescript 10-testing`
- Manual: start MinIO from its Dockerfile, post a request, list the bucket.
