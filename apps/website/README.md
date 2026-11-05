# Website

The LangDrift site, started as today's langdrift.md and reshaped into one agentic FinOps platform page. It is the public, localized LangDrift website: the home page, the platform pages, pricing and privacy. Demo calls to action lead to `/book-a-demo/`, which shows a return-soon card while demo requests are paused. Its form stays commented out in the page until demos open; it validates a request and forwards it to the demo-requests service under `apps/microservices`, and shows the return-soon notice while that service is not reachable.

The website holds no email credentials and never sends mail itself: the demo-requests service records a request and its mailer sends the emails. HTTP 202 confirms persistence, not delivery.

## Development

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm turbo run dev --filter=website` | Builds the tokens stylesheet the website depends on, then starts Next.js on port 3000. |
| `pnpm turbo run test --filter=website` | Runs the unit tests (Vitest; jsdom where a file asks for it). |
| `pnpm turbo run typecheck --filter=website` | Generate Next.js route types, then check TypeScript. |
| `pnpm turbo run build --filter=website` | Builds the app and the packages it depends on. |
| `drift check localization` | Checks the translation catalogs after a catalog or loader change. |
| `drift env setup` | Creates the local `.env` file from `.env.development`. |
| `drift env validate` | Checks that `.env.development` has its local `.env`. |

`.env.development` holds the shared development defaults and is tracked, so a fresh checkout runs without setup. In development, `CRYPTO_SECRET` is generated on the fly if not set. Git ignores `.env`. Put real and local values in `apps/website/.env` (run `drift env setup` to create it), which Next.js loads. Never put a production secret in a tracked file or in a `NEXT_PUBLIC_` variable.

| Variable | Requirement |
| --- | --- |
| `CRYPTO_SECRET` | Server-only 32-byte AES key in canonical base64; the locale proxy needs it. In development, it is generated on the fly if unset. In production, provide a 32-byte key from `openssl rand -base64 32`. Every instance of a deploy uses the same key. |
| `DEMO_REQUESTS_URL` | Origin of the demo-requests service the edge forwards accepted demo requests to; unset, the form answers that requests are paused. |
| `DEMO_REQUESTS_TOKEN` | Server-only bearer token the demo-requests service expects, at least 32 characters. |

`next start` serves the website correctly as its script runs it. Do not bind it with `--hostname 127.0.0.1`: Next renames a loopback host `localhost`, mistakes next-intl's rewrite to `/pt-BR` or `/zh-Hant` for an external one, and every `pt-br` and `zh-hant` page loops on 308 (ticket 0121). Vercel is unaffected.

## Security headers

Every response carries one header set, from `headers()` in `next.config.ts`; the proxy sets none. The pages stay statically prerendered, so the policy carries no nonce.

| Header | Value |
| --- | --- |
| `Content-Security-Policy` | `default-src 'self'; script-src 'self' 'unsafe-inline' https://va.vercel-scripts.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https://img.shields.io; font-src 'self' data:; connect-src 'self' https://va.vercel-scripts.com https://vitals.vercel-insights.com; worker-src 'self' blob:; manifest-src 'self'; frame-src 'none'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'` |
| `Strict-Transport-Security` | `max-age=63072000; includeSubDomains; preload` |
| `X-Content-Type-Options` | `nosniff` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `X-Frame-Options` | `DENY` |
| `Permissions-Policy` | `camera=(), microphone=(), geolocation=(), payment=(), usb=()` |
| `Cross-Origin-Opener-Policy` | `same-origin` |

`X-Powered-By` is off.

- **Why the policy looks like this.**
  - Scripts allow `'unsafe-inline'` because Next inlines its page payload, and nonces would render every page on demand; the site renders no user-supplied HTML.
  - Styles allow it for React `style` props, Monaco's injected styles and the footer orb's shadow root.
  - `font-src data:` and `worker-src` cover the Microruntimes editor's font and workers.
  - There is no `upgrade-insecure-requests`: HSTS already forces HTTPS, and the directive breaks plain-HTTP loopback runs.
- **Environment variants.** `next build` fixes the policy into the route manifest.
  - Development adds `'unsafe-eval'` and the development websocket.
  - A preview build (`VERCEL_ENV=preview`) allows `https://vercel.live` for the Vercel Toolbar. Turborepo's build task declares `VERCEL_ENV`, so a preview and a production build never share a cache entry.
  - Production allows nothing more.
- **Book a demo** checks the request before reading its body:
  - 405 for any method but POST;
  - 403 `INVALID_ORIGIN` for `Sec-Fetch-Site: cross-site` or an `Origin` other than the addressed origin (the request URL's scheme with the `Host` header's host);
  - 415 for any `Content-Encoding` or a type other than JSON;
  - 413 for a declared `Content-Length` over 4 KiB.

  A same-origin request keeps 400 `VALIDATION_ERROR` and 503 `UNAVAILABLE`.

## Search

`src/app/metadata/` is the search source, and `seo-contract.unit.ts` pins it.

- **Metadata.** Every indexable page in every locale has:
  - a unique title and description;
  - a self canonical on `https://langdrift.md` with a trailing slash;
  - hreflang for the five locales, with its own English URL as x-default;
  - Open Graph and Twitter cards and an index directive.
- **Sitemap and hreflang.** The sitemap lists exactly those canonicals with the same alternates. The HTML and the sitemap are the only hreflang sources: next-intl's `Link` header is off.
- **Structured data.** Every indexable page renders one JSON-LD graph with a `WebPage` node, built from the same JSON. The home page adds the `Organization` and `WebSite` nodes.
- **Social images.** The images live in the `[locale]` segment, under the layout's `metadataBase`. Next serves them once at `/-/opengraph-image.png` and `/-/twitter-image.png`, absolute on langdrift.md, outside the proxy.
- **robots.txt** keeps search and AI answer crawlers in, except on `/api/`. It disallows the whole site for the training crawlers GPTBot, ClaudeBot, Google-Extended, CCBot, Applebot-Extended and Bytespider.
- **`/.well-known/security.txt`** (RFC 9116) names `security@langdrift.md` and expires on 2027-10-29. The contract test fails once it expires, so renew it before then.

## Locales

English, Brazilian Portuguese, Traditional Chinese, Japanese and Vietnamese are served at `/en/`, `/pt-br/`, `/zh-hant/`, `/ja/` and `/vi/`.

- A language change keeps the current route, query and fragment.
- A Simplified Chinese request gets English with a visible availability notice.
- A localized missing path, dotted or not, returns 404 with recovery copy.

The proxy stores the language preference for 180 days in an encrypted, httpOnly `langdrift-website-locale` cookie:

1. A URL with a locale prefix keeps that language.
2. A URL without one uses a valid cookie.
3. Without a valid cookie, it uses `Accept-Language`.
4. Otherwise, English.

Copy lives in JSON under `src/app/template/i18n/messages/`; localized metadata lives under `src/app/metadata/`. Some catalog families are read directly by server components rather than through the next-intl provider. [The localization rule](/.claude/rules/08-localization.rule.md) and [the SEO rule](/.claude/rules/09-seo.rule.md) set the requirements.

## Privacy

The localized `/privacy/` page says what the site receives: the connection information every request carries, an email you choose to send, and the theme preference kept in local storage. It changes, with a new date, whenever the site starts collecting something else.

## After deploying

These are the owner's steps; no code can verify them.

1. **Production domain.** Confirm langdrift.md is the production domain, that `og:image` on a production page starts with `https://langdrift.md/-/`, and that the social card previews in a link debugger.
2. **The vercel.app alias.** Redirect the production `*.vercel.app` alias to langdrift.md, or keep relying on the canonicals; Vercel marks only previews `noindex`.
3. **Headers.** On production, check that `/` and `/en/` carry exactly one `Content-Security-Policy` and this HSTS value: Vercel's own default lacks `includeSubDomains` and `preload`.
4. **HSTS preload.** Once every langdrift.md subdomain serves HTTPS, submit the domain at [hstspreload.org](https://hstspreload.org/). Leaving the list takes months.
5. **Security contact.** Keep `security@langdrift.md` monitored, and renew `security.txt` before 2027-10-29.
6. **Deployment Protection.** Keep it on the console and the stubs.
