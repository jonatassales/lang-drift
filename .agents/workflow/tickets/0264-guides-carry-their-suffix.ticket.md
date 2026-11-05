# 0264: the handbook and the skill references carry the guide suffix

**Spec:** .agents/workflow/specs/0036-filesystem-rules-compliance.spec.md

**What to build:** The seven unsuffixed documents take `.guide.md`: the repository handbook becomes `.agents/context.guide.md` and the skills' references become `mocking.guide.md`, `tests.guide.md`, `adr-format.guide.md`, `glossary-format.guide.md`, `deepening.guide.md` and `design-it-twice.guide.md`. `AGENTS.md`, every `SKILL.md` and skill document, the comments in the platform CLI that cite the handbook and the hooks follow the rename; published specs, tickets and research notes keep their wording. The Markdown-suffix detector accepts `guide`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] Every Markdown document outside the conventional names carries a type suffix.
- [ ] Every path cited in `AGENTS.md` and under `.claude/skills/` exists, and a grep for the old names outside `.agents/workflow/` returns nothing.
- [ ] `platform/cli/drift quality check filesystem` reports no unsuffixed document.

## Dev part

**Override:** No part qualifies: edits to rules, documents and references, with no problem of their own to solve (policy, Dev part). 2026-11-04.

## Verification

- `platform/cli/drift quality check filesystem`
