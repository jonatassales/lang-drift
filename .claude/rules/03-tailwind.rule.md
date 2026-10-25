# Tailwind rules

Rules for Tailwind CSS in every web app. Loaded in every session.

1. **Layout stays native.** Spacing, sizing, grid, flex, breakpoints, border radius and the type scale use Tailwind's defaults. `@theme` never overrides them.

2. **Only design tokens override the theme.** `@theme` overrides colors and the other values `@langdrift/tokens` defines, through the bridge that package ships. An app never declares a theme value of its own.

3. **Role utilities only.** Colors and fonts come from the tokens' roles and fixed groups: `bg-background`, `text-muted-foreground`, `border-border`, `text-primary`, `text-status-danger`, `text-classification-review`, `text-drift-negative`, `font-sans`, `font-mono`. Never a palette literal; the native palette is reset, so `text-red-500` produces nothing. Arbitrary values stay for one-off measurements, not for colors.

4. **No component classes.** No `.btn`, `.header-nav`, `.mobile-visible` or `.mobile-hidden` when a native utility expresses the behavior (`className="block sm:hidden"`). `@utility` exists only for a repeated single-purpose capability Tailwind lacks, and `@layer components` is not a styling strategy. Repeated multi-property UI becomes a React component.

5. **`dark:` follows `.dark`.** The tokens stylesheet declares the variant; an app never redeclares it and never switches themes through a `data-theme` attribute.

