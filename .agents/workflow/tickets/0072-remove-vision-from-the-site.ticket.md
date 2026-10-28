# 0072: Remove Vision from the site

**Spec:** .agents/workflow/specs/0011-finops-single-page-site.spec.md

**What to build:** Vision disappears from the website in every locale. The Vision page is gone and its address answers the localized "not found" page. The home page no longer shows the Vision section, and the hero demo, the evidence sections, the integrations, the ROI band and the Vision FAQ are gone with their data, images, state, styles, catalogs and the charting-graph dependency. The home's catalog goes after its footer description moves to the footer's own catalog. Every link that led to the Vision page or its FAQ now leads to the home's FAQ or to the FinOps page, the hero's main call to action opens the FinOps page, the platform FAQ keeps its FinOps and Microruntimes tabs, and the Book a demo copy names only Microruntimes and FinOps. The page metadata, the sitemap and the home's description no longer mention Vision. Shared stylesheets lose only the Vision rules.

**Blocked by:** 0071

**Status:** ready-for-agent

- [ ] The Vision address answers "not found" in every locale, and no sitemap entry or metadata entry names it.
- [ ] The home page renders the hero, the Microruntimes section, the FinOps section, a FAQ with two tabs and the Book a demo band, and no link on the site leads to a removed address.
- [ ] The footer and the pricing page link to the home's FAQ, and the footer still shows its description.
- [ ] The catalog parity test passes with the removed families dropped, and a production build of the website generates every page in the five locales.
- [ ] The footer, the focus ring, the plans and the other kept pages look unchanged after the stylesheets shrink.
- [ ] A search for the whole word Vision in the website finds nothing outside unrelated words and the shared color token names.

### Context

- [The spec](/.agents/workflow/specs/0011-finops-single-page-site.spec.md)
- [The website adapts to the repository](/.agents/workflow/adrs/0003-website-adapts-to-the-repository.adr.md)
- [The pivot](/.agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md)
