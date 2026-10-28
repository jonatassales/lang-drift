# 0073: The header links straight to FinOps and Microruntimes

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** The Platform dropdown is gone. On desktop the header navigation and on a phone the drawer list a direct link to the FinOps page and one to the Microruntimes page, beside Docs, GitHub and Pricing, which do not change. The dropdown's unused copy leaves the header catalog in the five locales. The home's product section frame and the product list accept only the two remaining products.

**Blocked by:** 0072

**Status:** ready-for-agent

- [ ] The desktop header and the drawer each show a FinOps link and a Microruntimes link that open those pages, and no dropdown exists.
- [ ] The header catalog has no key that only the dropdown used, and the parity test passes in the five locales.
- [ ] The home product section test and the other website unit tests pass.
- [ ] In a browser, the header, the drawer, the language pill and the demo button still work.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
