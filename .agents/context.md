# Context

Project knowledge for agents working in this repository. [`AGENTS.md`](/AGENTS.md) holds the ideas; this file holds the mechanics and the vocabulary. It grows as decisions are made; it does not get restructured into subfolders unless the owner asks.

## Repository

Monorepo with the internal packages `tokens`, `streaming` and `web`, the published `ex` CLI as `@langdrift/cli` with the Microruntimes contracts in `@langdrift/core`, and three apps. `apps/website` is where the platform site is built; it started as today's langdrift.md, copied from the earlier `langdrift` monorepo under spec 0003, and has grown into the platform site. `apps/extub` is a local Go API of FinOps console stubs and seeds, copied from Exphere and adapted to this repository. `apps/console` is the console imported from Exphere, reading its data from extub.

Four workspace roots, each with one meaning, and a naming rule that follows it (spec 0044). `pnpm-workspace.yaml` lists them all, and a commit under any of them carries a `Ticket:` trailer.

- `apps/<app>/`: one deployable app each, Next.js except `extub`. Picked up through `apps/*`.
- `apps/microservices/<service>/`: the services, picked up through `apps/microservices/*`: the Nest.js microservices (`demo-requests`, `mailer`), which share the event log through `streaming` (see `.claude/rules/11-nestjs.rule.md`), and, once spec 0044 lands, `console-stubs`, the Go stub API the console reads (today `apps/extub`).
- `packages/<name>/`: packages internal to the monorepo, shared across apps and never published. Their names are bare, like the apps' (`web`, `tokens`, `streaming`). The scope they carried left under spec 0044. `shell`, `account` and `runtime` are planned.
- `lib/<name>/`: what LangDrift publishes, and only that, under the `@langdrift/*` scope: `@langdrift/cli` at `lib/cli` (the `ex` CLI) and `@langdrift/core` at `lib/core`, whose Go surface, the Microruntimes contracts, is the module `lib/core/contracts/go`. The scope means "published"; nothing outside `lib/` carries it.
- `infrastructure/<component>/`: the infrastructure components the products run on, each a Dockerfile, a README and a minimal private manifest: `messaging` (Redis Streams, the event log) and `storage` (MinIO, the archive). A fence names a component after its folder (ADR 0004).
- The Go modules (`lib/cli`, `lib/core/contracts/go` and `apps/extub`) are tied together by the root `go.work`: Turborepo runs their `build`, `test` and `typecheck` scripts like any package's, and `pnpm install` never builds them (no Go toolchain required for frontend work).
- A package takes an app's name only while it has a single consumer. When a second consumer appears, it becomes a capability package.
- `.github/`: GitHub Actions. Entry workflows are named after their event (`pull-request.yml`, `main.yml`), `validate.yml` is the validation they share; composite actions sit in `actions/<name>/` and the pull request widgets in `widgets/`. Publishing gets its own workflow file. See Continuous integration.
- `platform/`: repository tooling, grouped by concern. `cli/` holds the `drift` CLI, its Git hooks (`.husky/`), its `src/exceptions/dead-code.allow` and the detectors `drift quality check` and `coverage` share (`src/quality/`), `docs/` the README's translations and, under `assets/images/`, its light and dark banners with the SVG sources the PNGs are rendered from. The English `README.md` sits at the root; each README links every language right after the banner and the headline.
- `turbo.json` declares the environment each task may see: Turborepo runs in strict mode, so a variable reaches a task only through `globalEnv`, a task's `env` or the pass-through lists. `apps/website/README.md` holds the website's variable contract and its tracked `.env.development` the shared development defaults; in development, `CRYPTO_SECRET` is generated on the fly if unset. Deployed environments set their own value, and `CRYPTO_SECRET` stays pass-through only in Turborepo. The `build` task declares `VERCEL_ENV` in its `env`: the website's header policy differs on preview builds, so a preview and a production build must never share a cache entry.

## Toolchain and commands

- pnpm is the only package manager. `devEngines.packageManager` in the root `package.json` makes other package managers fail, and an older pnpm, such as the one Vercel's build image ships, downloads the required version and runs it instead of failing (`onFail: download`). `engineStrict` in `pnpm-workspace.yaml` enforces the Node range in `engines`.
- Dependency install scripts run only when `allowBuilds` in `pnpm-workspace.yaml` allows them. Prebuilt binaries ship as optional dependencies, so the entries there are denials; approve a new one only when its build is required.
- Work that spans the monorepo runs through Turborepo. The only root scripts that wrap `turbo run` are the two app launchers below; do not add others.
- A package's scripts call its tool directly (for example `tsc --noEmit`). Turborepo tasks run those scripts by name. Do not add shell or JavaScript wrappers for what Turborepo or pnpm already do.

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm install` | Installs dependencies and the Git hooks, then links the `drift` launcher. Fails under another package manager or an unsupported Node version. |
| `pnpm website` | Builds the packages the website depends on, then starts it in dev mode on port 3000. |
| `pnpm console` | Builds the packages the console depends on, then starts it in dev mode on `http://127.0.0.1:3001` together with its stub API, extub, on `127.0.0.1:3401`. Extub needs a Go toolchain. |
| `pnpm turbo run build` | Builds every package that has a `build` script, upstream packages first. |
| `pnpm turbo run dev --filter=<package-name>` | Builds the packages an app depends on, then starts the app. |
| `pnpm turbo run typecheck` | Runs each package's `typecheck` script. |
| `pnpm turbo run test` | Runs each package's `test` script. |
| `pnpm turbo run check` | Runs Biome format, lint and import checks on the whole repository. Read-only. |
| `pnpm turbo run fix` | Same as `check`, writing safe fixes. Review the diff afterwards. |
| `lib/cli/dist/ex …` | The Microruntimes session in a project (`@langdrift/cli` ships a `bin/ex` launcher and prebuilt binaries on npm for `npx @langdrift/cli …` / `pnpm dlx @langdrift/cli …`): `ex init`, the `ex harness` and `ex finops` inputs, then `ex agent plan`, `ex agent run`, `ex agent confirm`, `ex agent status` and `ex agent reindex`. `lib/cli/README.md` has the command table. |

The website needs `CRYPTO_SECRET` for encrypted locale cookies. In development, it is generated on the fly when unset; deployed environments set a canonical 32-byte key; `apps/website/README.md` says how to generate one for a deployment.

`biome.json` owns formatting and import order: tabs, lines up to 80 columns, single quotes (double in JSX attributes), no semicolons, and imports grouped as packages, then aliases (`@/`, `#`, `~` and the website's `@app/*`, `@components/*`, `@i18n/*`, `@state`, …), then relative paths, separated by blank lines.

`check` and `fix` are root tasks with caching off on purpose: their inputs span every package, which a root task's cache key does not reliably cover, and Biome is fast enough without it.

## Platform CLI

`platform/cli/drift` is the repository's POSIX shell CLI. `pnpm install` links it into your bin directory (`PNPM_HOME`, `XDG_BIN_HOME` or `~/.local/bin`; skipped in CI), so `drift` runs from anywhere. `drift --help` lists the commands: `setup`, `doctor`, `cleanup`, `reset`, `check biome|localization|workflows`, `env setup|validate`, `git setup|doctor|pre-commit|commit-msg|lint-history`, `quality check|coverage|dead-code` and `widget check-run|sticky-comment`. The `widget` commands serve the workflows (see Continuous integration) and need `gh` and `jq`.

`drift env setup` creates each `.env` from its tracked `.env.development`, and `drift env validate` checks that each `.env.development` has its `.env`.

`drift quality dead-code` finds code that only served something a refactor removed: a tsconfig alias nothing imports through, a module or export nothing uses, a dependency nothing imports, a `turbo.json` variable nothing reads, a CSS class no markup carries. A reference in a comment counts, so code parked in comments keeps what it needs, and folders held by a `.gitkeep` are never examined. Run it after a refactor; CI runs it in the Check job. `platform/cli/src/exceptions/dead-code.allow` lists what the owner keeps on purpose.

The `quality` commands measure the tree against `.claude/rules`. Both run the detectors in `platform/cli/src/quality/rule-detectors.tsv` (`run-detectors.sh`): `rg` patterns, awk programs for statement spacing, Biome's import order and checks on the tracked paths. They are patterns, so what they find is an estimate to review, and neither command fails on it. Test files and testdata are skipped, except by the testing rules; a rule file no detector covers shows as not measured. Both need `rg`, which CI installs in `setup-workspace`.

- `drift quality check [rule]` prints, in red, each path that does not apply a rule file, with the rule and the violation; the pre-commit hook and CI run it.
- `drift quality coverage [rule]` prints, per rule file, the share of the files its detectors look at with no finding; the Code Quality widget shows the same table.

A rule narrows either command to one file of `.claude/rules`: `drift quality coverage react.rule.md`, `05-react.rule.md` or `react`.

When the shell resolves `drift` to another program (a link from an older checkout, say), `drift setup` says which one and how to fix `PATH`.

## Design and tokens

- `tokens` builds the tokens with Style Dictionary. Each app imports one self-contained stylesheet built from the shared tokens and its app theme: `tokens/website.css`, which the website and the console share.
- Purple (primary) surfaces such as buttons carry `primary-foreground` text, which is white in both color schemes.
- Components use roles, never a foundation or a brand color directly. Roles take shadcn's names (`background`, `foreground`, `primary`, `muted-foreground`, …) because the web apps use SmoothUI and shadcn components.
- Prefixes: shared tokens are `--lg-*`; the website's app theme is `--site-*`. Each app theme lives in `packages/tokens/themes/` and defines only primary and primary hover (the Exphere purple for both), which `--lg-color-primary` points to.
- Status colors (`--lg-status-*`), classification colors (`--lg-classification-*`), classification text tones (`--lg-classification-text-*`), drift colors (`--lg-drift-*`) and rainbow colors (`--lg-rainbow-1` to `--lg-rainbow-6`) are fixed across apps and never inherit a brand color; the rainbow colors serve the header's demo button highlight. See ADR 0003.
- Light is declared in `:root` and dark in `.dark`. Apps render with the `.dark` class, so dark is the default, and the stylesheet points Tailwind's `dark:` variant at that class.
- Each stylesheet also bridges the roles to shadcn's unprefixed variables and sets Tailwind's theme. The native color palette is reset, so only roles and the fixed color groups exist as Tailwind colors. See ADR 0001 and `.claude/rules/03-tailwind.rule.md`.

## Navigation and session

- Everything is served under `https://langdrift.md` with a locale prefix. Moving between pages, on the site and in the console, is internal routing on that domain; it never links to another domain or subdomain.
- Account and session come from a central account package, `packages/account`, which is planned. No app owns them.
- The shell will live in `packages/shell`, which is planned.
- The console lives here as `apps/console`, imported from `jonatassales/exphere`, and reads the local stub API `apps/extub`.

## Spec-driven development

Workflow: `research` when a question needs primary sources, then `grill-with-docs` (or `grill-me` for plans outside the code), `to-spec`, `to-tickets`, `implement` and `code-review`.

- `.agents/policies/human-in-the-loop.policy.md` sets where the owner takes part, which stages run on their own and what gates them. It belongs to this project and team; the skills read it instead of assuming.
- `.agents/policies/safe-seat.policy.md` sets the dev part that every ticket reserves for the owner to write by hand. `to-tickets` proposes it, `implement` leaves its paths alone and stops when only they remain, and `code-review` reviews it once the owner has written it.
- A direct request to build or change application code enters the workflow: start with the `grilling` skill, plus `domain-modeling` when terms or decisions settle, instead of writing code.
- `.claude/hooks/guard-implement-effort.sh` blocks `implement` for the current session until the `llm-effort-analyzer` skill (`.claude/skills/llm-effort-analyzer/SKILL.md`) has recorded its result at `.agents/workflow/.effort-log/<session_id>.md`.
- The commit-msg hook rejects a commit that changes `apps/`, `packages/`, `lib/` or `infrastructure/` without a `Ticket: <NNNN>` trailer naming an existing ticket.

### Agent skills

The skills in `.claude/skills/` are adapted from mattpocock/skills (MIT, see `.claude/skills/LICENSE`), except `llm-effort-analyzer`, which is this project's own. Each `SKILL.md` records its source in its frontmatter.

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

GitHub reads workflows only from the top of `.github/workflows/`, so file names carry the structure. An entry workflow is named after the event that starts it, a reusable one after the work it does, and it declares only `workflow_call`. Each composite action is `.github/actions/<name>/action.yml` and nothing else; scripts belong to `drift`, except the pull request widgets, which live in `.github/widgets/` and run through `drift widget`. All names are kebab-case.

| File | Runs on | Does |
| --- | --- | --- |
| `workflows/pull-request.yml` | a pull request into `main` | `validate.yml`, then Guardrails |
| `workflows/main.yml` | a push to `main` | `validate.yml`, then the `ex` binaries |
| `workflows/validate.yml` | a call from the two above | Build, Check, Test and Security |
| `actions/setup-workspace/` | a step | installs pnpm, Node.js, Go and the dependencies |
| `widgets/` | a step, through `drift widget` | the pull request widgets: a check run and a sticky comment (see Widgets) |

Sign-up is paused in the website's code, so CI builds it without any service values and never reaches a live service. A reusable workflow's checks read `Validate / Build`:

| Job | Runs |
| --- | --- |
| Build | `pnpm turbo run build` |
| Check | `pnpm turbo run check`, `pnpm turbo run typecheck`, `drift check localization`, then `drift quality dead-code` |
| Test | `pnpm turbo run test` |
| Security | `pnpm audit --prod --audit-level=high`, then govulncheck, pinned, in each Go module; it fails on a high or critical npm advisory or on a Go vulnerability our code reaches |
| Guardrails | pull requests only: `drift git lint-history` over the pull request's commits, then `drift check workflows` |
| ex binaries | `main` only, after Validate: `ex` for linux and darwin on amd64 and arm64, without cgo, with `SHA256SUMS` and a build provenance attestation (`actions/attest`) |

Actions, container images and downloaded tools are pinned to full commit SHAs, digests or checksums, with the version in a comment. The workflow tokens are read-only; a job that needs more, such as the attestation's `id-token` and `attestations`, asks for it itself.

`drift check workflows` enforces these conventions: actionlint, pinned, then the layout, a top-level `permissions:` in every workflow, a SHA pin with its version for every action outside the repository, no `${{ }}` inside a `run:` script (values reach scripts through `env:`) and `retention-days` on every uploaded artifact.

### Widgets

A widget is what a workflow shows on a pull request instead of plain job output. `.github/widgets/` holds one `drift widget` command per kind. Each reads its Markdown on stdin, so a job renders the content and pipes it in.

| Command | Shows | Needs |
| --- | --- | --- |
| `drift widget check-run <name> <conclusion> <title>` | a completed check run on commit `WIDGET_SHA`, with the title and the Markdown summary | `checks: write`, `GH_TOKEN`, `WIDGET_SHA` |
| `drift widget sticky-comment <key>` | one comment per key on pull request `WIDGET_PULL_REQUEST`, edited in place on every run | `pull-requests: write`, `GH_TOKEN`, `WIDGET_PULL_REQUEST` |

- The conclusion is one GitHub accepts: `success`, `failure`, `neutral`, `cancelled`, `skipped`, `timed_out` or `action_required`. On a pull request, `WIDGET_SHA` is `github.event.pull_request.head.sha`, the commit the pull request shows, not the merge commit in `GITHUB_SHA`.
- A sticky comment starts with the hidden marker `<!-- drift-widget:<key> -->`, and its key is kebab-case. Only a comment a bot wrote with that marker is edited, so a person quoting the marker is never overwritten. Publish it with a bot token, `github.token` or a GitHub App's: a personal token posts as a user, so its comments would never be found again and new ones would pile up. It belongs to `pull_request` runs; a `push` to `main` has no pull request to comment on.
- Content past GitHub's limit, 65,535 characters for a summary and 65,536 for a comment, is cut with a note.
- Values reach the commands through `env:` (Actions sets `GITHUB_REPOSITORY`), and only the job that publishes a widget asks for its write permission. A pull request from a fork runs with a read-only token, so its widgets cannot be published.

A job in `pull-request.yml`, whose top-level `permissions:` stays `contents: read`:

```yaml
  report:
    name: Report
    runs-on: ubuntu-24.04
    timeout-minutes: 5
    permissions:
      contents: read
      checks: write
      pull-requests: write
    steps:
      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1
        with:
          persist-credentials: false
      - name: Report
        env:
          GH_TOKEN: ${{ github.token }}
          WIDGET_SHA: ${{ github.event.pull_request.head.sha }}
          WIDGET_PULL_REQUEST: ${{ github.event.pull_request.number }}
        run: |
          platform/cli/drift widget check-run 'Validate report' success 'All checks passed' <report.md
          platform/cli/drift widget sticky-comment validate-report <report.md
```

### Artifacts

| Artifact | From | Kept | Holds |
| --- | --- | --- | --- |
| `ex-<sha>` | ex binaries | 30 days | the four `ex` binaries and `SHA256SUMS`; `gh attestation verify <binary> -R jonatassales/lang-drift` checks one. A build, not a release |

## Git

- No commit, pull request or contribution carries an AI-generated signature, trailer or footer. This includes `Co-Authored-By` and "Generated with".
- Inside the spec-driven workflow, commit, push and pull request follow `.agents/policies/human-in-the-loop.policy.md`. Outside it, each needs explicit authorization every time. Merge and release always need it.
- The hooks in `platform/cli/.husky/` call `drift`. `drift git pre-commit` runs Biome on the staged files through lint-staged. `drift git commit-msg` rejects AI signatures, enforces the ticket trailer above and requires Conventional Commits through commitlint; its checks live in `platform/cli/src/core/commit-rules.sh`, which `drift git lint-history` also applies to every commit in a range. Husky installs the hooks on `pnpm install` once the folder is a Git repository. Both commitlint and lint-staged are configured directly in `package.json`.

## Boundaries

- Do not read or edit values in `.env` files.
- Keep credentials out of code, logs and evidence.
- An agent does not close an open decision. Pricing, limits, auth provider and formulas stay open until the owner decides.

## Glossary

### Platform

**Platform product**:
LangDrift's offering under the shared brand, named `LangDrift <Product>`: today only LangDrift FinOps, with Microruntimes as its engine. It has its pages on the site and its workspaces in the console.
_Avoid_: Module, vertical, sub-app

**Console**:
The one authenticated LangDrift surface on langdrift.md where a team manages FinOps and its microruntimes, behind one account and one shell. It lives in `apps/console`, imported from Exphere's console.
_Avoid_: Dashboard (for the whole), admin, control panel

**Safe Seat**:
The LangDrift FinOps harness capability that keeps a person writing a part of every ticket by hand, so the team keeps understanding and debugging what it ships.
_Avoid_: SS, health, wellness, quota

### Apps

**App**:
One LangDrift offering with its own URL, such as the website. Apps share account, plan and shell.
_Avoid_: Product, module

**Event log**:
The Redis Stream the microservices append to and consume from through consumer groups: the source of truth of a demo request, replayed to rebuild its state and archived event by event into MinIO.
_Avoid_: Queue, message bus, topic

### Design tokens

**Shared token**:
A design token with the same value in every app.
_Avoid_: Global token, base token

**App theme**:
The tokens through which one app expresses its brand. It holds only the primary color and its hover state; the website's is the one theme today, shared with the console.
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
The color of one classification: intentional, aligned, unexplained or review. It is the same in every app and color scheme, and never takes the brand color.
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

**Component Dockerfile**:
The committed Dockerfile in a component's folder that a microruntime builds from. `ex` runs the fixed `docker build` against it from the project root and reads its first `EXPOSE` and its own `HEALTHCHECK` for the port and the readiness. See ADR 0004.
_Avoid_: Build recipe, build command, build step

**Execution ledger**:
The on-disk record of what an operation created, written before side effects and settled afterwards. Cleanup acts only on what the ledger owns.
_Avoid_: State file, journal

**Evidence**:
The token-free record ex keeps of every lifecycle operation and inspection, for later audit.
_Avoid_: Log, report

**Run**:
One execution of a microruntime toward one goal: a ticket, a test, a load test or an end-to-end run. The race of the metaphor in which the microruntime is the car.
_Avoid_: Job, session

**Lap**:
One pass inside a run: an attempt paused at a gate, an execution, a load step, a suite pass.
_Avoid_: Step, iteration (in new writing)

**Episode**:
One earlier run of the project as the RAG remembers it: the memory's fence, with its stacks, component paths, measured peaks and tokens, embedded for retrieval when a new ticket is planned.
_Avoid_: Example, case, memory (for the whole)

**Judge**:
The model the planner calls through a port to keep or discard each source with its facts, weigh the episodes and suggest a model and effort. Its calls are paid work under the money gates; without a credential the plan keeps the deterministic shell and the judged fields stay null.
_Avoid_: Planner model, LLM, evaluator

**Stage driver**:
The owner's coding agent as `ex` runs it for a stage's skills: a trusted process started with its configured arguments and environment, whose usage and one-line summary `ex` reads back.
_Avoid_: Agent runner, executor, worker

**Input record**:
One value that can block or allow a run, as the private store keeps it: who set it, when, through what, its value and the record it replaced. Gates read the record, never its mirror.
_Avoid_: Setting, config, preference

**Mirror**:
The JSON file `ex` writes from one of its input records so the consumer and a console can read it: the envelope, the workflow, the budget, the gates, the prices and the forecasts. No gate reads a mirror.
_Avoid_: Config file, export, copy
