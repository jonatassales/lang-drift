# Dev cycle and authorization

The [common definition](../definition.json) owns lifecycle identifiers, transitions, blocking, and verification results. Dev uses that model; the pages below group activities for readers rather than define a second state machine.

## Stage mapping

| Common stage IDs | Dev activity | Output or exit condition |
| --- | --- | --- |
| `captured`, `discovery` | [Planning](00-planning.md) | Known request, source, owner, evidence, and material unknowns. Reuse sufficient existing context. |
| `proposed`, `specified` | [Planning](00-planning.md) | Human priority decision where needed; reviewed scope and acceptance criteria in typed artifacts. |
| `ready` | [Planning](00-planning.md) → [execution](01-execution.md) | Applicable spec revision, resolved blocking dependencies, and authorization for the bounded work. |
| `implementing` | [Execution](01-execution.md) | Changes within scope, with available checks and limitations recorded. |
| `verifying` | [Review](02-review.md) and [audit](03-audit.md) | Standards and spec findings plus the evidence needed for each applicable criterion. |
| `delivery-approved` | [Audit and delivery handoff](03-audit.md#delivery-handoff) | Authorization and satisfied conditions for the actual destination and effect. |
| `observing` | [Observation](03-audit.md#observation) | Signals or artifact observations appropriate to the delivered change; unknown operational results stay unknown. |
| `reviewed` | [Refinement](04-refinement.md) | Outcomes assessed and improvement proposals recorded; another cycle is not automatic. |

Missing or failed verification can return to implementation. Changed scope or criteria return to specification. Record the reason and resume condition for a blocked activity while preserving its lifecycle stage. Continue independent authorized work when useful. Do not fabricate earlier-stage decisions merely because the user started at a later phase.

## Ticket states and authorization

The tracker (consumer-owned reference: `../../artifacts/design/0006-DD-artifact-tracker.dd.md`) is the authority for ticket states. They are separate from the lifecycle stage IDs above.

- `draft`: proposed work that has not met the conditions for readiness.
- `ready-for-agent`: the applicable input is accepted, execution is authorized, and blocking dependencies are satisfied.
- `in-progress`: authorized implementation is underway.
- `in-review`: the resulting scope is available for review and verification.
- `done`: the ticket's deliverable and applicable acceptance criteria are met, with review and evidence recorded.
- `blocked`: record the condition preventing progress and the state to return to.

Associate authorization with the user request, scope, artifact version, and date when known. A direct implementation instruction can authorize the current spec and tickets without another confirmation. Do not infer organizational adoption, production permission, or acceptance of materially changed scope from that instruction.

Passing a review does not authorize commits, pushes, remote tickets, pull requests, deployments, provisioning, or paid actions. Reuse valid authorization for the same effect; otherwise keep that handoff pending and complete independent local work. A documentation ticket can finish without a production deployment when deployment is outside its scope.

## Invocation and verification limits

Load only the skills required for the authorized phase. Explicit-invocation skills retain that boundary; a link or sequence does not trigger them. Apply dependency skills according to the invoked skill and project rules. Do not require grilling when the user requests synthesis or provides sufficient decisions.

SPEC-GIT-HOOKS-SDD authorizes build/test before push and the scoped hook-contract suite in disposable repositories, including a local bare remote. Unrelated application tests and service probes remain outside that authorization. Use relevant authorized static inspection, formatting, or compilation; do not run an unavailable task wrapper to simulate verification. If a criterion requires behavior that cannot be observed within the authorization, record it as `inconclusive` and do not treat missing mandatory evidence as a pass.

Review the exact changed scope, including uncommitted work. Without Git, record a file manifest and content hashes instead of inventing commits or initializing a repository. Reuse review only while that snapshot and its applicable criteria remain valid.

Local [SDD traceability](../../harness/rules/0001-sdd.rule.md) warns for now. Missing SDD evidence remains `not-verified` or `non-compliant`, with an advisory exception to delivery blocking. Other mandatory evidence and blocking quality failures retain their existing effect.
