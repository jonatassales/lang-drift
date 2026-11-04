# Safe Seat on the website

**Status:** ready-for-agent

## Problem Statement

Spec 0033 settled Safe Seat, the LangDrift FinOps harness capability that keeps a person writing part of every ticket by hand, but left the website out of scope. A visitor learns nothing about it.

- The Platform menu lists FinOps, Microruntimes, CI Gates and Voice AI Interfaces. Voice AI Interfaces names a product LangDrift no longer develops.
- No page explains why a FinOps platform would keep a person writing code, nor the evidence behind that.
- The simulated session in the home hero runs a whole ticket through the agent, as if nobody needed to write any of it.

## Solution

**The menu.** Safe Seat takes the place of Voice AI Interfaces in the Platform menu. Like CI Gates, it carries a "Coming soon" badge. Unlike CI Gates, it links to its own page, on the desktop menu and in the mobile drawer, and the drawer shows the badge on a linked item too.

**The page.** `/platform/safe-seat` explains the capability in the five locales. It is described as part of LangDrift FinOps and marked "Coming soon". It argues the case in this order:

- why writing code is a way of thinking;
- what the research and aviation practice show, each finding with its limit and a link to its primary source;
- what stays with the person;
- how the capability will work;
- the principles that keep it from becoming surveillance;
- what the evidence does not show;
- the sources and an FAQ.

The page paraphrases the ideas of a post by a colleague of the owner, never quoting it or naming its author. It makes no clinical claim. Burnout appears only as the association the research reports.

**The home demo.** The simulated session in the home hero, the same one the Microruntimes page plays, gains the dev ticket:

- the agent implements up to the reserved part and stops;
- the terminal waits while the person writes it, then shows it closed "without an agent", with no paid call;
- the agent's code review runs after the person;
- the end of the run shows "Safe Seat met", never a share;
- the Workflow card shows the dev ticket as a stage the person drives, apart from the stages a person approves.

## User Stories

1. As a visitor, I want the Platform menu to list Safe Seat instead of Voice AI Interfaces, so that the menu names what LangDrift is building.
2. As a visitor, I want the Safe Seat item to say it is part of LangDrift FinOps, so that I do not take it for a separate product.
3. As a visitor, I want the Safe Seat item to carry a "Coming soon" badge, so that I know it is not available yet.
4. As a visitor, I want the Safe Seat item to open its page, so that I can read about it before it ships.
5. As a visitor on a phone, I want the drawer to show the Safe Seat link with its "Coming soon" badge, so that the phone and the desktop agree.
6. As a visitor, I want the page to open by saying what Safe Seat is and that it is coming soon, so that I know where I am.
7. As a developer, I want the page to argue that writing code is a way of thinking, so that the capability reads as craft and not as a quota.
8. As an engineering lead, I want each piece of evidence stated with its limit and linked to its primary source, so that I can check the claims myself.
9. As an engineering lead, I want to see which kinds of problems stay with the person and which never qualify, so that I understand what the practice asks.
10. As an engineering lead, I want the planned flow in five steps, so that I see how the reserved part, the dev ticket, the measurement and the override fit together.
11. As a developer, I want the page to promise that my numbers stay mine and never feed a ranking or a performance review, so that the practice does not feel like surveillance.
12. As a reader, I want the page to say what the evidence does not show, including any clinical effect, so that it does not oversell.
13. As a reader, I want every source cited in full with a link, so that I can read it.
14. As a visitor, I want an FAQ that answers whether it is a quota, who sees my share, whether I can use AI in the dev ticket, why 5%, whether it is available and whether it sees tools outside `ex`, so that the common doubts have answers.
15. As a visitor in any of the five locales, I want the page in my language with the name Safe Seat untranslated and never abbreviated, so that the capability reads the same everywhere.
16. As a search engine, I want the page to carry its own title, description, canonical, hreflang alternates and sitemap entries, so that it can be indexed like the other product pages.
17. As a visitor watching the home hero, I want the simulated run to stop at a dev ticket that the person writes, so that I see Safe Seat at work.
18. As a visitor, I want the dev ticket to close "written without an agent", with no paid call, so that the money story stays true.
19. As a visitor, I want the agent's code review to run after the person, so that the session matches the planned flow.
20. As a visitor, I want the run's end to show "Safe Seat met" and no percentage, so that the session respects the privacy the page promises.
21. As a visitor, I want the Workflow card to mark the dev ticket as driven by the person, unlike a stage a person approves, so that the two never blur.
22. As a visitor with reduced motion, I want the whole run, dev ticket included, printed at once, so that I lose nothing.
23. As a visitor using a screen reader, I want the hero summary to mention the dev ticket and still say the session is simulated, so that I get the same story.
24. As the `ex` maintainer, I want the new run lines recorded as part of the contract that ticket 0202 makes `ex` print, so that the website and the CLI do not drift apart.

## Implementation Decisions

### The menu

- A `SafeSeatIcon` joins the platform icon family in `@langdrift/web/ui/icons/platform`: a monochrome seat drawn in `currentColor`, on the CI Gates icon pattern.
- The menu entry is `{ label: 'safeSeat', description: 'safeSeatDescription', icon: SafeSeatIcon, href: '/platform/safe-seat', badge: 'comingSoon' }`. The header catalog replaces the two Voice AI Interfaces keys with these two, in every locale.
- Voice AI Interfaces leaves the repository: the icon, its barrel export, its test row and the two keys. The menu was its only consumer.
- The desktop menu already renders a badge and a link together. The mobile drawer's linked branch gains the badge.

### The page

- **Route and boilerplate.** `[locale]/(pages)/platform/safe-seat/` follows the Microruntimes page: an async Server Component, the structured data, the skip link, header, main and footer. Metadata comes from a `page-metadata.map.ts` that calls `pageMetadata('safe-seat', …)`. The page's components, its ids and links (`safe-seat-content.data.ts`) and its catalog types (`safe-seat-messages.types.ts`) sit in its folder, not under `lib/`, as specs 0039 and 0041 place new route files. The page inherits the locale's social images and has no call to action, like the other product pages.
- **Hero.** An eyebrow reading "LangDrift FinOps", the "Coming soon" badge, a two-part `h1` and a lead.
- **Sections.** Built with the product-page primitives:
  - `#why`;
  - `#evidence`, with six cards: Bainbridge 1983, Casner et al. 2014, Shen and Tamkin 2026, METR 2025 and 2026, Lee et al. 2025, Feng et al. 2025;
  - `#reserved`, the three categories and the excluded kinds;
  - `#flow`, five planned steps;
  - `#principles`;
  - `#limits`;
  - `#sources`;
  - the product FAQ.
- **Copy and data.**
  - One catalog `safe-seat.json` holds the copy for the five locales in `en`'s key order, including the citation text.
  - A content module holds the id lists, the source links and the `SafeSeatMessages` type, as the FinOps content module does.
  - The catalog is registered wherever the FinOps catalog is: the loaders, the shape validator, the messages barrel and the catalog parity test.
- **Metadata.**
  - `pages.json` gains `safe-seat`, with a unique title and description per locale.
  - The nested page path maps it to `platform/safe-seat/`.
  - The sitemap gains five entries, one per locale, each with its alternates.
- **Copy rules.** Every rule below holds in all five locales:
  - "Safe Seat" is never translated and never abbreviated;
  - no clinical or health vocabulary;
  - burnout is an association, not an effect;
  - the 5% floor is a product principle, not a number taken from research;
  - findings are paraphrased, and any quote is a few words long;
  - the colleague's post is paraphrased, never quoted or attributed.

### The home demo

- **The workflow data.**
  - The run's stages become planning, spec, tickets, implement, dev-ticket, review and merge.
  - `dev-ticket` is a stage the person drives. Like `merge`, it is not one of the stages a person can switch on or off.
  - Implement keeps `implement` and `tdd`, and prints a line saying it stops at the reserved part and writes only the boundary.
  - Review gains `code-review`.
- **New run messages.**
  - `devTicketWaiting`: "waiting dev-ticket · the person writes the reserved part".
  - `devTicketClosed`: "closed dev-ticket · written without an agent · no paid call".
  - `safeSeatMet`: the Safe Seat state, "met".
- **Order.** The dev-ticket lines follow implement, and `safeSeatMet` follows the checks. The gate-pause branch, the plan session and every money figure stay as they are.
- **Stage progress.** It covers every stage, so the person-driven row can go active and then done.
- **The Workflow card.**
  - A fixed row for the dev ticket, marked as driven by the person, beside the merge row.
  - A Safe Seat cell in the result, shown once `safeSeatMet` has printed.
- **Copy.** The run lines keep their `console ›` English form. The `#` comments and the card copy are translated, and the hero summary still says that every step is simulated.
- **The `ex` contract.** Spec 0016 makes the website session the contract for `ex`. Ticket 0202 makes `ex` print the new dev-ticket lines and mirror them in its session catalog.

## Testing Decisions

A good test checks what a visitor sees or what a pure function returns, never markup structure.

1. **Page content** (Vitest, beside the FinOps content test):
   - every id list matches its catalog keys;
   - every source link is an absolute `https` URL.
2. **Catalogs.**
   - The catalog parity test covers the new family and the changed `header` and `microruntimes` catalogs: same keys, ICU arguments and tags in every locale, and the counts 0, 1, 2, 3 and 12 formatted.
   - The shape validator checks key order.
3. **SEO.** The SEO contract and structured-data tests read `pages.json`. They pass only with the five sitemap entries and unique copy.
4. **Menu.** The platform menu test:
   - Safe Seat links to `/platform/safe-seat`;
   - "Coming soon" appears twice;
   - Voice AI Interfaces is gone.

   The icon test gains the Safe Seat row and loses the Voice row.
5. **Demo.**
   - The workflow test: the dev-ticket lines follow implement, the code review follows the dev ticket, `safeSeatMet` precedes `ready`, and the pause branch is unchanged.
   - The session cards test: the person-driven row and the Safe Seat cell.
   - The plan's duration and line count stay untouched.
6. **Repository checks:**
   - `pnpm --filter website test` and `typecheck`;
   - the `@langdrift/web` tests;
   - `pnpm turbo run check`;
   - `drift check localization`;
   - `drift quality check` with no new finding in the changed files.

Manual verification with `pnpm website`:
- the page in `en` and `pt-br`, at 375px, in light and dark;
- the menu on desktop and in the drawer;
- the home hero run through the dev ticket, with and without reduced motion;
- the canonical, the hreflang and the 35 sitemap entries.

## Out of Scope

- `ex` printing the new run lines, which belongs to ticket 0202.
- A Safe Seat demo on the page itself, and any change to the plan side of the session.
- A home page section or FAQ tab for Safe Seat.
- A Safe Seat social image of its own.
- Any change to CI Gates or to the pricing and Book a demo pages.

## Further Notes

- The owner settled these decisions in a grilling on 2026-11-04 that ended with an approved plan.
- The evidence and its limits are in the research note `safe-seat.research.md`. The page cites only what the note supports.
