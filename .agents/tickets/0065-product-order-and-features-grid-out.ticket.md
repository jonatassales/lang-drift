# 0065: Products in the platform order; the features grid goes

**Spec:** .agents/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** the Platform menu, the drawer and the hero's product column list Vision, Microruntimes, FinOps and AI Voice in that order, and the home page ends after the hero: the "Four products. One platform." grid is gone with everything that only served it.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The platform products' list is in the order Vision, Microruntimes, FinOps, AI Voice; the menu, the drawer and the hero column follow it without a change of their own.
- [ ] `PlatformFeatures` is deleted and the home page no longer renders it.
- [ ] The vendored Features Grid and the Glow Hover Card leave the web package with their export entries; `motion` stays, since other vendored blocks use it.
- [ ] The platform catalog loses `features` in every locale and its message type loses the key.
- [ ] `drift check dead-code` and `drift check localization` pass.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check dead-code`, `drift check localization`.
- `/en/`: hero, then footer; the hero column, the Platform menu and the drawer read Vision, Microruntimes, FinOps, AI Voice.
