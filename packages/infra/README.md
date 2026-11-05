# infra

The infrastructure the LangDrift products run on, as one private package of this repository: the event log, the archive and the transport the services share through them. Nothing here is a microservice and nothing is published.

| Folder | What it is |
| --- | --- |
| [`messaging/`](/packages/infra/messaging/README.md) | The event log: a Redis image with append-only persistence, holding the stream `demo-requests:events`. A Dockerfile and a README. |
| [`storage/`](/packages/infra/storage/README.md) | The archive: MinIO, built from its source, where the `archive` consumer of `demo-requests` writes every event. A Dockerfile and a README. |
| `src/streaming/` | The Redis Streams transport and the event schemas the Nest services exchange events with: `infra/streaming` exports the events and their Zod schemas, `infra/streaming/transport` the custom strategy, the client and the stream connection. |

`messaging` and `storage` are microruntime components named after their folders (ADR 0004): only Docker builds them, from the commands in their READMEs. `streaming` is TypeScript, built to `dist/` with `tsc` and tested with Vitest:

| Script | What it does |
| --- | --- |
| `build` | `tsc -p tsconfig.build.json`, emitting `dist/streaming/` for the services. |
| `typecheck` | `tsc --noEmit` over `src/streaming/`. |
| `test` | `vitest run` over `src/**/*.unit.ts`. |

`demo-requests` and `mailer` depend on it as `"infra": "workspace:*"` and import `infra/streaming` and `infra/streaming/transport`; `.claude/rules/11-nestjs.rule.md` sets how they use it.
