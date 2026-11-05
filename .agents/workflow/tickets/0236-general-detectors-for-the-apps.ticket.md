# 0236: drift sees the general-rule violations in the apps

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** `drift quality check general` gains one row per class the audit found that a pattern or an awk program catches: a TypeScript guard directly after the declaration it tests; a chained ternary; a compound condition with three or more operands in an `if` or `while`, across lines; a swallowed rejection and a comment-only `catch`; a `_` parameter; an awaited result consumed before it is named in a test; test files in scope for the multi-function and return-spacing rows. The command keeps listing without failing.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] Each new row prints the known sites of its class in the current tree and nothing from a named negative set.
- [ ] The guard-after-declaration program excludes Go, whose `if err != nil` sits under its assignment by the rule's own example.
- [ ] `drift quality check` still exits 0 with findings.

## Dev part

**Category:** `hard-part`. Telling a condition's operands apart across lines, through parentheses, strings and template literals, is where this ticket's difficulty sits.

**Paths:**
- `platform/cli/src/quality/compound-condition.awk`
- `platform/cli/src/quality/testdata/compound-condition/`

**Boundary:** The program reads the TypeScript files `run-detectors.sh` passes it, as `return-spacing.awk` does, and prints one `file:line: message` per finding; the agent writes its row in `rule-detectors.tsv`.

**Goal:** Report every `if` or `while` whose condition combines three or more operands with `&&` or `||`, on one line or across lines, and nothing else.

- [ ] It reports the 31 unnamed compound conditions the audit counted under `apps/` (seven in the stub frame validator, three in the cookie service, two in the locale proxy among them).
- [ ] It reports no condition of two operands, no operand inside a string or a template literal, and no `&&` in JSX; its fixtures hold one file it must report and one it must not, which 0246's `drift quality detectors` runs.

## Verification

- `platform/cli/drift quality check general`
- `platform/cli/drift quality coverage general`
