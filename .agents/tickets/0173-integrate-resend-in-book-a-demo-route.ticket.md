# 0173: Integrate email sending in book-a-demo API route

**Spec:** .agents/specs/0029-resend-book-a-demo-integration.spec.md

**What to build:** Update `apps/website/src/app/api/book-a-demo/route.ts` to replace the 503 pause response with real email dispatch using the Resend adapters. Update `book-demo-reply.response.ts` to support success responses (`200 OK`) and failure response (`500` / `DELIVERY_FAILED`). Update `apps/website/src/app/api/book-a-demo/route.unit.ts` to test successful dispatch and failure handling.

**Blocked by:** 0172

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `apps/website/src/app/api/book-a-demo/route.ts` calls `sendDemoConfirmation` and `sendDemoLeadNotification`.
- [x] `apps/website/src/app/lib/components/book-a-demo/book-demo-reply.response.ts` handles success and error responses.
- [x] `apps/website/src/app/api/book-a-demo/route.unit.ts` tests successful submission and error handling.
- [x] Existing origin, validation, and size checks remain intact and tested.

### Verification

- `pnpm turbo run test --filter=website`
- `platform/cli/drift check biome`
- `platform/cli/drift quality check 00-general 06-nextjs`
