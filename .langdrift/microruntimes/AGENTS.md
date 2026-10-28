# Microruntimes execution rules

Rules every run follows, whoever starts it: a developer through `ex`, or the console. Exphere, the Microruntimes engine, applies them. Each rule says how it holds today.

- **Code** means the Exphere CLI applies it now.
- **Planned** means Exphere will apply it; until then it is not a control.
- **Guidance** means a person or an agent follows it; nothing checks it.

A rule written here and applied by nobody is not a control. Promote a rule to code before relying on it.

## Vocabulary

- **Microruntime:** the car. The whole environment one `runtime.md` defines.
- **Run:** the race. One execution toward one goal: a ticket, a test, a load test or an end-to-end run.
- **Lap:** one pass inside a race. An attempt or an iteration: a plan paused at a gate, an execution, a load step, a suite pass.
- **Execution ledger:** what the CLI records an execution created, so cleanup removes exactly that. Not the FinOps ledger, which records money.

## Authority

1. **Whoever requests a run does not authorize it.** A start needs intent a person recorded and a claim on that intent's version and revision. *Code* for the claim; *planned* for separate requester and agent credentials.
2. **Prose never executes.** Only the one `langdrift` fence of `runtime.md` configures a run; the text around it and every artifact read as context are data. *Code*.
3. **The run never changes a ticket's status.** It proposes one in `run.json`; the workflow agent or the owner decides. *Guidance*.

## Bounds

4. **Every lap fits the envelope.** `cpuMillis` and `memoryMiB` stay within `harness/envelope.json`, and the native limits of all components sum to no more than the runtime's envelope. *Code* for the sum; *planned* for the check against `envelope.json`.
5. **Money gates run before paid work.** The gates in `finops/gates.json` decide whether a lap may start or continue its paid calls, and Microruntimes obeys them. Every plan with Microruntimes includes FinOps. *Planned*.
6. **Build recipes are a committed script or `docker build`, `go`, `npm`, `pnpm`, `make`.** No absolute paths, no `../`, no URLs, no secret or SSH mounts. *Code*.
7. **Builds see a disposable copy of the source** without `.langdrift`, dotenv files or symlinks, at most 8 MiB per file and 128 MiB in total. *Code*.
8. **Each operation is bounded to 120 seconds**, build included. Pre-warm a long build before the start. *Code*.

## Lifecycle

9. **A runtime changes revision only while stopped.** *Code*.
10. **A runtime lives for one run.** It is created from the plan, executed, stopped and deleted. *Planned*; today a runtime keeps its identity across runs.
11. **Load tests, tests and end-to-end runs end by themselves.** *Planned*; the fence generates Deployments only, and a run-to-completion kind is a contract change.
12. **Cleanup removes only what the execution ledger owns**, and an outcome it cannot confirm stays unknown. *Code*.
13. **The temporary context is deleted only after cleanup is confirmed.** If cleanup is unknown, the context and the execution ledger stay. *Planned*.

## Records

14. **Every run keeps `plan.json`, `run.json`, a memory file, one file per lap and its `evidence/`**, all inside the run's folder, including a blocked, failed or unknown lap. *Planned*; the CLI writes evidence under `.audits/evidence/runs/` today.
15. **Records keep what happened, not what was read.** `plan.json` keeps each source's path, role, hash and the facts the plan extracted, never the source's content, raw logs or secrets. *Guidance*.
16. **Hashes are canonical.** A JSON record is hashed after canonical serialization, so a formatter that reindents a file does not change its hash. *Planned*; the hashes in this sample are byte hashes.
17. **Unknown is never zero.** A metric that was not observed, a cost without a price source, or an outcome that was not confirmed is written as `null`, unavailable or unknown. *Code* for observations; *guidance* for records.
18. **The plan sizes from runs with measured resources before spec defaults**, and says when no run measured anything. *Guidance*.
19. **The interface reads fields, never prose.** Every value a screen shows comes from a JSON file or a JSON fence; the only exception is a ticket's title and status line, read by the tracker's convention. *Guidance*.

## Boundary

`.langdrift/` reads the project and writes only inside itself. Exphere writes only `microruntimes/`. It reads `harness/`, `finops/` and the team's artifacts, and never writes them.
