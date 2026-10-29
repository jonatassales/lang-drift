# 0098: The session's lines, printed as the website prints them

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the CLI gains one catalog of the session's lines, mirroring the website's terminal keys for the plan and the run, and two renderers: a human one that prints the website's words with the `agent ›` actor, the ✓ and # markers and true figures, and a JSON one that prints one event per line with its time, kind, the website's line id, the reveal id of the card field it fills, and the data behind the figures. Money prints as dollars with cents, tokens as compact counts (`3.5M`, `250K`), resources as millicores and MiB. A test proves the catalog covers every key of the website's plan and run sessions, so a line cannot be forgotten.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] The catalog has an entry for every plan line id and every run message id the website defines, and the test fails when one is missing.
- [ ] `--json` on a session prints NDJSON events whose ids and reveal ids match the website's, in order.
- [ ] A rendered line matches the website's English text for the sample's figures (for example the estimate, gate and ready lines).

### Verification

- `GOTOOLCHAIN=local go test -race ./...`.
