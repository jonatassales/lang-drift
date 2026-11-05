# 0273: the vendored tabs and accordion draw with utilities

**Spec:** .agents/workflow/specs/0038-tailwind-rules-compliance.spec.md

**What to build:** `AnimatedTabs` and `BasicAccordion` keep their markup, ARIA and keyboard model and draw every declaration of their stylesheets as Tailwind utilities and token roles, with `motion-reduce:` turning their transitions off; the scrollbar uses arbitrary properties that reference the border role. Their component stylesheets and imports are deleted.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05. Before: ready-for-agent, delivered with ticket 0299 of spec 0046.

- [x] Both components' tests assert each variant's indicator classes, the selected and focused states and the reduced-motion utilities.
- [x] The affected website pages match their capture at 1280px and 375px, light and dark, recorded under `.audits/`.
- [x] `pnpm turbo run check typecheck test build --filter=web --filter=tokens` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
