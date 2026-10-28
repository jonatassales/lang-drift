# FinOps — guardrails

**Local policy v1, adopted October 10, 2026.** Project criteria for changes affecting resources or paid consumption: cloud, CI, storage, network, observability, services, and AI calls. These are work and review rules; automatic enforcement does not exist yet.

## Applicability

Assess change impact before requiring estimates. Documentation edits with no change in consumption may record `not-applicable` with a sentence. A proposal may advance with open assumptions; a paid action depends on the corresponding criteria and the authorization already established for its scope. This policy does not authorize provisioning, purchases, billing access, or tests.

The criteria below adapt the [FinOps principles](https://www.finops.org/framework/principles/) and [engineering responsibilities](https://www.finops.org/framework/persona/engineering/) to this project. Responsible owners define budgets and targets for each case; there are no universal values or assumed organizational targets.

## Criteria and evidence

| ID | Project rule | Required evidence at the relevant stage |
| --- | --- | --- |
| FIN-01 | Identify the responsible team, service, environment, and billing scope of affected resources. Shared costs require an explicit allocation method. | In the spec: owners and the connection between resources and a cost center/project through tags, labels, or another supported mechanism. Identify unallocated portions. |
| FIN-02 | Estimate the effect of creating, expanding, or changing paid consumption before acting. Separate estimates, budgets, and observed spend. | Price source and date, currency, region/provider when relevant, volume, duration, scenarios, and uncertainty. Include affected compute, storage, network, CI, telemetry, and AI items; justify exclusions. |
| FIN-03 | Define an authorized budget or consumption envelope, duration, and response to deviations before starting a paid execution. | In the spec/decision: scope, limits, monitoring owner, and permitted action. Budget alerts and delayed billing data do not prove automatic spend interruption. |
| FIN-04 | Size CPU, memory, storage, replicas, and concurrency against required load and quality. | A plan with units, conditions, and criteria; after execution, observed consumption and effects on latency, errors, and availability. A descriptor value does not prove enforcement or capacity. |
| FIN-05 | Temporary resources need an expected duration, ownership, and shutdown procedure; persistent resources need retention and an owner. | Inventory what will be created, retained, or removed. After execution: remaining resources and their disposition. Cleanup affects only authorized scope, preserving data and shared resources. |
| FIN-06 | Track consumption with a source, window, and owner; define how to investigate deviations. | Query/dashboard/report, known source delay, and agreed alert criteria. Relate cost changes to volume, price, configuration, and quality; missing data remains unknown. |
| FIN-07 | Compare economic results under explicit conditions. | Baseline, period, volume, and included costs. When useful, unit cost with a defined numerator and denominator. Separate estimated savings, realized savings, cost avoidance, and recovered work capacity. |
| FIN-08 | Optimization preserves service security and quality requirements. | The decision records the alternative, expected impact, change effort, and reversibility. Changes to region, redundancy, retention, or commercial commitments require assessment by affected owners; lower price alone does not demonstrate a better outcome. |

For agents, FIN-02/03/06 include calls, model/provider, tokens or billable units, retries, parallelism, and auxiliary services when changed. Tokens and artifact counts do not demonstrate delivered value. An estimated monetary limit must not be described as a technical gate without an implemented mechanism.

## Application in workflows and microruntimes

- Specification: FIN-01 through FIN-05 and FIN-08 when applicable. Record open issues before committing spend; use the existing template.
- Verification/observation: FIN-04 through FIN-07 with scoped verification records, distinguishing plans from results.
- Review: record the criterion, evidence, assessment, and limitations in the review format. The former shared states and exceptions index is unavailable; the format does not replace that contract. A violation or missing mandatory evidence prevents a compliance claim for the affected action; it does not prevent independent research or documentation.

The `ex` CLI does not calculate costs or enforce budgets, quotas, or expiration. Record contracts in the spec and verification plan; do not invent fields accepted by the current schema. A future adapter must demonstrate which control it performs. Stack and composition remain consumer choices.

## Additional reference

[FinOps Foundation — Planning & Estimating](https://www.finops.org/framework/capabilities/planning-estimating/), consulted October 10, 2026. FIN rules are local decisions informed by these sources, not Foundation certification or another company's policy.
