# 0322: `@langdrift/core` 0.3.0 publishes the new shapes

**Spec:** .agents/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** A skill author, a console developer or anyone running the engine installs `@langdrift/core` 0.3.0 and finds the fence's schema, the judge's three answers and the five stage reports, with the README describing the packages as they now stand. The Go module is tagged `lib/core/contracts/go/v0.3.0`, and the CLI requires it.

**Blocked by:** 0318, 0319, 0320, 0321.

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] The manifest's `files` and `exports` add the fence's schema, the `judgment` schemas and samples and the stage-report schemas, and `pnpm pack --dry-run` lists exactly them beside the current files.
- [ ] The README of `@langdrift/core` and of the Go module describes the fence v3, the judge's three answers, the stage reports and the hand-off, and says which older shapes remain until their contract tickets land (0323, 0324, 0325).
- [ ] The tag `lib/core/contracts/go/v0.3.0` sits on the merge commit, `lib/cli/go.mod` requires it, and the CLI builds outside the workspace; the owner publishes 0.3.0 on npm.

## Dev part

**Override:** No part qualifies: release configuration and documentation (policy, Categories: never configuration). 2026-11-05.
