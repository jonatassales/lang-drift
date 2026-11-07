# 0127: The virtual list leaves the console

**Spec:** .agents/specs/0020-console-activity-scrolls-with-the-page.spec.md

**What to build:** once both overviews scroll with the page, nothing uses the virtual list, the activity row skeleton or the virtualization dependency, so they leave the console. The console builds and installs without them, and the dead-code check reports nothing new.

**Blocked by:** 0125, The Executions overview scrolls with the page; 0126, The History overview scrolls with the page.

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] `VirtualList` is deleted.
- [x] `ActivityRowSkeleton` is deleted.
- [x] `@tanstack/react-virtual` is removed from the console's dependencies through pnpm, so the manifest and the lockfile change and nothing else in the lockfile moves; the copy in the `ex` package's sample project fixture is untouched.
- [x] `ActivityListSkeleton` stays, because the per-runtime History page still loads with it.
- [x] `drift check dead-code` passes with no new entry in the allow list.

### Verification

- `pnpm install --frozen-lockfile`, `pnpm --filter console typecheck`, `pnpm --filter console test`, `pnpm turbo run build --filter=console`, `pnpm turbo run check`, `drift check dead-code`, `drift check localization`.
- In the browser: both overviews and the per-runtime History page load as in tickets 0125 and 0126.
