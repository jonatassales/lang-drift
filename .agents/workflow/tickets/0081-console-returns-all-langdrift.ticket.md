# 0081: The console returns, all LangDrift

**Spec:** .agents/workflow/specs/0013-console-areas.spec.md

**What to build:** The console comes back from the last commit with every Exphere mark replaced: the LangDrift brand in the navigation and the sign-in form, the app icon's title, the locale cookie, the catalog strings and the `.langdrift` folder name. Routes that only redirect go, and the sidebar links straight to the secrets page.

**Blocked by:** 0080

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] A search for Exphere in the console finds nothing.
- [ ] Every remaining route renders, the console's tests and type check pass, and the dead-code check passes.

### Context

- [The spec](/.agents/workflow/specs/0013-console-areas.spec.md)
