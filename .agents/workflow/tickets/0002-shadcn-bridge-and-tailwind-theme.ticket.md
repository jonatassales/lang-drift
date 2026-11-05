# 0002: Shadcn bridge and Tailwind theme

**Spec:** .agents/workflow/specs/0001-design-tokens.spec.md

**What to build:** a Vision app that adds shadcn or SmoothUI components and Tailwind gets LangDrift colors with no extra setup. The product stylesheet declares every shadcn variable as a reference to its `--lg-*` role. It resets Tailwind's native color palette and exposes only roles and status colors as Tailwind colors. It maps `--font-sans`, `--shadow-sm` and `--ease-standard`, and leaves Tailwind's layout, radius and type scales native.

**Blocked by:** 0001 (Vision stylesheet with shared tokens and its product theme).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Every shadcn variable, from `--background` to `--sidebar-ring`, is declared and references the matching `--lg-*` role.
- [ ] `--radius` is not declared.
- [ ] The Tailwind theme resets the native palette before declaring any token color.
- [ ] Every role and status color is a Tailwind color. No gray foundation is.
- [ ] `--font-sans`, `--shadow-sm` and `--ease-standard` come from the tokens. Nothing else in Tailwind's theme is declared.
- [ ] The contract test covers all of the above on the published Vision stylesheet.
