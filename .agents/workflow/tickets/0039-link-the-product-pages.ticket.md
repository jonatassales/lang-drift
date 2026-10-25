# 0039: Link FinOps and Microruntimes from the header and the home page

**Spec:** .agents/workflow/specs/0005-finops-and-microruntimes-pages.spec.md

**What to build:** the Platform dropdown, the drawer, the home hero's product list and the home feature cards open the FinOps and Microruntimes pages in every locale; AI Voice still renders as text.

**Blocked by:** 0035, 0037.

**Status:** ready-for-agent

- [ ] The platform product model carries `href` for FinOps and Microruntimes, which links the dropdown, the drawer and the home hero list through the existing model.
- [ ] `FeaturesGrid` accepts an optional `href` per item and makes the whole card the link, with a visible focus ring; the home features section passes localized paths for the products that have a page.
- [ ] The dropdown comment about items without a page stays true for AI Voice.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check dead-code`.
- From `/en/` and `/pt-br/`: the dropdown, the drawer at 375px, the hero list and the feature cards navigate to both pages; the Vision link still works; AI Voice is text.
