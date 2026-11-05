# 0312: commands copy with one click, and the Platform menu divides its columns

**Spec:** None. The owner's requests of 2026-11-05: every command shown for copy and paste carries a copy-to-clipboard icon, starting with `npx @langdrift/cli init` on the LangDrift CLI page; the CLI page no longer lists `ANTHROPIC_API_KEY` among its requirements (the docs keep it for now); and the Platform menu gets a delicate separator between its two columns.

**What to build:** `CommandLine`, which sets every copyable command on the CLI and CI Gates pages, carries a button that copies its command, shows a check when the copy took and the failure label when the browser refuses it, and announces both. Its labels live in the cli and ci-gates catalogs in every locale. The CLI hero's requirements line drops `ANTHROPIC_API_KEY` in every locale. On wide screens, the Platform menu draws a one-pixel line between its columns that fades in and out along its height, in the border role.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] Every `CommandLine` shows a copy button; a unit test covers the copied and the failed states.
- [x] The copy labels exist in the five locales of `cli.json` and `ci-gates.json`, with key parity.
- [x] The CLI requirements line no longer names `ANTHROPIC_API_KEY` in any locale.
- [x] The Platform menu's two columns are divided by a fading hairline from `sm` up; the mobile list has none.

## Dev part

**Override:** the owner does not write this part ("I won't write anything", 2026-11-05); the agent writes it. A copy button, a catalog line and a separator hold no core-domain logic, hard part or new stack (policy, Categories). 2026-11-05.
