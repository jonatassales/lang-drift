# Frontend rules

Rules for React and Next.js code in this repository. Loaded in every session.

1. **Style components with Tailwind CSS.** Do not write plain CSS for components: no hand-written stylesheets, CSS Modules or CSS-in-JS.

2. **CSS files only configure Tailwind.** An app's CSS holds its Tailwind entry point, `@theme` tokens (`--lg-*`) and Tailwind directives, nothing else.
