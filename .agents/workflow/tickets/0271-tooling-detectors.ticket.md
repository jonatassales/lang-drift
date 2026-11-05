# 0271: drift measures the tooling rules

**Spec:** .agents/workflow/specs/0037-tooling-rules-compliance.spec.md

**What to build:** `drift quality check tooling` flags automation that is not shell under `platform/`, `.github/` and `.claude/`, plus any `build.<js|mjs|cjs|ts|mts>` or `scripts/*.<js|mjs|cjs|ts|mts>` in an app or package; an extensionless file directly in `platform/cli/` or a package `bin/`; a package script that chains tools with `&&`, `||`, `;` or `pnpm run`. The scaffold assets stay outside the npm/npx and setup-node roots, as spec 0044 decided.

**Blocked by:** 0267, 0268, 0269.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] Each new row reports its class in the tree before its fix and nothing after.
- [ ] `platform/cli/drift quality check tooling` reports nothing.

## Dev part

**Override:** No part qualifies: detector rows are patterns in the existing table, configuration rather than a problem with its own goal (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check tooling`
