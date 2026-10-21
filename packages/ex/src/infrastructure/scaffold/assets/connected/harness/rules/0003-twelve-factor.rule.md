# Applying Twelve-Factor rules

Load the [12F guardrails](../guardrails/0003-twelve-factor.gr.md) when creating or changing services, configuration, dependencies, build/release, persistence, connections, scaling, logs, or microruntime lifecycle.

1. Identify whether the target is an application, library, job, or stateful dependency. Select the affected factors; assess all 12 for a new service.
2. Record the necessary contracts and adaptations in the [spec](../../templates/spec.md), preserving consumer stack and composition choices. Consult the existing schema and adapters before proposing fields or capabilities.
3. Inspect files and configuration relevant to the scope. Perform only authorized checks; respect the scoped hook-test authorization and remaining restrictions in root AGENTS. Static review does not prove runtime behavior.
4. Use the [12F completion criteria](../guardrails/0003-twelve-factor.gr.md#completion-criteria) in the [review](../../templates/review.md), linking findings to 12F IDs and the versions assessed. Record preexisting failures outside the scope separately.
5. External skills are references; 15-factor extensions do not silently replace these 12 criteria. Markdown and the current CLI do not constitute a Twelve-Factor certification system.

Completion: affected factors, evidence, adaptations, and open issues are explicit. End the assessment without turning a limited change into an unrequested architectural migration.
