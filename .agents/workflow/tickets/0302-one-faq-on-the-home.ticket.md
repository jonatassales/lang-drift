# 0302: one FAQ, on the home

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** The home's FAQ has three tabs, Microruntimes, FinOps and Safe Seat, each reading its product catalog's `faq.items`. The Microruntimes, FinOps and Safe Seat pages no longer render a FAQ, and the page-level FAQ title and description keys nothing reads are removed. Links to the FAQ point to `/#faq`.

**Blocked by:** 0300.

**Status:** ready-for-agent

- [x] Only the home renders a `#faq` section, with three tabs.
- [x] No component left unused, and no catalog key left unread, by the removal.
- [x] `platform/cli/drift check localization` passes.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
