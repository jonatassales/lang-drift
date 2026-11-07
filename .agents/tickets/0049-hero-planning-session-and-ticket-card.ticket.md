# 0049: The hero opens with the heading and a planning session that reveals the ticket

**Spec:** .agents/specs/0007-microruntimes-hero.spec.md

**What to build:** the Microruntimes page opens with an `h1` that reads only "Microruntimes" and, under it, the terminal on the left and the ticket card on the right. The terminal types `ex agent plan` with the ticket number on its own; on Enter the card fades in with the ticket number large, the PRD, design doc and spec names and dashes everywhere else. As the session prints, the suggested model and effort go from dash to skeleton to value, the artifacts gain a read mark and the guardrail fills in with its note. The card's components, hardware and report regions arrive in tickets 0050 and 0051; the session already prints every line.

**Blocked by:** 0046, 0048.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

### Opening

- [ ] The page composes its own hero section labelled by the `h1` "Microruntimes", with the current heading classes and no gradient highlight, a screen-reader-only summary wired by `aria-describedby`, and the existing two-column grid; `ProductHero` is no longer rendered on this page and is not edited. The skip link keeps its target.
- [ ] The summary describes the final plan, names the ticket, the artifacts, the suggestion, the five components and the bounded and suggested figures in units, and says the session is simulated and nothing was built or paid for, in all five locales.

### Session

- [ ] The content module replaces the old session with the nineteen-line plan session (ids, kinds, delays): the typed command, the plan-only comment, the console's RAG line, the episodes line, the suggestion, the two `reads` lines, the bounded line, the resolve line, the five component lines, the size line, the limits warning, the estimate line, the plan line and the closing comment. The catalog carries each line's text in five locales; `ex`, `runtime.md`, the model placeholder, file names, paths, the guardrail criteria and the artifact numbers stay untranslated.
- [ ] `SessionTerminal` (client) maps the session to catalog text and renders `MacTerminal` keyed on the run counter, with the height clamped to the viewport so the ticket head fits in a phone's first screen, binding `onProgress` and `paused` to the session atom; the page passes no function.
- [ ] `microruntimeSessionAtom` in its own atom module, shaped `{ playedCount, paused, run }`, re-exported by name from the state barrel.

### Ticket card

- [ ] `TicketCard` (client) is in the DOM from the first paint, invisible with zero opacity so the grid row reserves its height, and becomes visible with a motion-safe opacity transition once the typed command completes.
- [ ] The head: the ticket eyebrow, the number as the only large type, the slug, the title and a `ready-for-agent` chip; a read mark appears after the `reads` line.
- [ ] The artifacts as a two-column definition list: PRD, design doc and spec with their names from the first frame and a read mark after the `reads` line; the guardrail from a dash to its file, criteria and the note that it is a review rule with no automatic enforcement.
- [ ] The suggested model and effort cells: a dash, a pulsing skeleton between the RAG line and the suggestion, then the model placeholder and `medium` with the note that the recommendation is the team's to authorize. Skeletons do not pulse under reduced motion.
- [ ] Every field derives from `playedCount` and a reveal table in the content module; the card owns no timer. Under reduced motion the card renders its final state at once.
- [ ] The card's labels, the ticket title and the notes live in the catalog's new session family in five locales, in the same key order as `en`.

### Tests

- [ ] The content module test: exactly one typed command; every reveal points to a real line; delays are positive; the session lasts between fifteen and twenty seconds at the default typing speed.
- [ ] A jsdom test for `TicketCard` with a fresh Jotai store: at count zero the card is invisible and shows the artifact names and dashes; after the RAG line the model and effort cells show skeletons; after the suggestion they show the placeholder and `medium`; after the guardrail line the guardrail row shows its criteria.
- [ ] A jsdom test for `SessionTerminal` with the stubs and fake timers of ticket 0046: the store's count increments as lines complete.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/microruntimes/` and `/pt-br/platform/microruntimes/`: the heading alone, the terminal typing the command, the card fading in on Enter and the model and effort filling; 375px without horizontal overflow; reduced motion shows the full transcript and the filled card at once.
