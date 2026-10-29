# 0118: Pull requests fail on known high or critical advisories

**Spec:** .agents/workflow/specs/0018-search-coverage-and-hardening.spec.md

**What to build:** The pull-request workflow gains a Security job, which:

- runs `pnpm audit --prod --audit-level=high`;
- runs govulncheck, pinned to a released version, over each Go module (extub, ex and contracts);
- follows the workflow's conventions: actions pinned to full SHAs, a read-only token, no persisted credentials, a timeout.

The context file's CI table lists the job.

**Blocked by:** 0107.

**Status:** ready-for-agent

- [ ] The job runs on every pull request into `main`, and fails on a high or critical npm advisory or on a reachable Go vulnerability.
- [ ] Its commands pass on this branch.
- [ ] The CI table in `.agents/context.md` names the job and what it runs.

### Verification

- The job's commands run locally from the repository root; CI on the pull request.
