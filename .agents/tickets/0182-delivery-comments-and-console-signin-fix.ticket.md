# 0182: Split delivery tables into separate comments and fix console sign-in

**What to build:**
1. In delivery widgets (`.github/widgets/dashboard.sh`, `workflow.sh`, `finops.sh`, `microruntimes.sh`, `code-quality.sh`, `agentic.sh`):
   - Remove the LangDrift logo from all table headings and rows.
   - Remove all emojis from all gate tables.
   - Link all status column cells to `https://console.langdrift.md/sign-in`.
   - In `dashboard.sh`, publish each gate table as its own individual sticky PR comment with dedicated keys (`langdrift-ci-gates`, `langdrift-workflow`, `langdrift-finops`, `langdrift-microruntimes`, `langdrift-code-quality`, `langdrift-agentic`).
2. In `apps/console`:
   - Fix the broken `/sign-in` route by adding `/[locale]/sign-in` redirecting to `/login` and configuring redirects in `next.config.ts` for `/sign-in` and `/signin`.

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] Remove LangDrift logo and emojis from `.github/widgets/dashboard.sh` and gate widgets.
- [x] Add `https://console.langdrift.md/sign-in` links to all status column cells across all tables.
- [x] Update `dashboard.sh` to publish each table as a separate sticky comment on PRs.
- [x] Add `/sign-in` page and redirects in `apps/console`.
- [x] Verify all checks and tests pass with `pnpm turbo run check typecheck test build`.
