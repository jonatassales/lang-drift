# Context

Project knowledge for agents working in this repository. [`AGENTS.md`](/AGENTS.md) holds the ideas; this file holds the mechanics and the vocabulary. It grows as decisions are made; it does not get restructured into subfolders unless the owner asks.

## Repository

Monorepo with three packages, `@langdrift/tokens`, `@langdrift/web` and `@langdrift/ex` (the Go agent-harness CLI; see `packages/ex/README.md`), and two apps. `apps/website` is where the platform site is built; it started as today's langdrift.md, copied from the earlier `langdrift` monorepo under spec 0003, and has grown into the platform site. `apps/docs` is the Nextra documentation for LangDrift and every platform product, ported from Exphere's docs under spec 0006 and served under `basePath: '/docs'`.

- `apps/<app>/`: one deployable Next.js app each. `pnpm-workspace.yaml` picks them up through `apps/*`.
- `packages/<capability>/`: shared across apps, never inside one. `tokens`, `web` and `ex` exist; `shell`, `account`, `runtime` and `contracts` are planned.
- `packages/ex` is Go: Turborepo runs its `build`, `test` and `typecheck` scripts like any package's, and `pnpm install` never builds it (no Go toolchain required for frontend work).
- A package takes an app's name only while it has a single consumer. When a second consumer appears, it becomes a capability package. Package scope: `@langdrift/*`.
- `.github/`: GitHub Actions. `workflows/pull-request.yml` validates pull requests; publishing gets its own workflow file. `actions/setup-workspace/` installs pnpm, Node.js and the dependencies.
- `platform/`: repository tooling, grouped by concern. `cli/` holds the `drift` CLI, `quality/` the lint-staged config, the Biome plugin and `dead-code.allow`, `versioning/` the commitlint config and the Git hooks (`.husky/`).
- `turbo.json` declares the environment each task may see: Turborepo runs in strict mode, so a variable reaches a task only through `globalEnv`, a task's `env` or the pass-through lists. `apps/website/README.md` holds the website's variable contract and its tracked `.env.development` the shared development defaults, including a development-only `CRYPTO_SECRET`; deployed environments set their own value, and `CRYPTO_SECRET` stays pass-through only in Turborepo.

## Toolchain and commands

- pnpm is the only package manager. `devEngines.packageManager` in the root `package.json` makes other package managers fail, and `engineStrict` in `pnpm-workspace.yaml` enforces the Node range in `engines`.
- Dependency install scripts run only when `allowBuilds` in `pnpm-workspace.yaml` allows them. Prebuilt binaries ship as optional dependencies, so the entries there are denials; approve a new one only when its build is required.
- Work that spans the monorepo runs through Turborepo. The only root scripts that wrap `turbo run` are the two app launchers below; do not add others.
- A package's scripts call its tool directly (for example `tsc --noEmit`). Turborepo tasks run those scripts by name. Do not add shell or JavaScript wrappers for what Turborepo or pnpm already do.

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm install` | Installs dependencies and the Git hooks, then links the `drift` launcher. Fails under another package manager or an unsupported Node version. |
| `pnpm website` | Builds the packages the website depends on, then starts it in dev mode on port 3000. |
| `pnpm docs` | Builds the packages the docs depend on, then starts them in dev mode at `http://localhost:3001/docs`. |
| `pnpm turbo run build` | Builds every package that has a `build` script, upstream packages first. |
| `pnpm turbo run dev --filter=<package-name>` | Builds the packages an app depends on, then starts the app. |
| `pnpm turbo run typecheck` | Runs each package's `typecheck` script. |
| `pnpm turbo run test` | Runs each package's `test` script. |
| `pnpm turbo run check` | Runs Biome format, lint and import checks on the whole repository. Read-only. |
| `pnpm turbo run fix` | Same as `check`, writing safe fixes. Review the diff afterwards. |

The website needs `CRYPTO_SECRET` even locally; without it every page answers 500. The tracked `.env.development` carries a development-only key, so a fresh checkout runs; `apps/website/README.md` says how to generate one for a deployment.

`biome.json` owns formatting and import order: tabs, lines up to 80 columns, single quotes (double in JSX attributes), no semicolons, and imports grouped as packages, then aliases (`@/`, `#`, `~` and the website's `@app/*`, `@components/*`, `@i18n/*`, `@state`, …), then relative paths, separated by blank lines. An object literal with more than one property puts each property on its own line; `platform/quality/multi-property-object.grit` reports violations.

`check` and `fix` are root tasks with caching off on purpose: their inputs span every package, which a root task's cache key does not reliably cover, and Biome is fast enough without it.

## Platform CLI

`platform/cli/drift` is the repository's POSIX shell CLI. `pnpm install` links it into your bin directory (`PNPM_HOME`, `XDG_BIN_HOME` or `~/.local/bin`; skipped in CI), so `drift` runs from anywhere. `drift --help` lists the commands: `setup`, `doctor`, `cleanup`, `reset`, `check dead-code|localization`, `env setup|validate` and `git setup|doctor|pre-commit|commit-msg|lint-history`.

`drift check dead-code` finds code that only served something a refactor removed: a tsconfig alias nothing imports through, a module or export nothing uses, a dependency nothing imports, a `turbo.json` variable nothing reads, a CSS class no markup carries. A reference in a comment counts, so code parked in comments keeps what it needs, and folders held by a `.gitkeep` are never examined. Run it after a refactor; CI runs it in the Check job. `platform/quality/dead-code.allow` lists what the owner keeps on purpose.

When the shell resolves `drift` to another program (a link from an older checkout, say), `drift setup` says which one and how to fix `PATH`.

## Design and tokens

- `@langdrift/tokens` builds the tokens with Style Dictionary. Each app imports one self-contained stylesheet built from the shared tokens and its app theme: `@langdrift/tokens/website.css` and `@langdrift/tokens/docs.css`.
- Purple (primary) surfaces such as buttons carry `primary-foreground` text, which is white in both color schemes.
- Components use roles, never a foundation or a brand color directly. Roles take shadcn's names (`background`, `foreground`, `primary`, `muted-foreground`, …) because the web apps use SmoothUI and shadcn components.
- Prefixes: shared tokens are `--lg-*`; the website's app theme is `--site-*` and the docs' is `--docs-*`. Each app theme lives in `packages/tokens/themes/` and defines only primary and primary hover (the Exphere purple for both), which `--lg-color-primary` points to.
- Status colors (`--lg-status-*`), classification colors (`--lg-classification-*`), classification text tones (`--lg-classification-text-*`), drift colors (`--lg-drift-*`) and rainbow colors (`--lg-rainbow-1` to `--lg-rainbow-6`) are fixed across apps and never inherit a brand color; the rainbow colors serve the header's demo button highlight. See ADR 0003.
- Light is declared in `:root` and dark in `.dark`. Apps render with the `.dark` class, so dark is the default, and the stylesheet points Tailwind's `dark:` variant at that class.
- Each stylesheet also bridges the roles to shadcn's unprefixed variables and sets Tailwind's theme. The native color palette is reset, so only roles and the fixed color groups exist as Tailwind colors. See ADR 0001 and `.claude/rules/03-tailwind.rule.md`.

## Navigation and session

- Everything is served under `https://langdrift.md` with a locale prefix. The product switcher moves between platform products as internal routing on that domain; it never links to another domain or subdomain.
- Account and session come from a central account package, `packages/account`, which is planned. No app owns them.
- The switcher and the shell will live in `packages/shell`, which is planned. Each product owns its workspaces, sidebar and content.
- The two consoles that merge into the LangDrift console live upstream for now: `jonatassales/langdrift` `apps/dashboard` (Vision) and `jonatassales/exphere` `apps/console` (FinOps).

## Spec-driven development

Workflow: `research` when a question needs primary sources, then `grill-with-docs` (or `grill-me` for plans outside the code), `to-spec`, `to-tickets`, `implement` and `code-review`.

- `.agents/policies/human-in-the-loop.policy.md` sets where the owner takes part, which stages run on their own and what gates them. It belongs to this project and team; the skills read it instead of assuming.
- A direct request to build or change application code enters the workflow: start with the `grilling` skill, plus `domain-modeling` when terms or decisions settle, instead of writing code.
- `.claude/hooks/guard-implement-effort.sh` blocks `implement` for the current session until `.agents/prompts/llm-effort-analyzer.prompt.md` has recorded its result at `.agents/workflow/.effort-log/<session_id>.md`.
- The commit-msg hook rejects a commit that changes `apps/` or `packages/` without a `Ticket: <NNNN>` trailer naming an existing ticket.

### Agent skills

The skills in `.claude/skills/` are adapted from mattpocock/skills (MIT, see `.claude/skills/LICENSE`). Each `SKILL.md` records its source in its frontmatter.

### Issue tracker

Local Markdown under `.agents/workflow/`:

| Artifact | Path |
| --- | --- |
| Spec | `specs/<NNNN>-<slug>.spec.md` |
| Ticket | `tickets/<NNNN>-<slug>.ticket.md` |
| ADR | `adrs/<NNNN>-<slug>.adr.md` |
| Research note | `research/<slug>.research.md` |

- Each kind numbers its own sequence from `0001`; take the next free number. Create `research/` when the first file is written.
- "Publish to the issue tracker" means writing the file. "Fetch a ticket" means reading it.
- A ticket names its spec in a `**Spec:**` line.

### Triage labels

A `**Status:**` line near the top of each spec or ticket holds one of `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human` or `wontfix`.

### Domain docs

Single context. The glossary is the `## Glossary` section of this file, and ADRs live in `.agents/workflow/adrs/`.

## Continuous integration

`.github/workflows/pull-request.yml` runs on every pull request into `main`, one job per concern. Sign-up is paused in the website's code, so CI builds it without any service values and never reaches a live service:

| Job | Runs |
| --- | --- |
| Build | `pnpm turbo run build` |
| Check | `pnpm turbo run check`, `pnpm turbo run typecheck`, then `drift check dead-code` |
| Test | `pnpm turbo run test` |
| Guardrails | `drift git lint-history` over the pull request's commits |

Actions, container images and downloaded tools are pinned to full commit SHAs, digests or checksums, with the version in a comment. The workflow tokens are read-only.

## Git

- No commit, pull request or contribution carries an AI-generated signature, trailer or footer. This includes `Co-Authored-By` and "Generated with".
- Inside the spec-driven workflow, commit, push and pull request follow `.agents/policies/human-in-the-loop.policy.md`. Outside it, each needs explicit authorization every time. Merge and release always need it.
- The hooks in `platform/versioning/.husky/` call `drift`. `drift git pre-commit` runs Biome on the staged files through lint-staged. `drift git commit-msg` rejects AI signatures, enforces the ticket trailer above and requires Conventional Commits through commitlint; its checks live in `platform/cli/src/core/commit-rules.sh`, which `drift git lint-history` also applies to every commit in a range. Husky installs the hooks on `pnpm install` once the folder is a Git repository. Neither lint-staged nor commitlint finds its config by itself: `drift` passes `--config`, so run them through `drift` rather than `pnpm exec`.

## Boundaries

- Do not read or edit values in `.env` files.
- Keep credentials out of code, logs and evidence.
- An agent does not close an open decision. Pricing, limits, auth provider and formulas stay open until the owner decides.

## Glossary

### Platform

**Platform product**:
One of LangDrift's offerings under the shared brand, named `LangDrift <Product>`: Vision and FinOps. Each has its own page on the site and its own workspaces in the console.
_Avoid_: Module, vertical, sub-app

**Console**:
The one authenticated LangDrift surface on langdrift.md where every platform product is managed, behind one account and one shell. It unifies the upstream LangDrift dashboard and Exphere console.
_Avoid_: Dashboard (for the whole), admin, control panel

**Product switcher**:
The shell control that moves between platform products inside the console, as routing on the same domain.
_Avoid_: App switcher, domain switcher

### Apps

**App**:
One LangDrift offering with its own URL, such as the website. Apps share account, plan and shell.
_Avoid_: Product, module

### Design tokens

**Shared token**:
A design token with the same value in every app.
_Avoid_: Global token, base token

**App theme**:
The tokens through which one app expresses its brand. It holds only the primary color and its hover state; the website and the docs each have one.
_Avoid_: Theme (alone), product theme, brand theme, skin

**Foundation**:
A raw value on a scale, such as gray 50 to 950, that roles draw from. Components never use a foundation directly.
_Avoid_: Primitive, palette color

**Role**:
A token named for its purpose, such as background, foreground or primary, rather than for its value. Components use roles only.
_Avoid_: Semantic color, alias

**Status color**:
A role that signals success, warning, danger or info. It is the same in every app and never takes the brand color.
_Avoid_: Feedback color, alert color

**Classification color**:
The color of one Product Vision classification: intentional, aligned, unexplained or review. It is the same in every app and color scheme, and never takes the brand color.
_Avoid_: Category color, label color

**Classification text tone**:
The classification color at a contrast fit for small text: darker in light, lighter in dark. It exists for intentional, review and unexplained.
_Avoid_: Classification ink, label tone

**Drift color**:
The color of a drift's sign: negative, positive or zero. It is the same in every app and never takes the brand color.
_Avoid_: Delta color, trend color

**Rainbow color**:
One of six fixed spectrum stops behind the demo button's highlight. It is the same in every app and color scheme, and never takes the brand.
_Avoid_: Spectrum color, gradient stop

**Color scheme**:
Light or dark. Dark is the default in every app.
_Avoid_: Theme, mode

### Ex CLI

**Microruntime**:
One console-defined workload unit that ex runs on the project's local cluster: its configuration, execution and cleanup.
_Avoid_: Runtime (alone, in new writing), service

**Connected agent**:
The long-running `ex agent` process that supervises a project's microruntimes from authenticated console intent.
_Avoid_: Daemon, worker

**Build recipe**:
The console-supplied command that builds a microruntime's image. Only a committed project script or an allowlisted tool may run.
_Avoid_: Build command, build step

**Execution ledger**:
The on-disk record of what an operation created, written before side effects and settled afterwards. Cleanup acts only on what the ledger owns.
_Avoid_: State file, journal

**Evidence**:
The token-free record ex keeps of every lifecycle operation and inspection, for later audit.
_Avoid_: Log, report
