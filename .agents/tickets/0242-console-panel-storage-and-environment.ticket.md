# 0242: the recommendation panel, theme storage and policy environment apply the general rules

**Spec:** .agents/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** The recommendation panel splits into its ticket upload handler, its recommendation request (a plain function that 0289 later wraps in a hook) and its view, the reading version and the abort controller staying with the handler that owns them. The console's theme preference is persisted through an adapter that returns whether the write took. The proxy's policy builder reads `NODE_ENV` and `VERCEL_ENV` where it chooses the policy's environment, and the policy-environment module is deleted. The remaining compound conditions are named and the chained ternaries of the recommendation accounting and the runtime list become lookups or guards.

**Blocked by:** 0233, 0241.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] The panel's upload, request and view live in three files, each with one reason to change.
- [ ] No empty or comment-only `catch` remains in the console, and the storage adapter answers `false` when `localStorage` throws, with a unit test.
- [ ] `policy-environment.compute.ts` is gone and the content-security-policy tests pass.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes and no screen changes in the built-in browser.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift quality check general`
