# 0161: The website's statements are spaced as the rules say

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** `return`, `await`, props destructuring and block boundaries follow the general and TypeScript spacing rules in the website's source and tests. The tests use readable names instead of `t` or `(a, b)`.

**Blocked by:** 0160

**Status:** completed

- [x] `drift quality check` reports no spacing finding in `apps/website`.
- [x] The review finds no spacing breach the detectors miss, in source or tests.
- [x] No test uses a one-letter name.
- [x] `drift quality coverage` reaches 100% for the measured rules in the website, or every gap is a documented false positive.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check` and `platform/cli/drift quality coverage`
