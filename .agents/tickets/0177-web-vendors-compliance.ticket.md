# 0177: Vendor components rules compliance in packages/web

**Spec:** .agents/specs/0030-web-package-rules-compliance.spec.md

**What to build:**
Refactor vendor components under `packages/web/src/vendors/smoothui/`:
1. `animated-tabs.tsx`:
   - Change signature to `AnimatedTabs(props: AnimatedTabsProps)` with destructuring on the first line followed by a blank line.
   - Replace inline event handlers (`onClick`, `onFocus`, `onKeyDown`) with named functions (`handleTabClick`, `handleTabFocus`, `handleTabKeyDown`).
   - Add blank lines before returns inside `tabs.map`.
2. `basic-accordion.tsx`:
   - Change signature to `BasicAccordion(props: BasicAccordionProps)` with destructuring on first line followed by a blank line.
   - Replace inline `onClick` with named function `handleTriggerClick`.
3. `mac-terminal.tsx`:
   - Replace arrow functions (`createMachine`, `startNextLine`, `tick`, `updateVisibility`) with function declarations.
   - Extract `createTerminalMachine` into its own file.
   - Add blank lines before returns across the file.
4. `text-morph.tsx`:
   - Split helper functions into single-function files: `segment-graphemes.ts`, `tokenize-text.ts`, `build-tokens.ts`.
   - Colocate them under `packages/web/src/vendors/smoothui/`.

**Blocked by:** 0176

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `packages/web/src/vendors/smoothui/animated-tabs.tsx` complies with 00-general and 05-react.
- [x] `packages/web/src/vendors/smoothui/basic-accordion.tsx` complies with 00-general and 05-react.
- [x] `packages/web/src/vendors/smoothui/mac-terminal.tsx` complies with 00-general and 05-react.
- [x] `packages/web/src/vendors/smoothui/text-morph.tsx` has one function, delegating helpers to separate single-function files.
- [x] Quality check passes for `packages/web/src/vendors/smoothui/`.

### Verification

- `pnpm turbo run check typecheck test build --filter=@langdrift/web`
- `platform/cli/drift quality check`
