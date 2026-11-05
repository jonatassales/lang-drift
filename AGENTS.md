# LangDrift

LangDrift is an agentic FinOps platform: one brand, one domain and one console. This repository, code name **Tao**, builds the next LangDrift. This file holds the ideas an agent must carry into any task here; the mechanics live in [`.agents/context.guide.md`](.agents/context.guide.md).

## The idea

### One platform, one domain

Everything is served under `https://langdrift.md` with a locale prefix (`/en/`, `/pt-br/`, `/zh-hant/`, `/ja/`, `/vi/`). The public site presents the platform on one page and routes to the FinOps and Microruntimes pages. Moving between pages, on the site and in the console, is routing on that one domain, never a redirect to another domain or subdomain.

### The product

The product carries the brand in its name: **LangDrift FinOps**, with **Microruntimes** as its engine.

| Part | The idea it serves |
| --- | --- |
| **FinOps** | Make agent-assisted delivery accountable to a team budget. It connects product intent, delivery context, execution evidence and economic outcome in one inspectable chain, on three pillars: FinOps as the economic purpose, the agentic harness as the delivery workflow and its evidence, and microruntimes as bounded, consumer-owned execution around software that already exists. |
| **Microruntimes** | The engine: plan, run and record a bounded execution around a team's own software. Its Go CLI, `dft`, lives in `lib/cli` beside the Microruntimes contracts in `lib/core`; every plan with Microruntimes includes FinOps. |

FinOps grew out of the Exphere project. Its value comes from data the platform produces itself.

**Pivot, 2026-10-28.** The owner pivoted the project from a platform of products to FinOps, with Microruntimes as its engine. Vision is no longer developed as a product, because it would need ever more of the client's data and producing artifacts for it would break FinOps. It may return as a connector, an MCP to tools such as GitHub, Linear or Obsidian, that brings their data into the console. Spec 0010 records the decision.

### One console

The LangDrift console is where a team manages FinOps: projects registered from a local folder, their microruntimes, executions, history and secrets, supervised by one connected `dft` agent per project, with FinOps readings and recommendations that stay advisory. It lives here as `apps/console`, imported from Exphere's console (`jonatassales/exphere`), and reads `apps/microservices/console-stubs`, a local stub API that seeds one project, LangDrift Studio, with its runtimes and their FinOps readings. Its sidebar switches between two areas, Microruntimes and FinOps, and the URL decides which. The console keeps requested state apart from observed state and labels illustrative or stub data as such. The account and the shell are planned.

### Principles the platform holds

- **The platform owns boundaries and mechanisms; the consumer owns the application and its domain choices.** A platform team may require observability, evidence, security and bounded execution without choosing every library.
- **Recommendation is not authorization.** An agent may estimate, recommend and execute permitted work. A deterministic control decides whether a paid or risky action is allowed. A model's suggestion never widens scope, funding or access.
- **The authorized budget belongs to the team**, not to a developer or a ticket. Allocations to work are subordinate to that envelope.
- **Unknown is never zero.** Missing or stale evidence is shown as missing, not as zero spend or a confirmed stop. Desired state and observed state stay distinct.
- **Evidence over assertion.** Every movement, recommendation and lifecycle operation keeps a token-free record with provenance. A claim of savings, delivery or readiness needs the evidence behind it.
- **Meaning is fixed everywhere.** Status, classification and drift colors never take a brand color, so a classification reads the same on every surface.

## Where things stand

A pnpm and Turborepo monorepo. `apps/website` is the platform site, one page with the FinOps and Microruntimes pages beside it; it started as langdrift.md, ported from the earlier `langdrift` monorepo under spec 0003. `apps/console` is the console and `apps/microservices/console-stubs` its local stub API. `packages/tokens`, `packages/web`, `lib/cli` (the LangDrift CLI) and `lib/core` (the Microruntimes contracts) exist; `shell`, `account` and `runtime` are planned. Each app and package has a README with its own contract.

## How we work

- **Rules first.** [`.claude/rules/`](.claude/rules/) holds the rules that apply to every task and are loaded in every session. Read them before touching code.
- **The owner writes application code.** Agents edit app code only when the task asks for it. Repetitive work may be delegated when the task names its scope.
- **A request to build or change application code enters the spec-driven workflow** instead of producing code: `grilling` (plus `domain-modeling` when terms or decisions settle), then `to-spec`, `to-tickets`, `implement`, `code-review`. [`.agents/policies/human-in-the-loop.policy.md`](.agents/policies/human-in-the-loop.policy.md) says where the owner takes part and what runs on its own.
- **Open decisions stay open.** Pricing, limits, auth provider and formulas are the owner's to close. Rules are added only when the owner asks, and the harness stays a draft until the owner promotes a piece of it to a rule.
- **Nothing carries an AI signature**: no `Co-Authored-By`, no "Generated with", in commits, pull requests or contributions. Inside the workflow, commit, push and pull request follow the policy above; outside it, each needs explicit authorization every time. Merge and release always do.
- **Secrets stay out of reach.** Do not read or edit values in `.env` files, and keep credentials out of code, logs and evidence.

## Where to look

| When you need | Read |
| --- | --- |
| Layout, toolchain, commands, the `drift` CLI, design tokens, CI and Git mechanics | [`.agents/context.guide.md`](.agents/context.guide.md) |
| The project vocabulary | [`.agents/context.guide.md`](.agents/context.guide.md), `## Glossary` |
| Specs, tickets, ADRs and research notes | `.agents/workflow/` |
| The website's environment contract and sign-up pause | [`apps/website/README.md`](apps/website/README.md) |
| The LangDrift CLI and the Microruntimes contracts | [`lib/cli/README.md`](lib/cli/README.md), [`lib/core/contracts/go/README.md`](lib/core/contracts/go/README.md) |

## Commands

Run from the repository root. `pnpm website` serves the site on port 3000, `pnpm console` the console on 3001 with its stub API on 3401, so both run at once. [`.agents/context.guide.md`](.agents/context.guide.md) has the full table and what each command touches.

```bash
pnpm install
```

```bash
pnpm website
```

```bash
pnpm console
```

```bash
pnpm turbo run check
```

```bash
pnpm turbo run typecheck
```

```bash
pnpm turbo run test
```

# This is NOT the Turborepo you know

Turborepo configuration, task behavior, and CLI commands can vary between installed versions and may differ from your training data. Resolve the `turbo` package from this file's directory or relevant workspace; in monorepos, it may not be visible from the repository root. For example, run `node -p "require.resolve('turbo/package.json')"` from a workspace that depends on `turbo`.

Read `docs/README.md` inside that installed package first, then read the relevant pages from its `docs/` directory before changing Turborepo configuration or commands. Heed deprecation notices. These bundled docs match the installed package version and are available without network access.

This block is written and re-added by `turbo` before repository-scoped commands when an AI agent is detected. In the Turborepo source repository, its template is defined in `crates/turborepo-cli/src/cli/agent_guidance.rs`. Removing the managed block while updates are enabled means a later qualifying invocation will add it again. Set `"agentGuidance": false` in the root `turbo.json` or `turbo.jsonc` to opt out; this does not remove an existing block. Keep the block committed with your work to avoid an uncommitted change on the next agent invocation.
<!-- END:turborepo-agent-rules -->
