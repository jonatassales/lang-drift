# 0004: Tailwind rule, frontend rule and design docs

**Spec:** .agents/specs/0001-design-tokens.spec.md

**What to build:** agents and people working on a web app find the token and Tailwind conventions written where they look first.

- A new Tailwind rule says layout stays native and only design tokens override the theme.
- The frontend rule says an app's CSS holds its Tailwind entry point, its product's tokens import and its Tailwind directives.
- The AGENTS.md design section describes the prefixes, the shadcn vocabulary and the dark default.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The Tailwind rule exists with the text approved during grilling.
- [ ] Rule 2 of the frontend rule names the product tokens import instead of declaring `--lg-*` tokens.
- [ ] The AGENTS.md "Design and tokens" section covers `--lg-*`, `--viz-*`, `--fp-*`, the shadcn role vocabulary, the dark default and the one-import-per-product stylesheets.
- [ ] The AGENTS.md general rules overview mentions the Tailwind rule.
