# 0314: the site, the console and the docs say `dft`

**Spec:** .agents/workflow/specs/0047-the-langdrift-cli-answers-to-dft.spec.md

**What to build:** A visitor reading the CLI, Microruntimes or CI Gates page in any of the five locales sees `dft` in every command and terminal demo, and can paste what they copy. The console's connection instructions read `dft init`, `dft agent connect` and `dft agent run`. The five docs READMEs call the tool the LangDrift CLI and show `dft`, and their terminal screenshots are retaken from the renamed binary.

**Blocked by:** 0313.

**Status:** ready-for-agent

- [ ] No catalog, content module or console view shows `ex` as the command; the catalogs keep parity in every locale.
- [ ] The five READMEs under `packages/docs/` show `dft`, and their screenshots show `dft` running and are named after it.
- [ ] `drift check localization` and the website's catalog test pass.
- [ ] `pnpm turbo run check typecheck test build --filter=website --filter=console` passes.

## Dev part

**Override:** No part qualifies: the ticket renames a command in copy and pictures without changing what any surface does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-05.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=console`
- `drift check localization`
