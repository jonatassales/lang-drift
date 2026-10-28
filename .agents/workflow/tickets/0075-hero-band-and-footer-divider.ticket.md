# 0075: The hero holds the planning session; the closing band and the footer line go

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** The home's hero holds the Microruntimes planning session: the moving headline at its top left and, under it, one frame that holds the session's terminal and its ticket, plan and workflow cards side by side. The black hole sun, the product column and the two calls to action are gone in every locale, and the home's Microruntimes section goes, since its session now lives in the hero. No page closes with the "Book a demo" band any more: the home, the FinOps page and the Microruntimes page end with their FAQ, while the header's demo button and the Book a demo page stay. The footer loses the line between its two rows, and the home gets the blurred bottom scroll hint the product pages already have. The vendored components, the copy and the tests that only these elements used leave with them.

**Blocked by:** 0073

**Status:** ready-for-agent

- [ ] The hero renders the headline and, under it, one frame with the terminal and the session cards side by side, with no sun, no product column and no calls to action; from `lg` the whole frame fits the first screen and the cards scroll inside it, and on a phone the cards stack under the terminal.
- [ ] The home, the FinOps page and the Microruntimes page end with their FAQ, and the header's demo button still opens the Book a demo page.
- [ ] The footer shows no line between its rows, and its other styles are unchanged.
- [ ] The home shows the bottom scroll hint until its end is on screen.
- [ ] The dead-code check, the catalog parity test, typecheck and the tests pass, and the production build generates every page in the five locales.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
