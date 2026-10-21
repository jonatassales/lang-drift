# Workflow skills

The [AIHero catalog](https://www.aihero.dev/skills), inspected October 13, 2026, lists 27 skills. This directory contains 26: the owner explicitly excluded `writing-for-agents`. The Go catalog is separately owned by `.agents/skills` (consumer-owned reference: `../../../.agents/skills/README.md`); unrelated FinOps skills remain under the shared harness (consumer-owned reference: `../../harness/skills/`).

Source: [mattpocock/skills](https://github.com/mattpocock/skills/tree/49dd158d1076134a641b33efb035946536778336), commit `49dd158d1076134a641b33efb035946536778336`. The original [MIT license](UPSTREAM-LICENSE) applies to the upstream material. Preserve this attribution when copying or updating it.

## Inventory

| Catalog group | Installed skills |
| --- | --- |
| Getting started | `setup-matt-pocock-skills`, `ask-matt` |
| Main flow | `grill-with-docs`, `to-spec`, `to-tickets`, `implement`, `code-review`, `retro` |
| Shaping | `wayfinder`, `prototype`, `research` |
| Engineering | `improve-codebase-architecture`, `diagnosing-bugs`, `triage`, `wizard`, `implement-spec` |
| Productivity | `grill-me`, `handoff`, `to-questionnaire`, `teach`, `wait-what` |
| Reference | `pr`, `codebase-design`, `domain-modeling`, `grilling`, `tdd` |

The 11 preexisting local adaptations (`code-review`, `codebase-design`, `domain-modeling`, `grill-with-docs`, `grilling`, `implement`, `research`, `retro`, `tdd`, `to-spec`, `to-tickets`) were relocated from `.exphere/harness/skills/` with their support files. They were not overwritten by upstream downloads. The other 15 were installed from the pinned repository using the skill installer during PR #19 completion.

## Local adaptations and limits

- Root `AGENTS.md`, the artifact tracker (consumer-owned reference: `../../artifacts/design/0006-DD-artifact-tracker.dd.md`) and domain configuration (consumer-owned reference: `../../artifacts/design/0007-DD-domain-model.dd.md`) take precedence over generic upstream defaults. Existing local tracker configuration does not require another setup round. New remote issues, comments, labels, commits, publication, tests or infrastructure changes need their corresponding authorization.
- `ask-matt`, `setup-matt-pocock-skills`, `implement-spec`, `wayfinder` and `triage` identify those local contracts at entry. `ask-matt` no longer recommends the explicitly removed `writing-for-agents` skill.
- Optional shell templates `diagnosing-bugs/scripts/hitl-loop.template.sh` and `wizard/template.sh` are deliberately omitted because the owner reserves helper scripts for the CLI. Their skill instructions explain the remaining manual procedure and the unavailable template; installation does not claim those optional helpers work.
- Preserve existing invocation settings. Installation does not invoke a skill, chain lifecycle phases, establish governance, or prove automatic discovery in an agent client. These files are available to read at their canonical paths; clients may refresh discovery on a subsequent turn.
- Upstream directories outside the public AIHero catalog (`in-progress` and `misc`) are not part of the requested catalog and are not installed.
