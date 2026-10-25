# 0037: Microruntimes page and copy

**Spec:** .agents/workflow/specs/0005-finops-and-microruntimes-pages.spec.md

**What to build:** `/[locale]/platform/microruntimes/` renders in the five locales with the Vision page's layout: header, hero with the question headline and calls to action, the annotated runtime file, the ownership table, the state cards, the execution ledger, the console-to-agent flow, the FAQ, the demo band and the footer, with its own metadata and sitemap entries. The hero's demo slot, the lifecycle scenarios and the per-step evidence arrive with ticket 0038.

**Blocked by:** 0035.

**Status:** ready-for-agent

### Copy and data

- [ ] `microruntimes.json` in the five locales, same key order everywhere: `skip`, `hero`, `descriptor`, `ownership`, `state`, `evidence`, `flow`, `faq`, plus the `lifecycle` and `scenarios` keys that 0038 consumes.
- [ ] The Microruntimes content module: the six lifecycle stages in order, the runtime summary (components, `cpuMillis`, `memoryMiB`, desired `running`, observed `unknown`), the runtime excerpt within the `ex` limits, its annotation ids, the scenario and phase ids, the execution ledger, the initial stage and the catalog type; a unit test for the stage order, desired versus observed and the excerpt limits.
- [ ] The catalog index asserts the shape and exports the typed catalog; the parity test's family list gains it.

### Sections

- [ ] `RuntimeDescriptorCard` (server): the runtime excerpt as monospace JSON beside its numbered annotations and an "illustrative" badge; the code scrolls sideways at 375px.
- [ ] `AuthoringFlow` (server): console → connected agent → kind, Terraform and Kubernetes with the existing icons and inline arrows that turn downward below the large breakpoint.
- [ ] The ownership matrix as a real table with column and row headers; the state section as three `InfoCards`; the execution ledger as a `DetailList`.

### Route and metadata

- [ ] The page with the Vision boilerplate and the sections above, the FAQ, the demo band and the footer.
- [ ] `pages.json` gains `microruntimes` in five locales; the nested-path map gains `platform/microruntimes/`; `sitemap.xml` gains the five Microruntimes URLs with their alternates (30 in total).

### Verification

- The checks of 0035.
- `/en/platform/microruntimes/` renders with the Microruntimes title and canonical; `/ja/platform/microruntimes/` renders translated copy; `/en/sitemap.xml` lists 30 URLs.
