<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="/platform/docs/assets/images/langdrift-banner-dark.png">
    <img src="/platform/docs/assets/images/langdrift-banner-light.png" alt="LangDrift: Agentic FinOps Platform" width="820">
  </picture>
</p>

<p align="center">
  <strong>Agent-assisted delivery, accountable to your team's budget.</strong>
</p>

<p align="center">
  <strong>English</strong> ·
  <a href="/platform/docs/README.pt-br.md">Português (Brasil)</a> ·
  <a href="/platform/docs/README.zh-hant.md">繁體中文</a> ·
  <a href="/platform/docs/README.ja.md">日本語</a> ·
  <a href="/platform/docs/README.vi.md">Tiếng Việt</a>
</p>

<p align="center">
  <a href="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml"><img alt="Pull request checks" src="https://github.com/jonatassales/lang-drift/actions/workflows/pull-request.yml/badge.svg?event=pull_request" height="20"></a>
  <a href="/packages/ex/README.md"><img alt="Go 1.26 or newer" src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/package.json"><img alt="Node.js 24" src="https://img.shields.io/badge/Node.js-24-5FA04E?logo=nodedotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/apps/website/README.md"><img alt="Next.js 16" src="https://img.shields.io/badge/Next.js-16-F5F5F5?logo=nextdotjs&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
  <a href="/.agents/context.md"><img alt="pnpm 12" src="https://img.shields.io/badge/pnpm-12-F69220?logo=pnpm&amp;logoColor=white&amp;labelColor=18181B" height="20"></a>
</p>

<p align="center">
  <a href="https://langdrift.md">Website</a> ·
  <a href="#the-product">The product</a> ·
  <a href="#the-envelope">The envelope</a> ·
  <a href="#get-started-with-ex">Get started with <code>ex</code></a> ·
  <a href="#project-status">Project status</a>
</p>

## The product

**LangDrift FinOps** makes agent-assisted delivery accountable to a team budget. It connects what a ticket asks for, the work an agent does, the evidence that work leaves and what it cost, in one chain a team can inspect.

**Microruntimes** is its engine. Its CLI, `ex`, reads a ticket and the documents it links, plans a bounded local runtime around the software you already have, and drives the delivery through your own coding agent, with every paid call reserved against the team's envelope before it runs.

A model call that times out, a retry that costs more than planned, a pull request whose checks never report: delivery leaves the plan in ordinary ways. LangDrift's aim is that none of them becomes spend nobody approved or a number nobody can explain.

### What a team gets

| Capability | Value for the team |
| --- | --- |
| **Plans from tickets** | `ex agent plan` reads the ticket and only the documents it links, retrieves comparable earlier runs and writes a runtime plan sized inside the team's envelope. |
| **Reservations before spend** | Every paid call, the planner's and the coding agent's, is estimated at dated list prices and reserved against the per-run cap, the ticket's allocation and what is left of the period. |
| **Delivery through your harness** | Your workflow stages and skills run through the coding agent you configure. People confirm where the team says they take part, and the merge always stays the owner's. |
| **Runtimes around existing software** | Your components come up on a bounded local kind, Terraform and Kubernetes lifecycle, with no LangDrift dependency in your code. |
| **Memory with evidence** | Every run leaves token-free evidence, a ledger and a memory that the next plan retrieves. |

LangDrift has its own word for a pass of a run: a microruntime is the car, a run is the race, and a **lap** is one pass of it, either an attempt paused at a gate or the execution that finished.

## The envelope

The authorized budget belongs to the team, not to a developer or a ticket. Tickets receive allocations inside it, and each run has a cap.

Before a lap opens, `ex` estimates the run at list prices and adds the team's reserve margin. A reservation that fits the per-run cap, the ticket's allocation and the period's remainder opens the lap; one that does not pauses it with nothing paid. While the lap plays, a gate runs before every paid call. When it ends, the unused reservation returns to the budget.

| Gate | What it checks | When it does not pass |
| --- | --- | --- |
| **FG-01** | A dated price source covers the paid work. | Blocks the lap. |
| **FG-02** | The ticket has an allocation. | Blocks the lap. |
| **FG-03** | The estimate plus margin fits the per-run cap, the allocation and the period's remainder. | Pauses before anything is paid. |
| **FG-04** | The lap's observed spend has not reached its reservation. | Pauses the paid calls. |
| **FG-05** | The owner acknowledged an earlier lap's overspend. | Requires the owner, with `--acknowledge-overspend`. |
| **FG-06** | The provider has billed the run. | Keeps the figures provisional. |

**Unknown is never zero.** A call that reports no usage, a cleanup nobody could confirm or a bill that has not arrived stays unknown, with its reason, and keeps its reservation held until it is reconciled.

## Your stack, your decisions

LangDrift owns boundaries and mechanisms; your team owns the application and its domain choices.

- **A small footprint.** In your project, `ex` writes only `.langdrift/`. It reads the artifacts a ticket links, never outside the project, and refuses secrets and symlinks.
- **Your delivery, observed.** `ex` never pushes, opens a pull request or merges. Delivery belongs to your harness; `ex` observes the branch, the pull request and its checks through `git` and `gh`.
- **Recommendation is not authorization.** The planner suggests a model and an effort; running the runtime is what authorizes the plan. A model's suggestion never widens scope, funding or access.
- **Credentials stay in the environment.** The planner reads `ANTHROPIC_API_KEY` at the call, and the coding agent receives a minimal environment plus only the variables you list. `ex` never takes a credential as a command argument.

Every input that can block or allow a run is recorded with who set it, when and through what. The JSON files under `.langdrift/` mirror those records for reading; editing a mirror changes nothing a gate reads.

## Project status

**LangDrift is in development.** The website is live at [langdrift.md](https://langdrift.md) in five languages. The console reads `extub`, a stub API that seeds one project, LangDrift Studio, and labels its FinOps readings as illustrative. You build `ex` from source; no binary is published.

`ex agent plan` and `ex agent run` work locally inside a project, and the money gates run inside `ex` against the project's own ledger; a team ledger shared through a server is not built yet. The connected agent (`ex agent connect`, `project`, `supervise` and `migrate`) needs a server this repository does not have yet. The test suites drive the CLI with scripted models, coding agents and lifecycles; they do not establish a live cluster start, real spend or production readiness. Sign-up on the website is paused.

| Path | What it holds |
| --- | --- |
| [`apps/website`](/apps/website/README.md) | The platform site at langdrift.md: home, FinOps, Microruntimes, pricing and early access. |
| [`apps/console`](/apps/console/README.md) | The console where a team manages FinOps and its microruntimes. |
| [`apps/extub`](/apps/extub/README.md) | The local Go stub API the console reads. |
| [`packages/ex`](/packages/ex/README.md) | The `ex` CLI, in Go. |
| [`packages/contracts`](/packages/contracts/README.md) | The Microruntimes contracts: the `runtime.md` fence and the agent protocol. |
| [`packages/tokens`](/packages/tokens) · [`packages/web`](/packages/web) | Design tokens and the shared UI. |

## Get started with `ex`

You need macOS or Linux and Go 1.26 or newer; that is enough to plan. The planner's model judge reads `ANTHROPIC_API_KEY`. Running a lap also needs Docker, kind, kubectl and Terraform (`ex doctor` checks them and installs nothing), the coding agent you configure, such as Claude Code, and `gh` to observe pull requests.

### 1. Build the CLI

```sh
git clone https://github.com/jonatassales/lang-drift.git
cd lang-drift/packages/ex
GOTOOLCHAIN=local go build -o dist/ex ./src
./dist/ex version
```

`version` prints `ex 0.0.0-dev`. Your system may already have an `ex`, the Vi line editor, so call this one by its path. The steps below write `ex` for `packages/ex/dist/ex`; for this shell only, `alias ex="$PWD/dist/ex"` does that.

### 2. Initialize a project

To try it on a ready-made project, copy the one the test suite plans: a ticket, the spec, research, ADR and guardrail it links, and two app manifests.

```sh
cp -R src/testdata/langdrift/project /tmp/ex-demo
cd /tmp/ex-demo
ex init
```

`ex init` creates the project identity in `.langdrift/config.json` and the execution rules in `.langdrift/microruntimes/AGENTS.md`, and leaves existing files unchanged. In your own repository, add `/.langdrift/` to `.gitignore`: it is machine-local state that `ex` writes, and nothing your team authors lives in it.

### 3. Record the harness

Save the [example input files](#example-input-files) in the project, then record the delivery harness: the hardware envelope, the workflow stages, the default model and effort, and the coding agent that plays the stages.

```sh
ex harness envelope set --per-execution-cpu-millis 4000 --per-execution-memory-mib 6144 \
  --default-cpu-millis 2000 --default-memory-mib 3072 --not-runtime-field gpu
ex harness workflow set --file workflow.json
ex harness model set claude-sonnet-5-5 --effort medium
ex harness driver set --file driver.json
```

### 4. Record the money rules

This budget gives the team $500.00 for October 2026, caps each run at $25.00, reserves a 20% margin over every estimate and allocates $25.00 to ticket 0070. Amounts are in cents.

```sh
ex finops budget set --period 2026-10 --period-cents 50000 --per-run-cents 2500 \
  --reserve-margin-percent 20 --allocate 0070-console-microruntime=2500
ex finops gates set --default
ex finops prices set --file prices.json --slug list
ex finops forecasts set --file forecasts.json
```

Each command prints the record it wrote, and every input command takes `--by` to name who set it.

### 5. Plan from a ticket

Name the components the runtime should build and run, then let `ex` read the ticket. `--no-judge` plans without a model call, so nothing is paid; without it, and with `ANTHROPIC_API_KEY` set, the planner judges the sources through a model under the same gates as a lap, and its cost lands on the ledger's planning line.

```sh
ex agent plan .agents/workflow/tickets/0070-console-microruntime.ticket.md --no-judge \
  --component "path=apps/console,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=docker,args=build -t langdrift-console:local ."
```

```text
ex agent plan 0070-console-microruntime.ticket.md
# plan only · nothing is built or started · the planner's calls are paid and gated
agent › tool call: read 0070-console-microruntime.ticket.md
agent › needs-triage · links 0010-exphere-and-the-langdrift-folder.spec.md · ex-microruntimes-capability.research.md · 0002-console-input-is-untrusted.adr.md · 0001-finops.gr.md · package.json · go.mod
judge      unavailable · --no-judge; verdicts, facts and the suggestion are the inputs' own
# analysis · reading the linked documents, then the project's episodes
agent › tool call: read 0010-exphere-and-the-langdrift-folder.spec.md · read
agent › tool call: read ex-microruntimes-capability.research.md · read
agent › tool call: read 0002-console-input-is-untrusted.adr.md · read
agent › tool call: read 0001-finops.gr.md · read
agent › tool call: read package.json · read
agent › tool call: read go.mod · read
✓ read       7 documents · nothing outside the ticket's links
reads      .langdrift/harness/guardrails/0001-finops.gr.md · FIN-01 · FIN-02 · FIN-03 · FIN-04 · FIN-05 · FIN-06 · FIN-07 · FIN-08
bounded    cpuMillis 2000 · memoryMiB 3072 · set by owner in the CLI · record envelope-b228
budget     $25.00 per run for model calls · set by owner in the CLI
agent › RAG over project episodes · a model and effort for this ticket
agent › embedding the ticket intent · 1 query vector
agent › vector search · project episodes · top 6 by cosine similarity
agent › 0 candidates · no episode indexed yet
agent › dropped 0 below 0.55 · kept 0
agent › 0 episodes cited · none with measured resources
agent › reasoning · no cited episode measured cpu or memory · sizing follows the inputs, not history
✓ suggests   claude-sonnet-5-5 · effort medium
# plan · resolving the components, then sizing them inside the envelope
agent › resolving 1 components named in the CLI · paths, images, ports
agent › tool call: inspect apps/console · package.json · :3001
✓ apps/console          Node.js+Next.js · web 500m · 1024 MiB
agent › cross-checking · 1 ports distinct · 1 images resolved · no collisions
✓ size       500m of 2000m · 1024 of 3072 MiB · 1 one-replica Deployments · proposed shares
limits     gpu and disk are not runtime.md fields · none requested
estimate   runtime has no price source · cpu and memory stay in millicores and MiB
agent › tool call: validate runtime.md draft against the schema · 0 errors
✓ plan       runtime.md draft, revision 1 · proposed · authorize with ex agent run
# next: once the team authorizes, ex agent run applies the recorded intent
```

The plan lands in `.langdrift/microruntimes/runtimes/console-microruntime/runtime.md`. Its `langdrift` fence is what the runtime executes; a second fence, `langdrift-plan`, holds the analysis. Prose in the file is kept as context and never executed.

```json
{
  "llmBudgetUsd": 25,
  "runtimeId": "console-microruntime",
  "revision": 1,
  "cpuMillis": 2000,
  "memoryMiB": 3072,
  "model": "claude-sonnet-5-5",
  "effort": "medium",
  "components": [
    {
      "path": "apps/console",
      "image": "langdrift-console:local",
      "port": 3001,
      "readinessPath": "/icon.svg",
      "build": { "executable": "docker", "args": ["build", "-t", "langdrift-console:local", "."] }
    }
  ]
}
```

This excerpt leaves out the fence's version, name and runtime and project identities, which `ex` manages. Planning the same ticket again, while the runtime is stopped, raises the revision.

### 6. Run the plan

Running the runtime authorizes its plan. A lap builds images, starts the project's local cluster and makes paid calls through your coding agent, so run `ex doctor` first.

```sh
ex agent run console-microruntime
```

1. **The gate.** The run prints the intent, the estimate at list prices and the per-run budget, then reserves. A reservation inside every bound opens the lap; one above them pauses with nothing paid.
2. **The stages.** Each stage plays in order through the coding agent and prints a one-line summary. A stage where a person takes part waits for the owner: answer at the terminal or with `--confirm`, or later with `ex agent confirm console-microruntime planning`. The `owner-only` stage never runs.
3. **The runtime.** The components come up from the fence on the project's kind cluster, and the run prints when they are ready.
4. **The delivery.** The run observes the branch, the pull request and its checks; `--checks-wait` bounds how long it follows them.
5. **The settlement.** The runtime goes down, the unused reservation returns to the budget, the memory `run.md` is written and the episode index is rebuilt. The run says whether it is ready for review or to merge, and what comes next.

Every change lands in the lap file, `run.json` and the ledger.

### 7. Follow the run

```sh
ex agent status console-microruntime
ex agent reindex
ex doctor
```

`status` shows the run record and the open lap, `reindex` rebuilds the episode index from every runtime's memory, and `doctor` checks the native tools. With `--json`, each session line is one NDJSON event, `{at, session, kind, id, text, card, field, data}`, so a console or a terminal UI can render the same session.

| Exit code | Meaning |
| --- | --- |
| `0` | Success, or a run waiting for a person. |
| `1` | The operation failed: a skill, a hook, a build or a start. |
| `2` | Invalid usage. |
| `3` | Refused by a guard or policy, such as a money gate. |
| `4` | A prerequisite is missing: a tool, an input, a cluster or a runtime record. |
| `130` | Interrupted. |

The full command reference, the `.langdrift/` layout and the recovery rules are in the [`ex` README](/packages/ex/README.md).

### Example input files

<details>
<summary><code>workflow.json</code>: the stages, their skills and where a person takes part</summary>

```json
{
  "version": 1,
  "stages": [
    { "id": "planning", "label": "Planning", "skills": ["grilling"], "human": "required" },
    { "id": "spec", "label": "Spec", "skills": ["to-spec"], "human": "optional" },
    { "id": "implement", "label": "Implement", "skills": ["implement", "tdd", "code-review"], "human": "optional" },
    { "id": "review", "label": "Review", "skills": [], "human": "optional" },
    { "id": "merge", "label": "Merge", "skills": [], "human": "owner-only" }
  ],
  "defaults": { "humanAt": ["planning", "review", "merge"] },
  "efforts": ["low", "medium", "high"]
}
```

A stage may also list `hooks`: committed scripts it runs first, through `sh`, inside the project.

</details>

<details>
<summary><code>driver.json</code>: the coding agent that plays the stages</summary>

```json
{
  "executable": "claude",
  "args": ["-p", "{{prompt}}", "--output-format", "json"],
  "passEnv": ["ANTHROPIC_API_KEY"],
  "timeoutSeconds": 1800,
  "usage": { "inputTokens": "usage.input_tokens", "outputTokens": "usage.output_tokens" },
  "summary": "result",
  "unmeasured": "pause"
}
```

The executable runs from the project root, never through a shell. `usage` names where the output reports tokens; a call that reports none leaves its spend unknown, and `"unmeasured": "pause"` makes FG-04 pause the next paid call.

</details>

<details>
<summary><code>prices.json</code>: the dated list prices behind every estimate</summary>

```json
{
  "version": 1,
  "illustrative": true,
  "note": "Replace with your provider's current list prices.",
  "readOn": "2026-10-06",
  "currency": "USD",
  "unit": "cents-per-million-tokens",
  "cache": "uncached",
  "models": [
    { "id": "claude-sonnet-5-5", "name": "Claude Sonnet 5.5", "inputCents": 200, "outputCents": 1000 }
  ]
}
```

These figures are illustrative and exercise the format; they are not a verified price list.

</details>

<details>
<summary><code>forecasts.json</code>: the tokens each effort is expected to use</summary>

```json
{
  "version": 1,
  "illustrative": true,
  "unit": "thousands-of-tokens",
  "byEffort": {
    "low": { "input": 2000, "output": 120 },
    "medium": { "input": 3500, "output": 250 },
    "high": { "input": 6000, "output": 500 }
  },
  "planning": { "input": 35, "output": 5 }
}
```

`planning` covers the planner's own calls. Forecasts are used when no earlier run of the ticket measured its tokens.

</details>

### Explore the website and the console

With Node.js 24 and pnpm 12, from the repository root:

```sh
pnpm install
pnpm website
```

The site runs on `http://localhost:3000`. `pnpm console` starts the console on `http://127.0.0.1:3001` with `extub` beside it on `127.0.0.1:3401`, which needs Go. Every number the console shows comes from the stub.

## The work behind LangDrift

Created by [Jonatas Sales](https://github.com/jonatassales), LangDrift brings together FinOps practice, agentic delivery and platform engineering. It grew out of the Exphere project and is built as a pnpm and Turborepo monorepo: the website and the console in TypeScript with Next.js, React and Tailwind CSS, and `ex` and the Microruntimes contracts in Go, with the official Anthropic Go SDK for the planner.

The repository is built spec-first: every change starts from a spec and its tickets under [`.agents/workflow/`](/.agents/workflow), the same kind of ticket `ex` plans. Feedback is welcome through [GitHub issues](https://github.com/jonatassales/lang-drift/issues).

## References

- [AI for FinOps: Agentic Use Cases in FinOps](https://www.finops.org/insights/ai-for-finops-agentic-use-cases/), Jonathan Morley, FinOps Foundation, March 27, 2026. Practitioners describe agents that reconcile spend against budgets, find waste and enforce cost policy earlier in the development lifecycle; that last theme, shifting FinOps left into delivery, is where LangDrift works.
- [Harness engineering for coding agent users](https://martinfowler.com/articles/harness-engineering.html), Birgitta Böckeler, martinfowler.com, April 2, 2026. It describes the outer harness a team builds around a coding agent: guides that steer the agent before it acts and sensors that let it correct itself. `ex harness` records that harness for a project, and the FinOps gates bound what it may spend.

## Acknowledgements

Thanks to [Matt Pocock](https://x.com/mattpocockuk). The agent skills this repository works with, from grilling to to-spec, to-tickets, implement, tdd and code-review, are adapted from [mattpocock/skills](https://github.com/mattpocock/skills) under the MIT License, and the example workflow above runs several of them.
