# General rules

These rules apply to every task, file and agent in this repository. They are loaded in every session. No exception applies unless the owner states it in the task.

1. **Application code must be elegant for a human to read.** Optimize for the next person who reads it, not for the tool that produced it.

2. **Platform and tooling scripts use POSIX shell (`.sh`), never `.mjs`, `.js`, `.ts` or `.mts`.** This covers monorepo automation, workflow execution, and anything that orchestrates the repository itself. Application source code and framework configuration keep their own required languages.
