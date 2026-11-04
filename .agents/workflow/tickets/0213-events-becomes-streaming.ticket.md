# 0213: the event log package is `streaming`

**Spec:** .agents/workflow/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The internal package that holds the demo-request events and the Redis Streams transport loses its scope and takes a name that does not collide with Node's `events` module: `packages/streaming`, workspace name `streaming`. The two Nest services depend on and import `streaming` and `streaming/transport`, their READMEs and the Nest rule name it, and the Docker build context admits the new folder so both service images still prune and build from the repository root.

**Blocked by:** 0211.

**Status:** ready-for-agent

- [ ] `packages/events` no longer exists; `packages/streaming` exports `.` and `./transport` under the name `streaming`.
- [ ] demo-requests and mailer depend on `streaming` through `workspace:*` and import nothing from `@langdrift/events`; `pnpm turbo run build typecheck test --filter=demo-requests --filter=mailer` passes.
- [ ] `docker build -f apps/microservices/demo-requests/Dockerfile .` and the mailer equivalent build from the repository root.
- [ ] The Nest rule's event rule and both service READMEs name `streaming`.
- [ ] `git grep -n '@langdrift/events\|packages/events'` finds only workflow history and CLI fixtures; the context document names `streaming`.

## Dev part

**Override:** No part qualifies: a package rename across its consumers (policy, Categories: never a rename). 2026-11-04.

## Verification

- `pnpm install --frozen-lockfile && pnpm turbo run build typecheck test check --filter=streaming --filter=demo-requests --filter=mailer`
- `docker build -f apps/microservices/demo-requests/Dockerfile -t langdrift-demo-requests:local . && docker build -f apps/microservices/mailer/Dockerfile -t langdrift-mailer:local .`
