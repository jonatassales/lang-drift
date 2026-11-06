---
status: accepted
---

# Design tokens follow the shadcn vocabulary and own every Tailwind color

The web apps are built from SmoothUI and shadcn components, which read unprefixed shadcn variables (`--background`, `--primary`) and switch to dark through a `.dark` class. `@langdrift/tokens` therefore names its roles after shadcn's, under our own prefixes (`--lg-*` shared, `--viz-*` and `--fp-*` per product theme). Each product's stylesheet bridges those roles to the shadcn variables and to Tailwind's theme. It puts the light color scheme in `:root` and dark in `.dark`; apps open with `.dark` set, so dark is the default. The tokens are built with Style Dictionary.

Tailwind is for layout only. Its spacing, sizing, radius and type scales stay native, while its default color palette is reset, so the only color utilities are the roles and status colors from the tokens.

## Considered Options

- **Terrazzo instead of Style Dictionary.** It validates token values and emits every color scheme into one file natively, but Style Dictionary was chosen for its far larger adoption and ecosystem. Value validation is dropped and a contract test covers the output.
- **`data-theme` attributes and our own role names**, as the earlier Exphere design package used. Both would fight shadcn: its `dark:` variants only fire under `.dark`, and its components would need a rename layer for every role.
- **Keeping Tailwind's native color palette.** It would let `bg-red-500` bypass the design system.
