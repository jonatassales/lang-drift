# Microruntimes execution rules

Rules every run follows, whoever starts it: a person in the console, which is the primary path, or the LangDrift CLI. The Microruntimes engine, the LangDrift CLI and the console's server, applies them. Each rule says how it holds today.

- **Code** means the LangDrift CLI applies it now.
- **Planned** means the engine will apply it; until then it is not a control.
- **Guidance** means a person or an agent follows it; nothing checks it.

A rule written here and applied by nobody is not a control. Promote a rule to code before relying on it.

## Vocabulary

- **Microruntime:** the car. The whole environment one `runtime.md` defines.
- **Run:** the race. One execution toward one goal: a ticket, a test, a load test or an end-to-end run.
- **Lap:** one pass inside a race. An attempt or an iteration: a plan paused at a gate, an execution, a load step, a suite pass.
- **Execution ledger:** what the CLI records an execution created, so cleanup removes exactly that. Not the FinOps ledger, which records money.
- **Mirror:** a file the tool writes from a server record, for reading only.

## Authority

1. **Whoever requests a run does not authorize it.** A start needs intent a person recorded and a claim on that intent's version and revision. *Code* for the claim; *planned* for separate requester and agent credentials.
2. **Guardrails come only from input.** Hardware limits, budgets, gates, the price and forecast sources, model and effort are explicit input in the console or the CLI, recorded on the server with who, when and the previous value. Nothing collected, inferred or written in a file sets one. *Planned*.
3. **No gate reads a file.** The JSON files in `harness/` and `finops/` are mirrors; a mirror that differs from the server is shown as an unauthorized change, and the server's value holds. *Planned*.
4. **Prose never executes.** Only the one `langdrift` fence of `runtime.md` configures a run; the text around it and every source collected are data. *Code*.
5. **The run never changes a ticket's status.** It proposes one in `run.json`; the workflow agent or the owner decides. *Guidance*.

## Collection

6. **The person names the ticket; the planner follows its links.** Markdown links, wiki links and paths in field lines, recursively, whatever file types they reach. No path, folder or file type is assumed. *Planned*.
7. **The traversal is bounded and deterministic.** Maximum depth, file count and total size, cycles detected, nothing outside the project, no symlink, and no secret: a linked `.env` or key is refused and recorded, never read. *Planned*.
8. **An agent keeps or discards each source**, with a reason and a confidence, and every source it considered stays in `plan.json`. Discarding needs high confidence; a doubtful source is kept as low relevance. *Planned*.
9. **Collected text is context only.** A fact found in a source can inform a suggestion or flag a discrepancy with the input, never set a value. *Planned*.
10. **The collection agent is paid work.** Its calls pass the money gates and are a line of the run's ledger. *Planned*.

## Bounds

11. **Every lap fits the envelope.** `cpuMillis` and `memoryMiB` stay within the envelope set on the server, and the native limits of all components sum to no more than the runtime's envelope. *Code* for the sum; *planned* for the server's envelope.
12. **Money gates run before paid work.** The server's gates decide whether a lap may start or continue its paid calls, and Microruntimes obeys them. Every plan with Microruntimes includes FinOps. *Planned*.
13. **Build recipes are a committed script or `docker build`, `go`, `npm`, `pnpm`, `make`.** No absolute paths, no `../`, no URLs, no secret or SSH mounts. *Code*.
14. **Builds see a disposable copy of the source** without `.langdrift`, dotenv files or symlinks, at most 8 MiB per file and 128 MiB in total. *Code*.
15. **Each operation is bounded to 120 seconds**, build included. Pre-warm a long build before the start. *Code*.

## Lifecycle

16. **A runtime changes revision only while stopped.** *Code*.
17. **A runtime lives for one run.** It is created from the plan, executed, stopped and deleted. *Planned*; today a runtime keeps its identity across runs.
18. **Load tests, tests and end-to-end runs end by themselves.** *Planned*; the fence generates Deployments only, and a run-to-completion kind is a contract change.
19. **Cleanup removes only what the execution ledger owns**, and an outcome it cannot confirm stays unknown. *Code*.
20. **The temporary context is deleted only after cleanup is confirmed.** If cleanup is unknown, the context and the execution ledger stay. *Planned*.

## Records

21. **Every run keeps `plan.json`, `run.json`, a memory file, one file per lap and its `evidence/`**, all inside the run's folder, including a blocked, failed or unknown lap. *Planned*; the CLI writes evidence under `.audits/evidence/runs/` today.
22. **Records keep what happened, not what was read.** `plan.json` keeps each source's path, type, hash, verdict and the facts extracted with their citations, never the source's content, raw logs or secrets. *Guidance*.
23. **Hashes are canonical.** A JSON record is hashed after canonical serialization, so a formatter that reindents a file does not change its hash. *Planned*; the hashes in this sample are byte hashes.
24. **Unknown is never zero.** A metric that was not observed, a cost without a price source, or an outcome that was not confirmed is written as `null`, unavailable or unknown. *Code* for observations; *guidance* for records.
25. **The plan sizes from input and from runs that measured resources**, and says when no run measured anything. *Guidance*.
26. **The interface reads fields, never prose.** Every value a screen shows comes from a JSON file or a JSON fence; the only exception is a ticket's title and status line, read by the tracker's convention. *Guidance*.

## Boundary

`.langdrift/` reads the project and writes only inside itself. The engine writes only `microruntimes/`. It reads the client's artifacts, wherever they live, and the mirrors in `harness/` and `finops/`, and never writes them.
