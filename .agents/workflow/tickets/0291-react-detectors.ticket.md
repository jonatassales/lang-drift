# 0291: drift measures the React rules

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** `drift quality check react` gains: an arrow function assigned to a constant inside a TSX function body; a `useEffect` without a comment above it; a props interface not two lines above its component; a scope out of order; a `'use client'` file with no client feature; a `template.tsx` that returns its children; an inline `<svg>` in an app; the loading-flag row extended to `reading`, `fetching` and `busy`.

**Blocked by:** 0246, 0286, 0287, 0288, 0289, 0290.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] `platform/cli/drift quality check react` reports nothing.
- [ ] `platform/cli/drift quality detectors` covers every new awk program with a compliant and a violating component.

## Dev part

**Category:** `hard-part`. Classifying each statement of a component scope into the rule's order without a parser is where this ticket's difficulty sits.

**Paths:**
- `platform/cli/src/quality/scope-order.awk`
- `platform/cli/src/quality/testdata/scope-order/`

**Boundary:** The program reads the TSX files `run-detectors.sh` passes it and prints one `file:line: message` per finding; the agent writes the other detectors and the rows.

**Goal:** Report a component whose scope places a category before one the rule orders earlier: props destructuring, React hooks, custom hooks from `.hook` modules, third-party hooks, handlers, derived constants, the return.

- [ ] Its fixtures cover a compliant component, each out-of-order pair, a guard returning after the hooks, and a nested function that must be ignored.
- [ ] Run over the tree after 0286 to 0290, it reports nothing; false positives are documented in the pull request.

## Verification

- `platform/cli/drift quality detectors`
- `platform/cli/drift quality check react`
