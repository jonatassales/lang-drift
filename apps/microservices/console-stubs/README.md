# extub

The local Go stub API the [console](/apps/console/README.md) reads. It serves one seeded project, LangDrift Studio, with its four runtimes, their executions, history, secrets and FinOps readings, the routes to create projects and runtimes, start, stop, revoke and load-test them, and the project's tickets. The tickets are generated, newest first, across the ticket, analysis, plan and workflow phases; `console/products/{id}/tickets?reached={phase}&cursor=` pages the ones that have reached a phase, twenty at a time, and `console/tickets/{id}` returns one. What a phase in progress does not know yet is `null`, never zero, and a project created at runtime has no tickets. Every reply is marked `"source": "stub"`; nothing is real.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: start extub on `127.0.0.1:3401` together with the console. |
| `pnpm --filter console-stubs dev` | Start extub alone with `go run`. |
| `pnpm --filter console-stubs test` | Run the Go tests with the race detector (`go test -race ./...`), which needs cgo and a C compiler on Linux. |

It needs a Go toolchain; `go.work` at the root ties it to the other Go modules. The entrypoint is `cmd/server`. `EXTUB_PORT` changes the port, 3401 by default, and extub listens on `127.0.0.1` only, refusing any other host or a foreign origin. The seed lives in `src/infrastructure/seed/console.json` and resets on every start. WebRTC peers keep their sockets on loopback, because multicast DNS is off. A browser's `.local` candidates are ignored, and the connection forms through peer-reflexive candidates. While a console runtime page streams, `lsof -nP -iUDP -a -p <extub pid>` shows only `127.0.0.1`.

## Responses and limits

- **Headers.** Every answer, in both modes, carries `Content-Type: application/json`, `Cache-Control: no-store` and `X-Content-Type-Options: nosniff`. It also carries:
  - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
  - `X-Frame-Options: DENY`
  - `Referrer-Policy: no-referrer`
  - `Cross-Origin-Resource-Policy: same-origin`
- **Admission comes first:** `local_only` (403) locally, and `token_required` (401) or `read_only_preview` (403) when hosted.
- **Body size.** A body over 64 KiB answers 413 `request_too_large`. A declared `Content-Length` is refused before routing; a chunked body is refused when a route reads past the limit.
- **Bounds that answer 400 `invalid_input`:**
  - the `console/local/folders` path is absolute, at most 1024 bytes and without NUL, as product paths are;
  - a component's port is 1 to 65535 when set;
  - its image is at most 256 bytes;
  - its readiness path starts with `/` and is at most 256 bytes.

## Hosted preview

Setting `EXTUB_TOKEN` switches extub to a read-only preview for Vercel's `langdrift-stubs` project (spec 0015). It listens on every interface at `PORT`, 3000 by default, as Vercel's Go preset requires. Every request but a `GET` of `/health` must carry `Authorization: Bearer <EXTUB_TOKEN>`, or it is answered `401 token_required` with `WWW-Authenticate: Bearer realm="extub"`. The scheme matches in any case, and the token is compared as SHA-256 digests in constant time, so neither its value nor its length leaks through timing. A token shorter than 32 characters stops extub at start. Locally, leave `EXTUB_TOKEN` unset: with it, extub binds every interface on port 3000. Any method but `GET` is answered `403 read_only_preview`, so the seed never changes and no instance diverges from another. There is no live stream: `webrtc/offer` answers `503 stream_unavailable`, and the console keeps polling. Without the token, nothing above applies. `vercel.json` builds `./cmd/server` with Go itself, so Vercel's Turborepo detection never runs `turbo run build` for this project.

The owner's steps on Vercel, which no agent takes:

1. Generate one token, for example with `openssl rand -base64 32`, and set it as `EXTUB_TOKEN` on `langdrift-stubs` and `langdrift-console`, for production and preview.
2. Set `EXTUB_URL` on `langdrift-console` to the stub's production URL.
3. Point the Root Directory of `langdrift-stubs` at `apps/microservices/console-stubs`, the folder this service lives in since spec 0044.
4. Turn on Deployment Protection for the console's production and preview deployments.
5. Replace `CRYPTO_SECRET` on `langdrift-website` with a canonical 32-byte base64 key, for example from `openssl rand -base64 32`; the website answers 500 until then.
