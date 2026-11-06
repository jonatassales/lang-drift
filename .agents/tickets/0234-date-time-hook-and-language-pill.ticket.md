# 0234: the date-time hook is a locale helper and the language pill navigates with Link

**Spec:** .agents/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** `useDateTime` moves to `template/i18n/` as `use-date-time.hook.ts`. The language pill renders the i18n `Link` with the target locale, keeping the query string and hash suffix and the drawer-focus hand-off, so switching the website's language becomes a client-side navigation; the rendered `href` values do not change.

**Blocked by:** 0229, 0230.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `useDateTime` lives at `template/i18n/use-date-time.hook.ts` in its app and every import follows.
- [x] The language links render the i18n `Link`; their `href` values, query strings and hashes are unchanged.
- [x] Switching language in the built-in browser keeps the drawer-focus hand-off and writes the preference cookie through the proxy.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website`
- `pnpm turbo run check typecheck test build --filter=console`
- `platform/cli/drift check localization`

> **Note, 2026-11-05.** The rendered hrefs, query strings and hashes are unchanged, the switch is a client-side navigation, and the proxy still writes the locale cookie. The drawer-focus hand-off did not return focus to the trigger before this change either (the full-page load reset its module state); with the client navigation it can now, and ticket 0290's drawer hook takes that up. The console's importers move again with 0230.
