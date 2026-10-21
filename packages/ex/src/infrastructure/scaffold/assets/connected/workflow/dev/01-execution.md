# Execution

**Input:** an identified spec or ticket revision, satisfied blocking dependencies, and authorization for the implementation scope. Follow the [cycle contract](CYCLE.md).

## Work

1. State the selected ticket or scope. Load [implement](../skills/implement/SKILL.md), relevant working rules (consumer-owned reference: `../../../.agents/rules/README.md`), and only the needed skills. Inspect current files and applicable local instructions before editing.
2. Preserve user changes and identify the comparison baseline. Use Git when available; otherwise capture the task's file manifest and relevant contents or hashes. An absent repository does not authorize initializing one.
4. Update affected references and state descriptions together with the change. Work artifacts follow the tracker (consumer-owned reference: `../../artifacts/design/0006-DD-artifact-tracker.dd.md`) and use en-US. Record new scope or assumptions before they become hidden dependencies.
5. Perform only relevant authorized checks. For documents, inspect structure, links, anchors, paths, terminology, and applicable JSON. For code, select available formatting or compilation checks if appropriate. Inspect manifest commands before using them; do not repair absent helper scripts outside the ticket's scope.
6. Make the final changed scope available for [review](02-review.md), including relevant uncommitted and newly created files. Record checks performed, failures, skipped verification, and limitations.

Use TDD at the authorized hook-contract seam for SPEC-GIT-HOOKS-SDD. The build/test push gate is authorized; unrelated tests and service probes still need corresponding authorization. No commit, publication, provisioning, or production action is implied by completing local work.

**Output:** the implemented deliverable, updated artifacts, and a reviewable snapshot with a concise verification record.

**Completion:** the authorized changes are ready for technical review. Advance the ticket to `in-review`; review and audit determine whether its criteria are satisfied.

**Return condition:** a scope change returns to planning; an implementation defect remains execution work. Record any external dependency or unavailable mandatory check instead of declaring success.

For an authorized commit, use a Conventional Commit message with one `Ticket:` trailer under the [SDD rule](../../harness/rules/0001-sdd.rule.md). Pre-commit checks staged eligible content and reports harness-score without a target; commit-msg resolves staged artifact IDs. Keep warnings visible without masking blocking lint or message errors.
