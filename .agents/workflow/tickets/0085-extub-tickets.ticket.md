# 0085: console-stubs serves the project's tickets by phase

**Spec:** .agents/workflow/specs/0014-console-ticket-phases.spec.md

**What to build:** console-stubs generates a deterministic set of about sixty tickets for the seeded project, spread across the ticket, analysis, plan and workflow phases, and serves them through `console/products/{id}/tickets?reached=&cursor=` and `console/tickets/{id}`.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] A ticket carries its number, slug, title, status, phase and artifacts, and an analysis, plan and workflow section once each phase is reached; unknown fields are null, never zero.
- [ ] The early-access ticket carries the website's figures; every plan stays inside its envelope.
- [ ] The listing pages twenty at a time, newest first, with a decimal cursor; `reached` is cumulative and defaults to `ticket`; an unknown phase or a malformed cursor is invalid input.
- [ ] The detail returns one ticket; an unknown id is not found; a project created at runtime has no tickets.

### Verification

- `go test ./...` and `go vet ./...` in `apps/microservices/console-stubs`.
