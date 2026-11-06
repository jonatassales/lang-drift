# 0294: the book-a-demo form keeps one responsibility per file

**Spec:** .agents/specs/0042-state-rules-compliance.spec.md

**What to build:** The book-a-demo submit handler posts the validated data with the locale and returns `sent`, `unavailable` or `failed`; a map turns the response into that outcome; the component navigates, shows the notice or the toast on the outcome; the toast dismissal gets its handler.

**Blocked by:** 0229, 0239.

**Status:** wontfix. The owner closed it on 2026-11-05: it will not be implemented.

- [ ] Each user event has a `.handler.ts` and each mutation a `.map.ts`, with tests where they hold logic.
- [ ] `pnpm turbo run check typecheck test build --filter=website` passes.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
