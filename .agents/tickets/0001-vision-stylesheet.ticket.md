# 0001: Vision stylesheet with shared tokens and its product theme

**Spec:** .agents/specs/0001-design-tokens.spec.md

**What to build:** a new `@langdrift/tokens` package, built by Style Dictionary, that publishes `@langdrift/tokens/vision.css`. A Vision app imports that one file and gets every shared token as `--lg-*` and Vision's product theme as `--viz-*`. Light comes in `:root` and dark in `.dark`. The shared primary role resolves to Vision's brand in both color schemes.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] The package builds and tests through the repository's Turborepo `build` and `test` tasks.
- [ ] Token sources are DTCG JSON. Shared tokens and Vision's product theme live in separate files, the theme in the package's themes folder.
- [ ] Gray foundations 50 to 950, including gray 400, are emitted as `--lg-color-gray-*`.
- [ ] Every role from the spec has its light value under `:root` and its dark value under `.dark`, with aliases kept as `var(...)`. That covers the shadcn names plus primary hover.
- [ ] `foreground` is kept, and every `*-foreground` pair aliases it where the spec says so.
- [ ] Status colors are emitted as `--lg-status-*` with their light and dark values, and destructive aliases danger.
- [ ] The body font family, the small shadow and the standard easing are emitted as `--lg-*`. There are no dimension or font-size tokens.
- [ ] `--viz-color-primary` and `--viz-color-primary-hover` carry the approved light and dark values. `--lg-color-primary` and `--lg-color-primary-hover` reference them.
- [ ] Ring, sidebar primary, sidebar ring and chart 1 reference `--lg-color-primary`.
- [ ] A contract test builds an isolated copy and resolves `@langdrift/tokens/vision.css` from a disposable consumer. It checks the values above, an identical rebuild, and a failed build on a broken alias.
