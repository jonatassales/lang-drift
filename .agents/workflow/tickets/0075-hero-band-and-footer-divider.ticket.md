# 0075: The hero, the closing band and the footer line go

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** The home's hero holds only its moving headline and its two calls to action, at its top left: the black hole sun and the product column on its right are gone in every locale. No page closes with the "Book a demo" band any more: the home, the FinOps page and the Microruntimes page end with their FAQ, while the header's demo button and the Book a demo page stay. The footer loses the line between its two rows. The vendored components, the copy and the tests that only these elements used leave with them.

**Blocked by:** 0073

**Status:** ready-for-agent

- [ ] The hero renders the headline and the two calls to action at its top left, with no sun and no product column, on desktop and on a phone.
- [ ] The home, the FinOps page and the Microruntimes page end with their FAQ, and the header's demo button still opens the Book a demo page.
- [ ] The footer shows no line between its rows, and its other styles are unchanged.
- [ ] The dead-code check, the catalog parity test, typecheck and the tests pass, and the production build generates every page in the five locales.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
