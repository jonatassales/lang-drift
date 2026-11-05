# The web package complies with the repository rules

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

`packages/web` contains violations of `.claude/rules/00-general.rule.md`, `.claude/rules/05-react.rule.md`, `.claude/rules/04-typescript.rule.md`, and `.claude/rules/10-testing.rule.md`:
- `packages/web/src/ui/brand/brand.tsx` declares more than one function (`LogoMark` and `Brand`) in a single file.
- `packages/web/src/ui/hero-shader/palette.ts` declares five functions (`mix`, `rotateHue`, `derivePalette`, `normalizeColor`, `readPalette`) in a single file.
- `packages/web/src/ui/hero-shader/shaders.ts` has `return` statements without preceding blank lines in GLSL shader strings.
- `packages/web/src/ui/hero-shader/surface.ts` declares multiple functions (`compileShader`, `createProgram`, `mountSurface`) in a single file.
- `packages/web/src/ui/language-switcher/language-switcher.tsx` declares a `type` alias for an object shape and uses inline event handlers (`onClick`, `onBlur`, `onKeyDown`, `onPointerEnter`, `onPointerLeave`).
- `packages/web/src/vendors/smoothui/animated-tabs.tsx` destructures its props in the function parameter list, uses inline event handlers (`onClick`, `onFocus`, `onKeyDown`), and omits blank lines before returns.
- `packages/web/src/vendors/smoothui/basic-accordion.tsx` destructures its props in the function parameter list and uses an inline event handler (`onClick`).
- `packages/web/src/vendors/smoothui/mac-terminal.tsx` uses arrow functions and has statements without blank lines before `return`.
- `packages/web/src/vendors/smoothui/text-morph.tsx` declares multiple helper functions (`segmentGraphemes`, `tokenizeText`, `buildTokens`) alongside `TextMorph`.
- Unit test files across `packages/web/src/ui/` and `packages/web/src/vendors/` describe their suites with bare symbol names rather than the required package and file path prefix (`web/ui/...` or `web/vendors/...`).

## Solution

Refactor `packages/web` to bring all components, utilities, and tests into complete compliance with `00-general`, `05-react`, `04-typescript`, and `10-testing`, keeping exact runtime behavior and visual fidelity.

## User Stories

1. As the platform owner, I want `drift quality check` to report zero rule violations across `packages/web` for rules 00, 04, 05, and 10.
2. As a maintainer, I want every source file in `packages/web` to declare exactly one function or component named after its file, so concerns are isolated and easy to locate.
3. As a consumer of `@langdrift/web`, I want existing public entry points and exports to remain completely backward-compatible.

## Implementation Decisions

- **Brand UI (00.2, 00.4):**
  - Extract `LogoMark` into `logo-mark.tsx`.
  - Keep `Brand` in `brand.tsx`.
  - Re-export `Brand` and `LogoMark` and their props interfaces in `index.ts`.
  - Add unit test `brand.unit.tsx`.
- **Hero Shader UI (00.1, 00.2, 00 Statement Spacing):**
  - Extract `mix` into `mix.ts`.
  - Extract `rotateHue` into `rotate-hue.ts`.
  - Extract `derivePalette` into `derive-palette.ts`.
  - Extract `normalizeColor` into `normalize-color.ts`.
  - Extract `readPalette` into `read-palette.ts`.
  - Re-export palette functions and keep types in `palette.ts`.
  - Add blank lines before returns in `shaders.ts`.
  - Extract `compileShader` into `compile-shader.ts` and `createProgram` into `create-program.ts`, keeping `mountSurface` in `surface.ts`.
- **Language Switcher UI (04-typescript, 05.4, 05.5, 05.6):**
  - Replace `type LanguageSwitcherProps = { ... } & (...)` with interface definitions `BaseLanguageSwitcherProps`, `NavigationLanguageSwitcherProps`, and `ActionLanguageSwitcherProps`, unioned into `LanguageSwitcherProps`.
  - Replace all inline event handlers (`onClick`, `onBlur`, `onKeyDown`, `onPointerEnter`, `onPointerLeave`) with named functions in component scope.
- **Vendors SmoothUI (00.2, 05.3, 05.4, 05.5, 05.6):**
  - `animated-tabs.tsx`: Take `props: AnimatedTabsProps`, destructure on first line followed by a blank line, replace inline event handlers (`handleTabClick`, `handleTabFocus`, `handleTabKeyDown`), and ensure blank lines before returns.
  - `basic-accordion.tsx`: Take `props: BasicAccordionProps`, destructure on first line followed by a blank line, replace inline event handler with `handleTriggerClick`.
  - `mac-terminal.tsx`: Replace arrow functions (`createMachine`, `startNextLine`, `tick`, `updateVisibility`) with function declarations, and insert blank lines before returns.
  - `text-morph.tsx`: Extract `segmentGraphemes` into `segment-graphemes.ts`, `tokenizeText` into `tokenize-text.ts`, and `buildTokens` into `build-tokens.ts`.
- **Test Names (10.4):**
  - Update all `describe(...)` descriptions in `packages/web` test suites to begin with `web/ui/<name>/<file>` or `web/vendors/<name>/<file>`.
