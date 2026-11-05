# 0231: the console's route handlers define their methods inline

**Spec:** .agents/workflow/specs/0041-nextjs-rules-compliance.spec.md

**What to build:** The catch-all console route defines `GET`, `POST`, `PATCH` and `DELETE` itself: each extracts the console path, refuses an unsafe path with the 404 reply and calls the console-stubs request with the `console/` prefix. The three product routes (`executions`, `history`, `secrets`) go, because the catch-all answers them identically. The console-stubs route defines `GET` and `POST` with its allowed-path check inline and the allowed paths as a constant of the file. The server read builds its read and calls the bridge itself. `consoleHandler` is deleted. Every response is unchanged.

**Blocked by:** 0230.

**Status:** ready-for-agent

- [ ] No `route.ts` in the console is a forwarder: each HTTP method's dispatch, validation and reply sit in the file.
- [ ] The product `executions`, `history` and `secrets` routes are removed and the catch-all answers their paths with the same responses.
- [ ] A node test beside the catch-all route covers the 404 for an unsafe path and the forwarding of a safe one.
- [ ] `consoleHandler` and its file are gone and nothing is left dead.

## Dev part

**Override:** No part qualifies: the ticket moves, renames or reshapes existing code without changing what it does, and a rename or a migration is never a dev part (policy, Categories). 2026-11-04.

## Verification

- `pnpm turbo run check typecheck test build --filter=console`
