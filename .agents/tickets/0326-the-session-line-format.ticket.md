# 0326: the session's line format

**Spec:** .agents/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** `dft agent plan` and `dft agent run` print their lines the way the simulation lays them out.
- The typed command is never echoed.
- Every actor line starts with `console › `.
- The plan opens on its new scope comment, and its three acts open on their comments.
- Labels line up in the simulation's columns.
- On a terminal, comments are dimmed, success lines green and warnings yellow, and the window title reads `dft — project`. Piped, the bytes are plain text.

Tests read each printed line's kind from an injected session sink, not from `--json`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** change and build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] No printed line repeats the typed command: neither `dft agent plan …` nor `dft agent run …` is printed.
- [ ] Every actor line starts with `console › `, as in `console › tool call: read 0031-TICKET-book-a-demo-runtime`. No line starts with `agent › `.
- [ ] The plan's first line is `# plan only · nothing is built, started or paid for in this step`.
- [ ] The four comments sit at their moments:
  - `# analysis · reading the linked documents, then the project's episodes` after the links line;
  - `# plan · resolving the components, then sizing them inside the envelope` after `✓ suggests   claude-sonnet-5-5 · effort medium`;
  - `# next: once the team authorizes, dft agent run applies the recorded intent` last.
- [ ] Every line takes one of spec 0048's six shapes, laid out by one function that counts characters (`›`, `·`, `✓` and `—` count one each):
  - **comment, narrative actor line:** as written;
  - **labeled actor line, result line, status line:** the label padded to 10 columns, then one space;
  - **component row:** the path padded to 19 columns, three spaces, then `<stack> · <role>` padded to 15 columns, then one space.

  So `console › to-tickets 3 tickets · blocking edges recorded` keeps one space, `✓ pr #19     feat(demo): run book a demo as a microruntime` reaches its value at the same column as `✓ read       …`, and `✓ apps/signup           Next.js · web   500m · 1024 MiB` aligns with `✓ services/signup-api   Go · api        250m · 256 MiB`.
- [ ] Each line's kind is the simulation's, keyed by the website's line ids:
  - `limits     gpu and disk are not runtime.md fields · none requested` is a warning;
  - `console › gate       reserved $11.40 of $25.00 · opened` is a success line;
  - `console › gate       reservation $40.80 exceeds the $25.00 budget · paused` is a warning;
  - `✓ …` lines are success and `# …` lines are comments;
  - every other line is output.
- [ ] On a terminal, a comment is dimmed, a success line green and a warning yellow, and an output line keeps the default color. When standard output is not a terminal, no escape sequence is written.
- [ ] Plan and run set the terminal title `dft — <project folder>` (`dft — project` for the fixture) when standard output is a terminal, and only then.
- [ ] The session sink is a port in the CLI's dependencies, so a test reads each line's id, kind and text.

## Dev part

**Category:** `hard-part`. Every element spec prints through this layout. Counting characters rather than bytes, padding labels shorter than, equal to and longer than 10 columns, and the component row's three columns are where the ticket's difficulty sits.

**Paths:**
- `lib/cli/src/presentation/cli/session/layout.go`
- `lib/cli/src/presentation/cli/session/layout_test.go`

**Boundary:** `func Layout(shape Shape, cells Cells) string`. `Shape` names the six shapes (`Comment`, `Narrative`, `Labeled`, `Result`, `Status`, `ComponentRow`). `Cells` carries `Text`, `Label`, `Value`, `Path`, `StackRole` and `Share`. The agent declares both, each in its own file, and the catalog calls `Layout`.

**Goal:** Lay out one line of each shape as spec 0048's table says, counting characters, with no trailing space.

- [ ] `Labeled` with `gate` gives `console › gate       …`.
- [ ] `Labeled` with `to-tickets` (10 characters) and `code-review` (11) each give exactly one space before the value.
- [ ] `Result` with `pr #19` and with `read` reach the value at the same column.
- [ ] `ComponentRow` gives `✓ services/signup-api   Go · api        250m · 256 MiB` and `✓ apps/signup           Next.js · web   500m · 1024 MiB` byte for byte.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `Layout` and its test.
