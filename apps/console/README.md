# Console

The LangDrift console, where a team manages FinOps: projects registered from a local folder, their microruntimes, executions, history and secrets, with FinOps readings and recommendations that stay advisory. It was imported from Exphere's console (`jonatassales/exphere`) and carries the LangDrift brand. Every screen reads [console-stubs](/apps/microservices/console-stubs/README.md), the local stub API, so nothing here reaches a live service.

The sidebar opens with an area switcher, after the old LangDrift dashboard's product switcher. **Microruntimes** holds home, projects and each project's runtimes, tickets, analysis, plans, workflows, executions, history and secrets. Tickets, Analysis, Plans and Workflows each list every ticket that has reached that phase, newest first and loading older ones as the page scrolls: a card on a narrow screen and a compact row from `lg`, each with the panel's one-line summary for that phase. An entry opens the ticket's page, a single document with no card that shows where the ticket stands, its linked artifacts and the body of every phase it has reached. **FinOps**, under `/finops`, opens with a cost overview: the LLM cost, CPU and memory each runtime reported against its limits and budget, where a cost no execution reported stays unavailable. The URL decides the area.

## Local development

| Command | Purpose |
| --- | --- |
| `pnpm console` | From the repository root: build the shared packages, then start the console on `http://127.0.0.1:3001` with console-stubs alongside it on `127.0.0.1:3401`. |
| `pnpm turbo run typecheck --filter=console` | Generate Next.js route types, then check TypeScript. |
| `pnpm --filter console test` | Run the view tests, runtimes and tickets, with Node's test runner. |

The console listens on loopback only and answers only loopback hosts and same-origin requests. `CONSOLE_STUBS_URL` points it at console-stubs, configured in `.env.development` as `http://127.0.0.1:3401`; it must stay plain `http` on loopback. When console-stubs is down, data requests answer `console_stubs_unavailable` and the screens say so.

## Security and search

- **A nonce policy on every document.**
  - The proxy builds a `Content-Security-Policy` per request around a fresh nonce: `script-src 'self' 'nonce-…' 'strict-dynamic' https://va.vercel-scripts.com`, with no inline grant. `connect-src` adds `https://va.vercel-scripts.com` and `https://vitals.vercel-insights.com` for Vercel Analytics, which the layout renders with `<Analytics />`.
  - It sets the policy on the request, which next-intl forwards, so Next stamps the nonce on its scripts. It also sets the policy on the response, with `Cache-Control: private, no-store`.
  - Styles keep `'unsafe-inline'` for React `style` props and sinapsi's shadow root.
  - Development adds `'unsafe-eval'` and the development websocket, and preview builds allow `https://vercel.live`.
  - Do not enable Partial Prerendering or cache components: their static shell cannot carry a nonce.
- **The matcher.** The proxy runs on every path except `/api`, `/_next` and the root files `favicon.ico`, `icon.svg`, `apple-icon.png` and `robots.txt`. A new root static file joins that exclusion list.
- **The config's headers.**
  - Every response: HSTS (`max-age=63072000; includeSubDomains; preload`), `nosniff`, `Referrer-Policy: same-origin`, `X-Frame-Options: DENY`, the Permissions-Policy, `Cross-Origin-Opener-Policy` and `Cross-Origin-Resource-Policy: same-origin`, and `X-Robots-Tag: noindex, nofollow`.
  - API responses add `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`.

  No response carries two policies.
- **Server reads.** They build their request from the incoming host alone, so a page reached from another site, or rendered over HTTPS, reads normally. Browser calls to the route handlers keep the cross-site rule.
- **Entity ids.** Every id that reaches a console-stubs path keeps to letters, digits, hyphen and underscore, 1 to 128 characters. A route id outside that alphabet answers 404 with no console-stubs call, and no id is encoded twice. The console-stubs bridge is `server-only`.
- **The login form** posts, so a submit before hydration never puts the password in the URL.
- **Out of search.**
  - Every page says `noindex`, and `robots.txt` disallows everything. Crawlers that obey robots.txt therefore never fetch a page or see its noindex. A linked URL could still be listed bare, but Deployment Protection is the real gate.
  - Titles come from the catalogs: the site name, the template, and each screen after its project, runtime or ticket.
  - The 404 and the global error are localized.

## Hosted preview

On Vercel's `langdrift-console` project the console reads the hosted preview of console-stubs (spec 0015). Setting `CONSOLE_STUBS_TOKEN` lets the console serve its own host, accept an `https` `CONSOLE_STUBS_URL` on any host, and send the token to console-stubs as `Authorization: Bearer`. The preview is read-only: a change answers `read_only_preview`, and without a live stream the screens keep polling. Without the token, the loopback contract above holds. It needs two variables on Vercel:

| Variable | Requirement |
| --- | --- |
| `CONSOLE_STUBS_TOKEN` | The secret shared with `langdrift-stubs`; switches the console to the hosted preview. |
| `CONSOLE_STUBS_URL` | The stub's production URL, over `https`, without credentials. |

The owner's steps on Vercel, which no agent takes, are in [the stubs' README](/apps/microservices/console-stubs/README.md#hosted-preview). They set the token on both projects and `CONSOLE_STUBS_URL` on the console, and keep the console behind Deployment Protection for production and preview.

## The judge

`/api/judge` is the platform's judge for the LangDrift CLI (spec 0045). `dft agent plan` posts each of its questions there: the sources, the episodes, the suggestion and the reserved part. The console asks the platform's model with LangDrift's own key, so no consumer brings a provider. The question and the answer follow the contract in `@langdrift/core` (`lib/core/contracts/go/judge/schema/`).

- **`GET`** is the connection check `dft agent connect` runs. It answers 204 with the deployment's secret and 401 without it.
- **`POST`** takes one question as JSON, with the secret as `Authorization: Bearer`.
  - 401 without the secret.
  - 400 for a question outside the contract, naming the field. 400 for a model the console does not offer, with the offered set.
  - 413 over a mebibyte, and 415 for anything but plain JSON.
  - 503 while the judge is not configured, the platform key is missing or the provider cannot be reached.
  - 200 with the contract's response: the answer text the console validated against the kind's schema, or a refusal with its reason. It carries the tokens used, a cost in cents with provenance `console` marked illustrative (the list price of the tokens, until the owner prices a plan), and the judge's version.

The console stores no question. It logs one line per judgment with the kind, the model, the tokens and the cost, never the prompt or the answer.

It needs three variables on Vercel, documented here by name:

| Variable | Requirement |
| --- | --- |
| `RUNTIME_SECRET` | The deployment's secret. `dft agent connect` records it, and the route accepts only it. One per deployment until the account issues one per project. |
| `ANTHROPIC_API_KEY` | The platform's provider key, read at each call. Without it no call starts. |
| `JUDGE_MODELS` | The models the console offers, comma-separated, each one with a list price in `model-prices.data.ts`, for example `claude-sonnet-5-5,claude-opus-5-5`. |

The owner's steps on Vercel, which no agent takes:
- Set `RUNTIME_SECRET`, `ANTHROPIC_API_KEY` and `JUDGE_MODELS` on `langdrift-console` for production and preview.
- Let the CLI reach `/api/judge` through Deployment Protection, which otherwise answers its requests before the route does.
