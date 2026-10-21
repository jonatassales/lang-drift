# Tailwind rules

Rules for Tailwind CSS in every web app. Loaded in every session.

1. **Layout stays native.** Spacing, sizing, grid, flex, breakpoints, border radius and the type scale use Tailwind's defaults. `@theme` never overrides them.

2. **Only design tokens override the theme.** `@theme` overrides colors and the other values `@langdrift/tokens` defines, through the bridge that package ships. An app never declares a theme value of its own.
