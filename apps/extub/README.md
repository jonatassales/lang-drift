# extub

The local Go stub API the [console](/apps/console/README.md) reads. It serves FinOps and runtimes only: one seeded project, LangDrift Studio, with its four runtimes, their executions, history, secrets and FinOps readings, and the routes to create projects and runtimes, start, stop, revoke and load-test them. Every reply is marked `"source": "stub"`; nothing is real.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: start extub on `127.0.0.1:3401` together with the console. |
| `pnpm --filter @langdrift/extub dev` | Start extub alone with `go run`. |
| `pnpm --filter @langdrift/extub test` | Run the Go tests. |

It needs a Go toolchain; `go.work` at the root ties it to the other Go modules. `EXTUB_PORT` changes the port, 3401 by default, and extub listens on `127.0.0.1` only, refusing any other host or a foreign origin. The seed lives in `src/infrastructure/seed/console.json` and resets on every start.
