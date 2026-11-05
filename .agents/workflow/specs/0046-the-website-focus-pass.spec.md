# The website speaks for FinOps and its engine

**Status:** ready-for-agent

## Tickets

- [x] 0299: the site loads without console warnings
- [x] 0300: Early Access leaves and Book a Demo says we will be back soon
- [x] 0301: pricing reads Starter, Plus and Pro, coming soon
- [x] 0302: one FAQ, on the home
- [x] 0303: the copy speaks for today's business
- [x] 0304: platform headings carry the product icon
- [x] 0305: the LangDrift CLI has its own page
- [x] 0306: CI Gates has its own page
- [x] 0307: the Platform menu lists five products in two columns

## Problem Statement

The owner has set two priorities for the coming weeks: the website first, then the publishable packages. The console is set aside for now. A visitor to `https://langdrift.md` today sees a site that still describes an earlier business and leaves the products it ships half shown:

- **Console warnings.** The production console shows an unrecognized `browsing-topics` feature in the Permissions-Policy header, a Content Security Policy report for `eval`, stylesheets preloaded and never used, and form fields without an `id` or `name`.
- **Early Access.** The site still invites people to join early access, with a paused form, an API route, a confirmation page and a privacy section about Turnstile. Book a Demo is the way in now, but it is not ready to take requests either.
- **Pricing.** The pricing page shows one card, a "15-day free trial" with an "Early adopter offer" badge, under "Plans & early access". The Plus and Pro cards are hidden, and their copy describes product reviews from the Vision days.
- **FAQs.** Each platform page carries its own FAQ, and the home repeats only two of them.
- **Copy.** The footer says "From vision to reality, and everything in between." The FinOps demo labels its figures "illustrative". The CLI is called "the ex CLI".
- **Headings.** Microruntimes opens with its product icon before the title; FinOps and Safe Seat open with a plain "LangDrift FinOps" eyebrow.
- **Products without pages.** The LangDrift CLI has no page and no menu entry. CI Gates ships with `ex init` but sits in the menu as "Coming soon" with no page.
- **Menu.** The Platform menu is a single column in an order that does not lead with the engine.

## Solution

One pass over the website. The console app changes only in the Permissions-Policy header both apps share:

- The site loads with no console warnings: the Permissions-Policy names only features browsers know, no client bundle needs `eval`, no stylesheet is prefetched for a page that does not use it, and every form field has a name.
- Early Access is gone from the site. Book a Demo shows a "We'll be back soon" card until demos open; its form stays in the page as commented code for that day.
- Pricing reads "Plans" with a "Coming soon" badge and shows three cards, Starter, Plus and Pro, all "To be announced", each describing what LangDrift ships or plans: the CLI and budget-gated runs, CI Gates and team budgets, Safe Seat and multi-team envelopes.
- The home holds the only FAQ, with a tab per product: Microruntimes, FinOps and Safe Seat.
- The copy speaks for today's business. The footer carries the CLI's tagline, nothing is labelled illustrative, no Vision content remains, and the CLI is the LangDrift CLI in prose while its command stays `ex`.
- The FinOps and Safe Seat pages open with the FinOps icon and "FinOps" in place of the "LangDrift FinOps" eyebrow; the new pages open like Microruntimes, with their icon before the title.
- The LangDrift CLI and CI Gates each get a page under `/platform/`.
- The Platform menu lists Microruntimes, FinOps, LangDrift CLI, CI Gates and Safe Seat, in two columns on desktop.

## User Stories

1. As a visitor, I want the site to load without errors or warnings in the browser console, so that the site of a platform about evidence looks careful.
2. As a visitor on a browser that dropped the Topics API, I want the Permissions-Policy header to name only features my browser knows, so that the header does not raise an error.
3. As a visitor, I want no script on the site to need `eval`, so that the Content Security Policy never reports a violation.
4. As a visitor, I want no stylesheet preloaded for a page I did not open, so that the console stays clean and the network carries only what the page uses.
5. As a visitor using autofill, I want every form field to have a name, so that the browser can identify it.
6. As a visitor, I want no invitation to an early access that is closed, so that the site tells me the truth about how to get started.
7. As a visitor who wants a demo, I want the Book a Demo page to tell me demos will be back soon, so that I do not send a request nobody will receive.
8. As the owner, I want the book-a-demo form kept in the page as commented code, so that it returns as soon as demos open.
9. As a visitor reading the privacy page, I want it to describe only what the site collects today, so that I trust it.
10. As a visitor on the pricing page, I want to see that plans are coming soon, so that I know nothing is for sale yet.
11. As a visitor, I want three plans, Starter, Plus and Pro, each marked "To be announced", so that I can see how LangDrift will grow with a team.
12. As a visitor, I want each plan to describe what LangDrift does (the CLI, budget-gated runs, CI Gates, team budgets, Safe Seat), so that the plans match the products.
13. As a visitor, I want no free-trial promotion that is not on offer, so that I am not misled.
14. As a visitor, I want every question about the platform answered in one FAQ on the home, so that I do not hunt for it on each page.
15. As a visitor, I want the home FAQ to have a tab for Microruntimes, FinOps and Safe Seat, so that I find a product's questions quickly.
16. As a visitor on a platform page, I want no repeated FAQ at the bottom, so that the page ends on its product.
17. As a visitor, I want the footer to say what LangDrift does today, so that its line matches the product.
18. As a visitor, I want the FinOps demo shown without an "illustrative" label, so that the page reads as the product's own demonstration, as the owner decided.
19. As a visitor, I want no trace of the Vision product, so that the site describes one business.
20. As a visitor, I want the CLI called the LangDrift CLI, so that its name carries the brand, while the command I type stays `ex`.
21. As a visitor, I want every platform page to open with its product's icon, so that I recognize where I am.
22. As a visitor on FinOps or Safe Seat, I want the FinOps icon and "FinOps" above the title, so that I see which product the page belongs to.
23. As a developer evaluating LangDrift, I want a LangDrift CLI page that says how to install it, what it governs and how a run goes, so that I can try it.
24. As a team lead, I want a CI Gates page that shows the report a pull request receives and how a team installs it, so that I can judge whether to adopt it.
25. As a security reviewer, I want the CI Gates page to state the workflow's posture (read-only token, pinned actions, a pinned CLI version, no overwrite), so that I can approve it.
26. As a visitor, I want screenshots of a real CI Gates report from this repository, in light and dark, so that I see what it looks like.
27. As a visitor, I want the Platform menu to lead with Microruntimes and list FinOps, the LangDrift CLI, CI Gates and Safe Seat after it, so that the engine comes first.
28. As a visitor on desktop, I want the Platform menu in two columns, so that five products fit without a long list.
29. As a visitor on mobile, I want the drawer to list the same products in the same order, so that both menus agree.
30. As a visitor in any of the five locales, I want every change in my language, so that no page falls back to English.
31. As a search engine, I want the new pages in the sitemap with localized titles, descriptions and hreflang alternates, so that they are indexed like the other pages.

## Implementation Decisions

- **Headers.** The Permissions-Policy of both Next apps drops `browsing-topics=()` and keeps `camera=(), microphone=(), geolocation=(), payment=(), usb=()`. The header stays in each app's Next configuration, where it is today.
- **`eval`.** zod probes for `eval` when it builds an object schema. The only client schema on the site is the book-a-demo form's. With the form commented out of the page, no zod reaches the client. The Content Security Policy keeps `'unsafe-eval'` for development only.
- **Prefetched stylesheets.** The tab and accordion stylesheets of the vendored SmoothUI components are separate CSS files, which Next preloads when it prefetches the home. Ticket 0273 of spec 0038 redraws both components with Tailwind utilities and deletes the files, and this spec takes that ticket into its first ticket's scope.
- **Field names.** The Microruntimes workflow card's two selects take the names `model` and `effort`, and each human-stage checkbox takes `human-<stage>`.
- **Early Access out.** The early-access components, API route, confirmation page, catalogs, image, environment variables, Turnstile documentation and tests are deleted. The book-a-demo email schema stops importing the early-access schema and holds its own. The demo ticket in the Microruntimes session is renamed to a book-a-demo runtime.
- **Book a Demo paused.** The page renders a "We'll be back soon" card with a lead and a link to the home. The form's import and element stay in the page as commented code, with a comment that names when they return. The form's success, which led to the deleted confirmation page, shows a sent notice in the form's place. When the form returns, zod must run jitless (`z.config({ jitless: true })`) before its schema is built, or the Content Security Policy reports `eval` again. The book-a-demo API route stays.
- **Pricing.** The eyebrow is "Plans" and a `ProductBadge` reads "Coming soon". The cards are `starter`, `plus` and `pro`; each value reads "To be announced", and no card carries a promotion badge. The call to action goes to Book a Demo. The `pricing` description in the page metadata follows.
- **One FAQ.** The home's FAQ adds a Safe Seat tab beside Microruntimes and FinOps, reading each product catalog's `faq.items`. The platform pages no longer render a FAQ. Page-level FAQ title and description keys that nothing reads are removed.
- **Copy.** The footer description becomes "Agent-assisted delivery, accountable to your team's budget." "Illustrative" and its translations leave every catalog and comment. This follows the owner's latest instruction, which outranks localization rule 5's mention of illustrative data. "The ex CLI" and "ex connected agent" become the LangDrift CLI in prose, and terminal lines keep the `ex` command.
- **Product heading.** The icon tile Microruntimes draws before its title becomes one component in the platform folder, taking an icon. Microruntimes, the CLI page and the CI Gates page use it in their h1. FinOps and Safe Seat replace the "LangDrift FinOps" eyebrow with the tile holding the FinOps icon and the word "FinOps"; Safe Seat keeps its Coming soon badge on that line.
- **The CLI page.** It lives at `/platform/cli`, on the Microruntimes page pattern. It shows:
  - the h1 with a new CLI icon in the `packages/web` icon family;
  - the install line `npx @langdrift/cli init` and the requirements;
  - the flow from `init` to `agent status`;
  - what it governs: budget reservations with the gates FG-01 to FG-06, bounded runtimes, and the agent driver;
  - the `ex agent plan` screenshot from the docs;
  - links to npm and GitHub.
- **The CI Gates page.** It lives at `/platform/ci-gates`. It shows:
  - what `ex init --ci-gates` writes;
  - the workflow's posture;
  - the report a pull request receives, with FinOps readings that read "unknown" in CI and never zero;
  - the run-time FinOps gates.
  GitHub's pull request page does not render headlessly, so the images are renderings, in light and dark, served from `public/ci-gates/`: the real `ex gates report` output on the commits of pull request #122, and the repository's own CI Gates comment on that pull request, each rendered from its Markdown through GitHub's Markdown API. The comment's FinOps and Agentic figures are set by this repository's dashboard script, and the page says so. The page shows the one that matches the theme, through `dark:` utilities.
- **New pages, wiring.** Each new page has:
  - a catalog in five locales, registered in the loaders, the shape validator, the barrel and the parity test;
  - an entry in the page metadata with a unique localized title and description;
  - a route path in the page-path computation;
  - five sitemap entries with hreflang alternates.
  Both inherit the locale's social images.
- **Menu.** The products are Microruntimes, FinOps, LangDrift CLI, CI Gates and Safe Seat. The CLI and CI Gates link to their pages; CI Gates loses its Coming soon badge, since it ships today, and Safe Seat keeps its badge. The desktop panel is a two-column grid; the mobile drawer is one column in the same order. The header catalog gains the CLI's name and description.
- **Localization.** Every new or changed string exists in en, pt-BR, zh-Hant, ja and vi, in the same key order, with the same placeholders.

## Testing Decisions

- **The behaviour seam is the website's existing unit tests.** `pnpm turbo run check typecheck test build --filter=website --filter=web` passes after every ticket. A good test asserts what renders and where it links, never the file a component lives in.
- **Tests that change with the behaviour:**
  - the Next configuration test asserts the new Permissions-Policy;
  - the platform menu test asserts the order, the links and one Coming soon badge;
  - the catalog parity test covers the new catalogs and the removed ones;
  - the SEO contract test covers the new pages in the metadata and the sitemap;
  - the content data tests of the new pages assert their section ids;
  - the tab and accordion tests assert their utilities, as ticket 0273 sets out.
- **Localization.** `platform/cli/drift check localization` passes after every ticket that touches a catalog.
- **The browser is the visual seam.** The built-in browser checks, at 1280px and 375px, light and dark, in `en` and one other locale:
  - the pages, the menu, the FAQ tabs, the headings and the pricing cards;
  - that a full load of each page leaves the console free of the four warnings.
  No browser runner is added.
- **Prior art:**
  - `next.config.unit.ts`
  - `website-header/platform-menu.unit.tsx`
  - `template/i18n/catalogs.unit.ts`
  - `metadata/seo-contract.unit.ts`
  - `platform/safe-seat/safe-seat-content.data.unit.ts`

## Out of Scope

- The console's features, its stubs and the rules-compliance specs 0034 to 0042, which the owner has set aside. Spec 0038's ticket 0273 is the one exception, pulled in for the prefetched stylesheets.
- Making Book a Demo take requests again.
- Pricing figures and plan limits, which stay the owner's to decide.
- New social images per page.
- The CI Gates workflow and the `ex gates report` command themselves.

## Further Notes

- **Superseded tickets.** Tickets 0232, 0239, 0240 and 0293 reshape early-access code that this spec deletes. Each one is marked superseded by this spec, keeping its number.
- **CI Gates badge.** The menu drops CI Gates' Coming soon badge because `ex init` has generated the workflow since spec 0024. The owner may keep the badge at review.
- **A dead end, for now.** The pricing call to action leads to Book a Demo, which says demos are paused, until demos open.
- **Tabs and accordion.** Redrawn with Tailwind's defaults, they keep every computed style but two: weight 500 instead of 550, and a 4px trigger radius instead of 5px.
- **Plan copy.** The three cards describe what LangDrift ships or plans. Their exact items are a first draft for the owner to edit at review.

> **Note, 2026-11-05.** FinOps and Safe Seat no longer show the FinOps icon and name as a label above a sentence title: every product page now opens with the same heading, its own icon and product name at the Microruntimes size, so Safe Seat takes its own mark, and the sentence follows as a tagline (ticket 0311).
