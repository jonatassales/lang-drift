# Jotai rules

Rules for shared state in the apps that chose Jotai. Loaded in every session.

1. **Local first.** Ephemeral interaction state stays in `useState` while it does not need to leave the component. Atoms hold state that is shared, independently consumed, coordinated or observed outside that boundary.

2. **Atoms by concern.** One atom per `src/app/atoms/<name>.atom.ts`, re-exported by name from the `@state` barrel. A concern's primitive, derived and write atoms stay together; unrelated concerns are not spread across technical buckets.

3. **No `'use client'` in an atom module.** The client boundary belongs to the consuming UI.

4. **Immer only when drafts help.** Ordinary atoms for primitive values; `produce` for complex object or array updates only. Complex state does not require Immer by itself.

5. **No automatic logging.** Shared state stays inspectable through a diagnostic that fits the app; a central atom registry or a logger of every transition is not one.
