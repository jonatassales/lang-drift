# 0160: The website's client forms are split by responsibility

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** The early-access form's logic gets one `<event>.handler.ts` per user event, one `<mutation>.map.ts` per state mutation and its Zod schema in `form.validation.ts`, beside `form.types.ts` and `form-state.data.ts`. The book-a-demo form gets `form.types.ts`, `form-state.data.ts` and `form.validation.ts`. Both forms behave as before.

**Blocked by:** 0159

**Status:** completed

- [x] The early-access form has one handler file per event and one map file per mutation, each with one function.
- [x] The book-a-demo form has its types, its defaults and its schema in their own files.
- [x] Both forms' tests pass, moved with the code they cover.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check general`
