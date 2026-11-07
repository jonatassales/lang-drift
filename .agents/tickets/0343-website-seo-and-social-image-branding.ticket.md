# 0343: replace the website's social images and SEO metadata with the monochrome LangDrift branding

**Spec:** None. The owner's instruction of 2026-11-07: replace the website's legacy Exphere-colored social images (which still carried Vision and AI Voice) and update the SEO metadata to reflect the current LangDrift platform (`LLMOps`, `FinOps`, `Safe Seat`, `Microruntimes`, `LangDrift CLI`, `CI Gates`).

**What to build:** `apps/website/src/app/[locale]/opengraph-image.png` and `apps/website/src/app/[locale]/twitter-image.png` are redrawn in the monochrome LangDrift grid-and-crosshair aesthetic at `1200 × 630` with the drift mark, the `LangDrift` wordmark, `"Agentic LLMOps & FinOps Microruntimes"` and the six current platform items (`LLMOps · FinOps · Safe Seat · Microruntimes · LangDrift CLI · CI Gates`), with no Exphere purple corona and no Vision or AI Voice copy. `opengraph-image.alt.txt` and `twitter-image.alt.txt` match the new headline, `organization.json` points to `https://github.com/jonatassales/lang-drift`, and `pages.json` (plus the corresponding `conversion.json` and `book-a-demo.json` strings) reflects `LLMOps`, `FinOps` and `Microruntimes` across all five locales (`en`, `pt-BR`, `zh-Hant`, `ja`, `vi`).

**Blocked by:** None.

**Status:** completed.

- [x] `apps/website/src/app/[locale]/opengraph-image.png` and `twitter-image.png` are `1200 × 630` PNGs in the monochrome LangDrift branding with no Exphere purple corona and no Vision or AI Voice references.
- [x] `apps/website/src/app/[locale]/opengraph-image.alt.txt` and `twitter-image.alt.txt` read `LangDrift — Agentic LLMOps & FinOps Microruntimes`.
- [x] `apps/website/src/app/metadata/organization.json` names `https://github.com/jonatassales/lang-drift`.
- [x] `apps/website/src/app/metadata/pages.json`, `apps/website/src/app/template/i18n/messages/conversion.json` and `apps/website/src/app/template/i18n/messages/book-a-demo.json` reflect the current platform (`LLMOps`, `FinOps` and `Microruntimes`, with Safe Seat under `LangDrift LLMOps`) in all five locales.
- [x] `pnpm turbo run check`, `pnpm turbo run typecheck --filter=website`, `pnpm turbo run test --filter=website`, `packages/drift/drift.sh check localization` and `packages/drift/drift.sh quality dead-code` pass.

## Dev part

**Override:** the owner directed the agent to replace the social share images and update the SEO metadata strings and open a pull request; redrawing static social PNGs and updating metadata/i18n JSON strings holds no core-domain logic, hard part or new stack (policy, Categories). 2026-11-07.
