# The session panel's primitives move to web

**Status:** ready-for-agent

**Requested by:** the owner, 2026-11-05: move the hero's session panel and its progress bars into `packages/web`, so the website and the console build the panel from the same pieces. **Decided in:** grilling, 2026-11-05 (rounds 1 and 2).

## Tickets

- 0337: the reveal value, check mark and card section move to web.
- 0338: the expandable card moves to web.
- 0339: the select field moves to web, blocked by 0337.

## Problem Statement

The home hero's session panel is the stack of four cards beside the terminal: Ticket, Analysis, Plan and Workflow. ADR 0005 and console rule 2 require the console's runtime view to show exactly that panel, with real data. Today the panel is built only inside the website's Microruntimes route. Its cards, the card that expands, the value that reveals itself, the section under a hairline rule, the read mark and the labelled select all live there, next to the simulation's data and atoms.

The console cannot reuse any of them, so it already keeps its own copies in its tickets view: its own read mark, its own unknown value, its own widget section and its own document, episode and component rows. Those copies have drifted from the website's. For example, the console's episode row has phone columns the website's lacks. Every new copy is one more place where the console can stop showing what the simulation shows.

## Solution

`packages/web` takes the panel's generic, presentational pieces as dumb components. They hold no state, no copy, no simulation data and no localization. Each one receives everything it shows through props and composes with the others:

- the card that expands;
- the value that reveals itself;
- the section under a hairline rule;
- the check mark;
- the labelled select.

The progress bar, `ShareBar`, is already there.

The four cards, the stack, the rows and the skeleton of the stack are business composition, and they stay in each app. The website keeps building its panel in the Microruntimes route, now from the web pieces, and looks and behaves exactly as before. The console builds its own panel from the same pieces in its own spec, once the hand-off records carry the fields.

## User Stories

1. As the owner, I want the panel's generic pieces in `packages/web`, so that the website and the console show one panel built from one set of components.
2. As the owner, I want the web pieces to be dumb components that receive props and compose with each other, so that each app keeps its business composition, messages and behavior (React rule 11).
3. As the owner, I want no copy, catalog or translated label in `packages/web`, so that localization rule 1 holds: every label comes from the consuming app's JSON catalog.
4. As the owner, I want the web pieces free of Jotai, next-intl and the simulation's data, so that the console can feed them from the hand-off records without the website's state.
5. As a website visitor, I want the hero's panel to look and behave exactly as it does today, so that moving the pieces changes nothing I can see.
6. As a website visitor, I want every card to open when its act arrives and collapse to its one summary line, so that I follow the session as before.
7. As a website visitor, I want to open a card by hand until the next act takes it back, as before.
8. As a website visitor, I want an unknown value to read as a dash, a resolving one as a pulsing skeleton and a known one to fade in, so that unknown is never shown as zero.
9. As a screen-reader user, I want an unknown value announced as "Not set yet" and a resolving one as its resolving label, as today.
10. As a screen-reader user, I want a collapsed card's body to be inert and hidden from assistive technology, and its header to say whether it is expanded and which region it controls.
11. As a screen-reader user, I want the check mark announced by its label while the ✓ itself stays decorative.
12. As a website visitor, I want the Workflow card's model and effort selects to keep their 16px text on phones, so that focusing one does not zoom the page.
13. As a website visitor, I want a select that a run is reading to keep its focus but refuse changes, as today.
14. As a website visitor on a reduced-motion setting, I want the skeleton pulse, the fade-in and the card's slide to stay off, as today.
15. As the console's developer, I want the card that expands to be controlled from outside through `expanded` and `onToggle`, so that the console decides which card is open, whether one at a time or several, as its free layout needs (console rule 3).
16. As the console's developer, I want the card that expands to keep its body mounted while it is collapsed, so that values keep their state.
17. As the console's developer, I want each web piece exported from its own `web/ui/<name>` entry, as the existing pieces are, so that I import only what I compose.
18. As the console's developer, I want the reveal states `unset`, `resolving` and `set` typed by the web piece that draws them, so that my hand-off mapper can say "not known" as `unset` and never as zero.
19. As the console's developer, I want the labelled select typed by its option values, so that a model or effort choice stays type-safe when I wire it.
20. As the owner, I want the console's existing copies listed in this spec, so that the console's own spec replaces them with the web pieces.
21. As a reviewer, I want the website's existing session tests to pass with no changed assertion, so that "exactly as before" is proven rather than claimed.
22. As a reviewer, I want each web piece tested on its own states, tones, labels and callbacks, so that its contract holds for both apps.
23. As the owner, I want the pieces named generically, without FinOps vocabulary, so that they read as UI primitives and not as the website's business.

## Implementation Decisions

- **Boundary.** `packages/web` holds only generic presentational pieces. The four cards (Ticket, Analysis, Plan, Workflow), the stack that opens one card and collapses the others, the rows (document, episode, component), the stack's skeleton, the reveal timing, the summaries, the status line, the formatter and the atoms stay in the website's Microruntimes route. The terminal is out of this work: `MacTerminal` is already a web piece and does not change.
- **The pieces.** Each piece gets its own folder under the package's `ui/`, its own barrel and its own export entry, like the existing pieces:
  - **`ExpandableCard`** comes from the website's session card. It is an article with a header button carrying a label (through `Eyebrow`), a one-line summary and a caret that turns when expanded, over a body region that slides open and shut on its grid rows. It is controlled: it takes `label`, `summary`, `body`, `expanded` and an `onToggle` click handler. An optional `value` lands on the button as a data attribute, so one handler can serve several cards; the website uses the plan phase. Shut, the body stays mounted, inert and `aria-hidden`. The header controls the region and the region is labelled by the header, with ids from `useId`.
  - **`RevealValue`** comes from the website's card value. It takes `state`, `labels`, optional `children`, `className` and `width`. `unset` draws a decorative dash and the `notSet` label for screen readers. `resolving` draws a pulsing bar of the given width (default `w-24`), with no pulse under reduced motion, and the `resolving` label for screen readers. `set` fades the children in. The piece owns the type of its states (`'unset' | 'resolving' | 'set'`) and the interface of its labels (`notSet`, `resolving`). The website's `revealState` compute stays in the website and returns the web type.
  - **`CardSection`** is a block under a hairline rule, `mt-5 border-t border-border pt-4`, merged with a consumer `className` through `cn()`.
  - **`CheckMark`** comes from the website's read mark: a ✓ in the success role, decorative, with its `label` for screen readers.
  - **`SelectField`** comes from the website's workflow choice. It is a native select with a label (through `Eyebrow`), a caret, an optional note and a `locked` state that keeps focus and marks the select `aria-disabled`. Refusing the change stays with the consumer, as the website's choice atoms refuse it while a run plays. It is generic over its option value type. `onChange` receives the chosen option's value. Its text is 16px on phones and 12px from `sm`.
  - **`ShareBar`** stays as it is.
- **No visual change.** Each piece keeps the exact class strings, attributes, ARIA and motion of the website component it replaces. Only names and comments change: the comments drop the "ticket card" wording and describe the piece generically.
- **Composition stays local.** The website's cards, rows and stack import the web pieces in place of their local copies, and the local copies are deleted. The catalog, `microruntimes.json`, does not change: the website passes the same labels from the same keys. The website's session card file disappears, and the stack composes `ExpandableCard` directly.
- **Prop shape.** The web pieces take display values: localized text, shares from 0 to 1 and reveal states. They never take raw units or hand-off records. The hand-off records are still empty sections until specs 0050–0066 fill them. Turning a record into display values is each app's mapper; the console's belongs to the console's spec.
- **Dependencies.** The pieces use only what `packages/web` already depends on: React, `cn()`, `Eyebrow` and the arrow icons. No dependency is added to `packages/web` or to the apps.

## Testing Decisions

- A good test drives a component the way its consumer does, through props and the rendered DOM, and asserts what a person or assistive technology perceives: text, roles, ARIA state, the data attribute, the callback's argument. It does not assert class lists, except where a class is the behavior (the share bar's width, the 16px phone text).
- **Seam 1, the existing one, proves "exactly as before".** The website's session cards test and session terminal test, through their existing mount mocks, keep every assertion unchanged. Only their imports may change. This is the highest seam: it plays the whole session through the panel.
- **Seam 2, one per new web piece.** Each piece gets a unit test beside it, described by its path (`web/ui/expandable-card/expandable-card`, and so on), covering:
  - `ExpandableCard`: the expanded and collapsed ARIA, the inert body, the label and summary, and `onToggle` receiving the click with the `value` attribute;
  - `RevealValue`: the three states and their screen-reader labels, and the resolving width;
  - `CardSection`: the merged class;
  - `CheckMark`: the decorative ✓ and the label;
  - `SelectField`: the options, the value, `onChange` with the typed value, `aria-disabled` while locked, and the note.
- **Prior art.** The web unit tests already beside `share-bar`, `eyebrow`, `skeleton`, `detail-list` and `card` (`*.unit.tsx`, jsdom, Testing Library or `renderToString`). On the website, the session cards and session terminal tests with their `testdata/` mounts.

## Out of Scope

- The console. Building the console's panel, its mapper from the hand-off records and its runtime view is its own spec, written when its turn comes. The console's existing copies stay as they are until then. They are, in `apps/console`'s tickets view: the read mark, the unknown value, the widget section, the document row, the episode row and the component row. That spec replaces them with the web pieces.
- The terminal (`SessionTerminal`, `MacTerminal`).
- Generic rows. Document, episode and component rows stay in each app. A shared row waits for a second use that is really the same.
- The stack's skeleton, the play button's style and the large figures' typography.
- Any change to the simulation's behavior, timing, data or catalog.
- `lib/cli` and `lib/core`.

## Further Notes

- ADR 0005 holds: the simulation's behavior does not change. This spec moves code and changes no step, line or field.
- This spec runs outside the 0048–0066 order, because the owner asked for it on 2026-11-05. It touches only `apps/website` and `packages/web`, while another session changes `lib/cli` for spec 0048.
- The owner settled the boundary in grilling: "the web components have to be dumb components, which compose with each other and receive props."
