# 0335: LLMOps takes the FinOps page and the menu splits

**Spec:** .agents/specs/0067-llmops-joins-the-platform.spec.md

**What to build:** A visitor opens the Platform menu and reads two columns: LLMOps, FinOps and Safe Seat on the left; Microruntimes, LangDrift CLI and CI Gates on the right; one column in the same order on a phone. LLMOps links to `/platform/llmops`, which shows everything today's FinOps page shows (ledger, receipt, gate flow, failure scenarios, budget terms, sources), named LLMOps, in five locales, with its own metadata and sitemap entries. The home's ledger section is the LLMOps section and links there; the home FAQ's LLMOps tab takes the FinOps tab's place. Safe Seat says it is part of LangDrift LLMOps. The menu already describes FinOps as hardware; the `/platform/finops` route itself is left to ticket 0336, which gives it its new content in the same pull request.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

**Alignment:** none. The spec changes the website's marketing pages, not the console's runtime view, so ADR 0005's simulation does not apply.

- [ ] The menu lists LLMOps, FinOps, Safe Seat, Microruntimes, LangDrift CLI, CI Gates, in that DOM order, the first three in the left column from `sm`.
- [ ] LLMOps carries its own icon from the platform icon family and the description "Reasoning spend accountable to a team budget" (or its catalog wording) in five locales.
- [ ] `/en/platform/llmops/` renders the moved sections with the LLMOps heading; its title, description, canonical and hreflang are its own; `sitemap.xml` lists it in five locales.
- [ ] The home section reads LLMOps, has the `llmops` anchor and links to `/platform/llmops`.
- [ ] Safe Seat's menu description, flow lead and availability answer name LangDrift LLMOps.
- [ ] Catalog parity, the moved tests and the menu test pass; `drift check localization` passes.

## Dev part

**Override:** No part qualifies: the ticket moves an existing page under a new name, reorders a list and rewrites copy (policy, Categories: never a rename or a migration). The FinOps page's compute envelope in ticket 0336 holds this spec's dev part. 2026-11-05.
