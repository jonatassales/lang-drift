# Console

The LangDrift console, where a team manages FinOps: projects registered from a local folder, their microruntimes, executions, history and secrets, with FinOps readings and recommendations that stay advisory. It was imported from Exphere's console (`jonatassales/exphere`) and carries the LangDrift brand. Every screen reads [extub](/apps/extub/README.md), the local stub API, so nothing here reaches a live service.

The sidebar opens with an area switcher, after the old LangDrift dashboard's product switcher. **Microruntimes** holds home, projects and each project's runtimes, tickets, analysis, plans, workflows, executions, history and secrets. Tickets, Analysis, Plans and Workflows each list every ticket that has reached that phase, newest first and loading older ones as the page scrolls: a card on a narrow screen and a compact row from `lg`, each with the panel's one-line summary for that phase. An entry opens the ticket's page, a single document with no card that shows where the ticket stands, its linked artifacts and the body of every phase it has reached. **FinOps**, under `/finops`, opens with a cost overview: the LLM cost, CPU and memory each runtime reported against its limits and budget, where a cost no execution reported stays unavailable. The URL decides the area.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: build the shared packages, then start the console on `http://127.0.0.1:3001` with extub alongside it on `127.0.0.1:3401`. |
| `pnpm --filter console typecheck` | Generate Next.js route types and check TypeScript. |
| `pnpm --filter console test` | Run the view tests, runtimes and tickets, with Node's test runner. |

The console listens on loopback only and answers only loopback hosts and same-origin requests. `EXTUB_URL` points it at another extub; it must stay plain `http` on loopback, and without it the console uses `http://127.0.0.1:3401`. When extub is down, data requests answer `extub_unavailable` and the screens say so.

## Hosted preview

On Vercel's `langdrift-console` project the console reads extub's hosted preview (spec 0015). Setting `EXTUB_TOKEN` lets the console serve its own host, accept an `https` `EXTUB_URL` on any host, and send the token to extub as `Authorization: Bearer`. The preview is read-only: a change answers `read_only_preview`, and without a live stream the screens keep polling. Without the token, the loopback contract above holds. It needs two variables on Vercel:

| Variable | Requirement |
| --- | --- |
| `EXTUB_TOKEN` | The secret shared with `langdrift-stubs`; switches the console to the hosted preview. |
| `EXTUB_URL` | The stub's production URL, over `https`, without credentials. |

The owner's steps on Vercel, which no agent takes, are in [extub's README](/apps/extub/README.md#hosted-preview). They set the token on both projects and `EXTUB_URL` on the console, and keep the console behind Deployment Protection for production and preview.
