# LangDrift

AI harness web OS for your agentic journey. LangDrift is a suite of products under one account: each app is independent, with its own URL, and all apps share account, plan and shell.

## General rules

`.claude/rules/` holds the rules that always apply, read first in every session. `00-general.rule.md` covers code, `01-filesystem.rule.md` naming and paths, and `02-frontend.rule.md` React and Next.js styling.

## Status

Monorepo skeleton with no apps or packages yet. Planned apps:

- `apps/vision`: LangDrift Vision
- `apps/finops`: LangDrift FinOps

## Toolchain and commands

- pnpm is the only package manager. `devEngines.packageManager` in the root `package.json` makes other package managers fail, and `engineStrict` in `pnpm-workspace.yaml` enforces the Node range in `engines`.
- Work that spans the monorepo runs through Turborepo. Do not add root scripts that wrap `turbo run`.
- A package's scripts call its tool directly (for example `tsc --noEmit`). Turborepo tasks run those scripts by name.
- Do not add shell or JavaScript wrappers for what Turborepo or pnpm already do.

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm install` | Installs dependencies and the Git hooks, then links the `drift` launcher. Fails under another package manager or an unsupported Node version. |
| `pnpm turbo run build` | Builds every package that has a `build` script, upstream packages first. |
| `pnpm turbo run dev --filter=<package-name>` | Starts one app and the packages it depends on. |
| `pnpm turbo run typecheck` | Runs each package's `typecheck` script. |
| `pnpm turbo run test` | Runs each package's `test` script. |
| `pnpm turbo run check` | Runs Biome format, lint and import checks on the whole repository. Read-only. |
| `pnpm turbo run fix` | Same as `check`, writing safe fixes. Review the diff afterwards. |

`biome.json` owns formatting and import order: tabs, lines up to 80 columns, single quotes (double in JSX attributes), no semicolons, and imports grouped as packages, then aliases (`@/`, `#`, `~`), then relative paths, separated by blank lines. An object literal with more than one property puts each property on its own line; `platform/quality/multi-property-object.grit` reports violations.

`check` and `fix` are root tasks with caching off on purpose: their inputs span every package, which a root task's cache key does not reliably cover, and Biome is fast enough without it.

## Layout

- `apps/<product>/`: one app per product.
- `packages/<capability>/`: `ui`, `tokens`, `shell`, `account`, `runtime`, `contracts`.
- A package takes a product name only while it has a single consumer. When a second consumer appears, it becomes a capability package.
- Package scope: `@langdrift/*`.
- `platform/`: repository tooling, grouped by concern. `cli/` holds the `drift` CLI, `quality/` the lint-staged config and the Biome plugin, `versioning/` the commitlint config and the Git hooks (`.husky/`).

## Platform CLI

`platform/cli/drift` is the repository's POSIX shell CLI. `pnpm install` links it into your bin directory (`PNPM_HOME`, `XDG_BIN_HOME` or `~/.local/bin`; skipped in CI), so `drift` runs from anywhere. `drift --help` lists the commands: `setup`, `doctor`, `cleanup`, `env setup|validate` and `git setup|doctor|pre-commit|commit-msg|lint-history`.

When the shell resolves `drift` to another program (a link from an older checkout, say), `drift setup` says which one and how to fix `PATH`.

## Design and tokens

- Components use semantic roles, never a brand color directly.
- Neutral scales and roles are shared. Each product defines only its brand in `packages/tokens/themes/`.
- Token prefix: `--lg-*`.
- Status colors (`--lg-status-*`) are fixed across products and never inherit a brand color.

## Navigation and session

- The app switcher links to another domain or subdomain. It is not internal routing.
- Account and session come from a central account app, not from each app.
- The switcher and the shell live in `packages/shell`. Each app owns its sidebar and its content.

## Spec-driven development

Workflow: `research` when a question needs primary sources, then `grill-with-docs` (or `grill-me` for plans outside the code), `to-spec`, `to-tickets`, `implement` and `code-review`.

- `grill-with-docs`, `grill-me`, `to-spec`, `to-tickets` and `implement` run only when the owner invokes them.
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

## Git

- No commit, pull request or contribution carries an AI-generated signature, trailer or footer. This includes `Co-Authored-By` and "Generated with".
- Commit, push, pull request, merge and release each need explicit authorization every time.
- The hooks in `platform/versioning/.husky/` call `drift`. `drift git pre-commit` runs Biome on the staged files through lint-staged. `drift git commit-msg` rejects AI signatures, enforces the ticket trailer above and requires Conventional Commits through commitlint; its checks live in `platform/cli/src/commands/git-commit-msg.sh`. Husky installs the hooks on `pnpm install` once the folder is a Git repository. Neither lint-staged nor commitlint finds its config by itself: `drift` passes `--config`, so run them through `drift` rather than `pnpm exec`.

## Boundaries

- Do not read or edit values in `.env` files.
- Keep credentials out of code, logs and evidence.
- An agent does not close an open decision. Pricing, limits, auth provider and formulas stay open until the owner decides.

## Open

- Domain layout: subdomains under one domain, or separate domains per product.
