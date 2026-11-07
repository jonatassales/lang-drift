# 0292: the Microruntimes session is fine-grained atoms

**Spec:** .agents/specs/0042-state-rules-compliance.spec.md

**What to build:** The session module gains one primitive atom per value (the stage, the played count, whether the visitor paused, the run counter, the card opened by hand, the launched choices); the derived and write atoms keep their names and meaning and read or set the primitives; the session interface and its initial object go. Each consumer subscribes to the atoms it reads, and nothing but the `@state` barrel imports an atom module.

**Blocked by:** 0229.

**Status:** wontfix. The owner closed it on 2026-11-05: it will not be implemented.

- [ ] A unit test drives the write atoms through a Jotai store and reads the derived atoms.
- [ ] The session-cards and session-terminal tests pass.
- [ ] `pnpm turbo run check typecheck test build --filter=website` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
