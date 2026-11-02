# 0170: Implement route handlers directly in website API routes

**Spec:** .agents/workflow/specs/0028-route-handlers-own-implementation.spec.md

**What to build:** Refactor `apps/website/src/app/api/book-a-demo/route.ts` to implement `export async function POST(request: Request)` directly in the file, removing `apps/website/src/app/lib/components/book-a-demo/submission.route.ts` and moving its unit test to `apps/website/src/app/api/book-a-demo/route.unit.ts`. In `apps/website/src/app/api/early-access/`, remove `apps/website/src/app/lib/components/early-access/submission.route.ts` and ensure any early-access handler logic is in `route.ts`. Update associated unit tests (`route.unit.ts`, `provider-configuration.unit.ts`).

**Blocked by:** 0169

**Status:** ready-for-agent

- [x] `apps/website/src/app/api/book-a-demo/route.ts` declares `export async function POST(request: Request)` directly.
- [x] `apps/website/src/app/lib/components/book-a-demo/submission.route.ts` is deleted.
- [x] `apps/website/src/app/api/book-a-demo/route.unit.ts` tests `POST` directly on `route.ts`.
- [x] `apps/website/src/app/lib/components/early-access/submission.route.ts` is deleted.
- [x] All website tests (539 unit tests) pass.
- [x] No quality detector violations in `apps/website`.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check`
