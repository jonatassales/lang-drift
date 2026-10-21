# Design tokens for Vision and FinOps

**Status:** ready-for-agent

## Problem Statement

LangDrift's products, Vision and FinOps, have no design tokens yet. The owner wants both to start from the palette proven in the earlier Exphere design package. Their web apps will be built from SmoothUI and shadcn components and styled with Tailwind, and those components expect shadcn's variables and a `.dark` class. Without a shared token package, each product would copy colors by hand, drift apart, and wire shadcn and Tailwind its own way.

## Solution

A new package, `@langdrift/tokens`, built with Style Dictionary, gives each product one stylesheet to import. That stylesheet carries:

- the shared tokens and the product's own product theme;
- both color schemes, with dark as the default;
- the bridge to shadcn's variables;
- the Tailwind theme, which keeps Tailwind's layout scales native and replaces its color palette with the tokens' roles and status colors.

Vision and FinOps start with identical values and can diverge later by editing only their product theme.

## User Stories

1. As a product developer, I want to import a single stylesheet for my product, so that I get every token, both color schemes, the shadcn bridge and the Tailwind theme in one line.
2. As a product developer on Vision, I want `@langdrift/tokens/vision.css`, so that my app gets Vision's product theme.
3. As a product developer on FinOps, I want `@langdrift/tokens/finops.css`, so that my app gets FinOps's product theme.
4. As a product developer, I want each product stylesheet to be self-contained, so that I never have to import shared tokens separately or in a particular order.
5. As a product developer, I want dark to be the default color scheme, so that a new page looks right in dark without extra setup.
6. As a product developer, I want the light color scheme in `:root` and dark under `.dark`, so that shadcn's `dark:` variants and the tokens switch together.
7. As a product developer, I want shadcn components to render with LangDrift colors as soon as I add them, so that I do not restyle each component.
8. As a product developer, I want SmoothUI components, which follow shadcn's variables, to pick up the same colors, so that both libraries look consistent.
9. As a product developer, I want Tailwind color utilities such as `bg-primary`, `text-muted-foreground` and `border-border`, so that I style with roles.
10. As a product developer, I want utilities for status colors, such as `text-status-danger`, so that feedback states use the fixed status colors.
11. As a product developer, I want Tailwind's native color palette to be unavailable, so that a stray `bg-red-500` cannot bypass the design system.
12. As a product developer, I want Tailwind's spacing, sizing, grid, flex, breakpoints, border radius and type scale to stay native, so that layout follows Tailwind's defaults.
13. As a product developer, I want `font-sans`, `shadow-sm` and `ease-standard` to come from the tokens, so that typography and effects match the design.
14. As a component author in a shared package, I want to style only with shared roles (`--lg-*`), so that the same component renders correctly in every product.
15. As a component author, I want `--lg-color-primary` to resolve to the current product's brand, so that a shared component shows Vision's or FinOps's primary without knowing which product loaded it.
16. As a component author, I want ring, sidebar primary, sidebar ring and the first chart color to follow the primary role, so that focus and highlights match the brand.
17. As a designer, I want each product theme to hold only the primary color and its hover state today, so that the brand is the one thing that differs between products.
18. As a designer, I want to change a product's brand by editing only its product theme, so that the other product and the shared tokens stay untouched.
19. As a designer, I want product theme tokens prefixed `--viz-*` for Vision and `--fp-*` for FinOps, so that each product's own values are recognizable.
20. As a designer, I want shared tokens prefixed `--lg-*`, so that anything common to all products is recognizable.
21. As a designer, I want the roles named after shadcn's vocabulary (background, foreground, card, popover, primary, secondary, muted, accent, destructive, border, input, ring, chart 1 to 5, sidebar and their foreground pairs), so that our names and the component libraries' names match.
22. As a designer, I want primary hover kept as an extra role, so that hover states have a defined color.
23. As a designer, I want the Exphere light palette preserved:
    - background `#FFFFFF`
    - card `#F9F9F9`
    - foreground `#27262A`
    - muted foreground `#535257`
    - border `#DDDDDD`
    - primary `#6D0AD6` and primary hover `#5408A5`
24. As a designer, I want a complete dark palette:
    - background gray 950
    - card, popover and sidebar gray 900
    - foreground gray 50
    - secondary, muted, accent, border and input gray 800
    - muted foreground gray 400
    - primary `#A66BF0` and primary hover `#B98AF4`
    - primary foreground gray 950
25. As a designer, I want the remaining light roles derived from the palette:
    - popover white
    - secondary, muted and accent gray 100
    - input and sidebar border matching border
    - sidebar matching card
    - every foreground pair aliasing foreground
    - primary foreground white
26. As a designer, I want the gray foundation from 50 to 950, including a new gray 400 `#a3a3a3`, so that roles draw from one scale.
27. As a designer, I want foundations available only as variables, never as Tailwind utilities, so that components use roles and not raw scale values.
28. As a designer, I want status colors fixed across products (success `#15803D`/`#4ADE80`, warning `#B45309`/`#FBBF24`, danger `#B91C1C`/`#F87171`, info `#1D4ED8`/`#60A5FA`, light/dark), so that feedback means the same everywhere.
29. As a designer, I want destructive to alias the danger status color, so that destructive actions and danger feedback look the same.
30. As a designer, I want chart colors 2 to 5 to be their own hues (`#0D9488`, `#D97706`, `#0284C7`, `#E11D48` in light; `#2DD4BF`, `#FBBF24`, `#38BDF8`, `#FB7185` in dark), so that a data series is never mistaken for a status.
31. As a designer, I want the body font family `system-ui, sans-serif`, the small shadow `0px 1px 2px #00000014` and the standard easing `cubic-bezier(0.2, 0, 0, 1)` kept from Exphere, so that the existing look carries over.
32. As a designer, I want no dimension or font-size tokens, so that layout and type scale are Tailwind's.
33. As a designer, I want aliases kept as `var(...)` references in the output, so that changing one role updates everything that points to it.
34. As a maintainer, I want the build to fail when a token references a missing or circular target, so that broken aliases never ship.
35. As a maintainer, I want the build output to be deterministic, so that rebuilding unchanged sources yields identical files.
36. As a maintainer, I want the package to build and test through the repository's Turborepo `build` and `test` tasks, so that CI covers it with no extra wiring.
37. As a maintainer, I want the token sources in the Design Tokens Community Group format, so that the files stay portable between tools.
38. As a maintainer, I want the `03-tailwind` rule, so that agents and people keep Tailwind to layout and let only tokens override the theme.
39. As a maintainer, I want the frontend rule to say an app's CSS holds its Tailwind entry point, its product tokens import and Tailwind directives, so that the rule matches how tokens are delivered.
40. As a maintainer, I want the AGENTS.md design section to describe the prefixes, the shadcn vocabulary and the dark default, so that the next session starts from the current design.

## Implementation Decisions

- **New package `@langdrift/tokens`** in the shared packages area, not inside a product, because both products consume it.
- **Token sources** are DTCG JSON. They keep Exphere's string values: hex colors, the font-family list, the shadow string and the `cubic-bezier()` string. Shared tokens and each product theme live in separate source files. Product themes sit in the package's themes folder, per AGENTS.md.
- **Naming.** The CSS variable is the prefix plus the token path in kebab-case: `--lg-color-gray-950`, `--lg-color-muted-foreground`, `--lg-status-danger` (mirroring the existing `--lg-status-*` convention), `--viz-color-primary`, `--fp-color-primary-hover`.
- **Color schemes.** Every color token that changes with the scheme has a light and a dark value. Light is emitted under `:root` and dark under `.dark`. Apps render with the `.dark` class by default; the package itself does no scheme detection, toggling or persistence.
- **Product theme binding.** Each product stylesheet declares its product theme tokens, then binds the shared primary and primary hover roles to them in both schemes. Ring, sidebar primary, sidebar ring and chart 1 alias the shared primary role, never the product theme directly.
- **Shadcn bridge.** Each product stylesheet declares every shadcn variable (`--background`, `--foreground`, `--card`, …, `--sidebar-ring`) as a `var()` reference to the matching `--lg-*` role. It does not declare `--radius`.
- **Tailwind theme.** Each product stylesheet resets the native palette with `--color-*: initial`, then maps:
  - every role and status color to a `--color-*` theme variable;
  - `--font-sans` to the body font family;
  - `--shadow-sm` to the small shadow;
  - `--ease-standard` to the standard easing.

  Theme variables reference the shadcn or `--lg-*` variables, so they follow the color scheme at runtime. Nothing else in Tailwind's theme is touched; spacing, radius and the type scale stay native.
- **Outputs.** Two public exports, `@langdrift/tokens/vision.css` and `@langdrift/tokens/finops.css`, each self-contained. No JSON or JavaScript export.
- **Build.** Style Dictionary builds both stylesheets through its API, from a build program that creates one dictionary per product. Its CLI builds a single dictionary, and each product needs its own because both bind the shared primary role. A custom format emits the schemes, the bridge and the Tailwind theme together, which the predefined CSS format cannot. Build output goes to the package's `dist`, which Turborepo already caches and Git ignores. The build fails on broken or circular references; there is no other value validation.
- **Rules and docs.**
  - Add `03-tailwind.rule.md` with the approved text: layout stays native, and only design tokens override the theme.
  - Update rule 2 of `02-frontend.rule.md` to say that an app's CSS holds its Tailwind entry point, its product's tokens import and its Tailwind directives.
  - Update the "Design and tokens" section of AGENTS.md.
- **Domain model.** Terms follow the glossary in `.agents/context.md` (product, shared token, product theme, foundation, role, status color, color scheme). ADR 0001 records the shadcn vocabulary, the Tailwind palette reset and the choice of Style Dictionary.

## Testing Decisions

- **One seam: the published stylesheets, consumed as an app would consume them.** A contract test builds the package into an isolated copy, links it into a disposable consumer, resolves both public exports through the package's `exports` map and inspects the CSS text. It never looks inside the Style Dictionary config or its formats.
- **What the contract checks:**
  - each product's primary and primary hover in both schemes;
  - the shared primary role bound to the right product theme;
  - the light scheme under `:root` and dark under `.dark`;
  - the Exphere light values;
  - the dark palette and status colors;
  - destructive aliasing danger;
  - every shadcn variable present and referencing an `--lg-*` role;
  - the palette reset preceding the token colors;
  - no `--radius` and no dimension or font-size tokens;
  - no foundation exposed as a Tailwind color;
  - identical output on rebuild;
  - a failed build on a broken alias.
- **A good test** asserts on what a consumer receives (variable names, values, selectors and references), never on the order of internal steps.
- **Prior art:** the contract test of Exphere's design package, which builds an isolated copy, resolves the public exports from a disposable consumer and checks deterministic rebuilds and rejected graphs.
- The test runs under the repository's `test` task with Node's built-in test runner. No new test framework.

## Out of Scope

- Any web app, component, or `packages/ui`. Components come from SmoothUI and shadcn, and specific components get their own package later.
- Running `shadcn init` or SmoothUI setup in an app; this spec delivers what they will consume.
- Light/dark detection, toggling, persistence, or `prefers-color-scheme`.
- Dimension, spacing, radius and font-size tokens.
- A JSON or JavaScript export of the tokens.
- Validating token value formats beyond what Style Dictionary does.
- Publishing the package to a registry.
- Product theme tokens beyond primary and primary hover.

## Further Notes

- Commits that touch the package need a `Ticket: <NNNN>` trailer naming a ticket from `to-tickets`.
- Dark brand values (`#A66BF0`, `#B98AF4`), gray 400, chart hues and status colors were proposed during design and approved by the owner. Contrast was reasoned about but not measured.
- Source: the owner's design session and the Exphere design package, read from the owner's machine.
