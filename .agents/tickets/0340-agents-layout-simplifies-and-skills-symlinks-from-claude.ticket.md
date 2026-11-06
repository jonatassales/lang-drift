# 0340: agents layout simplifies and skills symlinks from claude

**Spec:** None. The owner's instruction of 2026-11-06: `.agents` simplifies directly into `adr/`, `policy/`, `prd/`, `specs/`, `tickets/`, `dd/`, `research/` and `context.md`, `.claude/skills` holds skills with a symlink at `.agents/skills`, and every path reference across tooling and docs follows.

**What to build:** The `.agents` directory holds artifacts directly without intermediate workflow nesting: `adr/`, `policy/`, `prd/`, `specs/`, `tickets/`, `dd/`, `research/` and `context.md`. Skills are defined in `.claude/skills` and linked into `.agents/skills`. References in `AGENTS.md`, `.agents/context.md`, `packages/drift`, `.claude/hooks` and skills update to the new paths.

**Blocked by:** None.

**Status:** completed.

- [x] `.agents/skills` symlinks to `.claude/skills`.
- [x] All document directories sit directly under `.agents/`.
- [x] `drift doctor`, `drift check workflows`, `drift check biome` pass cleanly.
- [x] References across scripts, hooks, and documentation point to `.agents/tickets/`, `.agents/specs/`, etc.

## Dev part

**Override:** a move of tooling and documentation layout; no part qualifies under the Safe Seat policy. 2026-11-06.
