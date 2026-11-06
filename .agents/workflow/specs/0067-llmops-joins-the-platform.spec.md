# LLMOps joins the platform

**Status:** ready-for-agent

**Decided by:** the owner, 2026-11-05: FinOps splits into LLMOps (reasoning and AI processing spend) and FinOps (hardware and computational resources).

## Tickets

- 0335: LLMOps takes the FinOps page and the menu splits
- 0336: FinOps holds hardware spend (blocked by 0335)

## Problem Statement

The website names every cost FinOps. Its FinOps page is about model calls: a team envelope, a receipt with the model tier and effort, a gate that reserves before every paid call, a call that stays unknown after a timeout, model providers as sources. A visitor who runs agents on real hardware has no page for the CPU, memory and GPU those runs occupy, and a visitor who reads "FinOps" expecting compute finds token spend. The Platform menu lists the engine first and mixes products with utilities, so the two kinds of spend cannot be told apart from it.

## Solution

The platform has two spend products. **LLMOps** holds reasoning and AI processing spend; **FinOps** holds hardware and computational resources. The Platform menu lists six items in two columns: on the left the products, LLMOps, FinOps and Safe Seat; on the right the utilities, Microruntimes, LangDrift CLI and CI Gates. On a narrow screen the menu is one column in that same order.

`/platform/llmops` carries today's FinOps page, every section and demo, named LLMOps. `/platform/finops` carries new hardware content: a compute ledger over the team's hardware envelope, an environment receipt, two failure scenarios, the budget terms of compute and the sources that feed it; its FAQ answers in the home FAQ's FinOps tab, as the LLMOps FAQ does in its own. The home's FinOps section becomes the LLMOps section. Safe Seat is part of LangDrift LLMOps. Every page and catalog exists in `en`, `pt-BR`, `zh-Hant`, `ja` and `vi`.

## User Stories

1. As a visitor, I want the Platform menu to list LLMOps, so that I can find the page about model spend.
2. As a visitor, I want the menu's left column to list LLMOps, FinOps and Safe Seat, so that the products read together.
3. As a visitor, I want the menu's right column to list Microruntimes, LangDrift CLI and CI Gates, so that the utilities read together.
4. As a visitor on a phone, I want the menu in one column, LLMOps first and CI Gates last, so that the order matches the desktop's reading order.
5. As a visitor, I want each menu item to carry its icon and a one-line description, so that I can tell LLMOps from FinOps before I click.
6. As a visitor, I want LLMOps described as reasoning spend accountable to a team budget, so that I know it covers model calls.
7. As a visitor, I want FinOps described as hardware spend metered while it runs, so that I know it covers compute.
8. As a visitor, I want Safe Seat to stay marked Coming soon and to say it is part of LangDrift LLMOps, so that I know where it belongs.
9. As a visitor, I want `/platform/llmops` to show the team envelope ledger and the delivery receipt, so that I see how model spend is held.
10. As a visitor, I want the LLMOps page to keep the gate flow (recommend, reserve, reconcile), so that I see that a control, not the model, authorizes a paid call.
11. As a visitor, I want the LLMOps page to keep the pause-and-replan and unknown-after-timeout scenarios, so that I see how the envelope holds when the plan does not.
12. As a visitor, I want the LLMOps page to keep the budget terms and the sources, so that I know who funds the work and what feeds the ledger.
13. As a visitor, I want the LLMOps page's FAQ to speak of LLMOps, so that its answers match the product I am reading about.
14. As a visitor, I want `/platform/finops` to open with a compute envelope split into used, held, idle and available, so that I see the team's hardware budget at a glance.
15. As a visitor, I want the compute ledger to list environments (a microruntime, a CI runner, a preview cluster, a GPU job) with CPU-hours, memory GB-hours, GPU-hours and cost, so that I see where hardware spend goes.
16. As a visitor, I want to select an environment and read its receipt (requested and observed usage, peak, idle time, cleanup state, source), so that I can follow one environment's spend.
17. As a visitor, I want an environment whose meter stopped reporting to read Unknown, never zero, so that missing evidence stays visible.
18. As a visitor, I want a scenario where an environment is left running after its run with cleanup unconfirmed, so that I see how FinOps catches orphaned spend.
19. As a visitor, I want a scenario where a runner's meter stops reporting, so that I see unknown usage held as unknown until a record settles it.
20. As a visitor, I want the FinOps budget terms (compute envelope, quota, cleanup as a spend control), so that I know how hardware spend is governed.
21. As a visitor, I want the FinOps sources named (Microruntimes observations, cluster metrics, CI runner minutes, cloud billing), so that I know what feeds the compute ledger.
22. As a visitor, I want a FinOps FAQ that says what FinOps covers and what LLMOps covers, so that I do not confuse the two.
23. As a visitor on the home, I want the product section that shows the ledger and receipt to be named LLMOps and to link to `/platform/llmops`, so that the home and the page agree.
24. As a visitor on the home, I want the FAQ tabs to offer LLMOps and FinOps, so that I can read the questions of each.
25. As a Portuguese, Chinese, Japanese or Vietnamese reader, I want both pages, the menu and the home in my language, so that the split reads the same in every locale.
26. As a search engine, I want `/platform/llmops/` in every locale with a unique title, description, canonical, hreflang alternates and a sitemap entry, so that the new page is indexed correctly.
27. As a search engine, I want `/platform/finops/`'s title and description to describe hardware spend, so that its listing matches its content.
28. As a keyboard user, I want the menu's focus order to follow LLMOps, FinOps, Safe Seat, Microruntimes, LangDrift CLI, CI Gates, so that the order I hear matches the order I see.
29. As a screen-reader user, I want each page's hero to keep its summary and its skip link, so that both pages stay navigable.

## Implementation Decisions

- **The menu's data owns the order.** The platform products list gains LLMOps and is ordered LLMOps, FinOps, Safe Seat, Microruntimes, CLI, CI Gates. The menu's grid flows by column with three rows from `sm`, so the list fills the left column first; one column keeps the same order. The centre divider stays.
- **An LLMOps icon joins the platform icon family** in `packages/web`, drawn in the family's style; FinOps keeps its icon.
- **LLMOps is today's FinOps page, moved.** Its route folder, sections, content, test data, catalog and identifiers move to the LLMOps names, with history kept. Its copy changes only where it names the product. The selected-work-item atom follows the move.
- **The header catalog** gains `llmops` and `llmopsDescription`; `finopsDescription` describes hardware; `safeSeatDescription` names LangDrift LLMOps.
- **The home section is LLMOps.** Its section, anchor and link name LLMOps; the platform catalog's section copy becomes `sections.llmops`. The home FAQ shows an LLMOps tab, then a FinOps tab from the new FinOps catalog. The home hero's rotating words stay.
- **Safe Seat's catalog** names LangDrift LLMOps where it named FinOps.
- **The FinOps page is new.** It composes the shared product-page pieces (product section, pillar scenarios, record card, FAQ shell, heading, tagline) with its own sections: a compute ledger demo, an environment receipt, failure scenarios, budget terms and sources. Its selection state is a Jotai atom of its own concern. Its amounts and formatting reuse the moved formatters and segment computations by alias where they fit.
- **The FinOps catalog is new**, in five locales, and keeps the platform's principles: unknown is never zero, evidence over assertion, recommendation is not authorization. Its demos carry no "simulated" labels.
- **Metadata.** The pages metadata gains `llmops` with a unique localized title and description and rewrites `finops` for hardware; the home description names both products. The page-path map gains `llmops`. The static sitemap gains the `/platform/llmops/` entries for five locales with hreflang alternates.
- **Glossary.** LLMOps and FinOps are defined; Safe Seat is an LLMOps harness capability; the platform products are LangDrift LLMOps and LangDrift FinOps.

## Testing Decisions

- A good test renders a component or loads a catalog and asserts what a visitor or a crawler sees: names, links, order, values and states. It never asserts class names or internal structure.
- **The platform menu** (prior art: the existing platform menu test): six items in order, each linking to its page, LLMOps linking to `/platform/llmops`, Safe Seat alone Coming soon, Safe Seat's description naming LLMOps.
- **The LLMOps sections** keep their moved tests (ledger demo, gate flow, failure scenarios, content data), renamed to their new paths.
- **The FinOps sections** get tests at the same seam (prior art: the budget ledger demo and failure scenarios tests): selecting an environment updates the receipt; the unknown environment reads Unknown, not zero; the envelope segments add up to the authorized amount.
- **The home** (prior art: the home product section and platform FAQ section tests): the section reads LLMOps and links to its page; the FAQ offers LLMOps then FinOps.
- **Catalogs** (prior art: the catalog parity test): `llmops` and `finops` carry the same keys, placeholders and valid ICU in every locale, and the counts 0, 1, 2, 3 and 12 format.
- **Icons** (prior art: the platform icons test): the LLMOps icon renders like the others.

## Out of Scope

- The CLI and CI wording: "FinOps gates" and "FinOps readings" in the CLI, CI Gates and Microruntimes catalogs follow `dft` and the simulation, and stay.
- The conversion plans and the book-a-demo copy, which name FinOps as the platform.
- The console, which keeps its Microruntimes and FinOps areas.
- AGENTS.md's product table.
- A FinOps section on the home.

## Further Notes

- The owner chose new hardware content for FinOps, the home section as LLMOps only, and Safe Seat under LLMOps.
- Spec 0010 made FinOps the product; this spec splits that product in two on the website.
