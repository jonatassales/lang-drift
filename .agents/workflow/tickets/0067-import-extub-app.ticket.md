# 0067: Import the `extub` app

**Spec:** None

**What to build:** copy Exphere's `apps/extub`, the local Go API of FinOps console stubs and seeds, into `apps/extub`, keeping its code as it is and adapting only what ties it to the Exphere repository: the package name becomes `@langdrift/extub`, the Go module path becomes `github.com/jonatassales/tao/apps/extub` (with its imports), the Go directive and toolchain match `packages/ex`, and `.agents/context.md` describes it. The seed then holds products instead of projects: `console.json` seeds `Microruntimes` first, which keeps the commerce fixture and its runtimes, and `Vision` second, with no components or runtimes yet, and the domain, the `/console/products` routes, the `productId` fields and the tests follow. Only files tracked by the source repository are copied, minus its Markdown files; its README is written when the app is adapted.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `apps/extub` holds the tracked files of the source, without any `.md` file.
- [ ] The seed's `products` array is `Microruntimes` then `Vision`, and no `project` name remains in the Go code, routes or seed.
- [ ] No reference to the Exphere repository's module path, package scope or root scripts remains in the code.
- [ ] Biome formatting applied to the files it flags; nothing else is reformatted.
- [ ] `pnpm turbo run check`, `typecheck` and `test` pass with `extub` in the workspace, `pnpm install --frozen-lockfile` passes, and `drift check dead-code` passes.

### Verification

- `pnpm turbo run check typecheck test`, `drift check dead-code`.
