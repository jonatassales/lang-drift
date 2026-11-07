# 0045: The Microruntimes page in the Vidro editor, without illustrative disclosures

**Spec:** None

**What to build:** the owner's review of the Microruntimes page asked for three changes. The hero sets the `ex` terminal beside the lifecycle card, the terminal only as wide as its content. The runtime.md card opens its langdrift fence in the owner's Vidro editor with JSON highlighting. Every illustrative, simulated or demo disclosure on the page goes, in every locale; the "demo requests are paused" status stays.

**Blocked by:** None (can start immediately).

**Status:** done

- [x] From 1024px up the terminal and the lifecycle card share one row. The terminal is as wide as its longest line from the first frame, never more than half the row, and does not widen while it types; the card's grids follow the card's own width. Below 1024px they stack.
- [x] The runtime.md fence renders read only in `<vidro-editor>` with JSON highlighting, the file path in its header and a copy button, in the site's color scheme. The server renders the plain text, which stays when the editor cannot load.
- [x] `vidropro@0.1.0-preview.3` (MIT) is vendored in `packages/web/vendor/` with its provenance until it is published, and wrapped by `@langdrift/web/vendors/vidro/code-editor`. A `zh-Hant` page gets the editor's English controls; Simplified Chinese never stands in.
- [x] The hero caption, the descriptor badge and the ledger period are gone, and the terminal label, scenario selector, hero summary, state lead and availability answer no longer call anything illustrative, in all five locales. FinOps keeps its own caption and period.

The page no longer follows `08-localization.rule.md` point 4 and spec 0005 on illustrative disclosures; the owner decided so for this page.

### Verification

- `/platform/microruntimes/` renders in `en` and `zh-hant` with the editor ready and JSON highlighted.
- `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test` and `drift check localization`.
