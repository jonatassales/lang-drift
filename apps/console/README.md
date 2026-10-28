# Console

The LangDrift console, where a team manages FinOps: projects registered from a local folder, their microruntimes, executions, history and secrets, with FinOps readings and recommendations that stay advisory. It was imported from Exphere's console (`jonatassales/exphere`) and carries the LangDrift brand. Every screen reads [extub](/apps/extub/README.md), the local stub API, so nothing here reaches a live service.

The sidebar opens with an area switcher, after the old LangDrift dashboard's product switcher. **Microruntimes** holds home, projects and each project's runtimes, executions, history and secrets. **FinOps**, under `/finops`, opens with a cost overview: the LLM cost, CPU and memory each runtime reported against its limits and budget, where a cost no execution reported stays unavailable. The URL decides the area.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: build the shared packages, then start the console on `http://127.0.0.1:3001` with extub alongside it on `127.0.0.1:3401`. |
| `pnpm --filter console typecheck` | Generate Next.js route types and check TypeScript. |
| `pnpm --filter console test` | Run the runtime view tests with Node's test runner. |

The console listens on loopback only and answers only loopback hosts and same-origin requests. `EXTUB_URL` points it at another extub; it must stay plain `http` on loopback, and without it the console uses `http://127.0.0.1:3401`. When extub is down, data requests answer `extub_unavailable` and the screens say so.
