# Client state complies with the state rules

**Status:** ready-for-agent

## Tickets

- [ ] 0292: the Microruntimes session is fine-grained atoms
- [ ] 0293: the early-access form updates its state without Immer
- [ ] 0294: the book-a-demo form keeps one responsibility per file
- [ ] 0295: the console's runtime form keeps one responsibility per file
- [ ] 0296: the product creation and folder picker forms keep one responsibility per file
- [ ] 0297: the login and signup forms keep one responsibility per file
- [ ] 0298: drift measures the state rules

## Problem Statement

`.claude/rules/07-jotai.rule.md` sets how client state, Jotai and form logic look in the apps. An audit on 2026-11-04 measured the rule at 82% applied, and `drift quality check jotai` reports nothing, because its two detectors only look at the atoms folder. What the audit found:

- **One coarse atom (rule 3).** The website's Microruntimes hero keeps its simulated session in one six-field object atom (`apps/website/src/app/atoms/microruntime-session.atom.ts`), replaced whole on every write. The session control reads only `paused` and the terminal only the count and the stage, yet every consumer subscribes to the whole object and re-renders on every terminal tick. The derived and write atoms are fine-grained; the primitive is not.
- **Atoms imported past the barrel (rule 3).** Two computes in the Microruntimes components (`play-label.compute.ts`, `run-outcome.compute.ts`) import the atom module through `@app/atoms/…` instead of `@state`, so the barrel is not the one entry to shared state.
- **Immer without drafts to help (rule 5).** The early-access form updates a four-field state with one nested object through `produce`; every map sets one to three fields that a spread expresses as plainly. The website carries `immer` for that alone.
- **An inline default (rule 7).** The console's runtime form seeds `useState` with a five-field object literal built in the component (`apps/console/src/app/lib/views/runtimes/runtime-form.tsx`), not from a `*-state.data` file.
- **Forms not split by responsibility (rule 8).** One client form in seven complies. The website's book-a-demo form has its types, defaults and schema in their files (ticket 0160 stopped there) but keeps its submit and its toast dismissal inline; the early-access form is split except for the focus event, whose logic stays in the component. In the console, none of the runtime form, the product creation form, the folder picker, the login form and the signup form has a folder, a handler file, a map file, a types file, a schema or a defaults file: the runtime form alone holds seven inline handlers and builds its request body inside the submit. Spec 0026 listed the state rules in its scope for the console and decided nothing about its forms; the rules outrank that spec (General, Authority 2), so the forms are fixed here.

Until this is remediated, the website's hero re-renders more than it reads, two files reach past the `@state` barrel, and five console forms cannot be tested without rendering them.

## Solution

Every app applies the state rules, with no visible change:

- The Microruntimes session becomes one concern of fine-grained atoms: one primitive per value the session holds, the derived atoms as today, and the write atoms that coordinate several values staying with them. Each component subscribes to the atoms it reads. Nothing but the `@state` barrel imports an atom module.
- The early-access form updates its state with plain functions that return the next state; `immer` leaves the website.
- Every client form has its folder with one file per responsibility: one handler per user event, one map per state mutation, its types, its Zod schema and its defaults. The handlers take their inputs and collaborators as parameters, so each app's test runner covers them without rendering.
- `drift quality check jotai` gains a detector for an atom module imported outside the barrel and one for a `useState` seeded with an object literal, and reports nothing.

## User Stories

1. As the owner, I want `drift quality check jotai` to report zero findings, so that the state rules are measured rather than asserted.
2. As a visitor of the Microruntimes hero, I want the session to play, pause, replay and run exactly as today, so that the refactor costs me nothing.
3. As a developer, I want the session's stage, count, pause, run counter, opened card and launched choices each in its own atom, so that a component subscribes to what it reads.
4. As a developer, I want the writes that move the session between acts to stay in the session's atom module, so that coordinated state is coordinated in one place.
5. As a developer, I want every shared atom imported from `@state`, so that the barrel is the one entry to shared state.
6. As a developer, I want the Microruntimes computes to take the session values they read, so that they do not depend on an atom module.
7. As a developer, I want the early-access form's maps to return the next state, so that the state rules' Immer clause holds and the website carries no `immer`.
8. As a visitor, I want the early-access form to validate, submit and report as today, so that dropping Immer changes nothing I see.
9. As a developer, I want the book-a-demo form's submit in its own handler file, so that posting the demo request is tested without rendering the form.
10. As a visitor, I want a demo request to reach the edge, to show the return-soon notice on 503 and the failure toast otherwise, exactly as today.
11. As a developer, I want the console's runtime form to seed its state from a defaults file built from the runtime being edited, so that defaults are data.
12. As a developer, I want the runtime form's configuration changes, component toggles, applied recommendations and submit each in its own handler, with the request body built by a map, so that each can be tested with Node's test runner.
13. As a console user, I want creating and editing a microruntime to send the same request as today, so that console-stubs and `ex` see no difference.
14. As a developer, I want the product creation form's folder selection, drop, browse and submit each in a handler, with the request built by a map, so that registering a project from a folder is testable logic.
15. As a developer, I want the folder picker's path submission in a handler and its listing request in a map, so that the picker's dialog holds only presentation.
16. As a developer, I want the login and signup forms to post their credentials through a handler with a schema, so that the two forms read alike and their logic is testable.
17. As a developer, I want every console form's schema to mirror the constraints its inputs already declare, so that the browser's validation and the form's agree.
18. As a developer, I want a form's handler to receive the transport and the navigation it needs as parameters, so that a test passes fakes without mocking modules.
19. As a developer, I want every form folder to carry only the files for the responsibilities the form has, so that no empty file exists to satisfy a list.
20. As a developer, I want the existing form, session and card tests to keep passing, moved with the code they cover, so that behaviour is proven unchanged.
21. As the owner, I want each strict reading of the rule recorded, so that I can overturn it at review.

## Implementation Decisions

- **The session concern.** The Microruntimes session module keeps one concern and gains one primitive atom per value: the stage, the count of played lines, whether the visitor paused, the run counter, the card opened by hand and the choices the playing run was launched with. The workflow choices atom stays. The derived atoms (run script, plan and run counts, session ended, run in progress, can play, open card) keep their meaning and read the primitives they need. The write atoms (record played count, open a card by hand, choose model, effort and human stages, play the workflow, start the run, replay, toggle pause) keep their names and set the primitives they coordinate, in the same module. The session interface and its initial object go; each primitive carries its own initial value. The `@state` barrel re-exports the new names.
- **Consumers read what they use.** The terminal, the cards, the workflow card and the session control subscribe to the atoms they read. The play-label and run-outcome computes take the stage, the launched choices and the count as values; their tests pass those values. No module outside the barrel imports an atom module; the computes that need an atom's value type take it from the data module that owns it or declare it locally.
- **Early-access without Immer.** The form's state updater takes a function from the current state to the next state, as React's functional `setState` does. Each map takes the state and returns the next one with a spread; the handlers and the submit handler call the maps through the updater. The handler tests apply the updater to a local state. `immer` leaves the website manifest and the dead-code allowlist is untouched.
- **Form folders.** Each client form lives in a folder named after it, beside the view that composes it: the website's `book-a-demo` and `early-access` folders as today; in the console, a `runtime-form` folder under the runtimes views, a `product-create` folder and a `folder-picker` folder under the products views, and a `login` and a `signup` folder under the auth views, with the auth shell staying beside them. The component file keeps its name. Beside it sit `form.types.ts`, `form.validation.ts`, `form-state.data.ts`, one `<event>.handler.ts` per user event and one `<mutation>.map.ts` per state mutation. A form gets a file only for a responsibility it has: a form with no defaults has no defaults file, and a form whose state is a primitive has no map.
- **Handlers and maps.** A handler is the logic of one user event: it takes the event's values, the form's updater or setters, and the collaborators it needs (the console's `mutation` transport, a navigation callback, an abort signal) as parameters. A map is pure: it takes the state and the event's values and returns the next state, or takes the form's values and returns the request. The component's named functions adapt the DOM event to the handler and hold nothing else.
- **Book-a-demo.** The submit handler posts the validated data with the locale to the edge and returns an outcome, `sent`, `unavailable` or `failed`; a map turns the response into that outcome; the component navigates, shows the notice or shows the toast on the outcome. The toast dismissal gets its handler. The `useState` primitives for the outcome and the toast message stay local, as rule 1 and rule 5 allow.
- **Early-access focus.** The focus event's logic, preparing the challenge, moves to its handler file beside the blur and change handlers.
- **Runtime form.** Its state, the selected component paths and the configuration fields as strings, has a types file and a defaults file whose function builds the defaults from the runtime being edited, or empty for a new one. Maps set one configuration field, apply a recommendation's values, and toggle a component. A request map builds the mutation from the state, the submitted name and native runtime id, the selected components and the runtime being edited: `PATCH` with the revision for an existing runtime, `POST` with the product id for a new one, the optional budget and native runtime id only when set. The schema validates the submitted values, the component fields included, replacing the ad-hoc parse of the form data; an invalid build keeps today's alert. Handlers cover each field change, the effort change, the component toggle, the applied recommendation and the submit.
- **Product creation.** Its state, the path, the name and whether a folder was dropped, has a types file and a defaults file. Maps derive the name from a selected path, record a dropped file, and build the request: the folder endpoint of an existing product or the products endpoint, with the trimmed name or none. Handlers cover folder selection, drop, browse, path and name change, picker close and submit. The schema requires a non-empty path and bounds the name.
- **Folder picker.** A handler turns the submitted path into the listing request through a map that builds the request from the path and the abort signal. The input value stays a local primitive.
- **Login and signup.** Each has a types file for its credentials, a schema that mirrors the inputs' constraints (signup's twelve-character minimum included) and a submit handler that posts the credentials through the transport and navigates home through the callback it receives. The console adds `zod` for its schemas.
- **Tests.** Handler and map tests sit beside their files: Vitest in the website, Node's test runner in the console. The session concern gets its own unit test beside the atom module, which the atoms detector already allows. Component tests are moved with the code they cover and read the new atoms.
- **Detectors.** `drift quality check jotai` gains: an atom module imported outside the `@state` barrel (a line pattern on `@app/atoms/` imports, excluding the barrel); a `useState` seeded with an object literal (a line pattern on `useState(` or `useState<…>(` followed by `{`). Whether Immer earns its place, whether a form folder is complete and whether a handler or a map holds one responsibility stay review-enforced.
- **Statement spacing, naming and suffixes** follow the general and TypeScript rules in every file this spec touches.

## Testing Decisions

- A good test asserts what a form or a session does, not how: the request a handler sends and the outcome it returns, the state a map produces, the values a derived atom yields after a write atom runs, the markup a component renders.
- **The existing suites are the behaviour seam.** `pnpm turbo run check typecheck test build --filter=website` and `--filter=console` pass after every ticket. The website's session-cards and session-terminal tests and the two form tests keep passing against the new atoms and handlers.
- **Handler and map tests are the new seam**, one file per handler or map that holds logic: the website's with Vitest, as the early-access handler test does; the console's with Node's test runner, as the ticket-phases test does, passing fakes as parameters.
- **The session atoms get a unit test** that drives the write atoms through a Jotai store and reads the derived atoms, the seam the component tests use today.
- **The detectors are the rules seam.** `platform/cli/drift quality check jotai` reports nothing at the end, and `drift quality coverage jotai` stays at 100%.
- Prior art: `form.handlers.unit.ts` and `session-cards.unit.tsx` in the website; `ticket-phases.test.ts` in the console.

## Out of Scope

- The dissolution of `src/app/lib/` in both apps, which spec 0041 owns; every app ticket here waits on it and takes the folders where 0041 places them.
- The React rules inside the forms (effects, scope order, loading flags), which spec 0040 owns, and the translator names and file suffixes, which spec 0039 owns.
- The SmoothUI class on the book-a-demo button (spec 0038) and the recommendation panel's `reading` flag (spec 0040).
- console-stubs and `ex`, whose request contracts do not change.
- Any new shared state: the console has no Jotai today and gains none.

## Further Notes

- **Immer is removed (strict reading of rule 5).** The early-access state is four fields with one nested object; the maps set one to three fields. A spread expresses that as plainly, so `produce` does not earn its place. The owner can overturn this at review and keep `immer`.
- **Console forms get Zod schemas (strict reading of rule 8).** Rule 8 names `form.validation.ts` for a form's Zod schema. The console's forms validate through HTML attributes today; each gains a schema that mirrors those constraints, so the browser's validation and the form's agree and the handlers have a testable contract. The console adds `zod`. The owner can overturn this and keep HTML validation alone.
- **One handler file per user event, however small.** The toast dismissal, the picker close and a field change each get a handler file, as rule 8 reads. The component's named function stays as the adapter from the DOM event.
- **A form gets a file only for a responsibility it has.** Rule 8 lists where each responsibility lives; it does not require an empty file. The owner can ask for every file in every folder instead.
- **The component file keeps its name** (`runtime-form.tsx`, `login-view.tsx`); rule 8 names the logic files, not the component.
- **Nothing here changes `lib/`.** Spec 0041 moves it; the folders named above are relative to the views and components wherever 0041 places them.
