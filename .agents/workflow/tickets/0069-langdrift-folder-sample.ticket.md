# 0069: An illustrative `.langdrift` folder covers the console

**Spec:** .agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md

**What to build:** the repository carries a versioned `.langdrift/` folder with one illustrative run of ticket 0070 of this repository's own tracker, executed as the microruntime `console-local` in two laps: the first paused at a FinOps gate, the second completed. It holds data only, with the layout of spec 0010: no Vision folder, no artifact templates, and every JSON file in `harness/` and `finops/` a declared mirror of a server record. Its `plan.json` is the graph of sources the ticket reached, each kept or discarded with a reason, and limits come from recorded console input, never from a document. Every field of the Microruntimes page's Ticket, Analysis, Plan and Workflow cards has a source in it. The research note on what the real `ex` reproduces of the page joins the tracker.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] The tree matches spec 0010; only `agent/state/` is ignored.
- [x] Nothing in it points outside the project's own folder except the workload it runs and the client's artifacts it reads.
- [x] No guardrail value has a document as its source; each one names the console record that set it.
- [x] Every hash in `plan.json` and the laps matches its file, and `runtime.md` passes the CLI's parser.
- [x] Each card field maps to a JSON field or a JSON fence, as the spec's coverage table lists.

### Verification

- `pnpm turbo run check`; the sample's hashes recomputed with `shasum -a 256`; `runtime.md` parsed with the CLI's `connection.Parse`.
