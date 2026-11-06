# 0067: Import the `extub` app

**Spec:** None

**What to build:** copy Exphere's `apps/extub`, the local Go API of FinOps console stubs and seeds, into `apps/extub`, keeping its code as it is and adapting only what ties it to the Exphere repository: the package name becomes `@langdrift/extub`, the Go module path becomes `github.com/jonatassales/tao/apps/extub` (with its imports), the Go directive and toolchain match `packages/ex`, and `.agents/context.md` describes it. The seed then holds products instead of projects: `console.json` seeds `Microruntimes` first, which keeps the commerce fixture and its runtimes, and `Vision` second, carrying the canonical Vision dashboard and workspace fixtures extracted from `ld-stub`, and the domain, the `/console/products` routes, the `productId` fields and the tests follow. Copy Exphere's `apps/console` into `apps/console`, adapt it to `tao`'s layout, `@langdrift/tokens` roles and `@langdrift/web` utilities, and connect it to `extub`'s `/console/products` contract. Only files tracked by the source repository are copied, minus its Markdown files; its README is written when the app is adapted.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `apps/extub` and `apps/console` hold the tracked files of the source, without any `.md` file.
- [x] The seed's `products` array is `Microruntimes` then `Vision` (with its `visionFixture` from `ld-stub`), and no `project` name remains in the Go code, routes or seed.
- [x] `apps/console` connects to `extub`, builds, and passes all repository checks.
- [x] No reference to the Exphere repository's module path, package scope or root scripts remains in the code.
- [x] Biome formatting applied to the files it flags; nothing else is reformatted.
- [x] `pnpm turbo run check`, `typecheck` and `test` pass with `extub` and `console` in the workspace, `pnpm install --frozen-lockfile` passes, and `drift check dead-code` passes.

### Verification

- `pnpm turbo run check typecheck test`, `drift check dead-code`.
