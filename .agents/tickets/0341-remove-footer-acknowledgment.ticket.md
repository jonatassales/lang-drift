# 0341: remove the external inspiration credit from the README and website footer

**Spec:** None. The owner's instruction of 2026-11-07: remove the external inspiration credit from the root `README.md` Acknowledgements section and from the website footer across all locales.

**What to build:** `README.md` drops the external inspiration line under `## Acknowledgements`. `WebsiteFooter` (`apps/website/src/app/[locale]/website-footer.tsx`), `FooterMessages` (`apps/website/src/app/template/i18n/messages/messages.types.ts`) and `footer.json` (`apps/website/src/app/template/i18n/messages/footer.json`) drop the `acknowledgment` link, type field and localized strings in all five locales.

**Blocked by:** None.

**Status:** completed.

- [x] The external inspiration credit is removed from `README.md` and the website footer.
- [x] `WebsiteFooter` and `FooterMessages` no longer carry `acknowledgment` or `creditLink`.
- [x] `footer.json` keeps key parity across all five locales (`en`, `pt-BR`, `zh-Hant`, `ja`, `vi`).
- [x] `pnpm turbo run check`, `pnpm turbo run typecheck`, `pnpm turbo run test` and `drift check localization` pass.

## Dev part

**Override:** the owner directed the agent to remove this credit and open a pull request; removing a footer link, a catalog key and a README line holds no core-domain logic, hard part or new stack (policy, Categories). 2026-11-07.
