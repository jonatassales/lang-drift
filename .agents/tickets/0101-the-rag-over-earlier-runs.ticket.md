# 0101: The plan retrieves the project's episodes

**Spec:** .agents/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the plan's second act prints for real: the ticket's intent is embedded with a deterministic local embedder, the project's episodes under `agent/state/rag/` are searched by cosine similarity, the top six candidates print with their scores, each is kept or dropped against the threshold with its structured criteria (same component path, same stack, measured resources, measured tokens), the cited episodes and whether any measured resources print, and the sizing basis follows the cited measured runs or the input. `ex agent reindex` rebuilds the index from every runtime's memory and plan fences, recording the embedder's id; a plan rebuilds it when the id changed. The judge's reason per candidate stays `null` in this ticket.

**Blocked by:** 0099.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] With an index of the website's six episodes, lines 14 to 21 print with the sample's scores, kept and dropped counts and the threshold, and the plan fence's `retrieval` and `sizing` equal the expected fixture.
- [ ] The embedder is deterministic across runs and platforms, and cosine similarity of a text with itself is 1.
- [ ] A second plan in a project with one finished run cites that run with its matched criteria and a sizing basis of measured runs when it measured resources.
- [ ] With no episodes the lines say so and the retrieval record notes the first run.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; `./dist/ex agent reindex` then a plan in a scratch project with a finished run.
