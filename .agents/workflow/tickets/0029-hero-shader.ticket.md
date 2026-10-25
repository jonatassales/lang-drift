# 0029: HeroShader component

**Spec:** .agents/workflow/specs/0004-platform-landing-hero.spec.md

**What to build:** `@langdrift/web` exports `./ui/hero-shader`: a client component that fills its parent with a WebGL2 shader adapted from Exphere's `exphere-shader`, colored from the token roles at runtime with extra purple shades and a light theme, recolored when `.dark` toggles on `<html>`, still under reduced motion, and hidden behind a token-only gradient when WebGL2 is unavailable.

**Blocked by:** none.

**Status:** ready-for-agent

### Files

- [ ] `packages/web/src/ui/hero-shader/shaders.ts`, `palette.ts`, `surface.ts`, `hero-shader.tsx`, `index.ts` and `hero-shader.unit.tsx`; `packages/web/package.json` exports `./ui/hero-shader`.

### Behavior

- [ ] Props: `className?`, `paused?`, `reducedMotion?`. Markup: an `aria-hidden` wrapper, `pointer-events-none absolute inset-0`, with a gradient from `--primary`, `--primary-hover`, `--background` and `--card`, holding a `<canvas>` hidden by `data-unavailable`.
- [ ] Colors: `--background`, `--card`, `--foreground`, `--primary` and `--primary-hover` read from the canvas's computed style and normalized through a 2D canvas; `shadeDeep`, `shadeSoft`, `glow` and `violet` derived in JS; `scheme` from `<html class="dark">`.
- [ ] A `MutationObserver` on `<html class>` recolors and redraws.
- [ ] Reduced motion, `paused` and `reducedMotion` give one static frame; hidden tabs and off-screen canvases stop the loop; DPR is capped at 2 and pixels at about 1920×1080; context loss sets `data-unavailable`; dispose deletes the program and detaches every listener and observer.
- [ ] No CSS file, no `ex:` prefix, no hex literal, no `../` import.

### Tests

- [ ] SSR markup and class merge; `mix`, `rotateHue` and `derivePalette` on literal colors.

### Verification

- `pnpm turbo run typecheck test --filter=@langdrift/web`
- `pnpm turbo run check`
