# General rules

These rules apply to every task, file and agent in this repository. They are loaded in every session. No exception applies unless the owner states it in the task.

1. **Application code must be elegant for a human to read.** Optimize for the next person who reads it, not for the tool that produced it.

2. **Platform and tooling scripts use POSIX shell (`.sh`), never `.mjs`, `.js`, `.ts` or `.mts`.** This covers monorepo automation, workflow execution, and anything that orchestrates the repository itself. Application source code and framework configuration keep their own required languages.

3. **pnpm is the only package manager.** Install, run and set up tooling through pnpm, locally and in CI: `pnpm exec` or `pnpm dlx`, never npm, npx, Yarn, Bun or Corepack. In GitHub Actions, `pnpm/setup` installs pnpm and Node.js together, so workflows do not use `actions/setup-node`.

4. **Mandatory two-axis `code-review` before every commit.** Every code modification must execute the `code-review` skill against the spec and against all repository rules in `.claude/rules/*.rule.md`. The agent must output its findings under `## Standards` and `## Spec` before committing or delivering. Any breach of documented rules or spec requirements must be fixed before delivery; no commit is made with unaddressed findings.
