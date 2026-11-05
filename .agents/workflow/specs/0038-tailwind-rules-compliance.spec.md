# Every web surface complies with the styling rules

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

> **Note, 2026-11-05.** The owner closed 0276 as won't fix: the work goes to the published packages, `@langdrift/cli` and `@langdrift/core`, until the runtime runs.

## Tickets

- [ ] 0272: the pill button and the text link are shared primitives
- [x] 0273: the vendored tabs and accordion draw with utilities
- [ ] 0274: every color and static style reads from a role or a utility
- [ ] 0275: the website uses the shared primitives and one cn
- [ ] 0276: no hand-written stylesheet remains

## Problem Statement

`.claude/rules/03-tailwind.rule.md` says Tailwind is the styling system, that no hand-written stylesheet exists, that colors come from the tokens' roles, that `cn()` has one owner and that repeated multi-property UI becomes a React component. The audit of 2026-11-04 measured the rule at 81.2% and found:

- Three vendored SmoothUI stylesheets in `packages/web/src/vendors/smoothui/` (`animated-tabs.css`, `basic-accordion.css`, `smoothui.css`, 209 lines of selectors, hover, focus and media rules). `smoothui.css` reaches every app through `packages/web/src/styles.css`, which both locale layouts import; the two component stylesheets are imported by `animated-tabs.tsx` and `basic-accordion.tsx`. Rule 8 already records them as violations for remediation; spec 0027 left them in place and spec 0030 did not touch them.
- Ten uses of `smooth-primary-button`, `smooth-text-link` and `smooth-footer` across nine website components (`book-a-demo/form.tsx`, `plans-section.tsx`, `website-footer.tsx`, `home/home-product-section.tsx`, `finops/source-list.tsx`, `product-page/section-outro.tsx`, `microruntimes/session-control.tsx`), where the unlayered stylesheet still beats Tailwind's utilities. `plans-section.tsx` carries `!` overrides marked "Transitional" for that reason (spec 0027, SmoothUI precedence).
- A second `cn()` in `apps/website/src/app/template/formatters/cn.fmt.ts`, unused, with `clsx` and `tailwind-merge` still listed in the website's manifest; rule 6 names the one in `web/utilities` as the only one. Spec 0027 kept it because the owner keeps the formatters; the rule outranks that choice.
- Arbitrary colors written as `color-mix(…)` of a role variable in `packages/web/src/ui/language-switcher/language-switcher.tsx` (the trigger and current option backgrounds, the container border, background and shadow) and a radial gradient mixing the primary with `white` and `black` in `packages/web/src/ui/agent-orb/agent-orb.tsx`; the existing arbitrary-color detector does not see `color-mix` or `var(--…)` values.
- Fixed hex colors and inline style objects in `packages/web/src/ui/brand/brand.tsx` (`#f5f5f5`, `#111111`, the wordmark's layout, font and spacing as a `style` object), static `display` and `white-space` style objects in `packages/web/src/vendors/smoothui/text-morph.tsx`, and `style={{ colorScheme: 'dark' }}` on the website's `<html>` where the console uses the `scheme-*` utilities and the header writes `colorScheme` imperatively on top of it.

The vendored stylesheets are unlayered, so every utility that must beat them needs `!`; a role color written as `color-mix` escapes the detectors and the tokens; and two `cn()` implementations can drift apart. Until these are remediated, a new component in `packages/web` or an app starts from a tree the rule rejects.

## Solution

No hand-written stylesheet remains anywhere in the tree. `AnimatedTabs` and `BasicAccordion` draw with Tailwind utilities and token roles; the three SmoothUI stylesheets and their imports are deleted. The website's pill button and text link become shared primitives in `packages/web`, used through `asChild` by links, buttons and the trial plan's call to action, and the `!` overrides and their "Transitional" comments go with them. The footer carries its measurements as utilities. Every color reads from a role: translucent role surfaces take Tailwind's opacity modifier on the role, the orb's fallback gradient takes two new app-theme tokens that `tokens` builds, and the brand's fixed tones become the foreground roles of the scheme they imitate. Static style objects become utilities; a `style` prop keeps only a measurement computed at runtime. The website drops its `cn()` and the two dependencies behind it. What a visitor sees does not change, within the computed-style tolerance recorded below.

## User Stories

1. As the owner, I want `drift quality check tailwind` to report zero findings, so that the styling rule is verified by the seam the repository already has.
2. As the owner, I want no `.css` file in the tree other than Tailwind entries, the tokens import and the approved utilities file, so that rule 7 and rule 9 hold without exception.
3. As the owner, I want every `smooth-*` class gone, so that no unlayered stylesheet forces `!` on a utility.
4. As a maintainer of `packages/web`, I want `AnimatedTabs` to express its underline, pill and segment variants as utilities, so that a variant is readable in the component that renders it.
5. As a maintainer of `packages/web`, I want `BasicAccordion` to carry its trigger, chevron and region styles as utilities, so that the chevron's rotation and the reduced-motion branch are visible where the markup is.
6. As a website developer, I want a pill variant of the shared `Button`, so that the demo form's submit and the trial plan's call to action share one recipe and a consumer override composes through `cn()` instead of `!`.
7. As a website developer, I want a shared `TextLink`, so that the six text links on the home, FinOps, product and pricing pages and the session control render the same underline without a class from a stylesheet.
8. As a website developer, I want the footer's width, padding and top border as utilities on the footer, so that its layout is read in its component.
9. As a developer of any app, I want one `cn()` from `web/utilities`, so that conditional classes merge the same way everywhere and the website's manifest lists no `clsx` or `tailwind-merge`.
10. As a maintainer, I want the language switcher's translucent surfaces written as role utilities with opacity modifiers, so that a color is a role and a detector can see it.
11. As a maintainer, I want the orb's fallback gradient drawn from tokens the build emits, so that the brand's highlight and shade live in `tokens` and never in a component.
12. As a maintainer, I want the brand wordmark styled by utilities and colored by the surrounding text, so that no hex literal lives in a component.
13. As a maintainer, I want the text morph's static display and white-space as utilities, so that the only `style` props left hold runtime measurements.
14. As a website visitor, I want the page's color scheme declared through the `scheme-*` utilities, as in the console, so that the scheme follows the `dark` class without an inline style or an imperative write.
15. As a visitor of the website, I want every page to look as it did, so that the refactor costs me nothing.
16. As a keyboard user, I want tabs, accordion triggers, buttons and links to keep their focus ring, so that focus stays visible after the rewrite.
17. As a visitor who prefers reduced motion, I want the tabs' indicator and the accordion's chevron to stay still, so that the rewrite keeps the reduced-motion branch.
18. As a visitor on a narrow viewport, I want the tabs to keep scrolling horizontally with a thin scrollbar, so that the tab list fits.
19. As the owner, I want a detector for each class of violation found, so that a `smooth-*` class, a stylesheet import in a component, a `color-mix` or variable color in an arbitrary value, a static value in a `style` prop and a styling dependency in an app manifest are reported by the seam.
20. As the owner, I want the added tokens recorded in the tokens package and the glossary question raised, so that the app theme's definition stays decided by me.

## Implementation Decisions

- **Shared primitives (rules 10, 11; React rule 11).**
  - `Button` in `packages/web` gains a `pill` variant that reproduces the SmoothUI primary button: full radius, 38px minimum height, 17px horizontal padding, 13px text at weight 650, 7px gap, the primary roles, and on hover a saturation of 1.18 with a one-pixel lift over a 180ms transition, off under reduced motion. The variant's classes compose through `cn()`, so a consumer's `flex`, `min-h-11`, `justify-between`, `px-[15px]` and `py-3` win without `!`.
  - A `TextLink` primitive in `packages/web` reproduces the SmoothUI text link: inherited color, 13px text, an underline in the current color at 25%, offset 5px. It renders an anchor by default and lends its look to its one child through `asChild`, so a next-intl `Link` and a `button` keep their semantics. It exports its props interface.
  - The website's demo form submit, the trial plan's call to action and the six text links consume these primitives; the "Transitional" comments and the `!` modifiers written against SmoothUI are removed.
- **Vendored components (rules 7, 8, 11).**
  - `AnimatedTabs` keeps its markup, ARIA relations, keyboard model and `data-variant` attribute, and draws every declaration of its stylesheet as utilities: the list's bottom border, gap, padding, horizontal overflow, thin scrollbar colored by the border role and inline overscroll containment; the tab's 44px minimum height, padding, weight 550, muted color turning to the foreground on hover and when selected, and the 150ms color transition; the underline indicator as a 2px primary bar at the bottom, and for `pill` and `segment` an indicator that fills the tab with the card role and a border; the focus ring as a 2px foreground outline. `motion-reduce:` turns the transitions off. The scrollbar declarations, which Tailwind has no utility for, are arbitrary properties that reference the border role's theme variable.
  - `BasicAccordion` keeps its markup and ARIA and draws its stylesheet as utilities: the divider between items from the border role, the trigger's full width, 44px minimum height, weight 550, space-between layout, hover underline with a 4px offset and 3px-offset focus outline, the chevron's 200ms rotation to 180° when expanded, the region's hidden overflow, the content's padding and `anywhere` overflow wrapping. `motion-reduce:` turns the rotation's transition off.
  - The component stylesheets are deleted with their imports, then `smoothui.css` and its import in the package stylesheet, which keeps the Tailwind import, the tokens import and the `@source` directive.
- **Colors (rules 3, 4).**
  - A translucent role surface is the role with Tailwind's opacity modifier: the switcher's current option and trigger take the background role at 88%, its container the card role at 72% with a border of the border role at 80%, and its shadow keeps its measurements with the foreground role at 8% as the shadow color. Mixing a color with transparent gives the same color at that alpha in any interpolation space, so the result is the same.
  - The orb's fallback gradient mixes the primary with white and black, which the tokens do not carry. The app theme gains two roles, `primary-highlight` (the primary at 35% over white, today's center) and `primary-shade` (the primary at 60% over black, today's rim), with the values today's mix computes, in both schemes; the tokens build emits them like the other theme roles and `packages/web` exposes them as `bg-primary-highlight` and `bg-primary-shade`. The orb's gradient becomes a radial utility from the highlight through the primary at 50% to the shade.
  - The brand's `tone` prop is removed. Its only caller, the website header, passes the active theme as the tone, so the fixed hex it selected always equals the scheme's foreground; the header passes the foreground role through `className`, as the console already does, and the wordmark takes the surrounding text color. The wordmark's layout, font family, weight, size, letter spacing, line height and gap become utilities.
- **Style props (rule 7).** A static declaration in a `style` object is styling written as CSS and becomes a utility: the text morph's inline-block display and preserved white space, the brand's layout and typography. A `style` prop stays only for a value computed at runtime (the orb's dimension, the share bar's width, the ledger segment's basis, the editor's maximum height).
- **Color scheme (rule 5).** The website's `<html>` takes `scheme-light dark:scheme-dark` beside its `dark` class, as the console's does, and the header's theme effect stops writing `colorScheme` on the document element, since the class toggle now carries it.
- **One `cn()` (rule 6).** The website's `cn.fmt.ts` and its `clsx` and `tailwind-merge` dependencies are removed; every website import of `cn` already comes from `web/utilities`.
- **Detectors (02-tooling rule 4).** `drift quality check` gains, for `03-tailwind`: a `smooth-` class in a TypeScript or TSX file; a stylesheet import in a TSX or TS file; an arbitrary value of a color-bearing utility that contains `color-mix(`, `var(--`, `white` or `black`; a static string or number in a `style` object; `clsx` or `tailwind-merge` in an app manifest. The arbitrary-color detector keeps its hex, `rgb`, `hsl` and `oklch` forms.
- **Order.** The primitives land first, then the two vendored components, then the colors and style props in `packages/web` and `tokens`, then the website consumers once spec 0041 has moved the website tree, then the contract that deletes `smoothui.css` and adds the detectors.

## Testing Decisions

- **Visual parity is the main seam**, as in spec 0027: before each ticket, the affected pages are captured in the built-in browser at 1280px and 375px, in light and dark, in `en`, and every affected element's computed styles and box are recorded under `.audits/`; after the ticket, the same capture is compared. A difference is accepted only where this spec names it: the orb's gradient interpolated by the browser's radial gradient instead of `color-mix`, and the brand wordmark drawn from the foreground role instead of its fixed tone. No browser runner is added and nothing of the capture enters the repository; the pull request summarizes the comparison.
- **The unit tests are the behavior seam.** `Button` gets a test that renders the pill variant with the plans section's overrides and asserts the merged class list holds the override and not the variant's conflicting utility. `TextLink` gets a test for its anchor form and its `asChild` form. `AnimatedTabs` and `BasicAccordion` keep their existing tests and gain assertions for each variant's indicator classes, the selected and focused states and the reduced-motion utilities; a good test asserts the rendered markup and attributes, never a stylesheet. The tokens test asserts the two new roles are emitted in both schemes. Prior art: `packages/web/src/ui/button/button.unit.tsx`, `language-switcher.unit.tsx`, `packages/tokens/tokens.test.js`.
- **The detectors are the rules seam.** `platform/cli/drift quality check tailwind` reports zero findings at the end, and `pnpm turbo run check typecheck test build` for `web`, `tokens` and `website` stays green after every ticket. `drift check localization` is unaffected and stays green.
- **Review-enforced.** The remaining `!` modifiers that beat a third-party web component or a browser default, and the arbitrary properties for the scrollbar, are not rule violations and no detector reports them; the review checks that none beats a stylesheet this spec deleted.

## Out of Scope

- The console, which uses no `smooth-*` class and whose Tailwind entry already complies; it only receives the shared primitives.
- Rules other than `03-tailwind`: the SmoothUI components' effects and arrow functions (spec 0040), their file structure (spec 0030), and the website tree move (spec 0041).
- The `themeColor` metadata hex in both locale layouts and in the header's meta sync: metadata values, not styling.
- The hex sentinel in the hero shader's color normalizer and the hex fallback in the console's palette computation: canvas probes, not utilities; the fallback's invented value belongs to the general rules.
- The `vidropro`, `orbo-voice` and `sinapsi` web components' own styles, which live in their shadow trees.
- New design tokens beyond the two orb roles.

## Further Notes

- **Strict readings taken.** (1) A static value in a `style` object is treated as authored CSS under rule 7, although the rule names stylesheets, modules and CSS-in-JS libraries; a runtime measurement stays in `style`. (2) The scrollbar's width and color are written as arbitrary properties referencing the border role's theme variable, since Tailwind ships no scrollbar utility and an `@utility` would need an app's utilities file; this is not an arbitrary color because the value is a role. (3) The brand's `tone` prop is removed rather than mapped to two new fixed tokens, because its only caller selects the active scheme's color; the owner may prefer fixed tokens if a caller ever needs the dark wordmark on a light surface.
- **Open question for the owner.** The glossary defines the app theme as "only the primary color and its hover state". The orb's highlight and shade are derived from the primary and are added to the theme here; either the glossary entry is amended through `domain-modeling`, or the two roles become a fixed group beside `rainbow` and stop following the primary. The spec takes the theme reading; the owner can overturn it at PR review.
- **Owner-approved deviation to confirm.** Spec 0027 approximated two colors without a token; this spec adds tokens instead of approximating, because the brief asks for fidelity and the tokens build is the place a role is born.
- **Spec 0027's choice on `cn()`** (kept for the formatters) is overturned by rule 6, as the general rules' precedence requires.
