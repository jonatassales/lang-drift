# 0315: the repository's documents say `dft`

**Spec:** .agents/workflow/specs/0047-the-langdrift-cli-answers-to-dft.spec.md

**What to build:** An agent or contributor reading this repository's guides meets the LangDrift CLI and `dft`. That covers AGENTS.md, the context guide, the root README, the Safe Seat policy, and the rules `10-testing` and `11-nestjs` where they name the command. The glossary tells `dft`, `drift` and the former `ex` apart. Every open spec and ticket says `dft`, and ticket 0228 waits for this rename. Closed specs and tickets keep `ex`.

**Blocked by:** 0313.

**Status:** ready-for-agent

- [ ] The glossary defines `dft`, `drift` and `ex` (former name of `dft`).
- [ ] AGENTS.md, the context guide, the root README, the Safe Seat policy and rules 10 and 11 name the command `dft`; the rules change only the name.
- [ ] Every spec with an unchecked ticket, and every unchecked ticket, says `dft`; no checked ticket or closed spec changes.
- [ ] Ticket 0228 is blocked by 0313, 0314 and 0315.

## Dev part

**Override:** No part qualifies: the ticket renames a command in documents without changing a decision, and a rename or a migration is never a dev part (policy, Categories). 2026-11-05.

## Verification

- `git grep` for `ex` as the command outside `.agents/workflow/` closed documents and 0.4.0 fixtures finds nothing.
