# 0287: the icon family has one owner

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** `packages/web` gains the technology icons the console draws (React, Python, PostgreSQL, Kafka) and the header icons the website draws inline (menu, sun, moon, close). The console's `StackIcon` becomes a lookup from a stack id to the shared icon, falling back to the neutral cube; `stack-icons.json` is deleted. Where both drew a technology, the shared artwork wins. The platform menu takes the shared caret.

**Blocked by:** 0229, 0230.

**Status:** ready-for-agent

- [ ] No app draws an inline `<svg>` icon.
- [ ] Captures differ only in the four technologies that take the shared artwork.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
- `pnpm turbo run check typecheck test build --filter=website`
- `pnpm turbo run check typecheck test build --filter=console`
