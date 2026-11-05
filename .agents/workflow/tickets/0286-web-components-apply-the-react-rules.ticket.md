# 0286: the web package's components apply the React rules

**Spec:** .agents/workflow/specs/0040-react-rules-compliance.spec.md

**What to build:** In `packages/web`, every function in component scope is a declaration, every props interface sits directly above its component with one blank line, every scope follows the rule's order, and every effect that stays says why. The terminal's in-view, page-visibility and animation loop move into a terminal playback hook, and the terminal reports its progress where it completes a line.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] No arrow function is assigned in a component scope of `packages/web`.
- [ ] Every remaining `useEffect` has its reason on the line above.
- [ ] The terminal's tests pass with the hook in place.
- [ ] `pnpm turbo run check typecheck test build --filter=web --filter=tokens` passes.

## Dev part

**Category:** `hard-part`. Coordinating an intersection observer, page visibility and an animation clock with their cleanup is where this ticket's difficulty sits.

**Paths:**
- `packages/web/src/vendors/smoothui/terminal-playback.hook.ts`
- `packages/web/src/vendors/smoothui/terminal-playback.hook.unit.ts`

**Boundary:** `useTerminalPlayback(machine: TerminalMachine, options: TerminalPlaybackOptions): TerminalPlayback`, called by the Mac terminal, where `TerminalPlaybackOptions` and `TerminalPlayback` are declared by the agent in `terminal-playback.types.ts`.

**Goal:** Play the terminal only while it is in view and the page is visible, advance it on the animation clock, and release every subscription on unmount.

- [ ] Tests with a fake observer, visibility and clock cover play in view, pause out of view, pause on a hidden page, resume, and cleanup on unmount.
- [ ] No subscription survives the unmount.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
