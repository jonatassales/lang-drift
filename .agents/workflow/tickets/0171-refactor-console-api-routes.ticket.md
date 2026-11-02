# 0171: Implement route handlers directly in console API routes

**Spec:** .agents/workflow/specs/0028-route-handlers-own-implementation.spec.md

**What to build:** In `apps/console/src/app/api/extub/[...path]/route.ts`, implement `export async function GET(request: Request)` and `export async function POST(request: Request)` directly and delete `extub-proxy.route.ts`. In `apps/console/src/app/api/console/products/[id]/secrets/route.ts`, declare explicit `GET` and `POST` async functions that invoke `consoleHandler(request)`. In `apps/console/src/app/api/console/[...path]/route.ts`, declare explicit `GET`, `POST`, `PATCH`, `DELETE` async functions that invoke `consoleHandler(request)`.

**Blocked by:** 0169

**Status:** ready-for-agent

- [ ] `apps/console/src/app/api/extub/[...path]/route.ts` implements `GET` and `POST` directly.
- [ ] `apps/console/src/app/api/extub/[...path]/extub-proxy.route.ts` is deleted.
- [ ] `apps/console/src/app/api/console/products/[id]/secrets/route.ts` declares explicit `GET` and `POST` functions.
- [ ] `apps/console/src/app/api/console/[...path]/route.ts` declares explicit `GET`, `POST`, `PATCH`, `DELETE` functions.
- [ ] All console tests (40 unit tests) pass.
- [ ] No quality detector violations in `apps/console`.

### Verification

- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift quality check`
