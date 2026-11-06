# 0318: the runtime.md fence v3, beside v2

**Spec:** .agents/specs/0049-the-core-contracts-carry-the-simulation.spec.md

**What to build:** The Core contracts carry the fence the simulation's plan drafts and its run applies: version 3, with the runtime's id and revision, the bound, the suggestion, and per component its path, stack, role, port, image, share and replicas. A JSON Schema describes it, and validation returns every error, so the plan can print a count. It lands beside the v2 `Configuration`, which the CLI still writes and reads until specs 0055, 0066 and 0062 move it, so the CLI keeps building in the workspace.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

**Alignment:** build. Canonical rule (ADR 0005): what is done beyond the simulation is removed; what the simulation shows and the code lacks becomes a new ticket that builds it.

- [ ] A v3 fence type and its component type exist in `configuration` with exactly the fields of spec 0049's table (`version` 3, `runtimeId`, `revision`, `cpuMillis`, `memoryMiB`, `model`, `effort`, and per component `path`, `stack`, `role`, `port`, `image`, `cpuMillis`, `memoryMiB`, `replicas`), and parsing one `langdrift` block refuses unknown fields.
- [ ] The fixture fence for the simulation's five components holds the Plan card's COMPONENTS rows: `services/signup-api Go api 250m 256 MiB`, `services/notify-api Node.js api 300m 512 MiB`, `platform/nats NATS broker 200m 256 MiB`, `platform/redis Redis cache 200m 256 MiB`, `apps/signup Next.js web 500m 1,024 MiB`, with ports `8080 · 8081 · 4222 · 6379 · 3000`, a bound of `2,000m · 3,072 MiB` and `5 one-replica Deployments`, and validates with 0 errors, so `console › tool call: validate runtime.md draft against the schema · 0 errors` and VALIDATION `schema · 0 errors` have a real count.
- [ ] Validation collects every error: two equal ports plus shares over the bound return exactly two, which is what `console › cross-checking · 5 ports distinct · 5 images resolved · no collisions` and `✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares` rely on.
- [ ] The fence's JSON Schema is published beside the judge's schemas, and the Go validation and the schema are tested against the same sample: the schema states every rule a JSON Schema can state, the same as Validate, and the rules it cannot (unique paths, distinct ports and images, shares within the bound) are Validate's alone, since the module stays on the standard library and runs no schema validator.
- [ ] The v2 `Configuration`, its `Parse` and `Validate` keep their behaviour (the reading of the one `langdrift` block is shared with the v3 parse, with the same refusals), and `lib/cli` builds and passes its tests in the workspace and with `GOWORK=off`.

## Dev part

**Category:** `core-domain`. Which drafts the platform accepts is the contract that keeps a plan from proposing more than the design doc allows; every rule is a promise the Plan card prints.

**Paths:**
- `lib/core/contracts/go/configuration/validate-fence.go`
- `lib/core/contracts/go/configuration/validate-fence_test.go`
- the rules it composes, one function per file: `header-problems.go`, `component-problems.go`, `repeated-problems.go`, `share-problems.go`, with `component-field.go` and `problem-at.go`

**Boundary:** `func (fence Fence) Validate() []Problem`, where `Fence`, `FenceComponent` and `Problem` (the field and the rule it breaks) are declared by the agent in their own files beside it; the plan counts the returned problems.

**Goal:** Return every problem of a v3 fence, never stopping at the first: the version, identifier and revision; the bound within 1 to 64,000 millicores and 1 to 1,048,576 MiB; the effort one of `low`, `medium`, `high` and the model not empty; 1 to 16 components; each path clean, relative, outside `.langdrift` and unique; each stack in `Go`, `Node.js`, `NATS`, `Redis`, `Next.js` and each role in `api`, `broker`, `cache`, `web`; each port valid and distinct; each image a valid reference and distinct; each share positive and their sum within the bound; each `replicas` equal to 1.

- [ ] One case per rule, and a fence breaking two rules returns exactly two problems.
- [ ] The simulation's fixture fence returns none.

**Override:** on 2026-11-05 the owner asked the agent to write this dev part ("You write it"), so the agent wrote `Validate` and its test.
