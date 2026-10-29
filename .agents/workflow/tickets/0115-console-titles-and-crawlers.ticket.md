# 0115: Every console screen has a unique, localized title, and the console stays out of search

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:**

- The site name, the default title, the title template and the application name come from the catalogs. The viewport declares both color schemes and the theme color.
- Entity pages title themselves after their project, runtime or ticket, through the request-cached read, and fall back to the section label. The new-runtime page and the project's runtime list use the keys their headings use, and the zh-Hant configuration label stops repeating settings.
- A static robots.txt disallows everything.
- A favicon and an Apple touch icon from the LangDrift mark answer at the root.
- A localized not-found and a localized global error replace Next's English fallbacks. They use Tailwind only, with no inline style.

**Blocked by:** 0113.

**Status:** ready-for-agent

- [ ] In each of the four locales, every route's title is unique for what it shows, and entity pages name their entity.
- [ ] No title or site-name string remains in TypeScript.
- [ ] `/robots.txt` answers `User-agent: *` / `Disallow: /`, and `/favicon.ico` and `/apple-icon.png` answer 200.
- [ ] An unknown path renders the localized 404, with noindex and no CSP violation.
- [ ] The document declares `application-name`, `color-scheme` and `theme-color`.
- [ ] `drift check localization` passes.

### Verification

- `next build`; `curl` on the convention files and on an unknown path; a browser pass reading `document.title` on every route in two locales.
