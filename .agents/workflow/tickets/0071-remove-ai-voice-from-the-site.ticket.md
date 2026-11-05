# 0071: Remove AI Voice from the site

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** AI Voice disappears from the website in every locale. The home page no longer shows the AI Voice section or its scripted voice preview, the Platform menu and the navigation drawer no longer list it, and the hero's rotating words and product column no longer name it. Every catalog key, type and test assertion that existed only for it goes with it, so the catalogs stay in parity in the five locales. The AI Voice icon leaves the shared web package, because nothing else uses it. The website README stops describing the voice preview.

**Blocked by:** None (can start immediately)

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The home page renders without an AI Voice section, and the voice preview component and its test are gone.
- [ ] The header menu, the drawer and the hero list no longer contain AI Voice, and the hero's rotating words no longer include the voice word.
- [ ] The catalog parity test passes in the five locales with the voice keys and the Vietnamese voice assertions removed.
- [ ] The dead-code check passes, so no AI Voice export is left without a user.
- [ ] Typecheck and the website's unit tests pass.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
- [The home sections this removes one of](/.agents/workflow/specs/0009-home-product-sections-and-book-a-demo.spec.md)
