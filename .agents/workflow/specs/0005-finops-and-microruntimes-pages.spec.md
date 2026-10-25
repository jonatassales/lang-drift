# FinOps and Microruntimes product pages

**Status:** ready-for-agent

## Problem Statement

The platform site names four products in the Platform menu, the drawer, the home hero list and the home cards, but only Vision has a page. A visitor who came for FinOps or Microruntimes finds plain text with nothing to open, and a search engine finds nothing to index. The material that explains both products exists, in Exphere's knowledge base and documentation, written in a voice that keeps the intended model apart from what is implemented; the site carries none of it.

## Solution

Two product pages, `/[locale]/platform/finops/` and `/[locale]/platform/microruntimes/`, that follow the Vision page's layout: the same header, footer and demo band; a hero with a question headline, a lead, two calls to action and an illustrative interactive demo; then sections with a kicker, a heading, a lead and cards; a FAQ; the demo band.

- **FinOps** presents the team budget model: the envelope split into spent, reserved, protected and available; work items with their estimate, reservation and observed spend; a delivery receipt for the selected item; the chain from intent to outcome in three scenarios; the budget that belongs to the team; recommendation that is not authorization; a weekly delivery and cost report; the sources that feed the ledger.
- **Microruntimes** presents the lifecycle (configure → validate → prepare → execute → inspect → clean up) around a runtime summary with desired and observed state; three scenarios from compose to clean up; the annotated runtime file; the ownership split between consumer and platform; requested, observed and unknown side by side; the execution ledger after cleanup; the console-to-connected-agent flow.

Every preview says it is illustrative, an unknown value is shown as unknown rather than zero, and no page claims a budget gate, a ledger, savings or an integration that does not exist. The Platform menu, the drawer, the home hero list and the home cards link to both pages.

## User Stories

1. As a visitor, I want a FinOps page at `/platform/finops/`, so that I can learn what LangDrift FinOps does before asking for a demo.
2. As a visitor, I want a Microruntimes page at `/platform/microruntimes/`, so that I can learn what a microruntime is and what running one involves.
3. As a visitor, I want both pages to share the Vision page's header, footer, section rhythm and demo band, so that the platform reads as one site.
4. As a visitor, I want the FinOps hero to show an illustrative team budget ledger, so that I see at a glance how an envelope splits into spent, reserved, protected and available.
5. As a visitor, I want to select a work item in the ledger and see its delivery receipt further down the page, so that I understand what one piece of work consumed and whether it was verified and cleaned up.
6. As a visitor, I want a work item whose spend is unknown to say so, so that a missing charge never looks like zero.
7. As a visitor, I want three chain scenarios, an allocation within the envelope, a pause and replan, and an unknown after a timeout, so that I see how intent, delivery, evidence and outcome connect even when things go wrong.
8. As a visitor, I want the FinOps page to say that the budget belongs to the team and to explain envelope, allocation and reservation, so that I do not read it as per-developer quotas.
9. As a visitor, I want the page to state that an agent recommends and a deterministic control decides, so that I know a model never widens funding, scope or access.
10. As a visitor, I want a weekly delivery and cost report preview, so that I see what leadership would read.
11. As a visitor, I want to know what feeds the ledger and that these sources are the intended model, so that I do not mistake the model for a shipped product.
12. As a visitor, I want a FinOps FAQ, so that the questions I would ask in a demo are answered on the page.
13. As a visitor, I want the Microruntimes hero to show the six lifecycle steps and a runtime summary with desired and observed state, so that I see what a bounded run looks like.
14. As a visitor, I want to select a lifecycle step and see what it does, what evidence it leaves and its record in the evidence section, so that I understand the lifecycle end to end.
15. As a visitor, I want three lifecycle scenarios, adding a component, a failed preparation and an agent that loses its lease, so that I see how failures and unknowns stay visible.
16. As a visitor, I want to read an annotated runtime file, so that I know the configuration lives in one file and never copies my source.
17. As a visitor, I want an ownership table, so that I know what I own and what LangDrift owns.
18. As a visitor, I want requested, observed and unknown state side by side, so that I understand why a request is not a confirmation.
19. As a visitor, I want the execution ledger preview, so that I know cleanup leaves evidence and reports leftovers.
20. As a visitor, I want to see that the console authors and the connected agent executes, so that I know nothing starts because a page was opened.
21. As a visitor, I want a Microruntimes FAQ, so that the questions about copying code, supported stacks and sandboxing are answered on the page.
22. As a visitor, I want the Platform menu, the drawer, the home hero list and the home cards to open the two pages, so that I can reach them from anywhere on the site.
23. As a visitor on a narrow screen, I want every section to stack and the runtime file to scroll sideways, so that nothing overflows at 375px.
24. As a keyboard user, I want the ledger items, the stepper and the scenario tabs reachable with Tab and operable with Enter, Space and the arrow keys, and the receipt and evidence sections announced when they change, so that the demos work without a pointer.
25. As a visitor who asked for reduced motion, I want the ledger bar, the stepper and the tabs still, so that nothing moves without my consent.
26. As a visitor reading pt-BR, zh-Hant, ja or vi, I want both pages in my language with the same illustrative labels, so that nothing depends on English.
27. As a visitor in either colour scheme, I want both pages to read correctly in light and dark, so that the theme toggle never breaks a section.
28. As a search engine, I want each page to carry its own localized title and description, a self canonical, hreflang, Open Graph and Twitter metadata and a sitemap entry, so that both index correctly.
29. As a frontend developer, I want the new sections built from Tailwind role utilities only, so that the frozen stylesheets keep shrinking.
30. As a frontend developer, I want the generic pieces, a section heading, a detail list and info cards, owned by `@langdrift/web`, so that a future product page composes them.
31. As a frontend developer, I want every new string in JSON catalogs with parity across the five locales, so that the localization checks stay green.
32. As a frontend developer, I want each hero demo and its mirrored section to share one atom, so that selection stays in sync without prop drilling across the page.
33. As the owner, I want the copy drafted in all five locales, so that I edit the catalogs rather than start from nothing.
34. As the owner, I want every claim on the pages to match the maturity the Exphere material supports, so that the site never promises a budget gate, savings or integrations that do not exist.

## Implementation Decisions

- **Routes.** Two App Router pages under the `(pages)` group, `platform/finops` and `platform/microruntimes`, with the Vision page's boilerplate: locale validation, request locale, metadata through the page-metadata helper with the inherited social images, skip link, header, `main`, footer. Fragment calls to action stay plain anchors.
- **Copy.** Two catalogs, `finops` and `microruntimes`, five locales in one file each under top-level locale keys, in the same key order as `en` because the catalog shape check joins keys in order. They are read on the server and handed to client demos as props; they do not join the client messages. The highlighted headline is split into start, highlight and end keys and the highlight is wrapped in `ChromaText`, since rich formatting is unavailable to a server-only catalog. The only plural is the receipt's retries; amounts and percentages arrive pre-formatted. Product names, `ex`, `runtime.md`, `cpuMillis`, `memoryMiB`, kind, Terraform and Kubernetes stay untranslated inside translated sentences; the illustrative disclosures appear in every locale.
- **Data.** A content module per page holds what is not a sentence: ids, amounts, counts, enums, paths, the runtime excerpt and the initial selection, plus the catalog's TypeScript shape typed from the id arrays, as the vision-loop content module does. Sentences stay in the catalogs.
- **Shared primitives in `@langdrift/web`.** `SectionHeading` (kicker, heading with id, optional lead), `DetailList` (a `dl` of term and detail rows with an optional status tone) and `InfoCards` (a grid of cards with eyebrow, title, description, optional icon and link, ordered or not). Each is a Server Component, Tailwind-only, exported by name from its own entry in the package export map, and lands in the commit that first uses it so the dead-code check stays green.
- **Shared page shell in the website.** `ProductHero` (the Vision hero layout in Tailwind with a demo slot, caption and screen-reader summary), `ProductSection` (the section shell with scroll margin, border and `SectionHeading`), `PillarScenarios` (a client scenario-tabs demo over `AnimatedTabs` with pillar columns, fed by the page), `ProductFaqSection` (a Tailwind-only FAQ over `BasicAccordion`; the Vision FAQ stays as it is).
- **FinOps demos.** `BudgetLedgerDemo` (client) draws the envelope as a stacked bar with a legend and lists the work items as toggle buttons; `DeliveryReceipt` (client) renders the selected item's receipt in a live region. Both share `selectedWorkItemAtom`.
- **Microruntimes demos.** `LifecycleDemo` (client) draws the stepper as buttons with `aria-current="step"`, the runtime summary and the selected step's panel; `StageEvidence` (client) renders the selected step's record. Both share `selectedLifecycleStageAtom`. `RuntimeDescriptorCard` and `AuthoringFlow` are Server Components.
- **State.** One primitive atom per page, in its own atom module and re-exported by name from the state barrel. Scenario tabs keep local state because nothing else consumes them.
- **Illustrative values.** Unknown spend and unknown observed state render copy in the warning tone, never a number. Currency is formatted with `Intl.NumberFormat` for the locale, in USD, and the caption says the amounts are simulated.
- **Accessibility.** Work items and lifecycle steps are real buttons (`aria-pressed`, `aria-current`); the receipt, the evidence record and the stage panel are live regions; the stacked bar is an image with a label and a visible legend; the ownership matrix is a real table; tab panels carry the ids `AnimatedTabs` emits; transitions are `motion-safe` only.
- **Metadata.** `pages.json` gains `finops` and `microruntimes`; the nested-path map gains both routes with trailing slashes; `sitemap.xml` gains ten entries with hreflang alternates, thirty in total.
- **Links.** The platform product model gains an `href` for FinOps and Microruntimes, which links the dropdown, the drawer and the home hero list. The home cards need `FeaturesGrid` to accept an optional `href` per item; the home features section passes localized paths. AI Voice stays text.
- **Workflow.** One commit per ticket on `website-product-redesign`, each with its `Ticket:` trailer.

## Testing Decisions

A good test checks what a visitor or a build sees: rendered text, ARIA state, the store value after a click, a catalog's parity. It never asserts class names or internal structure.

Seams, highest first:

1. **Catalog parity** (existing seam): both catalogs join the family list of the website's catalog test, which checks keys, ICU arguments and rich tags against `en` and formats every message with counts 0, 1, 2, 3 and 12 in every locale.
2. **Content modules** (new tests at an existing kind of seam): the FinOps envelope segments sum to the authorized amount and exactly one work item has unknown spend; the six lifecycle stages are in order, desired differs from observed, and the runtime excerpt stays within the `ex` limits for components, millicores and MiB.
3. **Web primitives** (existing seam): `renderToString` tests like the hero-shader test: the heading id and kicker order, the tone mapping of a detail row, ordered versus unordered cards and the column count.
4. **Demos** (existing seam): jsdom tests with Testing Library like the voice-preview test, wrapped in a Jotai `Provider` with a fresh store: clicking the second work item puts its id in the store and the receipt shows it; the unknown item shows the warning copy and never `0`; clicking a lifecycle step moves `aria-current` and updates the store.
5. **Repository checks**: `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.

Manual verification in the built-in browser: both pages in every locale, both colour schemes, 375px without horizontal overflow, keyboard operation of the ledger, stepper and tabs, reduced motion, the links from the menu, the drawer and the home page, and a sitemap of thirty URLs with matching titles, canonicals and hreflang.

## Out of Scope

- The AI Voice page, and destinations for Docs, GitHub and Book a Demo.
- Porting the Vision page's sections and FAQ to Tailwind so more of the frozen stylesheets can leave.
- Growing AGENTS.md's product table beyond Vision and FinOps.
- Real data, a budget gate, a financial ledger, savings claims or any console feature; the pages describe the intended model with illustrative values.

## Further Notes

The source material is Exphere's knowledge base export 1.0.0 (2026-10-21) and the documentation at exphere.vercel.app/docs. Both separate product principles and proposals from current source; the pages keep that separation. The connected agent the Microruntimes page describes is the `ex` CLI in `packages/ex`.
