# 0305: the LangDrift CLI has its own page

**Spec:** .agents/workflow/specs/0046-the-website-focus-pass.spec.md

**What to build:** A developer opens `/platform/cli` and reads how to install the LangDrift CLI (`npx @langdrift/cli init`, its requirements), how a run goes from `init` to `agent status`, what it governs (budget reservations with the gates FG-01 to FG-06, bounded runtimes, the agent driver), with the `ex agent plan` screenshot and links to npm and GitHub. The h1 is the shared tile with a new CLI icon in the `packages/web` icon family. The page has a catalog in five locales, registered in the loaders, shape validator, barrel and parity test, a page-metadata entry, a route path and five sitemap entries.

**Blocked by:** 0304.

**Status:** ready-for-agent

- [x] `/platform/cli` renders in every locale with its icon, install line, flow and image.
- [x] The SEO contract test covers its metadata and sitemap entries.
- [x] The CLI icon has a row in the platform icons test.

## Dev part

**Override:** No part qualifies: copy, content, layout and removals on the existing page pattern, with no problem of its own to solve (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=website --filter=web`
- `platform/cli/drift check localization`
