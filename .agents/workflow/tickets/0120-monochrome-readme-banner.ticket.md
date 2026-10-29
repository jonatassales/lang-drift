# 0120: The README banner is monochrome and the opening keeps only the stack badges

**What to build:** the README banner follows Vercel's restraint, and no color is left in it. A light and a dark variant: pure white or black, a 64 px cell grid fading from the center, a soft spotlight from the top, one framed cell with crosshairs at its corners and a beam along its top edge. In it, the drift mark and the LangDrift wordmark large and centered, over "Agentic FinOps Platform" in the muted tone. Every README (English at the root, pt-BR, zh-Hant, ja and vi under `platform/docs/`) drops the "development candidate" status badge and the languages badge, and the status section no longer calls `ex` a development candidate.

**Status:** ready-for-agent

### Banner

- [ ] `platform/docs/assets/images/langdrift-banner-{dark,light}.{svg,png}` use only black, white and grays: no rainbow, purple or peach.
- [ ] Each PNG is rendered at 2560×1280 from its SVG, and the `<picture>` in every README serves the dark variant to dark-theme readers.
- [ ] The alt text of the banner reads "LangDrift: Agentic FinOps Platform", the text the image shows.

### Badges and status

- [ ] No README carries the status badge or the languages badge; the pull request checks, Go, Node.js, Next.js and pnpm badges remain.
- [ ] The status section says that `ex` is built from source with no binary published, in every language.

### Verification

- `pnpm exec biome check .` and `drift check dead-code` pass.
- Every in-page anchor and repository path in the five READMEs resolves, and their code blocks stay identical.
- Both banners are inspected rendered, at 2×.
