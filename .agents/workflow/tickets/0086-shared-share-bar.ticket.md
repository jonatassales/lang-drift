# 0086: The share bar moves to @langdrift/web

**Spec:** .agents/workflow/specs/0014-console-ticket-phases.spec.md

**What to build:** the two-pixel share bar the website's panels draw moves to `packages/web`, so the console's widgets draw it the same way. The website imports it from there and looks the same.

**Blocked by:** None (can start immediately).

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `ShareBar` lives in `packages/web/src/ui/share-bar/`, exported by name and in the package export map.
- [ ] The website's panels import it from `@langdrift/web/ui/share-bar`; its local copy is gone.

### Verification

- `pnpm turbo run check typecheck test --filter=website --filter=@langdrift/web`, `drift check dead-code`.
