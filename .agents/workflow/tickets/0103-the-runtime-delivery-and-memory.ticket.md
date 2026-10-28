# 0103: The run brings the runtime up, observes delivery, stops and remembers

**Spec:** .agents/workflow/specs/0015-ex-performs-the-microruntimes-session.spec.md

**What to build:** after the implement stage, the run materializes the runtime from its fence through the existing native lifecycle with a local operator record, builds the images, claims an execution, starts, observes readiness and prints the runtime-ready line; it observes the branch, the short commit and the upstream from the repository and the pull request and its checks from the forge's CLI within a bounded wait, printing pushed, pull request and checks lines or their unknown forms; it stops and cleans up by the execution ledger, printing what was removed or that the outcome is unknown; it releases the unused reservation, settles the ledger's lap line and totals, prints ready for review or ready to merge with the amount released and the next-step line; it writes the lap's resources table and fence from the observations, the run record's outcome and window, the memory file with its fence, and rebuilds the RAG index. Evidence lands in the runtime's folder keyed by execution.

**Blocked by:** 0102, 0099.

**Status:** ready-for-agent

- [ ] With a fake lifecycle reporting 410 millicores and 690 MiB peaks and a fake delivery observer returning the website's branch, hash, pull request 19 and four checks, lines 15 and 17 to 23 print with the sample's figures, and the whole `.langdrift/` tree of the reproduced session equals the expected fixture except the private state.
- [ ] The lap's fence carries the runtime id, revision and markdown hash, and the memory fence's observed peaks and tokens match the ledger.
- [ ] An unconfirmed cleanup leaves the lap outcome unknown and the reservation held; a missing forge CLI leaves the pull request and checks unknown without failing the run.
- [ ] A connected start without a plan still works under the moved paths, and a second plan after this run cites it.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; the end-to-end run of the spec's verification on a scratch copy of this repository with kind available.
