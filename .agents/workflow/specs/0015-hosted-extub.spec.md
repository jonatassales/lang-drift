# Hosted extub: the console and its stub on Vercel

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

## Problem Statement

The repository is now connected to Vercel, with three projects: `langdrift-website`, `langdrift-console` and `langdrift-stubs`. The console and extub were built to run only on the owner's machine. extub listens on `127.0.0.1:3401`, refuses every host but `localhost`, keeps its state in one process's memory and streams over WebRTC with loopback candidates only. The console forwards only to an extub on loopback and serves only loopback hosts. On Vercel the stub never builds, because its `main.go` is not where Vercel looks. The console builds, but every data read answers `403 local_only`, so the hosted console is a shell with no data.

## Solution

extub and the console gain a hosted mode, chosen by one secret, `EXTUB_TOKEN`. Without it, nothing changes: both stay on loopback, as the README contract says. With it, extub becomes a read-only preview server on Vercel:
- It listens on Vercel's port and answers only requests that carry the token.
- It serves every read and refuses every change with its own error.
- It has no live stream.

The console, with the same token and an HTTPS `EXTUB_URL`, forwards to that extub from any host. The owner keeps the hosted console private behind Vercel's Deployment Protection. When a WebRTC offer fails, the console falls back to the polling it already does.

## User Stories

1. As the owner, I want the hosted console to show the seeded project, its runtimes, executions, history, tickets, analysis, plans and workflows, so that I can review the console without running it locally.
2. As the owner, I want the hosted console to stay private behind Vercel's Deployment Protection, so that the stub data and the stub sign-in are not public.
3. As the owner, I want every action that changes data to be refused on the hosted console with a clear message, so that no state appears, diverges between instances or disappears.
4. As the owner, I want the hosted extub to refuse every request without the shared token, so that only my console reads it.
5. As the owner, I want the local console and extub to work exactly as today when no token is set, so that development keeps its loopback guarantees.
6. As the owner, I want the live runtime stream to step aside when hosted and the screens to keep refreshing by polling, so that nothing waits on a WebRTC peer that Vercel cannot reach.
7. As the owner, I want the stub project on Vercel to build from the repository as it is, so that its deployment check passes on every pull request.
8. As a developer, I want the hosted mode decided in one place on each side and tested, so that a missing or wrong token never opens the stub.
9. As a developer, I want the READMEs to list the variables each project needs and the owner's steps on Vercel, so that wiring a new environment is a checklist.

## Implementation Decisions

- **One switch.** On each side, hosted mode is `EXTUB_TOKEN` being set and non-empty. Every hosted behavior follows from it. Without it, the local contract holds unchanged.
- **extub entrypoint.** `main.go` and the composition root move from `src/` to `cmd/server/`, one of the paths Vercel's Go preset looks for. The package scripts build and run `./cmd/server`. The packages under `src/` do not move. The project's `vercel.json` sets the Go preset and builds `./cmd/server` with `go build`, because Vercel's Turborepo detection would otherwise run `turbo run build` without the workspace installed.
- **extub listening.** In local mode extub listens on `127.0.0.1:EXTUB_PORT`, default 3401, as today. In hosted mode it listens on all interfaces at `PORT`, default 3000, as Vercel's Go preset requires.
- **extub admission.** In hosted mode the loopback host and origin checks give way to the token. A request without `Authorization: Bearer <token>` is answered `401 token_required`. The token is compared in constant time. Reading `/health` stays open so the platform can probe it.
- **Read-only.** In hosted mode every method other than `GET` is answered `403 read_only_preview` before any route runs, except the WebRTC offer, which answers `503 stream_unavailable` below. This covers creating projects and runtimes, start, stop, revoke, delete, secrets, recommendations, load tests and auth posts. The seed is generated per instance from the embedded data, and nothing writes to it.
- **No stream when hosted.** In hosted mode the WebRTC hub and the stream loop do not start, and `webrtc/offer` answers `503 stream_unavailable`. The console's stream hook already treats a failed offer as no stream and keeps polling.
- **Console link.** In hosted mode `extubRequest` accepts an `https` `EXTUB_URL` on any host, without credentials in the URL. It sends the token as `Authorization: Bearer` and skips the loopback host check for the console's own host, keeping the origin and cross-site checks. Server-side reads skip the loopback host check the same way. In local mode both guards stay as they are. The rule lives in one pure function the Node tests cover.
- **Messages.** The console's `errors` catalog gains `read_only_preview`, `token_required` and `stream_unavailable` in its four locales, so a refused action reads as the preview's rule, not as a failure.
- **Owner's steps on Vercel**, documented in the READMEs and not done by an agent:
  - generate one `EXTUB_TOKEN` and set it on `langdrift-stubs` and `langdrift-console`;
  - set `EXTUB_URL` on `langdrift-console` to the stub's production URL;
  - turn on Deployment Protection for the console's production and preview deployments;
  - replace the website's `CRYPTO_SECRET` with a canonical 32-byte base64 key.
- **Turbo.** `EXTUB_TOKEN` joins the pass-through environment beside `EXTUB_URL`, so local tasks see it when set.

## Testing Decisions

A good test observes what a client or the console sees, never internal structure.

1. **extub** (Go, `httptest`, beside the existing handler tests):
   - Without a token: a request from a non-loopback host is still `403 local_only`.
   - With a token: a request without the header is `401`, and with the wrong token `401`. The right token and `GET` gets the products. The right token and `POST` gets `403 read_only_preview`. `/health` answers without the token. `webrtc/offer` is `503 stream_unavailable`.
   - The hosted listen address uses `PORT`, and the local one `127.0.0.1:EXTUB_PORT`.
2. **Console logic** (Node test runner):
   - The pure link rule refuses an `https` `EXTUB_URL` without a token and accepts it with one.
   - Without a token it keeps refusing non-loopback hosts. With a token it accepts any host and never accepts credentials in the URL.
3. **Repository checks**: `go test ./...` and `go vet ./...` in `apps/extub`, `pnpm turbo run check typecheck test`, `drift check localization`, `drift check dead-code`.

Manual verification: run extub with a token and `PORT`, and the console against it with the token and an `https`-shaped `EXTUB_URL` on loopback. Then check four things. Local mode without the token behaves as before. On Vercel, once the owner has set the variables, the stub's deployment check passes and the console shows data.

## Out of Scope

- Live streaming on Vercel, which would need a TURN server.
- Persisting changes on the hosted stub, or any write path.
- Real authentication in the console; the sign-in stays a stub behind Deployment Protection.
- Setting secrets, environment variables or protection on Vercel; those are the owner's steps.
- A custom domain for any project.

## Further Notes

The decisions come from the owner's grilling round on 2026-10-28:
- the hosted console is a protected preview;
- hosted writes are refused;
- a shared token links the two projects;
- polling replaces the stream.

`.agents/workflow/research/ex-microruntimes-capability.research.md` describes the loopback guards this spec relaxes only under the token.
