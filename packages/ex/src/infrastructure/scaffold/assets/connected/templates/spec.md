# Technical specification

Save the artifact using the SPEC pattern and sequence in the artifact tracker (consumer-owned reference: `../../.exphere/artifacts/design/0006-DD-artifact-tracker.dd.md`).

- ID: `SPEC-<id>`
- State: `draft`
- PRD: `<ID and approved version; record the gap if absent>`
- Technical owners: `<consumer and platform>`

## Contract and composition

- Requirements addressed: `<REQ IDs>`
- Components selected by the consumer: `<names, interfaces, and versions>`
- Environment and executor selected by the consumer: `<decision or open question>`
- Extension points and execution lifecycle: `<contracts; avoid implicit dependencies>`
- Configuration, secret references, and data ownership: `<fill in>`

## Responsibilities

- Consumer: `<composition, adapters, scenarios, and configuration>`
- Platform: `<policies, gates, evidence, and execution mechanisms>`
- Service owners: `<approval of applicable targets and limits>`

## Operation and limits

- Resources, isolation, permissions, and capacity: `<MET IDs/decisions>`
- Failures, timeouts, cancellation, and cleanup: `<fill in>`
- Installation, updates, and removal: `<expected behavior>`
- Observability and execution correlation: `<evidence contract>`

## Verification and decisions

- Affected guardrails: `<FIN/12F IDs and version>`; consult the policy-specific guidance (consumer-owned reference: `../rules/README.md#policies-by-context`).
- Applicability, criteria, owners, and planned evidence: `<record only the affected scope; justify non-applicability>`
- Cost/consumption impact: `<when applicable: allocation, estimate with sources/assumptions, authorized budget, monitoring, and shutdown>`
- Open issues and exceptions: `<decision section IDs or unverified points; a plan does not prove compliance>`

- Verification plan: `<VER section ID and containing artifact>`
- Trade-offs, decisions, and blocking questions: `<decision section IDs>`
- Technical review: `pending; owners and version to be recorded`
