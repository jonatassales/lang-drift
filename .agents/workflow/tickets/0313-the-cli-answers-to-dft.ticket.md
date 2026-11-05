# 0313: the CLI answers to `dft`

**Spec:** .agents/workflow/specs/0047-the-langdrift-cli-answers-to-dft.spec.md

**What to build:** A developer who installs `@langdrift/cli` types `dft`. Help, usage and error messages say `dft`. The release binaries are `dft-<os>-<arch>`. Files the CLI writes into a project call `dft`, and new records name `dft` as their provenance. `ex` no longer resolves, and a record written by 0.4.0 with `ex` still reads. The widgets call the tool the LangDrift CLI, and the READMEs and comments of `lib/cli` and `lib/core` say `dft`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] `package.json` maps `bin` to `dft` alone; `pnpm pack --dry-run` lists `bin/dft.sh` and no `bin/ex.sh`.
- [x] `dft --help` names `dft` in its usage; the launcher test runs `bin/dft.sh`.
- [x] The main workflow builds `dft-<os>-<arch>` and their `SHA256SUMS`; the widgets say "LangDrift CLI".
- [x] The golden fixtures hold `dft` as provenance, and a test reads a 0.4.0 record that names `ex`.
- [x] `GOTOOLCHAIN=local go vet ./...` and `go test -race ./...` pass in `lib/cli` and `lib/core/contracts/go`.
- [x] `git grep` finds `ex` as the command in `lib/`, `.github/` and the widgets only in 0.4.0 fixtures.

> **Note, 2026-11-05.** Fixtures that stand for a project an earlier version wrote keep `ex`. That covers `testdata/langdrift/project/`, its guardrail, which the expected tree mirrors (the CLI ships no guardrail of its own and leaves the team's in place), and the new `testdata/langdrift/v0.4.0/` ledger that `TestLedgerOfVersion040Decodes` reads. The CI Gates golden copy is generator output and says `dft`. The host fallback is now `dist/dft`, and `LANGDRIFT_SKIP_EX_BUILD` is `LANGDRIFT_SKIP_DFT_BUILD`.

## Dev part

**Override:** No part qualifies: the ticket renames a command without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-05.

## Verification

- `GOTOOLCHAIN=local go -C lib/cli vet ./...` and `go -C lib/cli test -race ./...`
- `pnpm --filter @langdrift/cli exec pnpm pack --dry-run`
