# 0296: the product creation and folder picker forms keep one responsibility per file

**Spec:** .agents/specs/0042-state-rules-compliance.spec.md

**What to build:** The product creation form gets its types, defaults, maps (the name from a selected path, a dropped file, the request to the folder or products endpoint), a schema requiring a non-empty path and bounding the name, and handlers for selection, drop, browse, path and name change, picker close and submit. The folder picker's submit handler builds the listing request through a map with the abort signal.

**Blocked by:** 0230, 0295.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] Each map holding logic has a test with the console's node runner.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
