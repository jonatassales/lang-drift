# 0212: the event log and the archive become `infrastructure/messaging` and `infrastructure/storage`

**Spec:** .agents/workflow/specs/0044-published-and-internal-workspaces.spec.md

**What to build:** The Redis Streams component and the MinIO component leave `apps/microservices/`, where they read as services, and become the two product infrastructure components `infrastructure/messaging` and `infrastructure/storage`. Each keeps its Dockerfile and README, gains a manifest with only a name and `private: true`, and its build command in the Dockerfile comment and the README names the new path. The Nest rule says an infrastructure component is not a microservice and lives under `infrastructure/`; the context document and the service READMEs point at the new folders. Both images build from their folders as before.

**Blocked by:** 0211.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] `apps/microservices/` holds no `redis` and no `minio`; `infrastructure/messaging` and `infrastructure/storage` hold the Dockerfile, the README and a manifest named `messaging` and `storage`.
- [ ] `pnpm ls -r --depth -1` lists `messaging` and `storage`.
- [ ] `docker build -t langdrift-messaging:local infrastructure/messaging` and `docker build -t langdrift-storage:local infrastructure/storage` succeed, and the README of each shows that command.
- [ ] The Nest rule's scope line and tree item place infrastructure components under `infrastructure/`; the demo-requests README's local Redis instructions name `infrastructure/messaging`; both Dockerfiles and READMEs describe the components as the event log and the archive of the products, not of the microservices.
- [ ] ADR 0004 gains a note that the component name is its folder, so a fence of this repository names `messaging` and `storage`.
- [ ] `git grep -n 'microservices/\(redis\|minio\)'` finds only workflow history and CLI fixtures.

## Dev part

**Override:** No part qualifies: a folder move with its documentation (policy, Categories: never a rename or a migration run by a tool). 2026-11-04.

## Verification

- `docker build -t langdrift-messaging:local infrastructure/messaging && docker build -t langdrift-storage:local infrastructure/storage`
- `pnpm turbo run check`
- `platform/cli/drift quality check`
