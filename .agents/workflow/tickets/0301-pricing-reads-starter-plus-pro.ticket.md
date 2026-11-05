# 0301: pricing reads Starter, Plus and Pro, coming soon

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** The pricing page reads "Plans" with a "Coming soon" badge and shows three cards, Starter, Plus and Pro, each valued "To be announced" with no promotion badge. Starter covers the LangDrift CLI, budget-gated runs, the ledger and evidence; Plus adds CI Gates on every pull request, team budgets and the console; Pro adds Safe Seat, deeper reporting and multi-team envelopes. The introduction, availability line and call to action are rewritten, the call to action goes to Book a Demo, and the page's metadata description follows.

**Blocked by:** 0300.

**Status:** ready-for-agent

- [x] Three cards render in every locale, each with "To be announced".
- [x] No "trial", "15-day", "early adopter" or "early access" text remains in the conversion catalog or the pricing metadata.
- [x] The eyebrow reads "Plans" followed by the Coming soon badge.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
