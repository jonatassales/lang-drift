# 0239: the website's files hold one responsibility each

**Spec:** .agents/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** The website closes its structural general-rule findings: a test's helpers move to `testdata/` beside it, one function per file (fakes and drivers as `.mock.ts`, computations as `.compute.ts`); the files whose stem does not name their export are renamed after it with their responsibility's suffix, and `form-state.data` renames its export (`bookDemoFormState`); `BodyError` leaves the bounded body reader; every compound condition gets a name and every chained ternary becomes a lookup or guards; the language preference is persisted through an adapter that returns whether the write took; the body-reader cancel and the Turnstile prepare get a named rejection handler that reports through the existing observability path; the early-access form refuses a missing site key as a guard, so its `try` wraps only the Turnstile controller.

**Blocked by:** 0229, 0232, 0234, 0236.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment). Before: ready-for-agent. Its early-access items are superseded by spec 0046, which deletes that code. 2026-11-04.

- [ ] No test file under `apps/website` declares a helper function.
- [x] Every website file is named after its export.
- [x] No empty or comment-only `catch` and no discarded rejection remains in the website.
- [x] The storage adapter answers `false` when `localStorage` throws, with a unit test.
- [x] `pnpm turbo run check typecheck test build --filter=website` passes and no page changes in the built-in browser.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check general`

> **Note, 2026-11-05.** Every test helper under `apps/website/src` now lives in a `testdata/` folder, one per file. `next.config.unit.ts` keeps its two helpers, `headersAt` and `policyAt`: it sits at the app root, where no alias reaches, and the filesystem rules forbid both `../` and `./testdata/…` imports, so moving them needs an alias the owner decides on. That criterion stays open. The early-access items are obsolete; the multi-export `.types`, `.data` and metadata files are 0281's.
