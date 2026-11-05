# 0176: UI components rules compliance in packages/web

**Spec:** .agents/workflow/specs/0030-web-package-rules-compliance.spec.md

**What to build:**
Refactor UI components under `packages/web/src/ui/`:
1. `brand`: Extract `LogoMark` into `logo-mark.tsx` with its own `LogoMarkProps`. Keep `Brand` in `brand.tsx`. Re-export both from `index.ts`. Add `brand.unit.tsx`.
2. `hero-shader`:
   - Split `palette.ts` into single-function files: `mix.ts`, `rotate-hue.ts`, `derive-palette.ts`, `normalize-color.ts`, `read-palette.ts`. Re-export and declare types in `palette.ts`.
   - Add blank lines before returns inside GLSL string functions in `shaders.ts`.
   - Split `surface.ts` into single-function files: `compile-shader.ts`, `create-program.ts`, keeping `mountSurface` in `surface.ts`.
3. `language-switcher`:
   - Replace object shape type alias with named interfaces (`BaseLanguageSwitcherProps`, `NavigationLanguageSwitcherProps`, `ActionLanguageSwitcherProps`).
   - Eliminate all inline event handlers in `language-switcher.tsx` by defining named functions in component scope (`handleOptionClick`, `handleTriggerClick`, `handleContainerBlur`, `handleContainerKeyDown`, `handlePointerEnter`, `handlePointerLeave`).

**Blocked by:** none

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/web/src/ui/brand/logo-mark.tsx` declares `LogoMark`.
- [x] `packages/web/src/ui/brand/brand.tsx` declares `Brand`.
- [x] `packages/web/src/ui/brand/brand.unit.tsx` tests both components.
- [x] `packages/web/src/ui/hero-shader/` has single-function files (`mix.ts`, `rotate-hue.ts`, `derive-palette.ts`, `normalize-color.ts`, `read-palette.ts`, `compile-shader.ts`, `create-program.ts`).
- [x] `packages/web/src/ui/hero-shader/shaders.ts` satisfies return-spacing.
- [x] `packages/web/src/ui/language-switcher/language-switcher.tsx` uses interfaces and named handlers.
- [x] Quality check passes for `packages/web/src/ui/`.

### Verification

- `pnpm turbo run check typecheck test build --filter=@langdrift/web`
- `platform/cli/drift quality check`
