---
status: accepted
---

# The console's build recipes are allowlisted and subprocesses get a minimal environment

The connected agent executes build recipes that arrive from the console, so a compromised or misbehaving console is code execution on the user's machine. ex therefore refuses every recipe that is not either a script committed inside the project checkout or an allowlisted build tool (`docker build`, `go`, `npm`, `pnpm`, `make`) under the existing argument constraints, and refuses a flag-style first argument for every interpreter. Blocklisting was rejected: it is what previously allowed `python3 -c` and `node -e` through.

For the same reason every subprocess ex starts — hooks, terraform, kind, kubectl, git, build recipes — receives a minimal explicit environment (`PATH`, `HOME`, `TMPDIR`, locale, plus what that tool family needs) instead of inheriting the agent's, which holds `PROJECT_SECRET`. Hooks see only the variables the project's hook configuration declares, so a hook that needs `AWS_PROFILE` must say so.

## Consequences

- A project whose build genuinely needs another tool wraps it in a committed script, which code review then sees.
- Hooks relying on ambient variables stop working until the variables are declared. This is intended: the declaration is the audit trail.
