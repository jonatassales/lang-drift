# 0100: The plan resolves the components, sizes them and checks the draft

**Spec:** .agents/specs/0016-ex-performs-the-microruntimes-session.spec.md

**What to build:** the plan's third act prints for real: the components named by the facts or given as input are inspected by their manifests (a Go module, a Node.js or Next.js package, a NATS or Redis configuration, a Python project, a Dockerfile) for their technology, kind and port, each gets a proposed share from a sizing table by kind unless a fact states figures, the totals and the headroom are printed against the envelope, the executed shares are the equal split the connected agent applies, the cross-checks run (distinct ports, resolved images, allowed recipes), gpu and disk requests and the absence of a runtime price source are reported, the fence is validated against the schema, and the draft is written with its revision and proposed for authorization. Images, readiness paths and build recipes come from facts or input; a missing one blocks the plan with the field named.

**Blocked by:** 0099.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Against the fixture project, lines 25 to 43 print with the sample's component, totals, checks and revision, and the plan fence's `components`, `totals`, `materialization`, `price`, `checks` and `runtime` equal the expected fixture.
- [ ] The manifest inspector is unit-tested on each manifest kind, including the port sources, and a path with no manifest is recorded as unknown technology.
- [ ] A duplicate port, an unresolvable image or a recipe outside the allowlist leaves the plan `blocked` with the failing check and exit code 3.
- [ ] The plan's equal split and the connected agent's materialization share one rule, so the fence the plan writes materializes without conflict.

### Verification

- `GOTOOLCHAIN=local go test -race ./...`; the plan command on this repository with `--no-judge` shows the five-line component block for the components given.
