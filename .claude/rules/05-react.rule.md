# React rules

Rules for React components in apps and in `packages/web`. Loaded in every session.

1. **Server Components by default.** `'use client'` appears only for browser APIs, interactivity, client state or a client-only library, at the narrowest boundary that needs it.

2. **One component per file, one hook per `.hook.ts` file.** Related files sit together by concern. This is a component and hook boundary, not a one-function-per-file rule.

3. **Props are a named interface.** A component with props declares `interface <Name>Props`, receives a `props` argument and destructures it at the start of the body, followed by a blank line. An anonymous object annotation, a props type alias or a destructured parameter list does not express this contract. The interface stays with its component unless it has a real shared owner.

4. **Transitions, not loading flags.** `useTransition` or `startTransition` for a non-urgent update that must keep the interface responsive. A `useState` loading flag is for a real application state, not for a transition React already models with `isPending`.

5. **Scoped `<Suspense>`.** It wraps an asynchronous section that can render on its own: not the whole page for one section, and not every component. A route that can suspend during navigation gets a `loading.tsx`; a segment that needs a runtime error boundary gets an `error.tsx`, which handles render failures while `try/catch` handles an imperative external call. `template.tsx` exists only when a route must remount on navigation; an effect or a key hack never imitates it.

6. **Routes compose pages.** Independently meaningful UI is composed at its owning page, layout or feature. No central component selects whole pages through a mode or section prop. Configurable components, render props and dependency injection stay allowed.

7. **Shared primitives live in `packages/web`.** Generic UI, the icon family and their dependencies have one owner there, so every app uses one icon identity; apps keep business composition. A barrel exports its symbols by name, never with `export *`.

8. **Domain types stay with their owner** and are imported with `import type`. Props, and types derived only from props and used only in that file, stay local.
