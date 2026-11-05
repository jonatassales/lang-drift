# 0263: drift's commands declare one function per file

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** A command's local functions move to a folder named after the command (`commands/check-biome/<function>.sh`), one each, and the command sources them by name. A `return` or `exit` that follows a statement gets its blank line. `|| true` stays on a probe and is removed where it hides a failure, decided per site.

**Blocked by:** 0262.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] No shell file under `platform/` declares more than one function.
- [ ] Every `drift` command behaves as before, including the fixture-driven `drift check biome` and `drift quality detectors`.
- [ ] `platform/cli/drift quality check general` prints no finding under `platform/`.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `platform/cli/drift check biome`
- `platform/cli/drift quality detectors`
- `platform/cli/drift quality check general`
