# 0063: The Book a demo page, its form and its route

**Spec:** .agents/specs/0009-home-product-sections-and-book-a-demo.spec.md

**What to build:** `/book-a-demo/` in every locale: a title, a lead and a form asking for a company email, a company name and a position from a list of twelve. An empty or invalid field shows its error in the visitor's language; a free-mail address is refused as such; a valid submission reaches the server, which validates it the same way and answers that requests are paused, and the form shows the return-soon notice.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The website depends on react-hook-form and the hookform zod resolver, pinned exactly, the newest releases the release-age policy accepts; no exclusion is added.
- [ ] One zod schema validates the three fields: the email through the early-access email rules plus a free-mail domain refusal, the company trimmed, required and at most 120 characters, the position one of twelve ids held in code. Every issue's message is an error code; a submission schema adds the locale and refuses unknown keys.
- [ ] `BookDemoForm` renders the fields with labels, placeholders, `aria-invalid` and `aria-describedby`, validates on touch and on submit, maps codes to the catalog's errors, disables submit while submitting, posts JSON with the locale, swaps the fieldset for the return-soon notice on 503 with the heading focused, and keeps the form with a live failure message on any other error.
- [ ] The route handler bounds and parses the body as the early-access route does, answers 400 with the field codes on failure and otherwise 503 with `Retry-After`, `Cache-Control: no-store`, `X-Content-Type-Options: nosniff` and `{ ok: false, code: 'UNAVAILABLE' }`; it stores, logs and forwards nothing.
- [ ] The page carries skip link, header, the title, the lead, the form and the compact footer; its metadata entry gives it a localized title and description, a self canonical and hreflang; the sitemap gains its five URLs.
- [ ] A new `book-a-demo` catalog, five locales in one file, holds the page, band, form, position and error copy, joins the catalog index and the parity test; two tests tie the error codes and the position ids to the English keys.
- [ ] Schema, route and form tests cover the cases the spec lists.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website`, `drift check localization`, `drift check dead-code`, `pnpm turbo run build --filter=website`.
- `/en/book-a-demo/` and `/pt-br/book-a-demo/`: an empty submit shows three errors; a Gmail address is refused; a company address gets the return-soon notice; `/en/sitemap.xml` lists 35 URLs.
