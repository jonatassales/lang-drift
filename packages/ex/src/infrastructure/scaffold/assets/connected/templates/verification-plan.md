# Verification plan

Use this format as a section of the relevant spec or DD, following the artifact tracker (consumer-owned reference: `../../.langdrift/artifacts/design/0006-DD-artifact-tracker.dd.md`).

- Section ID: `VER-<id>`; state: `draft`
- PRD/spec and versions: `<references>`
- Owners of targets and execution: `<people/teams>`
- Execution profile: `<components, versions, resources, data, and environment>`
- Affected guardrails: `<FIN/12F IDs and version; required method and evidence>` according to the policy-specific guidance (consumer-owned reference: `../rules/README.md#policies-by-context`).

This template does not authorize execution. Respect the scoped build/test and hook-contract authorization in root AGENTS. For unrelated behavior without test authorization, record checks as pending and limit work to authorized inspection.

## Coverage and evidence

| ID | Requirement | Type/scenario | Executor | Criterion | Evidence |
| --- | --- | --- | --- | --- | --- |
| TEST-001 | REQ-001 | `<smoke/e2e/load/stress/soak/security>` | `<to be defined>` | `<condition>` | `<result and artifacts>` |

Gherkin describes expectations; each scenario needs a link to its executable test before it can serve as execution evidence.

## Metric contracts

| ID | Dimension | Indicator/unit | Approved limit | Window/conditions | Source |
| --- | --- | --- | --- | --- | --- |
| MET-001 | CPU | `<usage, peak, throttling>` | `<pending>` | `<define>` | `<define>` |
| MET-002 | Memory | `<peak, growth, retention, OOMs>` | `<pending>` | `<define>` | `<define>` |
| MET-003 | Storage | `<usage, growth, I/O latency>` | `<pending>` | `<define>` | `<define>` |
| MET-004 | Load | `<arrival rate, concurrency, achieved throughput>` | `<pending>` | `<define>` | `<define>` |
| MET-005 | Quality | `<errors, timeouts, p95/p99 per operation>` | `<pending>` | `<define>` | `<define>` |
| MET-006 | Availability | `<SLI, SLO, error budget>` | `<pending>` | `<actual operation>` | `<define>` |
| MET-007 | Cost | `<currency/period or cost per defined unit; estimated or observed>` | `<authorized budget when applicable>` | `<volume, duration, and included items>` | `<price/billing/consumption, date, and limitations>` |

## Result and gate

- Execution and results: `absent; plan not executed`.
- Approval requires approved applicable criteria, achieved load, and complete evidence; missing mandatory evidence keeps the gate blocked.
- Record the code/configuration revision, environment, timing, observations, and artifacts correlated by execution.
- Soak tests investigate leaks; a local test does not prove the universal absence of leaks or production availability.
