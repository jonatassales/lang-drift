# 0097: `ex agent plan` reads the ticket, follows its links and writes the runtime file

**Spec:** .agents/workflow/specs/0015-ex-performs-the-microruntimes-session.spec.md

**What to build:** `ex agent plan <ticket>` prints the first act of the session for real: it parses the ticket's number, slug, title and status, reads it, lists the artifacts it links, follows markdown links, wiki links and field-line paths breadth-first under the limits (depth, file count, total bytes), refuses secrets, symlinks and paths outside the project and records them, records every source with its type, size, hash, depth and the link that reached it, hashes the context, finds a guardrail document's criteria, prints the envelope and budget as the CLI input that set them with their records, flags a document that states different figures as a discrepancy, and writes `microruntimes/runtimes/<id>/runtime.md` with the `langdrift` fence (from the inputs and the components given on the command line) and the `langdrift-plan` fence, state `proposed`, plus the ledger's planning line with nothing paid. The runtime id derives from the ticket's slug. Judged fields stay `null` in this ticket; the lines that need a judge print their unjudged form.

**Blocked by:** 0095, 0096.

**Status:** ready-for-agent

- [ ] Against the frozen fixture project, the plan prints lines 1 to 13 of the session with the sample's counts and hashes, and the written plan fence's `ticket`, `collection`, `sources`, `inputs` and `discrepancies` equal the expected fixture.
- [ ] A linked `.env`, a symlink and a path above the project are recorded as refused and never read; a link to nothing is recorded as unresolved; a binary is recorded by hash.
- [ ] The context hash equals the sha256 over the sorted `path<TAB>sha256` lines with a trailing newline, as the sample records it.
- [ ] Replanning raises the revision and is refused while the runtime runs; the fence passes the contracts' parser.
- [ ] Without a ticket at the path, exit code 4; without an envelope, budget or model input, exit code 4 naming the missing input.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; `./dist/ex agent plan .agents/workflow/tickets/0070-console-microruntime.ticket.md --no-judge` in a scratch copy of this repository.
