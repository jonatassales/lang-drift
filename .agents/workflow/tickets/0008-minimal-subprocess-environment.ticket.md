# 0008: Minimal environment for every subprocess

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** every child process ex starts receives an explicit allowlist environment (PATH, HOME, TMPDIR, locale, plus its tool family's declared needs) instead of inheriting the agent's. Hooks receive only variables the project's hook configuration declares. Every git call runs with GIT_* scrubbed. Credential stripping becomes allowlist-based.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] One environment-builder helper in infrastructure/process, used by hooks, terraform, kind, kubectl, git and build recipes.
- [ ] Recorded-runner tests assert PROJECT_SECRET and AWS_SECRET_ACCESS_KEY never reach a child that did not declare them.
- [ ] Terraform keeps its existing cleaned TF_* policy, now including TF_CLI_CONFIG_FILE and TERRAFORMRC.
- [x] Hooks outcome, adjusted during implementation: hooks are reachable only through the retired runtime tree (deleted by ticket 0015), so no declaration schema was added yet. Hooks run on the minimal base plus the LANGDRIFT_* correlation variables, and the declaration syntax arrives with the first live hook surface.
