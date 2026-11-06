# 0211: the workspace knows `lib/` and `infrastructure/`

**Spec:** .agents/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The monorepo gains its two new roots before anything moves into them. pnpm lists `lib/*` and `infrastructure/*` as workspace globs; the commit rule asks a `Ticket:` trailer of any change under `lib/` or `infrastructure/` as it does for `apps/` and `packages/`; the `drift` dead-code scan loops over the new roots; and the context document describes the four workspaces and the naming rule: the `@langdrift/*` scope only under `lib/`, bare names everywhere else. Nothing moves yet, so every check stays green with the globs empty.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `pnpm install` succeeds with `lib/*` and `infrastructure/*` in the workspace globs and no folder under them.
- [ ] A commit touching a file under `lib/` or `infrastructure/` without a `Ticket:` trailer is refused by the commit hook with the same message as for `apps/` and `packages/`; the tooling rule says so.
- [ ] `drift quality dead-code` iterates `lib/*` and `infrastructure/*` and still skips Go modules.
- [ ] The context document's Repository section names the four workspaces, their meaning and the naming rule, and no longer says the package scope is `@langdrift/*` for every package.

## Dev part

**Override:** No part qualifies: mechanical configuration of workspace globs, a commit rule and documentation (policy, Categories: never configuration). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile`
- `platform/cli/drift quality check 02-tooling`
- `platform/cli/drift quality dead-code`
