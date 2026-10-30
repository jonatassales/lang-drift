# 0140: Align footer across pages, comment last two pricing cards, remove paused card, and link GitHub in header

**What to build:**
1. Align the website footer across all pages so every page renders the same full footer as the home page (removing `compact` footer usage on `/pricing`, `/privacy`, `/book-a-demo`, and `/early-access-sent`).
2. Comment out the last two pricing cards (Plus and Pro) on `/pricing` so only the 15-day free trial card is displayed.
3. Remove the "We’ll be back soon." early access card from `/pricing`.
4. Link GitHub to `https://github.com/jonatassales/lang-drift` in the header navigation bar and drawer.

**Blocked by:** none

**Status:** completed

- [x] All website pages render the full home footer (`<WebsiteFooter locale={locale} />`).
- [x] Last 2 pricing cards (`plus`, `pro`) are commented out on `/pricing`.
- [x] "We’ll be back soon." card (`EarlyAccessSection`) is removed from `/pricing`.
- [x] GitHub links to `https://github.com/jonatassales/lang-drift` in the header desktop nav and mobile drawer.
- [x] Verification with Biome check, typecheck, test suite, and drift checks passes.
