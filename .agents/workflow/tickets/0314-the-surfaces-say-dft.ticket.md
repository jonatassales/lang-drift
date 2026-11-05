# 0314: the site, the console and the docs say `dft`

**Spec:** .agents/workflow/specs/0047-the-langdrift-cli-answers-to-dft.spec.md

**What to build:** A visitor reading the CLI, Microruntimes or CI Gates page in any of the five locales sees `dft` in every command and terminal demo, and can paste what they copy. The console's connection instructions read `dft init`, `dft agent connect` and `dft agent run`. The root README and the five docs READMEs call the tool the LangDrift CLI and show `dft`, and their terminal screenshot is retaken from the renamed binary. The pull request merges with the 0.5.0 release, so that no visitor reads `dft` while npm still serves `ex`.

**Blocked by:** 0313.

**Status:** ready-for-agent

- [x] No catalog, content module or console view shows `ex` as the command; the catalogs keep parity in every locale.
- [ ] The root README and the READMEs under `packages/docs/` name the LangDrift CLI and show `dft`, and their screenshot shows `dft` running and is named after it.
- [ ] The CI Gates page's two screenshots (`report-light.png`, `report-dark.png`) show a summary that says `dft`.
- [x] `drift check localization` and the website's catalog test pass.
- [x] `pnpm turbo run check typecheck test build --filter=website --filter=console` passes.

> **Note, 2026-11-05.** The copy, the alt texts and the captions say `dft`, and the screenshot file is now `dft-agent-plan.webp` in `apps/website/public/cli/` and `packages/docs/assets/images/`. It still holds the 0.4.0 capture of `ex agent plan`, made in the owner's own terminal against the repository's local records, and no terminal renderer here can retake it faithfully. The owner replaces both files with a capture of `dft agent plan` before this pull request merges with the release. The CI Gates page's `report-light.png` and `report-dark.png` also show the job summary saying `.langdrift/` stays on the machine that ran `ex`, so the owner retakes them from a `dft gates report` run too. Their alt texts already describe `dft`. Those two criteria stay open until then.

## Dev part

**Override:** No part qualifies: the ticket renames a command in copy and pictures without changing what any surface does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-05.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=console`
- `drift check localization`
