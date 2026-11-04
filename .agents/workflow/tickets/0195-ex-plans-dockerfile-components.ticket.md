# 0195: ex plans Dockerfile components

**Spec:** .agents/workflow/specs/0032-components-by-dockerfile-path.spec.md

**What to build:** `ex agent plan` takes a component as `name=,path=` plus the optional keys (`dockerfile`, `port`, `readiness`, `env-file`, `pass-env`, `kind`, `tech`), reads the component's Dockerfile, takes the port from its first `EXPOSE` when the flag omits it, parses its `HEALTHCHECK` in both forms with its options into the probe, derives the image tag from the runtime id, the component name and the revision, and writes a version 2 fence. The cross-checks name a missing Dockerfile, a Dockerfile without its own `HEALTHCHECK`, a missing env file and colliding ports. The golden fixtures, the planner's test project (which gains a Dockerfile with a `HEALTHCHECK`), the plan record and the CLI README's fence and flag examples change with it. Shares the integration branch of 0194 and 0196.

**Blocked by:** 0194.

**Status:** ready-for-agent

- [x] `ex agent plan <ticket> --component name=web,path=apps/web` on the test project writes a fence whose component has the derived image, the `EXPOSE` port and the `HEALTHCHECK` probe.
- [x] `HEALTHCHECK` in shell and exec form, with and without interval, timeout, retries and start period, maps to the probe fields; `HEALTHCHECK NONE` and no `HEALTHCHECK` fail the `readiness-known` check with the file named.
- [x] `--component readiness=tcp` and `readiness=http:/path` produce the explicit kinds; `env-file=` and `pass-env=` are recorded and checked.
- [x] The findings name `ports-distinct`, `dockerfiles-exist`, `readiness-known`, `env-file-exists` and `within-ceiling`.
- [x] The golden fixtures reproduce the expected `runtime.md`, `plan.json` and `run.md`; the README shows the new fence and flags.
- [x] `go test -race ./...` passes in the CLI module's plan and session packages.

### Verification

- `pnpm turbo run test typecheck --filter=@langdrift/ex`
- `platform/cli/drift quality check 00-general`
