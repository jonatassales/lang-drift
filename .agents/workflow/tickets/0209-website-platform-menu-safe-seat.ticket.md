# 0209: Safe Seat replaces Voice AI Interfaces in the Platform menu

**Spec:** .agents/workflow/specs/0043-safe-seat-on-the-website.spec.md

**What to build:** The Platform menu lists Safe Seat where Voice AI Interfaces was, with a seat icon. Its description says it is part of LangDrift FinOps, and it carries a "Coming soon" badge. The item links to `/platform/safe-seat` both on the desktop menu and in the mobile drawer, and the drawer now shows the badge on a linked item. Voice AI Interfaces leaves the repository: its icon, its barrel export, its icon test row and its two header keys.

**Blocked by:** 0208.

**Status:** ready-for-agent

- [ ] `SafeSeatIcon` is in the platform icon family, monochrome in `currentColor` on the CI Gates pattern, and exported from the barrel. The icon test has a Safe Seat row and no Voice row.
- [ ] The menu entry is `safeSeat`, with its description, the icon, `href` set to `/platform/safe-seat` and the `comingSoon` badge.
- [ ] The header catalog replaces the two Voice AI Interfaces keys with `safeSeat` and `safeSeatDescription` in five locales, and `HeaderMessages` follows.
- [ ] The mobile drawer renders the badge on a linked item.
- [ ] The platform menu test checks that Safe Seat links to `/platform/safe-seat`, that "Coming soon" appears twice and that Voice AI Interfaces is gone.
- [ ] Nothing in the repository references Voice AI Interfaces any more, except ticket 0185 and the specs that record its history.

## Dev part

**Override:** no part of this ticket qualifies. It is an icon drawing, a menu entry, a key swap and one conditional badge, with no domain logic, hard part or new stack. Proposed by `to-tickets` on 2026-11-04; the owner may replace it.

## Verification

- `pnpm --filter website test`
- `pnpm turbo run test --filter @langdrift/web`
- `pnpm --filter website typecheck`
- `pnpm turbo run check`
