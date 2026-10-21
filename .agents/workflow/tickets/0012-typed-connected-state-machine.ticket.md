# 0012: Typed state machine for the connected agent

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** the supervisor's per-runtime loop becomes a typed domain state model (declared states and causes, guarded transitions) with table-driven tests at the Execution standard; the service keeps orchestration only. ex agent migrate takes the project lock and parses its runtime argument.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] States/causes are declared types in the domain; illegal transitions return errors; tests cover the table.
- [ ] connected/service.go's runOne shrinks to orchestration over the model.
- [ ] Migrate locks the project and rejects an unparseable runtime id.
- [ ] Behavior parity: existing service tests still pass unchanged where they assert outcomes.
