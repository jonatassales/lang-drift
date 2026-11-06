# 0265: no path climbs with ../

**Spec:** .agents/specs/0036-filesystem-rules-compliance.spec.md

**What to build:** `drift` derives the project root from its own absolute directory with `dirname`, so the path never contains `..` and resolves the same from any working directory and through the symlink `pnpm install` creates; the pre-commit command calls the dispatcher through its absolute `DRIFT_CLI_DIR`. The console proxy's matcher test becomes a sibling of `proxy.ts`. The two Nest READMEs link the Nest rule root-absolute. The `../` detectors widen: every quoted `../` string in TypeScript and JavaScript, tests included; the shell detector's roots add `.claude`, `apps`, `lib` and `packages`; a new row flags a `../` inside a Markdown link.

**Blocked by:** 0233, 0263.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] `platform/cli/drift quality check filesystem` reports no `../` anywhere.
- [ ] `drift --help` and `drift quality check filesystem` give the same result from the root, a nested directory and the installed symlink.
- [ ] The console proxy's matcher test passes beside `proxy.ts`.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift quality check filesystem`
- `pnpm turbo run check typecheck test build --filter=console`

> **Note, 2026-11-05.** `drift` moved from `platform/cli` to `packages/drift` (ticket 0309) at the same depth, so the `../../..` root climb this ticket removes is unchanged; read its paths under `packages/drift/`.
