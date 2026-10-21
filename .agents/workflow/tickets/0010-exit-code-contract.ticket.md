# 0010: Exit-code contract

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** ex exits 0 success, 1 operation failure, 2 usage, 3 policy/guard refusal, 4 missing prerequisite, 130 interrupted — mapped from sentinels in one place and documented in the README. Agent commands' argument errors exit 2 like every other command.

**Blocked by:** 0007.

**Status:** ready-for-agent

- [ ] One mapper in presentation/cli; ErrAdoptionRefused and recipe-policy refusals exit 3; missing tools/cluster exit 4.
- [ ] cli.Execute tests assert each code through injected streams.
- [ ] README documents the contract.
