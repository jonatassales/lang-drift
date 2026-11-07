# 0293: the early-access form updates its state without Immer

**Spec:** .agents/specs/0042-state-rules-compliance.spec.md

**What to build:** The form's updater takes a function from the current state to the next; each map returns the next state with a spread, and the handlers call the maps through the updater. The focus event's logic moves to its handler file. `immer` leaves the website.

**Blocked by:** 0229, 0239.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] The website manifest no longer lists `immer`.
- [ ] The handler tests apply the updater to a local state and pass.
- [ ] `pnpm turbo run check typecheck test build --filter=website` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
