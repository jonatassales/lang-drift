# 0246: drift measures one unit per Go file, nesting, Go tests and shell functions

**Spec:** .agents/specs/0035-general-rules-compliance-in-packages.spec.md

**What to build:** `drift quality check general` gains: one unit per Go file as an awk program (more than one free function, a free function beside a method, or methods on two receiver types), replacing the per-file count of `^func`; conditional nesting in Go at three or more tab stops; more than one function in a shell file; one-letter names in Go tests except `t`; one-line keyed literals in Go tests outside a table; one `Test` per Go test file; the discarded-error row extended to tests. Each awk detector ships with fixtures under `platform/cli/src/quality/testdata/<detector>/`, and `drift quality detectors` runs every awk program over its fixtures and fails when the reported set differs.

**Blocked by:** 0236, 0245.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] `platform/cli/drift quality detectors` passes, and fails when a fixture's expected set is edited.
- [ ] Every awk detector, including the existing ones, has a fixture set it must report and one it must not; 0236's compound-condition fixtures are the owner's and are run, not written, here.
- [ ] `drift quality check` still exits 0 with findings, and the review-enforced classes are listed.

## Dev part

**Category:** `hard-part`. Reading Go's top-level structure (receivers, free functions, generic receivers, comments and raw strings) without a parser is where this ticket's difficulty sits.

**Paths:**
- `platform/cli/src/quality/go-one-unit.awk`
- `platform/cli/src/quality/testdata/go-one-unit/`

**Boundary:** The program reads the Go files `run-detectors.sh` passes it and prints one `file:line: message` per finding; the agent writes its row in `rule-detectors.tsv` and the `drift quality detectors` runner that feeds it the fixtures.

**Goal:** Report a Go file that declares more than one free function, a free function beside a method, or methods on more than one receiver type; accept a file of one free function, or of one type and its methods, or of plain types and interfaces with no methods.

- [ ] Its fixtures cover a constructor beside its type, a generic receiver, a pointer and a value receiver of one type, a `func` inside a raw string and a comment, and a test file.
- [ ] Run over `lib/cli`, it reports the files the spec counts and none of the files already split.

## Verification

- `platform/cli/drift quality detectors`
- `platform/cli/drift quality check general`
