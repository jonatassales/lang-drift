# The site becomes one agentic FinOps platform page

**Status:** ready-for-agent

## Problem Statement

The owner pivoted LangDrift to an agentic FinOps platform with Microruntimes as its engine, and Vision stopped being a product (spec 0010). The website still presents a platform of products. Its header opens a Platform dropdown that lists Vision, Microruntimes, FinOps and AI Voice. The home page opens with a hero whose product column and rotating words name all four, then shows a Vision section with the Product Vision curve and the Sinapsi graph, an AI Voice section with a scripted sample, and a FAQ with a Vision tab. A Vision page of its own holds the longest demos. AI Voice was never built: it is a menu item with no destination and a home section with a scripted preview. A visitor reads a story the owner no longer tells.

## Solution

The site tells one story: an agentic FinOps platform, with Microruntimes as its engine. The Vision page, the Vision and AI Voice sections of the home, every Vision demo and every AI Voice preview are gone, along with the data, copy in the five locales, styles, images and dependencies that only they used. The Platform dropdown is gone: the header links straight to the FinOps and Microruntimes pages, which stay at their addresses. The home page is the hero, the FinOps section and a FAQ with a FinOps tab and a Microruntimes tab. The hero holds the Microruntimes planning session: its moving headline at the top left and, under it, one frame that holds the session's terminal and cards side by side, without the black hole sun, the product column or calls to action, and no page closes with the "Book a demo" band; the header's demo button and the Book a demo page remain the way to ask for one. A removed address answers "not found", as any unknown address does.

Pricing, early access, Book a demo and privacy stay as they are. Only the words that named Vision change: the link that led to the Vision FAQ now leads to the home's FAQ, and the copy that listed Vision beside the other products lists the two that remain.

The generic UI that served only Vision and AI Voice leaves the shared web package too: the Product Vision curve, the target distribution chart, the Vision and AI Voice icons, the avatar set, two vendored SmoothUI components and the charting dependency. The shared meaning colors for classification and drift stay in the tokens, because the console and other components use them.

## User Stories

1. As a visitor, I want the home page to describe an agentic FinOps platform, so that the site says what the product is today.
2. As a visitor, I want the header to link directly to FinOps and Microruntimes, so that I reach either page without opening a menu.
3. As a visitor on a phone, I want the navigation drawer to list the same two product links, so that the phone and the desktop agree.
4. As a visitor, I want the hero's rotating words to name only things the platform does, so that nothing promises a voice product that does not exist.
5. As a visitor, I want the hero to play the Microruntimes planning session under the headline, so that the first screen shows the platform at work.
6. As a visitor, I want the hero's main call to action to open the FinOps page, so that the first click explains the platform's purpose.
7. As a visitor, I want the home page to show the Microruntimes section and then the FinOps section, so that I see the engine and the accountability it serves.
8. As a visitor, I want the home page and the product pages to end with their FAQ, and the header's demo button to stay in reach, so that every page closes on answers and a demo is one click away.
9. As a visitor, I want the FAQ to switch between FinOps and Microruntimes questions, so that I find the answer for the product I am reading about.
10. As a visitor, I want the footer's FAQ link to take me to the home's FAQ, so that the link never lands on a page that no longer exists.
11. As a visitor on the pricing page, I want the FAQ link under the plans to take me to the home's FAQ, so that I can read the answers before I choose.
12. As a visitor, I want no menu item that leads nowhere, so that every item I can see is something I can use.
13. As a visitor who follows an old Vision address, I want the "not found" page in my language, so that I know the page is gone and can go home.
14. As a visitor, I want the FinOps and Microruntimes pages to stay at their addresses, so that bookmarks and shared links keep working.
15. As a visitor, I want the Book a demo, pricing and privacy pages to work as before, so that the cleanup removes nothing I still need.
16. As a visitor in any of the five locales, I want the same removals and the same navigation, so that no language still shows Vision or AI Voice.
17. As a visitor reading the Book a demo page, I want its copy to name Microruntimes and FinOps only, so that the offer matches the site.
18. As a search engine, I want the sitemap to list only pages that exist, so that I do not index a removed page.
19. As a search engine, I want each remaining page to keep its own title, description, canonical and alternate languages, so that removing a page does not weaken the others.
20. As a search engine, I want the home's description to describe the agentic FinOps platform, so that results match the page.
21. As the owner, I want no style rule, image, data file, translation key or dependency that only Vision or AI Voice used, so that dead weight does not hide in the repository.
22. As the owner, I want the shared styles that the footer, the focus ring and the plans use to survive the removal, so that nothing unrelated changes appearance.
23. As the owner, I want the classification and drift color tokens to stay, so that the console and the other products keep one meaning for a color.
24. As the owner, I want the agent orb and its dependency to stay, so that the footer keeps its mark.
25. As a developer, I want the website's catalogs to stay in parity across the five locales after the removal, so that a missing key fails a test and not a visitor.
26. As a developer, I want the checks for dead code, localization and history to pass, so that the cleanup is proved and not assumed.
27. As a developer, I want the shared web package to export only what an app uses, so that its public surface matches its consumers.
28. As a developer, I want the lockfile to lose the removed dependencies and nothing else, so that the change is easy to review.

## Implementation Decisions

- **Pages.** Pages that stay: home, FinOps, Microruntimes, pricing, Book a demo, privacy and the early-access confirmation. The Vision page is removed. Page identifiers, metadata entries and the hand-written sitemap lose Vision. FinOps and Microruntimes keep their nested addresses, so the metadata's address mapping keeps those two and drops Vision.
- **Removed addresses.** No redirects. The proxy and the framework configuration do not change; an unknown path already renders the localized "not found" page.
- **Home composition.** The hero, the FinOps section and the platform FAQ. The Microruntimes section leaves the home because its session moved into the hero. The home's section frame accepts only the two remaining products. The FAQ keeps its tabbed form with two tabs, named as the header names the products.
- **Hero.** The rotating words lose the voice word. The black hole sun, the product column and both calls to action are removed, with the vendored components only they used and the copy only they read. The session terminal and cards share one frame under the headline, which fits the first screen from `lg`; the session's summary describes the hero for screen readers. Copy in the five locales is edited only where it names a removed product, and the owner reviews it in the pull request. The home keeps "Agentic Platform" as its headline and its title.
- **Header.** The dropdown component is removed. The desktop navigation and the drawer render the two product links beside Docs, GitHub and Pricing, which are unchanged. The product list keeps its two entries and the helper that returns one entry.
- **Book a demo band.** Removed from the home and from the FinOps and Microruntimes pages, with its component, its copy and the FinOps sentence that only it read. The Book a demo page and the header's demo button stay.
- **Footer divider.** The line between the footer's two rows is removed; the rows keep their spacing.
- **Footer and plans.** The footer's description moves into the footer's own catalog, because the home catalog goes away and the footer was its last reader. The FAQ link in the footer and the one under the plans point to the home's FAQ anchor. The footer link's label is renamed if it still talks about Vision.
- **Catalogs.** Removed whole: the home, demo, integrations, Vision loop and executive review catalogs. Trimmed: the header catalog (platform menu, Vision and AI Voice entries, dropdown descriptions), the platform catalog (Vision section, AI Voice product, voice word, Vision tab), the conversion catalog (ROI and the Vision FAQ; the plans stay) and the Book a demo catalog (product list in its lead and band). The provider no longer receives the home and demo namespaces. The type definitions, the shape assertions and the catalog parity test drop the removed families and the two assertions that read their copy.
- **Shared styles.** The responsive stylesheet is removed. The conversion stylesheet keeps only the plans rules. The global and hero stylesheets lose the Vision rules and keep the reset, the focus ring, the section kicker, the footer and the scroll padding. The screen-reader-only class goes too: only the removed Vision components carried it. The layout drops the imports of the removed stylesheets. The temporary-CSS exception can only shrink, so no new rule, file or media query appears.
- **Assets and state.** The avatar and integration logo folders, the three hero-selection atoms and their exports, and the charting-graph dependency of the website are removed.
- **Shared web package.** Removed with their exports and style imports: the Product Vision curve, the target distribution chart, the classification text helper, the avatar set, the Vision and AI Voice icons, two vendored SmoothUI components (the avatar group and the chroma text) and the chroma text rules in the SmoothUI stylesheet, plus the charting dependency and its dead-code allow entry. The drift sign helper and its attribute rules go only if the dead-code check reports them. Kept: the agent orb and its voice dependency, tabs, accordion, brand, info cards, detail list, the language switcher and every token.
- **Docs inside the website.** The website README stops describing the Vision demos and the voice preview.
- **No change in behavior elsewhere.** FinOps and Microruntimes pages, Book a demo, pricing, early access, privacy, the proxy, locale routing and every API route keep their behavior.

## Testing Decisions

- A good test here checks what a visitor or a build can observe: which catalogs and keys exist in every locale, which links a component renders, which pages the metadata lists, and whether the repository's checks pass. It does not assert how a component is built.
- **Seams.** The existing seams are enough: the catalog parity test (keys, ICU arguments, rich-text tags and formatting in the five locales), the component unit tests under the website, the repository's dead-code, localization and history checks, and the production build, whose static generation fails on a message that is missing in any locale. No new seam is needed.
- **Tests that change.** The catalog parity test drops the removed families and the Vision and voice assertions. The home product section test and the FAQ tabs test stop using Vision and AI Voice and cover the two remaining products. The voice preview test goes with the component. The proxy test uses the pricing path only as an arbitrary example and needs no change.
- **Tests that must keep passing unchanged.** The FinOps and Microruntimes component and data tests, the Book a demo tests and the early-access route test.
- **Residue.** A search for the whole words Vision and AI Voice across the website and the shared web package finds only the shared color tokens' names and unrelated words such as revision.
- **Prior art.** The catalog parity test and the product section tests added under spec 0009.

## Out of Scope

- The docs app: its Vision page, menu entry, home and Platform page.
- The extub stub's Vision product and fixtures.
- AGENTS.md, the context file's glossary and navigation text, and the older specs and tickets that describe what is removed. This spec supersedes them in practice.
- Rule files. The CSS exception list in the frontend rule keeps naming files that no longer exist, and the localization rule that keeps "Product Vision" in English loses its subject; the owner edits both when ready.
- Pricing, early access, Book a demo and privacy beyond the words that named Vision.
- Folding the FinOps and Microruntimes pages into the home and deleting them.
- Rewriting the hero's wording to the owner's final positioning.
- Redirects from removed addresses.

## Further Notes

- The Vision page imports nearly every Vision component, so its removal is one atomic change; the evidence sections are orphaned the moment the page goes.
- The repository holds hundreds of stale copies under the worktrees folder of the Claude configuration. Searches must ignore it.
- ADR 0003 (the website adapts to the repository) still holds; this spec only removes content.
- Related: spec 0010 (the pivot), spec 0009 (the home sections this spec trims), spec 0004 (the landing hero and the Vision page this spec removes).
