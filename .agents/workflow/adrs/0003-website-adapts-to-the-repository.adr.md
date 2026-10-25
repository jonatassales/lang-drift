---
status: accepted
---

# The LangDrift website adapts to this repository, not the other way round

The marketing website is copied from the earlier `langdrift` monorepo, where it used its own token package (`--ld-*` variables, `data-theme`, its own role names) and a component package, `@pack/react`. This repository stays the source of truth: workspaces live under `apps/*` (the never-used `products/` layout and the Vision and FinOps themes are removed), `@langdrift/tokens` owns every color, and the website keeps its architecture and look while its token layer is remapped onto shadcn roles and the `.dark` class (ADR 0001).

To keep the look, the website's neutrals become the shared roles' values, and the website's domain colors become two fixed groups beside status colors: classification colors and drift colors. They never take the brand. The brand is the Exphere purple, in a single `website` theme. Components the website imports move to `@langdrift/web`, and the website imports them under that name, at the subpaths `@pack/react` used. Its first version kept the `@pack/react` specifiers through a pnpm workspace alias so the copy could be diffed against the original; once the copy was verified, the alias went, because tao names its packages `@langdrift/*`.

The website brings about 1,700 lines of hand-written CSS, which the frontend rule forbids. Those existing files, and the ones `@langdrift/web` inherits, are allowed by a scoped exception; no new CSS may be added, and converting them is left for later.

## Considered Options

- **Porting `@pack/design-tokens` beside `@langdrift/tokens`.** It would keep the website byte-identical, but two token systems would own colors, and the `data-theme` and custom role names ADR 0001 rejected would return.
- **Mapping classification and drift colors onto status colors.** Fewer tokens, but the website's look would change and Product Vision's meanings would blur into generic success, warning and danger.
- **Converting the CSS to Tailwind during the copy.** It honors the frontend rule at once, but the look could drift from the original and the copy could no longer be checked against it.
- **Keeping the `@pack/react` alias.** It kept the website's imports identical to langdrift's, but left a foreign package scope in a repository that names its own `@langdrift/*`.
