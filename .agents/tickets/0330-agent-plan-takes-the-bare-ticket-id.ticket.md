# 0330: `dft agent plan` takes the bare ticket id

**Spec:** .agents/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** A person types `dft agent plan 0031-TICKET-book-a-demo-runtime` with the ticket's bare id, as the simulation does, from any folder inside the project. The plan reads what it needs from the ticket and its links instead of from flags.

**Blocked by:** 0329, and the tickets of specs:
- 0053, which resolves the id to the project's one ticket file and reads its links;
- 0054, which takes the envelope from the design doc;
- 0055, which takes the runtime id from the ticket's slug;
- 0066, which takes the components from the spec.

**Status:** ready-for-agent

**Alignment:** change and remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] `dft agent plan` takes exactly one argument, the bare id. None, or more than one, is a stop (0327).
- [ ] The id goes to spec 0053's resolution, for `console › tool call: read 0031-TICKET-book-a-demo-runtime`.
- [ ] The project-relative path as the argument leaves.
- [ ] These flags leave: `--runtime`, `--name`, `--component`, `--cpu-millis`, `--memory-mib`, `--max-depth`, `--max-files`, `--max-bytes`, `--no-judge` and `--by`.
- [ ] The golden test plays the plan with the bare id, and its pending list drops no line it had matched.

## Dev part

**Override:** No part qualifies: an argument count and the removal of flags, with the resolution itself in spec 0053's tickets (policy, Categories: never configuration or a rename). 2026-11-05.
