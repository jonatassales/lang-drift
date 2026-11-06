# 0327: the stop rule

**Spec:** .agents/workflow/specs/0048-the-langdrift-cli-runs-init-plan-and-run.spec.md

**What to build:** Any failure, refusal or interruption the simulation does not show ends the command in one way:
- no further session line is printed;
- one plain sentence naming the cause goes to standard error;
- the command exits 1.

The over-budget pause is a shown ending and exits 0. A person or a script sees that the command did not finish, and standard output stays a prefix of the simulation.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** change and remove. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] A stop writes one sentence to standard error: no `dft:` prefix, no `hint:` line, no JSON and no secret.
- [ ] A stop exits 1, whatever the cause. Today's exit statuses 2 (usage), 3 (refused), 4 (missing prerequisite) and 130 (interrupt) become 1. An interrupt (Ctrl-C) is a stop.
- [ ] A stop prints no session line of its own. The failure lines the simulation never shows leave: the plan's blocked line, the run's blocked gate, the exhausted gate, the failed checks and the unknown runtime. The lines already printed stay as they were.
- [ ] Nothing is retried and nothing resumes inside the command. If the run had started the runtime, it is taken down without a line.
- [ ] A command that ends on its last simulation line exits 0.
- [ ] The over-budget run ends on `console › gate       reservation $40.80 exceeds the $25.00 budget · paused` and `# needs a decision: raise the budget or lower the effort · nothing was paid for`, and exits 0.
- [ ] A scripted failure after some lines leaves standard output exactly as printed up to that point, writes one sentence to standard error and exits 1.

## Dev part

**Category:** `core-domain`. Spec 0048 states one rule for every failure the simulation does not show. What the sentence says, that it never leaks a secret or a wrapped chain, and which endings are not stops: an owner must be able to explain all of it during an incident.

**Paths:**
- `lib/cli/src/presentation/cli/stop.go`
- `lib/cli/src/presentation/cli/stop_test.go`

**Boundary:** `func Stop(err error, stderr io.Writer) int`, called by the CLI's entry point where today's exit-code mapping and error report are. The over-budget pause reaches it as the run's existing pause sentinel, which the agent keeps distinct from every stop, outside these paths.

**Goal:** Turn the error that ended a command into its exit status, and, for a stop, its one sentence.

- [ ] `nil` and the over-budget pause return 0 and write nothing.
- [ ] Every other error, an interrupt included, returns 1 and writes one line: a plain sentence naming the cause, with no prefix, no hint and no secret.

**Override:** on 2026-11-06 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `Stop` and its test.
