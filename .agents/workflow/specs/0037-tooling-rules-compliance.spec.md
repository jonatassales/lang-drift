# The tooling complies with the tooling rules

**Status:** ready-for-agent

> **Note, 2026-11-05.** The owner closed 0267, 0270, 0271 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

> **Note, 2026-11-05.** From `@langdrift/cli` 0.5.0 the CLI's command is `dft` (spec 0047). This spec names it `ex`, as it was when the launchers took `.sh`.

## Tickets

- [ ] 0267: the tokens stylesheet is built by shell and jq
- [ ] 0268: the drift and ex launchers take .sh
- [x] 0269: next typegen runs as its own Turborepo task
- [ ] 0270: pre-commit runs the three checks and the commit header has its rules
- [ ] 0271: drift measures the tooling rules

## Problem Statement

The audit of 2026-11-04 measured `.claude/rules/02-tooling.rule.md` at 90.8% while `drift quality check tooling` prints zero findings: the detectors see none of the remaining classes.

- **Automation that is not POSIX shell (rule 1).** `packages/tokens/build.js` builds the tokens stylesheet with Style Dictionary and runs as `node build.js`. The entry points `platform/cli/drift` and `lib/cli/bin/ex` are POSIX shell without the `.sh` extension. The detector only looks under `platform/`, so a `.js` build script in a package is never reported.
- **A package manager other than pnpm (rule 2).** The CI Gates workflow that `ex init` writes into a consumer repository, `lib/cli/src/infrastructure/scaffold/assets/github/langdrift-gates.yml`, uses `actions/setup-node` and runs `npx --yes @langdrift/cli@{{version}} gates report`; its golden copy under `testdata/ci-gates/` matches. Spec 0044 settled that this stays: the tooling rule governs this repository's automation, not a workflow generated into a consumer repository, so it is not a finding here.
- **Chained scripts (rule 3).** `apps/console/package.json` and `apps/website/package.json` define `typecheck` as `next typegen && tsc --noEmit`: two tools in one leaf.
- **Pre-commit (rule 4).** The pre-commit hook runs Biome through lint-staged and `drift quality check`; `drift check localization` runs only in CI.
- **Commit messages (Commit messages).** Of the last 100 non-merge commits, 24 subjects are not lowercase `type(scope): subject` (multi-scope `(web, website)`, capitals such as "MinIO" and "Redis Streams"), 13 subjects exceed 72 characters and 50 body lines exceed 72. The commit-msg hook delegates to commitlint's conventional preset, whose header limit is 100 and whose case rule allows these subjects. History is not rewritten; enforcement belongs in `drift git commit-msg` and `drift git lint-history`.

## Solution

Make every automation path in the repository POSIX shell, every tool invocation pnpm, every package script a single tool, and every commit message the rule's shape, enforced by `drift`: rebuild the tokens stylesheet with a shell script and `jq`, give the `drift` and `ex` launchers their `.sh` extension, move `next typegen` into its own Turborepo task, run `drift check localization` in pre-commit, and add the header rules to the commit guardrails that both the hook and the history lint apply. Each class gains a detector so `drift quality check tooling` keeps measuring it.

## User Stories

1. As the platform owner, I want `drift quality check tooling` to detect each class the audit found, so that zero findings means compliance rather than blindness.
2. As a maintainer of the design tokens, I want the stylesheet built by a shell script and `jq` from the same JSON sources, so that the build obeys the tooling rule and drops a JavaScript dependency.
3. As a consumer of `tokens`, I want the built `website.css` byte-identical before and after the rebuild, so that no app changes.
4. As a tokens author, I want a reference to a missing token or a reference cycle to fail the build and name the token, so that a broken alias never reaches an app silently.
5. As a developer, I want `drift` to keep running as `drift` from anywhere after its launcher becomes `drift.sh`, so that the rename changes nothing I type.
6. As an `ex` user installing from npm, I want `ex`, `npx @langdrift/cli` and `pnpm dlx @langdrift/cli` to keep working after the launcher becomes `bin/ex.sh`, so that the published command name is stable.
7. As a developer running typecheck, I want `next typegen` to run as its own task before `tsc --noEmit`, so that each script invokes one tool and Turborepo orders and caches them.
8. As a developer committing a catalog change, I want pre-commit to run `drift check localization`, so that an invalid catalog or copy in a loader never reaches CI.
9. As a developer committing, I want the hook to reject a subject that is not lowercase `type(scope): subject`, carries more than one scope, exceeds 72 characters or ends with a period, so that history from here on follows the rule.
10. As a reviewer of a pull request, I want `drift git lint-history` to apply the same header rules to every commit in the range, so that `--no-verify` cannot skip them.
11. As a maintainer, I want the hook to leave existing history alone, so that no commit is rewritten and old subjects never fail the lint of a new range.
12. As a team adopting `ex`, I want the CI Gates workflow it writes to stay as spec 0044 left it, with `npx`, so that the generated workflow remains consumer product rather than this repository's automation.
13. As an `ex` maintainer, I want the scaffold's golden test to keep comparing the template byte for byte, so that the generated workflow stays pinned and reproducible.
14. As a reviewer, I want a `.js` or `.ts` build script, an extensionless launcher, a chained package script to be flagged by detectors, so that the rule is checked mechanically.
15. As a documentation reader, I want `.agents/context` and the READMEs to name `drift.sh`, `ex.sh` and the `typegen` task, so that the docs match the tree.

## Implementation Decisions

- **Tokens build: shell plus `jq`.** `packages/tokens/build.sh` replaces `build.js` and the `build` script becomes `sh build.sh`. The script checks for `jq` at the operation and fails naming it, concatenates the sources in the order Style Dictionary read them (the shared files in name order, then the theme), and pipes them through one `jq` program kept in its own `.jq` file beside the script, the way the quality detectors keep their `.awk` programs. The program flattens the DTCG tree (skipping `$type` and `$extensions` keys as names), names a token by its path without its scheme group joined with `-`, resolves a `{path}` alias to `var(--name)` of its target, emits the `:root`, `.dark`, bridge, `@custom-variant` and `@theme inline` blocks in the current order, and exits non-zero naming the token path on a missing reference or a cycle. `style-dictionary` leaves `devDependencies`. The output is proven identical to the current stylesheet before the JavaScript is deleted.
- **Launchers take `.sh`.** `platform/cli/drift` becomes `platform/cli/drift.sh`; `drift setup` links it into the bin directory under the name `drift`, so the command does not change. `lib/cli/bin/ex` becomes `lib/cli/bin/ex.sh`; the package's `bin` map keeps the command `ex` pointing at the new file, `files` lists it, and the launcher test reads it by the new name. Every caller follows: the root `postinstall`, the husky hooks, the workflows, the widgets, `drift doctor`'s required paths and the handbook.
- **`typegen` is a task.** The console and the website gain a `typegen` script (`next typegen`) and `typecheck` becomes `tsc --noEmit`. The root Turborepo configuration gains a `typegen` task that depends on `transit` and caches `.next/types/**`, and `typecheck` depends on `typegen` as well as `transit`; a package with no `typegen` script is skipped by Turborepo. Verification commands name `pnpm turbo run typecheck --filter=<app>` rather than the package script alone, because `tsc` reads the generated route types.
- **Pre-commit runs the three mechanical checks.** `drift git pre-commit` runs lint-staged's Biome, `drift check localization` and `drift quality check`, in that order, on every commit.
- **Header rules in the commit guardrails.** The shared commit rules gain one check that reads a message and accepts only a header of the form `type(scope): subject` or `type: subject`, optionally with the breaking `!` marker, in which type, scope and subject are lowercase, the scope is one kebab-case word, the header is at most 72 bytes and does not end with a period; its error names the violated part. `drift git commit-msg` runs it before commitlint, and `drift git lint-history` applies it to every commit in the range. commitlint keeps the conventional structure, the body line limit of 100 the rule allows, and raises the blank-line-before-body and blank-line-before-trailers rules to errors through the configuration already in the root manifest. A `drift check commit-rules` command runs fixture messages through the guardrails, the way `drift check biome` runs fixtures through Biome.
- **The scaffold template keeps `npx`.** The CI Gates workflow `ex init` writes is consumer product, not this repository's automation, as spec 0044 recorded at the owner's instruction; the template and its golden copy do not change, and the npm/npx and setup-node detectors keep their roots.
- **Detectors.** Automation that is not shell is tracked under `platform/`, `.github/` and `.claude/`, plus any `build.<js|mjs|cjs|ts|mts>` or `scripts/*.<js|mjs|cjs|ts|mts>` in an app or package; an extensionless file directly in `platform/cli/` or in a package `bin/` is flagged; a package script whose value chains tools with `&&`, `||`, `;` or `pnpm run` is flagged. The commit-message rules are verified by the guardrails themselves, which run in the hook and on every pull request's range.
- **`ex` version.** The launcher rename alters the published package; whether they ship as a patch release is the owner's call at review, and `prepack` is unchanged.

## Testing Decisions

- **Seam: `drift quality check tooling`.** New detector rows are proven by running them against the tree before and after each fix; the target is zero findings with every class measured.
- **Seam: the package's own tests.** The tokens contract test keeps building through `pnpm run build` and reading the stylesheet through the package exports, so it covers the shell build unchanged; the build-failure cases (missing reference, cycle, determinism) move to a test of the `jq` program and assert the new error wording. The `ex` launcher test covers the renamed launcher; the CI Gates golden test stays unchanged.
- **Seam: `drift check commit-rules`.** Fixture messages, accepted and rejected, run through the header check, with the error naming the violated part; both the hook and the history lint call the same function, so one fixture set proves both.
- **Behaviour stays green.** `pnpm turbo run build`, `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test`, `drift check localization`, `drift check biome` and `drift check workflows` pass after every ticket; a Vercel preview build of the website succeeds after the tokens build changes.
- **Review-enforced.** The body wrap at 72 (the linter allows 100), the "ideally 60" subject length, and whether a `build.sh` is the tool for its operation rather than a wrapper stay with the review.

## Out of Scope

- Rewriting any commit already in history.
- Replacing the `.awk` programs under the quality detectors, or the `.jq` program this spec adds, with shell: they are programs for purpose-built POSIX tools, which rule 1 asks for.
- `lib/cli/build.sh` as the `build` script: it is the cross-compilation tool for the package, not a wrapper around one framework command.
- The `npx @langdrift/cli` lines in the `ex` README and the handbook: consumer documentation for a public npm package, which names both `npx` and `pnpm dlx`.
- Making `drift quality check` exit non-zero on findings (Further Notes).
- Ignoring `.audits/` (the testing compliance spec, not yet written) and the Biome import matcher (spec 0039).

## Further Notes

- **Lowercase, strictly.** The rule says the subject is lowercase, so the header check rejects any capital, including proper nouns: "archive every event into minio" passes and "into MinIO" does not. The owner may relax this to a lowercase first word at review; the check is one pattern.
- **One scope.** `(web, website)` is read as two scopes and rejected; a commit that spans surfaces names the leading one or uses no scope.
- **72 in the hook, 60 by hand.** The hook enforces the hard limits the rule states (72 for the header, 100 for body lines, as the rule says the linter allows); the ideal 60 and the 72 wrap stay with the review.
- **Open question: should `drift quality check` fail?** It lists findings and never fails, so pre-commit and CI run it without gating. Once the compliance specs reach zero findings, the owner may want it to exit non-zero; no spec owns that switch yet.
- **`jq` at build time.** `jq` is already required by `drift check localization` and the widgets, and GitHub's runner image ships it. The Vercel build image must be confirmed to ship it before the tokens ticket merges; if it does not, the alternative is a Go program in a dedicated package, which the rule allows as a purpose-built tool, and the ticket records the outcome.
- **The scaffold template as consumer product.** Spec 0044 kept `npx` in the workflow `ex init` writes; this spec follows it and leaves the scaffold assets out of the detectors' roots.
- **`typecheck` alone.** After the split, `pnpm --filter console typecheck` runs `tsc` without regenerating route types; Turborepo runs the chain. Tickets and docs cite the Turborepo form.
- **`.sh` on a published launcher.** npm accepts any path as a `bin` target, and the installed command stays `ex`; the extension is visible only inside the package.
