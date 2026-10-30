# 0141: Draw the black cat on the 404 and 500 pages and restyle the early-access confirmation note

**What to build:**
1. Render the early-access confirmation note in a readable accent color, smaller and shorter, keeping the trial clause in every locale.
2. Give the 404 page the website layout (header and footer), an illustration of the owner's black cat playing with a ball in the air and a back-to-home button.
3. Give the 500 page (the `(pages)` error boundary) the website layout, an illustration of the black cat in a beige postman uniform and cap beside a spilled mail bag, a back-to-home button and a retry button.
4. Draw the cat in an anime style: black fur, yellow eyes, no purple tint.

**Blocked by:** none

**Status:** completed

- [x] The confirmation note uses the new `primary-text` token role (`#A970FF` in dark, `#6D0AD6` in light), `text-sm`, with shorter copy in all five locales.
- [x] `not-found.tsx` and `(pages)/error.tsx` compose a shared `RouteState` view with header, footer, illustration and a back-to-home button; the error page retries with `retry()`.
- [x] The illustrations are SVG files under `apps/website/public/route-state/`, with localized alt text.
- [x] The confirmation page's hand-written stylesheet is replaced by Tailwind utilities.
- [x] Verification with Biome check, typecheck, test suite, token tests and drift checks passes.
