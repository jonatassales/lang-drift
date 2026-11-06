# 0336: FinOps holds hardware spend

**Spec:** .agents/workflow/specs/0067-llmops-joins-the-platform.spec.md

**What to build:** A visitor opens `/platform/finops` and reads hardware spend: a compute envelope split into used, held, idle and available; a ledger of environments (a microruntime, a CI runner, a preview cluster, a GPU job) with CPU-hours, memory GB-hours, GPU-hours and cost; selecting one updates its environment receipt (requested and observed usage, peak, idle time, cleanup state, source). One environment's meter stopped reporting and it reads Unknown, never zero. Below come two failure scenarios (an environment left running after its run, a runner whose meter stopped), the budget terms of compute, the sources and a FAQ, in five locales. The page's metadata describes hardware, and the home FAQ gains a FinOps tab after LLMOps.

**Blocked by:** 0335 (the FinOps route and catalog name must be free, and the formatters live under LLMOps)

**Status:** ready-for-agent

**Alignment:** none. The spec changes the website's marketing pages, not the console's runtime view, so ADR 0005's simulation does not apply.

- [ ] `/en/platform/finops/` shows the compute ledger and the environment receipt; selecting an environment updates the receipt.
- [ ] The unmetered environment reads Unknown in the ledger and the receipt, and the envelope never counts it as zero.
- [ ] The failure scenarios, budget terms, sources and FAQ render in five locales; catalog parity passes.
- [ ] `/platform/finops/`'s title and description speak of hardware in five locales.
- [ ] The home FAQ offers LLMOps, then FinOps.

## Dev part

**Category:** `core-domain`. How a compute envelope splits into used, held, idle and available, and how an unmetered environment holds its requested ceiling instead of counting as zero, is the FinOps domain rule a visitor reads first and the owner must be able to explain.

**Paths:** `apps/website/src/app/[locale]/(pages)/platform/finops/content/compute-envelope.compute.ts` and `compute-envelope.compute.unit.ts` beside it.

**Boundary:** `computeEnvelope(envelope: ComputeEnvelope): ComputeEnvelopeSegments`, with `ComputeEnvelope` (authorized amount and environments, each with requested ceiling, observed cost or `null` when unmetered, idle cost and cleanup state) and `ComputeEnvelopeSegments` (`used`, `held`, `idle`, `available`) declared in `finops-content.types.ts`. The agent writes the types and the ledger that calls the function.

**Goal:** Given the envelope, return the four segments: used is the observed cost of metered environments less their idle cost, idle is their idle cost, held is the requested ceiling of every unmetered environment, and available is the authorized amount less the other three, never below zero.

- [ ] The four segments add up to the authorized amount whenever spend stays inside it.
- [ ] An unmetered environment adds its requested ceiling to held and nothing to used.
- [ ] Spend above the authorized amount leaves available at zero, never negative.

**Override:** The owner does not write this part and asked the agent to write it ("Override, you write it"). 2026-11-05.
