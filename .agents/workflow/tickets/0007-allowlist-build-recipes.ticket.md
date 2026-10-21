# 0007: Allowlist the console's build recipes

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** a console-sent build recipe runs only if it is a committed project script (resolved inside the checkout) or an allowlisted tool — docker build, go, npm, pnpm, make — under the existing argument constraints; every interpreter refuses a flag-style first argument. Anything else is refused with a message naming the policy, before any subprocess starts.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] Validation lives in domain/connection with table-driven tests: python3 -c, node -e, sh -c, bare curl all refused; ./build.sh (committed), docker build ., go build, make all pass.
- [ ] A refused recipe surfaces a typed policy error (exit-code work maps it later).
- [ ] native build path proves the script shape resolves within the checkout (no traversal, no symlink escape).
