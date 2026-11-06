# 0133: ex gates report writes the pull request's report

**Spec:** `.agents/specs/0024-ex-init-asks-finops-and-generates-ci-gates.spec.md`

**What to build:** `ex gates report` turns a pull request's commit range into a Markdown report, whether or not a project exists. It reads `BASE_SHA` and `HEAD_SHA` and lists the range's commits through git. The report shows the range, its commit count and each commit's short SHA and subject. Its FinOps section marks the budget, the spend and the ledger as unknown, because CI has no committed FinOps source. The report goes to the file `GITHUB_STEP_SUMMARY` names, or to standard output when that variable is unset. It never opens `.langdrift/` and posts nothing.

**Blocked by:** None (can start immediately).

**Status:** completed. Its acceptance criteria were all checked when ADR 0005 (Perfect Simulation Alignment) reset the backlog on 2026-11-05.

- [x] In a git repository without `.langdrift/`, `ex gates report` exits 0 and prints the range, the commit count and each commit's short SHA and subject.
- [x] Every FinOps reading is marked unknown, with the reason that no committed FinOps source exists. None reads as zero.
- [x] With `GITHUB_STEP_SUMMARY` set, the report is appended to that file and standard output stays empty.
- [x] A missing `BASE_SHA` or `HEAD_SHA` is a usage error that names the variable.
- [x] A SHA that git cannot resolve fails with the external-failure status and the reason.
- [x] `ex --help` lists the `gates` group, and the README documents `ex gates report`.
