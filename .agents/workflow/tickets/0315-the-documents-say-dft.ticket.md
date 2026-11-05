# 0315: the repository's documents say `dft`

**Spec:** .agents/workflow/specs/0047-the-langdrift-cli-answers-to-dft.spec.md

**What to build:** An agent or contributor reading this repository's guides meets the LangDrift CLI and `dft`. That covers AGENTS.md, the context guide, the Safe Seat policy, `SECURITY.md`, the infrastructure READMEs, and the rules `10-testing` and `11-nestjs` where they name the command. The glossary tells `dft`, `drift` and the former `ex` apart. Every open spec and ticket says `dft`, and ticket 0228 waits for this rename. Closed specs and tickets keep `ex`.

**Blocked by:** 0313.

**Status:** ready-for-agent

- [x] The glossary defines `dft`, `drift` and `ex` (former name of `dft`).
- [x] AGENTS.md, the context guide, the Safe Seat policy, `SECURITY.md`, the infrastructure READMEs and rules 10 and 11 name the command `dft`; the rules change only the name.
- [x] Every open spec and ticket says `dft` where it plans work; history keeps `ex` under a dated note; no closed document changes.
- [x] Ticket 0228 is blocked by 0313, 0314 and 0315.

> **Note, 2026-11-05.** The lists lag the merges (0222, 0225, 0229 and 0268 have landed unchecked), so a ticket counts as open when it is unchecked and no commit on `main` names it. That leaves specs 0035, 0036, 0037, 0042 and 0045 and tickets 0198, 0223, 0224, 0226, 0227, 0228 and 0260. Spec 0037 describes 0268's launcher rename, and spec 0045 describes what 0.4.0 did and the `ex agent init` that 0221 removed, so both keep `ex` there under a dated note. The context guide's row for the root's pinned 0.4.0 keeps `pnpm exec ex` until 0228 moves the pin. The root README moved to 0314, because it ships with the release. Ticket files keep their slugs.

## Dev part

**Override:** No part qualifies: the ticket renames a command in documents without changing a decision, and a rename or a migration is never a dev part (policy, Categories). 2026-11-05.

## Verification

- `git grep` for `ex` as the command outside `.agents/workflow/` closed documents and 0.4.0 fixtures finds nothing.
