# 0138: Restore the platform nav with dropdown

**What to build:** Restore the Platform navigation dropdown in the website header. The dropdown and the mobile drawer list the platform products: FinOps, Microruntimes and CI Gates. FinOps and Microruntimes link to their respective pages, and CI Gates lists as an upcoming capability with its description.

**Blocked by:** none

**Status:** ready-for-agent

- [x] Add `CiGatesIcon` to `@langdrift/web/ui/icons/platform` and export it.
- [x] Restore `platform-menu.tsx` and integrate `<PlatformMenu />` in `website-header.tsx`.
- [x] Update `platformProducts` in `platform-products.ts` with FinOps, Microruntimes and CI Gates.
- [x] Add `platform`, `platformMenu`, `ciGates`, `finopsDescription`, `microruntimesDescription` and `ciGatesDescription` to `header.json` across all five locales and update `HeaderMessages` type.
- [x] Mobile drawer in `website-header.tsx` lists the platform section with the product icons and labels.
- [x] Verification with typecheck, Biome check, test suite, and drift checks passes.
