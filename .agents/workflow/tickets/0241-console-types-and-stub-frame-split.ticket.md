# 0241: the console's types and stub frame split by concern

**Spec:** .agents/workflow/specs/0034-general-rules-compliance-in-apps.spec.md

**What to build:** The console's data types are split into one `.types.ts` per concern (account, product, runtime, execution, recommendation, ticket, page), imported with `import type`, nothing renamed. The stub frame's generic guards (`isRecord`, `isText`, `isDate`, `isNumber`, `isNullable`, `hasFields`) become one validation file each, the component and widget guards keep their frame concern, the parser composes them, and its seven compound conditions get names. A test's helpers move to `testdata/` beside it; files are renamed after their export (the server read after `readConsole`, the locale action after `setLocale`).

**Blocked by:** 0230, 0231, 0236.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] The console's types file is split by concern and no file declares types of two concerns.
- [ ] The stub-frame test keeps refusing every malformed frame it pins today.
- [ ] No console test file declares a helper function and every console file is named after its export.
- [ ] `pnpm turbo run check typecheck test build --filter=console` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift quality check general`
