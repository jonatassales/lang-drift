# 0102: The run drives the workflow stages and waits for a person

**Spec:** .agents/workflow/specs/0015-ex-performs-the-microruntimes-session.spec.md

**What to build:** after the gate opens, the run plays the workflow's stages in order: a stage marked human prints the waiting line and pauses, with a prompt on a terminal and otherwise an exit that `ex agent confirm <runtime> <stage>` resolves and a second `ex agent run` resumes; each stage's skills run through the coding-agent driver the owner configured, as a trusted process with placeholders substituted per argument, a minimal environment and a bounded log, and print their one-line summary; a stage's declared hooks run first as committed scripts; the money gate runs before every paid call, pausing when observed spend reaches the reservation or usage is unknown; usage reaches the ledger as observed tokens and cents; the owner-only merge never runs; the completed line prints observed against reserved. Stage states, confirmations and summaries land in the run record, and a token-free evidence file per call in the runtime's evidence.

**Blocked by:** 0101.

**Status:** ready-for-agent

- [ ] With a fake driver returning the website's summaries and 3.7M input and 240K output tokens, lines 7 to 14 and 16 print with the sample's figures, and the run record's `stages` and the ledger's observed figures equal the expected fixture.
- [ ] A human stage on a terminal prompts and continues on confirmation; without a terminal the run exits 0 as waiting, `ex agent confirm` records the confirmation, and a second run resumes at the next stage.
- [ ] The driver is tested with a recording runner: the arguments, the passed environment and the timeout are the configured ones, usage is read from the configured fields, and a missing usage pauses the next paid call unless the configuration allows unmeasured calls.
- [ ] `merge` stays pending and is never executed.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; a run on this repository with the `claude` driver configured confirms planning on the prompt and prints the stage lines.
