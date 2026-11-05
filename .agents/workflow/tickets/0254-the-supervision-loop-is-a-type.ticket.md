# 0254: the connected agent's supervision loop is a type

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `runOne` in `lib/cli/src/application/connected` becomes a supervision type, one file per method: the state it carried in local variables (the active operation, its cancel, the lease, the last contact, the stop cause, the intent) becomes its fields; `save`, `appendEvent`, `launch` and `shutdown` become methods; the `select` loop dispatches to one method per channel (a worker report, a worker result, a console tick, the context's end). Before the split, scenarios are added at the service's public methods for a lease lost with an active stop, a purge epoch advance and a migration with a Job unit. The rest of the package applies the general rules as the other layers do.

**Blocked by:** 0246, 0222, 0223, 0224, 0225, 0226, 0227.

**Status:** ready-for-agent

- [ ] The supervisor scenarios and service tests pass before and after the split, with the three new scenarios added first.
- [ ] No method of the supervision type is longer than one screen, and none nests beyond two conditional levels.
- [ ] `platform/cli/drift quality check general` prints no finding under `application/connected`.

## Dev part

**Category:** `core-domain`. Reconciling the console's intent with the local record decides whether a runtime is started, stopped or purged; the owner must be able to explain it during an incident.

**Paths:**
- `lib/cli/src/application/connected/reconcile.go`
- `lib/cli/src/application/connected/reconcile_test.go`

**Boundary:** `func (supervision *supervision) reconcile(ctx context.Context, intent Intent) error`, named here and called by the console-tick method after `Console.Sync` answers; the `supervision` type, its fields and `Intent` are declared outside the paths, the type by the agent in `supervision.go`.

**Goal:** Apply the console's answer to the record: adopt a newer revision, start or stop the operation the intent asks for, advance the purge epoch, and record the cause, without ever acting on a stale revision.

- [ ] Each branch of today's reconciliation (newer revision, same revision, stop requested, purge epoch advanced, lease lost) has one test.
- [ ] A stale intent changes nothing in the record.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
