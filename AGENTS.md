# LangDrift

LangDrift is an AI harness web OS for the agentic journey: one brand, one domain and a set of platform products that share account, plan and shell. This repository, code name **Tao**, builds the next LangDrift. This file holds the ideas an agent must carry into any task here; the mechanics live in [`.agents/context.md`](.agents/context.md).

## The idea

### One platform, one domain

LangDrift is a platform of products, not one product. The public site presents the platform and routes to each product's own page, the way a `Platform` menu does on a platform vendor's site. Everything is served under `https://langdrift.md` with a locale prefix (`/en/`, `/pt-br/`, `/zh-hant/`, `/ja/`, `/vi/`). Moving between products, on the site and in the console's product switcher, is routing on that one domain, never a redirect to another domain or subdomain.

Today's langdrift.md is a single-product site. Once the platform site is ready it takes over the domain, and today's content becomes Vision's product page.

### The products

Products carry the brand in their name: **LangDrift Vision** and **LangDrift FinOps**.

| Product | The idea it serves |
| --- | --- |
| **Vision** | Keep the why. A product moves away from the vision it set out with; Vision measures that movement as a **Product Vision** score and keeps the context of every move: what changed, who moved it, why, whether a decision was recorded. Movement is classified as intentional evolution, unexplained drift or under review. Drift is not a verdict: it can reveal a gap, a trade-off or a better direction. The goal is shared understanding, not blame. |
| **FinOps** | Make agent-assisted delivery accountable to a team budget. It connects product intent, delivery context, execution evidence and economic outcome in one inspectable chain, on three pillars: FinOps as the economic purpose, the agentic harness as the delivery workflow and its evidence, and microruntimes as bounded, consumer-owned execution around software that already exists. |

Vision's metric keeps its English name, "Product Vision", in every locale. FinOps grew out of the Exphere project; its Go CLI, `ex`, is already here as `packages/ex`.

### One console

Each product has a management surface today, in its upstream repository. LangDrift's dashboard (`jonatassales/langdrift`, `apps/dashboard`) is Vision's: Overview, Vision, Evolution and Evidence workspaces, with attention, briefs, targets, revisions, drifts, decisions, relations and sources beneath them. Exphere's console (`jonatassales/exphere`, `apps/console`) is FinOps's: projects registered from a local folder, their microruntimes, executions, history and secrets, supervised by one connected `ex` agent per project, with FinOps readings and recommendations that stay advisory.

They become **one LangDrift console** on langdrift.md: one account, one shell, a product switcher between Vision and FinOps, and under it each product's own workspaces, sidebar and content. The dashboard's workspace switcher is the seed of that product switcher. Both surfaces already keep requested state apart from observed state and label illustrative or stub data as such; the unified console keeps both habits. Neither console is in this repository yet.

### Principles the products share

- **The platform owns boundaries and mechanisms; the consumer owns the application and its domain choices.** A platform team may require observability, evidence, security and bounded execution without choosing every library.
- **Recommendation is not authorization.** An agent may estimate, recommend and execute permitted work. A deterministic control decides whether a paid or risky action is allowed. A model's suggestion never widens scope, funding or access.
- **The authorized budget belongs to the team**, not to a developer or a ticket. Allocations to work are subordinate to that envelope.
- **Unknown is never zero.** Missing or stale evidence is shown as missing, not as zero spend or a confirmed stop. Desired state and observed state stay distinct.
- **Evidence over assertion.** Every movement, recommendation and lifecycle operation keeps a token-free record with provenance. A claim of savings, delivery or readiness needs the evidence behind it.
- **Meaning is fixed across products.** Status, classification and drift colors never take a brand color, so a classification reads the same in every product.

## Where things stand

A pnpm and Turborepo monorepo. `apps/website` is where the platform site is built; it started as today's langdrift.md, ported from the earlier `langdrift` monorepo under spec 0003. Work happens in `apps/website`. `packages/tokens`, `packages/web` and `packages/ex` exist; `shell`, `account`, `runtime` and `contracts` are planned. Each app and package has a README with its own contract.

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
| Layout, toolchain, commands, the `drift` CLI, design tokens, CI and Git mechanics | [`.agents/context.md`](.agents/context.md) |
| The project vocabulary | [`.agents/context.md`](.agents/context.md), `## Glossary` |
| Specs, tickets, ADRs and research notes | `.agents/workflow/` |
| The website's environment contract and sign-up pause | [`apps/website/README.md`](apps/website/README.md) |
| The `ex` CLI and its lifecycle contracts | [`packages/ex/README.md`](packages/ex/README.md) |

## Commands

Run from the repository root. [`.agents/context.md`](.agents/context.md) has the full table and what each command touches.

```bash
pnpm install
```

```bash
pnpm website
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
