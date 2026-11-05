# 0267: the tokens stylesheet is built by shell and jq

**Spec:** .agents/workflow/specs/0037-tooling-rules-compliance.spec.md

**What to build:** `packages/tokens/build.sh` replaces `build.js` and the `build` script becomes `sh build.sh`. The script checks for `jq` at the operation and fails naming it, concatenates the sources in the order Style Dictionary read them (the shared files in name order, then the theme) and pipes them through one `jq` program beside it. `style-dictionary` leaves `devDependencies`. The output is proven byte-identical to the current stylesheet before `build.js` is deleted, and the Vercel build image is confirmed to ship `jq`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `pnpm turbo run build --filter=tokens` writes a stylesheet identical to the one `build.js` wrote, for every theme.
- [ ] `build.js` and `style-dictionary` are gone.
- [ ] A Vercel preview of the website builds; the ticket records whether the image ships `jq`.

## Dev part

**Category:** `new-stack`. The DTCG flattening and alias resolution in `jq` is the owner's first `jq` program of this size.

**Paths:**
- `packages/tokens/tokens.jq`
- `packages/tokens/tokens-jq.test.js`

**Boundary:** `jq -r -s -f tokens.jq`, reading the concatenated DTCG sources as a slurped array and writing the stylesheet to standard output; `build.sh`, written by the agent, calls it.

**Goal:** Flatten the DTCG tree (skipping `$type` and `$extensions` as names), name a token by its path without its scheme group joined with `-`, resolve a `{path}` alias to `var(--name)` of its target, emit the `:root`, `.dark`, bridge, `@custom-variant` and `@theme inline` blocks in the current order, and exit non-zero naming the token path on a missing reference or a cycle.

- [ ] Its output equals the current stylesheet byte for byte for every theme.
- [ ] A missing reference and a cycle each exit non-zero naming the token path; two runs give the same output.

## Verification

- `pnpm turbo run check typecheck test build --filter=web --filter=tokens`
- `pnpm turbo run build --filter=website`
