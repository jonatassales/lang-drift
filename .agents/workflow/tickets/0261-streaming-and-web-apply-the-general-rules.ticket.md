# 0261: the streaming transport and the web components apply the general rules

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** In `packages/streaming`, `parseStreamEntry` returns a parse result that is either the entry or the rejection with its reason, and `deliver` branches on it, with no `try/catch` around the parser; the strategy file is named after `RedisStreamsServer` and the transport barrel follows. In `packages/web`, the agent orb and the hero shader replace `.catch(() => {})` with a named handler that applies the product's fallback and records the failure in component state the markup exposes; the palette reader's predicate gets a name; the remaining spacing findings close.

**Blocked by:** 0236.

**Status:** ready-for-agent

- [ ] The transport parks an unreadable entry through the parse result, with a test.
- [ ] The orb and the shader expose their fallback in markup, with a test each.
- [ ] `platform/cli/drift quality check general` prints no finding under `packages/streaming` or `packages/web`.
- [ ] `pnpm turbo run check typecheck test --filter=streaming --filter=web` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test --filter=streaming --filter=web`
- `platform/cli/drift quality check general`
