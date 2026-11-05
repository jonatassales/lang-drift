# 0310: one infra package holds the infrastructure

**Spec:** None. The owner's instruction of 2026-11-05: the root `infrastructure/` folder moves under `packages/` as one repo package named `infra`, and the `streaming` package moves into it.

**What to build:** `packages/infra` is one private workspace package. `messaging/` and `storage/` keep their Dockerfile and README, without manifests of their own, and keep their component names (ADR 0004). `src/streaming/` holds the Redis Streams transport and the event schemas, built to `dist/` and exported as `infra/streaming` and `infra/streaming/transport`. `demo-requests` and `mailer` depend on `infra` and import those paths; the root `infrastructure/` folder and the `streaming` package are gone, and every caller (the workspace list, Docker's ignore file, Biome's import matcher, drift, the rules, the guide and the READMEs) follows.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [x] `pnpm-workspace.yaml` lists no `infrastructure/*`, and `pnpm install --frozen-lockfile` passes with `infra` as the services' dependency.
- [x] `infra` builds, typechecks and runs the streaming tests; `demo-requests` and `mailer` build, typecheck and test against it.
- [x] `turbo prune` for each service keeps `packages/infra`.
- [x] The ticket trailer covers `apps/`, `packages/` and `lib/`, the three roots left.

## Dev part

**Override:** a move of infrastructure files and one package's sources; no part qualifies under the Safe Seat policy. 2026-11-05.
