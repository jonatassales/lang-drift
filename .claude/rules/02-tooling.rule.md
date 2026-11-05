# Tooling rules

Scripts, package management, mechanical checks and commits. Loaded in every session.

1. **POSIX shell for automation.** Repository automation, orchestration and verification scripts, including CLI helpers and their script-based checks, are POSIX shell with a `.sh` extension. `.mjs`, `.mts`, `.js`, `.ts` and other non-shell automation scripts are prohibited. A `.sh` wrapper around embedded JavaScript, TypeScript or Python does not satisfy this rule; call purpose-built command-line tools for parsing and other operations. Application source and framework configuration keep their required languages and file conventions. When touching an existing automation script in another language, migrate it to `.sh`; unrelated legacy scripts do not authorize new exceptions.
2. **pnpm is the only package manager.** Install, run and set up tooling through pnpm (`pnpm exec`, `pnpm dlx`), locally and in CI. Never `npm`, `npx`, `yarn`, `bun` or `corepack`. In GitHub Actions, `pnpm/setup` installs pnpm and Node.js together, so workflows do not use `actions/setup-node`.
3. **Atomic leaf scripts.** A leaf script invokes the tool for that operation (`"build": "next build"`), and the tool's own flags are part of the command (`next build --webpack`). A script does not call another script, and no `scripts/build.sh`, `scripts/dev.sh` or `scripts/start.sh` wraps one framework command. The repository root uses Turborepo to coordinate packages; it calls atomic leaf commands, with no chain of package scripts and shell wrappers in front of the tool.
4. **Mechanical checks.** Biome (`pnpm turbo run check`), `drift quality check` and `drift check localization` run in pre-commit and CI. A rule they do not cover is enforced by review (General, Authority 7).

## Commit messages

- **Conventional Commits:** `type(scope): subject`, in lowercase.
- **Subject:** imperative mood, at most 72 characters (ideally 60), no trailing period.
- **Blank lines:** exactly one between the subject and the body, and between the body and the trailers.
- **Body:** hard-wrapped at 72 characters (the linter allows at most 100).
- **Ticket trailer:** a commit touching `apps/`, `packages/` or `lib/` ends with `Ticket: <NNNN>` after a blank line.
- **No AI signatures:** no `Co-Authored-By`, `Generated with` or assistant markers in commits or pull requests.
