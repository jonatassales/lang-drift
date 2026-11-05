# 0158: The website's control flow follows the general rules

**Spec:** .agents/workflow/specs/0027-website-rules-compliance.spec.md

**What to build:** Nested ternaries and conditionals deeper than two levels become guards, a `switch` or a lookup map. Compound predicates that combine several business or runtime rules get a name. Every empty `catch` states an explicit fallback with the same effect as today, so no error is swallowed in silence.

**Blocked by:** 0157

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] No ternary or conditional nests beyond two levels.
- [ ] Compound predicates are named.
- [ ] No `catch` block is empty.
- [ ] Behavior is unchanged: the test suite, the build and the captures still pass.

### Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `platform/cli/drift quality check general`
