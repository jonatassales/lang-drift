# State rules

Client state, Jotai and form logic in the apps. Loaded in every session.

1. **Local state first.** Ephemeral interaction state stays in local React state while it does not need to leave the component or its local composition boundary. Atoms hold state that must be shared, independently consumed, coordinated or observed outside that boundary.
2. **Shared client state is Jotai.** Shared client state uses Jotai atoms, never a React Context provider.
3. **Atoms by concern.** One concern per `src/app/atoms/<concern>.atom.ts`, re-exported by name from the `@state` barrel. A concern's primitive, derived and write atoms stay together; unrelated concerns are not spread across technical buckets. State units stay fine-grained and independently consumable, and coordinated writes stay with the values they coordinate.
4. **No `'use client'` in an atom module.** The client boundary belongs to the consuming UI.
5. **Immer only when drafts help.** Ordinary atoms and `useState` for primitive values; Immer (`produce`) only for complex object or array updates that benefit from drafts. Complex state does not require Immer by itself.
6. **No automatic logging.** Shared state stays inspectable through a diagnostic that fits the application. A central atom registry, a state logger or console logging of every transition is not one and is never added as an automatic compliance measure.
7. **Component state defaults.** Component-level defaults live in `*-state.data` files and may compose from the route-level `state.data`.
8. **Form logic by responsibility.** A client form's logic lives in its folder, one responsibility per file: one `<event>.handler.ts` per user event, one `<mutation>.map.ts` per state mutation, `form.types.ts` for its types, `form.validation.ts` for its Zod schema and `form-state.data.ts` for its defaults.

```ts
export const loginFormState: LoginFormState & LoginState = {
	...loginStateData,
	isPending: false,
}
```
