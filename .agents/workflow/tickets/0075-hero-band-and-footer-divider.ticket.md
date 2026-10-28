# 0075: The hero holds the planning session; the closing band and the footer line go

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** The home's hero holds the Microruntimes planning session: the moving headline alone, centered on the first screen over a black-and-white shader with a soft rainbow glow behind it that the header's demo button gives up, and, under it, one frame that holds the session's terminal and its ticket, plan and workflow cards side by side, with only its top showing at the bottom of the first screen. The header is transparent at the top of the page and runs over the hero's background; it turns to glass once the page scrolls. The black hole sun, the product column and the two calls to action are gone in every locale, and the home's Microruntimes section goes, since its session now lives in the hero. No page closes with the "Book a demo" band any more: the home, the FinOps page and the Microruntimes page end with their FAQ, while the header's demo button and the Book a demo page stay. The footer loses the line between its two rows, and the home gets the blurred bottom scroll hint the product pages already have. The vendored components, the copy and the tests that only these elements used leave with them.

**Blocked by:** 0073

**Status:** ready-for-agent

- [ ] The hero renders the headline centered on the first screen over a grayscale shader, and under it one frame with the terminal and the session cards side by side, whose top peeks at the bottom of the first screen; there is no sun, no product column and no call to action.
- [ ] From `lg` the frame grows with the cards so they never scroll, and the terminal matches their height; on a phone the cards stack under the terminal.
- [ ] The header has no background at the top of the page and blends with the hero; it shows its glass once the page scrolls.
- [ ] The home, the FinOps page and the Microruntimes page end with their FAQ, and the header's demo button still opens the Book a demo page.
- [ ] The footer shows no line between its rows, and its other styles are unchanged.
- [ ] The home shows the bottom scroll hint until its end is on screen.
- [ ] The dead-code check, the catalog parity test, typecheck and the tests pass, and the production build generates every page in the five locales.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
