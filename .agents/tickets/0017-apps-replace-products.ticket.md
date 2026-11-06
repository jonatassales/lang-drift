# 0017: Apps replace products

**Spec:** .agents/specs/0003-langdrift-website-on-apps.spec.md

**What to build:** the repository's workspaces live under `apps/` and `packages/`. The empty `products/` folders are gone, the ticket-trailer guardrail and `drift doctor` speak of `apps/`, and the docs and glossary describe apps instead of products.

**Blocked by:** none.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The pnpm workspace picks up `apps/*` and `packages/*`; `products/` no longer exists.
- [x] A commit that changes `apps/` or `packages/` without a `Ticket:` trailer is rejected, by the commit-msg hook and by `drift git lint-history`.
- [x] `drift doctor --ci` passes and requires `apps`.
- [x] AGENTS.md (status, layout, design, commit hook) and the glossary use apps; the planned `ui` package is named `web`.
