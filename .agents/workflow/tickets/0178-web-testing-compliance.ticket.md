# 0178: Testing suite rules compliance across packages/web

**Spec:** .agents/workflow/specs/0030-web-package-rules-compliance.spec.md

**What to build:**
Update test descriptions in unit tests under `packages/web/` so every `describe(...)` starts with `web/ui/...` or `web/vendors/...` per `.claude/rules/10-testing.rule.md` (Rule 4: "A TypeScript test description starts with the application or package name, follows the module path and ends with the file name, never the exported symbol."):
- `packages/web/src/ui/agent-orb/agent-orb.unit.tsx`
- `packages/web/src/ui/detail-list/detail-list.unit.tsx`
- `packages/web/src/ui/elevated-frame/elevated-frame.unit.tsx`
- `packages/web/src/ui/eyebrow/eyebrow.unit.tsx`
- `packages/web/src/ui/hero-shader/hero-shader.unit.tsx`
- `packages/web/src/ui/icons/technology/technology-icons.unit.tsx`
- `packages/web/src/ui/info-cards/info-cards.unit.tsx`
- `packages/web/src/ui/language-switcher/language-switcher.unit.tsx`
- `packages/web/src/ui/scroll-hint/scroll-hint.unit.tsx`
- `packages/web/src/ui/section-heading/section-heading.unit.tsx`
- `packages/web/src/ui/share-bar/share-bar.unit.tsx`
- `packages/web/src/vendors/smoothui/mac-terminal.unit.tsx`

**Blocked by:** 0177

**Status:** ready-for-agent

- [x] All `packages/web` test files start their `describe` calls with `web/ui/...` or `web/vendors/...`.
- [x] `drift quality check 10-testing` reports zero violations in `packages/web`.
- [x] All package and repo tests pass.

### Verification

- `pnpm turbo run test --filter=@langdrift/web`
- `platform/cli/drift quality check 10-testing`
