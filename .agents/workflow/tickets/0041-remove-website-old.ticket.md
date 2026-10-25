# 0041: Remove the old website

**Spec:** None

**What to build:** only the new website exists in the repository. `apps/website-old`, the copy of today's langdrift.md kept for reference, is deleted along with everything that existed to carry it.

**Blocked by:** None (0040, which replaced orbz, is merged).

**Status:** ready-for-agent

- [ ] `apps/website-old` is gone, and the lockfile loses its importer and nothing else.
- [ ] The root `website-old` script, the Biome ignore, the dead-code allowances for its formatters and react-dom, and the `--filter=!website-old` exclusions in the pull-request workflow are gone.
- [ ] `AGENTS.md`, `.agents/context.md`, the website README and the frontend rule no longer describe a reference copy; specs and tickets that mention it stay as history.

### Verification

- `pnpm install --frozen-lockfile`, `pnpm turbo run build check typecheck test` for the whole workspace, `drift check dead-code`, `drift check localization`.
- `git grep website-old` finds only specs and tickets.
