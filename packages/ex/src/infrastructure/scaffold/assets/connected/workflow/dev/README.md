# Dev — software delivery within SDD

Dev is the software specialization of the [common SDD workflow](../README.md) and its [declarative definition](../definition.json). It documents how to use the installed Pocock skills; it is not an executor, CI gate, approval service, or proof of organizational adoption.

## Entry points

Start with the phase the user has authorized. Reuse the existing conversation, artifacts, and valid authorization. A direct request for a spec does not require a new interview; an implementation request can cover an existing spec and its tickets. Resolve only gaps that materially prevent the requested work.

The [cycle and authorization contract](CYCLE.md) owns the stage mapping, ticket-state rules, and handoffs. Read it before performing a Dev phase. Each phase below states its required inputs, output, and completion condition.

| Document | Responsibility |
| --- | --- |
| [Planning](00-planning.md) | Clarify the problem when needed, specify scope, and prepare traceable tickets. |
| [Execution](01-execution.md) | Implement the authorized scope while preserving user changes and recording verification limits. |
| [Review](02-review.md) | Assess standards and spec alignment against the same artifact snapshot. |
| [Audit](03-audit.md) | Reconcile requirements, changes, reviews, evidence, and remaining delivery conditions. |
| [Refinement](04-refinement.md) | Learn from observed outcomes and propose bounded changes for a separately authorized cycle. |

## Project contracts

- AGENTS.md (consumer-owned reference: `../../../AGENTS.md`) and working rules (consumer-owned reference: `../../../.agents/rules/README.md`): scope, authorization, language, and scoped hook-test authorization and restrictions on unrelated tests.
- Context (consumer-owned reference: `../../../.agents/CONTEXT.md`) and domain configuration (consumer-owned reference: `../../artifacts/design/0007-DD-domain-model.dd.md`): consumer boundaries and available vocabulary.
- Artifact tracker (consumer-owned reference: `../../artifacts/design/0006-DD-artifact-tracker.dd.md`): typed filenames, numbering, directories, ticket states, and decision records.
- [Guardrails](../../harness/guardrails): applicable local policies and evidence criteria.
- [Verification plan](../../templates/verification-plan.md) and [review format](../../templates/review.md): existing formats for scoped criteria and observations. The former evidence README is absent from the current tree.

All project content uses en-US. User conversation may use another language. Preserve microruntime terminology, public third-party attribution, and the distinction between proposals and observed results.

## Installed skills

These links locate the relevant instructions; they do not invoke skills or authorize actions. Read the applicable `SKILL.md` when no Skill tool is available.

| Need | Skill |
| --- | --- |
| Primary-source investigation | [research](../skills/research/SKILL.md) |
| Explicit discovery discussion | [grill-with-docs](../skills/grill-with-docs/SKILL.md), using [grilling](../skills/grilling/SKILL.md) and [domain-modeling](../skills/domain-modeling/SKILL.md) |
| Specification | [to-spec](../skills/to-spec/SKILL.md) |
| Ticket decomposition | [to-tickets](../skills/to-tickets/SKILL.md) |
| Authorized implementation | [implement](../skills/implement/SKILL.md) |
| Module and interface design | [codebase-design](../skills/codebase-design/SKILL.md) |
| Technical review | [code-review](../skills/code-review/SKILL.md) |
| Explicit retrospective | [retro](../skills/retro/SKILL.md) |
| Test-first development at an authorized seam, including the hook-contract suite | [tdd](../skills/tdd/SKILL.md) |

Skill availability does not authorize a later phase. Current project restrictions override skill defaults, including automatic tests or commits. The installed versions and their limitations remain in the skill configuration (consumer-owned reference: `../../artifacts/design/0008-DD-agent-skills.dd.md`).

## Local developer checks

The command contract (consumer-owned reference: `../../../README.md#developer-checks`) configures Husky through its native prepare command without user PATH mutation. Pre-commit blocks staged lint failures and reports an informational harness score; commit-msg blocks invalid Conventional Commits and warns about ticket links; pre-push reports SDD findings before blocking build/test checks. Apply the [SDD rule](../../harness/rules/0001-sdd.rule.md) and [guardrail](../../harness/guardrails/0002-sdd.gr.md). Review checks actual ticket/spec scope; audit records unavailable evidence and bypass limitations.
