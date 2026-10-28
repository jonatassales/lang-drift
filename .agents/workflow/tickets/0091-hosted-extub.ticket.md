# 0091: extub runs as a read-only preview on Vercel

**Spec:** .agents/workflow/specs/0015-hosted-extub.spec.md

**What to build:** with `EXTUB_TOKEN` set, extub serves on Vercel. It builds from `cmd/server`, listens on `PORT`, admits only requests carrying the token, serves reads, refuses every change and has no live stream. Without the token it behaves exactly as today.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `main.go` and the composition root live in `cmd/server/`; the package scripts build and run `./cmd/server`.
- [ ] Hosted mode, decided once from `EXTUB_TOKEN`, listens on all interfaces at `PORT` (default 3000); local mode keeps `127.0.0.1:EXTUB_PORT`.
- [ ] Hosted mode admits a request only with `Authorization: Bearer <token>`, compared in constant time, and answers `401 token_required` otherwise; `/health` stays open.
- [ ] Hosted mode answers any method but `GET` and `HEAD` with `403 read_only_preview`, and `webrtc/offer` with `503 stream_unavailable`; the WebRTC hub and the stream loop do not start.
- [ ] The README lists the hosted mode, its variables and the owner's steps on Vercel.

### Verification

- `go test ./...` and `go vet ./...` in `apps/extub`; the Vercel stub deployment builds on the pull request.
