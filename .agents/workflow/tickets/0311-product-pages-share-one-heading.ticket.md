# 0311: every product page opens with the same heading

**Spec:** None. The owner's bug report of 2026-11-05: the icon and the headline differ in size across the platform product pages; they must match the Microruntimes page.

**What to build:** Microruntimes, LangDrift CLI, CI Gates, FinOps and Safe Seat open with one heading: the product's icon on its tile, sized in em so it scales with the title, then the product's name, at the Microruntimes page's size. FinOps reads "FinOps" and Safe Seat "Safe Seat" with its Coming soon badge; the sentence each used as its title stays under the heading as a tagline, and the home page's FinOps section keeps showing it. One component renders the heading on all five pages.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] One `ProductHeading` component renders the icon tile and the name on the five product pages, with a unit test.
- [x] The h1, the tile and the icon measure the same on the five pages at 1280px and 375px, equal to Microruntimes, with no horizontal overflow.
- [x] The FinOps and Safe Seat sentences move to `hero.taglineStart` and `hero.taglineEnd` in every locale, with key parity.

## Dev part

**Override:** a layout fix; no part qualifies under the Safe Seat policy. 2026-11-05.
