# 0003: FinOps stylesheet

**Spec:** .agents/specs/0001-design-tokens.spec.md

**What to build:** a FinOps app imports `@langdrift/tokens/finops.css` and gets the same shared tokens, shadcn bridge and Tailwind theme as Vision. The difference is FinOps's own product theme as `--fp-*`, which the shared primary role resolves to.

**Blocked by:** 0002 (Shadcn bridge and Tailwind theme).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] FinOps's product theme has its own source file in the themes folder, with primary and primary hover in both color schemes.
- [ ] `@langdrift/tokens/finops.css` is a public, self-contained export.
- [ ] In the FinOps stylesheet, `--lg-color-primary` and `--lg-color-primary-hover` reference `--fp-*`, and no `--viz-*` token appears.
- [ ] The contract test covers both product stylesheets.
