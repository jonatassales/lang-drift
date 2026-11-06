# 0329: the fixture `project` and the golden test

**Spec:** .agents/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** One in-process golden test plays the session against a frozen copy of the simulation's project, and compares what `dft` prints with the simulation, line by line and kind by kind. Every later spec adds its fixture content and its assertions to this test.

A line whose spec has not landed is listed with its owning spec, so the distance to the simulation is measured: the list shrinks to empty as the specs land.

**Blocked by:** 0326.

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

**Owner exception (2026-11-05):** the fixture reproduces the client project with the simulation's names, under an exception to filesystem rules 1 and 2 that the owner stated for this fixture only. The uppercase parts and the `.dd.md` and `.gr.md` suffixes are part of what the simulation prints (`0011-SPEC`).

- [ ] **The fixture.** A project folder named `project` in the CLI's test data, so the title reads `dft — project`. It holds:
  - in `.agents/`:
    - `0031-TICKET-book-a-demo-runtime.ticket.md`, with the heading `Run the book-a-demo request use case as a microruntime`, `**Status:** ready-for-agent`, and links to the three documents and the guardrail;
    - `0002-PRD-book-a-demo.prd.md`, with `**Use case:** book-a-demo request`;
    - `0005-DD-book-a-demo-runtime.dd.md`, with `**cpuMillis:** 2000`, `**memoryMiB:** 3072` and `**Budget per run:** $25.00`;
    - `0011-SPEC-book-a-demo-runtime.spec.md`, whose `## Components` lists the five paths in order (`services/signup-api` with `:8080`), and whose `## Acceptance criteria` has three items;
  - `.langdrift/harness/guardrails/0001-finops.gr.md`, with FIN-02 (`plan, run`), FIN-03 (`run`) and FIN-04 (`plan`);
  - five components with no Dockerfile:
    - `services/signup-api/go.mod`;
    - `services/notify-api/package.json`, with `PORT=8081`;
    - `platform/nats/nats.conf`, with `port: 4222`;
    - `platform/redis/redis.conf`, with `port 6379`;
    - `apps/signup/next.config.ts` beside a `package.json`;
  - one Git commit and a remote, made by the test on its copy.

  The episodes, the index and `driver.json` join with specs 0052 and 0058, as spec 0048 lets each spec add its content.
- [ ] **The goldens.** Frozen in the CLI's test data, transcribed from the website's English catalog with values filled:
  - `plan`: 42 lines, from `# plan only · nothing is built, started or paid for in this step` to `# next: once the team authorizes, dft agent run applies the recorded intent`;
  - `run`: 25 lines with Sonnet 5.5 at medium, Planning and Review checked, from `# run · 0031-TICKET-book-a-demo-runtime · the recorded intent, authorized by the team` to `# next: the owner reviews #19 and decides whether to merge`;
  - `run-over-budget`: 6 lines with Opus 5.5 at high, ending on `console › gate       reservation $40.80 exceeds the $25.00 budget · paused` and `# needs a decision: raise the budget or lower the effort · nothing was paid for`;
  - `run-to-merge`: 25 lines with Review unchecked, ending on `✓ ready      #19 to merge · $1.60 released to the envelope` and `# next: the merge stays the owner's`;
  - `run-all-stages`: 31 lines, each checked stage's `waiting`/`confirmed` pair before its work, and Implement's before `console › hook       effort recorded · implement unblocked`;
  - `run-haiku-low`: 25 lines with Haiku 5.5 at low, with `console › estimate   $0.26 · 2M input · 120K output · list prices · FIN-02`, `reserved $0.32`, `$0.27 observed` and `$0.05 released`.

  Each line carries its kind.
- [ ] **The test.** It runs `dft` in process on a fresh copy of the fixture: `dft init`, `dft agent plan`, the intent's authorization, then `dft agent run`. Until 0330 to 0333 land, it seeds through today's commands and the fake console.
- [ ] **This spec's own assertions:**
  - the four comment lines at their positions;
  - no line repeats the typed command;
  - every actor line starts with `console › `;
  - every line takes one of the six shapes;
  - every line's kind is the golden's.
- [ ] **The hand-off records.** The test reads `plan.json` after the plan and `intent.json` and `run.json` after the run, and compares them field by field with the four cards. Until spec 0050's tickets write them, that comparison is listed as pending under spec 0050.
- [ ] **The comparison.** Every golden line is compared with the printed line at its position. The pending list names each line that may still differ and the spec that owns it (0050 to 0066). The test fails when an unlisted line differs, and when a listed line already matches, so the list only shrinks.

## Dev part

**Category:** `hard-part`. Every element spec proves itself through this comparison. Aligning printed and golden lines, keeping a pending line from hiding a regression, and forcing the list to shrink are where the ticket's difficulty sits.

**Paths:**
- `lib/cli/src/golden-diff_test.go`
- `lib/cli/src/golden-diff-rules_test.go`

**Boundary:** `golden-diff_test.go` holds `goldenDiff`, and `golden-diff-rules_test.go` its rules. The signature is `func goldenDiff(golden []goldenLine, printed []printedLine, pending map[string]string) []string`, called by the golden test. The agent declares `goldenLine` and `printedLine` (id, kind, text), each in its own file.

**Goal:** Return every disagreement between the golden and the printed session, as a readable message:
- a line that differs and is not pending;
- a pending line that already matches;
- a missing or extra line;
- a kind that differs.

- [ ] Identical sessions with an empty pending list give no message.
- [ ] A differing unlisted line, a matching listed line, a missing line, an extra line and a different kind each give one message naming the line.
