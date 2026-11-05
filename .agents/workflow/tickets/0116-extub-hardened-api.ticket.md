# 0116: console-stubs answers like a hardened API

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:**

- **Headers.** Every console-stubs response adds:
  - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
  - `X-Frame-Options: DENY`
  - `Referrer-Policy: no-referrer`
  - `Cross-Origin-Resource-Policy: same-origin`
- **Authentication.**
  - 401 carries `WWW-Authenticate: Bearer realm="console-stubs"`.
  - The Bearer scheme is matched case-insensitively.
  - The token is compared as SHA-256 digests in constant time.
- **Body limit.** An oversized body answers 413 `request_too_large`.
- **Input bounds**, which answer 400 when broken:
  - the `local/folders` path must be absolute, at most 1024 bytes, without NUL;
  - a component port is 1 to 65535;
  - an image is at most 256 bytes;
  - a readiness path starts with `/` and is at most 256 bytes.
- **WebRTC.** The setting engine disables multicast DNS, so peers open loopback sockets only.
- **Tests.** The test script runs the race detector, and tests pin all of the above, together with the fetch-metadata check and the caps over HTTP.

**Blocked by:** 0107.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Every response, in both modes, carries the four headers next to today's three.
- [ ] A missing or wrong token answers 401 with `WWW-Authenticate`, and `bearer <token>` is admitted.
- [ ] A body over 64 KiB answers 413.
- [ ] `Sec-Fetch-Site: cross-site` and a foreign `Origin` answer 403 in local mode.
- [ ] The product, runtime and secret caps answer 429 over HTTP.
- [ ] Out-of-bounds folder paths and runtime fields answer 400.
- [ ] `go test -race ./...` passes.
- [ ] While a WebRTC peer is connected, console-stubs holds no UDP socket on a non-loopback address.

### Verification

- `go test -race ./...` in `apps/microservices/console-stubs`; `lsof -nP -iUDP -a -p <console-stubs pid>` during a console runtime stream; `curl` against both modes.
