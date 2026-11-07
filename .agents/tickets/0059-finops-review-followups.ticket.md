# 0059: FinOps review follow-ups

**Spec:** .agents/specs/0008-finops-opening.spec.md

**What to build:** the two-axis review of PR #24, run after it merged, found three places where the FinOps opening does not say or expose what spec 0008 asked for. A screen-reader user hears that the receipt follows the ledger, though it now sits beside it or under it. The gate's source line points to the ledger, though ticket 0053 asked for the receipt. And the hero's description never reaches assistive technology, because `aria-describedby` sits on a `div` without a role.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The hero summary says the receipt sits beside the ledger, or under it on a narrow screen, in all five locales.
- [ ] The gate's source line names the receipt at the top of the page as the origin of its figures, in all five locales, keeping the `{item}` placeholder.
- [ ] The hero section, already a region named by the hero heading, is described by the caption and the summary; the plain `div` around the ledger and receipt loses its `aria-describedby`.
- [ ] Vision and Microruntimes do not change. The Microruntimes hero carries the same `aria-describedby` pattern and is left to the branch that is reworking that page.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check localization`.
- `/en/platform/finops/` exposes the hero region with its name and description in the accessibility tree.
