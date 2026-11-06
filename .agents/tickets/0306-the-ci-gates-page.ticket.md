# 0306: CI Gates has its own page

**Spec:** .agents/specs/0046-the-website-focus-pass.spec.md

**What to build:** A team lead opens `/platform/ci-gates` and reads what `ex init --ci-gates` writes, the workflow's posture (read-only token, actions pinned by SHA, a pinned `npx @langdrift/cli gates report`, no overwrite), the report a pull request receives with FinOps readings that read "unknown" in CI and never zero, and the run-time FinOps gates. Renderings of pull request #122 of this repository, in light and dark, are served from `public/ci-gates/`, the one matching the theme shown through `dark:` utilities: the real `ex gates report` output on its commits, and the repository's own CI Gates comment, whose figures the page says are set by the dashboard script. GitHub's page does not render headlessly, so each is rendered from its Markdown through GitHub's Markdown API. The h1 is the shared tile with the CI Gates icon. Catalog, metadata, route path and sitemap as in 0305.

**Blocked by:** 0305.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `/platform/ci-gates` renders in every locale with both renderings, one per theme, each with localized alt text.
- [x] The SEO contract test covers its metadata and sitemap entries.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
