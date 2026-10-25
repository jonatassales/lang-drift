# LangDrift

AI harness web OS for your agentic journey. LangDrift is a suite of apps under one account: each app is independent, with its own URL, and all apps share account, plan and shell.

## General rules

`.claude/rules/` holds the rules that always apply, read first in every session. `00-general.rule.md` covers code and tooling, `01-filesystem.rule.md` naming and paths, `02-frontend.rule.md` styling and the CSS exception, `03-tailwind.rule.md` what Tailwind may do, `04-code.rule.md` external input and control flow, `05-react.rule.md` components, `06-nextjs.rule.md` the App Router tree, `07-jotai.rule.md` shared state, `08-localization.rule.md` copy and locales, and `09-seo.rule.md` metadata.

## Status

Monorepo with three packages, `@langdrift/tokens`, `@langdrift/web` and `@langdrift/ex` (the Go agent-harness CLI; see `packages/ex/README.md`), and one app, `apps/website`: the LangDrift website (see `apps/website/README.md`), copied from the earlier `langdrift` monorepo under spec 0003.

## Toolchain and commands

- pnpm is the only package manager. `devEngines.packageManager` in the root `package.json` makes other package managers fail, and `engineStrict` in `pnpm-workspace.yaml` enforces the Node range in `engines`.
- Dependency install scripts run only when `allowBuilds` in `pnpm-workspace.yaml` allows them. Prebuilt binaries ship as optional dependencies, so the entries there are denials; approve a new one only when its build is required.
- Work that spans the monorepo runs through Turborepo. Do not add root scripts that wrap `turbo run`.
- A package's scripts call its tool directly (for example `tsc --noEmit`). Turborepo tasks run those scripts by name.
- Do not add shell or JavaScript wrappers for what Turborepo or pnpm already do.

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm install` | Installs dependencies and the Git hooks, then links the `drift` launcher. Fails under another package manager or an unsupported Node version. |
| `pnpm turbo run build` | Builds every package that has a `build` script, upstream packages first. |
| `pnpm turbo run dev --filter=<package-name>` | Builds the packages an app depends on, then starts the app. |
| `pnpm turbo run typecheck` | Runs each package's `typecheck` script. |
| `pnpm turbo run test` | Runs each package's `test` script. |
| `pnpm turbo run check` | Runs Biome format, lint and import checks on the whole repository. Read-only. |
| `pnpm turbo run fix` | Same as `check`, writing safe fixes. Review the diff afterwards. |

`biome.json` owns formatting and import order: tabs, lines up to 80 columns, single quotes (double in JSX attributes), no semicolons, and imports grouped as packages, then aliases (`@/`, `#`, `~` and the website's `@app/*`, `@components/*`, `@i18n/*`, `@state`, …), then relative paths, separated by blank lines. An object literal with more than one property puts each property on its own line; `platform/quality/multi-property-object.grit` reports violations.

`check` and `fix` are root tasks with caching off on purpose: their inputs span every package, which a root task's cache key does not reliably cover, and Biome is fast enough without it.

## Layout

- `apps/<app>/`: one deployable app each, with its own URL. `pnpm-workspace.yaml` picks them up through `apps/*`.
- `packages/<capability>/`: shared across apps, never inside one.
- Packages: `tokens`, `web` and `ex` exist; `shell`, `account`, `runtime` and `contracts` are planned.
- `packages/ex` is Go: Turborepo runs its `build`, `test` and `typecheck` scripts like any package's, and `pnpm install` never builds it (no Go toolchain required for frontend work).
- A package takes an app's name only while it has a single consumer. When a second consumer appears, it becomes a capability package.
- Package scope: `@langdrift/*`.
- `.github/`: GitHub Actions. `workflows/pull-request.yml` validates pull requests; publishing gets its own workflow file. `actions/setup-workspace/` installs pnpm, Node.js and the dependencies.
- `turbo.json` declares the environment each task may see: Turborepo runs in strict mode, so a variable reaches a task only through `globalEnv`, a task's `env` or the pass-through lists. `apps/website/README.md` holds the website's variable contract, and its `.env.development` the local defaults; `CRYPTO_SECRET` is pass-through only and never committed.
- `platform/`: repository tooling, grouped by concern. `cli/` holds the `drift` CLI, `quality/` the lint-staged config and the Biome plugin, `versioning/` the commitlint config and the Git hooks (`.husky/`).

## Platform CLI

`platform/cli/drift` is the repository's POSIX shell CLI. `pnpm install` links it into your bin directory (`PNPM_HOME`, `XDG_BIN_HOME` or `~/.local/bin`; skipped in CI), so `drift` runs from anywhere. `drift --help` lists the commands: `setup`, `doctor`, `cleanup`, `check localization`, `env setup|validate` and `git setup|doctor|pre-commit|commit-msg|lint-history`.

When the shell resolves `drift` to another program (a link from an older checkout, say), `drift setup` says which one and how to fix `PATH`.

## Design and tokens

- `@langdrift/tokens` builds the tokens with Style Dictionary. Each app imports one self-contained stylesheet built from the shared tokens and its app theme; today the only one is `@langdrift/tokens/website.css`.
- Purple (primary) surfaces such as buttons carry `primary-foreground` text, which is white in both color schemes.
- Components use roles, never a foundation or a brand color directly. Roles take shadcn's names (`background`, `foreground`, `primary`, `muted-foreground`, …) because the web apps use SmoothUI and shadcn components.
- Prefixes: shared tokens are `--lg-*`; the website's app theme is `--site-*`. Each app theme lives in `packages/tokens/themes/` and defines only primary and primary hover (the Exphere purple for the website), which `--lg-color-primary` points to.
- Status colors (`--lg-status-*`), classification colors (`--lg-classification-*`), classification text tones (`--lg-classification-text-*`) and drift colors (`--lg-drift-*`) are fixed across apps and never inherit a brand color. See ADR 0003.
- Light is declared in `:root` and dark in `.dark`. Apps render with the `.dark` class, so dark is the default, and the stylesheet points Tailwind's `dark:` variant at that class.
- Each stylesheet also bridges the roles to shadcn's unprefixed variables and sets Tailwind's theme. The native color palette is reset, so only roles and the fixed color groups exist as Tailwind colors. See ADR 0001 and `03-tailwind.rule.md`.

## Navigation and session

- The app switcher links to another domain or subdomain. It is not internal routing.
- Account and session come from a central account app, not from each app.
- The switcher and the shell live in `packages/shell`. Each app owns its sidebar and its content.

## Spec-driven development

Workflow: `research` when a question needs primary sources, then `grill-with-docs` (or `grill-me` for plans outside the code), `to-spec`, `to-tickets`, `implement` and `code-review`.

- `.agents/workflow/human-in-the-loop.policy.md` sets where the owner takes part, which stages run on their own and what gates them. It belongs to this project and team; the skills read it instead of assuming.
- A direct request to build or change application code enters the workflow: start with the `grilling` skill, plus `domain-modeling` when terms or decisions settle, instead of writing code.
- `.claude/hooks/guard-implement-effort.sh` blocks `implement` for the current session until `.agents/prompts/llm-effort-analyzer.prompt.md` has recorded its result at `.agents/workflow/.effort-log/<session_id>.md`.
- The commit-msg hook rejects a commit that changes `apps/` or `packages/` without a `Ticket: <NNNN>` trailer naming an existing ticket.

## Agent skills

The skills in `.claude/skills/` are adapted from mattpocock/skills (MIT, see `.claude/skills/LICENSE`). Each `SKILL.md` records its source in its frontmatter.

### Issue tracker

Local Markdown under `.agents/workflow/`:

| Artifact | Path |
| --- | --- |
| Spec | `specs/<NNNN>-<slug>.spec.md` |
| Ticket | `tickets/<NNNN>-<slug>.ticket.md` |
| ADR | `adrs/<NNNN>-<slug>.adr.md` |
| Research note | `research/<slug>.research.md` |

- Each kind numbers its own sequence from `0001`; take the next free number. Create `adrs/` and `research/` when the first file is written.
- "Publish to the issue tracker" means writing the file. "Fetch a ticket" means reading it.
- A ticket names its spec in a `**Spec:**` line.

### Triage labels

A `**Status:**` line near the top of each spec or ticket holds one of `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human` or `wontfix`.

### Domain docs

Single context. The glossary is the `## Glossary` section of `.agents/context.md`, and ADRs live in `.agents/workflow/adrs/`.

## How we work

- The owner writes application code. Agents do not edit app code unless the task asks for it.
- Repetitive work may be delegated when the task names its scope.
- The harness (checks, guardrails, scripts) stays a draft until the owner promotes a piece of it to a rule.
- Rules are added only when the owner asks.

## Continuous integration

`.github/workflows/pull-request.yml` runs on every pull request into `main`, one job per concern. The workflow builds the website against safe values (early access paused, sign-up disabled), so CI never reaches a live service:

| Job | Runs |
| --- | --- |
| Build | `pnpm turbo run build` |
| Check | `pnpm turbo run check`, then `pnpm turbo run typecheck` |
| Test | `pnpm turbo run test` |
| Harness | `drift doctor --ci`, `drift check localization`, ShellCheck on the POSIX scripts, actionlint on the workflows |
| Guardrails | `drift git lint-history` over the pull request's commits |

Actions, container images and downloaded tools are pinned to full commit SHAs, digests or checksums, with the version in a comment. The workflow tokens are read-only.

## Git

- No commit, pull request or contribution carries an AI-generated signature, trailer or footer. This includes `Co-Authored-By` and "Generated with".
- Inside the spec-driven workflow, commit, push and pull request follow `.agents/workflow/human-in-the-loop.policy.md`. Outside it, each needs explicit authorization every time. Merge and release always need it.
- The hooks in `platform/versioning/.husky/` call `drift`. `drift git pre-commit` runs Biome on the staged files through lint-staged. `drift git commit-msg` rejects AI signatures, enforces the ticket trailer above and requires Conventional Commits through commitlint; its checks live in `platform/cli/src/core/commit-rules.sh`, which `drift git lint-history` also applies to every commit in a range. Husky installs the hooks on `pnpm install` once the folder is a Git repository. Neither lint-staged nor commitlint finds its config by itself: `drift` passes `--config`, so run them through `drift` rather than `pnpm exec`.

## Boundaries

- Do not read or edit values in `.env` files.
- Keep credentials out of code, logs and evidence.
- An agent does not close an open decision. Pricing, limits, auth provider and formulas stay open until the owner decides.

## Open

- Domain layout: subdomains under one domain, or separate domains per app.

<!-- BEGIN:turborepo-agent-rules -->

# This is NOT the Turborepo you know

Turborepo configuration, task behavior, and CLI commands can vary between installed versions and may differ from your training data. Resolve the `turbo` package from this file's directory or relevant workspace; in monorepos, it may not be visible from the repository root. For example, run `node -p "require.resolve('turbo/package.json')"` from a workspace that depends on `turbo`.

Read `docs/README.md` inside that installed package first, then read the relevant pages from its `docs/` directory before changing Turborepo configuration or commands. Heed deprecation notices. These bundled docs match the installed package version and are available without network access.

This block is written and re-added by `turbo` before repository-scoped commands when an AI agent is detected. In the Turborepo source repository, its template is defined in `crates/turborepo-cli/src/cli/agent_guidance.rs`. Removing the managed block while updates are enabled means a later qualifying invocation will add it again. Set `"agentGuidance": false` in the root `turbo.json` or `turbo.jsonc` to opt out; this does not remove an existing block. Keep the block committed with your work to avoid an uncommitted change on the next agent invocation.
<!-- END:turborepo-agent-rules -->
