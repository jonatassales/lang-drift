# 0136: ex init generates the CI Gates workflow

**Spec:** `.agents/workflow/specs/0024-ex-init-asks-finops-and-generates-ci-gates.spec.md`

**What to build:** `ex init` asks whether to generate CI Gates, `[y/N]`, or takes `--ci-gates`. On a yes it writes `.github/workflows/langdrift-gates.yml` from an embedded template. The workflow:
- opens with a header that names `ex init` and its version, and says the team owns the file;
- triggers on `pull_request` with `contents: read`;
- sets a timeout and a concurrency group per pull request;
- runs checkout pinned to a SHA, without credentials and with full history;
- runs setup-node pinned to a SHA;
- runs `npx --yes @langdrift/ex@<version> gates report`, with the SHAs passed through `env`.

An existing file is kept byte for byte. An unreleased build refuses to write it. The README's boundary says `init` may write this one file after an explicit yes.

**Blocked by:** 0133 (ex gates report writes the pull request's report), 0134 (ex init asks for the budget and the default gates).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] A yes writes the workflow, which matches a golden copy under testdata for a fixed version, and the result lists `ci-gates` as recorded.
- [x] A no, or no answer, writes nothing outside `.langdrift/`.
- [x] An existing workflow file is kept byte for byte, and the result lists it as kept.
- [x] A `0.0.0-dev` build refuses the workflow with the reason and writes nothing.
- [x] The golden copy pins every action to a SHA with a version comment, sets top-level permissions, and has no expression inside `run`.
- [x] The README's boundary and the init documentation describe the file.
