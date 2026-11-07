# 0074: Remove the Vision and voice UI from the shared web package

**Spec:** .agents/specs/0011-finops-single-page-site.spec.md

**What to build:** The shared web package exports only what an app uses. The Product Vision curve, the target distribution chart, the classification text helper, the avatar set, the Vision and AI Voice icons and two vendored SmoothUI components (the avatar group and the chroma text) are removed, with their exports, style imports, the chroma text rules and the charting dependency, including its dead-code allow entry and its lockfile lines. The agent orb, its voice dependency, tabs, accordion, brand, info cards, detail list, the language switcher and every token stay. The drift sign helper and its attribute rules go only if the dead-code check reports them.

**Blocked by:** 0072

**Status:** wontfix

**Resolution:** absorbed by 0071 and 0072. The dead-code check fails as soon as the website stops importing a module of the shared web package, so each removal left the package together with the website change that orphaned it.

- [ ] No app imports any removed export, and the package's export map lists only modules with a consumer.
- [ ] The dead-code check passes with no stale allow entry, and the lockfile loses only the removed dependencies.
- [ ] The footer still renders the agent orb, and the FinOps and Microruntimes pages render as before.
- [ ] Typecheck, check and tests pass for the package and for the apps that use it.

### Context

- [The spec](/.agents/specs/0011-finops-single-page-site.spec.md)
- [The classification and drift colors are shared meaning](/.agents/adr/0003-website-adapts-to-the-repository.adr.md)
