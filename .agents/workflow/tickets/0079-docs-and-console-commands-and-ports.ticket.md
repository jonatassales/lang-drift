# 0079: The docs and the console run on their own ports

**Spec:** .agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md

**What to build:** `pnpm docs` serves the docs at `http://localhost:3002/docs`. `pnpm console` starts the console on 3001 together with the stub on 3401. `pnpm website` keeps 3000, so all three run at once. The commands, their ports and the preview launch configurations are documented, and the rule on root launchers allows three.

**Blocked by:** 0076

**Status:** ready-for-agent

- [ ] The three commands run at the same time, and each app answers on its port.
- [ ] The console shows the stub's data without starting the stub separately.
- [ ] AGENTS.md, the context file and the docs app name each command and port.

### Context

- [The spec](/.agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md)
