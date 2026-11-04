# 0221: `ex agent init` and the LangGraph scaffold leave the CLI

**Spec:** .agents/workflow/specs/0045-the-judge-as-a-platform-service.spec.md

**What to build:** A consumer no longer finds a command that generates an agent nobody runs. `ex agent init`, the agent scaffolder port, the LangGraph templates and their tests are removed; `ex agent --help` lists `confirm`, `connect`, `migrate`, `plan`, `project`, `reindex`, `run`, `status` and `supervise`; the CLI README and the context document stop naming the scaffold; `.langdrift/agents/` is neither written nor read. Everything else behaves as before.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] `ex agent init` is unknown to the CLI, and no template under the scaffold folder mentions LangGraph or LangChain.
- [ ] `ex init` on a fresh project writes no `agents/` folder, and a project holding an old `.langdrift/agents/` is planned and run without reading it.
- [ ] The CLI's suite passes under the race detector; the composed CLI test's expected tree lists no scaffold file.
- [ ] The CLI README and the context document name no agent scaffold.

## Dev part

**Override:** No part qualifies: the ticket deletes a retired feature and its documentation; nothing is written by hand (policy, Dev part: a problem with its own goal and acceptance criteria; a deletion has none). 2026-11-04.

## Verification

- `cd lib/cli && GOTOOLCHAIN=local go test -race ./...`
- `lib/cli/dist/ex agent --help`
- `git grep -n -i 'langgraph\|langchain\|agent init' -- lib README.md .agents/context.md`
