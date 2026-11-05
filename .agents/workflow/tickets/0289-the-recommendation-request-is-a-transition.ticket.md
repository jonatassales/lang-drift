# 0289: the recommendation request runs in one transition

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** The recommendation panel's `reading` flag goes: reading the ticket file and requesting the recommendation run inside the existing transition, so `pending` models both, and a newer submission aborts the earlier one. The abortable read and the request function 0242 extracted live in a recommendation request hook beside the runtimes views, where 0230 places them.

**Blocked by:** 0242.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] No `useState` loading flag remains in the panel.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Category:** `hard-part`. Aborting a superseded file read and request inside a transition is where this ticket's difficulty sits.

**Paths:**
- `apps/console/src/app/[locale]/(menu)/runtimes/recommendation-request.hook.ts`
- `apps/console/src/app/[locale]/(menu)/runtimes/recommendation-request.hook.test.ts`

**Boundary:** `useRecommendationRequest(runtimeId: string): RecommendationRequest`, called by the panel, where `RecommendationRequest` (the submit function, `pending`, the result and the error) is declared by the agent in `recommendation-request.types.ts` beside it.

**Goal:** Read the ticket file and request the recommendation in one transition, abort the earlier read and request when a newer submission starts, and expose `pending` for both.

- [ ] Tests with a fake file and `fetch` show that a newer submission aborts the earlier and that `pending` covers the read.
- [ ] An aborted request never writes its result.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
