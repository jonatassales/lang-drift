# 0290: the website's components apply the React rules

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** The website header's drawer modal, body scroll lock, focus return and desktop close move into a navigation drawer hook. Functions in component scope are declarations, props interfaces sit above their components, scopes follow the rule's order, and every effect that stays says why.

**Blocked by:** 0239.

**Status:** ready-for-agent

- [x] No effect derives state in the website or the console.
- [x] The drawer hook has a test with a fake dialog element.
- [x] `pnpm turbo run check typecheck test build --filter=website` passes.

> **Note, 2026-11-05.** The drawer hook reads the locale and keeps the one the drawer opened in, so the drawer is open only in that locale and a locale switch closes it with no effect; the console's mobile navigation does the same with its path. The hand-off to the replacement trigger now focuses in the hook's effect instead of on the next animation frame: Next's router leaves focus alone, and a frame never comes in a hidden tab, which is how the hand-off looked broken under 0234. The 53 props interfaces that a comment or a declaration kept from their components now sit directly above them, with the comment above the interface. `global-error` reads its locale from the pathname during render instead of setting it in an effect. Rule 5 does not place effects, so the platform menu's dismissal effect follows the `visible` value it needs; session-terminal's effects read `session.stage` so that `run` and `stage` stay with the derived values.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
