# 0011: Human output where humans read

**Spec:** .agents/workflow/specs/0002-ex-hardening.spec.md

**What to build:** doctor, cluster inspect and version print human-readable text by default; --json keeps the machine form; every other command keeps JSON. Errors render as one line of cause plus a hint line. Adjusted during implementation: runtime inspect lives only in the retired tree that ticket 0015 deletes, so it keeps JSON.

**Blocked by:** 0005.

**Status:** ready-for-agent

- [ ] A presentation renderer over the same data the JSON path emits; no logic in templates.
- [ ] Golden-style tests for the four commands in both modes.
- [ ] Error rendering shows cause + hint; JSON error shape unchanged.
