# 0306: CI Gates has its own page

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** A team lead opens `/platform/ci-gates` and reads what `ex init --ci-gates` writes, the workflow's posture (read-only token, actions pinned by SHA, a pinned `npx @langdrift/cli gates report`, no overwrite), the report a pull request receives with FinOps readings that read "unknown" in CI and never zero, and the run-time FinOps gates. Screenshots of the "LangDrift CI Gates" comment and its checks, captured from a pull request of this repository with headless Chrome in light and dark, are served from `public/ci-gates/`, the one matching the theme shown through `dark:` utilities. The h1 is the shared tile with the CI Gates icon. Catalog, metadata, route path and sitemap as in 0305.

**Blocked by:** 0305.

**Status:** ready-for-agent

- [ ] `/platform/ci-gates` renders in every locale with both screenshots, one per theme, each with localized alt text.
- [ ] The SEO contract test covers its metadata and sitemap entries.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
