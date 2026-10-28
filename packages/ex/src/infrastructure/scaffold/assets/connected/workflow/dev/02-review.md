# Review

**Input:** the changed scope, applicable spec revision, documented standards, and one consistent review snapshot. Follow the [cycle contract](CYCLE.md).

## Work

1. Use [code-review](../skills/code-review/SKILL.md). Identify the spec and standards relevant to the requested change; include root AGENTS.md (consumer-owned reference: `../../../AGENTS.md`), applicable rules (consumer-owned reference: `../../../.agents/rules/README.md`), and [guardrails](../../harness/guardrails).
2. Give the standards and spec reviewers the same selected paths, contents, baseline, and snapshot identity. Follow the skill's parallel-review procedure. Exclude unrelated edits and third-party files unless the authorized change affects them.
3. Include committed and uncommitted changes when Git is available. Without Git, review the explicit manifest and captured contents as task additions or changes against an available local snapshot. State that no historical Git diff exists.
4. Keep the two findings lists separate. Standards review evaluates documented rules and applicable design concerns; spec review evaluates missing behavior, incorrect implementation, and unrequested scope. Cite actual locations and distinguish defects from suggestions.
5. Resolve actionable findings, then refresh the affected snapshot and review. Reuse an earlier review only when its evidence, criteria, and examined files still apply; do not repeat identical review without a reason.

Use the [review format](../../templates/review.md) within the relevant typed artifact. Record the reviewer role, spec version, inspected scope, findings, verification limits, and outcome. The hook-contract suite is authorized under SPEC-GIT-HOOKS-SDD; unrelated tests remain restricted. Review existing evidence before repeating checks.

**Output:** separate standards and spec assessments tied to the same snapshot, with remaining defects and evidence gaps visible.

**Completion:** applicable findings are resolved or their blocking impact is explicitly recorded. Hand the reviewed scope to [audit](03-audit.md); a review is advisory evidence, not human acceptance or delivery permission.

**Return condition:** implementation defects return to execution; changed requirements or criteria return to planning. Missing mandatory evidence keeps the dependent conclusion inconclusive.

Apply [SDD-01 through SDD-03](../../harness/guardrails/0002-sdd.gr.md): inspect each commit’s ticket and semantic scope. When an actual PR is available, inspect its `Spec:` declaration and compare every included ticket’s parent spec. A mismatched PR is advisory `non-compliant`; absent PR metadata is `not-verified`. A local candidate or a successful hook cannot prove PR association or grant approval.

The [SDD review rule](../../harness/rules/0001-sdd.rule.md#required-code-review-and-publication) requires this two-axis review before completion, commit, or push. Record the examined snapshot; rerun affected review after material changes. Publication uses a feature branch from main and a normal push; force pushes are prohibited.
