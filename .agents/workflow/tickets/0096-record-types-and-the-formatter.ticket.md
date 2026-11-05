# 0096: The records of a runtime and the formatter that writes them

**Spec:** .agents/workflow/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the CLI gains the types of every record the sample shows, in the sample's field order and names: the harness mirrors (envelope, workflow), the FinOps mirrors and ledger (budgets, gates, prices, forecasts with a planning entry, one ledger per runtime with a planning line, a line per lap and totals), the plan fence, the run record, the lap and run fences, and the private input record. A value not measured is `null`, arrays are never `null`, the illustrative flag round-trips but is never written by the CLI, and timestamps are RFC 3339 UTC to the second. A formatter writes files with tabs and short scalar arrays on one line, matching Biome's layout exactly, and fences with two spaces and small scalar objects on one line, so every sample file decodes into its type and re-encodes byte-equal.

**Blocked by:** 0093 (the fixtures).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Every JSON file and every JSON fence of the sample fixture decodes with unknown fields refused and re-encodes byte-equal through the formatter.
- [ ] The formatter's collapse rules are unit-tested, including an empty array, a long scalar array and a nested object, and `pnpm turbo run check` leaves a formatted file unchanged.
- [ ] The ledger's totals and the gate results are plain data; no arithmetic lives in the types.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; `pnpm turbo run check` over a file the formatter wrote into the fixtures.
