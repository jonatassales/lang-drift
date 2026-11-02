# Route handlers own their implementation

**Status:** ready-for-agent

## Problem Statement

Next.js API route files (`route.ts`) were delegating their HTTP method handling to helper files (such as `submission.route.ts` and `extub-proxy.route.ts`) or re-exporting external functions as `export const POST = ...` or `export { consoleHandler as GET }`. This pattern was introduced to circumvent the repository's single-function-per-file detector, resulting in hollow 3-line `route.ts` files and scattering route handling logic into UI component folders (`@components/book-a-demo/submission.route.ts`, `@components/early-access/submission.route.ts`).

Next.js App Router conventions specify that route handler files (`route.ts`) define their HTTP method functions (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`) directly in the route file.

## Solution

1. **Rule definition:** Add Rule 9 to `.claude/rules/06-nextjs.rule.md` mandating that route handlers define their HTTP methods directly in `route.ts` and never act as empty forwarders. Clarify in `.claude/rules/00-general.rule.md` that each exported HTTP method function in `route.ts` is an entrypoint. Update quality detectors in `platform/cli/` to accommodate `route.ts` HTTP method exports.
2. **Website routes refactoring:** Move `postBookDemo` and `postEarlyAccess` request handling logic directly into their corresponding `src/app/api/**/route.ts` files. Delete redundant `submission.route.ts` files under `lib/components/`. Update unit tests to verify `POST` directly on the route handlers.
3. **Console routes refactoring:** Define explicit HTTP method handler functions in `apps/console/src/app/api/extub/[...path]/route.ts` (deleting `extub-proxy.route.ts`), `apps/console/src/app/api/console/products/[id]/secrets/route.ts`, and `apps/console/src/app/api/console/[...path]/route.ts`.

## User Stories

1. As a developer, I want `route.ts` to implement its HTTP methods directly so that route logic is immediately readable and follows standard Next.js conventions.
2. As a developer, I want component directories to contain UI components rather than server route handlers.
3. As the owner, I want the rules and detectors to enforce that `route.ts` cannot be a hollow re-exporter.

## Implementation Decisions

- Each `route.ts` implements its HTTP method exports (`export async function GET(request: Request)`, `export async function POST(request: Request)`, etc.) directly.
- Server validation, computes, and responses remain in single-responsibility modules where appropriate, but the route request lifecycle is owned by `route.ts`.
- Files named `*.route.ts` outside of `route.ts` directories are removed.
- Route tests reside beside the route file as `route.unit.ts` or in the test suite covering the endpoint.

## Testing Decisions

- `pnpm turbo run check typecheck test build` passes across all workspaces.
- `platform/cli/drift quality check` reports no violations.
