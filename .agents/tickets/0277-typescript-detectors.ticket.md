# 0277: drift measures the TypeScript rules

**Spec:** .agents/specs/0039-typescript-rules-compliance.spec.md

**What to build:** `drift quality check` gains, under `04-typescript`: a translator-name awk program; a deep-relative-import pattern for `./dir/...`; an inline-object-type pattern for a generic argument, an `as` and a return annotation; a file-suffix awk program over tracked paths. The await-spacing and import-group detectors move to the `all` scope so tests are measured, and the `const t =` row is replaced.

**Blocked by:** 0246.

**Status:** wontfix. The owner closed it on 2026-11-05: the work goes to the published packages, @langdrift/cli and @langdrift/core, until the runtime runs.

- [ ] Each new detector prints exactly the known violations of its class in the tree at its commit.
- [ ] `platform/cli/drift quality detectors` covers the two new awk programs.

## Dev part

**Category:** `hard-part`. Deciding a file's responsibility from its path against two suffix tables and the framework's conventions is where this ticket's difficulty sits.

**Paths:**
- `platform/cli/src/quality/file-suffix.awk`
- `platform/cli/src/quality/testdata/file-suffix/`

**Boundary:** The program reads the tracked paths `run-detectors.sh` passes it and prints one `path: message` per finding; the agent writes its row and the other detectors.

**Goal:** Report a `.ts` file under an app or package `src/` whose suffix is not in the TypeScript table, the Nest table or the framework's file conventions, and is not a test, an `index.ts` barrel, a configuration file or a declaration file.

- [ ] Its fixtures cover every table suffix, each Next and Nest framework file, a test, a barrel, a config and a declaration file, and an unsuffixed module.
- [ ] Run over the tree, it reports the files the spec's suffix tickets rename and nothing else.

## Verification

- `platform/cli/drift quality detectors`
- `platform/cli/drift quality check 04-typescript`
