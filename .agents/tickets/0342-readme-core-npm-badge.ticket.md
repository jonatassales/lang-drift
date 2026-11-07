# 0342: add the @langdrift/core npm badge to the root README and its translations

**Spec:** None. The owner's instruction of 2026-11-07: add the `@langdrift/core` npm badge alongside the `@langdrift/cli` badge in `README.md` and in the translated READMEs under `packages/docs/`.

**What to build:** `README.md`, `packages/docs/README.pt-br.md`, `packages/docs/README.zh-hant.md`, `packages/docs/README.ja.md` and `packages/docs/README.vi.md` carry the `@langdrift/core` npm version badge next to the `@langdrift/cli` badge, linking to `https://www.npmjs.com/package/@langdrift/core`.

**Blocked by:** None.

**Status:** completed.

- [x] `README.md` displays the `@langdrift/core` npm badge beside `@langdrift/cli`.
- [x] `packages/docs/README.pt-br.md`, `packages/docs/README.zh-hant.md`, `packages/docs/README.ja.md` and `packages/docs/README.vi.md` display the `@langdrift/core` npm badge beside `@langdrift/cli`.
- [x] `pnpm turbo run check`, `pnpm turbo run typecheck` and `pnpm turbo run test` pass.

## Dev part

**Override:** the owner directed the agent to add the `@langdrift/core` badge to the root README and its translations and open a pull request; adding a static README badge holds no core-domain logic, hard part or new stack (policy, Categories). 2026-11-07.
