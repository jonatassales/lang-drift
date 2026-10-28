# 0001-TICKET: Bring the console preview up as a microruntime

**Status:** ready-for-agent
**PRD:** 0001-PRD-console-preview.prd.md
**Design doc:** 0001-DD-console-microruntime.dd.md
**Spec:** 0001-SPEC-console-microruntime.spec.md
**Guardrail:** 0001-finops.gr.md
**Blocked by:** None

**What to build:** run the console with extub as the microruntime `console-local`, inside the project's envelope and budget, and record the run.

- [ ] The plan proposes a `runtime.md` inside `harness/envelope.json`, and a person authorizes it.
- [ ] The run meets the spec's acceptance criteria.
- [ ] The run records its laps, the observed resources and anything the next run should know.
