# Frontend rules

Rules for React and Next.js code in this repository. Loaded in every session.

1. **Style components with Tailwind CSS.** Do not write plain CSS for components: no hand-written stylesheets, CSS Modules or CSS-in-JS.

2. **CSS files only configure Tailwind.** An app's CSS holds its Tailwind entry point, the import of its product's tokens file from `@langdrift/tokens` and Tailwind directives, nothing else.

3. **`cn()` composes conditional classes.** It is `clsx` plus `tailwind-merge`, exported by `@langdrift/web` (`@pack/react/utilities` in the website). A static class string needs no wrapper.

4. **shadcn and SmoothUI components are owned source.** Their classes change in place; nothing wraps them in a CSS abstraction.

5. **Temporary exception (ADR 0003).** The website arrived from langdrift with hand-written CSS, and `packages/web` inherited some. Only the files below are allowed, as they are. They may only shrink: no new `.css` file, selector, rule, hover or focus state, media query or layout rule is added anywhere, and a change that touches one of these files moves the touched styling to Tailwind when that is local and safe. The exception ends when the last file is gone.

   - `apps/website/src/app/[locale]/(pages)/early-access-sent/early-access-sent.css`
   - `apps/website/src/app/globals.css`
   - `apps/website/src/app/hero.css`
   - `apps/website/src/app/lib/components/conversion.css`
   - `apps/website/src/app/lib/components/early-access/early-access.css`
   - `apps/website/src/app/lib/components/executive-review.css`
   - `apps/website/src/app/lib/components/hero-demo.css`
   - `apps/website/src/app/lib/components/integrations-section.css`
   - `apps/website/src/app/lib/components/vision-loop-evidence.css`
   - `apps/website/src/app/lib/components/voice-preview/voice-preview.css`
   - `apps/website/src/app/lib/components/website-footer.css`
   - `apps/website/src/app/lib/components/website-header.css`
   - `apps/website/src/app/responsive.css`
   - `packages/web/src/styles.css`
   - `packages/web/src/ui/product-experience.css`
   - `packages/web/src/ui/product-vision-curve/product-vision-curve.css`
   - `packages/web/src/ui/ui-consistency.css`
   - `packages/web/src/vendors/smoothui/animated-tabs.css`
   - `packages/web/src/vendors/smoothui/basic-accordion.css`
   - `packages/web/src/vendors/smoothui/smoothui.css`

