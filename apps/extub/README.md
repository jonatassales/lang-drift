# extub

The local Go stub API the [console](/apps/console/README.md) reads. It serves one seeded project, LangDrift Studio, with its four runtimes, their executions, history, secrets and FinOps readings, the routes to create projects and runtimes, start, stop, revoke and load-test them, and the project's tickets. The tickets are generated, newest first, across the ticket, analysis, plan and workflow phases; `console/products/{id}/tickets?reached={phase}&cursor=` pages the ones that have reached a phase, twenty at a time, and `console/tickets/{id}` returns one. What a phase in progress does not know yet is `null`, never zero, and a project created at runtime has no tickets. Every reply is marked `"source": "stub"`; nothing is real.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: start extub on `127.0.0.1:3401` together with the console. |
| `pnpm --filter @langdrift/extub dev` | Start extub alone with `go run`. |
| `pnpm --filter @langdrift/extub test` | Run the Go tests. |

It needs a Go toolchain; `go.work` at the root ties it to the other Go modules. `EXTUB_PORT` changes the port, 3401 by default, and extub listens on `127.0.0.1` only, refusing any other host or a foreign origin. The seed lives in `src/infrastructure/seed/console.json` and resets on every start.
