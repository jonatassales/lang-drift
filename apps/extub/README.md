# extub

The local Go stub API the [console](/apps/console/README.md) reads. It serves one seeded project, LangDrift Studio, with its four runtimes, their executions, history, secrets and FinOps readings, the routes to create projects and runtimes, start, stop, revoke and load-test them, and the project's tickets. The tickets are generated, newest first, across the ticket, analysis, plan and workflow phases; `console/products/{id}/tickets?reached={phase}&cursor=` pages the ones that have reached a phase, twenty at a time, and `console/tickets/{id}` returns one. What a phase in progress does not know yet is `null`, never zero, and a project created at runtime has no tickets. Every reply is marked `"source": "stub"`; nothing is real.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: start extub on `127.0.0.1:3401` together with the console. |
| `pnpm --filter @langdrift/extub dev` | Start extub alone with `go run`. |
| `pnpm --filter @langdrift/extub test` | Run the Go tests. |

It needs a Go toolchain; `go.work` at the root ties it to the other Go modules. The entrypoint is `cmd/server`. `EXTUB_PORT` changes the port, 3401 by default, and extub listens on `127.0.0.1` only, refusing any other host or a foreign origin. The seed lives in `src/infrastructure/seed/console.json` and resets on every start.

## Hosted preview

Setting `EXTUB_TOKEN` switches extub to a read-only preview for Vercel's `langdrift-stubs` project (spec 0015). It listens on every interface at `PORT`, 3000 by default, as Vercel's Go preset requires. Every request but `/health` must carry `Authorization: Bearer <EXTUB_TOKEN>`, or it is answered `401 token_required`. Any method but `GET` and `HEAD` is answered `403 read_only_preview`, so the seed never changes and no instance diverges from another. There is no live stream: `webrtc/offer` answers `503 stream_unavailable`, and the console keeps polling. Without the token, nothing above applies.

The owner's steps on Vercel, which no agent takes:

1. Generate one token, for example with `openssl rand -base64 32`, and set it as `EXTUB_TOKEN` on `langdrift-stubs` and `langdrift-console`, for production and preview.
2. Set `EXTUB_URL` on `langdrift-console` to the stub's production URL.
3. Turn on Deployment Protection for the console.
