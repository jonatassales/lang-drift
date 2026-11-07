# 0300: Early Access leaves and Book a Demo says we will be back soon

**Spec:** .agents/specs/0046-the-website-focus-pass.spec.md

**What to build:** A visitor finds no early access anywhere, and the Book a Demo page shows a "We'll be back soon" card with a lead and a link to the home. The form's import and element stay in the page as commented code, with a comment naming when they return. The early-access components, API route, confirmation page and image, its catalogs (with their types, loaders, shape validator, barrel and parity entries), its environment and Turnstile documentation, and its tests are deleted. The book-a-demo email schema holds its own email rule. The privacy page describes only what the site collects. The Microruntimes session's demo ticket is renamed to a book-a-demo runtime.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `git grep -n -i 'early.access' apps/website` matches nothing outside history.
- [x] `/book-a-demo` renders the card in every locale and its client bundle holds no zod, so the CSP reports no `eval`.
- [x] The book-a-demo API route and its tests still pass.
- [x] `platform/cli/drift check localization` passes with the two catalogs removed.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
