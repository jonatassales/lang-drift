# LangDrift — design review (taste design expert)

**Date:** 2026-11-07. **Question:** how do all 10 English pages of `langdrift.md` measure against the `taste-design-expert` (`design-taste-frontend`) rules when inspected in a real browser? **Method:** review of all 10 English pages of `langdrift.md`, measured in a real browser at 1440×900 and 390×844 · Nov 7, 2026, with page-level script measurements at 1440×900, 1024×768, 768×1024 and 390×844 and screenshots read at 1440 and 390. Full method and limits are in section 8.

## 1. Verdict

The right subject is on the screen. Every page shows the product's own interface — terminal, ticket panel, budget ledger — in live HTML, with real numbers. There are no stock photographs, no invented social proof, no fictional customer logos, and no dead links: all 40 anchors on each page go somewhere. That is the hard part of a software site, and it is done.

The failure is one of discipline, on three fronts. The site's primary action — the **Book a Demo** button, present on all 10 pages — leads to a page that says requests are paused and does not have a single form field. No 900px screen shows a single section by itself, on any of the 10 pages, at any of the four measured widths. And the brand is built in violet `#6D0AD6`, the one color the rule forbids in the page's own design.

None of this calls for redesigning the product. These are four system fixes: a destination for the action, a section box with screen height, a color token swap, and a committed font.

## 2. The kind decides everything

The kind is **software** (site-kind software, page type software), developer-tool family. That unlocks what the site already uses: dark page, dark hero, one effect, up to two moments, zero to two pulsing diagram cards, one timed loop, CSS grounds. It also fixes the radius at `sharp` (4/6/8) and requires a section order: nav, hero, **app-frame** (the product large, the list beside the open item, in a single window, as section 2 or 3), two or three feature sections in different layouts, integrations as a diagram card, one quiet section after every two visual ones, pricing, FAQ, `footer.closing`.

None of the 10 pages declares its kind. There is no `data-type` on `<html>`, no `<meta name="page-type">`, no `<meta name="site-kind">`, no `data-radius`, and no `data-motion`. The site has the right ingredients — live UI is exactly what the software recipe calls for as the visual — but lacks the declared skeleton that makes the checkers and the layout rules themselves work.

**Already passes, and should stay that way:** no stock photographs, no invented testimonials or figures, no customer logo wall, no carousel, no pop-up, no floating chat, no award badge. Sideways scroll: none at 1440, 1024, 768 and 390. The Safe Seat page cites research and sends the reader to the sources — proof in the exact format the rule wants.

## 3. Critical failures

### 1. The site's primary action leads nowhere

**Rule:** the primary action names the transaction and reads the same in the nav, in the hero and in the closing; forms submit; "coming soon" is only for what is actually coming.

**Where:** `/en/book-a-demo/`, all widths. h1 "Book a demo", h2 "We'll be back soon.", content height of 0.40 screen, **zero form fields** (measured: no `input`, `select` or `textarea`, no `form`). All 10 pages point their single filled button there, and `/en/pricing/` repeats "Book a demo" three times, once per plan.

**Fix:** one of two options. Either a form with four to six fields (name, work email, team, question), 48px fields, `method="post"`, with a line stating how long a reply takes. Or swap the label across the entire site to the transaction that exists today — "Install the CLI", which is real and published on npm — and remove "Book a demo" from the nav and the three plans until the channel is back. A dead-end button is worse than one fewer button.

### 2. No screen shows a single section by itself

**Rule:** one section per screen. The screen at `scrollY = k × H` shows exactly one section; every `main > section` has `height: calc(var(--h,1) * 100svh)` with `--h` set to 1 or 2.

**Where:** across the entire site. On the home page at 1440×900, screen 1 = hero + llmops, screen 2 = llmops + faq, screen 3 = faq + footer. No section uses screen height: heights are whatever the content asks for.

| Width | Home section heights, in screens |
| --- | --- |
| 1440×900 | 1.58 · 1.18 · 0.74 · 0.27 |
| 1024×768 | 1.50 · 1.37 · 0.86 · 0.35 |
| 768×1024 | 1.63 · 1.58 · 0.63 · 0.27 |
| 390×844 | 1.86 · **2.61** · 0.78 · 0.60 |

At 390 the `#llmops` section exceeds two screens, which is the absolute ceiling. The same pattern holds across the other nine pages: `ci-gates#report` is 1.83 screens next to two half-screen sections; `cli#install` is 1.51; `microruntimes#evidence` is 1.37; safe-seat has eight sections and none of them closes in one screen.

**Fix:** a single section box — `height: calc(var(--h,1) * 100svh)`, 120 padding on top (where the h2 sits) and 72 on the bottom, one `.wrap` as the only child, h2 first and top-aligned, never vertically centered. Anything exceeding a 708px content window splits into two sections with its own h2; anything under 360px moves up into the section above. Start with `#report`, `#install` and the short openings of safe-seat, cli and ci-gates.

### 3. The brand is violet

**Rule:** never violet or lavender in the page's own design — accent, ground, gradient, effect color or tint. A token whose hue falls between 270° and 330° becomes blue (255°) or pink (345°).

**Where:** "Play workflow" button in `#6D0AD6` (hue \~273°); 20 elements in the violet range on the home page, 11 on microruntimes, 3 on finops and 3 on pricing; and the word "Platform" in the h1 painted with a blue→violet→pink gradient.

The h1 breaks three rules at once: violet, emphasis by color instead of contrast, and a two-tone headline — which is on the list of SaaS furniture to avoid.

**Fix:** swap the accent token to blue (\~255°), the closest neighbor to the current violet, and paint "Platform" in the same ink as the rest of the h1. The accent stays on the filled button, the focus ring, links, and active state inside visuals — under 3% of any other screen.

### 4. The hero neither fills the first screen nor has a button

**Rule:** the hero is exactly the first screen, `100svh`; the h1 comes first, with its top at y=120; buttons sit above the fold at every width.

**Where:** home at 1440×900 — 1425px hero (1.58 screens), h1 starting at y=251, and the first real button ("Play workflow") at y=1709, nearly two screens below the fold. At 390 the first screen shows the h1 and a terminal cut off by the fold. On the inner pages there is no hero: safe-seat opens at 394px, cli at 365px, ci-gates at 287px — a title band, not a screen.

**Fix:** a one-screen hero, h1 at the top at y=120, a 15-to-30-word subline stating what it does and for whom, at most two buttons inside the screen, and the UI as a layer behind or beside it, ending inside the hero. On the inner pages, the same hero with the module name and a single line stating what it solves.

### 5. The h1 names nothing — and swaps words on its own

**Rule:** the h1 names the category and, for an unknown brand, the buyer, in two to eight words. Never a riddle, a mood, or a capitalized coined concept.

**Where:** home — "Agentic DevEx Platform", with the middle word swapping to "LLMOps" in an animation. The DOM reads `AgenticDevEx PlatformDevEx Platform`: the text is duplicated for screen readers and search engines. On the inner pages the h1 is only the module label — "LLMOps", "FinOps", "Microruntimes", "CI Gates" — at 43px, smaller than several h2s on the same page.

**Fix:** the best sentence on the site is already written, in the footer: *"Agent-assisted delivery, accountable to your team's budget"*. It names the category and the buyer. Promote that to the h1, without animation and without duplicated text, and on the inner pages use "LLMOps for teams that pay per model call" and equivalents.

## 4. Serious failures

### 6. System typography

**Rule:** the font is dealt, never chosen — one family, or one display family paired with one text family; at most four font files; never Inter, Roboto, Open Sans, Lato, Montserrat, Poppins, Space Grotesk, Plus Jakarta Sans **nor the system look**.

**Where:** `ui-sans-serif` across all 10 pages, **zero font files loaded**. Only microruntimes brings in `Geist` at one spot — meaning two families on the site and neither committed. It is the detail that makes a page look like a template: nothing in the typography says this site belongs to this product.

**Fix:** a single family, loaded with `display=swap`. For a developer tool, Geist or IBM Plex Sans (which already pair well with the monospace in the terminals), in two or three weights, at most four files.

### 7. The type scale is inverted

**Rule:** h1 72 / h2 48 / h3 28 / lead 20 / body 16 at 1440, each level at least 1.3× the next, one size per element type across the entire site.

**Where:**

| Page | h1 | h2 | h3 |
| --- | --- | --- | --- |
| home | 72px | 43 and 46px on the same page | 28px |
| llmops | 43px | 28px | 18px |
| microruntimes, safe-seat, cli, ci-gates | 43px | **62px** | 16–18px |
| pricing | 50px | 21px | **11px** |

On four pages the h2 is larger than the h1. On pricing there are four levels within a 39px span.

**Fix:** a single scale applied identically across all 10 pages, with the h1 always the largest text on the page.

### 8. Text below 13px — and it is the very proof of the product

**Rule:** nothing below 13px; reading text at 16px or above.

**Where:** **219 nodes at 12px** on the home page — all the content inside the terminals and the budget ledger tables. At 390 the smallest text drops to **9px**. On pricing, the h3s ("Included", "Everything in Starter, plus") sit at 11px. The content that proves the product most is the smallest text on the page, and on mobile it is illegible.

**Fix:** 13px as the floor for monospace, 16px for running text in tables, and if it does not fit, show fewer rows instead of shrinking the font — text never shrinks to fit.

### 9. Accordion FAQ

**Rule:** the FAQ shows all answers open, never an accordion; nothing animates height, `max-height` or padding.

**Where:** home, `#faq` section — all six questions have `aria-expanded="false"`. Anyone reaching the section sees six question lines and no answers; the section takes up 0.74 screen because it is collapsed.

**Fix:** answers open, in rows with a 1px rule closing each one. That also solves half of the section's height problem.

### 10. The nav is transparent, is 81px tall, and uses blur

**Rule:** fixed `header.nav`, 72px, opaque background with a 1px hairline from first paint, identical at any scroll depth. No blur, no transparency, no shrinking, no floating pill.

**Where:** across all 10 pages — `background: rgba(0,0,0,0)`, no border, 81px tall, and four to six elements with `backdrop-filter` (the "EN" pill, the panels). `backdrop-filter` is on the banned list.

**Fix:** opaque background, 1px hairline, 72px height, the same `.wrap` as the page, logo on the left edge and the button ending on the right edge. Strip blur from all panels.

### 11. Three content edge lines on the same page

**Rule:** equal left and right edges; a single `.wrap` for header, sections and footer; at most two text lines — the left edge and the split column.

**Where, at 1440:** on the home page sections start at x=115 and x=116, and the footer at x=128; the right edges are 1324 and 1312. On llmops, x=115, x=166 and a block closing at 1126 coexist. There are three wrap widths on the site: 1180, 1240 and 1280px.

**Fix:** a single `.wrap { width: min(100% - 2*var(--gut), 1280px) }`, used by the header, every section and the footer. It is a 13px shift, invisible in isolation and visible when the page scrolls.

## 5. Medium failures

**12. Too many things moving on their own.** The rule gives a tech page one timed loop, one effect, one or two moments, and nothing else moving on its own. Measured at rest: **15 running animations** on the home page (plus one canvas) and **13 on microruntimes** (plus three canvases). Never two things moving on their own inside the same 900px window. The terminal's "Pause" button is right and must stay — it is exactly what the rule asks for when something runs past 5 seconds. Reduce to one loop per card group, 2.5 to 4 seconds, and leave the rest on their final frame.

**13. Banned button labels.** A button is a verb plus what the person gets; "Explore", "Learn more", "Discover" and "More" are banned. Measured: **"Explore LLMOps"** on the home page and on microruntimes. Replace with "See how spend is reserved".

**14. Headings that are sentences, not nouns.** Every heading needs a noun naming a thing, in one to six words; never a sentence ending in a period. Measured: "One team envelope. Every paid call reserved before it runs.", "The envelope holds when the plan does not.", "Nothing starts because a page was opened.", "Evidence over assertion." (repeated on llmops and finops), "A few things worth knowing." The sentences are good — they are just in the wrong place. The h2 becomes "Team envelope", "Failures and reversal", "Evidence", "Questions", and the sentence moves down to the line below the heading.

**15. Empty placeholders on the first screen.** No placeholders: either real text, or the element comes out. Measured on microruntimes, in the ticket panel visible on load — `LINKED ARTIFACTS` lists PRD, DESIGN DOC, SPEC and GUARDRAIL, each with an em dash in place of content. Four em dashes on the first screen of the page that explains the product. Fill them with filenames from the example itself, or remove the rows.

**16. The pricing page has neither prices nor other proof.** Three plans, all three reading "To be announced", and the entire page occupies 1.21 screens. "Coming soon" is allowed when it is true, and here it is. But the page is left with nothing that sells: while there is no number, state what can already be done today — the CLI is published on npm and is free — and move the three plans after that. Also: the primary action appears three times on the same screen, once per plan, whereas the rule asks for one action per screen.

**17. Printed email address in the footer.** The rule does not print an email address: it uses an "Email us" button with `mailto:`. Measured: `founders@langdrift.md` printed across all 10 pages. It is a real address and yours, so the cost here is only automated scraping — low priority, but a one-minute fix.

## 6. Page by page

Measured at 1440×900. "Sections" lists the height of each section in screens, in page order; the target is 1 or 2 for each.

| Page | Total | Sections (in screens) | h1 | Main finding |
| --- | --- | --- | --- | --- |
| `/en/` | 3.66 | 1.58 · 1.18 · 0.74 · 0.27 | "Agentic DevEx Platform", 72px | Hero without button and 1.58 screens; animated h1 duplicated in the DOM; violet; collapsed FAQ |
| `/en/platform/llmops/` | 5.08 | 0.97 · 0.95 · 1.17 · 0.84 · 0.79 | "LLMOps", 43px | No hero; three edge lines (115, 166, 1126); no section closes in one screen |
| `/en/platform/finops/` | 3.82 | 1.05 · 0.84 · 0.79 · 0.76 | "FinOps", 43px | Same structure; closes with no action at all before the footer |
| `/en/platform/microruntimes/` | 7.36 | 0.87 · 0.88 · 1.21 · 1.00 · 0.72 · 1.37 · 0.94 | "Microruntimes", 43px | 62px h2 larger than the h1; four placeholder em dashes at the opening; 13 animations |
| `/en/platform/safe-seat/` | 6.72 | 0.44 · 0.72 · 1.22 · 0.63 · 0.80 · 0.66 · 0.87 · 1.02 | "Safe Seat", 43px | Eight sections, none at one screen; 394px opening; the best content on the site, with the worst structure |
| `/en/platform/cli/` | 3.66 | 0.41 · 1.51 · 0.62 · 0.76 | "LangDrift CLI", 43px | `#install` at 1.51 screens next to a 0.41 opening |
| `/en/platform/ci-gates/` | 4.34 | 0.32 · 0.46 · 1.83 · 0.60 · 0.78 | "CI Gates", 43px | `#report` at 1.83 screens between four half-screen sections — the most irregular rhythm on the site |
| `/en/pricing/` | 1.21 | 0.85 + footer | "Find your starting point.", 50px | h1 that names nothing; h3 at 11px; primary action repeated three times; no price |
| `/en/book-a-demo/` | 1.00 | 0.40 + footer | "Book a demo", 43px | Destination of the primary action across all 10 pages, with no form |
| `/en/privacy/` | 1.72 | no `<section>` | "Privacy Policy", 48px | Appropriate for what it is; only inherits the nav and the edges |

The pattern across the columns is the diagnosis: the platform pages are long (3.7 to 7.4 screens) and nearly every section lands between 0.6 and 1.2 screens — always close to one screen, never at one screen. They were written by content, not by screen.

## 7. Order of remediation

In this order, because each step makes the next one visible.

**1. Give the action a destination.** A form on `/en/book-a-demo/`, or a new label across the entire site. Without that, the rest is decoration: the site has nowhere to send someone who has been convinced. *(Finding 1.)*

**2. Section box and hero.** `--h` of 1 or 2 with `100svh`, 120/72 padding, h2 first; hero with h1 at y=120 and a button inside the screen; split `#report`, `#install`, `#evidence`; merge the short openings of safe-seat, cli and ci-gates. It is a shared layout: it fixes the one-section-per-screen rule and the hero rule across all 10 pages at once. *(Findings 2 and 4.)*

**3. Tokens.** Violet accent to blue; h1 gradient removed; one loaded font; a single type scale with a 13px floor; opaque 72px nav with a hairline and no blur; one 1280px `.wrap`. All in CSS, without touching content. *(Findings 3, 6, 7, 8, 10 and 11.)*

**4. Copy.** An h1 that names the category and the buyer across all 10 pages; headings with a noun and the sentences moved down to the line below; "Explore" removed; em dashes in the ticket panel removed. *(Findings 5, 13, 14 and 15.)*

**5. Motion and FAQ.** One loop per group, the rest on their final frame; FAQ open. *(Findings 9 and 12.)*

After that, run again at 1440×900, 1024×768, 820×1180, 768×1024 and 390×844, capture every width, and read every screenshot. The measure that tells whether it worked is simple: at every multiple of the screen height, a single section appears.

## 8. Method and limits of this review

**Method.** Real browser, all 10 English pages in the sitemap. Measurements came from a script running inside the page itself (section heights, content edges, font sizes and weights, computed colors and their hues, `<img>` count, canvases, running animations, `backdrop-filter` and gradients, sideways scroll checks, `href` on every link, form fields) at 1440×900, 1024×768, 768×1024 and 390×844. Screenshots read at 1440 and 390. The cited rules are those from the `taste-design-expert` file provided — each finding names the rule, where it appears, and the fix.

**Passed and worth recording:** no sideways scroll at any of the four widths; no dead links (empty `#` or missing `href`) across the 10 pages; no stock photographs; no invented social proof; no pop-ups, newsletter modals, chat bubbles or back-to-top buttons.

**Not verified:**

- **820×1180** was not measured — one of the five quality-gate widths is missing.
- **Performance budget.** `transferSize` returned zero for all 31 resources (cached or cross-origin), so transferred weight, JS and CSS KB, request count, FCP, LCP and layout shift were not measured. There are no performance numbers in this report.
- **Contrast.** Computed colors were read, but contrast ratios were not calculated. The 12px and 9px text likely fails on size before failing on contrast, but that is an inference, not a measurement.
- **Hover layer.** The presence of animations was checked, not the interactive behavior: button lifting 1px, linked card lifting 2px, underline drawing, 2px focus ring, or the state under `prefers-reduced-motion`.
- **Other locales.** The pt-br, zh-hant, ja and vi versions were not reviewed; the sitemap has 50 URLs and the 10 under `/en/` were measured. Because the layout is shared, the structural findings almost certainly hold across all locales, but the copy findings were not checked outside English.
- **Source code.** The rendered DOM was inspected, not the repository. Fixes are described by effect, not by file.
