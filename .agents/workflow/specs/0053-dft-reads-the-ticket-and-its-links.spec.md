# `dft` reads the ticket and its links

**Status:** ready-for-agent

**Decided by:** ADR 0005 (Perfect Simulation Alignment). **Evidence:** `research/simulation-conformance.research.md`, element E6.

> **Main constraint (ADR 0005).** This element does exactly what the home simulation shows: nothing more, nothing less. What it does beyond the simulation is removed; what the simulation shows and it lacks is built; what differs is both.
>
> **Obligation.** Whoever writes this spec's tickets, or implements them, first plays the simulation in production at https://langdrift.md/en/. They watch both phases: `dft agent plan` plays on load, and `dft agent run` plays after Play workflow. They change the Workflow choices (checkpoints, model, effort, a reservation over the budget) and expand every card. They compare what it shows with this spec, line by line and field by field. Where the two differ, the simulation in production wins, and the difference is reported.
>
> **Execution order.** The specs run one at a time, in order: 0049, 0048, 0050, 0065, 0053, 0054, 0056, 0052, 0051, 0066, 0055, 0057, 0058, 0059, 0060, 0061, 0063, 0062, 0064. This spec is step 5 of 19, after spec 0065 and before spec 0054. Its tickets are not written in advance: when an agent starts this spec, it plays the simulation as the obligation above requires, then writes this spec's tickets with `to-tickets`, then implements them. The next spec starts only when this one's tickets are done.

## Tickets

## Problem Statement

The simulation opens `dft agent plan 0031-TICKET-book-a-demo-runtime` by reading the ticket, naming its four links by short id, and reading each linked document for one exact fact: the PRD's use case, the design doc's envelope and budget, the spec's component and acceptance-criteria counts. It reads nothing outside the ticket's links. The LangDrift CLI does something else. It takes a file path instead of a ticket id, follows links breadth-first to depth 3 under file and byte limits, prints full file names without short ids or the "guardrail" label, and summarizes each document with the judge's free text of up to eight words, or "read". Nothing deterministic extracts the use case, the envelope, the budget, "5 components named" or "3 acceptance criteria", so the line a team reads to check what the plan understood is a model's paraphrase, and the Ticket card and the Analysis card's documents have no data behind them.

## Solution

`dft agent plan <ticket id>` finds the one ticket document of that id in the client's project, reads its status and title, and reads exactly the documents the ticket links: one PRD, one design doc, one spec and one guardrail. Each link is named by its short id, the guardrail by its label. The PRD, the design doc and the spec are each read for the fields the simulation shows, by fixed rules over field lines and sections, with no model involved. The terminal prints the six lines this element owns, in the simulation's order and text, and the plan's card data carries the whole Ticket card and the Documents header and rows of the PRD, the design doc and the spec. Nothing is followed beyond the ticket's own links.

## Responsibility

This element:
- resolves the ticket id that the CLI (spec 0048) passes to `dft agent plan` to the one ticket document of the project, and reads its number, slug, status and title;
- reads the ticket's links, classifies each linked document by kind (PRD, design doc, spec, guardrail) and names it by its short id;
- reads the PRD, the design doc and the spec, and extracts the use case, the envelope (`cpuMillis`, `memoryMiB`), the budget per run, the component paths with their counts and the acceptance-criteria count;
- counts the documents read for the terminal's `✓ read` line and for the card's Documents header;
- stops the plan safely when the ticket or a link cannot be read as the simulation needs.

It does not:
- read or state the guardrail document's criteria, its enforcement or its `reads` line (spec 0054, E7), although it finds the guardrail link;
- state the envelope as the plan's bound (`bounded …`, spec 0054) or the budget (`budget …`, spec 0056, E9), although it extracts both from the design doc;
- inspect the components the spec names (spec 0066, E19), although it extracts their paths;
- judge, summarize or weigh anything (spec 0051, E4);
- print the comments, the command or the column layout of the lines (spec 0048, E1).

## Produces

This element writes no file of its own under `.langdrift/`. It produces values that the plan records through its neighbours:

- **Ticket** (into the plan's card data, spec 0050): `id` (`0031-TICKET-book-a-demo-runtime`), `number` (`0031`), `slug` (`book-a-demo-runtime`), `status` (`ready-for-agent`), `title` (`Run the book-a-demo request use case as a microruntime`), `path` (project-relative). Feeds the `tool call: read` and `ready-for-agent · links …` lines and the whole Ticket card. The id also feeds the run's header (spec 0050) and the slug the runtime id (spec 0055) and the branch (spec 0063).
- **Links**, in kind order PRD, design doc, spec, guardrail: `kind` (`prd`, `designDoc`, `spec`, `guardrail`), `file` (`0002-PRD-book-a-demo.prd.md`), `shortId` (`0002-PRD`, `0005-DD`, `0011-SPEC`, `0001-finops`), `path` (project-relative). Feeds the links line and the card's LINKED ARTIFACTS.
- **Extracted facts** (into the plan's card data, and handed to the neighbours named):
  - PRD `useCase` (`book-a-demo request`), feeding the PRD read line and the PRD document row;
  - design doc `envelope.cpuMillis` (integer, millicores, `2000`) and `envelope.memoryMiB` (integer, MiB, `3072`), feeding the design doc read line and row, and handed to spec 0054 for the `bounded` line;
  - design doc `budgetPerRunCents` (integer cents, `2500`), feeding the design doc row and handed to spec 0056 for the `budget` line and, before the plan's first line, for the Workflow card's header and spend (see Build);
  - spec `components` (ordered list of project-relative paths, with the port the spec names beside a path when it names one), feeding the count `5` and handed to spec 0066;
  - spec `acceptanceCriteria` (integer count, `3`), feeding the spec read line and row.
- **Counts**: `readCount` (`4`: the ticket, the PRD, the design doc and the spec, the documents read before the `✓ read` line) for the terminal; `documentsRead` (`4`: the four documents the ticket links, PRD, design doc, spec and guardrail, counted from the resolved links) for the card's Documents header and for the Analysis header's first state (spec 0051), both revealed at the `✓ read` line, before the guardrail row is read at the next line, as the simulation shows.

## Shows

Phase 1, `dft agent plan`, in this order. The first two lines follow `# plan only · nothing is built, started or paid for in this step` (spec 0048):

```
console › tool call: read 0031-TICKET-book-a-demo-runtime
console › ready-for-agent · links 0002-PRD · 0005-DD · 0011-SPEC · guardrail 0001-finops
```

The next four follow `# analysis · reading the linked documents, then the project's episodes` (spec 0048):

```
console › tool call: read 0002-PRD · use case: book-a-demo request
console › tool call: read 0005-DD · envelope cpuMillis 2000 · memoryMiB 3072 · set by the team
console › tool call: read 0011-SPEC · 5 components named · 3 acceptance criteria
✓ read       4 documents · nothing outside the ticket's links
```

The line after them is `reads      harness/guardrails/0001-finops.gr.md · FIN-02 · FIN-04` (spec 0054).

Card fields:

- **Ticket card**, whole: the header `0031 · book-a-demo-runtime`; the number `0031`, the slug `book-a-demo-runtime` and the status `ready-for-agent`; the title `Run the book-a-demo request use case as a microruntime`, all present from the moment the panel appears, before the plan's first line, because the plan reads the ticket's header before it prints anything (see Build); LINKED ARTIFACTS, filled when the links line prints: `PRD 0002-PRD-book-a-demo.prd.md`, `DESIGN DOC 0005-DD-book-a-demo-runtime.dd.md`, `SPEC 0011-SPEC-book-a-demo-runtime.spec.md`, `GUARDRAIL 0001-finops.gr.md`.
- **Analysis card, DOCUMENTS header**: `4 documents read`, once the `✓ read` line prints.
- **Analysis card, document rows**, each set by its own read line:
  - `0002-PRD-book-a-demo.prd.md ✓ Read — use case: book-a-demo request`
  - `0005-DD-book-a-demo-runtime.dd.md ✓ Read — envelope 2,000m · 3,072 MiB · budget $25.00 per run · set by the team`
  - `0011-SPEC-book-a-demo-runtime.spec.md ✓ Read — 5 components · 3 acceptance criteria`

The guardrail document row and its FIN line belong to spec 0054. The card renders integers with the reader's grouping (`2,000m`); the data carries plain integers and cents.

## User Stories

1. As a developer, I want to type `dft agent plan 0031-TICKET-book-a-demo-runtime` with the ticket's id, so that I do not have to know where the ticket file sits.
2. As a developer, I want the plan to print `tool call: read` with the ticket's id first, so that I see which ticket it is about to plan.
3. As a developer, I want the plan to print the ticket's status and its links by short id, so that I can check at a glance that it found the right documents.
4. As a team member, I want the guardrail named with the label `guardrail` and its own short id, so that it reads apart from the product documents.
5. As a team member, I want the links listed in the order PRD, design doc, spec, guardrail, so that every plan reads the same way whatever order the ticket wrote them in.
6. As a product owner, I want the PRD read for its use case, so that the plan states which use case it serves.
7. As a team lead, I want the design doc read for the envelope in `cpuMillis` and `memoryMiB` and marked "set by the team", so that the bound the plan respects is the one the team wrote down.
8. As a team lead, I want the design doc read for the budget per run, so that the money the run may spend comes from the team's document.
9. As a developer, I want the spec read for the components it names and for its acceptance criteria, so that the plan's scope matches the spec's.
10. As a reviewer, I want each read to extract the same fact every time from the same document, so that a plan is reproducible and no model paraphrases what the team wrote.
11. As a team lead, I want the plan to read nothing outside the ticket's links, so that the context of a plan is exactly what the ticket points at.
12. As a developer, I want `✓ read 4 documents · nothing outside the ticket's links` to close the reads, so that I know the reading act is over and bounded.
13. As a console user, I want the Ticket card to show the number, slug, status and title as soon as the ticket is read, so that the panel says which ticket the session is planning.
14. As a console user, I want the Ticket card's linked artifacts to fill in when the links line prints, so that the card and the terminal move together.
15. As a console user, I want each document row of the Analysis card to show ✓ Read and the fact its read extracted, so that I can check the plan's understanding document by document.
16. As a console user, I want the design doc row to show the envelope and the budget together, so that the team's two limits read in one place.
17. As a developer, I want a ticket id that matches no ticket, or more than one, to stop the plan with nothing written, so that a plan never runs on a guess.
18. As a developer, I want a ticket whose link cannot be read, or which lacks one of the four documents, to stop the plan with nothing written, so that a plan never runs on missing context.
19. As a security-minded team lead, I want a link that leaves the project, points at a secret or is a symlink to stop the plan, so that the plan never reads what it must not.
20. As a developer, I want a document that lacks a field the plan needs to stop the plan with the document and the field named, so that I know what to fix.
21. As an implementer of the component inspection, I want the spec's component paths handed to me in the spec's order, so that the inspected components are the ones the spec names.
22. As a maintainer of the fixture, I want the documents' fields written in the tracker's field-line convention, so that `dft init`'s templates and the fixture share one shape.

## Implementation Decisions

**Remove**
- Following links breadth-first to depth 3, and the `--max-depth`, `--max-files` and `--max-bytes` limits with their "unresolved beyond the limits" state. Serves `✓ read … nothing outside the ticket's links`.
- The judge's sources question: the per-source `kept`/`discarded` verdict, the confidence, the inferred role and the summary of up to eight words, and the fallback summary "read". Serves the three `tool call: read` document lines and the document rows.
- The judge's facts and stated figures collected from the documents. The facts the plan needs are the extracted fields above.
- The "N refused" variant of the reads line, and reading refused or binary sources into the count. Serves `✓ read       4 documents · …`.
- Full file names in the links line. Serves `ready-for-agent · links 0002-PRD · 0005-DD · 0011-SPEC · guardrail 0001-finops`.
- The "no status" fallback for a ticket without a status line.
- The `--runtime` override of the runtime's id from the ticket. The slug is the runtime's id (spec 0055).
- Wiki links and field-line links as link syntaxes. A ticket links its documents with Markdown links.
- The per-source record fields no shown field uses: media type, bytes, hash, depth, linked-from, link syntax, and the plan's context hash over the sources.

**Build**
- Ticket resolution: the id names the one file called `<id>.ticket.md` in the client's own artifact folder, `.agents/` or `.claude/` at the project root. The ticket links the PRD, the design doc and the spec, which are the client's and sit in that same folder, and LangDrift's guardrail, which sits in `.langdrift/harness/guardrails/` (spec 0054). `.langdrift/` holds only LangDrift's workflow, skills, harness with its guardrails, prices and records (spec 0048). A PRD, design doc or spec link must stay inside the client's artifact folder, and the guardrail link inside `.langdrift/harness/guardrails/`. The number is the id's first four digits. The slug is what follows `<number>-TICKET-`. The title is the document's first heading. The status is its `**Status:**` line. The ticket is resolved and its header (number, slug, status, title) read before the plan's first line, and handed to the plan's card data first (spec 0050); `tool call: read …` then prints. Serves the first two lines and the Ticket card.
- The early budget: before the plan's first line, the design doc link is resolved and its `**Budget per run:**` line read, so that the Workflow card's header and spend exist from the moment the panel appears (spec 0056). The design doc's read line still reveals the envelope and the budget on the Analysis card at its own moment.
- Document kinds by suffix: `.prd.md` PRD, `.dd.md` design doc, `.spec.md` spec, `.gr.md` guardrail. A ticket links exactly one of each.
- Short ids: a name of the form `<number>-<KIND>-<rest>` is shown as `<number>-<KIND>` (`0002-PRD`, `0005-DD`, `0011-SPEC`). Any other name is shown without its suffix (`0001-finops`). The guardrail is preceded by the label `guardrail`. Serves the links line and the three read lines.
- Deterministic extraction, by the tracker's field-line convention (`**Field:** value`) and named sections:
  - the PRD's `**Use case:**` line gives `use case: <value>`;
  - the design doc's `**cpuMillis:**`, `**memoryMiB:**` and `**Budget per run:**` lines (a dollar amount) give the envelope and the budget;
  - the spec's `## Components` section, one list item per component, each starting with the component's path in backticks and optionally followed by a port as `:<port>`, gives the paths and their count;
  - the spec's `## Acceptance criteria` section, one list item per criterion, gives the count.
- The provenance "set by the team" on the design doc's facts: the design doc is the team's document, so whatever it sets, the team set.
- The two counts: `readCount` counts the ticket and the three product documents; `documentsRead` counts the four linked documents.
- The card data entries listed under Produces, handed to the plan's card data (spec 0050).
- Safe stops: no ticket, or more than one, for the id; a ticket without a title or a status line; a missing, refused (outside the project, secret, symlink) or unreadable link; a link of another kind; a kind missing or linked twice; a missing or malformed field. The plan stops before the line this element would print next, writes no draft and exits non-zero. The reason reaches the person through the CLI's error channel (spec 0048). No line the simulation lacks is printed.

**Change**
- `dft agent plan` takes the ticket's id, not a project-relative file path. The command's argument shape is spec 0048's; resolving the id to the file is this element's.
- The design doc, not CLI input or a recorded input, is where the envelope and the budget are read. Their use as the bound and the budget is specs 0054 and 0056.
- Roles of PRD, design doc and spec come from the file suffix, never from a model.

**Interfaces with neighbours**
- To spec 0054: the guardrail link (path, file, short id) and the design doc's envelope with its source.
- To spec 0056: the budget per run in cents, with its source (the design doc's path), before the plan's first line.
- To spec 0066: the spec's component paths in order, each with the port the spec names beside it, if any.
- To spec 0051: the read documents' text and the extracted facts, the material the judge reasons over.
- To spec 0050: the Ticket and Documents card data and the ticket id for the run's header.

## Testing Decisions

- The seam is the golden test of spec 0048. `dft agent plan 0031-TICKET-book-a-demo-runtime` runs against the fixture project, and the printed lines and the card data are compared byte for byte with the simulation. This element adds:
  - its six lines;
  - the Ticket card;
  - the Documents header;
  - the PRD, design doc and spec rows.
- The fixture holds:
  - under `.agents/`, `0031-TICKET-book-a-demo-runtime.ticket.md`, with the title, `**Status:** ready-for-agent` and Markdown links to the three documents beside it and to the guardrail in `.langdrift/harness/guardrails/`;
  - `0002-PRD-book-a-demo.prd.md` with `**Use case:** book-a-demo request`;
  - `0005-DD-book-a-demo-runtime.dd.md` with `**cpuMillis:** 2000`, `**memoryMiB:** 3072` and `**Budget per run:** $25.00`;
  - `0011-SPEC-book-a-demo-runtime.spec.md`, whose Components section lists `services/signup-api` (with `:8080` beside it), `services/notify-api`, `platform/nats`, `platform/redis` and `apps/signup` in that order, and whose Acceptance criteria section has three items;
  - the guardrail at `.langdrift/harness/guardrails/0001-finops.gr.md`.
- Each safe stop gets one golden case: an unknown id, a missing link and a missing field. Each prints the lines up to the stop, writes no draft and exits non-zero.
- Unit tests only for the rules:
  - the short id of a name;
  - the kind of a suffix;
  - the slug and number of a ticket id;
  - each field and section extraction, including a missing field.
- Tests only observe the printed lines, the card data and the files written, never the reading internals. Prior art: the plan's existing golden and record round-trip tests in `lib/cli`.

## Out of Scope

- The guardrail's criteria, enforcement, `reads` line and FIN line, and the `bounded` line: spec 0054.
- The `budget` line and every money figure: spec 0056.
- The judge and the Analysis header: spec 0051.
- Component inspection, ports and images: spec 0066.
- The command, its argument shape, the comments, the `console ›` prefix, the column layout and the error channel: spec 0048.
- The record that carries the card data to the console: spec 0050.
- The console and its stub API.

## Open Questions

- **The "4 documents" count.** The terminal's `✓ read       4 documents` counts the ticket, the PRD, the design doc and the spec, and reads the guardrail on its own line after it. The card's `4 documents read` lists the PRD, the design doc, the spec and the guardrail. This spec follows both literally, with two counts that happen to agree on the simulation's ticket. Which set does "documents read" mean?
- **Field lines in the client's documents.** Reading the use case, the envelope, the budget, the components and the acceptance criteria relies on fixed field lines and sections in the client's documents (`**Use case:**`, `**Budget per run:**`, `## Components`, `## Acceptance criteria`). The simulation does not say how the documents are written; whether LangDrift imposes this shape is the owner's decision.

## Further Notes

- The field names (`**Use case:**`, `**cpuMillis:**`, `**memoryMiB:**`, `**Budget per run:**`, `## Components`, `## Acceptance criteria`) are the reading contract between this element, the fixture and the templates `dft init` scaffolds (spec 0048). A change to one is a change to all three.
- The port optionally named beside a spec component exists only for spec 0066's port rule, which is itself an open question there. It is not printed by this element.
- The ticket's status is printed as read and gates nothing. The simulation shows `ready-for-agent` and no refusal for other statuses.
