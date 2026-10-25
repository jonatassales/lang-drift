# Website (old)

The previous LangDrift website, kept while the platform site in `apps/website` is built. The public, localized LangDrift website: the product story, the Product Vision demos, pricing, privacy and the early-access flow. Demo calls to action lead to the return-soon section above the footer; requests are paused.

Early access posts to a private registry service outside this repository. Admission → Turnstile Siteverify → registration in the registry, which persists the contact and sends the confirmation. The website holds no email credentials and never sends mail itself. HTTP 202 confirms persistence, not delivery.

Sign-up is paused in code until feature flags have a proper home: the section shows the return-soon notice, the privacy page states the pause and `/api/early-access/` answers 503. The code that opens it stays commented where it lives: the form in `early-access/section.tsx`, the active notice in the privacy page and the real handler in `api/early-access/route.ts`.

## Development

Run from the repository root:

| Command | Effect |
| --- | --- |
| `pnpm turbo run dev --filter=website-old` | Builds the tokens stylesheet the website depends on, then starts Next.js on port 2000. |
| `pnpm turbo run test --filter=website-old` | Runs the unit tests (Vitest; jsdom where a file asks for it). |
| `pnpm turbo run typecheck --filter=website-old` | Type-checks the app. |
| `pnpm turbo run build --filter=website-old` | Builds the app and the packages it depends on. |
| `drift check localization` | Checks the translation catalogs after a catalog or loader change. |

`.env.development` holds the shared development defaults and is tracked, so it never holds a secret: it lists `CRYPTO_SECRET` without a value. Git ignores every other `.env*` file. Put local and secret values in `apps/website-old/.env.local`, which Next.js loads and which takes precedence over `.env.development`. Never put a secret in a `NEXT_PUBLIC_` variable.

| Variable | Requirement |
| --- | --- |
| `CRYPTO_SECRET` | Server-only 32-byte AES key in canonical base64; the locale proxy needs it, even locally. Generate with `openssl rand -base64 32`. Every instance of a deploy uses the same key. |
| `NEXT_PUBLIC_TURNSTILE_SITE_KEY` | The Invisible production widget; Cloudflare's official Invisible test key locally. |
| `TURNSTILE_SECRET_KEY` | Server-only widget secret; production rejects test keys. |
| `TURNSTILE_EXPECTED_HOSTNAME` | Exact approved hostname; `localhost` for local fixtures. |
| `TURNSTILE_VERIFY_URL` | In production, exactly `https://challenges.cloudflare.com/turnstile/v0/siteverify`. |
| `EARLY_ACCESS_PUBLIC_ORIGIN` | Exact canonical browser origin; HTTPS in production. |
| `EARLY_ACCESS_REGISTRY_URL` | Private authenticated registry origin; HTTPS in production. |
| `EARLY_ACCESS_REGISTRY_KEY` | Server-only website credential for the registry, at least 32 characters. |
| `EARLY_ACCESS_MODE` | `test` selects the test fixtures outside development; anything else behaves as production. |

Production fails closed until a trusted ingress contract for the host exists. Registry admission is never exposed to browsers, and client-supplied IP headers never count as proof.

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

## Early access

The voice interaction is a scripted preview with a prepared answer. Its description and accessible state feedback say it uses no microphone and generates no live answer. The early-access flow below stays paused until it is separately authorized.

**Submission:**

1. Focus prepares Invisible Turnstile lazily.
2. A valid submission takes a local lock, runs one challenge and sends JSON to `/api/early-access/`.
3. A spinner, live feedback, `aria-busy`, disabled controls and reduced-motion behavior cover verification, submission and navigation.
4. An error keeps the address, and a retry runs a fresh challenge.
5. Unmounting cancels client work. A cancellation after the server has persisted the request does not undo it; repeated requests keep the original intent.

**The route:**

- It checks the exact origin and host, a body of at most 4 KiB, strict field types and properties, and address limits before registry admission.
- A rejected admission never calls Siteverify.
- Siteverify has a 2.5 s timeout and strict success, hostname and action checks. Only verified tokens reach registration, and tokens never enter the registry, logs, URLs or browser storage.
- A registry 429 returns `Retry-After`; a 503 means unavailable.
- No result reveals whether an address already exists.

## Privacy and release

The localized `/privacy/` page and the footer links disclose Cloudflare Invisible Turnstile and link its [Privacy Addendum](https://www.cloudflare.com/en-gb/turnstile-privacy-policy/). The notice describes private contact persistence and the queued confirmation.

Before sign-up opens:

- Retention, deletion, operator and legal details need owner review. No global compliance claim is made.
- Production needs approved widget hostnames, pre-clearance off, TLS and private connectivity to the registry, proof of the WAF and proxy chain, and a controlled real-widget smoke test on the approved hostname.
- Local fixtures prove none of these.
- A future CSP follows [Cloudflare's guidance](https://developers.cloudflare.com/turnstile/reference/content-security-policy/).
