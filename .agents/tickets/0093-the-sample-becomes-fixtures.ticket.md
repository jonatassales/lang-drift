# 0093: The `.langdrift` sample becomes the CLI's fixtures and the folder is ignored

**Spec:** .agents/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the repository stops versioning `.langdrift/`: an anchored ignore rule at the root keeps the folder out of Git while the CLI package's fixtures stay tracked, and the committed sample moves, byte for byte and from the tree's history, into the CLI package's test fixtures: the sample itself, plus a frozen copy of ticket 0070 and every source its links reach, so later tests can reproduce the sample's run with stable hashes. The working tree's leftover `ex init` output under `.langdrift/` goes, since the CLI regenerates the folder. The CLI package's own ignore file drops the rules the move makes pointless, and the sample's README table moves into the CLI README.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `git ls-files .langdrift` prints nothing, and a file created under `.langdrift/` is ignored while the fixtures under the CLI package are not.
- [ ] The fixture copy of the sample is byte-identical to the sample at the commit before this ticket, checked with `git show`.
- [ ] The frozen sources are the ones the sample's plan records, and the plan's recorded hashes match the frozen copies, except the one the sample itself notes as having changed, which is regenerated once.
- [ ] `pnpm turbo run check` passes on the fixtures and the CLI's suite is unchanged.

### Verification

- `pnpm turbo run check`; `GOTOOLCHAIN=local go test -race ./...` in the CLI package; `git status --short .langdrift` empty after the move.
