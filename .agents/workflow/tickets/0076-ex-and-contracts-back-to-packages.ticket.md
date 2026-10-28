# 0076: The CLI and its contracts return to packages

**Spec:** .agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md

**What to build:** The `ex` CLI lives at `packages/ex` and the Microruntimes contracts at `packages/contracts`, published as `@langdrift/contracts`. Their Go module paths, every import, the Go and pnpm workspaces, the lockfile and the CI cache follow the folders. The Exphere workspace and its README go; what the README said about the boundary and the Go workspace lives in the two package READMEs. AGENTS.md, the context file, the docs pages that cite the CLI's path and the `.langdrift` sample stop naming Exphere.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `git log --follow` reaches the history of the moved files.
- [ ] The CLI's suite passes under the race detector, the contracts' suite passes, and the build script writes the binary at `packages/ex/dist/ex`.
- [ ] `pnpm install --frozen-lockfile` passes and the lockfile names the new paths and package.
- [ ] No live document names an `exphere/` path or the Exphere workspace.

### Context

- [The spec](/.agents/workflow/specs/0012-packages-back-and-a-finops-stub.spec.md)
- [The workspace it reverses](/.agents/workflow/specs/0010-exphere-and-the-langdrift-folder.spec.md)
