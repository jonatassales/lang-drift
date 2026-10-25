# 0047: Technology icons for Go, Node.js, NATS, Redis and Next.js

**Spec:** .agents/workflow/specs/0007-microruntimes-hero.spec.md

**What to build:** the web package offers five technology marks, Go, Node.js, NATS, Redis and Next.js, as icon components with the same interface as the platform icons, so a component row on the website can show a mark that follows the text color in light and dark without image files, an attribution README or a third-party icon dependency.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

### Icons

- [ ] Five icon components in their own icon family beside the platform icons, each a simplified monochrome rendering of the mark: a `className` prop, a 24-unit viewBox, strokes and fills in `currentColor`, `aria-hidden`.
- [ ] The family's props interface has the same shape as the platform icons' props and lives with the family.
- [ ] A barrel exports the five icons by name, and the package export map gains the family's entry, following the platform icons' entry.
- [ ] Each mark stays recognizable at 16px in both color schemes; none carries a brand color.

### Tests

- [ ] A `renderToString` unit test like the platform icons' test: each icon renders an `svg` with `aria-hidden` and the class prop applied.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=@langdrift/web`, `drift check dead-code` (the icons may be unused until ticket 0050; if the check reports them, note it in the ticket and let 0050 close it).
