# 0194: contracts v2: component by path, readiness from the Dockerfile, env by file

**Spec:** .agents/workflow/specs/0032-components-by-dockerfile-path.spec.md

**What to build:** The Microruntimes contracts describe a component by its name and folder, with the Dockerfile name, the port, the readiness kind, an env file and pass-through variable names as optional refinements, and refuse the old inline image, build recipe and readiness path. The fence version becomes 2 and any other version is refused with the re-plan message. The build-tool allowlist and the recipe-refused error leave the contracts; the image-reference check stays for the derived tag. This ticket shares the integration branch with 0195 and 0196: the CLI does not compile against the new type until 0196 lands, and green is promised there.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] A fence with `version: 2` and a component `{ name, path }` parses; `dockerfile` defaults to `Dockerfile`, `port` and `readiness` may be absent.
- [ ] `name` must be a unique DNS label; `dockerfile` has no slash; `readiness.kind` is `tcp` or `http`, and `http` requires a path starting with `/`; `envFile` is a clean relative path inside the checkout; env and pass-through names match an uppercase variable name and do not overlap; each rule has a refusing test with its reason.
- [ ] `image`, `build` and `readinessPath` are refused as unknown fields; `version: 1` is refused with "re-plan the runtime".
- [ ] The contracts README describes the new fence and the readiness derivation.
- [ ] `go test -race ./...` passes in the contracts module.

### Verification

- `pnpm turbo run test typecheck --filter=@langdrift/contracts`
- `platform/cli/drift quality check 00-general 01-filesystem`
