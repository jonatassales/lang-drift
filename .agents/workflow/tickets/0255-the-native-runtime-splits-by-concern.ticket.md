# 0255: the native runtime splits by concern

**Spec:** .agents/workflow/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `lib/cli/src/infrastructure/connected/native.go` splits by concern: verification of the repository, materialization, the snapshot, start, stop, observation, migration with its source lookup and intent files, the build tree and its copy, and the managed-state check, each in its own file, `Native` keeping its methods' names and signatures. One-letter names become words and discarded parses are handled. The rest of the package applies the general rules.

**Blocked by:** 0246, 0254.

**Status:** ready-for-agent

- [ ] `native.go` no longer exists; each concern is a file of `Native`'s methods.
- [ ] The native tests pass at every step of the split.
- [ ] `platform/cli/drift quality check general` prints no finding under `infrastructure/connected`.

## Dev part

**Category:** `hard-part`. `preserveNative` decides which rendered workload objects keep the native runtime's own values; its depth-3 branching is where this ticket's difficulty sits.

**Paths:**
- `lib/cli/src/infrastructure/connected/preserve-native.go`
- `lib/cli/src/infrastructure/connected/preserve-native_test.go`

**Boundary:** `func (native Native) preserveNative(owner project.Project, original microruntime.Descriptor, plan microruntime.Blueprint, rendered map[string][]byte) error`, its signature unchanged and named here, called by the materialization file the agent writes; `Native` and every type in the signature are declared outside the paths.

**Goal:** Flatten `preserveNative` into guards and named steps (one unit per workload object, a service port rewrite, a workload rewrite) with no `else` after an exit and at most two conditional levels, preserving exactly what it keeps today.

- [ ] Every case today's native tests pin passes unchanged, plus one test per named step.
- [ ] No branch nests beyond two levels.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go vet ./... && GOTOOLCHAIN=local go test -race ./...`
- `platform/cli/drift quality check general`
